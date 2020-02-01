package service

import (
	enjson "encoding/json"
	"errors"
	"fmt"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/cond"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/core/util/times"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/facade/resourcefacade"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"gopkg.in/fatih/set.v0"
	"strings"
	"sync"
	"time"
	"upper.io/db.v3"
)

func Projects(reqVo projectvo.ProjectsRepVo) (*vo.ProjectList, errs.SystemErrorInfo) {
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	page := reqVo.Page
	size := reqVo.Size
	params := reqVo.ProjectExtraBody.Params
	order := reqVo.ProjectExtraBody.Order
	input := reqVo.ProjectExtraBody.Input

	log.Infof(consts.UserLoginSentence, currentUserId, orgId)

	var joinParams db.Cond
	var joinErr errs.SystemErrorInfo
	if len(params) == 0 {
		joinParams, joinErr = GetProjectCondAssemblyByInput(input, currentUserId, orgId)
	} else {
		joinParams, joinErr = GetProjectCondAssemblyByParam(params, currentUserId, orgId)
	}
	joinParams[consts.TcOrgId] = orgId
	log.Infof("joinParams %v %v %v", params, len(params), joinParams)

	if joinErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError)
	}

	//获取项目列表
	var totalNumberOfEntries int64
	entities, totalNumberOfEntries, err := domain.GetProjectList(currentUserId, joinParams, order, size, page)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	resultList := &[]*vo.Project{}
	copyer.Copy(entities, resultList)

	if isFailure := getRedundancyInfo(resultList, orgId, reqVo.SourceChannel); isFailure {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError)
	}

	result := &vo.ProjectList{}
	result.Total = int64(totalNumberOfEntries)
	result.List = *resultList
	return result, nil
}

//获取冗余信息
func getRedundancyInfo(resultList *[]*vo.Project, orgId int64, sourceChannel string) (isFailure bool) {
	resourceIds := []int64{}
	creatorIds := []int64{}
	projectIds := []int64{}

	if len(*resultList) != 0 {
		for _, v := range *resultList {
			resourceIds = append(resourceIds, v.ResourceID)
			creatorIds = append(creatorIds, v.Creator)
			projectIds = append(projectIds, v.ID)
		}

		//ownerInfo, participantInfo, followerInfo, creatorInfo, err := domain.GetProjectMemberInfo(projectIds, orgId, creatorIds)
		//if err != nil {
		//	return resourcesRespVo.Failure()
		//}
		////获取任务信息
		//issueStat, _ := domain.GetIssueInfoByProject(projectIds, orgId)
		//
		////获取projectType localCache
		//projectTypeList, err := domain.GetProjectTypeList(orgId)
		//if err != nil {
		//	log.Error(err)
		//	return resourcesRespVo.Failure()
		//}
		//projectTypeLocalCache := maps.NewMap("Id", projectTypeList)
		////获取项目详情
		//projectDetails, detailErr := domain.GetProjectDetails(orgId, projectIds)
		//if detailErr != nil {
		//	log.Error(detailErr)
		//	return false
		//}
		//projectDetailById := maps.NewMap("ProjectId", projectDetails)

		type midType struct {
			creatorInfo           map[int64]bo.UserIDInfoBo
			ownerInfo             map[int64]bo.UserIDInfoBo
			participantInfo       map[int64][]bo.UserIDInfoBo
			followerInfo          map[int64][]bo.UserIDInfoBo
			resourceByPath        map[int64]bo.ResourceBo
			issueStat             map[int64]bo.IssueStatistic
			projectTypeLocalCache maps.LocalMap
			projectDetailById     maps.LocalMap
		}

		handlerFuncList := make([]func(midInfo *midType, wg *sync.WaitGroup), 0)
		midInfo := &midType{}

		handlerFuncList = append(handlerFuncList, func(midInfo *midType, wg *sync.WaitGroup) {
			defer wg.Add(-1)

			ownerInfo, participantInfo, followerInfo, creatorInfo, err := domain.GetProjectMemberInfo(projectIds, orgId, creatorIds, sourceChannel)
			if err != nil {
				log.Error(err)
				return
			}
			midInfo.creatorInfo = creatorInfo
			midInfo.ownerInfo = ownerInfo
			midInfo.followerInfo = followerInfo
			midInfo.participantInfo = participantInfo
		})

		handlerFuncList = append(handlerFuncList, func(midInfo *midType, wg *sync.WaitGroup) {
			defer wg.Add(-1)
			//资源列表
			resourcesRespVo := resourcefacade.GetResourceById(resourcevo.GetResourceByIdReqVo{GetResourceByIdReqBody: resourcevo.GetResourceByIdReqBody{ResourceIds: resourceIds}})
			if resourcesRespVo.Failure() {
				log.Error(resourcesRespVo.Message)
				return
			}
			resourceEntities := resourcesRespVo.ResourceBos
			resourceByPath := map[int64]bo.ResourceBo{}
			for _, v := range resourceEntities {
				resourceByPath[v.Id] = v
			}
			midInfo.resourceByPath = resourceByPath
		})

		handlerFuncList = append(handlerFuncList, func(midInfo *midType, wg *sync.WaitGroup) {
			defer wg.Add(-1)

			//获取任务信息
			issueStat, _ := domain.GetIssueInfoByProject(projectIds, orgId)
			midInfo.issueStat = issueStat
		})

		handlerFuncList = append(handlerFuncList, func(midInfo *midType, wg *sync.WaitGroup) {
			defer wg.Add(-1)

			//获取projectType localCache
			projectTypeList, err := domain.GetProjectTypeList(orgId)
			if err != nil {
				log.Error(err)
				return
			}
			projectTypeLocalCache := maps.NewMap("Id", projectTypeList)
			midInfo.projectTypeLocalCache = projectTypeLocalCache
		})

		handlerFuncList = append(handlerFuncList, func(midInfo *midType, wg *sync.WaitGroup) {
			defer wg.Add(-1)

			//获取项目详情
			projectDetails, detailErr := domain.GetProjectDetails(orgId, projectIds)
			if detailErr != nil {
				log.Error(detailErr)
				return
			}
			projectDetailById := maps.NewMap("ProjectId", projectDetails)
			midInfo.projectDetailById = projectDetailById
		})

		var wg sync.WaitGroup
		wg.Add(len(handlerFuncList))

		for _, handlerFunc := range handlerFuncList {
			tempHandlerFunc := handlerFunc
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Errorf("捕获到的错误：%s", r)
					}
				}()
				tempHandlerFunc(midInfo, &wg)
			}()
		}

		wg.Wait()

		dealResultList(resultList, midInfo.ownerInfo, midInfo.participantInfo, midInfo.followerInfo, midInfo.resourceByPath,
			midInfo.issueStat, midInfo.projectTypeLocalCache, midInfo.projectDetailById)
		addCreatorInfo(resultList, midInfo.creatorInfo)
	}

	return false
}

