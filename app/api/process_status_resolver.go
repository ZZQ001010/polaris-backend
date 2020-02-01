package api

import (
	"context"
	"github.com/galaxy-book/common/core/util/validator"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/processvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
)

func (r *queryResolver) ProcessStatuss(ctx context.Context, page *int, size *int) (*vo.ProcessStatusList, error) {
	pageA := uint(0)
	sizeA := uint(0)
	if page != nil && size != nil && *page > 0 && *size > 0 {
		pageA = uint(*page)
		sizeA = uint(*size)
	}
	resp := processfacade.ProcessStatusList(vo.BasicReqVo{
		Page: pageA,
		Size: sizeA,
	})
	return resp.ProcessStatusList, resp.Error()
}

func (r *mutationResolver) CreateProcessStatus(ctx context.Context, input vo.CreateProcessStatusReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := processfacade.CreateProcessStatus(processvo.CreateProcessStatusReqVo{
		CreateProcessStatusReq: input,
		UserId:                 cacheUserInfo.UserId,
		OrgId:                  cacheUserInfo.OrgId,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) UpdateProcessStatus(ctx context.Context, input vo.UpdateProcessStatusReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := processfacade.UpdateProcessStatus(processvo.UpdateProcessStatusReqVo{
		UpdateProcessStatusReq: input,
		UserId:                 cacheUserInfo.UserId,
		OrgId:                  cacheUserInfo.OrgId,
	})

	return resp.Void, resp.Error()
}

func (r *mutationResolver) DeleteProcessStatus(ctx context.Context, input vo.DeleteProcessStatusReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := processfacade.DeleteProcessStatus(processvo.DeleteProcessStatusReq{
		DeleteProcessStatusReq: input,
		UserId:                 cacheUserInfo.UserId,
		OrgId:                  cacheUserInfo.OrgId,
	})
	return resp.Void, resp.Error()
}
