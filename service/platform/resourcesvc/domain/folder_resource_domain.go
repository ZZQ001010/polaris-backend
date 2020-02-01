package domain

import (
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/service/platform/resourcesvc/dao"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func DeleteMidTable(folderIds []int64, orgId, userId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	upd := mysql.Upd{}
	upd[consts.TcIsDelete] = consts.AppIsDeleted
	upd[consts.TcUpdator] = userId
	upd[consts.TcUpdateTime] = time.Now()
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcFolderId: db.In(folderIds),
		consts.TcOrgId:    orgId,
	}
	err := dao.UpdateMidTableByCond(cond, upd, tx...)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
