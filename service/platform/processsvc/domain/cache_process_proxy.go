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
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/processsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/po"
	"upper.io/db.v3"
)

func GetProcess(orgId, id int64) (*bo.ProcessBo, errs.SystemErrorInfo) {
	list, err := GetProcessList(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, process := range *list {
		if process.Id == id {
			return &process, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.ProcessNotExist)
}

func GetProcessByLangCode(orgId int64, langCode string) (*bo.ProcessBo, errs.SystemErrorInfo) {
	processList, err1 := GetProcessList(orgId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}

	for _, process := range *processList {
		if process.LangCode == langCode {
			return &process, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.ProcessNotExist)
}

func GetProcessList(orgId int64) (*[]bo.ProcessBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProcessList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	processListJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	processListPo := &[]po.PpmPrsProcess{}
	processListBo := &[]bo.ProcessBo{}
	if processListJson != "" {
		err = json.FromJson(processListJson, processListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return processListBo, nil
	} else {
		err := mysql.SelectAllByCond(consts.TableProcess, db.Cond{
			consts.TcOrgId:    db.In([]int64{orgId, 0}),
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, processListPo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.ProcessNotExist)
		}
		_ = copyer.Copy(processListPo, processListBo)
		processListJson, err := json.ToJson(processListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, processListJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return processListBo, nil
	}
}
