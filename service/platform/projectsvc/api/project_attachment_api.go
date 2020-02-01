package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) DeleteProjectAttachment(reqVo projectvo.DeleteProjectAttachmentReqVo) projectvo.DeleteProjectAttachmentRespVo {
	res, err := service.DeleteProjectAttachment(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return projectvo.DeleteProjectAttachmentRespVo{Err: vo.NewErr(err), Output: res}
}

func (PostGreeter) GetProjectAttachment(reqVo projectvo.GetProjectAttachmentReqVo) projectvo.GetProjectAttachmentRespVo {
	res, err := service.GetProjectAttachment(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.GetProjectAttachmentRespVo{Err: vo.NewErr(err), Output: res}
}
