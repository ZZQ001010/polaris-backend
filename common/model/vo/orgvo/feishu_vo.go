package orgvo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type FeiShuAuthRespVo struct {
	vo.Err
	Auth *vo.FeiShuAuthResp `json:"data"`
}