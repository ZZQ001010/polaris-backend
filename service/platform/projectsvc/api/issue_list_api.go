package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) HomeIssues(reqVo projectvo.HomeIssuesReqVo) projectvo.HomeIssuesRespVo {
	res, err := service.HomeIssues(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.HomeIssuesRespVo{Err: vo.NewErr(err), HomeIssueInfo: res}
}

func (GetGreeter) IssueReport(reqVo projectvo.IssueReportReqVo) projectvo.IssueReportRespVo {
	res, err := service.IssueReport(reqVo.OrgId, reqVo.UserId, reqVo.ReportType)
	return projectvo.IssueReportRespVo{Err: vo.NewErr(err), IssueReport: res}
}

func (GetGreeter) IssueReportDetail(reqVo projectvo.IssueReportDetailReqVo) projectvo.IssueReportDetailRespVo {
	res, err := service.IssueReportDetail(reqVo.ShareID)
	return projectvo.IssueReportDetailRespVo{Err: vo.NewErr(err), IssueReportDetail: res}
}

func (PostGreeter) IssueStatusTypeStat(reqVo projectvo.IssueStatusTypeStatReqVo) projectvo.IssueStatusTypeStatRespVo {
	res, err := service.IssueStatusTypeStat(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.IssueStatusTypeStatRespVo{Err: vo.NewErr(err), IssueStatusTypeStat: res}
}

func (PostGreeter) IssueStatusTypeStatDetail(reqVo projectvo.IssueStatusTypeStatReqVo) projectvo.IssueStatusTypeStatDetailRespVo {
	res, err := service.IssueStatusTypeStatDetail(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.IssueStatusTypeStatDetailRespVo{Err: vo.NewErr(err), IssueStatusTypeStatDetail: res}
}

func (PostGreeter) GetSimpleIssueInfoBatch(reqVo projectvo.GetSimpleIssueInfoBatchReqVo) projectvo.GetSimpleIssueInfoBatchRespVo {
	res, err := service.GetSimpleIssueInfoBatch(reqVo.OrgId, reqVo.Ids)
	return projectvo.GetSimpleIssueInfoBatchRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) GetIssueRemindInfoList(reqVo projectvo.GetIssueRemindInfoListReqVo) projectvo.GetIssueRemindInfoListRespVo{
	res, err := service.GetIssueRemindInfoList(reqVo)
	return projectvo.GetIssueRemindInfoListRespVo{Err: vo.NewErr(err), Data: res}
}
