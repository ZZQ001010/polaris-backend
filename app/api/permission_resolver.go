package api

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
)

func (r *queryResolver) PermissionOperationList(ctx context.Context, roleID int64, projectID *int64) ([]*vo.PermissionOperationListResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := rolefacade.PermissionOperationList(rolevo.PermissionOperationListReqVo{
		OrgId:     cacheUserInfo.OrgId,
		RoleId:    roleID,
		UserId:    cacheUserInfo.UserId,
		ProjectId: projectID,
	})

	return resp.Data, resp.Error()
}

func (r *mutationResolver) UpdateRolePermissionOperation(ctx context.Context, input vo.UpdateRolePermissionOperationReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := rolefacade.UpdateRolePermissionOperation(rolevo.UpdateRolePermissionOperationReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  input,
	})

	return resp.Void, resp.Error()
}

func (r *queryResolver) GetPersonalPermissionInfo(ctx context.Context, projectID *int64, issueID *int64) (*vo.GetPersonalPermissionInfoResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := rolefacade.GetPersonalPermissionInfo(rolevo.GetPersonalPermissionInfoReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		ProjectId:projectID,
		IssueId:issueID,
		SourceChannel:cacheUserInfo.SourceChannel,
	})
	return &vo.GetPersonalPermissionInfoResp{Data:resp.Data}, resp.Error()
}