package service

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/date"
	"github.com/galaxy-book/common/core/util/encrypt"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"time"
	"upper.io/db.v3"
)

var getProcessError = "proxies.GetProcessStatusId: %v\n"

var log = *logger.GetDefaultLogger()

func HomeIssues(orgId, currentUserId int64, page int, size int, input *vo.HomeIssueInfoReq) (*vo.HomeIssueInfoResp, errs.SystemErrorInfo) {
	logger.GetDefaultLogger().Infof("当前登录用户 %d 组织 %d", currentUserId, orgId)

	issueCond, err1 := IssueCondAssembly(orgId, currentUserId, input)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueCondAssemblyError, err1)
	}

	logger.GetDefaultLogger().Infof("首页任务列表查询条件 %v", issueCond)

	var orderBy interface{} = consts.TcCreateTime + " desc"
	if input != nil && input.OrderType != nil {
		orderBy = domain.IssueCondOrderBy(orgId, *input.OrderType)
	}

	issueBos, total, err := domain.SelectList(issueCond, nil, page, size, orderBy)
	if err != nil {
		logger.GetDefaultLogger().Error(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	var actualTotal = total
	if input.EnableParentIssues != nil && *input.EnableParentIssues == 1 {
		parentIds := getIssueParentIds(*issueBos)
		if len(parentIds) > 0 {
			union := db.Or(issueCond, db.Cond{
				consts.TcId: db.In(parentIds),
			})
			issueBos, actualTotal, err = domain.SelectList(db.Cond{}, union, page, size, orderBy)
			if err != nil {
				logger.GetDefaultLogger().Error(strs.ObjectToString(err))
				return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
			}
		}
	}

	logger.GetDefaultLogger().Infof("首页任务列表命中总条数 %d", total)
	logger.GetDefaultLogger().Info(strs.ObjectToString(issueBos))

	homeIssueBos, err3 := domain.ConvertIssueBosToHomeIssueInfos(orgId, *issueBos)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}

	homeIssueVos := &[]*vo.HomeIssueInfo{}

	err2 := copyer.Copy(homeIssueBos, homeIssueVos)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err2)
	}

	resp := &vo.HomeIssueInfoResp{
		Total:       total,
		ActualTotal: actualTotal,
		List:        *homeIssueVos,
	}
	return resp, nil
}

func IssueReport(orgId, currentUserId int64, reportType int64) (*vo.IssueReportResp, errs.SystemErrorInfo) {

	logger.GetDefaultLogger().Infof("当前登录用户 %d 组织 %d", currentUserId, orgId)

	issueCond := db.Cond{}
	issueCond[consts.TcIsDelete] = consts.AppIsNoDelete
	issueCond[consts.TcOrgId] = orgId

	//获取我参与的和我负责的任务
	issueRelations, issueErr := domain.GetRelatedIssues(currentUserId, orgId)
	if issueErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, issueErr)
	}
	if len(issueRelations) > 0 {
		issueRelationIds := make([]int64, len(issueRelations))
		for i, entity := range issueRelations {
			issueRelationIds[i] = entity.IssueId
		}
		issueCond[consts.TcId] = db.In(issueRelationIds)
	} else {
		issueCond[consts.TcId] = -1
	}

	//时间范围
	var startTime, endTime string
	now := time.Now()
	switch reportType {
	case consts.DailyReport:
		startTime = date.Format(date.GetZeroTime(now))
		endTime = date.Format((date.GetZeroTime(now)).AddDate(0, 0, 1).Add(-1 * time.Second))
	case consts.WeeklyReport:
		startTime = date.Format(date.GetWeekStart(now))
		endTime = date.Format((date.GetWeekStart(now)).AddDate(0, 0, 7).Add(-1 * time.Second))
	case consts.MonthlyReport:
		startTime = date.Format(date.GetMonthStart(now))
		endTime = date.Format((date.GetMonthStart(now)).AddDate(0, 1, 0).Add(-1 * time.Second))
	}

	//进行中
	processingIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeProcessing)
	if err != nil {
		log.Errorf(getProcessError, err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}
	//已完成
	finishedId, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
	if err != nil {
		log.Errorf(getProcessError, err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	var union *db.Union = nil
	//获取进行中和在当前时间段内完成的任务
	union = db.Or(db.Cond{
		consts.TcStatus: db.In(processingIds),
	}).Or(
		db.And(
			db.Cond{consts.TcStatus: db.In(finishedId)},
			db.Cond{consts.TcEndTime: db.Gte(startTime)},
		),
	)

	logger.GetDefaultLogger().Infof("任务分享列表查询条件 %v", issueCond)

	orderBy := consts.TcPlanStartTime + " desc"

	issueBos, total, err := domain.SelectList(issueCond, union, -1, -1, orderBy)
	if err != nil {
		logger.GetDefaultLogger().Error(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}
	homeIssueBos, err := domain.ConvertIssueBosToHomeIssueInfos(orgId, *issueBos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	resp, insertErr := domain.InsertIssueReport(orgId, currentUserId, total, startTime, endTime, homeIssueBos)
	if insertErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, insertErr)
	}

	result := &vo.IssueReportResp{}

	copyErr := copyer.Copy(resp, result)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return result, nil
}

func IssueReportDetail(shareID string) (*vo.IssueReportResp, errs.SystemErrorInfo) {
	id, err := encrypt.AesDecrypt(shareID)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.DecryptError, err)
	}

	shareInfo, err := domain.GetIssueReport(id)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	homeIssue := &vo.IssueReportResp{}
	err = json.FromJson(*shareInfo.Content, homeIssue)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
	}

	return homeIssue, nil
}

