package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) ImportIssues(reqVo projectvo.ImportIssuesReqVo) projectvo.ImportIssuesRespVo {
	count, err := service.ImportIssues(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.ImportIssuesRespVo{Err: vo.NewErr(err), Data: count}
}

func (GetGreeter) ExportIssueTemplate(reqVo projectvo.ExportIssueTemplateReqVo) projectvo.ExportIssueTemplateRespVo {
	url, err := service.ExportIssueTemplate(reqVo.OrgId, reqVo.ProjectId)
	return projectvo.ExportIssueTemplateRespVo{Data: &vo.ExportIssueTemplateResp{URL: url}, Err: vo.NewErr(err)}
}

func (GetGreeter) ExportData(reqVo projectvo.ExportIssueTemplateReqVo) projectvo.ExportIssueTemplateRespVo {
	url, err := service.ExportData(reqVo.OrgId, reqVo.ProjectId)
	return projectvo.ExportIssueTemplateRespVo{Data: &vo.ExportIssueTemplateResp{URL: url}, Err: vo.NewErr(err)}
}
