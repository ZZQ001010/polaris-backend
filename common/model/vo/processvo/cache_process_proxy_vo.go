package processvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type GetNextProcessStepStatusListReqVo struct {
	OrgId         int64 `json:"orgId"`
	ProcessId     int64 `json:"processId"`
	StartStatusId int64 `json:"startStatusId"`
}

type GetNextProcessStepStatusListRespVo struct {
	vo.Err
	CacheProcessStatus *[]bo.CacheProcessStatusBo `json:"data"`
}
