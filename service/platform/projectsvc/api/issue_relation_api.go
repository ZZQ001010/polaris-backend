package api

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
	"strings"
)

func (PostGreeter) CreateIssueComment(reqVo projectvo.CreateIssueCommentReqVo) vo.CommonRespVo {
	reqVo.Input.Comment = strings.TrimSpace(reqVo.Input.Comment)
	//checkCommentLenErr := util.CheckIssueCommentLen(reqVo.Input.Comment)
	//if checkCommentLenErr != nil{
	//	log.Error(checkCommentLenErr)
	//	return vo.CommonRespVo{Err: vo.NewErr(checkCommentLenErr), Void: nil}
	//}
	isCommentRight := format.VerifyIssueCommenFormat(reqVo.Input.Comment)
	if !isCommentRight {
		log.Error(errs.IssueCommentLenError)
		return vo.CommonRespVo{Err: vo.NewErr(errs.IssueCommentLenError), Void: nil}
	}

	res, err := service.CreateIssueComment(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) CreateIssueResource(reqVo projectvo.CreateIssueResourceReqVo) vo.CommonRespVo {
	res, err := service.CreateIssueResource(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) CreateIssueRelationIssue(reqVo projectvo.CreateIssueRelationIssueReqVo) vo.CommonRespVo {
	res, err := service.CreateIssueRelationIssue(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) DeleteIssueResource(reqVo projectvo.DeleteIssueResourceReqVo) vo.CommonRespVo {
	res, err := service.DeleteIssueResource(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) RelatedIssueList(reqVo projectvo.RelatedIssueListReqVo) projectvo.RelatedIssueListRespVo {
	res, err := service.RelatedIssueList(reqVo.OrgId, reqVo.Input)
	return projectvo.RelatedIssueListRespVo{Err: vo.NewErr(err), RelatedIssueList: res}
}

func (PostGreeter) IssueResources(reqVo projectvo.IssueResourcesReqVo) projectvo.IssueResourcesRespVo {
	res, err := service.IssueResources(reqVo.OrgId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.IssueResourcesRespVo{Err: vo.NewErr(err), IssueResources: res}
}

func (PostGreeter) CreateIssueRelationTags(reqVo projectvo.CreateIssueRelationTagsReqVo) vo.CommonRespVo {
	err := service.CreateIssueRelationTags(reqVo)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: &vo.Void{
		ID: reqVo.Input.ID,
	}}
}

func (GetGreeter) GetIssueMembers(reqVo projectvo.GetIssueMembersReqVo) projectvo.GetIssueMembersRespVo{
	res, err := service.GetIssueMembers(reqVo.OrgId, reqVo.IssueId)
	return projectvo.GetIssueMembersRespVo{Err: vo.NewErr(err), Data: *res}
}