package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/service"
)

func (GetGreeter) GetInviteCode(req orgvo.GetInviteCodeReqVo) orgvo.GetInviteCodeRespVo {
	res, err := service.GetInviteCode(req.CurrentUserId, req.OrgId, req.SourcePlatform)
	return orgvo.GetInviteCodeRespVo{Err: vo.NewErr(err), Data: res}
}

func (GetGreeter) GetInviteInfo(req orgvo.GetInviteInfoReqVo) orgvo.GetInviteInfoRespVo {
	res, err := service.GetInviteInfo(req.InviteCode)
	return orgvo.GetInviteInfoRespVo{Err: vo.NewErr(err), Data: res}
}
