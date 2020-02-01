package service

import (
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/random"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/orgsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/domain"
)

func UserLogin(req vo.UserLoginReq) (*vo.UserLoginResp, errs.SystemErrorInfo) {
	loginType := req.LoginType
	var userBo *bo.UserInfoBo = nil
	var err errs.SystemErrorInfo = nil
	if loginType == 1 || loginType == 3{
		userBo, err = UserAuthCodeLogin(req)
	} else if loginType == 2 {
		userBo, err = UserPwdLogin(req)
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if userBo == nil {
		//暂时不支持的登录类型
		return nil, errs.BuildSystemErrorInfo(errs.UnSupportLoginType)
	}

	token := random.Token()
	res := &vo.UserLoginResp{
		Token:       token,
		UserID:      userBo.ID,
		OrgID:       userBo.OrgID,
		Name:        userBo.Name,
		Avatar:      userBo.Avatar,
		NeedInitOrg: userBo.OrgID == 0,
	}

	if userBo.OrgID != 0 {
		organizationBo, err := domain.GetOrgBoById(userBo.OrgID)
		if err != nil {
			log.Error(err)
		} else {
			res.OrgCode = organizationBo.Code
			res.OrgName = organizationBo.Name
		}
	}

	cacheUserBo := &bo.CacheUserInfoBo{
		UserId:        userBo.ID,
		SourceChannel: req.SourceChannel,
		OrgId:         userBo.OrgID,
	}
	cacheErr := cache.SetEx(sconsts.CacheUserToken+res.Token, json.ToJsonIgnoreError(cacheUserBo), consts.CacheUserTokenExpire)
	if cacheErr != nil {
		log.Error(cacheErr)
	}

	//更新上次登录时间
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("捕获到的错误：%s", r)
			}
		}()
		if err := domain.UserLoginHook(userBo.ID, userBo.OrgID); err != nil {
			log.Error(err)
		}
	}()

	return res, nil
}

func UserAuthCodeLogin(req vo.UserLoginReq) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	loginName := req.LoginName
	authCode := req.AuthCode
	loginType := req.LoginType

	addressType := consts.ContactAddressTypeEmail
	if loginType == consts.LoginTypeSMSCode{
		addressType = consts.ContactAddressTypeMobile
	}

	name := ""
	inviteCode := ""
	if req.Name != nil {
		name = *req.Name
	}
	if req.InviteCode != nil {
		inviteCode = *req.InviteCode
	}
	if authCode == nil {
		log.Error("验证码不能为空")
		return nil, errs.AuthCodeIsNull
	}

	//台湾跳过短信验证 && config.GetEnv() != consts.RunEnvStag
	if config.GetEnv() != consts.RunEnvProdTw {
		err1 := domain.AuthCodeVerify(consts.AuthCodeTypeLogin, addressType, loginName, *authCode)
		if err1 != nil {
			log.Error(err1)
			return nil, err1
		}
	}

	bos, _ := domain.GetUserInfoListByOrg(1001)

	fmt.Printf("测试数据%v \n", bos)

	var loginUserBo bo.UserInfoBo
	//手机登录
	if loginType == consts.LoginTypeSMSCode{
		//做登录和自动注册逻辑
		userBo, err := domain.GetUserInfoByMobile(loginName)
		if err != nil {
			log.Infof("用户%s未注册，开始注册....", loginName)
			userBo, err = domain.UserRegister(bo.UserSMSRegisterInfo{
				PhoneNumber:    loginName,
				SourceChannel:  req.SourceChannel,
				SourcePlatform: req.SourcePlatform,
				Name:           name,
				InviteCode:     inviteCode,
			})
			if err != nil {
				log.Error(err)
				return nil, err
			}
		}
		loginUserBo = *userBo
	}else if loginType == consts.LoginTypeMail{
		userBo, err := domain.GetUserInfoByEmail(loginName)
		if err != nil {
			//log.Infof("用户%s未注册，开始注册....", loginName)
			//userBo, err = domain.UserRegister(bo.UserSMSRegisterInfo{
			//	Email: loginName,
			//	SourceChannel:  req.SourceChannel,
			//	SourcePlatform: req.SourcePlatform,
			//	Name:           name,
			//	InviteCode:     inviteCode,
			//})
			//if err != nil {
			//	log.Error(err)
			//	return nil, err
			//}
			//暂时不支持自动注册
			if err != nil{
				log.Error(err)
				return nil, errs.EmailNotRegisterError
			}
		}
		loginUserBo = *userBo
	}
	err := UserAlreadyRegisterHandleHook(req, loginUserBo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &loginUserBo, nil
}

func UserAlreadyRegisterHandleHook(req vo.UserLoginReq, userBo bo.UserInfoBo) errs.SystemErrorInfo {
	//判断邀请逻辑
	if req.InviteCode != nil {
		inviteCodeInfo, err := domain.GetUserInviteCodeInfo(*req.InviteCode)
		if err != nil {
			log.Error(err)
			return err
		}
		orgId := inviteCodeInfo.OrgId
		inviterId := inviteCodeInfo.InviterId
		userId := userBo.ID

		//这里为用户切换组织
		userBo.OrgID = orgId

		//修改用户默认组织
		updateUserInfoErr := domain.UpdateUserDefaultOrg(userId, orgId)
		if updateUserInfoErr != nil {
			log.Error(updateUserInfoErr)
		}

		err = domain.AddOrgMember(orgId, userId, inviterId, true, false)
		if err != nil {
			log.Error(err)
			return err
		}

		//增加动态
		asyn.Execute(func() {
			domain.PushOrgTrends(bo.OrgTrendsBo{
				OrgId: orgId,
				PushType: consts.PushTypeApplyJoinOrg,
				TargetMembers: []int64{userId},
				SourceChannel: userBo.SourceChannel,
				OperatorId: inviterId,
			})
		})
	}
	return nil
}

//用户密码登录
func UserPwdLogin(req vo.UserLoginReq) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	if req.Password == nil {
		return nil, errs.PasswordEmptyError
	}
	loginName := req.LoginName
	pwd := *req.Password

	//账号密码登录暂时不用验证码
	//authCode := req.AuthCode

	//localCode, err := domain.GetPwdLoginCode(loginName)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}

	//if !strings.EqualFold(localCode, *authCode) {
	//	//	_ = domain.ClearPwdLoginCode(loginName)
	//	//	return nil, errs.BuildSystemErrorInfo(errs.PwdLoginCodeNotMatch)
	//	//}

	userBo, err := domain.GetUserInfoByPwd(loginName, pwd)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return userBo, nil
}

//用户退出
func UserQuit(req orgvo.UserQuitReqVo) errs.SystemErrorInfo {
	return domain.ClearUserCacheInfo(req.Token)
}
