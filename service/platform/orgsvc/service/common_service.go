package service

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
)

func AuthOrgRole(orgId, userId int64, path string, operation string) errs.SystemErrorInfo {
	return rolefacade.AuthenticateRelaxed(orgId, userId, nil, nil, path, operation)
}