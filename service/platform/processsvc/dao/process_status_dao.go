package dao

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

func InsertProcessStatus(po po.PpmPrsProcessStatus, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("priority dao Insert err %v", err)
	}
	return nil
}

//func InsertProcessStatusBatch(pos []po.PpmPrsProcessStatus, tx ...sqlbuilder.Tx) error {
//	var err error = nil
//	if tx != nil && len(tx) > 0 {
//		batch := tx[0].InsertInto(consts.TableProcessStatus).Batch(len(pos))
//		//go func() {
//		//	defer batch.Done()
//		//	for i := range pos {
//		//		batch.Values(pos[i])
//		//	}
//		//}()
//		//go processStatusBatchDone(pos,batch)
//
//		go mysql.BatchDone(slice.ToSlice(pos), batch)
//
//		err = batch.Wait()
//		if err != nil {
//			return err
//		}
//	} else {
//		conn, err := mysql.GetConnect()
//		defer func() {
//			if err := conn.Close(); err != nil {
//				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
//			}
//		}()
//		if err != nil {
//			return err
//		}
//		batch := conn.InsertInto(consts.TableProcessStatus).Batch(len(pos))
//		//go func() {
//		//	defer batch.Done()
//		//	for i := range pos {
//		//		batch.Values(pos[i])
//		//	}
//		//}()
//		//go processStatusBatchDone(pos,batch)
//		go mysql.BatchDone(slice.ToSlice(pos), batch)
//
//		err = batch.Wait()
//		if err != nil {
//			return err
//		}
//	}
//	if err != nil {
//		log.Errorf("priority dao InsertBatch err %v", err)
//	}
//	return nil
//}

func InsertProcessStatusBatch(pos []po.PpmPrsProcessStatus, tx ...sqlbuilder.Tx) error {
	var err error = nil

	isTx := tx != nil && len(tx) > 0

	var batch *sqlbuilder.BatchInserter

	if !isTx {
		//没有事务
		conn, err := mysql.GetConnect()
		defer func() {
			if conn != nil {
				if err := conn.Close(); err != nil {
					logger.GetDefaultLogger().Info(strs.ObjectToString(err))
				}
			}
		}()

		if err != nil {
			return err
		}

		batch = conn.InsertInto(consts.TableProcessStatus).Batch(len(pos))
	}

	if batch == nil {
		batch = tx[0].InsertInto(consts.TableProcessStatus).Batch(len(pos))
	}

	go func() {
		defer batch.Done()
		for i := range pos {
			batch.Values(pos[i])
		}
	}()

	err = batch.Wait()
	if err != nil {
		log.Errorf("Iteration dao InsertBatch err %v", err)
		return err
	}
	return nil
}

func UpdateProcessStatus(po po.PpmPrsProcessStatus, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("priority dao Update err %v", err)
	}
	return err
}

func UpdateProcessStatusById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateProcessStatusByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateProcessStatusByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateProcessStatusByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateProcessStatusByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableProcessStatus, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableProcessStatus, cond, upd)
	}
	if err != nil {
		log.Errorf("priority dao Update err %v", err)
	}
	return mod, err
}

func DeleteProcessStatusById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateProcessStatusByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteProcessStatusByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateProcessStatusByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectProcessStatusById(id int64) (*po.PpmPrsProcessStatus, error) {
	po := &po.PpmPrsProcessStatus{}
	err := mysql.SelectById(consts.TableProcessStatus, id, po)
	if err != nil {
		log.Errorf("priority dao SelectById err %v", err)
	}
	return po, err
}

func SelectProcessStatusByIdAndOrg(id int64, orgId int64) (*po.PpmPrsProcessStatus, error) {
	po := &po.PpmPrsProcessStatus{}
	err := mysql.SelectOneByCond(consts.TableProcessStatus, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("priority dao Select err %v", err)
	}
	return po, err
}

func SelectProcessStatus(cond db.Cond) (*[]po.PpmPrsProcessStatus, error) {
	pos := &[]po.PpmPrsProcessStatus{}
	err := mysql.SelectAllByCond(consts.TableProcessStatus, cond, pos)
	if err != nil {
		log.Errorf("priority dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneProcessStatus(cond db.Cond) (*po.PpmPrsProcessStatus, error) {
	po := &po.PpmPrsProcessStatus{}
	err := mysql.SelectOneByCond(consts.TableProcessStatus, cond, po)
	if err != nil {
		log.Errorf("priority dao Select err %v", err)
	}
	return po, err
}

func SelectProcessStatusByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmPrsProcessStatus, uint64, error) {
	pos := &[]po.PpmPrsProcessStatus{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableProcessStatus, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("priority dao SelectPage err %v", err)
	}
	return pos, total, err
}
