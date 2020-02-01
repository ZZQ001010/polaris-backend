package service

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
)

func UpdateIssueSort(reqVo projectvo.UpdateIssueSortReqVo) (*vo.Void, errs.SystemErrorInfo) {
	currentUserId := reqVo.UserId
	orgId := reqVo.OrgId
	reqInput := reqVo.Input
	issueId := reqInput.ID
	beforeId := reqInput.BeforeID
	afterId := reqInput.AfterID

	issueBo, err := domain.GetIssueBo(orgId, issueId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	err = domain.AuthIssue(orgId, currentUserId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	if reqInput.ProjectObjectTypeID != nil {
		//更新项目对象类型(相比sort优先)
		err1 := domain.UpdateIssueProjectObjectType(orgId, currentUserId, *issueBo, *reqInput.ProjectObjectTypeID)

		if err1 != nil {
			log.Error(err1)
			return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
		}
	}

	if beforeId != nil || afterId != nil{
		err = domain.UpdateIssueSort(*issueBo, currentUserId, beforeId, afterId)
		if err != nil{
			log.Error(err)
			return nil, err
		}
	}else{
		log.Errorf("更新任务 %d sort时，beforeId和afterId都为空， 不需要更新", issueId)
	}

	asyn.Execute(func() {
		PushModifyIssueNotice(issueBo.OrgId, issueBo.ProjectId, issueBo.Id, currentUserId)
	})
	return &vo.Void{ID:issueId}, nil
}
