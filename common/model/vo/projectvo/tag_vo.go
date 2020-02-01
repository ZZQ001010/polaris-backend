package projectvo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type CreateTagReqVo struct {
	UserId int64           `json:"userId"`
	OrgId  int64           `json:"orgId"`
	Input  vo.CreateTagReq `json:"data"`
}

type TagListReqVo struct {
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	OrgId int64          `json:"orgId"`
	Input vo.TagListReq `json:"input"`
}

type TagListRespVo struct {
	vo.Err
	Data *vo.TagList `json:"data"`
}

type TagDefaultStyleRespVo struct {
	vo.Err
	Data []string `json:"data"`
}

type HotTagListReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"input"`
}

type DeleteTagReqVo struct {
	OrgId     int64 `json:"orgId"`
	UserId    int64 `json:"userId"`
	Data vo.DeleteTagReq `json:"data"`
}

type UpdateTagReqVo struct {
	OrgId     int64 `json:"orgId"`
	UserId    int64 `json:"userId"`
	Data vo.UpdateTagReq `json:"data"`
}