func GetProjectCondAssemblyByParam(params map[string]interface{}, currentUserId int64, orgId int64) (db.Cond, errs.SystemErrorInfo) {
	var relationType interface{}
	if _, ok := params["relation_type"]; ok {
		if val, ok := params["relation_type"].(map[string]interface{}); ok {
			if val["type"] != nil && val["value"] != nil {
				relationType = val["value"]
			}
		} else {
			relationType = params["relation_type"]
		}
		delete(params, "relation_type")
	}
	var relateType int64 = 0

	converRelateType(&relateType, relationType)

	condParam, err := cond.HandleParams(params)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ConditionHandleError, err)
	}
	switch relateType {
	case 0:
		//所有
	case 1:
		//我发起的
		condParam[consts.TcCreator] = currentUserId
	case 2:
		//我负责的
		condParam[consts.TcOwner] = db.Eq(currentUserId)
	case 3:
		//我参与的
		need, err := domain.GetParticipantMembers(orgId, currentUserId)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.ConditionHandleError, err)
		}
		condParam[consts.TcId+" In"] = db.In(need)
	}
	//默认查询没有被删除的
	condParam[consts.TcIsDelete] = consts.AppIsNoDelete

	if params["id"] != nil {
		condParam[consts.TcId] = db.In(db.Raw("select p.id as id from ppm_pro_project p where p.id = ? and (p.public_status = 1 or p.id in (SELECT DISTINCT pr.project_id FROM ppm_pro_project_relation pr WHERE pr.relation_id = ? AND relation_type in (1,2,3) AND pr.is_delete = 2)) and p.is_delete = 2", params["id"], currentUserId))
	} else {
		condParam[consts.TcId] = db.In(db.Raw("select p.id as id from ppm_pro_project p where (p.public_status = 1 or p.id in (SELECT DISTINCT pr.project_id FROM ppm_pro_project_relation pr WHERE pr.relation_id = ? AND relation_type in (1,2,3) AND pr.is_delete = 2)) and p.is_delete = 2", currentUserId))
	}

	return condParam, nil
}

func converRelateType(relateType *int64, relationType interface{}) {
	if val, ok := relationType.(enjson.Number); ok {
		*relateType, _ = val.Int64()
	} else if val, ok := relationType.(int64); ok {
		*relateType = val
	} else {
		if relationType == float64(1) {
			*relateType = int64(1)
		} else if relationType == float64(2) {
			*relateType = int64(2)
		} else if relationType == float64(3) {
			*relateType = int64(3)
		}
	}
}

