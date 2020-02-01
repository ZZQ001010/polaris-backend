package bo

import "time"

type FeiShuAppTicketCacheBo struct {
	AppId string `json:"appId"`
	AppTicket string `json:"appTicket"`
	LastUpdateTime string `json:"lastUpdateTime"`
}

type FeiShuAppAccessTokenCacheBo struct {
	AppAccessToken string `json:"appAccessToken"`
	LastUpdateTime string `json:"lastUpdateTime"`
}

type FeiShuTenantAccessTokenCacheBo struct {
	TenantAccessToken string `json:"tenantAccessToken"`
	TenantKey string `json:"tenantKey"`
	LastUpdateTime string `json:"lastUpdateTime"`
}

type FeiShuCallBackMqBo struct {
	//类型
	Type string `json:"type"`
	//数据
	Data string `json:"data"`
	//消息时间
	MsgTime time.Time `json:"msgTime"`
}