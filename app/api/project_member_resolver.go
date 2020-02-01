package api

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
)

func (r *queryResolver) ProjectUserList(ctx context.Context, page *int, size *int, input vo.ProjectUserListReq) (*vo.ProjectUserListResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	defaultPage := 1
	defaultSize := 20
	if page == nil || *page < 1 {
		page = &defaultPage
	}

	if size == nil || *size <= 0 {
		size = &defaultSize
	}
	resp := projectfacade.ProjectUserList(projectvo.ProjectUserListReq{
		Size:  *size,
		Page:  *page,
		OrgId: cacheUserInfo.OrgId,
		Input: input,
	})

	return resp.Data, resp.Error()
}

func (r *mutationResolver) RemoveProjectMember(ctx context.Context, input vo.RemoveProjectMemberReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.RemoveProjectMember(projectvo.RemoveProjectMemberReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  input,
	})

	return resp.Void, resp.Error()
}

func (r *mutationResolver) AddProjectMember(ctx context.Context, input vo.RemoveProjectMemberReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.AddProjectMember(projectvo.RemoveProjectMemberReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  input,
	})

	return resp.Void, resp.Error()
}
