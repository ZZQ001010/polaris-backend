package projectvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type ProjectsRepVo struct {
	Page             int              `json:"page"`
	Size             int              `json:"size"`
	ProjectExtraBody ProjectExtraBody `json:"projectExtraBody"`
	OrgId            int64            `json:"orgId"`
	UserId           int64            `json:"userId"`
	SourceChannel    string           `json:"sourceChannel"`
}

type ProjectExtraBody struct {
	Params map[string]interface{} `json:"params"`
	Order  []*string              `json:"order"`
	Input  *vo.ProjectsReq        `json:"input"`
}

//项目列表
type ProjectsRespVo struct {
	vo.Err
	ProjectList *vo.ProjectList `json:"data"`
}

//项目统计入参
type ProjectDayStatsReqVo struct {
	Page   uint                  `json:"page"`
	Size   uint                  `json:"size"`
	Params *vo.ProjectDayStatReq `json:"params"`
	OrgId  int64                 `json:"orgId"`
}

//项目统计出参
type ProjectDayStatsRespVo struct {
	vo.Err
	ProjectDayStatList *vo.ProjectDayStatList `json:"data"`
}

type CreateProjectReqVo struct {
	Input         vo.CreateProjectReq `json:"input"`
	OrgId         int64               `json:"orgId"`
	UserId        int64               `json:"userId"`
	SourceChannel string              `json:"sourceChannel"`
}

type UpdateProjectReqVo struct {
	Input         vo.UpdateProjectReq `json:"input"`
	OrgId         int64               `json:"orgId"`
	UserId        int64               `json:"userId"`
	SourceChannel string              `json:"sourceChannel"`
}

type ProjectRespVo struct {
	vo.Err
	Project *vo.Project `json:"data"`
}

type ProjectIdReqVo struct {
	ProjectId     int64  `json:"projectId"`
	OrgId         int64  `json:"orgId"`
	UserId        int64  `json:"userId"`
	SourceChannel string `json:"sourceChannel"`
}

type GetCacheProjectInfoReqVo struct {
	ProjectId int64 `json:"projectId"`
	OrgId     int64 `json:"orgId"`
}

type GetCacheProjectInfoRespVo struct {
	vo.Err

	ProjectCacheBo *bo.ProjectAuthBo `json:"data"`
}

type QuitProjectRespVo struct {
	vo.Err
	QuitProject *vo.QuitResult `json:"data"`
}

type OperateProjectRespVo struct {
	vo.Err
	OperateProject *vo.OperateProjectResp `json:"data"`
}

type ProjectStatisticsRespVo struct {
	vo.Err
	ProjectStatistics *vo.ProjectStatisticsResp `json:"data"`
}

type ProjectInfoRespVo struct {
	vo.Err
	ProjectInfo *vo.ProjectInfo `json:"data"`
}

type UpdateProjectStatusReqVo struct {
	Input         vo.UpdateProjectStatusReq `json:"input"`
	OrgId         int64                     `json:"orgId"`
	UserId        int64                     `json:"userId"`
	SourceChannel string                    `json:"sourceChannel"`
}

type ProjectInfoReqVo struct {
	Input vo.ProjectInfoReq `json:"input"`
	OrgId int64             `json:"orgId"`
	SourceChannel string    `json:"sourceChannel"`
}

type GetProjectProcessIdReqVo struct {
	OrgId               int64 `json:"orgId"`
	ProjectId           int64 `json:"projectId"`
	ProjectObjectTypeId int64 `json:"projectObjectTypeId"`
}

type GetProjectProcessIdRespVo struct {
	ProcessId int64 `json:"data"`

	vo.Err
}

type GetProjectBoListByProjectTypeLangCodeReqVo struct {
	OrgId               int64   `json:"orgId"`
	ProjectTypeLangCode *string `json:"projectTypeLangCode"`
}

type GetProjectBoListByProjectTypeLangCodeRespVo struct {
	ProjectBoList []bo.ProjectBo `json:"data"`
	vo.Err
}

type AppendProjectDayStatReqVo struct {
	ProjectBo bo.ProjectBo `json:"projectBo"`
	Date      string       `json:"date"`
}

type ProjectObjectTypesReqVo struct {
	Page   uint                      `json:"page"`
	Size   uint                      `json:"size"`
	Params *vo.ProjectObjectTypesReq `json:"params"`
	OrgId  int64                     `json:"orgId"`
}

type ProjectObjectTypeWithProjectVo struct {
	ProjectId int64 `json:"projectId"`
	OrgId     int64 `json:"orgId"`
}

type GetSimpleProjectInfoReqVo struct {
	OrgId int64   `json:"orgId"`
	Ids   []int64 `json:"ids"`
}

type GetSimpleProjectInfoRespVo struct {
	vo.Err
	Data *[]vo.Project `json:"data"`
}

type ProjectRelationList struct {
	Id           int64 `json:"id"`
	RelationType int   `json:"relationType"`
	RelationId   int64 `json:"relationId"`
}

type GetProjectRelationReqVo struct {
	ProjectId    int64   `json:"projectId"`
	RelationType []int64 `json:"relationType"`
}

type GetProjectRelationRespVo struct {
	vo.Err
	Data []ProjectRelationList `json:"data"`
}

type GetProjectInfoListByOrgIdsReqVo struct {
	OrgIds []int64 `json:"orgIds"`
}

type GetProjectInfoListByOrgIdsListRespVo struct {
	vo.Err
	ProjectInfoListByOrgIdsRespVo []GetProjectInfoListByOrgIdsRespVo `json:"data"`
}

type GetProjectInfoListByOrgIdsRespVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
	Owner     int64 `json:"Owner"`
}

type OrgProjectListReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type OrgProjectMemberReqVo struct {
	OrgId     int64 `json:"orgId"`
	UserId    int64 `json:"userId"`
	ProjectId int64 `json:"projectId"`
}

type OrgProjectMemberListRespVo struct {
	vo.Err
	OrgProjectMemberRespVo *OrgProjectMemberRespVo `json:"data"`
}

type OrgProjectMemberRespVo struct {
	Owner        OrgProjectMemberVo   `json:"owner"`
	Participants []OrgProjectMemberVo `json:"participants"`
	Follower     []OrgProjectMemberVo `json:"follower"`
	AllMembers   []OrgProjectMemberVo `json:"allMembers"`
}

type OrgProjectMemberVo struct {
	UserId        int64  `json:"userId"`
	OutUserId     string `json:"outUserId"` //有可能为空
	OrgId         int64  `json:"orgId"`
	OutOrgId      string `json:"outOrgId"` //有可能为空
	Name          string `json:"name"`
	NamePy		  string `json:"namePy"`
	Avatar        string `json:"avatar"`
	HasOutInfo    bool   `json:"hasOutInfo"`
	HasOrgOutInfo bool   `json:"hasOrgOutInfo"`

	OrgUserIsDelete    int `json:"orgUserIsDelete"`          //是否被组织移除
	OrgUserStatus      int `json:"orgUserStatus"`      //用户组织状态
	OrgUserCheckStatus int `json:"orgUserCheckStatus"` //用户组织审核状态
}

type RemoveProjectMemberReqVo struct {
	Input  vo.RemoveProjectMemberReq `json:"input"`
	OrgId  int64                     `json:"orgId"`
	UserId int64                     `json:"userId"`
}

type ProjectUserListReq struct {
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
	OrgId int64                 `json:"orgId"`
	Input vo.ProjectUserListReq `json:"input"`
}

type ProjectUserListRespVo struct {
	vo.Err
	Data *vo.ProjectUserListResp `json:"data"`
}
