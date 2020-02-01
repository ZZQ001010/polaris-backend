package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) IssueAssignRank(reqVo projectvo.IssueAssignRankReqVo) projectvo.IssueAssignRankRespVo {
	input := reqVo.Input
	projectId := input.ProjectID
	rankTop := 5
	if input.RankTop != nil {
		rt := *input.RankTop
		if rt >= 1 && rt <= 100 {
			rankTop = rt
		}
	}
	res, err := service.IssueAssignRank(reqVo.OrgId, projectId, rankTop)
	return projectvo.IssueAssignRankRespVo{Err: vo.NewErr(err), IssueAssignRankResp: res}
}


func (GetGreeter) IssueAndProjectCountStat(reqVo projectvo.IssueAndProjectCountStatReqVo) projectvo.IssueAndProjectCountStatRespVo {
	res, err := service.IssueAndProjectCountStat(reqVo)
	return projectvo.IssueAndProjectCountStatRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) IssueDailyPersonalWorkCompletionStat(reqVo projectvo.IssueDailyPersonalWorkCompletionStatReqVo) projectvo.IssueDailyPersonalWorkCompletionStatRespVo {
	res, err := service.IssueDailyPersonalWorkCompletionStat(reqVo)
	return projectvo.IssueDailyPersonalWorkCompletionStatRespVo{Err: vo.NewErr(err), Data: res}
}