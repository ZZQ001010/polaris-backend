package orgvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type GetBaseOrgInfoReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	OrgId         int64  `json:"orgId"`
}

type GetBaseOrgInfoByOutOrgIdReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	OutOrgId         string  `json:"outOrgId"`
}

type GetBaseOrgInfoRespVo struct {
	vo.Err
	BaseOrgInfo *bo.BaseOrgInfoBo `json:"data"`
}

type GetBaseOrgInfoByOutOrgIdRespVo struct {
	vo.Err
	BaseOrgInfo *bo.BaseOrgInfoBo `json:"data"`
}

type GetDingTalkBaseUserInfoByEmpIdReqVo struct {
	OrgId int64  `json:"orgId"`
	EmpId string `json:"empId"`
}

type GetBaseUserInfoByEmpIdReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	OrgId int64  `json:"orgId"`
	EmpId string `json:"empId"`
}

type GetBaseUserInfoByEmpIdRespVo struct {
	vo.Err
	BaseUserInfo *bo.BaseUserInfoBo `json:"data"`
}

type GetDingTalkBaseUserInfoByEmpIdRespVo struct {
	vo.Err
	DingTalkBaseUserInfo *bo.BaseUserInfoBo `json:"data"`
}

type GetUserConfigInfoReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetUserConfigInfoRespVo struct {
	vo.Err
	UserConfigInfo *bo.UserConfigBo `json:"data"`
}

type GetBaseUserInfoReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	OrgId         int64  `json:"orgId"`
	UserId        int64  `json:"userId"`
}

type GetBaseUserInfoRespVo struct {
	vo.Err
	BaseUserInfo *bo.BaseUserInfoBo `json:"data"`
}

type GetDingTalkBaseUserInfoReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetBaseUserInfoBatchReqVo struct {
	SourceChannel string  `json:"sourceChannel"`
	OrgId         int64   `json:"orgId"`
	UserIds       []int64 `json:"userIds"`
}

type GetBaseUserInfoBatchRespVo struct {
	vo.Err
	BaseUserInfos []bo.BaseUserInfoBo `json:"data"`
}
