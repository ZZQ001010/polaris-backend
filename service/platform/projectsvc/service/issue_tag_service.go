package service

import (
	"fmt"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
)

func CreateIssueRelationTags(reqVo projectvo.CreateIssueRelationTagsReqVo) errs.SystemErrorInfo {
	input := reqVo.Input
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId

	issueId := input.ID

	issueBo, err := domain.GetIssueBo(orgId, issueId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.IllegalityIssue)
	}

	err = domain.AuthIssue(orgId, currentUserId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	newUUID := uuid.NewUuid()
	lockKey := fmt.Sprintf("%s%d", consts.IssueRelateOperationLock, issueId)
	if issueBo.ParentId != 0 {
		lockKey = fmt.Sprintf("%s%d", consts.IssueRelateOperationLock, issueBo.ParentId)
	}
	suc, lockErr := cache.TryGetDistributedLock(lockKey, newUUID)
	if lockErr != nil{
		log.Error(lockErr)
		return errs.TryDistributedLockError
	}
	if suc{
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, newUUID); err != nil{
				log.Error(err)
			}
		}()
	}else{
		//未获取到锁，直接响应错误信息
		return errs.IssueRelateTagFail
	}

	err1 := domain.IssueRelateTags(orgId, issueBo.ProjectId, issueId, currentUserId, ConvertIssueTagReqVoToBo(input.AddTags), ConvertIssueTagReqVoToBo(input.DelTags))
	if err1 != nil{
		log.Error(err1)
		return err1
	}

	asyn.Execute(func() {
		PushModifyIssueNotice(issueBo.OrgId, issueBo.ProjectId, issueBo.Id, currentUserId)
	})
	return nil
}