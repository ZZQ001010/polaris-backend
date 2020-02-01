package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) RemoveProjectMember(req projectvo.RemoveProjectMemberReqVo) vo.CommonRespVo {
	res, err := service.RemoveProjectMember(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) ProjectUserList(req projectvo.ProjectUserListReq) projectvo.ProjectUserListRespVo {
	res, err := service.ProjectUserList(req.OrgId, req.Page, req.Size, req.Input)
	return projectvo.ProjectUserListRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) AddProjectMember(req projectvo.RemoveProjectMemberReqVo) vo.CommonRespVo {
	res, err := service.AddProjectMember(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}
