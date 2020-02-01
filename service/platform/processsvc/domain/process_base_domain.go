package domain

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetProcessBo(cond db.Cond, tx sqlbuilder.Tx) (bo.ProcessBo, errs.SystemErrorInfo) {
	//获取默认项目流程id
	processInfo := &po.PpmPrsProcess{}
	processInfoBo := &bo.ProcessBo{}
	err := tx.Collection(processInfo.TableName()).Find(cond).One(processInfo)
	if err != nil {
		return *processInfoBo, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	copyErr := copyer.Copy(processInfo, processInfoBo)
	if copyErr != nil {
		log.Error(copyErr)
		return *processInfoBo, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return *processInfoBo, nil
}

func InitProcess(orgId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	count, err := mysql.SelectCountByCond((&po.PpmPrsProcess{}).TableName(), db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if count != 0 {
		logger.GetDefaultLogger().Infof("默认流程已初始化")
		return nil
	}

	processInsert := []interface{}{}
	for _, v := range po.Process {
		id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProcess)
		if err != nil {
			log.Error(err)
			return err
		}
		temp := v
		temp.OrgId = orgId
		temp.Id = id
		processInsert = append(processInsert, temp)
	}
	err = mysql.TransBatchInsert(tx, &po.PpmPrsProcess{}, processInsert)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

func AssignValueToField(processRes *map[string]int64, tx sqlbuilder.Tx, orgId int64) errs.SystemErrorInfo {
	process := &[]po.PpmPrsProcess{}
	err := tx.Select(consts.TcId, consts.TcLangCode).From((&po.PpmPrsProcess{}).TableName()).Where(db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}).All(process)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if len(*process) == 0 {
		logger.GetDefaultLogger().Infof("初始化项目类型与项目对象类型关联：无对应默认流程")
		return errs.BuildSystemErrorInfo(errs.ProcessNotExist)
	}

	for _, v := range *process {
		switch v.LangCode {
		case po.ProcessDefaultTestTask.LangCode:
			(*processRes)["test_task"] = v.Id
		case po.ProcessDefaultBug.LangCode:
			(*processRes)["bug"] = v.Id
		case po.ProcessDefaultTask.LangCode:
			(*processRes)["task"] = v.Id
		case po.ProcessDefaultDemand.LangCode:
			(*processRes)["demand"] = v.Id
		case po.ProcessDefaultFeature.LangCode:
			(*processRes)["feature"] = v.Id
		case po.ProcessDefaultIteration.LangCode:
			(*processRes)["iteration"] = v.Id
		case po.ProcessDefaultAgileTask.LangCode:
			(*processRes)["agile_task"] = v.Id
		}
	}

	return nil
}
