package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetIterationBoList(page uint, size uint, cond db.Cond) (*[]bo.IterationBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectIterationByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.IterationBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetIterationBo(id int64) (*bo.IterationBo, errs.SystemErrorInfo) {
	po, err := dao.SelectIterationById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.IterationBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func GetIterationBoByOrgId(id int64, orgId int64) (*bo.IterationBo, errs.SystemErrorInfo) {
	po, err := dao.SelectIterationByIdAndOrg(id, orgId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.IterationBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func CreateIteration(bo *bo.IterationBo) errs.SystemErrorInfo {
	po := &po.PpmPriIteration{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.InsertIteration(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateIteration(updateBo *bo.IterationUpdateBo) errs.SystemErrorInfo {
	_, err := dao.UpdateIterationById(updateBo.Id, updateBo.Upd)
	//TODO 处理动态

	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func DeleteIteration(bo *bo.IterationBo, operatorId int64) errs.SystemErrorInfo {
	err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除迭代
		_, err := dao.DeleteIterationById(bo.Id, operatorId, tx)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}

		//删除任务和迭代关联
		_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
			consts.TcOrgId:       bo.OrgId,
			consts.TcIterationId: bo.Id,
			consts.TcIsDelete:    consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcIterationId: 0,
			consts.TcUpdator:     operatorId,
		})
		if err2 != nil {
			log.Error(err2)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
		return nil
	})

	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}
	return nil
}