//状态类型,1未开始2进行中3已完成4未完成
func condStatusAssembly(cond db.Cond, orgId int64, status int) errs.SystemErrorInfo {
	var statusIds []int64 = nil
	if status == 4 {
		notStartedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryProject, consts.ProcessStatusTypeNotStarted)
		if err != nil {
			log.Errorf(getProcessError, err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		processingIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryProject, consts.ProcessStatusTypeProcessing)
		if err != nil {
			log.Errorf(getProcessError, err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		statusIds = append(*notStartedIds, *processingIds...)
	} else if status == 3 {
		finishedId, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryProject, consts.ProcessStatusTypeCompleted)
		if err != nil {
			log.Errorf(getProcessError, err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		statusIds = *finishedId
	} else if status == 1 {
		notStartedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryProject, consts.ProcessStatusTypeNotStarted)
		if err != nil {
			log.Errorf(getProcessError, err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		statusIds = *notStartedIds
	} else if status == 2 {
		processingIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryProject, consts.ProcessStatusTypeProcessing)
		if err != nil {
			log.Errorf(getProcessError, err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		statusIds = *processingIds
	}
	cond[consts.TcStatus+" in"] = db.In(statusIds)
	return nil
}

func GetProjectCondAssemblyByInput(input *vo.ProjectsReq, currentUserId int64, orgId int64) (db.Cond, errs.SystemErrorInfo) {
	condParam := make(db.Cond)

	if input == nil {
		input = &vo.ProjectsReq{}
	}
	err := dealInputRelateType(condParam, input, currentUserId, orgId)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ConditionHandleError, err)
	}

	//默认查询没有被删除的
	condParam[consts.TcIsDelete] = consts.AppIsNoDelete
	if input.ID != nil {
		condParam[consts.TcId] = input.ID
	}
	if input.Name != nil {
		condParam[consts.TcName] = db.Like("%" + *input.Name + "%")
	}
	if input.Status != nil {
		condParam[consts.TcStatus] = input.Status
	}
	if input.StatusType != nil {
		err := condStatusAssembly(condParam, orgId, *input.StatusType)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.ConditionHandleError, err)
		}
	}
	if len(input.OwnerIds) > 0 {
		condParam[consts.TcOwner+" "] = db.In(input.OwnerIds)
	}
	if len(input.CreatorIds) > 0 {
		condParam[consts.TcCreator+" "] = db.In(input.CreatorIds)
	}
	if input.Owner != nil {
		condParam[consts.TcOwner] = input.Owner
	}
	if input.ProjectTypeID != nil {
		condParam[consts.TcProjectTypeId] = input.ProjectTypeID
	}
	if input.IsFiling != nil {
		if *input.IsFiling == 1 || *input.IsFiling == 2 {
			condParam[consts.TcIsFiling] = input.IsFiling
		}
	} else {
		//默认查未归档
		condParam[consts.TcIsFiling] = consts.AppIsNotFilling
	}
	if input.PriorityID != nil {
		condParam[consts.TcPriorityId] = input.PriorityID
	}
	if input.PlanEndTime != nil {
		condParam[consts.TcPlanEndTime] = db.Lte(input.PlanEndTime)
	}
	if input.PlanStartTime != nil {
		condParam[consts.TcPlanStartTime] = db.Gte(input.PlanStartTime)
	}

	//拿到当前用户的管理员flag
	adminFlag, err := rolefacade.GetUserAdminFlagRelaxed(orgId, currentUserId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	args := []interface{}{orgId}
	sql := "select p.id as id from ppm_pro_project p where p.org_id = ?"
	if input.ID != nil {
		sql += " and p.id = ?"
		args = append(args, *input.ID)
	}

	//不是超级管理员莫得私有项目查看权
	if !adminFlag.IsAdmin {
		sql += " and (p.public_status = 1 or p.id in (SELECT DISTINCT pr.project_id FROM ppm_pro_project_relation pr WHERE pr.relation_id = ? AND relation_type in (1,2,3) AND pr.is_delete = 2))"
		args = append(args, currentUserId)
	}

	condParam[consts.TcId] = db.In(db.Raw(sql, args...))

	if len(input.Participants) > 0 {
		idStr := strings.Replace(strings.Trim(fmt.Sprint(input.Participants), "[]"), " ", ",", -1)
		sql := "select distinct(project_id) as id from ppm_pro_project_relation where is_delete = 2 and relation_type = 2 and relation_id in (" + idStr + ")"
		condParam[consts.TcId+" "] = db.In(db.Raw(sql))
	}
	if len(input.Followers) > 0 {
		idStr := strings.Replace(strings.Trim(fmt.Sprint(input.Followers), "[]"), " ", ",", -1)
		sql := "select distinct(project_id) as id from ppm_pro_project_relation where is_delete = 2 and relation_type = 3 and relation_id in (" + idStr + ")"
		condParam[consts.TcId+"  "] = db.In(db.Raw(sql))
	}

	return condParam, nil
}

func dealInputRelateType(condParam db.Cond, input *vo.ProjectsReq, currentUserId int64, orgId int64) errs.SystemErrorInfo {
	if input.RelateType != nil {
		switch *input.RelateType {
		case 0:
			//所有
		case 1:
			//我发起的
			condParam[consts.TcCreator] = currentUserId
		case 2:
			//我负责的
			condParam[consts.TcOwner] = db.Eq(currentUserId)
		case 3:
			//我参与的
			need, err := domain.GetParticipantMembers(orgId, currentUserId)
			if err != nil {
				return err
			}
			condParam[consts.TcId+" In"] = db.In(need)
		case 4:
			//我参与的和我负责的
			need, err := domain.GetParticipantMembersAndOwner(orgId, currentUserId)
			if err != nil {
				return err
			}
			condParam[consts.TcId+" In"] = db.In(need)
		}
	}
	return nil
}

func addCreatorInfo(resultList *[]*vo.Project, creatorInfo map[int64]bo.UserIDInfoBo) {
	for k, v := range *resultList {
		if _, ok := creatorInfo[v.Creator]; ok {
			creatorInfoModel := &vo.UserIDInfo{}
			copyer.Copy(creatorInfo[v.Creator], creatorInfoModel)
			(*resultList)[k].CreatorInfo = creatorInfoModel
		}
	}
}
func dealResultList(resultList *[]*vo.Project, ownerInfo map[int64]bo.UserIDInfoBo, participantInfo map[int64][]bo.UserIDInfoBo, followerInfo map[int64][]bo.UserIDInfoBo, resourceByPath map[int64]bo.ResourceBo, issueStat map[int64]bo.IssueStatistic, projectTypeLocalMap, projectDetailById maps.LocalMap) {
	for k, v := range *resultList {
		if _, ok := ownerInfo[v.ID]; ok {
			ownerInfoModel := &vo.UserIDInfo{}
			copyer.Copy(ownerInfo[v.ID], ownerInfoModel)
			(*resultList)[k].OwnerInfo = ownerInfoModel
		}
		if _, ok := participantInfo[v.ID]; ok {
			participantInfoModel := &[]*vo.UserIDInfo{}
			copyer.Copy(participantInfo[v.ID], participantInfoModel)
			(*resultList)[k].MemberInfo = *participantInfoModel
		}
		if _, ok := followerInfo[v.ID]; ok {
			followerInfoModel := &[]*vo.UserIDInfo{}
			copyer.Copy(followerInfo[v.ID], followerInfoModel)
			(*resultList)[k].FollowerInfo = *followerInfoModel
		}
		if _, ok := resourceByPath[v.ResourceID]; ok {
			resource := resourceByPath[v.ResourceID]
			coverUrl := util.JointUrl(resource.Host, resource.Path)

			thumbnailUrl := util.GetCompressedPath(coverUrl, resource.Type)
			(*resultList)[k].ResourcePath = thumbnailUrl
			(*resultList)[k].ResourceCompressedPath = thumbnailUrl
		}
		if _, ok := issueStat[v.ID]; ok {
			(*resultList)[k].AllIssues = issueStat[v.ID].All
			(*resultList)[k].FinishIssues = issueStat[v.ID].Finish
			(*resultList)[k].OverdueIssues = issueStat[v.ID].Overdue
		}
		if times.GetUnixTime(*v.PlanStartTime) <= 0 {
			(*resultList)[k].PlanStartTime = nil
		}
		if times.GetUnixTime(*v.PlanEndTime) <= 0 {
			(*resultList)[k].PlanEndTime = nil
		}
		if projectTypeInterface, ok := projectTypeLocalMap[v.ProjectTypeID]; ok {
			projectType := projectTypeInterface.(bo.ProjectTypeBo)
			(*resultList)[k].ProjectTypeName = projectType.Name
			(*resultList)[k].ProjectTypeLangCode = projectType.LangCode
		}
		if projectDetailInterface, ok := projectDetailById[v.ID]; ok {
			projectDetail := projectDetailInterface.(bo.ProjectDetailBo)
			(*resultList)[k].IsSyncOutCalendar = projectDetail.IsSyncOutCalendar
		}
		//获取项目状态
		allProjectStatus, err := domain.GetProjectStatus(v.OrgID, v.ID)
		if err != nil {
			log.Error(err)
			continue
		}
		statusInfo := []*vo.HomeIssueStatusInfo{}
		//项目状态去除未开始
		statusNeedUpdate := false
		var processingStatus int64
		var statusIds []int64
		for _, val := range *allProjectStatus {
			statusIds = append(statusIds, val.StatusId)
			if val.StatusType == consts.ProcessStatusTypeNotStarted {
				if val.StatusId == v.Status {
					statusNeedUpdate = true
				}
				continue
			}
			if val.StatusType == consts.ProcessStatusTypeProcessing {
				processingStatus = val.StatusId
			}
			info := vo.HomeIssueStatusInfo{
				Type:        val.StatusType,
				ID:          val.StatusId,
				Name:        val.Name,
				DisplayName: &val.DisplayName,
				BgStyle:     val.BgStyle,
				FontStyle:   val.FontStyle,
			}
			statusInfo = append(statusInfo, &info)
		}
		(*resultList)[k].AllStatus = statusInfo
		//如果项目是未开始则改为进行中
		if ok, _ := slice.Contain(statusIds, v.Status); !ok {
			statusNeedUpdate = true
		}
		if statusNeedUpdate && processingStatus != 0 {
			(*resultList)[k].Status = processingStatus
		}
	}
}

func CreateProject(reqVo projectvo.CreateProjectReqVo) (*vo.Project, errs.SystemErrorInfo) {
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	input := reqVo.Input
	sourceChannel := reqVo.SourceChannel

	isError, err := checkAuth(&currentUserId, &orgId)

	if isError {
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProject)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}

	blankString := consts.BlankString
	var zero int64 = 0

	checkErr := assignmentInput(zero, &input, blankString)
	if checkErr != nil {
		log.Error(checkErr)
		return nil, checkErr
	}

	entity := &bo.ProjectBo{
		Id:            id,
		OrgId:         orgId,
		Code:          *input.Code,
		Name:          input.Name,
		PreCode:       *input.PreCode,
		Owner:         input.Owner,
		PriorityId:    *input.PriorityID,
		PublicStatus:  input.PublicStatus,
		ProjectTypeId: *input.ProjectTypeID,
		IsFiling:      2,
		Remark:        *input.Remark,
		Creator:       currentUserId,
		CreateTime:    types.NowTime(),
		Updator:       currentUserId,
		UpdateTime:    types.NowTime(),
		Version:       1,
		IsDelete:      consts.AppIsNoDelete,
	}
	initStatusErr := initProjectTypeAndProcessStatus(orgId, entity)
	if initStatusErr != nil {
		return nil, initStatusErr
	}

	isRepeat, err := checkRepeat(err, input, orgId, entity)

	if isRepeat {
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	//创建project
	//new, addedMemberIds, err := domain.CreateProject(*entity, orgId, currentUserId, input.ResourcePath, input.ResourceType, input.MemberIds, input.FollowerIds)
	new, addedMemberIds, err := domain.CreateProject(*entity, orgId, currentUserId, input)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	result := &vo.Project{}
	err = copyer.Copy(new, result)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}

	//创建日历
	if sourceChannel == consts.AppSourceChannelFeiShu {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("捕获到的错误：%s", r)
				}
			}()
			domain.CreateCalendar(input.IsSyncOutCalendar, orgId, id, currentUserId, addedMemberIds)
		}()
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = input.Name
		projectTrendsBo := bo.ProjectTrendsBo{
			PushType:            consts.PushTypeCreateProject,
			OrgId:               orgId,
			ProjectId:           id,
			OperatorId:          currentUserId,
			BeforeChangeMembers: []int64{},
			AfterChangeMembers:  addedMemberIds,
			NewValue:            json.ToJsonIgnoreError(entity),
			Ext:                 ext,
			SourceChannel: sourceChannel,
		}
		domain.PushProjectTrends(projectTrendsBo)
	})
	asyn.Execute(func() {
		PushAddProjectNotice(orgId, result.ID)
	})
	return result, nil
}

