package orgvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type CacheUserInfoVo struct {
	vo.Err

	CacheInfo bo.CacheUserInfoBo `json:"data"`
}
