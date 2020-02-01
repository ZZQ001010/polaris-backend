package service

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

var Log = logger.GetDefaultLogger()

func AssemblyUserIdInfo(baseUserInfo *bo.BaseUserInfoBo) *vo.UserIDInfo {
	return &vo.UserIDInfo{
		UserID:  baseUserInfo.UserId,
		Name:    baseUserInfo.Name,
		Avatar:  baseUserInfo.Avatar,
		EmplID:  baseUserInfo.OutUserId,
		IsDeleted: baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted,
		IsDisabled: baseUserInfo.OrgUserStatus == consts.AppStatusDisabled,
	}
}

func NeedUpdate(updateFields []string, field string) bool {
	if updateFields == nil || len(updateFields) == 0 {
		return true
	}
	bol, err := slice.Contain(updateFields, field)
	if err != nil {
		return false
	}
	return bol
}
