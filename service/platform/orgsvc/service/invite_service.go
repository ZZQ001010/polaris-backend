package service

import (
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/core/util/uuid"
	consts2 "github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/rand"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/domain"
	"strconv"
)

func GetInviteCode(currentUserId int64, orgId int64, sourcePlatform string) (*orgvo.GetInviteCodeRespVoData, errs.SystemErrorInfo) {
	//用户角色权限校验
	authErr := AuthOrgRole(orgId, currentUserId, consts2.RoleOperationPathOrgUser, consts2.RoleOperationInvite)
	if authErr != nil {
		log.Error(authErr)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, authErr)
	}
	inviteInfo := bo.InviteInfoBo{
		InviterId:      currentUserId,
		OrgId:          orgId,
		SourcePlatform: sourcePlatform,
	}

	inviteCode := rand.RandomInviteCode(uuid.NewUuid() + strconv.FormatInt(currentUserId, 10) + sourcePlatform)
	err := domain.SetUserInviteCodeInfo(inviteCode, inviteInfo)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, err
	}
	return &orgvo.GetInviteCodeRespVoData{InviteCode: inviteCode, Expire: consts.CacheUserInviteCodeExpire}, nil
}

func GetInviteInfo(inviteCode string) (*vo.GetInviteInfoResp, errs.SystemErrorInfo) {
	inviteInfo, err := domain.GetUserInviteCodeInfo(inviteCode)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, err
	}

	orgBaseInfo, err := domain.GetBaseOrgInfo("", inviteInfo.OrgId)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, err
	}
	userBaseInfo, err := domain.GetBaseUserInfo("", inviteInfo.OrgId, inviteInfo.InviterId)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, err
	}

	return &vo.GetInviteInfoResp{
		OrgID:       orgBaseInfo.OrgId,
		OrgName:     orgBaseInfo.OrgName,
		InviterID:   userBaseInfo.UserId,
		InviterName: userBaseInfo.Name,
	}, nil
}
