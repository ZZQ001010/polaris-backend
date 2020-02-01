package api

import (
	"context"
	"github.com/galaxy-book/common/core/util/validator"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
)

func (r *queryResolver) IssueSources(ctx context.Context, page *int, size *int, params *vo.IssueSourcesReq) (*vo.IssueSourceList, error) {
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

	respVo := projectfacade.IssueSourceList(projectvo.IssueSourceListReqVo{
		Page:  pageA,
		Size:  sizeA,
		Input: params,
		OrgId: cacheUserInfo.OrgId,
	})
	return respVo.IssueSourceList, respVo.Error()
}

func (r *mutationResolver) CreateIssueSource(ctx context.Context, input vo.CreateIssueSourceReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.CreateIssueSource(projectvo.CreateIssueSourceReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Void, respVo.Error()
}

func (r *mutationResolver) UpdateIssueSource(ctx context.Context, input vo.UpdateIssueSourceReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.UpdateIssueSource(projectvo.UpdateIssueSourceReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Void, respVo.Error()
}

func (r *mutationResolver) DeleteIssueSource(ctx context.Context, input vo.DeleteIssueSourceReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.DeleteIssueSource(projectvo.DeleteIssueSourceReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.Void, respVo.Error()
}
