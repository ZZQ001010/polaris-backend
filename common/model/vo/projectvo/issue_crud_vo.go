package projectvo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type CreateIssueReqVo struct {
	CreateIssue   vo.CreateIssueReq `json:"createIssue"`
	UserId        int64             `json:"userId"`
	OrgId         int64             `json:"orgId"`
	SourceChannel string            `json:"sourceChannel"`
}

type IssueRespVo struct {
	vo.Err
	Issue *vo.Issue `json:"data"`
}

type UpdateIssueReqVo struct {
	Input         vo.UpdateIssueReq `json:"input"`
	UserId        int64             `json:"userId"`
	OrgId         int64             `json:"orgId"`
	SourceChannel string            `json:"sourceChannel"`
}

type UpdateIssueRespVo struct {
	vo.Err
	UpdateIssue *vo.UpdateIssueResp `json:"data"`
}

type IssueInfoReqVo struct {
	IssueID int64 `json:"issueId"`
	UserId  int64 `json:"userId"`
	OrgId   int64 `json:"orgId"`
	SourceChannel string `json:"sourceChannel"`
}

type IssueInfoRespVo struct {
	vo.Err
	IssueInfo *vo.IssueInfo `json:"data"`
}

type GetIssueRestInfosReqVo struct {
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
	Input *vo.IssueRestInfoReq `json:"input"`
	OrgId int64                `json:"orgId"`
}

type GetIssueRestInfosRespVo struct {
	vo.Err
	GetIssueRestInfos *vo.IssueRestInfoResp `json:"data"`
}

type DeleteIssueReqVo struct {
	Input         vo.DeleteIssueReq `json:"input"`
	UserId        int64             `json:"userId"`
	OrgId         int64             `json:"orgId"`
	SourceChannel string            `json:"sourceChannel"`
}

type UpdateIssueStatusReqVo struct {
	Input         vo.UpdateIssueStatusReq `json:"input"`
	UserId        int64                   `json:"userId"`
	OrgId         int64                   `json:"orgId"`
	SourceChannel string                  `json:"sourceChannel"`
}

type UpdateIssueProjectObjectTypeReqVo struct {
	Input         vo.UpdateIssueProjectObjectTypeReq `json:"input"`
	UserId        int64                              `json:"userId"`
	OrgId         int64                              `json:"orgId"`
	SourceChannel string                             `json:"sourceChannel"`
}

type LarkIssueInitReqVo struct {
	OrgId      int64 `json:"orgId"`
	ZhangsanId int64 `json:"zhangsanId"`
	LisiId     int64 `json:"lisiId"`
	ProjectId  int64 `json:"projectId"`
	OperatorId int64 `json:"operatorId"`
}

type IssueInfoListReqVo struct {
	IssueIds []int64 `json:"issueIds"`
}

type IssueInfoListRespVo struct {
	vo.Err
	IssueInfos []vo.Issue `json:"data"`
}

type UpdateIssueSortReqVo struct {
	Input  vo.UpdateIssueSortReq `json:"input"`
	UserId int64                 `json:"userId"`
	OrgId  int64                 `json:"orgId"`
}

type DailyProjectIssueStatisticsCountRespVo struct {
	//完成数量
	DailyFinishCount int `json:"dailyFinishCount"`
	//剩余未完成
	RemainingCount int `json:"dailyRemainingCount"`
	//逾期任务数量
	OverdueCount int `json:"dailyOverdueCount"`
}

type DailyProjectIssueStatisticsRespVo struct {
	vo.Err
	DailyProjectIssueStatisticsCountRespVo *DailyProjectIssueStatisticsCountRespVo `json:"data"`
}

type DailyProjectIssueStatisticsReqVo struct {
	ProjectId int64 `json:"projectId"`
	OrgId     int64 `json:"orgId"`
}
