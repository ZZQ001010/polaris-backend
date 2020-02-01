package orgfacade

import (
	"context"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
)

var log = logger.GetDefaultLogger()

func GetCurrentUserRelaxed(ctx context.Context) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	respVo := GetCurrentUser(ctx)
	if respVo.Failure() {
		return nil, respVo.Error()
	}

	return &respVo.CacheInfo, nil
}

func GetCurrentUserWithoutOrgVerifyRelaxed(ctx context.Context) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	respVo := GetCurrentUserWithoutOrgVerify(ctx)
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return &respVo.CacheInfo, nil
}


func GetBaseOrgInfoRelaxed(sourceChannel string, orgId int64) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	respVo := GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{
		SourceChannel: sourceChannel,
		OrgId:         orgId,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}

	return respVo.BaseOrgInfo, nil
}

func GetBaseUserInfoRelaxed(sourceChannel string, orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	respVo := GetBaseUserInfo(orgvo.GetBaseUserInfoReqVo{
		SourceChannel: sourceChannel,
		OrgId:         orgId,
		UserId:        userId,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}

	return respVo.BaseUserInfo, nil
}

func GetDingTalkBaseUserInfoRelaxed(orgId int64, userId int64) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	respVo := GetDingTalkBaseUserInfo(orgvo.GetDingTalkBaseUserInfoReqVo{
		OrgId:  orgId,
		UserId: userId,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}

	return respVo.BaseUserInfo, nil
}

func GetUserConfigInfoRelaxed(orgId int64, userId int64) (*bo.UserConfigBo, errs.SystemErrorInfo) {
	respVo := GetUserConfigInfo(orgvo.GetUserConfigInfoReqVo{
		OrgId:  orgId,
		UserId: userId,
	})

	if respVo.Failure() {
		return nil, respVo.Error()
	}

	return respVo.UserConfigInfo, nil
}

func GetUserInfoRelaxed(orgId int64, userId int64, sourceChannel string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	respVo := GetUserInfo(orgvo.GetUserInfoReqVo{
		OrgId:  orgId,
		UserId: userId,
		SourceChannel: sourceChannel,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}

	return respVo.UserInfo, nil
}

func VerifyOrgRelaxed(orgId int64, userId int64) bool {
	respVo := VerifyOrg(orgvo.VerifyOrgReqVo{
		OrgId:  orgId,
		UserId: userId,
	})
	if respVo.Failure() {
		return false
	}

	return respVo.IsTrue
}

func VerifyOrgUsersRelaxed(orgId int64, userIds []int64) bool{
	respVo := VerifyOrgUsers(orgvo.VerifyOrgUsersReqVo{
		OrgId:  orgId,
		Input: orgvo.VerifyOrgUsersReqData{
			UserIds: userIds,
		},
	})
	if respVo.Failure() {
		return false
	}
	return respVo.IsTrue
}
