package format

import "regexp"

//组织名
func VerifyOrgNameFormat(input string) bool {
	reg := regexp.MustCompile(OrgNamePattern)
	return reg.MatchString(input)
}

//网址后缀
func VerifyOrgCodeFormat(input string) bool {
	reg := regexp.MustCompile(OrgCodePattern)
	return reg.MatchString(input)
}

//详细地址
func VerifyOrgAdressFormat(input string) bool {
	reg := regexp.MustCompile(OrgAdressPattern)
	return reg.MatchString(input)
}
