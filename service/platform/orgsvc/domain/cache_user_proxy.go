package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/orgsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

func GetBaseUserOutInfoBatch(sourceChannel string, orgId int64, userIds []int64) ([]bo.BaseUserOutInfoBo, errs.SystemErrorInfo) {
	keys := make([]interface{}, len(userIds))
	for i, userId := range userIds {
		key, _ := util.ParseCacheKey(sconsts.CacheBaseUserOutInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeySourceChannelConstName: sourceChannel,
			consts.CacheKeyUserIdConstName: userId,
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
	baseUserOutInfoList := make([]bo.BaseUserOutInfoBo, 0)
	validUserIds := map[int64]bool{}
	for _, userInfoJson := range resultList {
		userOutInfoBo := &bo.BaseUserOutInfoBo{}
		err := json.FromJson(userInfoJson, userOutInfoBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		baseUserOutInfoList = append(baseUserOutInfoList, *userOutInfoBo)
		validUserIds[userOutInfoBo.UserId] = true
	}

	missUserIds := make([]int64, 0)
	//找不存在的
	if len(userIds) != len(validUserIds) {
		for _, userId := range userIds {
			if _, ok := validUserIds[userId]; !ok{
				missUserIds = append(missUserIds, userId)
			}
		}
	}

	//批量查外部信息
	outInfos, userErr := GetBaseUserOutInfoByUserIds(sourceChannel, orgId, missUserIds)
	if userErr != nil{
		log.Error(userErr)
		return nil, userErr
	}

	if len(outInfos) > 0{
		baseUserOutInfoList = append(baseUserOutInfoList, outInfos...)
	}

	return baseUserOutInfoList, nil
}

func GetBaseUserInfoBatch(sourceChannel string, orgId int64, userIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	//去重
	userIds = slice.SliceUniqueInt64(userIds)

	keys := make([]interface{}, len(userIds))
	for i, userId := range userIds {
		key, _ := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:  orgId,
			consts.CacheKeyUserIdConstName: userId,
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
	baseUserInfoList := make([]bo.BaseUserInfoBo, 0)
	validUserIds := map[int64]bool{}
	for _, userInfoJson := range resultList {
		userInfoBo := &bo.BaseUserInfoBo{}
		err := json.FromJson(userInfoJson, userInfoBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		baseUserInfoList = append(baseUserInfoList, *userInfoBo)
		validUserIds[userInfoBo.UserId] = true
	}

	log.Infof("from cache %s", json.ToJsonIgnoreError(baseUserInfoList))
	missUserIds := make([]int64, 0)
	//找不存在的
	if len(userIds) != len(validUserIds) {
		for _, userId := range userIds {
			if _, ok := validUserIds[userId]; !ok{
				missUserIds = append(missUserIds, userId)
			}
		}
	}

	missUserInfos, userErr := getLocalBaseUserInfoBatch(orgId, missUserIds)
	if userErr != nil{
		log.Error(userErr)
		return nil, userErr
	}
	if len(missUserInfos) > 0{
		baseUserInfoList = append(baseUserInfoList, missUserInfos...)
	}

	if sourceChannel != ""{
		//获取用户外部信息
		baseUserOutInfos, err := GetBaseUserOutInfoBatch(sourceChannel, orgId, userIds)
		if err != nil{
			log.Error(err)
			return nil, err
		}
		outInfoMap := maps.NewMap("UserId", baseUserOutInfos)
		for i, _ := range baseUserInfoList{
			userInfo := baseUserInfoList[i]
			if outInfoInterface, ok := outInfoMap[userInfo.UserId]; ok{
				outInfo := outInfoInterface.(bo.BaseUserOutInfoBo)

				userInfo.OutUserId = outInfo.OutUserId
				userInfo.OutOrgId = outInfo.OutOrgId
				userInfo.HasOutInfo = outInfo.OutUserId != ""
				userInfo.HasOrgOutInfo = outInfo.OutOrgId != ""
			}
			baseUserInfoList[i] = userInfo
		}
	}
	return baseUserInfoList, nil
}

func ClearBaseUserInfo(orgId, userId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:         orgId,
		consts.CacheKeyUserIdConstName:        userId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

//批量清楚用户缓存信息
func ClearBaseUserInfoBatch(orgId int64, userIds []int64) errs.SystemErrorInfo {
	keys := make([]interface{}, 0)
	for _, userId := range userIds{
		key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:         orgId,
			consts.CacheKeyUserIdConstName:        userId,
		})
		if err5 != nil {
			log.Error(err5)
			return err5
		}
		keys = append(keys, key)
	}
	_, err := cache.Del(keys...)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

//sourceChannel可以为空
func GetBaseUserInfo(sourceChannel string, orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	if userId == 0{
		//系统创建
		return &bo.BaseUserInfoBo{
			OrgId: orgId,
			Name: "系统",
		}, nil
	}

	key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:         orgId,
		consts.CacheKeyUserIdConstName:        userId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	baseUserInfoJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	baseUserInfo := &bo.BaseUserInfoBo{}
	if baseUserInfoJson != "" {
		err := json.FromJson(baseUserInfoJson, baseUserInfo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
	} else {
		userInfo, errorInfo := getLocalBaseUserInfo(orgId, userId, sourceChannel, key)
		if errorInfo != nil {
			log.Error(errorInfo)
			return nil, errorInfo
		}
		baseUserInfo = userInfo
	}
	//这里不存缓存，动态获取
	if sourceChannel != ""{
		baseUserOutInfo, err := GetBaseUserOutInfo(sourceChannel, orgId, userId)
		if err != nil{
			log.Error(err)
			return nil, err
		}
		baseUserInfo.OutUserId = baseUserOutInfo.OutUserId
		baseUserInfo.OutOrgId = baseUserOutInfo.OutOrgId
		baseUserInfo.HasOutInfo = baseUserInfo.OutUserId != ""
		baseUserInfo.HasOrgOutInfo = baseUserInfo.OutOrgId != ""
	}

	return baseUserInfo, nil
}

func GetBaseUserOutInfo(sourceChannel string, orgId int64, userId int64) (*bo.BaseUserOutInfoBo, errs.SystemErrorInfo) {
	if userId == 0{
		//系统创建
		return &bo.BaseUserOutInfoBo{
			OrgId: orgId,
		}, nil
	}

	key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserOutInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:         orgId,
		consts.CacheKeySourceChannelConstName: sourceChannel,
		consts.CacheKeyUserIdConstName:        userId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}
	baseUserOutInfoJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if baseUserOutInfoJson != "" {
		baseUserOutInfo := &bo.BaseUserOutInfoBo{}
		err := json.FromJson(baseUserOutInfoJson, baseUserOutInfo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		return baseUserOutInfo, nil
	} else {
		//用户外部信息
		userOutInfo := &po.PpmOrgUserOutInfo{}
		_ = mysql.SelectOneByCond(consts.TableUserOutInfo, db.Cond{
			consts.TcIsDelete:      consts.AppIsNoDelete,
			consts.TcOrgId:         orgId,
			consts.TcSourceChannel: sourceChannel,
			consts.TcUserId:        userId,
		}, userOutInfo)

		//组织外部信息
		orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
		err = mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
			consts.TcIsDelete:      consts.AppIsNoDelete,
			consts.TcOrgId:         orgId,
			consts.TcSourceChannel: sourceChannel,
		}, orgOutInfo)

		outInfo := bo.BaseUserOutInfoBo{
			UserId: userId,
			OrgId: orgId,
			OutUserId: userOutInfo.OutUserId,
			OutOrgId: orgOutInfo.OutOrgId,
		}
		baseUserOutInfoJson, err := json.ToJson(outInfo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		err = cache.SetEx(key, baseUserOutInfoJson, consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return &outInfo, nil
	}
}

func GetBaseUserOutInfoByUserIds(sourceChannel string, orgId int64, userIds []int64) ([]bo.BaseUserOutInfoBo, errs.SystemErrorInfo) {
	log.Infof("批量获取用户外部信息 %d, %s", orgId, json.ToJsonIgnoreError(userIds))

	//用户外部信息
	userOutInfos := &[]po.PpmOrgUserOutInfo{}
	err := mysql.SelectAllByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcOrgId:         orgId,
		consts.TcSourceChannel: sourceChannel,
		consts.TcUserId:        db.In(userIds),
	}, userOutInfos)
	if err != nil{
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	//组织外部信息
	orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
	err = mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcOrgId:         orgId,
		consts.TcSourceChannel: sourceChannel,
	}, orgOutInfo)

	resultList := make([]bo.BaseUserOutInfoBo, 0)

	msetArgs := map[string]string{}
	keys := make([]string, 0)
	for _, userOutInfo := range *userOutInfos{
		key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserOutInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:         orgId,
			consts.CacheKeySourceChannelConstName: sourceChannel,
			consts.CacheKeyUserIdConstName:        userOutInfo.UserId,
		})
		if err5 != nil {
			log.Error(err5)
			return nil, err5
		}
		keys = append(keys, key)

		outInfo := bo.BaseUserOutInfoBo{
			UserId: userOutInfo.UserId,
			OrgId: orgId,
			OutUserId: userOutInfo.OutUserId,
			OutOrgId: orgOutInfo.OutOrgId,
		}

		resultList = append(resultList, outInfo)
		msetArgs[key] = json.ToJsonIgnoreError(outInfo)
	}

	if len(msetArgs) > 0{
		err = cache.MSet(msetArgs)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("捕获到的错误：%s", r)
			}
		}()

		for _, key := range keys{
			_, _ = cache.Expire(key, consts.GetCacheBaseUserInfoExpire())
		}
	}()

	return resultList, nil
}


