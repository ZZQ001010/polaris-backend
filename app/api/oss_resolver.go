package api

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/resourcefacade"
)

func (r *queryResolver) GetOssSignURL(ctx context.Context, input vo.OssApplySignURLReq) (*vo.OssApplySignURLResp, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := resourcefacade.GetOssSignURL(resourcevo.OssApplySignURLReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.GetOssSignURL, respVo.Error()
}

func (r *queryResolver) GetOssPostPolicy(ctx context.Context, input vo.OssPostPolicyReq) (*vo.OssPostPolicyResp, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	req := resourcevo.GetOssPostPolicyReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	}

	respVo := resourcefacade.GetOssPostPolicy(req)
	return respVo.GetOssPostPolicy, respVo.Error()
}
