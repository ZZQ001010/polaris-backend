package processvo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type GetDefaultProcessIdReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProcessId int64 `json:"processId"`
	Category  int   `json:"category"`
}

type GetDefaultProcessIdRespVo struct {
	vo.Err
	ProcessId int64 `json:"data"`
}
