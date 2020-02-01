package service

import (
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/core/util/excel"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/bo/mqbo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"github.com/tealeg/xlsx/v2"
	"os"
	"strconv"
	"strings"
	"time"
	"upper.io/db.v3"
)

const maxRawSupport = 300

func ImportIssues(orgId, currentUserId int64, input vo.ImportIssuesReq) (int64, errs.SystemErrorInfo) {
	//获取url数据
	data, importErr := excel.GenerateCSVFromXLSXFile(input.URL, input.URLType, 0, 2, []int64{6, 7})
	//importData, importErr := excel.GenerateCSVFromXLSXFile("F:\\polaris-backend-clone\\service\\platform\\projectsvc\\service\\issue.xlsx", input.URL, 0, true)
	if importErr != nil {
		log.Error(importErr)
		return 0, errs.InvalidImportFile
	}
	//手动处理空行
	importData := [][]string{}
	for _, v := range data {
		tempLen := len(v)
		for _, val := range v {
			if val == "" {
				tempLen--
			}
		}
		if tempLen > 0 {
			importData = append(importData, v)
		}
	}
	if len(importData) == 0 {
		return 0, errs.ImportDataEmpty
	}
	if len(importData) > maxRawSupport {
		return 0, errs.TooLargeImportData
	}

	projectInfo, projectErr := domain.GetProject(orgId, input.ProjectID)
	if projectErr != nil {
		log.Error(projectErr)
		return 0, errs.ProjectNotExist
	}
	//获取优先级
	issueType := 2
	priorityList, err := domain.GetPriorityList(orgId)
	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}
	newPriority := map[string]int64{}
	for _, v := range *priorityList {
		if v.Type == issueType {
			newPriority[v.Name] = v.Id
		}
	}
	//获取任务类型
	typeList, err := domain.GetProjectObjectTypeList(orgId, input.ProjectID)
	if err != nil {
		log.Error(err)
		return 0, errs.SystemError
	}

	newIssueType := map[string]int64{}
	for _, v := range *typeList {
		if v.ObjectType == 2 {
			newIssueType[v.Name] = v.Id
		}
	}
	//人员信息
	userInfo := orgfacade.GetUserInfoListByOrg(orgvo.GetUserInfoListReqVo{
		OrgId: orgId,
	})

	if userInfo.Failure() {
		log.Error(userInfo.Error())
		return 0, errs.SystemError
	}
	newUserInfo := map[string]int64{}
	for _, v := range userInfo.SimpleUserInfo {
		newUserInfo[v.Name] = v.Id
	}
	count, err := handleImportData(orgId, currentUserId, importData, newPriority, newIssueType, newUserInfo, input.ProjectID)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = projectInfo.Name
		domain.PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeCreateIssueBatch,
			OrgId:      orgId,
			ProjectId:  input.ProjectID,
			OperatorId: currentUserId,
			Ext:        ext,
		})
	})
	return count, nil
}

func handleImportData(orgId, userId int64, importData [][]string, newPriority, newIssueType, newUserInfo map[string]int64, projectId int64) (int64, errs.SystemErrorInfo) {
	allIssue := []vo.CreateIssueReq{}
	i := int64(0)
	var errStr []string
	//errIds := map[int64][]int64{}
	for k, v := range importData {
		log.Info(v)
		if len(v) < 10 {
			arrTemp := make([]string, 10-len(v))
			v = append(v, arrTemp...)
		}
		errIdList := []int64{}
		//是否子任务v[9]
		if i == 0 && v[9] == "子任务" {
			return 0, errs.ChildIssueForFirst
		}
		if v[9] == "子任务" {
			child := buildChildIssues(i, v, newUserInfo, newPriority, newIssueType, &errIdList, userId)
			if child != nil {
				allIssue[i-1].Children = append(allIssue[i-1].Children, child)
			}
		} else {
			issueData := vo.CreateIssueReq{}
			issueData.ProjectID = projectId

			//标题
			issueData.Title = strings.TrimSpace(v[0])
			if issueData.Title == "" || strs.Len(issueData.Title) > 50 {
				errIdList = append(errIdList, 0+1)
			}

			assemblyPriorityAndIssueType(newPriority, newIssueType, &issueData, v, &errIdList)

			//人员处理
			assemblyImportIssueUser(newUserInfo, &issueData, v, &errIdList, userId)

			//时间处理
			assemblyImportIssueTime(v, &issueData, &errIdList)

			//备注
			remark := v[8]
			issueData.Remark = &remark
			allIssue = append(allIssue, issueData)
			i++
		}
		if len(errIdList) > 0 {
			errStr = append(errStr, fmt.Sprintf(" 第%d行,第%s列文本格式错误", k+3, strings.Replace(strings.Trim(fmt.Sprint(errIdList), "[]"), " ", "、", -1)))
			//errIds[int64(k)+3] = errIdList
		}
	}
	if len(errStr) > 0 {
		log.Error(errStr)
		//return errs.BuildSystemErrorInfoWithMessage(errs.FileParseFail, errStr)
		return 0, errs.BuildSystemErrorInfoWithMessage(errs.FileParseFail, json.ToJsonIgnoreError(errStr))
	}
	if err := dealCreateIssue(allIssue, orgId, userId); err != nil {
		log.Error(err)
		return 0, errs.SystemError
	}

	return int64(len(importData)), nil
}

