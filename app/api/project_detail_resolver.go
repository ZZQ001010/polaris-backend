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

func (r *queryResolver) ProjectDetail(ctx context.Context, projectID int64) (*vo.ProjectDetail, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.ProjectDetail(projectvo.ProjectDetailReqVo{
		ProjectId: projectID,
		OrgId:     cacheUserInfo.OrgId,
	})
	return resp.ProjectDetail, resp.Error()
}

func (r *mutationResolver) CreateProjectDetail(ctx context.Context, input vo.CreateProjectDetailReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.CreateProjectDetail(projectvo.CreateProjectDetailReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) UpdateProjectDetail(ctx context.Context, input vo.UpdateProjectDetailReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.UpdateProjectDetail(projectvo.UpdateProjectDetailReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) DeleteProjectDetail(ctx context.Context, input vo.DeleteProjectDetailReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.DeleteProjectDetail(projectvo.DeleteProjectDetailReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
	})
	return resp.Void, resp.Error()
}
