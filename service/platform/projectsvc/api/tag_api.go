package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) CreateTag(reqVo projectvo.CreateTagReqVo) vo.CommonRespVo {
	res, err := service.CreateTag(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Void: res, Err: vo.NewErr(err)}
}

func (PostGreeter) TagList(reqVo projectvo.TagListReqVo) projectvo.TagListRespVo {
	res, err := service.TagList(reqVo.OrgId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.TagListRespVo{Err: vo.NewErr(err), Data: res}
}

func (GetGreeter) TagDefaultStyle() projectvo.TagDefaultStyleRespVo {
	res := service.GetTagDefaultStyle()
	return projectvo.TagDefaultStyleRespVo{Data: res}
}

func (PostGreeter) HotTagList(reqVo projectvo.HotTagListReqVo) projectvo.TagListRespVo {
	res, err := service.HotTagList(reqVo.OrgId, reqVo.ProjectId)
	return projectvo.TagListRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) DeleteTag(reqVo projectvo.DeleteTagReqVo) vo.CommonRespVo {
	res, err := service.DeleteTag(reqVo.OrgId, reqVo.UserId, reqVo.Data)
	return vo.CommonRespVo{Void: res, Err: vo.NewErr(err)}
}

func (PostGreeter) UpdateTag(reqVo projectvo.UpdateTagReqVo) vo.CommonRespVo {
	res, err := service.UpdateTag(reqVo.OrgId, reqVo.UserId, reqVo.Data)
	return vo.CommonRespVo{Void: res, Err: vo.NewErr(err)}
}

