package api

import (
	"context"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/validator"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
)

func (r *queryResolver) Projects(ctx context.Context, page int, size int, params map[string]interface{}, order []*string, input *vo.ProjectsReq) (*vo.ProjectList, error) {
	maxPageSize := config.GetParameters().MaxPageSize
	if size > maxPageSize {
		return nil, errs.PageSizeOverflowMaxSizeError
	}

	err := validator.ValidateConds("Projects", &params)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.OutOfConditionError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.Projects(projectvo.ProjectsRepVo{
		Page: page,
		Size: size,
		ProjectExtraBody: projectvo.ProjectExtraBody{
			Params: params,
			Order:  order,
			Input:  input,
		},
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		SourceChannel:cacheUserInfo.SourceChannel,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.ProjectList, nil
}

func (r *mutationResolver) CreateProject(ctx context.Context, input vo.CreateProjectReq) (*vo.Project, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.CreateProject(projectvo.CreateProjectReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.Project, nil
}

func (r *mutationResolver) UpdateProject(ctx context.Context, input vo.UpdateProjectReq) (*vo.Project, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.UpdateProject(projectvo.UpdateProjectReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.Project, nil
}

func (r *mutationResolver) QuitProject(ctx context.Context, projectID int64) (*vo.QuitResult, error) {
	validate, err := validator.Validate(projectID)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.QuitProject(projectvo.ProjectIdReqVo{
		ProjectId:     projectID,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})

	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.QuitProject, nil
}

func (r *mutationResolver) ConvertCode(ctx context.Context, input vo.ConvertCodeReq) (*vo.ConvertCodeResp, error) {
	if input.Name == "" {
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "code转换：名称不能为空")
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.ConvertCode(projectvo.ConvertCodeReqVo{
		Input: input,
		OrgId: cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.ConvertCode, nil
}

func (r *mutationResolver) StarProject(ctx context.Context, projectID int64) (*vo.OperateProjectResp, error) {
	validate, err := validator.Validate(projectID)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.StarProject(projectvo.ProjectIdReqVo{
		ProjectId:     projectID,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.OperateProject, nil
}

func (r *mutationResolver) UnstarProject(ctx context.Context, projectID int64) (*vo.OperateProjectResp, error) {
	validate, err := validator.Validate(projectID)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.UnstarProject(projectvo.ProjectIdReqVo{
		ProjectId:     projectID,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.OperateProject, nil
}

func (r *queryResolver) ProjectIssueRelatedStatus(ctx context.Context, input vo.ProjectIssueRelatedStatusReq) ([]*vo.HomeIssueStatusInfo, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.ProjectIssueRelatedStatus(projectvo.ProjectIssueRelatedStatusReqVo{
		Input: input,
		OrgId: cacheUserInfo.OrgId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.ProjectIssueRelatedStatus, nil
}

func (r *queryResolver) ProjectStatistics(ctx context.Context, id int64) (*vo.ProjectStatisticsResp, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.ProjectStatistics(projectvo.ProjectIdReqVo{
		ProjectId: id,
		UserId:    cacheUserInfo.UserId,
		OrgId:     cacheUserInfo.OrgId,
	})

	return resp.ProjectStatistics, resp.Error()
}

func (r *mutationResolver) UpdateProjectStatus(ctx context.Context, input vo.UpdateProjectStatusReq) (*vo.Void, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.UpdateProjectStatus(projectvo.UpdateProjectStatusReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})

	return resp.Void, resp.Error()
}

func (r *queryResolver) ProjectInfo(ctx context.Context, input vo.ProjectInfoReq) (*vo.ProjectInfo, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.ProjectInfo(projectvo.ProjectInfoReqVo{
		Input: input,
		OrgId: cacheUserInfo.OrgId,
		SourceChannel:cacheUserInfo.SourceChannel,
	})

	return resp.ProjectInfo, resp.Error()
}

func (r *mutationResolver) ArchiveProject(ctx context.Context, projectID int64) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.ArchiveProject(projectvo.ProjectIdReqVo{
		OrgId:     cacheUserInfo.OrgId,
		UserId:    cacheUserInfo.UserId,
		ProjectId: projectID,
	})
	return resp.Void, resp.Error()
}

func (r *mutationResolver) CancelArchivedProject(ctx context.Context, projectID int64) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := projectfacade.CancelArchivedProject(projectvo.ProjectIdReqVo{
		OrgId:     cacheUserInfo.OrgId,
		UserId:    cacheUserInfo.UserId,
		ProjectId: projectID,
	})
	return resp.Void, resp.Error()
}

func (r *queryResolver) OrgProjectMember(ctx context.Context, input vo.OrgProjectMemberReq) (*vo.OrgProjectMemberResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.OrgProjectMember(projectvo.OrgProjectMemberReqVo{
		OrgId: cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		ProjectId: input.ProjectID,
	})

	if resp.Failure() {
		return nil, resp.Error()
	}

	result := &vo.OrgProjectMemberResp{}

	copyError := copyer.Copy(resp.OrgProjectMemberRespVo, result)

	if copyError != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyError)
	}

	return result, nil

}
