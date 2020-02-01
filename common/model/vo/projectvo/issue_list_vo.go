package projectvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type HomeIssuesReqVo struct {
	Page   int                  `json:"page"`
	Size   int                  `json:"size"`
	Input  *vo.HomeIssueInfoReq `json:"input"`
	UserId int64                `json:"userId"`
	OrgId  int64                `json:"orgId"`
}

type HomeIssuesRespVo struct {
	vo.Err
	HomeIssueInfo *vo.HomeIssueInfoResp `json:"data"`
}

type IssueReportReqVo struct {
	ReportType int64 `json:"reportType"`
	UserId     int64 `json:"userId"`
	OrgId      int64 `json:"orgId"`
}

type IssueReportRespVo struct {
	vo.Err
	IssueReport *vo.IssueReportResp `json:"data"`
}

type IssueReportDetailReqVo struct {
	ShareID string `json:"shareId"`
}

type IssueReportDetailRespVo struct {
	vo.Err
	IssueReportDetail *vo.IssueReportResp `json:"data"`
}

type IssueStatusTypeStatRespVo struct {
	vo.Err
	IssueStatusTypeStat *vo.IssueStatusTypeStatResp `json:"data"`
}

type IssueStatusTypeStatDetailRespVo struct {
	vo.Err
	IssueStatusTypeStatDetail *vo.IssueStatusTypeStatDetailResp `json:"data"`
}

type IssueStatusTypeStatReqVo struct {
	Input  *vo.IssueStatusTypeStatReq `json:"input"`
	UserId int64                      `json:"userId"`
	OrgId  int64                      `json:"orgId"`
}

type GetIssueStatusStatReqVo struct {
	Input bo.IssueStatusStatCondBo `json:"input"`
}

type GetIssueStatusStatRespVo struct {
	vo.Err
	Data GetIssueStatusStatRespData `json:"data"`
}

type GetIssueStatusStatRespData struct {
	List []bo.IssueStatusStatBo `json:"list"`
}

type IssueStatusTypeStatDetailReqVo struct {
	Input *vo.IssueStatusTypeStatReq `json:"input"`
}

type GetSimpleIssueInfoBatchReqVo struct {
	OrgId int64   `json:"orgId"`
	Ids   []int64 `json:"ids"`
}

type GetSimpleIssueInfoBatchRespVo struct {
	vo.Err
	Data *[]vo.Issue `json:"data"`
}

type GetIssueRemindInfoListReqVo struct {
	Page int `json:"page"`
	Size int `json:"size"`
	Input GetIssueRemindInfoListReqData `json:"input"`
}

type GetIssueRemindInfoListRespVo struct {
	Data *GetIssueRemindInfoListRespData `json:"data"`
	vo.Err
}

type GetIssueRemindInfoListRespData struct {
	Total int64 `json:"total"`
	List []bo.IssueRemindInfoBo `json:"issueRemindInfoList"`
}

type GetIssueRemindInfoListReqData struct {
	//计划开始时间开始范围
	BeforePlanEndTime *string `json:"beforePlanEndTime"`
	//计划开始时间结束范围
	AfterPlanEndTime *string `json:"afterPlanEndTime"`
}
