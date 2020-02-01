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

func GetProcessStepList(orgId int64, processId int64) (*[]bo.ProcessStepBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProcessStepList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProcessIdConstName: processId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	processStepListJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	processStepListPo := &[]po.PpmPrsProcessStep{}
	processStepListBo := &[]bo.ProcessStepBo{}
	if processStepListJson != "" {
		err = json.FromJson(processStepListJson, processStepListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return processStepListBo, nil
	} else {
		err := mysql.SelectAllByCond(consts.TableProcessStep, db.Cond{
			consts.TcOrgId:     db.In([]int64{orgId, 0}),
			consts.TcProcessId: processId,
			consts.TcIsDelete:  consts.AppIsNoDelete,
			consts.TcStatus:    consts.AppStatusEnable,
		}, processStepListPo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.ProcessNotExist)
		}
		_ = copyer.Copy(processStepListPo, processStepListBo)
		processStepListJson, err := json.ToJson(processStepListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, processStepListJson, consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return processStepListBo, nil
	}
}
