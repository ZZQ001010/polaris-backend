package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/trendsvo"
	"github.com/galaxy-book/polaris-backend/service/platform/trendssvc/service"
)

func (PostGreeter) UnreadNoticeCount(req trendsvo.UnreadNoticeCountReqVo) trendsvo.UnreadNoticeCountRespVo {
	res, err := service.UnreadNoticeCount(req.OrgId, req.UserId)
	return trendsvo.UnreadNoticeCountRespVo{Err: vo.NewErr(err), Count: int64(res)}
}

func (PostGreeter) NoticeList(req trendsvo.NoticeListReqVo) trendsvo.NoticeListRespVo {
	res, err := service.NoticeList(req.OrgId, req.UserId, req.Page, req.Size, req.Input)
	return trendsvo.NoticeListRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) GetMQTTChannelKey(req trendsvo.GetMQTTChannelKeyReqVo) trendsvo.GetMQTTChannelKeyRespVo {
	res, err := service.GetMQTTChannelKey(req.OrgId, req.UserId, req.Input)
	return trendsvo.GetMQTTChannelKeyRespVo{Err: vo.NewErr(err), Data: res}
}

