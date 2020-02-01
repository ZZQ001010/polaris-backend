package api

import (
	"context"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/core/util/validator"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"strings"
	"time"
)

func (r *queryResolver) Iterations(ctx context.Context, page *int, size *int, params *vo.IterationListReq) (*vo.IterationList, error) {
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

	respVo := projectfacade.IterationList(projectvo.IterationListReqVo{
		Page:  pageA,
		Size:  sizeA,
		Input: params,
		OrgId: cacheUserInfo.OrgId,
	})
	return respVo.IterationList, respVo.Error()
}

func (r *mutationResolver) CreateIteration(ctx context.Context, input vo.CreateIterationReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}
	if input.PlanStartTime.IsNull() || input.PlanEndTime.IsNull() {
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "开始和结束时间不能为空")
	}
	if time.Time(input.PlanEndTime).Before(time.Time(input.PlanStartTime)) {
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "开始时间应该在结束时间之前")
	}

	name := strings.Trim(input.Name, " ")
	if name == "" || strs.Len(name) > 200 {
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.ReqParamsValidateError, "迭代名称不能为空且限制在200字以内")
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.CreateIteration(projectvo.CreateIterationReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Void, respVo.Error()
}

func (r *mutationResolver) UpdateIteration(ctx context.Context, input vo.UpdateIterationReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.UpdateIteration(projectvo.UpdateIterationReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Void, respVo.Error()
}

func (r *mutationResolver) DeleteIteration(ctx context.Context, input vo.DeleteIterationReq) (*vo.Void, error) {
	validate, err := validator.Validate(input)
	if !validate {
		return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError, err)
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.DeleteIteration(projectvo.DeleteIterationReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Void, respVo.Error()
}

func (r *queryResolver) IterationStatusTypeStat(ctx context.Context, input *vo.IterationStatusTypeStatReq) (*vo.IterationStatusTypeStatResp, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.IterationStatusTypeStat(projectvo.IterationStatusTypeStatReqVo{
		Input: input,
		OrgId: cacheUserInfo.OrgId,
	})
	return respVo.IterationStatusTypeStat, respVo.Error()
}

func (r *mutationResolver) UpdateIterationIssueRelate(ctx context.Context, input vo.IterationIssueRealtionReq) (*vo.Void, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.IterationIssueRelate(projectvo.IterationIssueRelateReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Void, respVo.Error()
}

func (r *mutationResolver) UpdateIterationStatus(ctx context.Context, input vo.UpdateIterationStatusReq) (*vo.Void, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.UpdateIterationStatus(projectvo.UpdateIterationStatusReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Void, respVo.Error()
}

func (r *queryResolver) IterationInfo(ctx context.Context, input vo.IterationInfoReq) (*vo.IterationInfoResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.IterationInfo(projectvo.IterationInfoReqVo{
		Input: input,
		OrgId: cacheUserInfo.OrgId,
	})
	return respVo.IterationInfo, respVo.Error()
}
