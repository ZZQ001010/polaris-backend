package service

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/domain"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InitOrg(initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	return domain.InitOrg(initOrgBo, tx)
}

func GeneralInitOrg(initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	return domain.GeneralInitOrg(initOrgBo, tx)
}
