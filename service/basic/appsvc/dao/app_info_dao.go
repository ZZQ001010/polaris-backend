package basedao

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/basic/appsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var (
	log = logger.GetDefaultLogger()
)

func InsertAppInfo(po po.PpmBasAppInfo, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransInsert(tx[0], &po)
	} else {
		err = mysql.Insert(&po)
	}
	if err != nil {
		log.Errorf("AppInfo dao Insert err %v", err)
	}
	return nil
}

func InsertAppInfoBatch(pos []po.PpmBasAppInfo, tx ...sqlbuilder.Tx) error {
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
		batch = conn.InsertInto(consts.TableAppInfo).Batch(len(pos))
	}
	if batch == nil {
		batch = tx[0].InsertInto(consts.TableAppInfo).Batch(len(pos))
	}
	go func() {
		defer batch.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("捕获到的错误：%s", r)
			}
		}()
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

func UpdateAppInfo(po po.PpmBasAppInfo, tx ...sqlbuilder.Tx) error {
	var err error = nil
	if tx != nil && len(tx) > 0 {
		err = mysql.TransUpdate(tx[0], &po)
	} else {
		err = mysql.Update(&po)
	}
	if err != nil {
		log.Errorf("AppInfo dao Update err %v", err)
	}
	return err
}

func UpdateAppInfoById(id int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateAppInfoByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateAppInfoByOrg(id int64, orgId int64, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	return UpdateAppInfoByCond(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func UpdateAppInfoByCond(cond db.Cond, upd mysql.Upd, tx ...sqlbuilder.Tx) (int64, error) {
	var mod int64 = 0
	var err error = nil
	if tx != nil && len(tx) > 0 {
		mod, err = mysql.TransUpdateSmartWithCond(tx[0], consts.TableAppInfo, cond, upd)
	} else {
		mod, err = mysql.UpdateSmartWithCond(consts.TableAppInfo, cond, upd)
	}
	if err != nil {
		log.Errorf("AppInfo dao Update err %v", err)
	}
	return mod, err
}

func DeleteAppInfoById(id int64, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateAppInfoByCond(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func DeleteAppInfoByCode(code string, operatorId int64, tx ...sqlbuilder.Tx) (int64, error) {
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if operatorId > 0 {
		upd[consts.TcUpdator] = operatorId
	}
	return UpdateAppInfoByCond(db.Cond{
		consts.TcCode:     code,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd, tx...)
}

func SelectAppInfoById(id int64) (*po.PpmBasAppInfo, error) {
	po := &po.PpmBasAppInfo{}
	err := mysql.SelectById(consts.TableAppInfo, id, po)
	if err != nil {
		log.Errorf("AppInfo dao SelectById err %v", err)
	}
	return po, err
}

func SelectAppInfoByCode(code string) (*po.PpmBasAppInfo, error) {
	po := &po.PpmBasAppInfo{}
	err := mysql.SelectOneByCond(consts.TableAppInfo, db.Cond{
		consts.TcCode:     code,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, po)
	if err != nil {
		log.Errorf("AppInfo dao Select err %v", err)
	}
	return po, err
}

func SelectAppInfo(cond db.Cond) (*[]po.PpmBasAppInfo, error) {
	pos := &[]po.PpmBasAppInfo{}
	err := mysql.SelectAllByCond(consts.TableAppInfo, cond, pos)
	if err != nil {
		log.Errorf("AppInfo dao SelectList err %v", err)
	}
	return pos, err
}

func SelectOneAppInfo(cond db.Cond) (*po.PpmBasAppInfo, error) {
	po := &po.PpmBasAppInfo{}
	err := mysql.SelectOneByCond(consts.TableAppInfo, cond, po)
	if err != nil {
		log.Errorf("AppInfo dao Select err %v", err)
	}
	return po, err
}

func SelectAppInfoByPage(cond db.Cond, pageBo bo.PageBo) (*[]po.PpmBasAppInfo, uint64, error) {
	pos := &[]po.PpmBasAppInfo{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableAppInfo, cond, nil, pageBo.Page, pageBo.Size, pageBo.Order, pos)
	if err != nil {
		log.Errorf("AppInfo dao SelectPage err %v", err)
	}
	return pos, total, err
}
