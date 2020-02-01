package domain

import (
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"upper.io/db.v3"
)

func IssueObjectTypeExist(orgId, typeId int64) bool {
	isExist, err := mysql.IsExistByCond(consts.TableIssueObjectType, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       typeId,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err != nil {
		return false
	}

	return isExist
}