func initProjectTypeAndProcessStatus(orgId int64, entity *bo.ProjectBo) errs.SystemErrorInfo {
	projectTypeId, status, err := domain.GetTypeAndStatus(orgId, entity.ProjectTypeId, 0)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	entity.ProjectTypeId = projectTypeId
	entity.Status = status

	return nil
}

//校验重复
func checkRepeat(err error, input vo.CreateProjectReq, orgId int64, entity *bo.ProjectBo) (isRepeatError bool, repeatErr errs.SystemErrorInfo) {
	_, err = domain.JudgeRepeatProjectName(&input.Name, orgId, nil)
	if err != nil {
		return true, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	_, err = domain.JudgeRepeatProjectPreCode(input.PreCode, orgId, nil)
	if err != nil {
		return true, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	if ok, _ := slice.Contain([]int{consts.PublicProject, consts.PrivateProject}, entity.PublicStatus); !ok {
		return true, errs.BuildSystemErrorInfo(errs.ProjectDomainError, errors.New("项目可见性选择有误"))
	}
	//entity.PlanStartTime input.PlanStartTime的地址
	if input.PlanStartTime != nil && input.PlanStartTime.IsNotNull() {

		(*entity).PlanStartTime = *input.PlanStartTime
	} else {

		PlanStartTime := types.Time(consts.BlankTimeObject)

		(*entity).PlanStartTime = PlanStartTime
	}

	//entity.planEndTime的指针变量等于 input.PlanEndTime的地址
	if input.PlanEndTime != nil && input.PlanEndTime.IsNotNull() {

		(*entity).PlanEndTime = *input.PlanEndTime
	} else {
		BlankTime := types.Time(consts.BlankTimeObject)

		(*entity).PlanEndTime = BlankTime
	}

	if time.Time(entity.PlanEndTime).After(consts.BlankTimeObject) && time.Time(entity.PlanStartTime).After(time.Time(entity.PlanEndTime)) {
		return true, errs.BuildSystemErrorInfo(errs.CreateProjectTimeError)
	}
	return false, nil
}

//校验权限
func checkAuth(currentUserId *int64, orgId *int64) (isError bool, error error) {

	err := domain.AuthOrg(*orgId, *currentUserId, consts.RoleOperationPathOrgProject, consts.RoleOperationCreate)
	if err != nil {
		log.Error(err)
		return true, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	return false, nil
}

func assignmentInput(zero int64, input *vo.CreateProjectReq, blankString string) errs.SystemErrorInfo {
	if strings.Trim(input.Name, " ") == "" {
		return errs.ProjectNameEmpty
	}

	//if strs.Len(input.Name) > 256 {
	//	return errs.BuildSystemErrorInfo(errs.ProjectNameLenError)
	//}
	isNameRight := format.VerifyProjectNameFormat(input.Name)
	if !isNameRight {
		log.Error(errs.InvalidProjectNameError)
		return errs.InvalidProjectNameError
	}

	if input.Code == nil {
		input.Code = &blankString
	}
	if strs.Len(*input.Code) > 64 {
		return errs.BuildSystemErrorInfo(errs.ProjectCodeLenError)
	}

	if input.PreCode == nil {
		input.PreCode = &blankString
	}
	//if strs.Len(*input.PreCode) > 16 {
	//	return errs.BuildSystemErrorInfo(errs.ProjectPreCodeLenError)
	//}
	isPreCodeRight := format.VerifyProjectPreviousCodeFormat(*input.PreCode)
	if !isPreCodeRight {
		log.Error(errs.InvalidProjectPreCodeError)
		return errs.InvalidProjectPreCodeError
	}

	if input.PriorityID == nil {
		input.PriorityID = &zero
	}
	if input.ProjectTypeID == nil {
		input.ProjectTypeID = &zero
	}
	if input.Remark == nil {
		input.Remark = &blankString
	}
	//if strs.Len(*input.Remark) > 512 {
	//	return errs.BuildSystemErrorInfo(errs.ProjectRemarkLenError)
	//}
	isRemarkRight := format.VerifyProjectRemarkFormat(*input.Remark)
	if !isRemarkRight {
		log.Error(errs.InvalidProjectRemarkError)
		return errs.InvalidProjectRemarkError
	}

	return nil
}

func UpdateProject(reqVo projectvo.UpdateProjectReqVo) (*vo.Project, errs.SystemErrorInfo) {
	currentUserId := reqVo.UserId
	orgId := reqVo.OrgId
	input := reqVo.Input
	sourceChannel := reqVo.SourceChannel

	log.Infof(consts.UserLoginSentence, currentUserId, orgId)

	err := domain.AuthProject(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProProConfig, consts.RoleOperationModify)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	originProjectInfo, err := domain.GetProjectInfo(input.ID, orgId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, errors.New("项目不存在或已被删除"))
	}
	if originProjectInfo.IsFiling == consts.ProjectIsFiling {
		return nil, errs.ProjectIsFilingYet
	}
	entity := &bo.UpdateProjectBo{}
	newValue := &bo.ProjectBo{}
	copyer.Copy(input, entity)
	copyer.Copy(originProjectInfo, newValue)

	old := &map[string]interface{}{}
	new := &map[string]interface{}{}
	changeList := []bo.TrendChangeListBo{}

	upd, err := domain.UpdateProjectCondAssembly(*entity, orgId, old, new, originProjectInfo, &changeList)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	afterOwner := originProjectInfo.Owner
	if util.FieldInUpdate(input.UpdateFields, "owner") && input.Owner != nil && originProjectInfo.Owner != *input.Owner {
		(*old)["owner"] = originProjectInfo.Owner
		(*new)["owner"] = *input.Owner
		afterOwner = *input.Owner
	}

	if util.FieldInUpdate(input.UpdateFields, "resourcePath") && util.FieldInUpdate(input.UpdateFields, "resourceType") && input.ResourcePath != nil && input.ResourceType != nil {
		//资源列表
		resourcesRespVo := resourcefacade.GetResourceById(resourcevo.GetResourceByIdReqVo{GetResourceByIdReqBody: resourcevo.GetResourceByIdReqBody{ResourceIds: []int64{originProjectInfo.ResourceId}}})
		if resourcesRespVo.Failure() {
			log.Error(resourcesRespVo.Error())
			return nil, resourcesRespVo.Error()
		}
		oldResourcePath := ""
		if len(resourcesRespVo.ResourceBos) > 0 {
			oldResourcePath = resourcesRespVo.ResourceBos[0].Host + resourcesRespVo.ResourceBos[0].Path
		}
		changeList = append(changeList, bo.TrendChangeListBo{
			Field:     "resourcePath",
			FieldName: consts.ProjectResourcePath,
			NewValue:  *input.ResourcePath,
			OldValue:  oldResourcePath,
		})
	}

	beforeMemberIds, afterMemberIds, err := domain.UpdateProject(orgId, currentUserId, &upd, originProjectInfo.Owner, *entity)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	//更新关注人
	oldFollower, newFollower, err := domain.UpdateFollower(*entity, currentUserId, orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err1 := domain.RefreshProjectAuthBo(orgId, input.ID)
	if err1 != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
	}

	pushType := consts.PushTypeUpdateProject
	beforeChangeMembersSet := set.New(set.ThreadSafe)
	for _, member := range beforeMemberIds {
		beforeChangeMembersSet.Add(member)
	}
	afterChangeMembersSet := set.New(set.ThreadSafe)
	for _, member := range afterMemberIds {
		afterChangeMembersSet.Add(member)
	}

	deletedMembersSet := set.Difference(beforeChangeMembersSet, afterChangeMembersSet)
	addedMembersSet := set.Difference(afterChangeMembersSet, beforeChangeMembersSet)

	delFollower, addFollower := util.GetDifMemberIds(oldFollower, newFollower)
	if deletedMembersSet.Size() != 0 || addedMembersSet.Size() != 0 || originProjectInfo.Owner != afterOwner || len(delFollower) > 0 || len(addFollower) > 0 {
		pushType = consts.PushTypeUpdateProjectMembers
		(*old)["memberIds"] = deletedMembersSet
		(*new)["memberIds"] = addedMembersSet
	}

	//更新同步日历
	updateCalendarErr := updateProjectByUpdateCalendarSet(input, orgId, currentUserId)
	if updateCalendarErr != nil {
		return nil, updateCalendarErr
	}
	//更新日历
	oldAll := append(append(beforeMemberIds, oldFollower...), originProjectInfo.Owner)
	newAll := append(append(afterMemberIds, newFollower...), afterOwner)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("捕获到的错误：%s", r)
			}
		}()
		domain.UpdateCalendar(input, orgId, oldAll, newAll, currentUserId)
	}()

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = originProjectInfo.Name
		ext.ChangeList = changeList
		domain.PushProjectTrends(bo.ProjectTrendsBo{
			PushType:            pushType,
			OrgId:               orgId,
			ProjectId:           input.ID,
			OperatorId:          currentUserId,
			BeforeChangeMembers: beforeMemberIds,
			AfterChangeMembers:  afterMemberIds,
			//OldValue:            json.ToJsonIgnoreError(originProjectInfo),
			//NewValue:            json.ToJsonIgnoreError(newValue),
			BeforeOwner:           originProjectInfo.Owner,
			AfterOwner:            afterOwner,
			BeforeChangeFollowers: oldFollower,
			AfterChangeFollowers:  newFollower,
			OldValue:              json.ToJsonIgnoreError(old),
			NewValue:              json.ToJsonIgnoreError(new),
			Ext:                   ext,

			SourceChannel: sourceChannel,
		})
	})
	asyn.Execute(func() {
		PushModifyProjectNotice(orgId, input.ID)
	})
	return &vo.Project{ID: input.ID}, nil
}

