package processvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"upper.io/db.v3"
)

type InitProcessReqVo struct {
	OrgId int64 `json:"orgId"`
}

type AssignValueToFieldReqVo struct {
	ProcessRes *map[string]int64 `json:"processRes"`
	OrgId      int64             `json:"orgId"`
}

type GetProcessByLangCodeReqVo struct {
	OrgId    int64  `json:"orgId"`
	LangCode string `json:"langCode"`
}

type GetProcessByLangCodeRespVo struct {
	ProcessBo *bo.ProcessBo `json:"data"`

	vo.Err
}

type GetProcessBoReqVo struct {
	Cond db.Cond `json:"cond"`
}

type GetProcessBoRespVo struct {
	ProcessBo bo.ProcessBo `json:"data"`

	vo.Err
}

type GetProcessByIdReqVo struct {
	OrgId int64 `json:"orgId"`
	Id int64 `json:"id"`
}

type GetProcessByIdRespVo struct {
	ProcessBo *bo.ProcessBo `json:"data"`
	vo.Err
}