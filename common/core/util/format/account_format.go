package format

import (
	"regexp"
)

func VerifyPwdFormat(password string) bool {
	reg := regexp.MustCompile(PasswordPattern)
	return reg.MatchString(password)
}

//用户名
func VerifyUserNameFormat(input string) bool {
	blankReg := regexp.MustCompile(AllBlankPattern)
	if blankReg.MatchString(input) {
		return false
	}
	reg := regexp.MustCompile(ChinesePattern)
	formInput := reg.ReplaceAllString(input, "aa")
	reg = regexp.MustCompile(UserNamePattern)
	return reg.MatchString(formInput)
}
