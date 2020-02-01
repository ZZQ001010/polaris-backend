package domain

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"upper.io/db.v3"
)

//statusType定义（1：未开始，2：进行中，3：已完成）
func IterationCondStatusAssembly(cond *db.Cond, orgId int64, statusType int) errs.SystemErrorInfo {
	var statusIds []int64 = nil

	if statusType == 1 {
		notStartedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIteration, consts.ProcessStatusTypeNotStarted)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		statusIds = *notStartedIds
	} else if statusType == 2 {
		processingIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIteration, consts.ProcessStatusTypeProcessing)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		statusIds = *processingIds
	} else {
		finishedId, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIteration, consts.ProcessStatusTypeCompleted)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		statusIds = *finishedId
	}
	(*cond)[consts.TcStatus] = db.In(statusIds)
	return nil
}
