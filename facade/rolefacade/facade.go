package rolefacade

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
)

var log = logger.GetDefaultLogger()

func AuthenticateRelaxed(orgId int64, userId int64, projectAuthInfo *bo.ProjectAuthBo, issueAuthInfo *bo.IssueAuthBo, path string, operation string) errs.SystemErrorInfo {
	respVo := Authenticate(rolevo.AuthenticateReqVo{
		OrgId:  orgId,
		UserId: userId,
		AuthInfoReqVo: rolevo.AuthenticateAuthInfoReqVo{
			ProjectAuthInfo: projectAuthInfo,
			IssueAuthInfo:   issueAuthInfo,
		},
		Path:      path,
		Operation: operation,
	})
	if respVo.Failure() {
		return respVo.Error()
	}

	return nil
}

func RoleUserRelationRelaxed(orgId, userId, roleId int64) errs.SystemErrorInfo {
	respVo := RoleUserRelation(rolevo.RoleUserRelationReqVo{
		OrgId:  orgId,
		UserId: userId,
		RoleId: roleId,
	})
	if respVo.Failure() {
		return respVo.Error()
	}

	return nil
}
