package api

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
)

func (r *queryResolver) OrgRoleList(ctx context.Context) ([]*vo.Role, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := rolefacade.GetOrgRoleList(rolevo.GetOrgRoleListReqVo{
		OrgId: cacheUserInfo.OrgId,
	})
	return resp.Data, resp.Error()
}

func (r *mutationResolver) CreateRole(ctx context.Context, input vo.CreateRoleReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := rolefacade.CreateRole(rolevo.CreateOrgReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  input,
	})

	return resp.Void, resp.Error()
}

func (r *mutationResolver) DelRole(ctx context.Context, input vo.DelRoleReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := rolefacade.DelRole(rolevo.DelRoleReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  input,
	})

	return resp.Void, resp.Error()

}

func (r *mutationResolver) UpdateRole(ctx context.Context, input vo.UpdateRoleReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	resp := rolefacade.UpdateRole(rolevo.UpdateRoleReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  input,
	})

	return resp.Void, resp.Error()
}

func (r *queryResolver) ProjectRoleList(ctx context.Context, projectID int64) ([]*vo.Role, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := rolefacade.GetProjectRoleList(rolevo.GetProjectRoleListReqVo{
		OrgId:     cacheUserInfo.OrgId,
		ProjectId: projectID,
	})
	return resp.Data, resp.Error()
}
