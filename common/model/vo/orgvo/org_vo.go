package orgvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type UserOrganizationListReqVo struct {
	UserId int64 `json:"userId"`
}

type UserOrganizationListRespVo struct {
	vo.Err
	UserOrganizationListResp *vo.UserOrganizationListResp `json:"data"`
}

type SwitchUserOrganizationReqVo struct {
	UserId int64  `json:"userId"`
	OrgId  int64  `json:"orgId"`
	Token  string `json:"userToken"`
}

type SwitchUserOrganizationRespVo struct {
	vo.Err
	OrgId int64 `json:"data"`
}

type UpdateOrganizationSettingReqVo struct {
	//入参
	Input vo.UpdateOrganizationSettingsReq `json:"input"`
	//用户Id
	UserId int64 `json:"userId"`
}

type UpdateOrganizationSettingRespVo struct {
	vo.Err
	OrgId int64 `json:"data"`
}

type OrganizationInfoReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type OrganizationInfoRespVo struct {
	vo.Err
	OrganizationInfo *vo.OrganizationInfoResp `json:"data"`
}

type ScheduleOrganizationPageListReqVo struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ScheduleOrganizationPageListRespVo struct {
	vo.Err
	ScheduleOrganizationPageListResp *ScheduleOrganizationPageListResp `json:"data"`
}

type ScheduleOrganizationPageListResp struct {
	Total                        int64                            `json:"total"`
	ScheduleOrganizationListResp *[]*ScheduleOrganizationListResp `json:"list"`
}

type ScheduleOrganizationListResp struct {
	OrgId                      int64  `json:"orgId"`
	ProjectDailyReportSendTime string `json:"projectDailyReportSendTime"`
}

type GetOrgIdListBySourceChannelReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	Page int `json:"page"`
	Size int `json:"size"`
}

type GetOrgIdListBySourceChannelRespVo struct {
	Data GetOrgIdListBySourceChannelRespData `json:"data"`
	vo.Err
}

type GetOrgIdListBySourceChannelRespData struct {
	OrgIds []int64 `json:"orgIds"`
}