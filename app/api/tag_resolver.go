package api

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
)

func (r *queryResolver) TagList(ctx context.Context, page *int, size *int, params vo.TagListReq) (*vo.TagList, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	defaultPage := 0
	defaultSize := 5
	if page == nil || *page < 0 {
		page = &defaultPage
	}

	if size == nil || *size <= 0 {
		size = &defaultSize
	}

	resp := projectfacade.TagList(projectvo.TagListReqVo{
		OrgId: cacheUserInfo.OrgId,
		Page:  *page,
		Size:  *size,
		Input: params,
	})

	return resp.Data, resp.Error()
}

func (r *mutationResolver) CreateTag(ctx context.Context, input vo.CreateTagReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.CreateTag(projectvo.CreateTagReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Input:  input,
	})

	return resp.Void, resp.Error()
}

func (r *queryResolver) TagDefaultStyle(ctx context.Context) (*vo.StypeList, error) {
	resp := projectfacade.TagDefaultStyle()
	return &vo.StypeList{List: resp.Data}, nil
}

func (r *mutationResolver) DeleteTag(ctx context.Context, input vo.DeleteTagReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.DeleteTag(projectvo.DeleteTagReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Data:  input,
	})

	return resp.Void, resp.Error()
}

func (r *mutationResolver) UpdateTag(ctx context.Context, input vo.UpdateTagReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.UpdateTag(projectvo.UpdateTagReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
		Data:  input,
	})

	return resp.Void, resp.Error()
}

func (r *queryResolver) HotTagList(ctx context.Context, projectID int64) (*vo.TagList, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	resp := projectfacade.HotTagList(projectvo.HotTagListReqVo{
		OrgId:  cacheUserInfo.OrgId,
		ProjectId:projectID,
	})

	return resp.Data, resp.Error()
}