//sourceChannel可以为空
func getLocalBaseUserInfo(orgId, userId int64, sourceChannel, key string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	user := &po.PpmOrgUser{}
	err := mysql.SelectById(user.TableName(), userId, user)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	newestUserOrganization, err1 := GetUserOrganizationNewestRelation(orgId, userId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}
	baseUserInfo := &bo.BaseUserInfoBo{
		UserId:        user.Id,
		Name:          user.Name,
		NamePy:		   user.NamePinyin,
		Avatar:        user.Avatar,
		OrgId:         orgId,
		OrgUserIsDelete: newestUserOrganization.IsDelete,
		OrgUserStatus: newestUserOrganization.Status,
		OrgUserCheckStatus: newestUserOrganization.CheckStatus,
	}

	baseUserInfoJson, err := json.ToJson(baseUserInfo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
	}
	err = cache.SetEx(key, baseUserInfoJson, consts.GetCacheBaseUserInfoExpire())
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	return baseUserInfo, nil
}

func getLocalBaseUserInfoBatch(orgId int64, userIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	log.Infof("批量获取用户信息 %d, %s", orgId, json.ToJsonIgnoreError(userIds))

	baseUserInfos := make([]bo.BaseUserInfoBo, 0)

	users := &[]po.PpmOrgUser{}
	err := mysql.SelectAllByCond(consts.TableUser, db.Cond{
		consts.TcId: db.In(userIds),
	}, users)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	//获取关联列表，要做去重
	userOrganizationPos := &[]po.PpmOrgUserOrganization{}
	_, selectErr := mysql.SelectAllByCondWithPageAndOrder(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcUserId: db.In(userIds),
	}, nil, 0, -1, "id asc", userOrganizationPos)
	if selectErr != nil{
		log.Error(selectErr)
		return nil, errs.MysqlOperateError
	}

	//id升序，保留最新: key: userId, value: po
	userOrgMap := map[int64]po.PpmOrgUserOrganization{}
	for _, userOrg := range *userOrganizationPos{
		userOrgMap[userOrg.UserId] = userOrg
	}

	for _, user := range *users{
		baseUserInfo := bo.BaseUserInfoBo{
			UserId:        user.Id,
			Name:          user.Name,
			NamePy:		   user.NamePinyin,
			Avatar:        user.Avatar,
			OrgId:         orgId,
		}

		if userOrg, ok := userOrgMap[user.Id]; ok{
			baseUserInfo.OrgUserIsDelete = userOrg.IsDelete
			baseUserInfo.OrgUserStatus = userOrg.Status
			baseUserInfo.OrgUserCheckStatus = userOrg.CheckStatus
		}

		baseUserInfos = append(baseUserInfos, baseUserInfo)
	}

	msetArgs := map[string]string{}
	keys := make([]string, 0)
	for _, baseUserInfo := range baseUserInfos{
		key, err5 := util.ParseCacheKey(sconsts.CacheBaseUserInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:         orgId,
			consts.CacheKeyUserIdConstName:        baseUserInfo.UserId,
		})
		if err5 != nil {
			log.Error(err5)
			return nil, err5
		}
		msetArgs[key] = json.ToJsonIgnoreError(baseUserInfo)
		keys = append(keys, key)
	}

	if len(msetArgs) > 0{
		err = cache.MSet(msetArgs)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("捕获到的错误：%s", r)
			}
		}()

		for _, key := range keys{
			_, _ = cache.Expire(key, consts.GetCacheBaseUserInfoExpire())
		}
	}()
	return baseUserInfos, nil
}

func GetUserConfigInfo(orgId int64, userId int64) (*bo.UserConfigBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheUserConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	userConfigJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	userConfigBo := &bo.UserConfigBo{}
	if userConfigJson != "" {
		err := json.FromJson(userConfigJson, userConfigBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		return userConfigBo, nil
	} else {
		userConfig := &po.PpmOrgUserConfig{}
		err = mysql.SelectOneByCond(userConfig.TableName(), db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcUserId:   userId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, userConfig)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_ = copyer.Copy(userConfig, userConfigBo)
		userConfigJson, err = json.ToJson(userConfigBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		err = cache.Set(key, userConfigJson)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return userConfigBo, nil
	}
}

func DeleteUserConfigInfo(orgId int64, userId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheUserConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	return nil
}

func ClearUserCacheInfo(token string) errs.SystemErrorInfo{
	userCacheKey := sconsts.CacheUserToken + token
	_, err := cache.Del(userCacheKey)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}