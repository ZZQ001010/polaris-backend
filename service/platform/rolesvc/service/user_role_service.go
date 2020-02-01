package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/domain"
)

func GetOrgRoleUser(orgId int64, projectId int64) ([]rolevo.RoleUser, errs.SystemErrorInfo) {
	roleBo, err := domain.GetOrgRoleUser(orgId, projectId, nil)
	if err != nil {
		return nil, err
	}
	roleVo := &[]rolevo.RoleUser{}
	copyErr := copyer.Copy(roleBo, roleVo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return *roleVo, nil
}

func GetOrgAdminUser(orgId int64) ([]int64, errs.SystemErrorInfo) {
	roleBos, err := domain.GetOrgRoleUser(orgId, 0, []string{consts.RoleGroupOrgAdmin, consts.RoleGroupOrgManager})
	if err != nil {
		return nil, err
	}
	userIds := make([]int64, len(*roleBos))
	for i, roleBo := range *roleBos{
		userIds[i] = roleBo.UserId
	}
	userIds = slice.SliceUniqueInt64(userIds)
	return userIds, nil
}

func UpdateUserOrgRole(req rolevo.UpdateUserOrgRoleReqVo) (*vo.Void, errs.SystemErrorInfo) {
	orgId := req.OrgId
	operatorId := req.CurrentUserId
	targetUserId := req.UserId

	//这里根据id查角色
	role, err := domain.GetRole(0, 0, req.RoleId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if role.OrgId == 0 {
		//如果是全局角色，只能修改默认角色(暂时特殊角色只能修改为组织成员)
		if role.LangCode != consts.RoleGroupSpecialMember {
			log.Error("目标特殊角色只能是组织成员")
			return nil, errs.NoOperationPermissions
		}
	} else {
		if role.OrgId != orgId || role.LangCode == consts.RoleGroupOrgAdmin {
			log.Error("不能修改为组织管理员")
			return nil, errs.NoOperationPermissions
		}
	}

	if req.ProjectId == nil || *req.ProjectId == 0 {
		//判断组织权限
		authErr := AuthOrgRole(req.OrgId, req.CurrentUserId, consts.RoleOperationPathOrgUser, consts.RoleOperationBind)
		if authErr != nil {
			log.Error(authErr)
			return nil, authErr
		}

		// 组织角色逻辑判断，当前用户必须角色属性高于被修改角色（目前只有超管可以修改角色的权限）
		targetUserAdminFlag, err := GetUserAdminFlag(orgId, targetUserId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if targetUserAdminFlag.IsAdmin {
			log.Error("超管的角色不允许修改")
			return nil, errs.OrgUserRoleModifyError
		}

		operatorAdminFlag, err := GetUserAdminFlag(orgId, operatorId)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		if operatorAdminFlag.IsManager && !operatorAdminFlag.IsAdmin {
			log.Error("没有权限修改，因为只有超管才能修改管理员的角色")
			return nil, errs.NoOperationPermissions
		}
	} else {
		if role.ProjectId != 0 && role.ProjectId != *req.ProjectId {
			log.Error("角色不属于该项目")
			return nil, errs.NoOperationPermissions
		}
		//判断项目权限
		authResp := projectfacade.AuthProjectPermission(projectvo.AuthProjectPermissionReqVo{
			Input: projectvo.AuthProjectPermissionReqData{
				OrgId:      req.OrgId,
				UserId:     req.CurrentUserId,
				ProjectId:  *req.ProjectId,
				Path:       consts.RoleOperationPathOrgProRole,
				Operation:  consts.RoleOperationBind,
				AuthFiling: true,
			},
		})
		if authResp.Failure() {
			log.Error(authResp.Message)
			return nil, authResp.Error()
		}
	}

	id, updErr := domain.UpdateUserOrgRole(role, orgId, req.CurrentUserId, req.UserId, req.ProjectId)
	if updErr != nil {
		log.Error(updErr)
		return nil, updErr
	}

	return &vo.Void{
		ID: id,
	}, nil
}

func GetUserAdminFlag(orgId, userId int64) (*bo.UserAdminFlagBo, errs.SystemErrorInfo) {
	userRole, err := GetUserRoleList(orgId, userId, 0)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	//roleList, err := GetRoleList(orgId)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var adminRole *bo.RoleBo = nil
	//var managerRole *bo.RoleBo = nil
	//for _, role := range roleList {
	//	currentRole := role
	//	if strings.EqualFold(role.LangCode, consts.RoleGroupOrgAdmin) {
	//		adminRole = &currentRole
	//	}
	//	if strings.EqualFold(role.LangCode, consts.RoleGroupOrgManager) {
	//		managerRole = &currentRole
	//	}
	//}
	adminRole, err := GetRoleByLangCode(orgId, consts.RoleGroupOrgAdmin)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
	}
	managerRole, err := GetRoleByLangCode(orgId, consts.RoleGroupOrgManager)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
	}
	isAdmin := false
	isManager := false
	for _, v := range *userRole {
		if adminRole != nil && v.RoleId == adminRole.Id {
			isAdmin = true
		}
		if managerRole != nil && v.RoleId == managerRole.Id {
			isManager = true
		}
	}
	return &bo.UserAdminFlagBo{
		IsAdmin:   isAdmin,
		IsManager: isManager,
	}, nil
}

//获取组织角色列表
func GetOrgRoleList(orgId int64) ([]*vo.Role, errs.SystemErrorInfo) {
	groupRole, err := GetRoleListByGroup(orgId, consts.RoleGroupOrg, 0)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	orgMember, err := GetRoleByLangCode(0, consts.RoleGroupSpecialMember)
	if err != nil {
		return nil, err
	}
	if orgMember != nil {
		//替换下名称（奇怪吧）
		remark := orgMember.Remark
		name := orgMember.Name
		orgMember.Name = remark
		orgMember.Remark = name
		groupRole = append(groupRole, *orgMember)
	}

	resVo := &[]*vo.Role{}
	copyErr := copyer.Copy(groupRole, resVo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return *resVo, nil
}

//获取项目所有角色
func GetProjectRoleList(orgId int64, projectId int64) ([]*vo.Role, errs.SystemErrorInfo) {
	groupRole, err := GetRoleListByGroup(orgId, consts.RoleGroupPro, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//负责人
	ownerRole, err := GetRoleByLangCode(0, consts.RoleGroupSpecialOwner)
	if err != nil {
		return nil, err
	}
	newGroup := []bo.RoleBo{}
	if ownerRole != nil {
		//替换下名称（奇怪吧）
		remark := ownerRole.Remark
		name := ownerRole.Name
		ownerRole.Name = remark
		ownerRole.Remark = name
		newGroup = append(newGroup, *ownerRole)
	}

	for _, roleBo := range groupRole {
		//该项目的角色和项目成员
		if roleBo.ProjectId == projectId || roleBo.LangCode == consts.RoleGroupProMember {
			newGroup = append(newGroup, roleBo)
		}
	}

	resVo := &[]*vo.Role{}
	copyErr := copyer.Copy(newGroup, resVo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return *resVo, nil
}
