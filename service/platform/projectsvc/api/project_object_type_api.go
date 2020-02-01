package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) ProjectObjectTypeList(reqVo projectvo.ProjectObjectTypesReqVo) projectvo.ProjectObjectTypeListRespVo {
	res, err := service.ProjectObjectTypeList(reqVo.OrgId, reqVo.Page, reqVo.Size, reqVo.Params)
	return projectvo.ProjectObjectTypeListRespVo{Err: vo.NewErr(err), ProjectObjectTypeList: res}
}

func (PostGreeter) CreateProjectObjectType(reqVo projectvo.CreateProjectObjectTypeReqVo) vo.CommonRespVo {
	res, err := service.CreateProjectObjectType(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateProjectObjectType(reqVo projectvo.UpdateProjectObjectTypeReqVo) vo.CommonRespVo {
	res, err := service.UpdateProjectObjectType(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) DeleteProjectObjectType(reqVo projectvo.DeleteProjectObjectTypeReqVo) vo.CommonRespVo {
	res, err := service.DeleteProjectObjectType(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) ProjectSupportObjectTypes(reqVo projectvo.ProjectSupportObjectTypesReqVo) projectvo.ProjectSupportObjectTypesRespVo {
	res, err := service.ProjectSupportObjectTypes(reqVo.OrgId, reqVo.Input)
	return projectvo.ProjectSupportObjectTypesRespVo{Err: vo.NewErr(err), ProjectSupportObjectTypes: res}
}

func (PostGreeter) ProjectObjectTypesWithProject(reqVo projectvo.ProjectObjectTypeWithProjectVo) projectvo.ProjectObjectTypeWithProjectListRespVo {
	res, err := service.ProjectObjectTypesWithProject(reqVo.OrgId, reqVo.ProjectId)
	return projectvo.ProjectObjectTypeWithProjectListRespVo{Err: vo.NewErr(err), ProjectObjectTypeWithProjectList: res}
}
