package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"upper.io/db.v3"
)

//获取未完成的的迭代列表
func GetNotCompletedIterationBoList(orgId int64, projectId int64) ([]bo.IterationBo, errs.SystemErrorInfo) {
	statusIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIteration, consts.ProcessStatusTypeCompleted)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}
	bos := &[]bo.IterationBo{}

	pos, err1 := dao.SelectIteration(db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
		consts.TcStatus:    db.NotIn(statusIds),
		consts.TcIsDelete:  consts.AppIsNoDelete,
	})
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}

	err2 := copyer.Copy(pos, bos)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *bos, nil
}
