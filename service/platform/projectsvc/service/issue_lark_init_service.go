package service

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"upper.io/db.v3/lib/sqlbuilder"
)

func LarkIssueInit(orgId int64, zhangsanId, lisiId, projectId, operatorId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	return domain.LarkIssueInit(orgId, zhangsanId, lisiId, projectId, operatorId, tx)
}
