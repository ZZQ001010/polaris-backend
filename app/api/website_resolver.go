package api

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/websitevo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/websitefacade"
)

func (r *mutationResolver) RegisterWebSiteContact(ctx context.Context, input vo.RegisterWebSiteContactReq) (*vo.Void, error) {
	registerWebSiteContactReqVo := websitevo.RegisterWebSiteContactReqVo{
		Input: input,
	}

	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err == nil && cacheUserInfo != nil{
		registerWebSiteContactReqVo.OrgId = cacheUserInfo.OrgId
		registerWebSiteContactReqVo.UserId = cacheUserInfo.UserId
	}

	respVo := websitefacade.RegisterWebSiteContact(registerWebSiteContactReqVo)
	return respVo.Void, respVo.Error()

}
