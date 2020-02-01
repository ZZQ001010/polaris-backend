package projectvo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type DeleteProjectAttachmentReqVo struct {
	Input  vo.DeleteProjectAttachmentReq `json:"input"`
	UserId int64                         `json:"userId"`
	OrgId  int64                         `json:"orgId"`
}

type DeleteProjectAttachmentRespVo struct {
	Output *vo.DeleteProjectAttachmentResp `json:"input"`
	vo.Err
}

type GetProjectAttachmentReqVo struct {
	Input  vo.ProjectAttachmentReq `json:"input"`
	UserId int64                   `json:"userId"`
	OrgId  int64                   `json:"orgId"`
	Page   int                     `json:"page"`
	Size   int                     `json:"size"`
}

type GetProjectAttachmentRespVo struct {
	vo.Err
	Output *vo.AttachmentList `json:"data"`
}
