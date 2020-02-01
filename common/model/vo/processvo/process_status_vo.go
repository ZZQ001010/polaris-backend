package processvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type ProcessStatusListRespVo struct {
	vo.Err
	ProcessStatusList *vo.ProcessStatusList `json:"data"`
}

type CreateProcessStatusReqVo struct {
	CreateProcessStatusReq vo.CreateProcessStatusReq `json:"createProcessStatusReq"`
	UserId                 int64                     `json:"userId"`
	OrgId                  int64                     `json:"orgId"`
}

type UpdateProcessStatusReqVo struct {
	UpdateProcessStatusReq vo.UpdateProcessStatusReq `json:"updateProcessStatusReq"`
	UserId                 int64                     `json:"userId"`
	OrgId                  int64                     `json:"orgId"`
}

type DeleteProcessStatusReq struct {
	DeleteProcessStatusReq vo.DeleteProcessStatusReq `json:"deleteProcessStatusReq"`
	UserId                 int64                     `json:"userId"`
	OrgId                  int64                     `json:"orgId"`
}

type GetProcessStatusReqVo struct {
	OrgId int64 `json:"orgId"`
	Id    int64 `json:"id"`
}

type GetProcessStatusRespVo struct {
	CacheProcessStatusBo *bo.CacheProcessStatusBo `json:"data"`
	vo.Err
}

type ProcessStatusInitReqVo struct {
	OrgId      int64                  `json:"orgId"`
	ContextMap map[string]interface{} `json:"contextMap"`
}

type GetProcessStatusByCategoryReqVo struct {
	OrgId    int64 `json:"orgId"`
	StatusId int64 `json:"statusId"`
	Category int   `json:"category"`
}

type GetProcessStatusByCategoryRespVo struct {
	CacheProcessStatusBo *bo.CacheProcessStatusBo `json:"data"`

	vo.Err
}

type GetProcessStatusListByCategoryReqVo struct {
	OrgId      int64 `json:"orgId"`
	Category   int   `json:"category"`
	StatusType int   `json:"statusType"`
}

type GetProcessStatusListByCategoryRespVo struct {
	CacheProcessStatusBoList []bo.CacheProcessStatusBo `json:"data"`

	vo.Err
}

type GetProcessStatusIdsReqVo struct {
	OrgId    int64 `json:"orgId"`
	Category int   `json:"category"`
	Typ      int   `json:"typ"`
}

type GetProcessStatusIdsRespVo struct {
	ProcessStatusIds *[]int64 `json:"data"`

	vo.Err
}

type GetProcessStatusListReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProcessId int64 `json:"processId"`
}

type GetProcessStatusListRespVo struct {
	ProcessStatusBoList *[]bo.CacheProcessStatusBo `json:"data"`

	vo.Err
}

type GetProcessInitStatusIdReqVo struct {
	OrgId               int64 `json:"orgId"`
	ProjectId           int64 `json:"projectId"`
	ProjectObjectTypeId int64 `json:"projectObjectTypeId"`
	Category            int   `json:"category"`
}

type GetProcessInitStatusIdRespVo struct {
	ProcessInitStatusId int64 `json:"data"`

	vo.Err
}
