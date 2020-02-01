package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/domain"
)

//更新成员状态
func UpdateOrgMemberStatus(reqVo orgvo.UpdateOrgMemberStatusReq) (*vo.Void, errs.SystemErrorInfo) {
	orgId := reqVo.OrgId
	userId := reqVo.UserId
	input := reqVo.Input

	//校验当前用户是否具有修改成员状态的权限
	authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgUser, consts.RoleOperationModifyStatus)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}

	//如果有权限，修改成员状态
	modifyErr := domain.ModifyOrgMemberStatus(orgId, input.MemberIds, input.Status, userId)
	if modifyErr != nil {
		log.Error(modifyErr)
		return nil, modifyErr
	}
	return &vo.Void{
		ID: orgId,
	}, nil
}

//更新成员检查状态
func UpdateOrgMemberCheckStatus(reqVo orgvo.UpdateOrgMemberCheckStatusReq) (*vo.Void, errs.SystemErrorInfo) {
	orgId := reqVo.OrgId
	userId := reqVo.UserId
	input := reqVo.Input

	//校验当前用户是否具有修改成员状态的权限
	authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgUser, consts.RoleOperationModifyStatus)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}

	//如果有权限，修改成员状态
	modifyErr := domain.ModifyOrgMemberCheckStatus(orgId, input.MemberIds, input.CheckStatus, userId)
	if modifyErr != nil {
		log.Error(modifyErr)
		return nil, modifyErr
	}
	return &vo.Void{
		ID: orgId,
	}, nil
}

func RemoveOrgMember(reqVo orgvo.RemoveOrgMemberReq) (*vo.Void, errs.SystemErrorInfo) {
	orgId := reqVo.OrgId
	userId := reqVo.UserId
	input := reqVo.Input

	//校验当前用户是否具有修改删除成员的权限
	authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgUser, consts.RoleOperationDelete)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}
	//如果有权限，移除成员
	modifyErr := domain.RemoveOrgMember(orgId, input.MemberIds, userId)
	if modifyErr != nil {
		log.Error(modifyErr)
		return nil, modifyErr
	}
	return &vo.Void{
		ID: orgId,
	}, nil
}

func OrgUserList(orgId, userId int64, page, size int, input *vo.OrgUserListReq) (*vo.UserOrganizationList, errs.SystemErrorInfo) {
	//校验当前用户是否具有查看成员的权限
	//authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgUser, consts.RoleOperationView)
	//if authErr != nil {
	//	log.Error(authErr)
	//	return nil, authErr
	//}
	//查询成员角色
	roleUserResp := rolefacade.GetOrgRoleUser(rolevo.GetOrgRoleUserReqVo{
		OrgId: orgId,
	})
	if roleUserResp.Failure() {
		log.Error(roleUserResp.Error())
		return nil, roleUserResp.Error()
	}
	//所有超管用户
	var allUserHaveRoleIds []int64
	for _, v := range roleUserResp.Data {
		if v.RoleLangCode == consts.RoleGroupOrgAdmin {
			allUserHaveRoleIds = append(allUserHaveRoleIds, v.UserId)
		}
	}

	total, info, err := domain.GetOrganizationUserList(orgId, page, size, input, allUserHaveRoleIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//if len(info) == 0 {
	//	return &vo.UserOrganizationList{
	//		Total: int64(total),
	//	}, nil
	//}
	var userIds []int64
	for _, v := range info {
		userIds = append(userIds, v.UserId, v.AuditorId)
	}

	infoVo := &[]*vo.OrganizationUser{}
	copyErr := copyer.Copy(info, infoVo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	userIdsInfo, err := domain.BatchGetUserDetailInfo(userIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	userInfoVos := &[]vo.PersonalInfo{}
	copyErr = copyer.Copy(userIdsInfo, userInfoVos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	userInfoMap := maps.NewMap("ID", *userInfoVos)

	roleGroup := rolefacade.GetOrgRoleList(rolevo.GetOrgRoleListReqVo{
		OrgId: orgId,
	})
	if roleGroup.Failure() {
		log.Error(roleGroup.Error())
		return nil, roleGroup.Error()
	}
	orgMemberRoleInfo := vo.Role{}
	for _, v := range roleGroup.Data {
		if v.LangCode == consts.RoleGroupSpecialMember {
			orgMemberRoleInfo = *v
		}
	}
	//暂时默认一个用户最多只有一个组织角色)
	//userRoleMap := maps.NewMap("UserId", roleUserResp.Data)
	userRoleMap := map[int64]rolevo.RoleUser{}
	for _, datum := range roleUserResp.Data {
		if _, ok := userRoleMap[datum.UserId]; !ok {
			userRoleMap[datum.UserId] = datum
		}
	}
	for k, v := range *infoVo {
		if _, ok := userRoleMap[v.UserID]; ok {
			role := userRoleMap[v.UserID]
			(*infoVo)[k].UserRole = &vo.UserRoleInfo{
				ID:       role.RoleId,
				Name:     role.RoleName,
				LangCode: role.RoleLangCode,
			}
		} else {
			(*infoVo)[k].UserRole = &vo.UserRoleInfo{
				ID:       orgMemberRoleInfo.ID,
				Name:     orgMemberRoleInfo.Name,
				LangCode: orgMemberRoleInfo.LangCode,
			}
		}
		if _, ok := userInfoMap[v.UserID]; ok {
			user := userInfoMap[v.UserID].(vo.PersonalInfo)
			(*infoVo)[k].UserInfo = &user
		}
		if _, ok := userInfoMap[v.AuditorID]; ok {
			user := userInfoMap[v.AuditorID].(vo.PersonalInfo)
			(*infoVo)[k].AuditorInfo = &user
		}
		if v.CheckStatus == consts.AppCheckStatusSuccess && v.AuditTime.String() <= consts.BlankElasticityTime {
			(*infoVo)[k].AuditTime = v.CreateTime
		}
	}
	return &vo.UserOrganizationList{
		Total: int64(total),
		List:  *infoVo,
	}, nil
}

func GetOrgUserInfoListBySourceChannel(reqVo orgvo.GetOrgUserInfoListBySourceChannelReq) (*orgvo.GetOrgUserInfoListBySourceChannelRespData, errs.SystemErrorInfo) {
	userInfoList, total, err := domain.GetOrgUserInfoListBySourceChannel(reqVo.OrgId, reqVo.SourceChannel, reqVo.Page, reqVo.Size)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &orgvo.GetOrgUserInfoListBySourceChannelRespData{
		Total: total,
		List:  userInfoList,
	}, nil
}

func BatchGetUserDetailInfo(userIds []int64) ([]vo.PersonalInfo, errs.SystemErrorInfo) {
	res, err := domain.BatchGetUserDetailInfo(userIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	vos := &[]vo.PersonalInfo{}
	copyErr := copyer.Copy(res, vos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return *vos, nil
}
