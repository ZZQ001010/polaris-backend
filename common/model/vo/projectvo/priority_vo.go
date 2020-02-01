package projectvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type PriorityListRespVo struct {
	vo.Err
	PriorityList *vo.PriorityList `json:"data"`
}

type CreatePriorityReqVo struct {
	CreatePriorityReq vo.CreatePriorityReq `json:"createPriorityReq"`
	UserId            int64                `json:"userId"`
}

type UpdatePriorityReqVo struct {
	UpdatePriorityReq vo.UpdatePriorityReq `json:"updatePriorityReq"`
	UserId            int64                `json:"userId"`
}

type DeletePriorityReqVo struct {
	DeletePriorityReq vo.DeletePriorityReq `json:"deletePriorityReq"`
	UserId            int64                `json:"userId"`
	OrgId             int64                `json:"orgId"`
}

type VerifyPriorityReqVo struct {
	OrgId      int64 `json:"orgId"`
	Typ        int   `json:"typ"`
	BeVerifyId int64 `json:"beVerifyId"`
}

type VerifyPriorityRespVo struct {
	Successful bool `json:"data"`

	vo.Err
}

type InitPriorityReqVo struct {
	OrgId int64 `json:"orgId"`
}

type GetPriorityByIdReqVo struct {
	OrgId int64 `json:"orgId"`
	Id    int64 `json:"id"`
}

type GetPriorityByIdRespVo struct {
	PriorityBo *bo.PriorityBo `json:"data"`
	vo.Err
}

type PriorityListReqVo struct {
	Page  *int  `json:"page"`
	Size  *int  `json:"size"`
	Type  *int  `json:"type"`
	OrgId int64 `json:"orgId"`
}
