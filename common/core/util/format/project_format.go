package format

import (
	"github.com/galaxy-book/common/core/util/strs"
	"regexp"
)

//项目名
func VerifyProjectNameFormat(input string) bool {
	reg := regexp.MustCompile(ProjectNamePattern)
	return reg.MatchString(input)
}

//项目前缀编号
func VerifyProjectPreviousCodeFormat(input string) bool {
	reg := regexp.MustCompile(ProjectPreviousCodePattern)
	return reg.MatchString(input)
}

//项目简介
func VerifyProjectRemarkFormat(input string) bool {
	reg := regexp.MustCompile(ProjectRemarkPattern)
	return reg.MatchString(input)
}

//任务名
func VerifyIssueNameFormat(input string) bool {
	reg := regexp.MustCompile(IssueNamePattern)
	return reg.MatchString(input)
}

//任务简介
func VerifyIssueRemarkFormat(input string) bool {
	//reg := regexp.MustCompile(IssueRemarkPattern)
	//return reg.MatchString(input)
	inputLen := strs.Len(input)
	if inputLen > 10000 {
		return false
	} else {
		return true
	}
}

//任务评论
func VerifyIssueCommenFormat(input string) bool {
	reg := regexp.MustCompile(IssueCommenPattern)
	return reg.MatchString(input)
}

//任务栏名字
func VerifyProjectObjectTypeNameFormat(input string) bool {
	reg := regexp.MustCompile(ProjectObjectTypeNamePattern)
	return reg.MatchString(input)
}

//项目公告
func VerifyProjectNoticeFormat(input string) bool {
	inputLen := strs.Len(input)
	if inputLen > 2000 {
		return false
	} else {
		return true
	}
}
