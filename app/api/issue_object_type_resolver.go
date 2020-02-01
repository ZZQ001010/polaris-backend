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

func (r *queryResolver) IssueObjectTypes(ctx context.Context, page *int, size *int, params *vo.IssueObjectTypesReq) (*vo.IssueObjectTypeList, error) {
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

	respVo := projectfacade.IssueObjectTypeList(projectvo.IssueObjectTypeListReqVo{
		Page:   pageA,
		Size:   sizeA,
		Params: params,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.IssueObjectTypeList, respVo.Error()
}

func (r *mutationResolver) CreateIssueObjectType(ctx context.Context, input vo.CreateIssueObjectTypeReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.CreateIssueObjectType(projectvo.CreateIssueObjectTypeReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Void, respVo.Error()
}

func (r *mutationResolver) UpdateIssueObjectType(ctx context.Context, input vo.UpdateIssueObjectTypeReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.UpdateIssueObjectType(projectvo.UpdateIssueObjectTypeReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Void, respVo.Error()
}

func (r *mutationResolver) DeleteIssueObjectType(ctx context.Context, input vo.DeleteIssueObjectTypeReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.DeleteIssueObjectType(projectvo.DeleteIssueObjectTypeReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.Void, respVo.Error()
}
