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

func SourceExist(orgId, sourceId int64) bool {
	isExist, err := mysql.IsExistByCond(consts.TableIssueSource, db.Cond{
		consts.TcId:       sourceId,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
	})
	if err != nil {
		return false
	}

	return isExist
}

func GetIssueSourceInfo(orgId int64, sourceIds []int64) ([]bo.IssueSourceBo, errs.SystemErrorInfo) {
	info := &[]po.PpmPrsIssueSource{}
	err := mysql.SelectAllByCond(consts.TableIssueSource, db.Cond{
		consts.TcId:       db.In(sourceIds),
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
	}, info)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	infoBo := &[]bo.IssueSourceBo{}
	copyErr := copyer.Copy(info, infoBo)
	if copyErr != nil {
		return nil, errs.ObjectCopyError
	}

	return *infoBo, nil
}
