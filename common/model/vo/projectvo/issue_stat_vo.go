package projectvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type StatisticDailyTaskCompletionProgressReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type IssueAssignRankReqVo struct {
	Input vo.IssueAssignRankReq
	OrgId int64
}

type IssueAndProjectCountStatReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type IssueDailyPersonalWorkCompletionStatReqVo struct {
	Input *vo.IssueDailyPersonalWorkCompletionStatReq `json:"input"`
	
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type StatisticDailyTaskCompletionProgressRespVo struct {
	IssueDailyNoticeBo *bo.IssueDailyNoticeBo `json:"data"`
	vo.Err
}

type IssueAssignRankRespVo struct {
	IssueAssignRankResp []*vo.IssueAssignRankInfo `json:"data"`
	vo.Err
}

type IssueAndProjectCountStatRespVo struct {
	Data *vo.IssueAndProjectCountStatResp `json:"data"`
	
	vo.Err
}

type IssueDailyPersonalWorkCompletionStatRespVo struct {
	Data *vo.IssueDailyPersonalWorkCompletionStatResp `json:"data"`

	vo.Err
}