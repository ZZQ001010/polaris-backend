package basedao

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/websitesvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertContact(po po.PpmWstContact, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("Contact dao Insert err %v", err)
	}
	return nil
}

func InsertContactBatch(pos []po.PpmWstContact, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
		conn, err := mysql.GetConnect()
		defer func() {
			if conn != nil {
				if err := conn.Close(); err != nil {
					logger.GetDefaultLogger().Info(err)
				}
			}
		}()
		if err != nil {
			return err
		}
		batch = conn.InsertInto(consts.TableContact).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableContact).Batch(len(pos))
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

func UpdateContact(po po.PpmWstContact, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("Contact dao Update err %v", err)
	}
	return err
}

func UpdateContactById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateContactByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateContactByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateContactByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateContactByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableContact, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableContact, cond, upd)
	}
	if err != nil {
		log.Errorf("Contact dao Update err %v", err)
	}
	return mod, err
}

func DeleteContactById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateContactByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteContactByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateContactByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectContactById(id int64) (*po.PpmWstContact, error) {
	po := &po.PpmWstContact{}
	err := mysql.SelectById(consts.TableContact, id, po)
	if err != nil {
		log.Errorf("Contact dao SelectById err %v", err)
	}
	return po, err
}

func SelectContactByIdAndOrg(id int64, orgId int64) (*po.PpmWstContact, error) {
	po := &po.PpmWstContact{}
	err := mysql.SelectOneByCond(consts.TableContact, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("Contact dao Select err %v", err)
	}
	return po, err
}

func SelectContact(cond db.Cond) (*[]po.PpmWstContact, error) {
	pos := &[]po.PpmWstContact{}
	err := mysql.SelectAllByCond(consts.TableContact, cond, pos)
	if err != nil {
		log.Errorf("Contact dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneContact(cond db.Cond) (*po.PpmWstContact, error) {
	po := &po.PpmWstContact{}
	err := mysql.SelectOneByCond(consts.TableContact, cond, po)
	if err != nil {
		log.Errorf("Contact dao Select err %v", err)
	}
	return po, err
}

func SelectCountContact(cond db.Cond) (int64, error) {
	total, err := mysql.SelectCountByCond(consts.TableContact, cond)
	if err != nil {
		log.Errorf("Contact dao Select err %v", err)
	}
	return int64(total), err
}

func SelectContactByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmWstContact, uint64, error) {
	pos := &[]po.PpmWstContact{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableContact, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("Contact dao SelectPage err %v", err)
	}
	return pos, total, err
}