func IssueCondBaseAssembly(issueCond db.Cond, input *vo.HomeIssueInfoReq) {
	if input.IterationID != nil {
		issueCond[consts.TcIterationId] = input.IterationID
	}
	if input.ProjectID != nil {
		issueCond[consts.TcProjectId] = input.ProjectID
	}
	if input.PriorityID != nil {
		issueCond[consts.TcPriorityId] = input.PriorityID
	}
	if input.PlanType != nil {
		planType := *input.PlanType
		dealIssueCondTcIterationId(planType, issueCond)
	}
	if input.ProjectObjectTypeID != nil {
		issueCond[consts.TcProjectObjectTypeId] = input.ProjectObjectTypeID
	}
	if input.ProjectObjectTypeIds != nil && len(input.ProjectObjectTypeIds) > 0 {
		issueCond[consts.TcProjectObjectTypeId] = db.In(input.ProjectObjectTypeIds)
	}
	if input.SearchCond != nil {
		issueCond[consts.TcTitle] = db.Like("%" + *input.SearchCond + "%")
	}
	if input.ProcessStatusID != nil {
		issueCond[consts.TcProcessStatusId] = input.ProcessStatusID
	}
	//这个是之前本周和全部的逻辑，现在去掉
	//if input.TimeScope != nil {
	//	union = db.Or(db.Cond{"create_time": db.Gte(date.FormatTime(*input.TimeScope))}).
	//		Or(db.Cond{consts.TcPlanStartTime: db.Gte(date.FormatTime(*input.TimeScope))}).
	//		Or(db.Cond{consts.TcPlanEndTime: db.Gte(date.FormatTime(*input.TimeScope))})
	//}
	if input.Type != nil {
		dealIssueCondTcParentId(input, issueCond)
	}
	if input.StartTime != nil && input.EndTime != nil {
		issueCond[consts.TcPlanEndTime] = db.Between(date.FormatTime(*input.StartTime), date.FormatTime(*input.EndTime))
	}
	if input.StartTime != nil && input.EndTime == nil {
		issueCond[consts.TcPlanEndTime] = db.Gte(date.FormatTime(*input.StartTime))
	}
	if input.StartTime == nil && input.EndTime != nil {
		issueCond[consts.TcPlanEndTime] = db.Lte(date.FormatTime(*input.EndTime))
	}
	if input.OwnerIds != nil {
		issueCond[consts.TcOwner] = db.In(input.OwnerIds)
	}
	if input.CreatorIds != nil {
		issueCond[consts.TcCreator] = db.In(input.CreatorIds)
	}
	if input.ParentID != nil {
		issueCond[consts.TcParentId] = input.ParentID
	}
	if input.PeriodStartTime != nil && input.PeriodEndTime != nil{
		//issueCond[consts.TcPlanStartTime] = db.Gte(date.FormatTime(*input.PeriodStartTime))
		//issueCond[consts.TcPlanEndTime] = db.Lte(date.FormatTime(*input.PeriodEndTime))
		//issueCond[consts.TcPlanEndTime + " "] = db.Gte(date.FormatTime(*input.PeriodStartTime))
		issueCond[consts.TcPlanStartTime] = db.Between(consts.BlankElasticityTime, date.FormatTime(*input.PeriodEndTime))
		issueCond[consts.TcPlanEndTime] = db.Gte(date.FormatTime(*input.PeriodStartTime))
	}
}

func dealIssueCondTcParentId(input *vo.HomeIssueInfoReq, issueCond db.Cond) {
	if *input.Type == 1 {
		issueCond[consts.TcParentId] = 0
	} else if *input.Type == 2 {
		issueCond[consts.TcParentId] = db.Gt(0)
	}
}

