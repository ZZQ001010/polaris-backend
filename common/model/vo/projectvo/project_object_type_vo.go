package projectvo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type ProjectObjectTypeListRespVo struct {
	vo.Err
	ProjectObjectTypeList *vo.ProjectObjectTypeList `json:"data"`
}

type ProjectSupportObjectTypesRespVo struct {
	vo.Err
	ProjectSupportObjectTypes *vo.ProjectSupportObjectTypeListResp `json:"data"`
}

type CreateProjectObjectTypeReqVo struct {
	Input  vo.CreateProjectObjectTypeReq `json:"input"`
	OrgId  int64                         `json:"orgId"`
	UserId int64                         `json:"userId"`
}

type UpdateProjectObjectTypeReqVo struct {
	Input  vo.UpdateProjectObjectTypeReq `json:"input"`
	OrgId  int64                         `json:"orgId"`
	UserId int64                         `json:"userId"`
}

type DeleteProjectObjectTypeReqVo struct {
	Input  vo.DeleteProjectObjectTypeReq `json:"input"`
	OrgId  int64                         `json:"orgId"`
	UserId int64                         `json:"userId"`
}

type ProjectSupportObjectTypesReqVo struct {
	Input vo.ProjectSupportObjectTypeListReq `json:"input"`
	OrgId int64                              `json:"orgId"`
}

type ProjectObjectTypeWithProjectListRespVo struct {
	vo.Err
	ProjectObjectTypeWithProjectList *vo.ProjectObjectTypeWithProjectList `json:"data"`
}
