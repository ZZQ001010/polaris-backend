package domain

import (
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func UpdateIterationStatus(iterationBo bo.IterationBo, nextStatusId, operationId int64) errs.SystemErrorInfo {
	orgId := iterationBo.OrgId
	iterationId := iterationBo.Id

	if iterationBo.Status == nextStatusId {
		log.Error("更新迭代状态-要更新的状态和当前状态一样")
		return errs.BuildSystemErrorInfo(errs.IterationStatusUpdateError)
	}

	//验证状态有效性
	nextStatus, err := processfacade.GetProcessStatusByCategoryRelaxed(orgId, nextStatusId, consts.ProcessStatusCategoryIteration)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	needInsertStat := false
	//判断迭代下是否有未完成的任务
	if nextStatus.StatusType == consts.ProcessStatusTypeCompleted {
		finishedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
		if err != nil {
			log.Errorf("proxies.GetProcessStatusId: %c\n", err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}

		cond := db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcIsDelete:    consts.AppIsNoDelete,
			consts.TcStatus:      db.NotIn(finishedIds),
			consts.TcIterationId: iterationId,
		}

		count, err2 := mysql.SelectCountByCond(consts.TableIssue, cond)
		if err2 != nil {
			log.Errorf("mysql.SelectAllByCond: %c\n", err2)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
		}

		if count > 0 {
			return errs.BuildSystemErrorInfo(errs.IterationExistingNotFinishedTask)
		}

		//是否创建统计基准数据
		needInsertStat = true
	}

	//处理db
	err3 := dealTx(needInsertStat, iterationBo, iterationId, orgId, nextStatusId)

	if err3 != nil {
		log.Error(err3)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err3)
	}
	return nil
}

func dealTx(needInsertStat bool, iterationBo bo.IterationBo, iterationId, orgId, nextStatusId int64) error {

	err3 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if needInsertStat {
			//创建统计基准数据
			err1 := AppendIterationStat(iterationBo, consts.BlankDate, tx)
			if err1 != nil {
				log.Error(err1)
				return errs.BuildSystemErrorInfo(errs.IterationDomainError, err1)
			}
		}
		_, err2 := dao.UpdateIterationByOrg(iterationId, orgId, mysql.Upd{
			consts.TcStatus: nextStatusId,
		}, tx)
		if err2 != nil {
			log.Error(err2)
			return errs.BuildSystemErrorInfo(errs.IterationStatusUpdateError)
		}
		return nil
	})

	return err3
}
