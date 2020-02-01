package orgvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type SendSMSLoginCodeReqVo struct {
	Input vo.SendSmsLoginCodeReq `json:"input"`
}

type SendAuthCodeReqVo struct {
	Input vo.SendAuthCodeReq `json:"input"`
}

type UserLoginReqVo struct {
	UserLoginReq vo.UserLoginReq `json:"userLoginReq"`
}

type UserSMSLoginRespVo struct {
	vo.Err
	Data *vo.UserLoginResp `json:"data"`
}

type GetInviteCodeReqVo struct {
	CurrentUserId  int64  `json:"userId"`
	OrgId          int64  `json:"orgId"`
	SourcePlatform string `json:"sourcePlatform"`
}

type GetInviteInfoReqVo struct {
	InviteCode string `json:"inviteCode"`
}

type GetInviteInfoRespVo struct {
	Data *vo.GetInviteInfoResp `json:"data"`
	vo.Err
}

type GetInviteCodeRespVo struct {
	Data *GetInviteCodeRespVoData `json:"data"`
	vo.Err
}

type GetInviteCodeRespVoData struct {
	InviteCode string `json:"inviteCode"`
	Expire int `json:"expire"`
}

type UserQuitReqVo struct {
	Token string `json:"token"`
}

type GetPwdLoginCodeReqVo struct {
	CaptchaId string `json:"captchaId"`
}

type GetPwdLoginCodeRespVo struct {
	CaptchaPassword string `json:"data"`
	vo.Err
}

type SetPwdLoginCodeReqVo struct {
	CaptchaId string `json:"captchaId"`
	CaptchaPassword string `json:"captchaPassword"`
}