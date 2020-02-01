package domain

import (
	"errors"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	appconsts "github.com/galaxy-book/polaris-backend/service/basic/appsvc/consts"
	basedao "github.com/galaxy-book/polaris-backend/service/basic/appsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/basic/appsvc/po"
	"github.com/opentracing/opentracing-go/log"
)

//func GetAppInfoBoList(page uint, size uint, cond db.Cond) (*[]bo.AppInfoBo, int64, errs.SystemErrorInfo) {
//	pos, total, err := basedao.SelectAppInfoByPage(cond, bo.PageBo{
//		Page:  int(page),
//		Size:  int(size),
//		Order: "",
//	})
//	if err != nil {
//		log.Error(err)
//		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	bos := &[]bo.AppInfoBo{}
//
//	copyErr := copyer.Copy(pos, bos)
//	if copyErr != nil {
//		log.Error(copyErr)
//		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
//	}
//	return bos, int64(total), nil
//}

func GetAppInfoBoNoCache(id int64) (*bo.AppInfoBo, errs.SystemErrorInfo) {
	po, err := basedao.SelectAppInfoById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.AppInfoBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func GetAppInfoBoByCode(code string) (*bo.AppInfoBo, errs.SystemErrorInfo) {

	if code == "" {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, errors.New(" code "))
	}

	cacheKey := appconsts.CacheAppInfoCodeKey + code

	appInfoJson, err := cache.Get(cacheKey)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if appInfoJson != "" {
		appInfoBo := &bo.AppInfoBo{}
		err = json.FromJson(appInfoJson, appInfoBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return appInfoBo, nil
	}

	po, err := basedao.SelectAppInfoByCode(code)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.AppInfoBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	appInfoJson, err = json.ToJson(bo)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
	}
	err = cache.SetEx(cacheKey, appInfoJson, consts.GetCacheBaseExpire())

	return bo, nil
}

func CreateAppInfo(bo *bo.AppInfoBo) errs.SystemErrorInfo {
	po := &po.PpmBasAppInfo{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := basedao.InsertAppInfo(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateAppInfo(bo *bo.AppInfoBo) errs.SystemErrorInfo {
	po := &po.PpmBasAppInfo{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err := basedao.UpdateAppInfo(*po)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	cacheKey := appconsts.CacheAppInfoCodeKey + bo.Code
	cache.Del(cacheKey)
	return nil
}

func DeleteAppInfo(bo *bo.AppInfoBo, operatorId int64) errs.SystemErrorInfo {
	_, err := basedao.DeleteAppInfoById(bo.Id, operatorId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}