func assemblyImportIssueTime(v []string, issueData *vo.CreateIssueReq, errIdList *[]int64) error {
	//计划开始时间0
	//planStartTime := convertToFormatDay(strings.Split(v[6], ".")[0])
	if v[6] == "" {
		v[6] = consts.BlankTime
	}
	planStartTime, err := time.Parse(consts.AppTimeFormat, v[6])
	if err != nil {
		*errIdList = append(*errIdList, 6+1)
		log.Error(err)
	}

	if planStartTime.After(consts.BlankTimeObject) {
		start := types.Time(planStartTime)
		issueData.PlanStartTime = &start
	}
	//计划结束时间
	//planEndTime := convertToFormatDay(strings.Split(v[7], ".")[0])
	if v[7] == "" {
		v[7] = consts.BlankTime
	}
	planEndTime, err := time.Parse(consts.AppTimeFormat, v[7])
	if err != nil {
		*errIdList = append(*errIdList, 7+1)
		log.Error(err)
	}
	if planEndTime.After(consts.BlankTimeObject) {
		end := types.Time(planEndTime)
		issueData.PlanEndTime = &end
	}

	if planEndTime.After(consts.BlankTimeObject) && planEndTime.Before(planStartTime) {
		*errIdList = append(*errIdList, 7+1)
	}

	return nil
}

func assemblyPriorityAndIssueType(newPriority, newIssueType map[string]int64, issueData *vo.CreateIssueReq, v []string, errIdList *[]int64) {
	//优先级
	if newPriority[v[1]] != 0 {
		issueData.PriorityID = newPriority[v[1]]
	} else {
		*errIdList = append(*errIdList, 1+1)
	}
	//问题类型
	if newIssueType[v[2]] != 0 {
		typeId := newIssueType[v[2]]
		issueData.TypeID = &typeId
	} else {
		*errIdList = append(*errIdList, 2+1)
	}
}

func assemblyImportIssueUser(newUserInfo map[string]int64, issueData *vo.CreateIssueReq, v []string, errIdList *[]int64, userId int64) {

	//负责人
	if newUserInfo[v[3]] != 0 {
		issueData.OwnerID = newUserInfo[v[3]]
	} else {
		issueData.OwnerID = userId
	}
	//参与人
	participant := strings.Split(v[4], "|")
	for _, name := range participant {
		if newUserInfo[name] != 0 {
			issueData.ParticipantIds = append(issueData.ParticipantIds, newUserInfo[name])
		}
	}
	//关注人
	followers := strings.Split(v[5], "|")
	for _, name := range followers {
		if newUserInfo[name] != 0 {
			issueData.FollowerIds = append(issueData.FollowerIds, newUserInfo[name])
		}
	}
}

func dealCreateIssue(allIssue []vo.CreateIssueReq, orgId, userId int64) error {
	issueListSize := len(allIssue)

	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssue, issueListSize)
	if err != nil {
		log.Error(err)
		return err
	}

	for i, v := range allIssue {
		//defer wg.Add(-1)
		reqVo := projectvo.CreateIssueReqVo{
			CreateIssue: v,
			UserId:      userId,
			OrgId:       orgId,
		}
		createIssueBo := mqbo.PushCreateIssueBo{
			IssueId:          ids.Ids[i].Id,
			CreateIssueReqVo: reqVo,
		}
		domain.PushCreateIssue(createIssueBo)
	}

	return nil
}

