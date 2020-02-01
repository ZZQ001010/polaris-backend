package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetProjectObjectTypeProcessByCond(projectObjectTypeId, orgId int64) (*[]bo.PpmPrsProjectObjectTypeProcessBo, errs.SystemErrorInfo) {

	pos := &[]po.PpmPrsProjectObjectTypeProcess{}

	err := mysql.SelectAllByCond(consts.TableProjectObjectTypeProcess, db.Cond{
		consts.TcIsDelete:            consts.AppIsNoDelete,
		consts.TcProjectObjectTypeId: projectObjectTypeId,
		consts.TcOrgId:               db.In([]int64{orgId, 0}),
	}, pos)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	bos := &[]bo.PpmPrsProjectObjectTypeProcessBo{}

	err = copyer.Copy(pos, bos)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}
	return bos, nil

}
