package coding_carefree

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	cache "github.com/allstar/coding-carefree/components/cache"
	datasources "github.com/allstar/coding-carefree/components/datasources"
	consts "github.com/allstar/coding-carefree/consts"
	domains "github.com/allstar/coding-carefree/modules"
	utils "github.com/allstar/coding-carefree/utils"

	db "upper.io/db.v3"
)

func RegisterService(ctx context.Context, input RegisterInVo) (*Void, error) {
	usr := input.Usr
	pwd := input.Pwd
	code := input.Code

	if !utils.VerifyEmailFormat(usr) {
		return nil, errors.New("邮箱格式错误")
	}

	if !utils.VerifyPwdFormat(pwd) || len(pwd) < 8 || len(pwd) > 16 {
		return nil, errors.New("密码必须同时包含字母或数字，且长度在8~16")
	}

	rp := cache.GetProxy()
	defer rp.Close()

	mailType := string(consts.MAIL_SEND_TYPE_REGISTER)
	mailCodeKey := consts.MAIL_CODE_KEYS + mailType + usr
	mailAuthCountKey := consts.MAIL_AUTH_COUNT_KEYS + mailType + usr
	mailSendLimitKey := consts.MAIL_SEND_LIMIT_KEYS + mailType + usr
	authCode := rp.Get(mailCodeKey)
	if authCode == "" {
		return nil, errors.New("验证码已失效，请重新发送")
	}

	if !strings.EqualFold(authCode, code) {
		count := rp.Incrby(mailAuthCountKey, 1)
		if count > consts.MAIL_AUTO_COUNT_LIMIT {
			rp.Del(mailCodeKey)
			rp.Del(mailAuthCountKey)
			return nil, errors.New("验证错误次数过多，请重新发送验证")
		}
		return nil, errors.New("验证码错误，请重新验证，剩余次数：" + strconv.Itoa(int(consts.MAIL_AUTO_COUNT_LIMIT-count)))
	}

	sess := datasources.GetMysqlConnect()
	defer sess.Close()

	exist, _ := sess.Collection("user").Find(db.Cond{
		"username": usr,
	}).Exists()

	if exist {
		return nil, errors.New("该邮箱已经被注册")
	}

	userPo := domains.User{
		Username: usr,
		Nickname: usr,
		Password: utils.Md5V(pwd),
	}

	_, err := sess.Collection("user").Insert(userPo)
	if err != nil {
		return nil, err
	}
	rp.Del(mailCodeKey)
	rp.Del(mailAuthCountKey)
	rp.Del(mailSendLimitKey)

	return &Void{200}, nil
}

func LoginService(ctx context.Context, input LoginInVo) (*LoginOutVo, error) {
	sess := datasources.GetMysqlConnect()
	defer sess.Close()

	rp := cache.GetProxy()
	defer rp.Close()

	user := domains.User{}

	err := sess.Collection("user").Find(db.Cond{
		"username": input.Usr,
		"password": utils.Md5V(input.Pwd),
	}).One(&user)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("用户不存在或用户名密码错误")
	}

	token := utils.GetToken()

	b, err := json.Marshal(user)

	rp.SetEx(consts.USER_TOKEN_KEYS+strconv.FormatInt(user.Id, 10), token, consts.USER_TOKEN_EXPIRE)
	rp.SetEx(consts.TOKEN_USER_KEYS+token, string(b), consts.USER_TOKEN_EXPIRE)

	fmt.Println(string(b))
	return &LoginOutVo{
		200,
		token,
	}, nil
}
