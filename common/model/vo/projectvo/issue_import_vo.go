package projectvo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type ImportIssuesReqVo struct {
	UserId int64              `json:"userId"`
	OrgId  int64              `json:"orgId"`
	Input  vo.ImportIssuesReq `json:"data"`
}

type ExportIssueTemplateReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type ExportIssueTemplateRespVo struct {
	vo.Err
	Data *vo.ExportIssueTemplateResp `json:"data"`
}

type ImportIssuesRespVo struct {
	vo.Err
	Data int64 `json:"data"`
}