func updateProjectByUpdateCalendarSet(input vo.UpdateProjectReq, orgId int64, currentUserId int64) errs.SystemErrorInfo {
	if !util.FieldInUpdate(input.UpdateFields, "isSyncOutCalendar") || input.IsSyncOutCalendar == nil {
		return nil
	}
	if ok, _ := slice.Contain([]int{consts.IsSyncOutCalendar, consts.IsNotSyncOutCalendar}, *input.IsSyncOutCalendar); !ok {
		return nil

	}
	projectDetail, err := domain.GetProjectDetailByProjectIdBo(input.ID, orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	if projectDetail.IsSyncOutCalendar != *input.IsSyncOutCalendar {
		updateErr := domain.UpdateProjectDetail(&bo.ProjectDetailBo{
			Id:                projectDetail.Id,
			IsSyncOutCalendar: *input.IsSyncOutCalendar,
			Updator:           currentUserId,
		})
		if updateErr != nil {
			log.Error(updateErr)
			return updateErr
		}
		cacheErr := domain.DeleteProjectCalendarInfo(orgId, input.ID)
		if cacheErr != nil {
			log.Error(cacheErr)
			return cacheErr
		}
		if *input.IsSyncOutCalendar == consts.IsSyncOutCalendar {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Errorf("捕获到的错误：%s", r)
					}
				}()
				domain.SyncCalendarConfirm(orgId, currentUserId, projectDetail.ProjectId)
			}()
		}
	}
	return nil
}

