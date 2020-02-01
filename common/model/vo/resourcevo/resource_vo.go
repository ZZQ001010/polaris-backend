package resourcevo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type CreateResourceReqVo struct {
	CreateResourceBo bo.CreateResourceBo `json:"createResourceBo"`
}

type CreateResourceRespVo struct {
	ResourceId int64 `json:"data"`

	vo.Err
}

type UpdateResourceInfoReqVo struct {
	Input bo.UpdateResourceInfoBo `json:"updateResourceBo"`
}
type UpdateResourceInfoResVo struct {
	*UpdateResourceData
	vo.Err
}
type UpdateResourceData struct {
	OldBo             []bo.ResourceBo
	NewBo             []bo.ResourceBo
	CurrentFolderName *string
	TargetFolderName  *string
}

type UpdateResourceFolderReqVo struct {
	Input bo.UpdateResourceFolderBo `json:"updateResourceBo"`
}

type DeleteResourceReqVo struct {
	Input bo.DeleteResourceBo `json:"deleteResourceBo"`
}

type GetResourceReqVo struct {
	Input bo.GetResourceBo `json:"getResourceBo"`
}
type InsertResourceReqVo struct {
	Input InsertResourceReqData `json:"input"`
}

type InsertResourceReqData struct {
	ResourcePath  string `json:"resourcePath"`
	OrgId         int64  `json:"orgId"`
	CurrentUserId int64  `json:"currentUserId"`
	ResourceType  int    `json:"resourceType"`
	FileName      string `json:"fileName"`
}

type InsertResourceRespVo struct {
	ResourceId int64 `json:"data"`

	vo.Err
}

type GetResourceByIdReqVo struct {
	GetResourceByIdReqBody GetResourceByIdReqBody `json:"getResourceByIdReqBody"`
}

type GetResourceByIdReqBody struct {
	ResourceIds []int64 `json:"resourceIds"`
}

type GetResourceByIdRespVo struct {
	ResourceBos []bo.ResourceBo `json:"data"`

	vo.Err
}

type GetIdByPathReqVo struct {
	OrgId        int64  `json:"orgId"`
	ResourcePath string `json:"resourcePath"`
	ResourceType int    `json:"resourceType"`
}

type GetIdByPathRespVo struct {
	ResourceId int64 `json:"data"`

	vo.Err
}

type GetResourceBoListReqVo struct {
	Page  uint                  `json:"page"`
	Size  uint                  `json:"size"`
	Input GetResourceBoListCond `json:"cond"`
}

type GetResourceBoListCond struct {
	OrgId       int64    `json:"orgId"`
	ResourceIds *[]int64 `json:"resourceIds"`
	IsDelete    *int     `json:"isDelete"`
}

type GetResourceBoListRespVo struct {
	GetResourceBoListRespData `json:"data"`

	vo.Err
}

type GetResourceBoListRespData struct {
	ResourceBos *[]bo.ResourceBo `json:"resourceBos"`
	Total       int64            `json:"total"`
}

type GetResourceVoListRespVo struct {
	*vo.ResourceList `json:"data"`

	vo.Err
}

type CreateIssueResourceReqVo struct {
	Input  vo.CreateProjectResourceReq `json:"input"`
	UserId int64                       `json:"userId"`
	OrgId  int64                       `json:"orgId"`
}
