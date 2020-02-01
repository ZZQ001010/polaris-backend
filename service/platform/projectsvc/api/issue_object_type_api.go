package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) IssueObjectTypeList(reqVo projectvo.IssueObjectTypeListReqVo) projectvo.IssueObjectTypeListRespVo {
	res, err := service.IssueObjectTypeList(reqVo.OrgId, reqVo.Page, reqVo.Size, reqVo.Params)
	return projectvo.IssueObjectTypeListRespVo{Err: vo.NewErr(err), IssueObjectTypeList: res}
}

func (PostGreeter) CreateIssueObjectType(reqVo projectvo.CreateIssueObjectTypeReqVo) vo.CommonRespVo {
	res, err := service.CreateIssueObjectType(reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateIssueObjectType(reqVo projectvo.UpdateIssueObjectTypeReqVo) vo.CommonRespVo {
	res, err := service.UpdateIssueObjectType(reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) DeleteIssueObjectType(reqVo projectvo.DeleteIssueObjectTypeReqVo) vo.CommonRespVo {
	res, err := service.DeleteIssueObjectType(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}
