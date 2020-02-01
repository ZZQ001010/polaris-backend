package projectvo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type ProjectTypesRespVo struct {
	vo.Err
	ProjectTypes []*vo.ProjectType `json:"data"`
}

type ProjectTypesReqVo struct {
	OrgId int64 `json:"orgId"`
}
