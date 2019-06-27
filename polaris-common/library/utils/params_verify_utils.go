package utils

import (
	"regexp"
)

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//密码格式校验：必须同时包含字母或数字，且长度在8~16
func VerifyPwdFormat(pwd string) bool {
	pattern := `[A-Za-z]+[0-9]+|[0-9]+[A-Za-z]+`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(pwd)
}
