package api

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/trendsvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/trendsfacade"
)

func (r *queryResolver) NoticeList(ctx context.Context, page *int, size *int, params *vo.NoticeListReq) (*vo.NoticeList, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	var defaultPage, defaultSize int
	if page == nil {
		page = &defaultPage
	}
	if size == nil {
		size = &defaultSize
	}

	resp := trendsfacade.NoticeList(trendsvo.NoticeListReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Page:   *page,
		Size:   *size,
		Input:  params,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.Data, nil
}

func (r *queryResolver) UnreadNoticeCount(ctx context.Context) (*vo.NoticeCountResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := trendsfacade.UnreadNoticeCount(trendsvo.UnreadNoticeCountReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return &vo.NoticeCountResp{Total: resp.Count}, nil
}

func (r *queryResolver) GetMQTTChannelKey(ctx context.Context, input vo.GetMQTTChannelKeyReq) (*vo.GetMQTTChannelKeyResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := trendsfacade.GetMQTTChannelKey(trendsvo.GetMQTTChannelKeyReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input: input,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Data, nil
}