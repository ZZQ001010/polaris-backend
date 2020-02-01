package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/domain"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) SendSMSLoginCode(req orgvo.SendSMSLoginCodeReqVo) vo.VoidErr {
	phoneNumber := req.Input.PhoneNumber
	//phone format check
	verifyErr := service.VerifyCaptcha(req.Input.CaptchaID, req.Input.CaptchaPassword, req.Input.PhoneNumber)
	if verifyErr != nil {
		return vo.VoidErr{Err: vo.NewErr(verifyErr)}
	}

	err := service.SendSMSLoginCode(phoneNumber)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) SendAuthCode(req orgvo.SendAuthCodeReqVo) vo.VoidErr {
	verifyErr := service.VerifyCaptcha(req.Input.CaptchaID, req.Input.CaptchaPassword, req.Input.Address)
	if verifyErr != nil {
		return vo.VoidErr{Err: vo.NewErr(verifyErr)}
	}
	err := service.SendAuthCode(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) GetPwdLoginCode(req orgvo.GetPwdLoginCodeReqVo) orgvo.GetPwdLoginCodeRespVo {
	res, err := domain.GetPwdLoginCode(req.CaptchaId)
	return orgvo.GetPwdLoginCodeRespVo{Err:vo.NewErr(err), CaptchaPassword:res}
}

func (PostGreeter) SetPwdLoginCode(req orgvo.SetPwdLoginCodeReqVo) vo.VoidErr {
	return vo.VoidErr{Err:vo.NewErr(domain.SetPwdLoginCode(req.CaptchaId, req.CaptchaPassword))}
}