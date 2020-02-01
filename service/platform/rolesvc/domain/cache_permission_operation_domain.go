package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	bo2 "github.com/galaxy-book/polaris-backend/common/model/bo"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/rolesvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/po"
	"github.com/prometheus/common/log"
	"upper.io/db.v3"
)

func GetPermissionOperationList() (*[]bo2.PermissionOperationBo, errs.SystemErrorInfo) {
	key := sconsts.CachePermissionOperationList
	listJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	bo := &[]bo2.PermissionOperationBo{}
	if listJson != "" {
		err := json.FromJson(listJson, bo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	} else {
		po := &[]po.PpmRolPermissionOperation{}
		selectErr := mysql.SelectAllByCond(consts.TablePermissionOperation, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcStatus:   consts.AppStatusEnable,
		}, po)
		if selectErr != nil {
			log.Error(selectErr)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
		}
		_ = copyer.Copy(po, bo)
		listJson, err = json.ToJson(bo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, listJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
	}

	return bo, nil
}

func GetPermissionOperationListByPermissionId(permissionId int64) ([]bo2.PermissionOperationBo, errs.SystemErrorInfo) {
	res := []bo2.PermissionOperationBo{}

	list, err := GetPermissionOperationList()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, v := range *list {
		//过滤所有读的操作项，默认拥有
		if v.IsShow == consts.AppShowEnable && v.PermissionId == permissionId && v.OperationCodes != consts.RoleOperationView {
			res = append(res, v)
		}
	}

	return res, nil
}
