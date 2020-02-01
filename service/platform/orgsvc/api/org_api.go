package api

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/service"
)

func (GetGreeter) GetOrgBoList() orgvo.GetOrgBoListRespVo {
	res, err := service.GetOrgBoList()
	return orgvo.GetOrgBoListRespVo{Err: vo.NewErr(err), OrganizationBoList: res}
}

func (PostGreeter) CreateOrg(req orgvo.CreateOrgReqVo) orgvo.CreateOrgRespVo {
	sourceChannel := ""
	sourcePlatform := ""
	if req.Data.CreateOrgReq.SourceChannel != nil {
		sourceChannel = *req.Data.CreateOrgReq.SourceChannel
	}
	if req.Data.CreateOrgReq.SourcePlatform != nil {
		sourcePlatform = *req.Data.CreateOrgReq.SourcePlatform
	}
	orgId, err := service.CreateOrg(req, sourceChannel, sourcePlatform)
	if err != nil {
		log.Error(err)
		return orgvo.CreateOrgRespVo{Err: vo.NewErr(err), Data: orgvo.CreateOrgRespVoData{OrgId: orgId}}
	}

	if req.Data.ImportSampleData == 1 {
		//初始化示例数据
		if sourceChannel != consts.AppSourceChannelFeiShu {
			err1 := service.LarkInit(orgId, req.UserId, sourceChannel, sourcePlatform)
			if err1 != nil {
				log.Error(err1)
			}
		} else {
			resp := projectfacade.DataInitForLarkApplet(vo.BasicInfoReqVo{
				OrgId:  orgId,
				UserId: req.UserId,
			})
			if resp.Failure() {
				log.Error(resp.Error())
			}
		}
	}
	return orgvo.CreateOrgRespVo{Err: vo.NewErr(err), Data: orgvo.CreateOrgRespVoData{OrgId: orgId}}
}

func (PostGreeter) UserOrganizationList(req orgvo.UserOrganizationListReqVo) orgvo.UserOrganizationListRespVo {

	res, err := service.UserOrganizationList(req.UserId)

	return orgvo.UserOrganizationListRespVo{Err: vo.NewErr(err), UserOrganizationListResp: res}
}

func (PostGreeter) SwitchUserOrganization(req orgvo.SwitchUserOrganizationReqVo) orgvo.SwitchUserOrganizationRespVo {

	err := service.SwitchUserOrganization(req.OrgId, req.UserId, req.Token)

	return orgvo.SwitchUserOrganizationRespVo{Err: vo.NewErr(err), OrgId: req.OrgId}
}

func (PostGreeter) UpdateOrganizationSetting(req orgvo.UpdateOrganizationSettingReqVo) orgvo.UpdateOrganizationSettingRespVo {

	res, err := service.UpdateOrganizationSetting(req)

	return orgvo.UpdateOrganizationSettingRespVo{Err: vo.NewErr(err), OrgId: res}
}

func (PostGreeter) OrganizationInfo(req orgvo.OrganizationInfoReqVo) orgvo.OrganizationInfoRespVo {

	res, err := service.OrganizationInfo(req)

	return orgvo.OrganizationInfoRespVo{Err: vo.NewErr(err), OrganizationInfo: res}
}

func (PostGreeter) ScheduleOrganizationPageList(req orgvo.ScheduleOrganizationPageListReqVo) orgvo.ScheduleOrganizationPageListRespVo {

	res, err := service.ScheduleOrganizationPageList(req)

	return orgvo.ScheduleOrganizationPageListRespVo{Err: vo.NewErr(err), ScheduleOrganizationPageListResp: res}
}

//通过来源获取组织id列表
func (PostGreeter) GetOrgIdListBySourceChannel(req orgvo.GetOrgIdListBySourceChannelReqVo) orgvo.GetOrgIdListBySourceChannelRespVo{
	res, err := service.GetOrgIdListBySourceChannel(req.SourceChannel, req.Page, req.Size)
	return orgvo.GetOrgIdListBySourceChannelRespVo{Err: vo.NewErr(err), Data: orgvo.GetOrgIdListBySourceChannelRespData{OrgIds: res}}
}