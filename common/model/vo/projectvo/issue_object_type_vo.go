package projectvo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type IssueObjectTypeListReqVo struct {
	Page   uint                    `json:"page"`
	Size   uint                    `json:"size"`
	Params *vo.IssueObjectTypesReq `json:"params"`
	//UserId  int64           `json:"userId"`
	OrgId int64 `json:"orgId"`
}

type IssueObjectTypeListRespVo struct {
	vo.Err
	IssueObjectTypeList *vo.IssueObjectTypeList `json:"data"`
}

type CreateIssueObjectTypeReqVo struct {
	Input  vo.CreateIssueObjectTypeReq `json:"input"`
	UserId int64                       `json:"userId"`
	//OrgId   int64           `json:"orgId"`
}

type UpdateIssueObjectTypeReqVo struct {
	Input  vo.UpdateIssueObjectTypeReq `json:"input"`
	UserId int64                       `json:"userId"`
	//OrgId   int64           `json:"orgId"`
}

type DeleteIssueObjectTypeReqVo struct {
	Input  vo.DeleteIssueObjectTypeReq `json:"input"`
	UserId int64                       `json:"userId"`
	OrgId  int64                       `json:"orgId"`
}
