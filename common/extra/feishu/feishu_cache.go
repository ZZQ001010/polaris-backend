package feishu

import (
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/date"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/feishu-sdk-golang/sdk"
	"time"
)

//60s * 110 -> 110分钟
const CacheInvalidSeconds = 60 * 110

func ResendAppTicket(){
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil{
		log.Error("飞书配置为空")
		return
	}
	resp, err := sdk.AppTicketResend(fsConfig.AppId, fsConfig.AppSecret)
	if err != nil{
		log.Error(err)
		return
	}
	log.Infof("app_ticket 重新发送请求响应 %s", json.ToJsonIgnoreError(resp))
}

func GetAppTicket() (string, errs.SystemErrorInfo){
	cacheJson, err := cache.Get(consts.CacheFeiShuAppTicket)
	log.Infof("飞书AppTicket: %s", cacheJson)
	if err != nil{
		log.Error(err)
		return "", errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if cacheJson == ""{
		log.Error("app ticket为空")
		ResendAppTicket()
		return "", errs.BuildSystemErrorInfo(errs.FeiShuAppTicketNotExistError)
	}
	cacheBo := &bo.FeiShuAppTicketCacheBo{}
	_ = json.FromJson(cacheJson, cacheBo)
	return cacheBo.AppTicket, nil
}

func SetAppTicket(appId, appTicket string) errs.SystemErrorInfo{
	cacheJson := json.ToJsonIgnoreError(bo.FeiShuAppTicketCacheBo{
		AppId: appId,
		AppTicket: appTicket,
		LastUpdateTime: date.Format(time.Now()),
	})

	err := cache.Set(consts.CacheFeiShuAppTicket, cacheJson)
	if err != nil{
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func GetAppAccessToken() (*bo.FeiShuAppAccessTokenCacheBo, errs.SystemErrorInfo){
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil{
		log.Error("飞书配置为空")
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuConfigNotExistError)
	}
	appId := fsConfig.AppId
	appSecret := fsConfig.AppSecret
	key := getAppAccessTokenCacheKey(appId, appSecret)

	accessTokenBo, err := GetAppInfoFromCache(key)
	if err != nil{
		//防止缓存击穿
		uuid := uuid.NewUuid()
		suc, err := cache.TryGetDistributedLock(consts.FeiShuGetAppAccessTokenLockKey, uuid)
		if err != nil{
			log.Error("获取锁异常")
			return nil, errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
		}
		if suc{
			//释放锁
			defer func() {
				if _, e := cache.ReleaseDistributedLock(consts.FeiShuGetAppAccessTokenLockKey, uuid); e != nil{
					log.Error(e)
				}
			}()

			//二次判断
			accessTokenBo, err := GetAppInfoFromCache(key)
			if err == nil{
				return accessTokenBo, nil
			}

			//重新获取
			appTicket, err := GetAppTicket()
			if err != nil{
				log.Error(err)
				return nil, err
			}

			appInfoResp, fsErr := sdk.GetAppAccessToken(fsConfig.AppId, fsConfig.AppSecret, appTicket)
			if fsErr != nil{
				log.Error(fsErr)
				return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}
			if appInfoResp.Code != 0{
				log.Error(appInfoResp.Msg)
				if appInfoResp.Code == 99991666{
					ResendAppTicket()
				}
				return nil, errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, appInfoResp.Msg)
			}

			accessTokenBo = &bo.FeiShuAppAccessTokenCacheBo{}
			accessTokenBo.AppAccessToken = appInfoResp.AppAccessToken
			accessTokenBo.LastUpdateTime = date.Format(time.Now())

			accessTokenJson := json.ToJsonIgnoreError(accessTokenBo)
			cacheErr := cache.SetEx(key, accessTokenJson, CacheInvalidSeconds)
			if cacheErr != nil{
				log.Error(cacheErr)
			}
			return accessTokenBo, nil
		}else{
			//抢占锁超时后也获取一下
			return GetAppInfoFromCache(key)
		}
	}
	return accessTokenBo, nil
}

func GetAppInfoFromCache(key string) (*bo.FeiShuAppAccessTokenCacheBo, errs.SystemErrorInfo){
	accessTokenBo := &bo.FeiShuAppAccessTokenCacheBo{}
	accessTokenJson, err := cache.Get(key)
	if err != nil{
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if accessTokenJson != ""{
		jsonErr := json.FromJson(accessTokenJson, accessTokenBo)
		if jsonErr != nil{
			log.Error(jsonErr)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return accessTokenBo, nil
	}else{
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuGetAppAccessTokenError)
	}
}

func GetTenantAccessTokenFromCache(key string) (*bo.FeiShuTenantAccessTokenCacheBo, errs.SystemErrorInfo){
	accessTokenBo := &bo.FeiShuTenantAccessTokenCacheBo{}
	accessTokenJson, err := cache.Get(key)
	if err != nil{
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if accessTokenJson != ""{
		jsonErr := json.FromJson(accessTokenJson, accessTokenBo)
		if jsonErr != nil{
			log.Error(jsonErr)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return accessTokenBo, nil
	}else{
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuGetTenantAccessTokenError)
	}
}

func GetTenantAccessToken(tenantKey string) (*bo.FeiShuTenantAccessTokenCacheBo, errs.SystemErrorInfo){
	key := getTenantAccessTokenCacheKey(tenantKey)

	accessTokenBo, err := GetTenantAccessTokenFromCache(key)
	if err != nil{
		//防止缓存击穿
		uuid := uuid.NewUuid()
		suc, err := cache.TryGetDistributedLock(consts.FeiShuGetTenantAccessTokenLockKey, uuid)
		if err != nil{
			log.Error("获取锁异常")
			return nil, errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
		}
		if suc {
			//释放锁
			defer func() {
				if _, e := cache.ReleaseDistributedLock(consts.FeiShuGetTenantAccessTokenLockKey, uuid); e != nil {
					log.Error(e)
				}
			}()

			accessTokenBo, err := GetTenantAccessTokenFromCache(key)
			if err == nil{
				return accessTokenBo, nil
			}

			appAccessTokenBo, err := GetAppAccessToken()
			if err != nil{
				log.Error(err)
				return nil, err
			}
			tenantResp, fsErr := sdk.GetTenantAccessToken(appAccessTokenBo.AppAccessToken, tenantKey)
			if fsErr != nil{
				log.Error(fsErr)
				return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}
			if tenantResp.Code != 0{
				log.Error(tenantResp.Msg)
				if tenantResp.Code == 99991666{
					ResendAppTicket()
				}
				return nil, errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, tenantResp.Msg)
			}

			accessTokenBo = &bo.FeiShuTenantAccessTokenCacheBo{}
			accessTokenBo.TenantAccessToken = tenantResp.TenantAccessToken
			accessTokenBo.TenantKey = tenantKey
			accessTokenBo.LastUpdateTime = date.Format(time.Now())

			accessTokenJson := json.ToJsonIgnoreError(accessTokenBo)
			cacheErr := cache.SetEx(key, accessTokenJson, CacheInvalidSeconds)
			if cacheErr != nil{
				log.Error(cacheErr)
			}

			return accessTokenBo, nil
		}else{
			//抢占锁超时后也获取一下
			return GetTenantAccessTokenFromCache(key)
		}

	}
	return accessTokenBo, nil
}

func ClearFsScopeCache(tenantKey string) errs.SystemErrorInfo{
	key := consts.CacheFeiShuScope + tenantKey
	_, err := cache.Del(key)
	if err != nil{
		log.Error(err)
		return errs.CacheProxyError
	}
	return nil
}

func GetScopeOpenIdsFromCache(tenantKey string) ([]string, errs.SystemErrorInfo){
	key := consts.CacheFeiShuScope + tenantKey

	value, cacheErr := cache.Get(key)
	if cacheErr != nil{
		log.Error(cacheErr)
		return nil, errs.CacheProxyError
	}
	if value != ""{
		result := &[]string{}
		_ = json.FromJson(value, result)
		return *result, nil
	}else{
		scopeOpenIds, err := GetScopeOpenIds(tenantKey)
		if err != nil{
			log.Error(err)
			return nil, err
		}

		value := json.ToJsonIgnoreError(scopeOpenIds)
		cacheErr = cache.SetEx(key, value, consts.GetCacheBaseExpire())
		if cacheErr != nil{
			log.Error(cacheErr)
			return nil, errs.CacheProxyError
		}

		return scopeOpenIds, nil
	}
}

func getAppAccessTokenCacheKey(appId, appSecret string) string{
	return consts.CacheFeiShuAppAccessToken + appId + ":" + appSecret
}

func getTenantAccessTokenCacheKey(tenant string) string{
	return consts.CacheFeiShuTenantAccessToken + tenant
}