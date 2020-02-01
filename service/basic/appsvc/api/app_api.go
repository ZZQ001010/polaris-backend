package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/appvo"
	services "github.com/galaxy-book/polaris-backend/service/basic/appsvc/service"
)

func (GetGreeter) GetAppInfoByActive(req appvo.AppInfoReqVo) appvo.AppInfoRespVo {
	res, err := services.GetAppInfoByActive(req.AppCode)

	return appvo.AppInfoRespVo{
		Err:     vo.NewErr(err),
		AppInfo: res,
	}
}
