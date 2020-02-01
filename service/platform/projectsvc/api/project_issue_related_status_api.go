package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) ProjectIssueRelatedStatus(reqVo projectvo.ProjectIssueRelatedStatusReqVo) projectvo.ProjectIssueRelatedStatusRespVo {
	res, err := service.ProjectIssueRelatedStatus(reqVo.OrgId, reqVo.Input)
	return projectvo.ProjectIssueRelatedStatusRespVo{Err: vo.NewErr(err), ProjectIssueRelatedStatus: res}
}
