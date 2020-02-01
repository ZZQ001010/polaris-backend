package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/date"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"strings"
	"time"
	"upper.io/db.v3"
)

func IterationList(orgId int64, page uint, size uint, params *vo.IterationListReq) (*vo.IterationList, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//orgId := cacheUserInfo.OrgId

	cond := db.Cond{}
	if params != nil {

		isError, err := dealParam(params, &cond, orgId)
		if isError {
			return nil, err
		}

	}
	cond[consts.TcOrgId] = orgId
	cond[consts.TcIsDelete] = consts.AppIsNoDelete

	bos, total, err := domain.GetIterationBoList(page, size, cond)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}

	resultList := &[]*vo.Iteration{}
	copyErr := copyer.Copy(bos, resultList)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	//状态列表
	statusList, err := processfacade.GetProcessStatusListByCategoryRelaxed(orgId, consts.ProcessStatusCategoryIssue)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	ownerIds := make([]int64, 0)
	for _, result := range *resultList {
		ownerId := result.Owner
		ownerIdExist, _ := slice.Contain(ownerIds, ownerId)
		if !ownerIdExist {
			ownerIds = append(ownerIds, ownerId)
		}
	}
	ownerInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(consts.AppSourceChannelDingTalk, orgId, ownerIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//转换map
	ownerMap := maps.NewMap("UserId", ownerInfos)
	statusMap := maps.NewMap("StatusId", statusList)

	//获取负责人
	for _, result := range *resultList {
		if userCacheInfo, ok := ownerMap[result.Owner]; ok {
			baseUserInfo := userCacheInfo.(bo.BaseUserInfoBo)
			ownerInfo := &vo.HomeIssueOwnerInfo{}
			ownerInfo.ID = baseUserInfo.UserId
			ownerInfo.Name = baseUserInfo.Name
			ownerInfo.Avatar = &baseUserInfo.Avatar
			result.OwnerInfo = ownerInfo
		}
		if statusCacheInfo, ok := statusMap[result.Status]; ok {
			statusCacheBo := statusCacheInfo.(bo.CacheProcessStatusBo)
			homeStatusInfoBo := domain.ConvertStatusInfoToHomeIssueStatusInfo(statusCacheBo)
			statusInfo := &vo.HomeIssueStatusInfo{}
			_ = copyer.Copy(homeStatusInfoBo, statusInfo)
			result.StatusInfo = statusInfo
		}
	}

	return &vo.IterationList{
		Total: total,
		List:  *resultList,
	}, nil
}

func dealParam(params *vo.IterationListReq, cond *db.Cond, orgId int64) (isError bool, eoor errs.SystemErrorInfo) {

	if params.Name != nil {
		(*cond)[consts.TcName] = db.Like("%" + *params.Name + "%")
	}
	if params.ProjectID != nil {
		projectId := *params.ProjectID
		_, err := domain.LoadProjectAuthBo(orgId, projectId)
		if err != nil {
			log.Error(err)
			return true, errs.BuildSystemErrorInfo(errs.IllegalityProject, err)
		}
		(*cond)[consts.TcProjectId] = projectId
	}
	if params.StatusType != nil {
		err := domain.IterationCondStatusAssembly(cond, orgId, *params.StatusType)
		if err != nil {
			log.Error(err)
			return true, errs.BuildSystemErrorInfo(errs.IterationDomainError, err)
		}
	}

	return false, nil
}

