package mqbo

import "github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"

type PushCreateIssueBo struct {
	CreateIssueReqVo projectvo.CreateIssueReqVo `json:"createIssueReqVo"`
	IssueId int64 `json:"issueId"`
}