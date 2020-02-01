package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/po"
	"upper.io/db.v3"
)

func GetProcessStatusBoList(page uint, size uint, cond db.Cond) (*[]bo.ProcessStatusBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectProcessStatusByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.ProcessStatusBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetProcessStatusBo(id int64) (*bo.ProcessStatusBo, errs.SystemErrorInfo) {
	po, err := dao.SelectProcessStatusById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.ProcessStatusBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func CreateProcessStatus(bo *bo.ProcessStatusBo) errs.SystemErrorInfo {
	po := &po.PpmPrsProcessStatus{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.InsertProcessStatus(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateProcessStatus(bo *bo.ProcessStatusBo) errs.SystemErrorInfo {
	po := &po.PpmPrsProcessStatus{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := dao.UpdateProcessStatus(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func DeleteProcessStatus(bo *bo.ProcessStatusBo, operatorId int64) errs.SystemErrorInfo {
	_, err := dao.DeleteProcessStatusById(bo.Id, operatorId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}
