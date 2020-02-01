package api

import (
	"context"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
)

func (r *queryResolver) GetInviteCode(ctx context.Context, input *vo.GetInviteCodeReq) (*vo.GetInviteCodeResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := orgfacade.GetInviteCode(orgvo.GetInviteCodeReqVo{
		CurrentUserId:  cacheUserInfo.UserId,
		OrgId:          cacheUserInfo.OrgId,
		SourcePlatform: cacheUserInfo.SourceChannel,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return &vo.GetInviteCodeResp{InviteCode: respVo.Data.InviteCode, Expire: respVo.Data.Expire}, nil
}

func (r *queryResolver) GetInviteInfo(ctx context.Context, input vo.GetInviteInfoReq) (*vo.GetInviteInfoResp, error) {
	respVo := orgfacade.GetInviteInfo(orgvo.GetInviteInfoReqVo{
		InviteCode: input.InviteCode,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.Data, nil
}