func QuitProject(reqVo projectvo.ProjectIdReqVo) (*vo.QuitResult, errs.SystemErrorInfo) {
	currentUserId := reqVo.UserId
	orgId := reqVo.OrgId
	projectId := reqVo.ProjectId
	sourceChannel := reqVo.SourceChannel

	result := &vo.QuitResult{}
	log.Infof("当前登录用户 %d 组织 %d", currentUserId, orgId)

	project, err := domain.GetProjectInfo(projectId, orgId)
	if err != nil {
		result.IsQuitted = false
		return result, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	judgeErr := JudgeProjectFiling(orgId, projectId)
	if judgeErr != nil {
		log.Error(judgeErr)
		return nil, judgeErr
	}
	//负责人不允许退出项目
	if project.Owner == currentUserId {
		result.IsQuitted = false
		return result, errs.BuildSystemErrorInfo(errs.NotAllowQuitProject)
	}
	//如果不是项目成员不允许退出项目
	member, err := domain.JudgeIsProjectMember(currentUserId, orgId, projectId)

	if err != nil {
		result.IsQuitted = false
		return result, errs.BuildSystemErrorInfo(errs.NotProjectParticipant)
	}

	err = domain.QuitProject(currentUserId, orgId, project.Owner, projectId, member.Id)
	if err != nil {
		result.IsQuitted = false
		return result, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	err1 := domain.RefreshProjectAuthBo(orgId, projectId)
	if err1 != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
	}

	result.IsQuitted = true

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = project.Name
		domain.PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeUnbindProject,
			OrgId:      orgId,
			ProjectId:  projectId,
			OperatorId: currentUserId,
			Ext:        ext,

			SourceChannel: sourceChannel,
		})
	})

	asyn.Execute(func() {
		domain.UpdateCalendarAttendees(orgId, "", []int64{}, []int64{currentUserId}, projectId)
	})

	asyn.Execute(func() {
		PushModifyProjectNotice(orgId, projectId)
	})
	return result, nil
}

func StarProject(reqVo projectvo.ProjectIdReqVo) (*vo.OperateProjectResp, errs.SystemErrorInfo) {
	currentUserId := reqVo.UserId
	orgId := reqVo.OrgId
	projectId := reqVo.ProjectId
	sourceChannel := reqVo.SourceChannel

	log.Infof(consts.UserLoginSentence, currentUserId, orgId)

	//校验权限
	err := domain.AuthProject(orgId, currentUserId, projectId, consts.RoleOperationPathOrgProject, consts.RoleOperationAttention)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	projectInfo, err := domain.GetProject(orgId, projectId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist, err)
	}

	res := &vo.OperateProjectResp{
		IsSuccess: false,
	}

	//isExist, err := domain.JudgeIsFollower(projectId, currentUserId, orgId)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	//}
	//
	//if isExist {
	//	return res, errs.BuildSystemErrorInfo(errs.AlreadyStarProject)
	//}

	//isSuccess, err := domain.AddMember(projectId, orgId, currentUserId, currentUserId, consts.IssueRelationTypeFollower)
	//if err != nil {
	//	return res, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}

	updateProjectRelation := domain.UpdateProjectRelation(currentUserId, orgId, projectId, consts.IssueRelationTypeFollower, []int64{currentUserId})
	if updateProjectRelation != nil {
		log.Error(updateProjectRelation)
		return nil, updateProjectRelation
	}

	err1 := domain.RefreshProjectAuthBo(orgId, projectId)
	if err1 != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
	}
	res.IsSuccess = true

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = projectInfo.Name
		domain.PushProjectTrends(bo.ProjectTrendsBo{
			PushType:      consts.PushTypeStarProject,
			OrgId:         orgId,
			ProjectId:     projectId,
			OperatorId:    currentUserId,
			Ext:           ext,
			SourceChannel: sourceChannel,
		})
	})

	asyn.Execute(func() {
		domain.UpdateCalendarAttendees(orgId, "", []int64{currentUserId}, []int64{}, projectId)
	})

	asyn.Execute(func() {
		PushModifyProjectNotice(orgId, projectId)
	})
	return res, nil
}

func UnstarProject(orgId, currentUserId int64, projectId int64) (*vo.OperateProjectResp, errs.SystemErrorInfo) {

	log.Infof(consts.UserLoginSentence, currentUserId, orgId)

	err := domain.AuthProject(orgId, currentUserId, projectId, consts.RoleOperationPathOrgProject, consts.RoleOperationUnAttention)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	projectInfo, err := domain.GetProject(orgId, projectId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist, err)
	}

	isExist, err := domain.JudgeIsFollower(projectId, currentUserId, orgId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	res := &vo.OperateProjectResp{
		IsSuccess: false,
	}
	if !isExist {
		return res, errs.BuildSystemErrorInfo(errs.NotYetStarProject)
	}

	isSuccess, err := domain.DeleteMember(projectId, orgId, currentUserId, currentUserId, consts.IssueRelationTypeFollower)
	if err != nil {
		return res, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	err1 := domain.RefreshProjectAuthBo(orgId, projectId)
	if err1 != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
	}

	res.IsSuccess = isSuccess

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = projectInfo.Name
		domain.PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeUnstarProject,
			OrgId:      orgId,
			ProjectId:  projectId,
			OperatorId: currentUserId,
			Ext:        ext,
		})
	})

	asyn.Execute(func() {
		domain.UpdateCalendarAttendees(orgId, "", []int64{}, []int64{currentUserId}, projectId)
	})

	asyn.Execute(func() {
		PushModifyProjectNotice(orgId, projectId)
	})
	return res, nil
}

