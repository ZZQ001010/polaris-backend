package domain

import (
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/schedule/consts"
)

func SetIssuePlanEndTimeLastScanTime(lastScanTime string) errs.SystemErrorInfo{
	err := cache.Set(consts.CacheIssuePlanEndTimeLastScanTime, lastScanTime)
	if err != nil{
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func GetIssuePlanEndTimeLastScanTime() (string, errs.SystemErrorInfo){
	value, err := cache.Get(consts.CacheIssuePlanEndTimeLastScanTime)
	if err != nil{
		log.Error(err)
		return "", errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return value, nil
}