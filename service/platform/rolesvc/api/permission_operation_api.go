package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/service"
)

func (GetGreeter) PermissionOperationList(req rolevo.PermissionOperationListReqVo) rolevo.PermissionOperationListRespVo {
	res, err := service.PermissionOperationList(req.OrgId, req.RoleId, req.UserId, req.ProjectId)
	return rolevo.PermissionOperationListRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) UpdateRolePermissionOperation(req rolevo.UpdateRolePermissionOperationReqVo) vo.CommonRespVo {
	res, err := service.UpdateRolePermissionOperation(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (GetGreeter) GetPersonalPermissionInfo(req rolevo.GetPersonalPermissionInfoReqVo) rolevo.GetPersonalPermissionInfoRespVo {
	res, err := service.GetPersonalPermissionInfo(req.OrgId, req.UserId, req.ProjectId, req.IssueId, req.SourceChannel)
	return rolevo.GetPersonalPermissionInfoRespVo{Err:vo.NewErr(err), Data:res}
}