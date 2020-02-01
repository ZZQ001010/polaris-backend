package orgvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type PersonalInfoRespVo struct {
	vo.Err
	PersonalInfo *vo.PersonalInfo `json:"data"`
}

type PersonalInfoReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	UserId        int64  `json:"userId"`
	OrgId         int64  `json:"orgId"`
}

type GetUserIdsReqVo struct {
	SourceChannel string       `json:"sourceChannel"`
	EmpIdsBody    EmpIdsBodyVo `json:"empIdsBody"`
	CorpId        string       `json:"corpId"`
	OrgId         int64        `json:"orgId"`
}

type EmpIdsBodyVo struct {
	EmpIds []string `json:"empIds"`
}

type GetUserIdsRespVo struct {
	vo.Err
	GetUserIds []*vo.UserIDInfo `json:"data"`
}

type GetUserIdReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	EmpId         string `json:"empId"`
	CorpId        string `json:"corpId"`
	OrgId         int64  `json:"orgId"`
}

type GetUserIdRespVo struct {
	vo.Err
	GetUserId *vo.UserIDInfo `json:"data"`
}

type UserConfigInfoRespVo struct {
	vo.Err
	UserConfigInfo *vo.UserConfig `json:"data"`
}

type UserConfigInfoReqVo struct {
	UserId int64 `json:"userId"`
	OrgId  int64 `json:"orgId"`
}

type UpdateUserConfigRespVo struct {
	vo.Err
	UpdateUserConfig *vo.UpdateUserConfigResp `json:"data"`
}

type UpdateUserConfigReqVo struct {
	UpdateUserConfigReq vo.UpdateUserConfigReq `json:"updateUserConfigReq"`
	OrgId               int64                  `json:"orgId"`
	UserId              int64                  `json:"userId"`
}

type UpdateUserPcConfigReqVo struct {
	UpdateUserPcConfigReq vo.UpdateUserPcConfigReq `json:"updateUserConfigPcReq"`
	OrgId               int64                  `json:"orgId"`
	UserId              int64                  `json:"userId"`
}

type UpdateUserInfoReqVo struct {
	UpdateUserInfoReq vo.UpdateUserInfoReq `json:"updateUserInfoReq"`
	OrgId             int64                `json:"orgId"`
	UserId            int64                `json:"userId"`
}

type UpdateUserDefaultProjectIdConfigReqVo struct {
	UpdateUserDefaultProjectIdConfigReq vo.UpdateUserDefaultProjectConfigReq `json:"updateUserDefaultProjectIdConfigReq"`
	OrgId                               int64                                `json:"orgId"`
	UserId                              int64                                `json:"userId"`
}

type VerifyOrgReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type VerifyOrgUsersReqVo struct {
	OrgId  int64 `json:"orgId"`
	Input VerifyOrgUsersReqData `json:"input"`
}

type VerifyOrgUsersReqData struct {
	UserIds []int64 `json:"userIds"`
}

type GetUserInfoReqVo struct {
	OrgId         int64  `json:"orgId"`
	UserId        int64  `json:"userId"`
	SourceChannel string `json:"sourceChannel"`
}

type GetUserInfoRespVo struct {
	vo.Err
	UserInfo *bo.UserInfoBo `json:"data"`
}

type GetOrgBoListRespVo struct {
	vo.Err
	OrganizationBoList []bo.OrganizationBo `json:"data"`
}

type CreateOrgRespVo struct {
	Data CreateOrgRespVoData `json:"data"`
	vo.Err
}

type CreateOrgRespVoData struct {
	OrgId int64 `json:"orgId"`
}

type GetOutUserInfoListBySourceChannelReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	Page          int    `json:"page"`
	Size          int    `json:"size"`
}

type GetOutUserInfoBySourceChannelByUserInfoReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	OrgId         int64  `json:"orgId"`
	UserId        int64  `json:"userId"`
}

type GetOutUserInfoListBySourceChannelRespVo struct {
	UserOutInfoBoList []bo.UserOutInfoBo `json:"data"`
	vo.Err
}
type GetOutUserInfoBySourceChannelByUserInfoRespVo struct {
	UserOutInfoBoList bo.UserOutInfoBo `json:"data"`
	vo.Err
}

type GetUserInfoListReqVo struct {
	OrgId int64 `json:"orgId"`
}

type GetUserInfoListRespVo struct {
	vo.Err
	SimpleUserInfo []bo.SimpleUserInfoBo `json:"data"`
}

type CreateOrgReqVo struct {
	Data   CreateOrgReqVoData `json:"data"`
	OrgId  int64              `json:"orgId"`
	UserId int64              `json:"userId"`
}

type CreateOrgReqVoData struct {
	CreatorId        int64           `json:"creatorId"`
	CreateOrgReq     vo.CreateOrgReq `json:"createOrgReq"`
	UserToken        string          `json:"userToken"`
	ImportSampleData int             `json:"importSampleData"`
}

type GetUserInfoByUserIdsReqVo struct {
	UserIds []int64 `json:"userIds"`
	OrgId   int64   `json:"orgId"`
}

type GetUserInfoByUserIdsListRespVo struct {
	vo.Err
	GetUserInfoByUserIdsRespVo *[]GetUserInfoByUserIdsRespVo `json:"data"`
}

type GetUserInfoByUserIdsRespVo struct {
	UserId        int64  `json:"userId"`
	OutUserId     string `json:"outUserId"` //有可能为空
	OrgId         int64  `json:"orgId"`
	OutOrgId      string `json:"outOrgId"` //有可能为空
	Name          string `json:"name"`
	Avatar        string `json:"avatar"`
	HasOutInfo    bool   `json:"hasOutInfo"`
	HasOrgOutInfo bool   `json:"hasOrgOutInfo"`

	OrgUserIsDelete    int `json:"orgUserIsDelete"`           //是否被组织移除
	OrgUserStatus      int `json:"orgUserStatus"`      //用户组织状态
	OrgUserCheckStatus int `json:"orgUserCheckStatus"` //用户组织审核状态
}
