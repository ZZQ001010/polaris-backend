package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) Projects(reqVo projectvo.ProjectsRepVo) projectvo.ProjectsRespVo {
	res, err := service.Projects(reqVo)
	return projectvo.ProjectsRespVo{Err: vo.NewErr(err), ProjectList: res}
}

func (PostGreeter) CreateProject(reqVo projectvo.CreateProjectReqVo) projectvo.ProjectRespVo {
	res, err := service.CreateProject(reqVo)
	return projectvo.ProjectRespVo{Err: vo.NewErr(err), Project: res}
}

func (PostGreeter) UpdateProject(reqVo projectvo.UpdateProjectReqVo) projectvo.ProjectRespVo {
	res, err := service.UpdateProject(reqVo)
	return projectvo.ProjectRespVo{Err: vo.NewErr(err), Project: res}
}

func (PostGreeter) QuitProject(reqVo projectvo.ProjectIdReqVo) projectvo.QuitProjectRespVo {
	res, err := service.QuitProject(reqVo)
	return projectvo.QuitProjectRespVo{Err: vo.NewErr(err), QuitProject: res}
}

func (PostGreeter) StarProject(reqVo projectvo.ProjectIdReqVo) projectvo.OperateProjectRespVo {
	res, err := service.StarProject(reqVo)
	return projectvo.OperateProjectRespVo{Err: vo.NewErr(err), OperateProject: res}
}

func (PostGreeter) UnstarProject(reqVo projectvo.ProjectIdReqVo) projectvo.OperateProjectRespVo {
	res, err := service.UnstarProject(reqVo.OrgId, reqVo.UserId, reqVo.ProjectId)
	return projectvo.OperateProjectRespVo{Err: vo.NewErr(err), OperateProject: res}
}

func (GetGreeter) ProjectStatistics(reqVo projectvo.ProjectIdReqVo) projectvo.ProjectStatisticsRespVo {
	res, err := service.ProjectStatistics(reqVo.OrgId, reqVo.ProjectId)
	return projectvo.ProjectStatisticsRespVo{Err: vo.NewErr(err), ProjectStatistics: res}
}

func (PostGreeter) UpdateProjectStatus(reqVo projectvo.UpdateProjectStatusReqVo) vo.CommonRespVo {
	res, err := service.UpdateProjectStatus(reqVo)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) ProjectInfo(reqVo projectvo.ProjectInfoReqVo) projectvo.ProjectInfoRespVo {
	res, err := service.ProjectInfo(reqVo.OrgId, reqVo.Input, reqVo.SourceChannel)
	return projectvo.ProjectInfoRespVo{Err: vo.NewErr(err), ProjectInfo: res}
}

func (GetGreeter) GetProjectProcessId(req projectvo.GetProjectProcessIdReqVo) projectvo.GetProjectProcessIdRespVo {
	res, err := service.GetProjectProcessId(req.OrgId, req.ProjectId, req.ProjectObjectTypeId)
	return projectvo.GetProjectProcessIdRespVo{ProcessId: res, Err: vo.NewErr(err)}
}

//通过项目类型langCode获取项目列表
func (GetGreeter) GetProjectBoListByProjectTypeLangCode(req projectvo.GetProjectBoListByProjectTypeLangCodeReqVo) projectvo.GetProjectBoListByProjectTypeLangCodeRespVo {
	res, err := service.GetProjectBoListByProjectTypeLangCode(req.OrgId, req.ProjectTypeLangCode)
	return projectvo.GetProjectBoListByProjectTypeLangCodeRespVo{ProjectBoList: res, Err: vo.NewErr(err)}
}

func (PostGreeter) GetSimpleProjectInfo(req projectvo.GetSimpleProjectInfoReqVo) projectvo.GetSimpleProjectInfoRespVo {
	res, err := service.GetSimpleProjectInfo(req.OrgId, req.Ids)
	return projectvo.GetSimpleProjectInfoRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) GetProjectRelation(req projectvo.GetProjectRelationReqVo) projectvo.GetProjectRelationRespVo {
	res, err := service.GetProjectRelation(req.ProjectId, req.RelationType)
	return projectvo.GetProjectRelationRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) ArchiveProject(reqVo projectvo.ProjectIdReqVo) vo.CommonRespVo {
	res, err := service.ArchiveProject(reqVo.OrgId, reqVo.UserId, reqVo.ProjectId)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) CancelArchivedProject(reqVo projectvo.ProjectIdReqVo) vo.CommonRespVo {
	res, err := service.CancelArchivedProject(reqVo.OrgId, reqVo.UserId, reqVo.ProjectId)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (GetGreeter) GetCacheProjectInfo(reqVo projectvo.GetCacheProjectInfoReqVo) projectvo.GetCacheProjectInfoRespVo {
	res, err := service.GetCacheProjectInfo(reqVo)
	return projectvo.GetCacheProjectInfoRespVo{Err: vo.NewErr(err), ProjectCacheBo: res}
}

//通过组织id集合获取未删除 未归档的项目
func (PostGreeter) GetProjectInfoByOrgIds(req projectvo.GetProjectInfoListByOrgIdsReqVo) projectvo.GetProjectInfoListByOrgIdsListRespVo {
	res, err := service.GetProjectInfoByOrgIds(req.OrgIds)
	return projectvo.GetProjectInfoListByOrgIdsListRespVo{ProjectInfoListByOrgIdsRespVo: res, Err: vo.NewErr(err)}
}

func (GetGreeter) OrgProjectMember(reqVo projectvo.OrgProjectMemberReqVo) projectvo.OrgProjectMemberListRespVo {
	res, err := service.OrgProjectMembers(reqVo)
	return projectvo.OrgProjectMemberListRespVo{Err: vo.NewErr(err), OrgProjectMemberRespVo: res}
}
