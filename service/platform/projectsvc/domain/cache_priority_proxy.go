package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/projectsvc/consts"
	po2 "github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetPriorityList(orgId int64) (*[]bo.PriorityBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CachePriorityList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}
	priorityListJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	priorityListPo := &[]po2.PpmPrsPriority{}
	priorityListBo := &[]bo.PriorityBo{}
	if priorityListJson != "" {
		err = json.FromJson(priorityListJson, priorityListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return priorityListBo, nil
	} else {
		err := mysql.SelectAllByCond(consts.TablePriority, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, priorityListPo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_ = copyer.Copy(priorityListPo, priorityListBo)
		priorityListJson, err := json.ToJson(priorityListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, priorityListJson, consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return priorityListBo, nil
	}
}

func GetPriorityListByType(orgId int64, typ int) (*[]bo.PriorityBo, errs.SystemErrorInfo) {
	list, err := GetPriorityList(orgId)
	if err != nil {
		return nil, err
	}
	priorityList := &[]bo.PriorityBo{}
	for _, priority := range *list {
		if priority.Type == typ {
			*priorityList = append(*priorityList, priority)
		}
	}
	return priorityList, nil
}

func GetPriorityById(orgId int64, id int64) (*bo.PriorityBo, errs.SystemErrorInfo) {
	list, err := GetPriorityList(orgId)
	if err != nil {
		return nil, err
	}
	for _, priority := range *list {
		if priority.Id == id {
			return &priority, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.PriorityNotExist)
}