func CreateIteration(orgId, currentUserId int64, input vo.CreateIterationReq) (*vo.Void, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//currentUserId := cacheUserInfo.UserId
	//orgId := cacheUserInfo.OrgId
	projectId := input.ProjectID
	ownerId := input.Owner

	err := domain.AuthProject(orgId, currentUserId, input.ProjectID, consts.RoleOperationPathOrgProIteration, consts.RoleOperationCreate)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	//校验负责人
	if ownerId != currentUserId {
		suc := orgfacade.VerifyOrgRelaxed(orgId, ownerId)
		if !suc {
			return nil, errs.BuildSystemErrorInfo(errs.IllegalityOwner)
		}
	}

	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIteration)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}

	iterationObjectType, err3 := domain.GetProjectObjectTypeOfIteration(orgId, projectId)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err3)
	}

	//获取初始化状态
	initStatusId, err := processfacade.GetProcessInitStatusIdRelaxed(orgId, input.ProjectID, iterationObjectType.Id, consts.ProcessStatusCategoryIteration)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	entity := &bo.IterationBo{
		Id:            id,
		OrgId:         orgId,
		ProjectId:     projectId,
		Name:          input.Name,
		Owner:         ownerId,
		PlanStartTime: input.PlanStartTime,
		PlanEndTime:   input.PlanEndTime,
		Status:        initStatusId,
		Creator:       currentUserId,
		Updator:       currentUserId,
		IsDelete:      consts.AppIsNoDelete,
	}
	err1 := domain.CreateIteration(entity)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	return &vo.Void{
		ID: id,
	}, nil
}

func UpdateIteration(orgId, currentUserId int64, input vo.UpdateIterationReq) (*vo.Void, errs.SystemErrorInfo) {
	targetId := input.ID

	iterationBo, err1 := domain.GetIterationBo(targetId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}
	err := domain.AuthProject(orgId, currentUserId, iterationBo.ProjectId, consts.RoleOperationPathOrgProIteration, consts.RoleOperationModify)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	iterationNewBo := &bo.IterationBo{}
	_ = copyer.Copy(iterationBo, iterationNewBo)
	iterationUpdateBo := &bo.IterationUpdateBo{}
	upd := mysql.Upd{}

	isErr1, err1 := dealNeedUpdate(input, &upd, iterationNewBo, currentUserId, orgId)

	if isErr1 {
		return nil, err1
	}

	var planStartTime = iterationBo.PlanStartTime
	var planEndTime = iterationBo.PlanEndTime
	// && NeedUpdate(input.UpdateFields, "planEndTime")
	if NeedUpdate(input.UpdateFields, "planStartTime") {
		if input.PlanStartTime == nil {
			return nil, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "计划开始时间不能为空")
		}
		planStartTime = *input.PlanStartTime
		upd[consts.TcPlanStartTime] = date.FormatTime(planStartTime)
		iterationNewBo.PlanStartTime = planStartTime
	}
	if NeedUpdate(input.UpdateFields, "planEndTime") {
		if input.PlanEndTime == nil {
			return nil, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "计划结束时间不能为空")
		}
		planEndTime = *input.PlanEndTime
		upd[consts.TcPlanEndTime] = date.FormatTime(planEndTime)
		iterationNewBo.PlanEndTime = planEndTime
	}
	if time.Time(planEndTime).Before(time.Time(planStartTime)) {
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "开始时间应该在结束时间之前")
	}
	upd[consts.TcUpdator] = currentUserId

	iterationUpdateBo.Id = targetId
	iterationUpdateBo.Upd = upd
	iterationUpdateBo.IterationNewBo = *iterationNewBo

	err1 = domain.UpdateIteration(iterationUpdateBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	return &vo.Void{
		ID: input.ID,
	}, nil
}

//判断处理是否需要更新
func dealNeedUpdate(input vo.UpdateIterationReq, upd *mysql.Upd, iterationNewBo *bo.IterationBo, currentUserId int64, orgId int64) (bool, errs.SystemErrorInfo) {

	if NeedUpdate(input.UpdateFields, "name") {
		if input.Name == nil {
			return true, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "迭代名称不能为空")
		}
		//input.Name指针指向的值 去空格给name
		name := strings.Trim(*input.Name, " ")
		if name == "" || strs.Len(name) > 200 {
			return true, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "迭代名称不能为空且限制在200字以内")
		}
		(*upd)[consts.TcName] = name
		(*iterationNewBo).Name = name
	}

	if NeedUpdate(input.UpdateFields, "owner") {
		if input.Owner == nil {
			return true, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "负责人不能为空")
		}
		ownerId := *input.Owner
		if ownerId != currentUserId {
			suc := orgfacade.VerifyOrgRelaxed(orgId, ownerId)
			if !suc {
				return true, errs.BuildSystemErrorInfo(errs.IllegalityOwner)
			}

			(*upd)[consts.TcOwner] = ownerId
			(*iterationNewBo).Owner = ownerId
		}
	}

	return false, nil

}

