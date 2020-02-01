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

func (r *queryResolver) ProjectObjectTypes(ctx context.Context, page *int, size *int, params *vo.ProjectObjectTypesReq) (*vo.ProjectObjectTypeList, error) {
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

	resp := projectfacade.ProjectObjectTypeList(projectvo.ProjectObjectTypesReqVo{
		Page:   pageA,
		Size:   sizeA,
		Params: params,
		OrgId:  cacheUserInfo.OrgId,
	})
	return resp.ProjectObjectTypeList, resp.Error()
}

func (r *mutationResolver) CreateProjectObjectType(ctx context.Context, input vo.CreateProjectObjectTypeReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.CreateProjectObjectType(projectvo.CreateProjectObjectTypeReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) UpdateProjectObjectType(ctx context.Context, input vo.UpdateProjectObjectTypeReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.UpdateProjectObjectType(projectvo.UpdateProjectObjectTypeReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) DeleteProjectObjectType(ctx context.Context, input vo.DeleteProjectObjectTypeReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.DeleteProjectObjectType(projectvo.DeleteProjectObjectTypeReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})

	return resp.Void, resp.Error()
}

func (r *queryResolver) ProjectSupportObjectTypes(ctx context.Context, input vo.ProjectSupportObjectTypeListReq) (*vo.ProjectSupportObjectTypeListResp, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.ProjectSupportObjectTypes(projectvo.ProjectSupportObjectTypesReqVo{
		Input: input,
		OrgId: cacheUserInfo.OrgId,
	})

	return resp.ProjectSupportObjectTypes, resp.Error()
}

func (r *queryResolver) ProjectObjectTypesWithProject(ctx context.Context, projectID int64) (*vo.ProjectObjectTypeWithProjectList, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.ProjectObjectTypesWithProject(projectvo.ProjectObjectTypeWithProjectVo{
		ProjectId: projectID,
		OrgId:     cacheUserInfo.OrgId,
	})
	return resp.ProjectObjectTypeWithProjectList, resp.Error()
}
