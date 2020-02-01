package api

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
)

func (r *queryResolver) IterationStats(ctx context.Context, page *int, size *int, params vo.IterationStatReq) (*vo.IterationStatList, error) {
	pageA := uint(0)
	sizeA := uint(0)
	if page != nil && size != nil && *page > 0 && *size > 0 {
		pageA = uint(*page)
		sizeA = uint(*size)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.IterationStats(projectvo.IterationStatsReqVo{
		Page:  pageA,
		Size:  sizeA,
		Input: params,
		OrgId: cacheUserInfo.OrgId,
	})
	return respVo.IterationStats, respVo.Error()
}
