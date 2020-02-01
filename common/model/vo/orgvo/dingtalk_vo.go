package orgvo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type GetJsAPISignRespVo struct {
	vo.Err
	GetJsAPISign *vo.JsAPISignResp `json:"data"`
}

type AuthRespVo struct {
	vo.Err
	Auth *vo.AuthResp `json:"data"`
}
