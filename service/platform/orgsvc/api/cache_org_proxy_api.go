package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/service"
)

func (GetGreeter) GetBaseOrgInfo(reqVo orgvo.GetBaseOrgInfoReqVo) orgvo.GetBaseOrgInfoRespVo {
	res, err := service.GetBaseOrgInfo(reqVo.SourceChannel, reqVo.OrgId)
	return orgvo.GetBaseOrgInfoRespVo{Err: vo.NewErr(err), BaseOrgInfo: res}
}

func (GetGreeter) GetBaseOrgInfoByOutOrgId(reqVo orgvo.GetBaseOrgInfoByOutOrgIdReqVo) orgvo.GetBaseOrgInfoByOutOrgIdRespVo{
	res, err := service.GetBaseOrgInfoByOutOrgId(reqVo.SourceChannel, reqVo.OutOrgId)
	return orgvo.GetBaseOrgInfoByOutOrgIdRespVo{Err: vo.NewErr(err), BaseOrgInfo: res}
}

func (GetGreeter) GetDingTalkBaseUserInfoByEmpId(reqVo orgvo.GetDingTalkBaseUserInfoByEmpIdReqVo) orgvo.GetDingTalkBaseUserInfoByEmpIdRespVo {
	res, err := service.GetDingTalkBaseUserInfoByEmpId(reqVo.OrgId, reqVo.EmpId)
	return orgvo.GetDingTalkBaseUserInfoByEmpIdRespVo{Err: vo.NewErr(err), DingTalkBaseUserInfo: res}
}

func (GetGreeter) GetBaseUserInfoByEmpId(reqVo orgvo.GetBaseUserInfoByEmpIdReqVo) orgvo.GetBaseUserInfoByEmpIdRespVo {
	res, err := service.GetBaseUserInfoByEmpId(reqVo.SourceChannel, reqVo.OrgId, reqVo.EmpId)
	return orgvo.GetBaseUserInfoByEmpIdRespVo{Err: vo.NewErr(err), BaseUserInfo: res}
}

func (GetGreeter) GetUserConfigInfo(reqVo orgvo.GetUserConfigInfoReqVo) orgvo.GetUserConfigInfoRespVo {
	res, err := service.GetUserConfigInfo(reqVo.OrgId, reqVo.UserId)
	return orgvo.GetUserConfigInfoRespVo{Err: vo.NewErr(err), UserConfigInfo: res}
}

func (GetGreeter) GetBaseUserInfo(reqVo orgvo.GetBaseUserInfoReqVo) orgvo.GetBaseUserInfoRespVo {
	res, err := service.GetBaseUserInfo(reqVo.SourceChannel, reqVo.OrgId, reqVo.UserId)
	return orgvo.GetBaseUserInfoRespVo{Err: vo.NewErr(err), BaseUserInfo: res}
}

func (GetGreeter) GetDingTalkBaseUserInfo(reqVo orgvo.GetDingTalkBaseUserInfoReqVo) orgvo.GetBaseUserInfoRespVo {
	res, err := service.GetDingTalkBaseUserInfo(reqVo.OrgId, reqVo.UserId)
	return orgvo.GetBaseUserInfoRespVo{Err: vo.NewErr(err), BaseUserInfo: res}
}

func (PostGreeter) GetBaseUserInfoBatch(reqVo orgvo.GetBaseUserInfoBatchReqVo) orgvo.GetBaseUserInfoBatchRespVo {
	res, err := service.GetBaseUserInfoBatch(reqVo.SourceChannel, reqVo.OrgId, reqVo.UserIds)
	return orgvo.GetBaseUserInfoBatchRespVo{Err: vo.NewErr(err), BaseUserInfos: res}
}
