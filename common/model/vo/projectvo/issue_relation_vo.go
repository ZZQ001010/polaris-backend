package projectvo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type RelatedIssueListRespVo struct {
	vo.Err
	RelatedIssueList *vo.IssueRestInfoResp `json:"relatedIssueList"`
}

type CreateIssueRelationIssueReqVo struct {
	Input  vo.UpdateIssueAndIssueRelateReq `json:"input"`
	UserId int64                           `json:"userId"`
	OrgId  int64                           `json:"orgId"`
}

type CreateIssueRelationTagsReqVo struct {
	Input  vo.UpdateIssueTagsReq 			`json:"input"`
	UserId int64                           `json:"userId"`
	OrgId  int64                           `json:"orgId"`
}

type IssueResourcesReqVo struct {
	Page  uint                     `json:"page"`
	Size  uint                     `json:"size"`
	Input *vo.GetIssueResourcesReq `json:"input"`
	OrgId int64                    `json:"orgId"`
}

type IssueResourcesRespVo struct {
	vo.Err
	IssueResources *vo.ResourceList `json:"issueResources"`
}

type CreateIssueCommentReqVo struct {
	Input  vo.CreateIssueCommentReq `json:"input"`
	UserId int64                    `json:"userId"`
	OrgId  int64                    `json:"orgId"`
}

type CreateIssueResourceReqVo struct {
	Input  vo.CreateIssueResourceReq `json:"input"`
	UserId int64                     `json:"userId"`
	OrgId  int64                     `json:"orgId"`
}

type DeleteIssueResourceReqVo struct {
	Input  vo.DeleteIssueResourceReq `json:"input"`
	UserId int64                     `json:"userId"`
	OrgId  int64                     `json:"orgId"`
}

type RelatedIssueListReqVo struct {
	Input vo.RelatedIssueListReq `json:"input"`
	OrgId int64                  `json:"orgId"`
}

type GetIssueMembersReqVo struct {
	IssueId int64 `json:"issueId"`
	OrgId int64 `json:"orgId"`
}

type GetIssueMembersRespVo struct {
	Data GetIssueMembersRespData `json:"data"`

	vo.Err
}

type GetIssueMembersRespData struct {
	MemberIds []int64 `json:"memberIds"`
	OwnerId int64 `json:"ownerId"`
	ParticipantIds []int64 `json:"participantIds"`
	FollowerIds []int64 `json:"followerIds"`
}