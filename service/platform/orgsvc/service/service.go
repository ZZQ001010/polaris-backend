package service

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/orgsvc/consts"
	"github.com/gin-gonic/gin"
)

var log = *logger.GetDefaultLogger()

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func GetCtxParameters(ctx context.Context, key string) (string, error) {
	gc, err := GinContextFromContext(ctx)
	if err != nil {
		return "", err
	}
	v := gc.GetString(key)
	return v, nil
}

func GetCurrentUser(ctx context.Context) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	return GetCurrentUserWithCond(ctx, true)
}

func GetCurrentUserWithoutOrgVerify(ctx context.Context) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	return GetCurrentUserWithCond(ctx, false)
}

func GetCurrentUserWithCond(ctx context.Context, orgVerify bool) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	token, err := GetCtxParameters(ctx, consts.AppHeaderTokenName)

	if err != nil || token == "" {
		return nil, errs.BuildSystemErrorInfo(errs.TokenNotExist)
	} else {

		redisJson, _ := json.ToJson(config.GetRedisConfig())

		fmt.Println("redis配置", redisJson)

		cacheUserInfoJson, err := cache.Get(sconsts.CacheUserToken + token)
		if err != nil {
			logger.GetDefaultLogger().Error(strs.ObjectToString(err))
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		if cacheUserInfoJson == "" {
			log.Error("token失效")
			return nil, errs.BuildSystemErrorInfo(errs.TokenExpires)
		}
		cacheUserInfo := &bo.CacheUserInfoBo{}
		err = json.FromJson(cacheUserInfoJson, cacheUserInfo)
		_, _ = cache.Expire(sconsts.CacheUserToken+token, consts.CacheUserTokenExpire)
		if err != nil {
			logger.GetDefaultLogger().Error(strs.ObjectToString(err))
			return nil, errs.BuildSystemErrorInfo(errs.TokenExpires)
		}
		//判断用户组织状态
		if cacheUserInfo.OrgId != 0 && orgVerify{
			baseUserInfo, err := GetBaseUserInfo("", cacheUserInfo.OrgId, cacheUserInfo.UserId)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			err = baseUserInfoOrgStatusCheck(*baseUserInfo)
			if err != nil{
				log.Error(err)
				return nil, err
			}
		}
		return cacheUserInfo, nil
	}
}


//用户信息所在组织状态监测
func baseUserInfoOrgStatusCheck(baseUserInfo bo.BaseUserInfoBo) errs.SystemErrorInfo{
	if baseUserInfo.OrgUserStatus != consts.AppStatusEnable {
		return errs.OrgUserUnabled
	}
	if baseUserInfo.OrgUserCheckStatus != consts.AppCheckStatusSuccess{
		return errs.OrgUserCheckStatusUnabled
	}
	if baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted{
		return errs.OrgUserDeleted
	}
	return nil
}


func UpdateCacheUserInfoOrgId(token string, orgId int64) errs.SystemErrorInfo {
	cacheUserInfoJson, err := cache.Get(sconsts.CacheUserToken + token)
	if err != nil {
		logger.GetDefaultLogger().Error(strs.ObjectToString(err))
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if cacheUserInfoJson == "" {
		return errs.BuildSystemErrorInfo(errs.TokenExpires)
	}
	cacheUserInfo := &bo.CacheUserInfoBo{}
	_ = json.FromJson(cacheUserInfoJson, cacheUserInfo)

	//更新缓存用户的orgId
	cacheUserInfo.OrgId = orgId
	cacheUserInfoJson = json.ToJsonIgnoreError(cacheUserInfo)
	err = cache.SetEx(sconsts.CacheUserToken+token, cacheUserInfoJson, consts.CacheUserTokenExpire)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}