func dealIssueCondTcIterationId(planType int, issueCond db.Cond) {
	if planType == 1 {
		issueCond[consts.TcIterationId] = db.NotEq(0)
	} else if planType == 2 {
		issueCond[consts.TcIterationId] = db.Eq(0)
	}
}

func IssueCondAssembly(orgId, currentUserId int64, input *vo.HomeIssueInfoReq) (db.Cond, errs.SystemErrorInfo) {
	issueCond := db.Cond{}
	issueCond[consts.TcIsDelete] = consts.AppIsNoDelete
	issueCond[consts.TcOrgId] = orgId
	if input == nil {
		input = &vo.HomeIssueInfoReq{}
	}

	//封装基础条件
	IssueCondBaseAssembly(issueCond, input)

	if input.RelatedType != nil {
		domain.IssueCondRelatedTypeAssembly(issueCond, *input.RelatedType, currentUserId)
	} else {
		//拿到当前用户的管理员flag
		adminFlag, err := rolefacade.GetUserAdminFlagRelaxed(orgId, currentUserId)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		//会对私有项目做过滤，私有项目除外
		domain.IssueCondNoRelatedTypeAssembly(issueCond, currentUserId, orgId, adminFlag.IsAdmin)
	}

	if input.ProjectID != nil {
		issueCond[consts.TcProjectId+" "] = *input.ProjectID
	}

	if input.IsFiling == nil || *input.IsFiling == 0 || *input.IsFiling > 3 {
		//默认查询未归档的项目
		defaultFiling := consts.ProjectIsNotFiling
		input.IsFiling = &defaultFiling
	}
	domain.IssueCondFiling(issueCond, orgId, *input.IsFiling)

	if len(input.IssueTagID) != 0 {
		domain.IssueCondTagId(issueCond, orgId, input.IssueTagID)
	}

	if input.ResourceID != nil{
		domain.IssueCondResourceId(issueCond, orgId, *input.ResourceID)
	}

	if input.Status != nil {
		err1 := domain.IssueCondStatusAssembly(issueCond, orgId, *input.Status)
		if err1 != nil {
			log.Error(err1)
			return nil, errs.BuildSystemErrorInfo(errs.IssueCondAssemblyError, err1)
		}
	}
	domain.IssueCondRelationMemberAssembly(issueCond, input)
	//组合筛选类型封装
	err2 := domain.IssueCondCombinedCondAssembly(issueCond, input, currentUserId, orgId)
	if err2 != nil {
		log.Error(err2)
		return nil, err2
	}
	//增量查询条件封装
	domain.IssueCondLastUpdateTimeCondAssembly(issueCond, input)
	return issueCond, nil
}

func getIssueParentIds(issues []bo.IssueBo) []int64 {
	//子任务map, 子任务id -> 父任务id
	childMap := map[int64]int64{}
	for _, issue := range issues {
		childMap[issue.Id] = issue.ParentId
	}
	missingParentIds := make([]int64, 0)
	for _, v := range childMap {
		if v <= 0 {
			continue
		}
		if _, ok := childMap[v]; !ok {
			missingParentIds = append(missingParentIds, v)
		}
	}
	return missingParentIds
}

func IssueStatusTypeStat(orgId, currentUserId int64, input *vo.IssueStatusTypeStatReq) (*vo.IssueStatusTypeStatResp, errs.SystemErrorInfo) {
	if input.ProjectID != nil && !domain.JudgeProjectIsExist(orgId, *input.ProjectID) {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}
	if input.IterationID != nil {
		iterationBo, err := domain.GetIterationBoByOrgId(*input.IterationID, orgId)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.IterationNotExist)
		}
		input.ProjectID = &iterationBo.ProjectId
	}

	issueStatusStatBos, err1 := domain.GetIssueStatusStat(bo.IssueStatusStatCondBo{
		OrgId:        orgId,
		ProjectId:    input.ProjectID,
		IterationId:  input.IterationID,
		RelationType: input.RelationType,
		UserId:       currentUserId,
	})
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}
	result := &vo.IssueStatusTypeStatResp{}

	for _, issueStatusStatBo := range issueStatusStatBos {
		result.NotStartTotal += int64(issueStatusStatBo.IssueWaitCount)
		result.ProcessingTotal += int64(issueStatusStatBo.IssueRunningCount)
		result.CompletedTotal += int64(issueStatusStatBo.IssueEndCount)
		result.CompletedTodayTotal += int64(issueStatusStatBo.IssueEndTodayCount)
		result.Total += int64(issueStatusStatBo.IssueCount)
		result.OverdueCompletedTotal += int64(issueStatusStatBo.IssueOverdueEndCount)
		result.OverdueTodayTotal += int64(issueStatusStatBo.IssueOverdueTodayCount)
		result.OverdueTotal += int64(issueStatusStatBo.IssueOverdueCount)
		result.OverdueTomorrowTotal += int64(issueStatusStatBo.IssueOverdueTomorrowCount)
		result.TodayCount += int64(issueStatusStatBo.TodayCount)
		result.TodayCreateCount += int64(issueStatusStatBo.TodayCreateCount)
	}

	//即将到期 今天到期的主子任务数+明日逾期的主子任务数
	result.BeAboutToOverdueSum = result.OverdueTomorrowTotal + result.OverdueTodayTotal

	result.List = append(result.List, &vo.StatCommon{
		Name:  "已逾期",
		Count: result.OverdueTotal,
	})
	result.List = append(result.List, &vo.StatCommon{
		Name:  "进行中",
		Count: result.ProcessingTotal,
	})
	result.List = append(result.List, &vo.StatCommon{
		Name:  "未完成",
		Count: result.NotStartTotal + result.ProcessingTotal,
	})
	result.List = append(result.List, &vo.StatCommon{
		Name:  "已完成",
		Count: result.CompletedTotal,
	})

	return result, nil
}

