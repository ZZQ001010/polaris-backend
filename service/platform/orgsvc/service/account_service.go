package service

import (
	"github.com/galaxy-book/common/core/util/md5"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/domain"
	"strings"
)

func SetPassword(req orgvo.SetPasswordReqVo) errs.SystemErrorInfo {
	userId := req.UserId

	targetPassword := strings.TrimSpace(req.Input.Password)
	pwdLen := len(targetPassword)
	if pwdLen < 6 || pwdLen > 16{
		return errs.PwdFormatError
	}

	suc := format.VerifyPwdFormat(targetPassword)
	if ! suc{
		return errs.PwdFormatError
	}

	userBo, _, err := domain.GetUserBo(userId)
	if err != nil{
		log.Error(err)
		return err
	}

	if strings.TrimSpace(userBo.Password) != consts.BlankString{
		return errs.PwdAlreadySettingsErr
	}

	salt := md5.Md5V(uuid.NewUuid())
	pwd := util.PwdEncrypt(targetPassword, salt)
	err = domain.SetUserPassword(userId, pwd, salt, userId)
	if err != nil{
		log.Error(err)
		return err
	}

	return nil
}

func ResetPassword(req orgvo.ResetPasswordReqVo) errs.SystemErrorInfo {
	userId := req.UserId
	input := req.Input

	userBo, _, err := domain.GetUserBo(userId)
	if err != nil{
		log.Error(err)
		return err
	}

	if strings.TrimSpace(userBo.Password) == consts.BlankString{
		return errs.PasswordNotSetError
	}

	salt := userBo.PasswordSalt

	targetPassword := strings.TrimSpace(input.NewPassword)
	pwdLen := len(targetPassword)
	if pwdLen < 6 || pwdLen > 16{
		return errs.PwdFormatError
	}

	suc := format.VerifyPwdFormat(targetPassword)
	if ! suc{
		return errs.PwdFormatError
	}

	currentPwd := util.PwdEncrypt(input.CurrentPassword, salt)
	if currentPwd != userBo.Password{
		return errs.PasswordNotMatchError
	}

	newPassword := util.PwdEncrypt(targetPassword, salt)
	err = domain.SetUserPassword(userId, newPassword, salt, userId)
	if err != nil{
		log.Error(err)
		return err
	}
	return nil
}

func RetrievePassword(req orgvo.RetrievePasswordReqVo) errs.SystemErrorInfo {
	input := req.Input

	authCode := input.AuthCode
	username := input.Username
	password := input.NewPassword

	userBo, err := domain.GetUserInfoByLoginName(username)
	if err != nil{
		log.Error(err)
		if err.Code() == errs.UserNotExist.Code(){
			return errs.NotBindAccountError
		}else{
			return err
		}
	}

	userId := userBo.ID

	//未设置密码也允许他找回
	//if strings.TrimSpace(userBo.Password) == consts.BlankString{
	//	return errs.PasswordNotSetError
	//}

	addressType := consts.ContactAddressTypeMobile
	if format.VerifyEmailFormat(username){
		addressType = consts.ContactAddressTypeEmail
	}

	err1 := domain.AuthCodeVerify(consts.AuthCodeTypeRetrievePwd, addressType, username, authCode)
	if err1 != nil {
		log.Error(err1)
		return err1
	}

	salt := userBo.PasswordSalt

	targetPassword := strings.TrimSpace(password)
	pwdLen := len(targetPassword)
	if pwdLen < 6 || pwdLen > 16{
		return errs.PwdFormatError
	}

	suc := format.VerifyPwdFormat(targetPassword)
	if ! suc{
		return errs.PwdFormatError
	}

	newPassword := util.PwdEncrypt(targetPassword, salt)
	err = domain.SetUserPassword(userId, newPassword, salt, userId)
	if err != nil{
		log.Error(err)
		return err
	}

	return nil
}

func UnbindLoginName(req orgvo.UnbindLoginNameReqVo) errs.SystemErrorInfo {
	userId := req.UserId
	input := req.Input

	addressType := input.AddressType
	authCode := input.AuthCode

	userBo, _, err := domain.GetUserBo(userId)
	if err != nil{
		log.Error(err)
		return err
	}

	username := ""
	if addressType == consts.ContactAddressTypeEmail{
		if strings.TrimSpace(userBo.Email) == consts.BlankString{
			return errs.EmailNotBindError
		}
		username = userBo.Email
	}else if addressType == consts.ContactAddressTypeMobile{
		if strings.TrimSpace(userBo.Mobile) == consts.BlankString{
			return errs.MobileNotBindError
		}
		username = userBo.Mobile
	}else{
		return errs.NotSupportedContactAddressType
	}

	err1 := domain.AuthCodeVerify(consts.AuthCodeTypeUnBind, addressType, username, authCode)
	if err1 != nil {
		log.Error(err1)
		return err1
	}

	err = domain.UnbindUserName(userId, addressType)
	if err != nil{
		log.Error(err)
		return err
	}
	return nil
}

func BindLoginName(req orgvo.BindLoginNameReqVo) errs.SystemErrorInfo {
	userId := req.UserId
	input := req.Input

	username := input.Address
	addressType := input.AddressType
	authCode := input.AuthCode

	userBo, _, err := domain.GetUserBo(userId)
	if err != nil{
		log.Error(err)
		return err
	}

	if addressType == consts.ContactAddressTypeEmail{
		if strings.TrimSpace(userBo.Email) != consts.BlankString{
			return errs.EmailAlreadyBindError
		}
	}else if addressType == consts.ContactAddressTypeMobile{
		if strings.TrimSpace(userBo.Mobile) != consts.BlankString{
			return errs.MobileAlreadyBindError
		}
	}else{
		return errs.NotSupportedContactAddressType
	}

	err1 := domain.AuthCodeVerify(consts.AuthCodeTypeBind, addressType, username, authCode)
	if err1 != nil {
		log.Error(err1)
		return err1
	}

	err = domain.BindUserName(userId, addressType, username)
	if err != nil{
		log.Error(err)
		return err
	}
	return nil
}

func CheckLoginName(req orgvo.CheckLoginNameReqVo) errs.SystemErrorInfo {
	input := req.Input

	addressType := input.AddressType
	address := input.Address
	if addressType == consts.ContactAddressTypeEmail{
		_, err := domain.GetUserInfoByEmail(address)
		if err != nil{
			if err.Code() == errs.UserNotExist.Code(){
				return errs.EmailNotBindAccountError
			}
		}
	}

	if addressType == consts.ContactAddressTypeMobile{
		_, err := domain.GetUserInfoByMobile(address)
		if err != nil {
			if err.Code() == errs.UserNotExist.Code() {
				return errs.MobileNotBindAccountError
			}
		}
	}

	return nil
}