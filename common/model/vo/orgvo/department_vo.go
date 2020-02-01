package orgvo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type DepartmentsReqVo struct {
	Page          *int
	Size          *int
	Params        *vo.DepartmentListReq
	CurrentUserId int64 `json:"userId"`
	OrgId         int64 `json:"orgId"`
}

type DepartmentsRespVo struct {
	DepartmentList *vo.DepartmentList `json:"data"`
	vo.Err
}

type DepartmentMembersReqVo struct {
	CurrentUserId int64                      `json:"userId"`
	OrgId         int64                      `json:"orgId"`
	Params        vo.DepartmentMemberListReq `json:"params"`
}

type DepartmentMembersRespVo struct {
	DepartmentMemberInfos []*vo.DepartmentMemberInfo `json:"data"`
	vo.Err
}