func DeleteIteration(orgId, currentUserId int64, input vo.DeleteIterationReq) (*vo.Void, errs.SystemErrorInfo) {

	targetId := input.ID

	bo, err1 := domain.GetIterationBo(targetId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	err := domain.AuthProject(orgId, currentUserId, bo.ProjectId, consts.RoleOperationPathOrgProIteration, consts.RoleOperationDelete)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	err2 := domain.DeleteIteration(bo, currentUserId)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	}

	return &vo.Void{
		ID: targetId,
	}, nil
}

func IterationStatusTypeStat(orgId int64, input *vo.IterationStatusTypeStatReq) (*vo.IterationStatusTypeStatResp, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacorgIdade.GetCurrentUserRelaxed(ctx)
	//	//if err != nil {
	//	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//	//}
	//	//orgId := cacheUserInfo.OrgId

	projectId := int64(0)
	if input != nil {
		if input.ProjectID != nil {
			projectId := *input.ProjectID
			_, err := domain.LoadProjectAuthBo(orgId, projectId)
			if err != nil {
				log.Error(err)
				return nil, errs.BuildSystemErrorInfo(errs.IllegalityProject, err)
			}
		}
	}

	countBo, err1 := domain.StatisticIterationCountGroupByStatus(orgId, projectId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IterationDomainError, err1)
	}

	return &vo.IterationStatusTypeStatResp{
		NotStartTotal:   countBo.NotStartTotal,
		ProcessingTotal: countBo.ProcessingTotal,
		CompletedTotal:  countBo.FinishedTotal,
		Total:           countBo.NotStartTotal + countBo.ProcessingTotal + countBo.FinishedTotal,
	}, nil
}

func IterationIssueRelate(orgId, currentUserId int64, input vo.IterationIssueRealtionReq) (*vo.Void, errs.SystemErrorInfo) {
	iterationBo, err1 := domain.GetIterationBoByOrgId(input.IterationID, orgId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIteration)
	}

	err2 := domain.RelateIteration(orgId, iterationBo.ProjectId, iterationBo.Id, currentUserId, input.AddIssueIds, input.DelIssueIds)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
	}

	return &vo.Void{
		ID: iterationBo.Id,
	}, nil
}

func UpdateIterationStatus(orgId, currentUserId int64, input vo.UpdateIterationStatusReq) (*vo.Void, errs.SystemErrorInfo) {

	iterationBo, err1 := domain.GetIterationBoByOrgId(input.ID, orgId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIteration)
	}

	err := domain.AuthProject(orgId, currentUserId, iterationBo.ProjectId, consts.RoleOperationPathOrgProIteration, consts.RoleOperationModifyStatus)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	err2 := domain.UpdateIterationStatus(*iterationBo, input.NextStatusID, currentUserId)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IterationDomainError, err2)
	}

	return &vo.Void{
		ID: iterationBo.Id,
	}, nil
}

func IterationInfo(orgId int64, input vo.IterationInfoReq) (*vo.IterationInfoResp, errs.SystemErrorInfo) {

	iterationBo, err1 := domain.GetIterationBoByOrgId(input.ID, orgId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIteration)
	}

	iterationInfoBo, err2 := domain.GetIterationInfoBo(*iterationBo)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IterationDomainError, err2)
	}

	iterationInfoResp := &vo.IterationInfoResp{}
	err3 := copyer.Copy(iterationInfoBo, iterationInfoResp)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return iterationInfoResp, nil
}

//获取未完成的的迭代列表
func GetNotCompletedIterationBoList(orgId int64, projectId int64) ([]bo.IterationBo, errs.SystemErrorInfo) {
	return domain.GetNotCompletedIterationBoList(orgId, projectId)
}