func IssueStatusTypeStatDetail(orgId, currentUserId int64, input *vo.IssueStatusTypeStatReq) (*vo.IssueStatusTypeStatDetailResp, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//orgId := cacheUserInfo.OrgId
	projectId := input.ProjectID

	if projectId != nil && !domain.JudgeProjectIsExist(orgId, *projectId) {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}

	issueStatusStatBos, err1 := domain.GetIssueStatusStat(bo.IssueStatusStatCondBo{
		OrgId:        orgId,
		UserId:       currentUserId,
		ProjectId:    input.ProjectID,
		IterationId:  input.IterationID,
		RelationType: input.RelationType,
	})
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	result := &vo.IssueStatusTypeStatDetailResp{
		NotStart:   []*vo.IssueStatByObjectType{},
		Processing: []*vo.IssueStatByObjectType{},
		Completed:  []*vo.IssueStatByObjectType{},
	}

	for _, issueStatusStatBo := range issueStatusStatBos {
		if issueStatusStatBo.IssueWaitCount > 0 {
			result.NotStart = append(result.NotStart, &vo.IssueStatByObjectType{
				ProjectObjectTypeID:   &issueStatusStatBo.ProjectTypeId,
				ProjectObjectTypeName: &issueStatusStatBo.ProjectTypeName,
				Total:                 int64(issueStatusStatBo.IssueWaitCount),
			})
		}
		if issueStatusStatBo.IssueRunningCount > 0 {
			result.Processing = append(result.Processing, &vo.IssueStatByObjectType{
				ProjectObjectTypeID:   &issueStatusStatBo.ProjectTypeId,
				ProjectObjectTypeName: &issueStatusStatBo.ProjectTypeName,
				Total:                 int64(issueStatusStatBo.IssueRunningCount),
			})
		}
		if issueStatusStatBo.IssueEndCount > 0 {
			result.Completed = append(result.Completed, &vo.IssueStatByObjectType{
				ProjectObjectTypeID:   &issueStatusStatBo.ProjectTypeId,
				ProjectObjectTypeName: &issueStatusStatBo.ProjectTypeName,
				Total:                 int64(issueStatusStatBo.IssueEndCount),
			})
		}
	}

	return result, nil
}

func GetSimpleIssueInfoBatch(orgId int64, ids []int64) (*[]vo.Issue, errs.SystemErrorInfo) {
	list, _, err := domain.SelectList(db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       db.In(ids),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, nil, 0, 0, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	issueVo := &[]vo.Issue{}
	copyErr := copyer.Copy(list, issueVo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return issueVo, nil
}

func GetIssueRemindInfoList(reqVo projectvo.GetIssueRemindInfoListReqVo) (*projectvo.GetIssueRemindInfoListRespData, errs.SystemErrorInfo) {
	if reqVo.Page < 0 {
		return nil, errs.BuildSystemErrorInfo(errs.PageInvalidError)
	}
	if reqVo.Size < 0 || reqVo.Size > 100 {
		return nil, errs.BuildSystemErrorInfo(errs.PageSizeInvalidError)
	}

	selectIssueIdsCondBo := bo.SelectIssueIdsCondBo{}

	input := reqVo.Input
	//计划结束时间条件
	selectIssueIdsCondBo.BeforePlanEndTime = input.BeforePlanEndTime
	selectIssueIdsCondBo.AfterPlanEndTime = input.AfterPlanEndTime

	issueRemindInfos, total, err := domain.SelectIssueRemindInfoList(selectIssueIdsCondBo, reqVo.Page, reqVo.Size)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &projectvo.GetIssueRemindInfoListRespData{
		Total: total,
		List:  issueRemindInfos,
	}, nil
}
