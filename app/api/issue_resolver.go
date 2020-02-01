package api

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
)

func (r *mutationResolver) CreateIssue(ctx context.Context, input vo.CreateIssueReq) (*vo.Issue, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.CreateIssue(projectvo.CreateIssueReqVo{
		CreateIssue:   input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	return respVo.Issue, respVo.Error()
}

func (r *mutationResolver) UpdateIssue(ctx context.Context, input vo.UpdateIssueReq) (*vo.UpdateIssueResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.UpdateIssue(projectvo.UpdateIssueReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	return respVo.UpdateIssue, respVo.Error()
}

func (r *mutationResolver) DeleteIssue(ctx context.Context, input vo.DeleteIssueReq) (*vo.Issue, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.DeleteIssue(projectvo.DeleteIssueReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	return respVo.Issue, respVo.Error()
}

func (r *queryResolver) HomeIssues(ctx context.Context, page int, size int, input *vo.HomeIssueInfoReq) (*vo.HomeIssueInfoResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.HomeIssues(projectvo.HomeIssuesReqVo{
		Page:   page,
		Size:   size,
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.HomeIssueInfo, respVo.Error()
}

func (r *queryResolver) IssueInfo(ctx context.Context, issueID int64) (*vo.IssueInfo, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.IssueInfo(projectvo.IssueInfoReqVo{
		IssueID:       issueID,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	return respVo.IssueInfo, respVo.Error()
}
func (r *queryResolver) IssueRestInfos(ctx context.Context, page int, size int, input *vo.IssueRestInfoReq) (*vo.IssueRestInfoResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.GetIssueRestInfos(projectvo.GetIssueRestInfosReqVo{
		Page:  page,
		Size:  size,
		Input: input,
		OrgId: cacheUserInfo.OrgId,
	})
	return respVo.GetIssueRestInfos, respVo.Error()
}

func (r *mutationResolver) UpdateIssueStatus(ctx context.Context, input vo.UpdateIssueStatusReq) (*vo.Issue, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.UpdateIssueStatus(projectvo.UpdateIssueStatusReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})
	return respVo.Issue, respVo.Error()
}

func (r *mutationResolver) UpdateIssueProjectObjectType(ctx context.Context, input vo.UpdateIssueProjectObjectTypeReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.UpdateIssueProjectObjectType(projectvo.UpdateIssueProjectObjectTypeReqVo{
		Input:         input,
		UserId:        cacheUserInfo.UserId,
		OrgId:         cacheUserInfo.OrgId,
		SourceChannel: cacheUserInfo.SourceChannel,
	})

	return respVo.Void, respVo.Error()
}

func (r *queryResolver) IssueReport(ctx context.Context, reportType int64) (*vo.IssueReportResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.IssueReport(projectvo.IssueReportReqVo{
		ReportType: reportType,
		UserId:     cacheUserInfo.UserId,
		OrgId:      cacheUserInfo.OrgId,
	})
	return respVo.IssueReport, respVo.Error()
}

func (r *queryResolver) IssueReportDetail(ctx context.Context, shareID string) (*vo.IssueReportResp, error) {
	_, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.IssueReportDetail(projectvo.IssueReportDetailReqVo{
		ShareID: shareID,
	})
	return respVo.IssueReportDetail, respVo.Error()
}

func (r *queryResolver) IssueStatusTypeStat(ctx context.Context, input *vo.IssueStatusTypeStatReq) (*vo.IssueStatusTypeStatResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.IssueStatusTypeStat(projectvo.IssueStatusTypeStatReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.IssueStatusTypeStat, respVo.Error()
}

func (r *queryResolver) IssueAssignRank(ctx context.Context, input vo.IssueAssignRankReq) ([]*vo.IssueAssignRankInfo, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	reqVo := projectvo.IssueAssignRankReqVo{
		Input: input,
		OrgId: cacheUserInfo.OrgId,
	}

	respVo := projectfacade.IssueAssignRank(reqVo)
	return respVo.IssueAssignRankResp, respVo.Error()
}

func (r *queryResolver) IssueStatusTypeStatDetail(ctx context.Context, input *vo.IssueStatusTypeStatReq) (*vo.IssueStatusTypeStatDetailResp, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.IssueStatusTypeStatDetail(projectvo.IssueStatusTypeStatReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.IssueStatusTypeStatDetail, respVo.Error()
}

func (r *mutationResolver) CreateIssueComment(ctx context.Context, input vo.CreateIssueCommentReq) (*vo.Void, error) {

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.CreateIssueComment(projectvo.CreateIssueCommentReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.Void, respVo.Error()
}

func (r *mutationResolver) CreateIssueResource(ctx context.Context, input vo.CreateIssueResourceReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.CreateIssueResource(projectvo.CreateIssueResourceReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.Void, respVo.Error()
}

func (r *mutationResolver) UpdateIssueAndIssueRelate(ctx context.Context, input vo.UpdateIssueAndIssueRelateReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.CreateIssueRelationIssue(projectvo.CreateIssueRelationIssueReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.Void, respVo.Error()
}

func (r *mutationResolver) DeleteIssueResource(ctx context.Context, input vo.DeleteIssueResourceReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.DeleteIssueResource(projectvo.DeleteIssueResourceReqVo{
		Input:  input,
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
	})
	return respVo.Void, respVo.Error()
}

func (r *queryResolver) RelatedIssueList(ctx context.Context, input vo.RelatedIssueListReq) (*vo.IssueRestInfoResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.RelatedIssueList(projectvo.RelatedIssueListReqVo{
		Input: input,
		OrgId: cacheUserInfo.OrgId,
	})
	return respVo.RelatedIssueList, respVo.Error()
}

func (r *queryResolver) IssueResources(ctx context.Context, page *int, size *int, input *vo.GetIssueResourcesReq) (*vo.ResourceList, error) {
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

	respVo := projectfacade.IssueResources(projectvo.IssueResourcesReqVo{
		Page:  pageA,
		Size:  sizeA,
		Input: input,
		OrgId: cacheUserInfo.OrgId,
	})
	return respVo.IssueResources, respVo.Error()
}

func (r *mutationResolver) ImportIssues(ctx context.Context, input vo.ImportIssuesReq) (*vo.Void, error) {
	panic("implement me")
}

func (r *queryResolver) ExportIssueTemplate(ctx context.Context, projectID int64) (*vo.ExportIssueTemplateResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.ExportIssueTemplate(projectvo.ExportIssueTemplateReqVo{
		OrgId:     cacheUserInfo.OrgId,
		ProjectId: projectID,
	})

	return respVo.Data, respVo.Error()
}

func (r *queryResolver) ExportData(ctx context.Context, projectID int64) (*vo.ExportIssueTemplateResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.ExportData(projectvo.ExportIssueTemplateReqVo{
		OrgId:     cacheUserInfo.OrgId,
		ProjectId: projectID,
	})

	return respVo.Data, respVo.Error()
}

func (r *mutationResolver) UpdateIssueSort(ctx context.Context, input vo.UpdateIssueSortReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := projectfacade.UpdateIssueSort(projectvo.UpdateIssueSortReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Void, respVo.Error()
}

func (r *queryResolver) IssueDailyPersonalWorkCompletionStat(ctx context.Context, input *vo.IssueDailyPersonalWorkCompletionStatReq) (*vo.IssueDailyPersonalWorkCompletionStatResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.IssueDailyPersonalWorkCompletionStat(projectvo.IssueDailyPersonalWorkCompletionStatReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}

func (r *queryResolver) IssueAndProjectCountStat(ctx context.Context) (*vo.IssueAndProjectCountStatResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.IssueAndProjectCountStat(projectvo.IssueAndProjectCountStatReqVo{
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Data, respVo.Error()
}

func (r *mutationResolver) UpdateIssueTags(ctx context.Context, input vo.UpdateIssueTagsReq) (*vo.Void, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.CreateIssueRelationTags(projectvo.CreateIssueRelationTagsReqVo{
		Input:  input,
		OrgId:  cacheUserInfo.OrgId,
		UserId: cacheUserInfo.UserId,
	})
	return respVo.Void, respVo.Error()
}
