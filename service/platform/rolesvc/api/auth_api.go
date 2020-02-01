package api

import (
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/service"
	"upper.io/db.v3/lib/sqlbuilder"
)

func (PostGreeter) Authenticate(req rolevo.AuthenticateReqVo) vo.VoidErr {
	err := service.Authenticate(req.OrgId, req.UserId, req.AuthInfoReqVo.ProjectAuthInfo, req.AuthInfoReqVo.IssueAuthInfo, req.Path, req.Operation)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) RoleUserRelation(req rolevo.RoleUserRelationReqVo) vo.VoidErr{
	err := service.RoleUserRelation(req.OrgId, req.UserId, req.RoleId)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) RemoveRoleUserRelation(req rolevo.RemoveRoleUserRelationReqVo) vo.VoidErr{
	err := service.RemoveRoleUserRelation(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) RoleInit(req rolevo.RoleInitReqVo) rolevo.RoleInitRespVo {
	respVo := rolevo.RoleInitRespVo{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error{
		roleInitResp, err := service.RoleInit(req.OrgId, tx)
		respVo.RoleInitResp = roleInitResp
		respVo.Err = vo.NewErr(err)
		return err
	})
	return respVo
}

