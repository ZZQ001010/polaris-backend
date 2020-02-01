package trendsvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type UnreadNoticeCountReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type UnreadNoticeCountRespVo struct {
	vo.Err
	Count int64 `json:"count"`
}

type NoticeListReqVo struct {
	UserId int64             `json:"userId"`
	OrgId  int64             `json:"orgId"`
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Input  *vo.NoticeListReq `json:"input"`
}

type NoticeListRespVo struct {
	vo.Err
	Data *vo.NoticeList `json:"data"`
}

type GetMQTTChannelKeyReqVo struct {
	UserId int64             `json:"userId"`
	OrgId  int64             `json:"orgId"`
	Input vo.GetMQTTChannelKeyReq `json:"input"`
}

type GetMQTTChannelKeyRespVo struct {
	vo.Err
	Data *vo.GetMQTTChannelKeyResp `json:"data"`
}

type PushMQTTDataRefreshMsgReqVo struct {
	RefreshBo bo.MQTTDataRefreshNotice `json:"refreshBo"`
}

