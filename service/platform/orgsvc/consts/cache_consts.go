package consts

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
)

var (
	//用户配置缓存
	CacheUserConfig = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "config"
	//用户基础信息缓存key
	CacheBaseUserInfo = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "baseinfo"
	//用户外部信息缓存key
	CacheBaseUserOutInfo = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfSourceChannel + consts.CacheKeyOfUser + "outinfo"

	//组织基础信息
	CacheBaseOrgInfo = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + "baseinfo"
	//组织外部信息
	CacheBaseOrgOutInfo = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfSourceChannel + "outinfo"
	//获取外部组织id关联的内部组织id
	CacheOutOrgIdRelationId = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOutOrg + consts.CacheKeyOfSourceChannel + "org_id_info"

	//获取外部用户id关联的内部用户id
	CacheOutUserIdRelationId = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfSourceChannel + consts.CacheKeyOfOutUser + "user_id"
	//部门对应关系
	CacheDeptRelation = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + "dept_relation_list"
	//用户token
	CacheUserToken = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + "user:token:"

	//用户邀请code, 拼接 inviteCode
	CacheUserInviteCode = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + "invite_code:"
	//用户邀请code有效时间为: 1小时
	CacheUserInviteCodeExpire = 60 * 60 * 1

	// 短信验证码相关 + 手机号, 验证失败五次间隔调整为五分钟，验证失败50次冻结一天
	// 短信发送时间间隔: 一分钟，五分钟
	CacheSmsSendLoginCodeFreezeTime = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + consts.CacheKeyOfAuthType + consts.CacheKeyOfAddressType + consts.CacheKeyOfPhone + "sms_auth_code:freeze_time"
	// 短信验证码: 五分钟
	CacheSmsLoginCode = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + consts.CacheKeyOfAuthType + consts.CacheKeyOfAddressType + consts.CacheKeyOfPhone + "sms_auth_code"
	// 号码白名单
	CachePhoneNumberWhiteList = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + "sms_white_list"
	// 登录短信验证失败次数
	CacheSmsLoginCodeVerifyFailTimes = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + consts.CacheKeyOfAuthType + consts.CacheKeyOfAddressType + consts.CacheKeyOfPhone + "sms_auth_code:verify_times"
	//登录图形验证码：一分钟
	CacheLoginGraphCode = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + consts.CacheKeyOfLoginName + "graph_auth_code"
)

const (
	OrgCodeLength = 50
)
