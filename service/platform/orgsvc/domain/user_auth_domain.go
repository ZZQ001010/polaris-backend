package domain

import (
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"upper.io/db.v3/lib/sqlbuilder"
)

func DingAuth(corpId, openId, name string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo){
	//获取组织信息
	orgInfo, err := GetOrgInfoByOutOrgId(corpId, consts.AppSourceChannelDingTalk)
	if err != nil{
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.OrgNotInitError)
	}
	//获取用户信息
	baseUserInfo, err := GetDingTalkBaseUserInfoByEmpId(orgInfo.OrgId, openId)
	if err != nil {
		//这里做用户初始化的兜底
		lockKey := consts.InitUserLock + consts.AppSourceChannelDingTalk + openId
		suc, err := cache.TryGetDistributedLock(lockKey, openId)
		log.Infof("准备获取分布式锁 %v", suc)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		if suc {
			log.Infof("获取分布式锁成功 %v", suc)
			defer func() {
				if _, lockErr := cache.ReleaseDistributedLock(lockKey, openId); lockErr != nil{
					log.Error(lockErr)
				}
			}()

			//double check
			baseUserInfo, err = GetFeiShuBaseUserInfoByEmpId(orgInfo.OrgId, openId)
			if err != nil {
				err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
					_, err := InitDingTalkUser(orgInfo.OrgId, corpId, openId, tx)
					if err != nil{
						log.Error(err)
						return err
					}
					return nil
				})
				if err1 != nil{
					log.Error(err1)
					return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
				}
			}
		}
		if baseUserInfo == nil{
			baseUserInfo, err = GetFeiShuBaseUserInfoByEmpId(orgInfo.OrgId, openId)
			if err != nil {
				log.Error(err)
				return nil, errs.UserInitError
			}
		}
	}
	return baseUserInfo, nil
}
