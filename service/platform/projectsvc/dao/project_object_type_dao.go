package dao

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetProjectObjectType(orgId int64) ([]po.PpmPrsProjectObjectType, error) {
	objectType := &[]po.PpmPrsProjectObjectType{}
	err := mysql.SelectAllByCond(consts.TableProjectObjectType, db.Cond{
		"org_id":    orgId,
		"is_delete": consts.AppIsNoDelete,
	}, objectType)
	if err != nil {
		return *objectType, err
	}

	return *objectType, nil
}

func InsertProjectObjectType(po po.PpmPrsProjectObjectType, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("ProjectObjectType dao Insert err %v", err)
	}
	return nil
}

//func InsertProjectObjectTypeBatch(pos []po.PpmPrsProjectObjectType, tx ...sqlbuilder.Tx) error {
//	var err error = nil
//	if tx != nil && len(tx) > 0 {
//		batch := tx[0].InsertInto(consts.TableProjectObjectType).Batch(len(pos))
//		//go func() {
//		//	defer batch.Done()
//		//	for i := range pos {
//		//		batch.Values(pos[i])
//		//	}
//		//}()
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
//		batch := conn.InsertInto(consts.TableProjectObjectType).Batch(len(pos))
//		//go func() {
//		//	defer batch.Done()
//		//	for i := range pos {
//		//		batch.Values(pos[i])
//		//	}
//		//}()
//		go mysql.BatchDone(slice.ToSlice(pos), batch)
//
//		err = batch.Wait()
//		if err != nil {
//			return err
//		}
//	}
//	if err != nil {
//		log.Errorf("ProjectObjectType dao InsertBatch err %v", err)
//	}
//	return nil
//}

func InsertProjectObjectTypeBatch(pos []po.PpmPrsProjectObjectType, tx ...sqlbuilder.Tx) error {
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

		batch = conn.InsertInto(consts.TableProjectObjectType).Batch(len(pos))
	}

	if batch == nil {
		batch = tx[0].InsertInto(consts.TableProjectObjectType).Batch(len(pos))
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

func UpdateProjectObjectType(po po.PpmPrsProjectObjectType, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("ProjectObjectType dao Update err %v", err)
	}
	return err
}

func UpdateProjectObjectTypeById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateProjectObjectTypeByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateProjectObjectTypeByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateProjectObjectTypeByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateProjectObjectTypeByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableProjectObjectType, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableProjectObjectType, cond, upd)
	}
	if err != nil {
		log.Errorf("ProjectObjectType dao Update err %v", err)
	}
	return mod, err
}

func DeleteProjectObjectTypeById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateProjectObjectTypeByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteProjectObjectTypeByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateProjectObjectTypeByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectProjectObjectTypeById(id int64) (*po.PpmPrsProjectObjectType, error) {
	po := &po.PpmPrsProjectObjectType{}
	err := mysql.SelectById(consts.TableProjectObjectType, id, po)
	if err != nil {
		log.Errorf("ProjectObjectType dao SelectById err %v", err)
	}
	return po, err
}

func SelectProjectObjectTypeByIdAndOrg(id int64, orgId int64) (*po.PpmPrsProjectObjectType, error) {
	po := &po.PpmPrsProjectObjectType{}
	err := mysql.SelectOneByCond(consts.TableProjectObjectType, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("ProjectObjectType dao Select err %v", err)
	}
	return po, err
}

func SelectProjectObjectType(cond db.Cond) (*[]po.PpmPrsProjectObjectType, error) {
	pos := &[]po.PpmPrsProjectObjectType{}
	err := mysql.SelectAllByCond(consts.TableProjectObjectType, cond, pos)
	if err != nil {
		log.Errorf("ProjectObjectType dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneProjectObjectType(cond db.Cond) (*po.PpmPrsProjectObjectType, error) {
	po := &po.PpmPrsProjectObjectType{}
	err := mysql.SelectOneByCond(consts.TableProjectObjectType, cond, po)
	if err != nil {
		log.Errorf("ProjectObjectType dao Select err %v", err)
	}
	return po, err
}

func SelectProjectObjectTypeByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmPrsProjectObjectType, uint64, error) {
	pos := &[]po.PpmPrsProjectObjectType{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableProjectObjectType, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("ProjectObjectType dao SelectPage err %v", err)
	}
	return pos, total, err
}