func ProjectStatistics(orgId, id int64) (*vo.ProjectStatisticsResp, errs.SystemErrorInfo) {

	exist := domain.JudgeProjectIsExist(orgId, id)
	if !exist {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}

	result := &vo.ProjectStatisticsResp{}
	projectStat, err := domain.StatProject(orgId, id)
	if err != nil {
		return result, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	_ = copyer.Copy(projectStat, result)

	return result, nil
}

func UpdateProjectStatus(reqVo projectvo.UpdateProjectStatusReqVo) (*vo.Void, errs.SystemErrorInfo) {
	input := reqVo.Input
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	sourceChannel := reqVo.SourceChannel

	projectId := input.ProjectID

	projectBo, err1 := domain.GetProject(orgId, projectId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}
	if projectBo.IsFiling == consts.ProjectIsFiling {
		return nil, errs.ProjectIsFilingYet
	}

	err := domain.AuthProject(orgId, currentUserId, projectId, consts.RoleOperationPathOrgProProConfig, consts.RoleOperationModifyStatus)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	err2 := domain.UpdateProjectStatus(*projectBo, input.NextStatusID)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err2)
	}

	refreshErr := domain.RefreshProjectAuthBo(orgId, projectId)
	if refreshErr != nil {
		log.Error(refreshErr)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, refreshErr)
	}



	asyn.Execute(func() {
		operateObjProperty := consts.TrendsOperObjPropertyNameStatus
		oldValueMap := map[string]interface{}{
			operateObjProperty: projectBo.Status,
		}
		newValueMap := map[string]interface{}{
			operateObjProperty: input.NextStatusID,
		}

		//状态列表
		statusList, err := processfacade.GetProcessStatusListByCategoryRelaxed(orgId, consts.ProcessStatusCategoryProject)
		if err != nil {
			log.Error(strs.ObjectToString(err))
			return
		}

		change := bo.TrendChangeListBo{
			Field:     "status",
			FieldName: consts.Status,
		}
		for _, v := range statusList {
			if v.StatusId == projectBo.Status {
				change.OldValue = v.Name
			} else if v.StatusId == input.NextStatusID {
				change.NewValue = v.Name
			}
		}
		changeList := []bo.TrendChangeListBo{}
		changeList = append(changeList, change)
		ext := bo.TrendExtensionBo{
			ObjName:    projectBo.Name,
			ChangeList: changeList,
		}

		domain.PushProjectTrends(bo.ProjectTrendsBo{
			PushType:           consts.PushTypeUpdateProjectStatus,
			OrgId:              orgId,
			ProjectId:          projectId,
			OperatorId:         currentUserId,
			OperateObjProperty: operateObjProperty,
			OldValue:           json.ToJsonIgnoreError(oldValueMap),
			NewValue:           json.ToJsonIgnoreError(newValueMap),
			Ext:                ext,

			SourceChannel: sourceChannel,
		})
	})

	asyn.Execute(func() {
		PushModifyProjectNotice(orgId, projectId)
	})
	return &vo.Void{
		ID: projectId,
	}, nil
}

func ProjectInfo(orgId int64, input vo.ProjectInfoReq, sourceChannel string) (*vo.ProjectInfo, errs.SystemErrorInfo) {

	projectId := input.ProjectID

	projectBo, err1 := domain.GetProject(orgId, projectId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}

	result := &vo.ProjectInfo{}
	err2 := copyer.Copy(projectBo, result)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	resourcesRespVo := resourcefacade.GetResourceById(resourcevo.GetResourceByIdReqVo{GetResourceByIdReqBody: resourcevo.GetResourceByIdReqBody{ResourceIds: []int64{projectBo.ResourceId}}})
	if resourcesRespVo.Failure() {
		log.Error(resourcesRespVo.Message)
		return nil, resourcesRespVo.Error()
	}
	resourceEntities := resourcesRespVo.ResourceBos

	for _, v := range resourceEntities {
		result.ResourcePath = v.Host + v.Path
	}

	//项目相关人员
	ownerInfo, participantInfo, followerInfo, creatorInfo, err := domain.GetProjectMemberInfo([]int64{projectId}, orgId, []int64{projectBo.Creator}, sourceChannel)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	if _, ok := ownerInfo[projectId]; ok {
		ownerInfoModel := &vo.UserIDInfo{}
		copyer.Copy(ownerInfo[projectId], ownerInfoModel)
		result.OwnerInfo = ownerInfoModel
	}
	if _, ok := participantInfo[projectId]; ok {
		participantInfoModel := &[]*vo.UserIDInfo{}
		copyer.Copy(participantInfo[projectId], participantInfoModel)
		result.MemberInfo = *participantInfoModel
	}
	if _, ok := followerInfo[projectId]; ok {
		followerInfoModel := &[]*vo.UserIDInfo{}
		copyer.Copy(followerInfo[projectId], followerInfoModel)
		result.FollowerInfo = *followerInfoModel
	}
	if _, ok := creatorInfo[projectBo.Creator]; ok {
		creatorInfoModel := &vo.UserIDInfo{}
		copyer.Copy(creatorInfo[projectBo.Creator], creatorInfoModel)
		result.CreatorInfo = creatorInfoModel
	}

	//获取项目状态
	allProjectStatus, err := domain.GetProjectStatus(orgId, projectId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	statusInfo := []*vo.HomeIssueStatusInfo{}

	//项目状态去除未开始
	statusNeedUpdate := false
	var processingStatus int64
	var statusIds []int64
	for _, v := range *allProjectStatus {
		statusIds = append(statusIds, v.StatusId)
		if v.StatusType == consts.ProcessStatusTypeNotStarted {
			if v.StatusId == result.Status {
				statusNeedUpdate = true
			}
			continue
		}
		if v.StatusType == consts.ProcessStatusTypeProcessing {
			processingStatus = v.StatusId
		}
		info := vo.HomeIssueStatusInfo{
			Type:        v.StatusType,
			ID:          v.StatusId,
			Name:        v.Name,
			DisplayName: &v.DisplayName,
			BgStyle:     v.BgStyle,
			FontStyle:   v.FontStyle,
		}
		statusInfo = append(statusInfo, &info)
	}
	if ok, _ := slice.Contain(statusIds, result.Status); !ok {
		statusNeedUpdate = true
	}
	result.AllStatus = statusInfo
	if processingStatus != 0 && statusNeedUpdate {
		result.Status = processingStatus
	}

	return result, nil
}

func GetProjectProcessId(orgId int64, projectId int64, projectObjectTypeId int64) (int64, errs.SystemErrorInfo) {
	return domain.GetProjectProcessId(orgId, projectId, projectObjectTypeId)
}

//通过项目类型langCode获取项目列表
func GetProjectBoListByProjectTypeLangCode(orgId int64, projectTypeLangCode *string) ([]bo.ProjectBo, errs.SystemErrorInfo) {
	return domain.GetProjectBoListByProjectTypeLangCode(orgId, projectTypeLangCode)
}

func GetSimpleProjectInfo(orgId int64, ids []int64) (*[]vo.Project, errs.SystemErrorInfo) {
	list, err := domain.GetProjectBoList(orgId, ids)
	if err != nil {
		return nil, err
	}
	projectVo := &[]vo.Project{}
	copyErr := copyer.Copy(list, projectVo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return projectVo, nil
}

func GetProjectRelation(projectId int64, relationType []int64) ([]projectvo.ProjectRelationList, errs.SystemErrorInfo) {
	bos, err := domain.GetProjectRelationByType(projectId, relationType)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	res := []projectvo.ProjectRelationList{}
	for _, v := range *bos {
		res = append(res, projectvo.ProjectRelationList{
			Id:           v.Id,
			RelationId:   v.RelationId,
			RelationType: v.RelationType,
		})
	}

	return res, nil
}

func ArchiveProject(orgId, userId, projectId int64) (*vo.Void, errs.SystemErrorInfo) {
	err := domain.AuthProject(orgId, userId, projectId, consts.RoleOperationPathOrgProProConfig, consts.RoleOperationFiling)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	projectInfo, err := domain.GetProject(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if projectInfo.IsFiling == consts.ProjectIsFiling {
		return nil, errs.ProjectIsFilingYet
	}

	_, updateErr := dao.UpdateProjectByOrg(projectId, orgId, mysql.Upd{
		consts.TcIsFiling: consts.ProjectIsFiling,
	})
	if updateErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, updateErr)
	}

	err1 := domain.RefreshProjectAuthBo(orgId, projectId)
	if err1 != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
	}

	asyn.Execute(func() {
		PushModifyProjectNotice(orgId, projectId)
	})
	return &vo.Void{
		ID: projectId,
	}, nil
}

