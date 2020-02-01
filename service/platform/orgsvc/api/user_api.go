package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/service"
)

func (GetGreeter) PersonalInfo(req orgvo.PersonalInfoReqVo) orgvo.PersonalInfoRespVo {
	res, err := service.PersonalInfo(req.OrgId, req.UserId, req.SourceChannel)
	return orgvo.PersonalInfoRespVo{Err: vo.NewErr(err), PersonalInfo: res}
}

func (PostGreeter) GetUserIds(reqVo orgvo.GetUserIdsReqVo) orgvo.GetUserIdsRespVo {
	res, err := service.GetUserIds(reqVo.OrgId, reqVo.CorpId, reqVo.SourceChannel, reqVo.EmpIdsBody.EmpIds)
	return orgvo.GetUserIdsRespVo{Err: vo.NewErr(err), GetUserIds: res}
}

func (GetGreeter) GetUserId(reqVo orgvo.GetUserIdReqVo) orgvo.GetUserIdRespVo {
	res, err := service.GetUserId(reqVo.OrgId, reqVo.CorpId, reqVo.SourceChannel, reqVo.EmpId)
	return orgvo.GetUserIdRespVo{Err: vo.NewErr(err), GetUserId: res}
}

func (GetGreeter) UserConfigInfo(reqVo orgvo.UserConfigInfoReqVo) orgvo.UserConfigInfoRespVo {
	res, err := service.UserConfigInfo(reqVo.OrgId, reqVo.UserId)
	return orgvo.UserConfigInfoRespVo{Err: vo.NewErr(err), UserConfigInfo: res}
}

func (PostGreeter) UpdateUserConfig(req orgvo.UpdateUserConfigReqVo) orgvo.UpdateUserConfigRespVo {
	res, err := service.UpdateUserConfig(req.OrgId, req.UserId, req.UpdateUserConfigReq)
	return orgvo.UpdateUserConfigRespVo{Err: vo.NewErr(err), UpdateUserConfig: res}
}

func (PostGreeter) UpdateUserPcConfig(req orgvo.UpdateUserPcConfigReqVo) orgvo.UpdateUserConfigRespVo {
	res, err := service.UpdateUserPcConfig(req.OrgId, req.UserId, req.UpdateUserPcConfigReq)
	return orgvo.UpdateUserConfigRespVo{Err: vo.NewErr(err), UpdateUserConfig: res}
}

func (PostGreeter) UpdateUserDefaultProjectIdConfig(req orgvo.UpdateUserDefaultProjectIdConfigReqVo) orgvo.UpdateUserConfigRespVo {
	res, err := service.UpdateUserDefaultProjectIdConfig(req.OrgId, req.UserId, req.UpdateUserDefaultProjectIdConfigReq)
	return orgvo.UpdateUserConfigRespVo{Err: vo.NewErr(err), UpdateUserConfig: res}
}

func (GetGreeter) VerifyOrg(reqVo orgvo.VerifyOrgReqVo) vo.BoolRespVo {
	return vo.BoolRespVo{Err: vo.NewErr(nil), IsTrue: service.VerifyOrg(reqVo.OrgId, reqVo.UserId)}
}

func (PostGreeter) VerifyOrgUsers(reqVo orgvo.VerifyOrgUsersReqVo) vo.BoolRespVo {
	return vo.BoolRespVo{Err: vo.NewErr(nil), IsTrue: service.VerifyOrgUsers(reqVo.OrgId, reqVo.Input.UserIds)}
}

func (GetGreeter) GetUserInfo(reqVo orgvo.GetUserInfoReqVo) orgvo.GetUserInfoRespVo {
	res, err := service.GetUserInfo(reqVo.OrgId, reqVo.UserId, reqVo.SourceChannel)
	return orgvo.GetUserInfoRespVo{Err: vo.NewErr(err), UserInfo: res}
}

func (GetGreeter) GetOutUserInfoListBySourceChannel(req orgvo.GetOutUserInfoListBySourceChannelReqVo) orgvo.GetOutUserInfoListBySourceChannelRespVo {
	res, err := service.GetOutUserInfoListBySourceChannel(req.SourceChannel, req.Page, req.Size)
	return orgvo.GetOutUserInfoListBySourceChannelRespVo{Err: vo.NewErr(err), UserOutInfoBoList: res}
}

func (GetGreeter) GetUserInfoListByOrg(reqVo orgvo.GetUserInfoListReqVo) orgvo.GetUserInfoListRespVo {
	res, err := service.GetUserInfoListByOrg(reqVo.OrgId)
	return orgvo.GetUserInfoListRespVo{Err: vo.NewErr(err), SimpleUserInfo: res}
}

func (PostGreeter) UpdateUserInfo(req orgvo.UpdateUserInfoReqVo) vo.CommonRespVo {
	res, err := service.UpdateUserInfo(req.OrgId, req.UserId, req.UpdateUserInfoReq)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) GetUserInfoByUserIds(input orgvo.GetUserInfoByUserIdsReqVo) orgvo.GetUserInfoByUserIdsListRespVo {
	res, err := service.GetUserInfoByUserIds(input)
	return orgvo.GetUserInfoByUserIdsListRespVo{Err: vo.NewErr(err), GetUserInfoByUserIdsRespVo: res}
}

func (PostGreeter) BatchGetUserDetailInfo(reqVo orgvo.BatchGetUserInfoReq) orgvo.BatchGetUserInfoResp {
	res, err := service.BatchGetUserDetailInfo(reqVo.UserIds)
	return orgvo.BatchGetUserInfoResp{Data: res, Err: vo.NewErr(err)}
}
