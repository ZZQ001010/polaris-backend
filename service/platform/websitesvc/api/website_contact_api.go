package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/websitevo"
	"github.com/galaxy-book/polaris-backend/service/platform/websitesvc/service"
)

func (PostGreeter) RegisterWebSiteContact(reqVo websitevo.RegisterWebSiteContactReqVo) vo.CommonRespVo {
	res, err := service.RegisterWebSiteContact(reqVo)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}