func CancelArchivedProject(orgId, userId, projectId int64) (*vo.Void, errs.SystemErrorInfo) {
	err := domain.AuthProjectWithOutArchivedCheck(orgId, userId, projectId, consts.RoleOperationPathOrgProProConfig, consts.RoleOperationUnFiling)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	_, err = domain.GetProject(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	_, updateErr := dao.UpdateProjectByOrg(projectId, orgId, mysql.Upd{
		consts.TcIsFiling: consts.ProjectIsNotFiling,
	})
	if updateErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, updateErr)
	}

	err1 := domain.RefreshProjectAuthBo(orgId, projectId)
	if err1 != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
	}

	asyn.Execute(func() {
		PushModifyProjectNotice(orgId, projectId)
	})
	return &vo.Void{
		ID: projectId,
	}, nil
}

func GetProjectInfoByOrgIds(orgIds []int64) ([]projectvo.GetProjectInfoListByOrgIdsRespVo, errs.SystemErrorInfo) {

	bos, err := domain.GetProjectInfoByOrgIds(orgIds)

	if err != nil {
		return nil, err
	}

	result := []projectvo.GetProjectInfoListByOrgIdsRespVo{}

	for _, value := range bos {

		vo := projectvo.GetProjectInfoListByOrgIdsRespVo{
			OrgId:     value.OrgId,
			ProjectId: value.Id,
			Owner:     value.Owner,
		}
		result = append(result, vo)
	}

	return result, nil

}

func GetCacheProjectInfo(reqVo projectvo.GetCacheProjectInfoReqVo) (*bo.ProjectAuthBo, errs.SystemErrorInfo) {
	projectAuthBo, err := domain.LoadProjectAuthBo(reqVo.OrgId, reqVo.ProjectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return projectAuthBo, nil
}

func OrgProjectMembers(input projectvo.OrgProjectMemberReqVo) (*projectvo.OrgProjectMemberRespVo, errs.SystemErrorInfo) {

	relationBo, err := domain.GetProjectRelationByType(input.ProjectId, []int64{consts.IssueRelationTypeOwner, consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower})

	if err != nil {
		log.Error(err)
		return nil, err
	}
	//去重所有的用户id
	distinctUserIds, err := distinctUserIds(relationBo)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed("", input.OrgId, *distinctUserIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	allMembers := &[]projectvo.OrgProjectMemberVo{}
	_ = copyer.Copy(userInfos, allMembers)

	var ownerId int64
	//参与者
	participants := make(map[int64]bool)
	//关注者
	follower := make(map[int64]bool)

	//用户id分组
	for _, value := range *relationBo {

		if value.RelationType == consts.IssueRelationTypeOwner {
			ownerId = value.RelationId
			continue
		}

		if value.RelationType == consts.IssueRelationTypeParticipant {
			participants[value.RelationId] = true
			continue
		}

		if value.RelationType == consts.IssueRelationTypeFollower {
			follower[value.RelationId] = true
			continue
		}
	}
	//返回结果数组
	participantsMemberList := make([]projectvo.OrgProjectMemberVo, 0)
	followerMemberList := make([]projectvo.OrgProjectMemberVo, 0)

	ownerMember := projectvo.OrgProjectMemberVo{}

	for _, member := range *allMembers {
		//拥有者
		if member.UserId == ownerId {
			ownerMember = member
		}

		if _, ok := participants[member.UserId]; ok {
			participantsMemberList = append(participantsMemberList, member)
		}

		if _, ok := follower[member.UserId]; ok {
			followerMemberList = append(followerMemberList, member)
		}
	}

	return &projectvo.OrgProjectMemberRespVo{
		Owner: ownerMember,
		Participants: participantsMemberList,
		Follower: followerMemberList,
		AllMembers: *allMembers,
	}, nil

}

func distinctUserIds(relationBo *[]bo.ProjectRelationBo) (*[]int64, errs.SystemErrorInfo) {

	rlen := len(*relationBo)

	if relationBo == nil || rlen < 1 {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectRelationNotExist)
	}

	allUserIds := make([]int64, rlen)

	for i := 0; i < rlen; i++ {
		allUserIds[i] = (*relationBo)[i].RelationId
	}

	//去重
	uniqueInt64 := slice.SliceUniqueInt64(allUserIds)

	return &uniqueInt64, nil
}
