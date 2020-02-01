package service

import (
	"fmt"
	"github.com/galaxy-book/common/core/util/temp"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/rand"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/facade/msgfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/domain"
)

const defaultAuthCode = "000000"

func SendSMSLoginCode(phoneNumber string) errs.SystemErrorInfo {
	return SendAuthCode(orgvo.SendAuthCodeReqVo{
		Input: vo.SendAuthCodeReq{
			AuthType: consts.AuthCodeTypeLogin,
			AddressType: consts.ContactAddressTypeMobile,
			Address: phoneNumber,
		},
	})
}

func SendAuthCode(req orgvo.SendAuthCodeReqVo) errs.SystemErrorInfo {
	input := req.Input

	addressType := input.AddressType
	authType := input.AuthType
	contactAddress := input.Address

	if addressType != consts.ContactAddressTypeMobile && addressType != consts.ContactAddressTypeEmail{
		return errs.NotSupportedContactAddressType
	}

	limitErr := domain.CheckSMSLoginCodeFreezeTime(authType, addressType, contactAddress)
	if limitErr != nil {
		log.Error(limitErr)
		return limitErr
	}

	//如果不是注册，登录，绑定，该账户必须存在
	if authType != consts.AuthCodeTypeRegister && authType != consts.AuthCodeTypeLogin && authType != consts.AuthCodeTypeBind {
		_, err := domain.GetUserInfoByLoginName(contactAddress)
		if err != nil{
			log.Error(err)
			if err.Code() == errs.UserNotExist.Code(){
				return errs.NotBindAccountError
			}else{
				return err
			}
		}
	}

	//如果是注册或者绑定，该账户必须不存在
	if authType == consts.AuthCodeTypeRegister || authType == consts.AuthCodeTypeBind{
		_, err := domain.GetUserInfoByLoginName(contactAddress)
		if err == nil{
			return errs.AccountAlreadyBindError
		}
	}

	authCode := defaultAuthCode
	if !IsInWhiteList(contactAddress) {
		authCode = rand.RandomVerifyCode(6)
		//异步发送
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("捕获到的错误：%s\n", r)
				}
			}()
			switch addressType {
			case consts.ContactAddressTypeMobile:
				sendErr := sendSmsAuthCode(authType, contactAddress, authCode)
				if sendErr != nil{
					log.Error(sendErr)
				}
			case consts.ContactAddressTypeEmail:
				sendErr := sendMailAuthCode(authType, contactAddress, authCode)
				if sendErr != nil{
					log.Error(sendErr)
				}
			}
		}()
	}
	setFreezeErr := domain.SetSMSLoginCodeFreezeTime(authType, addressType, contactAddress, 1)
	if setFreezeErr != nil {
		//这里不要影响主流程
		log.Error(setFreezeErr)
	}
	setLoginCode := domain.SetSMSLoginCode(authType, addressType, contactAddress, authCode)
	if setLoginCode != nil {
		//这里不要影响主流程
		log.Error(setLoginCode)
	}
	return nil
}

func sendSmsAuthCode(authType int, mobile string, authCode string) errs.SystemErrorInfo{
	switch authType {
	case consts.AuthCodeTypeLogin:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeLoginAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeRegister:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeRegisterAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeResetPwd:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeResetPwdAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeRetrievePwd:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeRetrievePwdAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeBind:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeBindAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeUnBind:
		return msgfacade.SendSMSRelaxed(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeUnBindAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	}
	return errs.NotSupportedAuthCodeType
}

func sendMailAuthCode(authType int, email string, authCode string) errs.SystemErrorInfo{
	emails := []string{email}

	switch authType {
	case consts.AuthCodeTypeLogin:
		return msgfacade.SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeLogin, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionLogin,
		}))
	case consts.AuthCodeTypeRegister:
		return msgfacade.SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeRegister, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionRegister,
		}))
	case consts.AuthCodeTypeResetPwd:
		return msgfacade.SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeResetPwd, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionResetPwd,
		}))
	case consts.AuthCodeTypeRetrievePwd:
		return msgfacade.SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeRetrievePwd, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionRetrievePwd,
		}))
	case consts.AuthCodeTypeBind:
		return msgfacade.SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeBind, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionBind,
		}))
	case consts.AuthCodeTypeUnBind:
		return msgfacade.SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeUnBind, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionUnBind,
		}))
	}
	return errs.NotSupportedAuthCodeType
}

func IsInWhiteList(phoneNumber string) bool {
	whiteList, err := domain.GetPhoneNumberWhiteList()
	if err != nil {
		log.Error(err)
		return false
	}
	for _, v := range whiteList {
		if v == phoneNumber {
			return true
		}
	}

	return false
}

func VerifyCaptcha(captchaID, captchaPassword *string, phoneNumber string) errs.SystemErrorInfo {
	if IsInWhiteList(phoneNumber) {
		return nil
	}
	if captchaID == nil || captchaPassword == nil  {
		return errs.CaptchaError
	}

	res, err := domain.GetPwdLoginCode(*captchaID)
	if err != nil {
		log.Error(err)
		return err
	}

	clearErr := domain.ClearPwdLoginCode(*captchaID)
	if clearErr != nil {
		log.Error(clearErr)
		return clearErr
	}

	if res != *captchaPassword {
		return errs.CaptchaError
	}

	return nil
}