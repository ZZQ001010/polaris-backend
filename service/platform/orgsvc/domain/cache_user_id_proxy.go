package domain

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/orgsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"github.com/pkg/errors"
	"upper.io/db.v3"
)

func GetUserIdBatchByEmpId(sourceChannel string, orgId int64, empIds []string) ([]int64, errs.SystemErrorInfo) {
	keys := make([]interface{}, len(empIds))
	for i, empId := range empIds {
		key, _ := util.ParseCacheKey(sconsts.CacheOutUserIdRelationId, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:         orgId,
			consts.CacheKeySourceChannelConstName: sourceChannel,
			consts.CacheKeyOutUserIdConstName:     empId,
		})
		keys[i] = key
	}
	resultList := make([]string, 0)
	if len(keys) > 0{
		list, err := cache.MGet(keys...)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		resultList = list
	}
	userIds := make([]int64, 0)
	validEmpIds := make([]string, 0)
	for _, empInfoJson := range resultList {
		empIdInfo := &bo.UserEmpIdInfo{}
		err := json.FromJson(empInfoJson, empIdInfo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		userIds = append(userIds, empIdInfo.UserId)
		validEmpIds = append(validEmpIds, empIdInfo.EmpId)
	}
	//找不存在的
	if len(empIds) != len(validEmpIds) {
		for _, empId := range empIds {
			exist, _ := slice.Contain(validEmpIds, empId)
			if !exist {
				userId, err := GetUserIdByEmpId(sourceChannel, orgId, empId)
				if err != nil {
					log.Error(err)
					continue
				}
				userIds = append(userIds, userId)
			}
		}
	}
	return userIds, nil
}

func GetUserIdByEmpId(sourceChannel string, orgId int64, empId string) (int64, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheOutUserIdRelationId, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:         orgId,
		consts.CacheKeySourceChannelConstName: sourceChannel,
		consts.CacheKeyOutUserIdConstName:     empId,
	})
	if err5 != nil {
		log.Error(err5)
		return 0, err5
	}

	empInfoJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if empInfoJson != "" {
		empIdInfo := &bo.UserEmpIdInfo{}
		err := json.FromJson(empInfoJson, empIdInfo)
		if err != nil {
			log.Error(err)
			return 0, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return empIdInfo.UserId, nil
	} else {
		userOutInfo := &po.PpmOrgUserOutInfo{}
		err := mysql.SelectOneByCond(userOutInfo.TableName(), db.Cond{
			consts.TcOutUserId:  empId,
			consts.TcOrgId:         orgId,
			consts.TcSourceChannel: sourceChannel,
			consts.TcIsDelete:      consts.AppIsNoDelete,
			consts.TcStatus:        consts.AppStatusEnable,
		}, userOutInfo)
		if err != nil {
			log.Error(err)
			return 0, errs.BuildSystemErrorInfo(errs.UserNotExist, errors.New(" empId:"+empId))
		}
		empIdInfo := bo.UserEmpIdInfo{
			EmpId:  empId,
			UserId: userOutInfo.UserId,
		}
		err = cache.Set(key, json.ToJsonIgnoreError(empIdInfo))
		if err != nil {
			log.Error(err)
			return 0, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return userOutInfo.UserId, nil
	}
}

func GetDingTalkBaseUserInfoByEmpId(orgId int64, empId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	userId, err := GetUserIdByEmpId(consts.AppSourceChannelDingTalk, orgId, empId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return GetDingTalkBaseUserInfo(orgId, userId)
}

func GetBaseUserInfoByEmpId(sourceChannel string, orgId int64, empId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	userId, err := GetUserIdByEmpId(sourceChannel, orgId, empId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return GetBaseUserInfo(sourceChannel, orgId, userId)
}

func GetFeiShuBaseUserInfoByEmpId(orgId int64, empId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	userId, err := GetUserIdByEmpId(consts.AppSourceChannelFeiShu, orgId, empId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return GetFeiShuBaseUserInfo(orgId, userId)
}

func GetDingTalkBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return GetBaseUserInfo(consts.AppSourceChannelDingTalk, orgId, userId)
}

func GetFeiShuBaseUserInfo(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return GetBaseUserInfo(consts.AppSourceChannelFeiShu, orgId, userId)
}
