package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
)

func UpdateIssueStatus(reqVo projectvo.UpdateIssueStatusReqVo) (*vo.Issue, errs.SystemErrorInfo) {
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	input := reqVo.Input
	sourceChannel := reqVo.SourceChannel

	issueBo, err := domain.GetIssueBo(orgId, input.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = domain.AuthIssue(orgId, currentUserId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModifyStatus)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	parentId := issueBo.ParentId
	if parentId > 0{
		//验证父任务状态是否是已完成，父任务已完成不允许编辑子任务状态
		parentIssueBo, err := domain.GetIssueBo(orgId, parentId)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		parentStatus, err1 := processfacade.GetProcessStatusByCategoryRelaxed(orgId, parentIssueBo.Status, consts.ProcessStatusCategoryIssue)
		if err1 != nil {
			log.Error(err1)
			return nil, err1
		}

		if parentStatus.StatusType == consts.ProcessStatusTypeCompleted{
			log.Error("父任务已完成不允许编辑子任务状态")
			return nil, errs.CantUpdateStatusWhenParentIssueIsCompleted
		}

	}

	needModifyChildStatus := 2
	if input.NeedModifyChildStatus != nil {
		needModifyChildStatus = *input.NeedModifyChildStatus
	}

	var err1 errs.SystemErrorInfo = nil
	if input.NextStatusID != nil && *input.NextStatusID > 0 {
		err1 = domain.UpdateIssueStatus(*issueBo, currentUserId, *input.NextStatusID, needModifyChildStatus, sourceChannel)
	} else if input.NextStatusType != nil && *input.NextStatusType > 0 {
		err1 = domain.UpdateIssueStatusByStatusType(*issueBo, currentUserId, *input.NextStatusType, needModifyChildStatus, sourceChannel)
	} else {
		log.Error("要更新的状态无效")
		return nil, errs.BuildSystemErrorInfo(errs.IssueStatusUpdateError)
	}

	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	result := &vo.Issue{}
	copyErr := copyer.Copy(issueBo, result)
	if copyErr != nil {
		log.Errorf("copyer.Copy(): %q\n", copyErr)
	}

	asyn.Execute(func() {
		PushModifyIssueNotice(issueBo.OrgId, issueBo.ProjectId, issueBo.Id, currentUserId)
	})
	return result, nil
}
