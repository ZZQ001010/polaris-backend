package domain

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
)

func GetIterationInfoBo(iterationBo bo.IterationBo) (*bo.IterationInfoBo, errs.SystemErrorInfo) {
	orgId := iterationBo.OrgId

	statusInfo, err := GetHomeIssueStatusInfoBo(orgId, iterationBo.Status)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProcessStatusNotExist, err)
	}

	projectInfo, err := GetHomeProjectInfoBo(orgId, iterationBo.ProjectId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist, err)
	}

	//负责人信息
	ownerId := iterationBo.Owner
	ownerBaseInfo, err := orgfacade.GetDingTalkBaseUserInfoRelaxed(orgId, ownerId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}
	ownerInfo := &bo.UserIDInfoBo{
		UserID:  ownerBaseInfo.UserId,
		Name:    ownerBaseInfo.Name,
		Avatar:  ownerBaseInfo.Avatar,
		EmplID:  ownerBaseInfo.OutUserId,
	}

	//迭代对象类型
	iterationObjectType, err3 := GetProjectObjectTypeOfIteration(orgId, projectInfo.ID)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err3)
	}

	//流程
	process, err := GetProjectProcessBo(orgId, iterationBo.ProjectId, iterationObjectType.Id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	//流程步骤-接下来的状态
	nextStatusList, err := processfacade.GetNextProcessStepStatusListRelaxed(orgId, process.Id, statusInfo.ID)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	nextStatusInfoList := make([]bo.HomeIssueStatusInfoBo, len(*nextStatusList))
	for i, nextStatus := range *nextStatusList {
		nextStatusName := nextStatus.DisplayName
		nextStatusInfoList[i] = bo.HomeIssueStatusInfoBo{
			ID:          nextStatus.StatusId,
			Name:        nextStatus.Name,
			BgStyle:     nextStatus.BgStyle,
			FontStyle:   nextStatus.FontStyle,
			Type:        nextStatus.StatusType,
			DisplayName: &nextStatusName,
		}
	}

	return &bo.IterationInfoBo{
		Iteration:  iterationBo,
		Project:    projectInfo,
		Status:     statusInfo,
		Owner:      ownerInfo,
		NextStatus: &nextStatusInfoList,
	}, nil
}
