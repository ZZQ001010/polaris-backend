package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) Departments(req orgvo.DepartmentsReqVo) orgvo.DepartmentsRespVo {
	page := req.Page
	size := req.Size
	params := req.Params
	orgId := req.OrgId

	pageA := uint(0)
	sizeA := uint(0)
	if page != nil && size != nil && *page > 0 && *size > 0 {
		pageA = uint(*page)
		sizeA = uint(*size)
	}
	res, err := service.Departments(pageA, sizeA, params, orgId)
	return orgvo.DepartmentsRespVo{Err: vo.NewErr(err), DepartmentList: res}
}

func (PostGreeter) DepartmentMembers(req orgvo.DepartmentMembersReqVo) orgvo.DepartmentMembersRespVo {
	res, err := service.DepartmentMembers(req.Params, req.OrgId)
	return orgvo.DepartmentMembersRespVo{Err: vo.NewErr(err), DepartmentMemberInfos: res}
}