func buildChildIssues(i int64, value []string, newUserInfo, newPriority, newIssueType map[string]int64, errIdList *[]int64, userId int64) *vo.IssueChildren {
	child := vo.IssueChildren{}
	child.Title = strings.TrimSpace(value[0])
	if child.Title == "" || strs.Len(child.Title) > 50 {
		*errIdList = append(*errIdList, 0+1)
	}
	if newPriority[value[1]] != 0 {
		child.PriorityID = newPriority[value[1]]
	} else {
		*errIdList = append(*errIdList, 1+1)
	}
	//问题类型
	if newIssueType[value[2]] != 0 {
		typeId := newIssueType[value[2]]
		child.TypeID = &typeId
	} else {
		*errIdList = append(*errIdList, 2+1)
	}
	//负责人
	if newUserInfo[value[3]] != 0 {
		child.OwnerID = newUserInfo[value[3]]
	} else {
		child.OwnerID = userId
	}
	//计划开始时间0
	if value[6] == "" {
		value[6] = consts.BlankTime
	}
	planStartTime, err := time.Parse(consts.AppTimeFormat, value[6])
	if err != nil {
		*errIdList = append(*errIdList, 6+1)
		log.Error(err)
	}

	if planStartTime.After(consts.BlankTimeObject) {
		start := types.Time(planStartTime)
		child.PlanStartTime = &start
	}

	//计划结束时间
	if value[7] == "" {
		value[7] = consts.BlankTime
	}
	planEndTime, err := time.Parse(consts.AppTimeFormat, value[7])
	if err != nil {
		*errIdList = append(*errIdList, 7+1)
		log.Error(err)
	}

	if planEndTime.After(consts.BlankTimeObject) {
		end := types.Time(planEndTime)
		child.PlanEndTime = &end
	}

	if planEndTime.After(consts.BlankTimeObject) && planEndTime.Before(planStartTime) {
		*errIdList = append(*errIdList, 7+1)
	}
	return &child
}

