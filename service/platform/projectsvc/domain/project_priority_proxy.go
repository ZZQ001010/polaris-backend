package domain

import (
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/projectsvc/consts"
)

func DeletePriorityListCache(orgId int64) errs.SystemErrorInfo {

	key, err := util.ParseCacheKey(sconsts.CachePriorityList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})

	if err != nil {
		log.Error(err)
		return err
	}

	_, err1 := cache.Del(key)

	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}
