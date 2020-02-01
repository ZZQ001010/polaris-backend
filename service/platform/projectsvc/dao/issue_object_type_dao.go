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

func InsertIssueObjectType(po po.PpmPrsIssueObjectType, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("IssueObjectType dao Insert err %v", err)
	}
	return nil
}

func InsertIssueObjectTypeBatch(pos []po.PpmPrsIssueObjectType, tx ...sqlbuilder.Tx) error {
	var err error = nil
	isTx := tx != nil && len(tx) > 0
	var batch *sqlbuilder.BatchInserter
	if !isTx {
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
		batch = conn.InsertInto(consts.TableIssueObjectType).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableIssueObjectType).Batch(len(pos))
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

func UpdateIssueObjectType(po po.PpmPrsIssueObjectType, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("IssueObjectType dao Update err %v", err)
	}
	return err
}

func UpdateIssueObjectTypeById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIssueObjectTypeByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIssueObjectTypeByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIssueObjectTypeByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIssueObjectTypeByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableIssueObjectType, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableIssueObjectType, cond, upd)
	}
	if err != nil {
		log.Errorf("IssueObjectType dao Update err %v", err)
	}
	return mod, err
}

func DeleteIssueObjectTypeById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIssueObjectTypeByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteIssueObjectTypeByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIssueObjectTypeByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectIssueObjectTypeById(id int64) (*po.PpmPrsIssueObjectType, error) {
	po := &po.PpmPrsIssueObjectType{}
	err := mysql.SelectById(consts.TableIssueObjectType, id, po)
	if err != nil {
		log.Errorf("IssueObjectType dao SelectById err %v", err)
	}
	return po, err
}

func SelectIssueObjectTypeByIdAndOrg(id int64, orgId int64) (*po.PpmPrsIssueObjectType, error) {
	po := &po.PpmPrsIssueObjectType{}
	err := mysql.SelectOneByCond(consts.TableIssueObjectType, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("IssueObjectType dao Select err %v", err)
	}
	return po, err
}

func SelectIssueObjectType(cond db.Cond) (*[]po.PpmPrsIssueObjectType, error) {
	pos := &[]po.PpmPrsIssueObjectType{}
	err := mysql.SelectAllByCond(consts.TableIssueObjectType, cond, pos)
	if err != nil {
		log.Errorf("IssueObjectType dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneIssueObjectType(cond db.Cond) (*po.PpmPrsIssueObjectType, error) {
	po := &po.PpmPrsIssueObjectType{}
	err := mysql.SelectOneByCond(consts.TableIssueObjectType, cond, po)
	if err != nil {
		log.Errorf("IssueObjectType dao Select err %v", err)
	}
	return po, err
}

func SelectIssueObjectTypeByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmPrsIssueObjectType, uint64, error) {
	pos := &[]po.PpmPrsIssueObjectType{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableIssueObjectType, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("IssueObjectType dao SelectPage err %v", err)
	}
	return pos, total, err
}
