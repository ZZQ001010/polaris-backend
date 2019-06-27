package coding_carefree

import (
	"context"
	"errors"
	"strconv"

	cache "github.com/allstar/coding-carefree/components/cache"
	mail "github.com/allstar/coding-carefree/components/mail"
	consts "github.com/allstar/coding-carefree/consts"
	utils "github.com/allstar/coding-carefree/utils"
)

func SendMailService(ctx context.Context, input SendMailInVo) (*Void, error) {
	code := utils.GetRandomString5()

	email := input.Usr
	typ := input.Type

	mailSendLimitKey := consts.MAIL_SEND_LIMIT_KEYS + strconv.Itoa(typ) + email
	mailAuthCountKey := consts.MAIL_AUTH_COUNT_KEYS + strconv.Itoa(typ) + email
	mailCodeKey := consts.MAIL_CODE_KEYS + strconv.Itoa(typ) + email

	rp := cache.GetProxy()
	defer rp.Close()

	if rp.Exist(mailSendLimitKey) {
		return nil, errors.New("发送频繁")
	}
	err := mail.SendMail([]string{email}, consts.SUBJECT_OF_REGISTER, "验证码："+code)
	if err == nil {
		rp.SetEx(mailSendLimitKey, code, consts.MAIL_SEND_LIMIT)
		rp.SetEx(mailCodeKey, code, consts.MAIL_CODE_EXPIRE)
		rp.Incrby(mailAuthCountKey, 0)
		rp.Expire(mailAuthCountKey, consts.MAIL_CODE_EXPIRE)
		return &Void{200}, nil
	}
	return nil, err
}
