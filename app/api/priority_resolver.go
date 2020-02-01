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

func (r *queryResolver) Prioritys(ctx context.Context, page *int, size *int, params *vo.PriorityListReq) (*vo.PriorityList, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	reqBody := projectvo.PriorityListReqVo{
		Page:  page,
		Size:  size,
		OrgId: cacheUserInfo.OrgId,
	}
	if params != nil {
		reqBody.Type = params.Type
	}

	resp := projectfacade.PriorityList(reqBody)
	return resp.PriorityList, resp.Error()
}

func (r *mutationResolver) CreatePriority(ctx context.Context, input vo.CreatePriorityReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.CreatePriority(projectvo.CreatePriorityReqVo{
		CreatePriorityReq: input,
		UserId:            cacheUserInfo.UserId,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) UpdatePriority(ctx context.Context, input vo.UpdatePriorityReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.UpdatePriority(projectvo.UpdatePriorityReqVo{
		UpdatePriorityReq: input,
		UserId:            cacheUserInfo.UserId,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) DeletePriority(ctx context.Context, input vo.DeletePriorityReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.DeletePriority(projectvo.DeletePriorityReqVo{
		DeletePriorityReq: input,
		UserId:            cacheUserInfo.UserId,
		OrgId:             cacheUserInfo.OrgId,
	})
	return resp.Void, resp.Error()
}
