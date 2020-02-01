package dao

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InsertIssueTag(po po.PpmPriIssueTag, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("IssueTag dao Insert err %v", err)
	}
	return nil
}

func InsertIssueTagBatch(pos []po.PpmPriIssueTag, tx ...sqlbuilder.Tx) error {
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
		batch = conn.InsertInto(consts.TableIssueTag).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableIssueTag).Batch(len(pos))
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

func UpdateIssueTag(po po.PpmPriIssueTag, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("IssueTag dao Update err %v", err)
	}
	return err
}

func UpdateIssueTagById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIssueTagByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIssueTagByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateIssueTagByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateIssueTagByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableIssueTag, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableIssueTag, cond, upd)
	}
	if err != nil {
		log.Errorf("IssueTag dao Update err %v", err)
	}
	return mod, err
}

func DeleteIssueTagById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIssueTagByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteIssueTagByOrg(id int64, orgId int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateIssueTagByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectIssueTagById(id int64) (*po.PpmPriIssueTag, error) {
	po := &po.PpmPriIssueTag{}
	err := mysql.SelectById(consts.TableIssueTag, id, po)
	if err != nil {
		log.Errorf("IssueTag dao SelectById err %v", err)
	}
	return po, err
}

func SelectIssueTagByIdAndOrg(id int64, orgId int64) (*po.PpmPriIssueTag, error) {
	po := &po.PpmPriIssueTag{}
	err := mysql.SelectOneByCond(consts.TableIssueTag, db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("IssueTag dao Select err %v", err)
	}
	return po, err
}

func SelectIssueTag(cond db.Cond) (*[]po.PpmPriIssueTag, error) {
	pos := &[]po.PpmPriIssueTag{}
	err := mysql.SelectAllByCond(consts.TableIssueTag, cond, pos)
	if err != nil {
		log.Errorf("IssueTag dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneIssueTag(cond db.Cond) (*po.PpmPriIssueTag, error) {
	po := &po.PpmPriIssueTag{}
	err := mysql.SelectOneByCond(consts.TableIssueTag, cond, po)
	if err != nil {
		log.Errorf("IssueTag dao Select err %v", err)
	}
	return po, err
}

func SelectIssueTagByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmPriIssueTag, uint64, error) {
	pos := &[]po.PpmPriIssueTag{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableIssueTag, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("IssueTag dao SelectPage err %v", err)
	}
	return pos, total, err
}