// excel日期字段格式化 yyyy-mm-dd
func convertToFormatDay(excelDaysString string) time.Time {
	// 2006-01-02 距离 1900-01-01的天数
	baseDiffDay := 38719 //在网上工具计算的天数需要加2天，什么原因没弄清楚
	curDiffDay := excelDaysString
	b, _ := strconv.Atoi(curDiffDay)
	// 获取excel的日期距离2006-01-02的天数
	realDiffDay := b - baseDiffDay
	//fmt.Println("realDiffDay:",realDiffDay)
	// 距离2006-01-02 秒数
	realDiffSecond := realDiffDay * 24 * 3600
	//fmt.Println("realDiffSecond:",realDiffSecond)
	// 2006-01-02 15:04:05距离1970-01-01 08:00:00的秒数 网上工具可查出
	baseOriginSecond := 1136185445
	resultTime := time.Unix(int64(baseOriginSecond+realDiffSecond), 0)
	return resultTime
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func DeleteProjectExcel(orgId int64, projectId int64) {
	relatePath := "/org_" + strconv.FormatInt(orgId, 10) + "/project_" + strconv.FormatInt(projectId, 10)
	excelDir := config.GetOSSConfig().RootPath + relatePath
	excelPath := excelDir + "/任务批量导入模板.xlsx"

	if !IsFileExist(excelPath) {
		return
	}
	err := os.Remove(excelPath)
	if err != nil {
		log.Error(err)
	}
}

func ExportIssueTemplate(orgId int64, projectId int64) (string, errs.SystemErrorInfo) {
	_, projectErr := domain.GetProject(orgId, projectId)
	if projectErr != nil {
		log.Error(projectErr)
		return "", errs.ProjectNotExist
	}
	relatePath := "/org_" + strconv.FormatInt(orgId, 10) + "/project_" + strconv.FormatInt(projectId, 10)
	excelDir := config.GetOSSConfig().RootPath + relatePath
	mkdirErr := os.MkdirAll(excelDir, 0777)
	if mkdirErr != nil {
		log.Error(mkdirErr)
		return "", errs.BuildSystemErrorInfo(errs.IssueDomainError, mkdirErr)
	}
	fileName := "任务批量导入模板.xlsx"
	excelPath := excelDir + "/" + fileName
	url := config.GetOSSConfig().LocalDomain + relatePath + "/" + fileName
	if IsFileExist(excelPath) {
		log.Info("模板已存在")
		return url, nil
	}

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return "", errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	row = sheet.AddRow()
	cell = row.AddCell()
	sheet.SetColWidth(1, 1, 100)

	var hStyle = xlsx.NewStyle()
	hStyle.Font.Bold = true
	hStyle.Alignment.WrapText = true
	cell.SetStyle(hStyle)
	cell.Value = "Tip：" +
		"\n1.*为必填项" +
		"\n2.任务类型若为子任务，当前子任务会跟随上一个父级任务" +
		"\n3.任务总数不超过300条" +
		"\n4.任务标题不能超过50字" +
		"\n5.多选字段如“参与人”字段：请用“||”符号隔开，如“张三||李四||王五”" +
		"\n6.日期字段如“开始时间”字段：填写格式为YYYY-MM-DD XX:XX，如“2020-04-10 12:00”" +
		"\n7.截止时间不能小于开始时间" +
		"\n8.负责人若不填写或组织内不存在该人员姓名，那么本条任务的负责人将默认为上传者"

	row = sheet.AddRow()

	cell = row.AddCell()
	cell.Value = "*任务标题"

	cell = row.AddCell()
	cell.Value = "*优先级"

	cell = row.AddCell()
	cell.Value = "*任务栏"

	cell = row.AddCell()
	cell.Value = "*负责人"

	cell = row.AddCell()
	cell.Value = "参与人"

	cell = row.AddCell()
	cell.Value = "关注人"

	cell = row.AddCell()
	cell.Value = "开始时间"
	//cell.SetFormat("yyyy-m-d h:mm")
	//cell.SetFormat("yyyy\"年\"m\"月\"d\"日\"")

	cell = row.AddCell()
	cell.Value = "截止时间"
	//cell.SetFormat("yyyy-m-d h:mm")
	//cell.SetFormat("yyyy\"年\"m\"月\"d\"日\"")

	cell = row.AddCell()
	cell.Value = "任务描述"

	cell = row.AddCell()
	cell.Value = "*父子类型"

	//优先级
	priority, err := PriorityList(orgId, 0, 0, db.Cond{
		consts.TcType: 2,
	})
	if err != nil {
		return "", errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}
	priorityList := []string{}
	for _, v := range priority.List {
		priorityList = append(priorityList, v.Name)
	}
	dd := xlsx.NewDataValidation(2, 1, maxRawSupport+1, 1, false)
	dd.SetDropList(priorityList)
	sheet.AddDataValidation(dd)

	//获取任务类型
	typeList, err := domain.GetProjectObjectTypeList(orgId, projectId)
	if err != nil {
		return "", errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}
	newIssueType := []string{}
	for _, v := range *typeList {
		newIssueType = append(newIssueType, v.Name)
	}
	dd1 := xlsx.NewDataValidation(2, 2, maxRawSupport+1, 2, false)
	dd1.SetDropList(newIssueType)
	sheet.AddDataValidation(dd1)

	dd2 := xlsx.NewDataValidation(2, 9, maxRawSupport+1, 9, false)
	dd2.SetDropList([]string{"父任务", "子任务"})
	sheet.AddDataValidation(dd2)

	err = file.Save(excelPath)
	if err != nil {
		return "", errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	return url, nil
}

func ExportData(orgId int64, projectId int64) (string, errs.SystemErrorInfo) {
	projectInfo, projectErr := domain.GetProject(orgId, projectId)
	if projectErr != nil {
		log.Error(projectErr)
		return "", errs.ProjectNotExist
	}

	relatePath := "/org_" + strconv.FormatInt(orgId, 10) + "/project_" + strconv.FormatInt(projectId, 10)
	excelDir := config.GetOSSConfig().RootPath + relatePath
	mkdirErr := os.MkdirAll(excelDir, 0777)
	if mkdirErr != nil {
		log.Error(mkdirErr)
		return "", errs.BuildSystemErrorInfo(errs.IssueDomainError, mkdirErr)
	}
	fileName := projectInfo.Name + "_任务总览.xlsx"
	excelPath := excelDir + "/" + fileName
	url := config.GetOSSConfig().LocalDomain + relatePath + "/" + fileName

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return "", errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	row = sheet.AddRow()

	cell = row.AddCell()
	cell.Value = "任务栏名称"

	cell = row.AddCell()
	cell.Value = "任务标题"

	cell = row.AddCell()
	cell.Value = "负责人"

	cell = row.AddCell()
	cell.Value = "任务类型"

	cell = row.AddCell()
	cell.Value = "任务优先级"

	cell = row.AddCell()
	cell.Value = "任务状态"

	cell = row.AddCell()
	cell.Value = "任务描述"

	cell = row.AddCell()
	cell.Value = "任务开始时间"
	//cell.SetFormat("yyyy-m-d h:mm")
	//cell.SetFormat("yyyy\"年\"m\"月\"d\"日\"")

	cell = row.AddCell()
	cell.Value = "任务结束时间"

	cell = row.AddCell()
	cell.Value = "任务创建时间"

	cell = row.AddCell()
	cell.Value = "任务创建者"

	//优先级
	priority, err := PriorityList(orgId, 0, 0, db.Cond{
		consts.TcType: 2,
	})
	if err != nil {
		return "", errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}
	priorityList := map[int64]string{}
	for _, v := range priority.List {
		priorityList[v.ID] = v.Name
	}

	//获取任务类型
	typeList, err := domain.GetProjectObjectTypeList(orgId, projectId)
	if err != nil {
		return "", errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}
	newIssueType := map[int64]string{}
	for _, v := range *typeList {
		newIssueType[v.Id] = v.Name
	}

	//人员信息
	//userInfo := orgfacade.GetUserInfoListByOrg(orgvo.GetUserInfoListReqVo{
	//	OrgId: orgId,
	//})
	//if userInfo.Failure() {
	//	log.Error(userInfo.Error())
	//	return "", userInfo.Error()
	//}
	//newUserInfo := map[int64]string{}
	//for _, v := range userInfo.SimpleUserInfo {
	//	newUserInfo[v.Id] = v.Name
	//}

	//任务状态
	statusList, statusErr := processfacade.GetProcessStatusListByCategoryRelaxed(orgId, consts.ProcessStatusCategoryIssue)
	if statusErr != nil {
		log.Error(statusErr)
		return "", statusErr
	}
	newStatusList := map[int64]string{}
	for _, statusBo := range statusList {
		newStatusList[statusBo.StatusId] = statusBo.Name
	}

	//获取父任务
	issuesParent, issueErr := domain.AllIssueForProject(orgId, projectId, true)
	if issueErr != nil {
		log.Error(issueErr)
		return "", issueErr
	}
	//获取子任务
	issuesChildren, issueErr := domain.AllIssueForProject(orgId, projectId, false)
	if issueErr != nil {
		log.Error(issueErr)
		return "", issueErr
	}

	userIds := []int64{}
	for _, infoBo := range issuesParent {
		userIds = append(userIds, infoBo.Owner, infoBo.Creator)
	}
	for _, infoBo := range issuesChildren {
		userIds = append(userIds, infoBo.Owner, infoBo.Creator)
	}
	userIds = slice.SliceUniqueInt64(userIds)

	//获取人员信息
	userInfo, infoErr := orgfacade.GetBaseUserInfoBatchRelaxed("", orgId, userIds)
	if infoErr != nil {
		log.Error(infoErr)
		return "", infoErr
	}
	newUserInfo := map[int64]string{}
	for _, v := range userInfo {
		newUserInfo[v.UserId] = v.Name
	}

	//拼装sheet数据
	for _, bos := range issuesParent {
		insertCell(sheet, bos, priorityList, newIssueType, newUserInfo, newStatusList, true)
		for _, infoBo := range issuesChildren {
			if infoBo.ParentId == bos.Id {
				insertCell(sheet, infoBo, priorityList, newIssueType, newUserInfo, newStatusList, false)
			}
		}
	}

	err = file.Save(excelPath)
	if err != nil {
		return "", errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	return url, nil
}

func insertCell(sheet *xlsx.Sheet, issueInfo bo.IssueAndDetailInfoBo, priorityList, typeList, userInfo, statusInfo map[int64]string, isParent bool) {
	var row *xlsx.Row
	var cell *xlsx.Cell

	row = sheet.AddRow()

	cell = row.AddCell()
	if typeName, ok := typeList[issueInfo.ProjectObjectTypeId]; ok {
		cell.Value = typeName
	}

	cell = row.AddCell()
	cell.Value = issueInfo.Title

	cell = row.AddCell()
	if userName, ok := userInfo[issueInfo.Owner]; ok {
		cell.Value = userName
	}

	cell = row.AddCell()
	if isParent {
		cell.Value = "父任务"
	} else {
		cell.Value = "子任务"
	}

	cell = row.AddCell()
	if priorityName, ok := priorityList[issueInfo.PriorityId]; ok {
		cell.Value = priorityName
	}

	cell = row.AddCell()
	if statusName, ok := statusInfo[issueInfo.Status]; ok {
		cell.Value = statusName
	}

	cell = row.AddCell()
	cell.Value = issueInfo.Remark

	cell = row.AddCell()
	if issueInfo.PlanStartTime.String() >= consts.BlankElasticityTime {
		cell.Value = issueInfo.PlanStartTime.String()
	}
	//cell.SetFormat("yyyy-m-d h:mm")
	//cell.SetFormat("yyyy\"年\"m\"月\"d\"日\"")

	cell = row.AddCell()
	if issueInfo.PlanEndTime.String() >= consts.BlankElasticityTime {
		cell.Value = issueInfo.PlanEndTime.String()
	}

	cell = row.AddCell()
	cell.Value = issueInfo.CreateTime.String()

	cell = row.AddCell()
	if userName, ok := userInfo[issueInfo.Creator]; ok {
		cell.Value = userName
	}
}
