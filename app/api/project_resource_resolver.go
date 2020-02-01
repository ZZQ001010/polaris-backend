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

func (r *mutationResolver) CreateProjectResource(ctx context.Context, input vo.CreateProjectResourceReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.CreateProjectResource(projectvo.CreateProjectResourceReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Void, resp.Error()
}

func (r *mutationResolver) UpdateProjectResourceFolder(ctx context.Context, input vo.UpdateProjectResourceFolderReq) (*vo.UpdateProjectResourceFolderResp, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.UpdateProjectResourceFolder(projectvo.UpdateProjectResourceFolderReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Output, resp.Error()
}

func (r *mutationResolver) UpdateProjectResourceName(ctx context.Context, input vo.UpdateProjectResourceNameReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.UpdateProjectResourceName(projectvo.UpdateProjectResourceNameReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Void, resp.Error()
}

func (r *mutationResolver) DeleteProjectResource(ctx context.Context, input vo.DeleteProjectResourceReq) (*vo.DeleteProjectResourceResp, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.DeleteProjectResource(projectvo.DeleteProjectResourceReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Output, resp.Error()
}

func (r *queryResolver) ProjectResource(ctx context.Context, page *int, size *int, params vo.ProjectResourceReq) (*vo.ResourceList, error) {
	validate, err := validator.Validate(params)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.GetProjectResource(projectvo.GetProjectResourceReqVo{
		Input:  params,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Page:   *page,
		Size:   *size,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	return resp.Output, resp.Error()
}
