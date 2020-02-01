package service

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/domain"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

func AssignValueToField(processRes *map[string]int64, tx sqlbuilder.Tx, orgId int64) errs.SystemErrorInfo {
	return domain.AssignValueToField(processRes, tx, orgId)
}

func InitProcess(orgId int64) errs.SystemErrorInfo {
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		return domain.InitProcess(orgId, tx)
	})
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func GetProcessBo(cond db.Cond) (bo.ProcessBo, errs.SystemErrorInfo) {
	processBo := bo.ProcessBo{}
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		res, err := domain.GetProcessBo(cond, tx)
		processBo = res
		return err
	})
	if err != nil {
		log.Error(err)
		return processBo, errs.BuildSystemErrorInfo(errs.ProcessDomainError, err)
	}
	return processBo, nil
}

func GetProcessById(orgId, id int64) (*bo.ProcessBo, errs.SystemErrorInfo) {
	return domain.GetProcess(orgId, id)
}
