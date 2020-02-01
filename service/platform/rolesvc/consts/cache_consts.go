package consts

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
)

var (
	//角色列表
	CacheRoleList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + "role_list"
	////用户角色列表
	//CacheUserRoleList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "user_role_list"
	//角色权限列表
	CacheRolePermissionOperationList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfRole + "role_permission_list"
	//角色操作列表
	CacheRoleOperationList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfSys + "role_operation_list"

	//补偿的角色列表
	CacheCompensatoryRolePermissionPathList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + "compensatory_role_permission_path_list"
	//角色组信息
	CacheRoleGroupList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + "role_group_list"
	//权限项列表
	CachePermissionList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfSys + "permission_list"
	//权限项操作列表
	CachePermissionOperationList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfSys + "permission_operation_list"
	//用户角色列表(hash)
	CacheUserRoleListHash = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "user_role_list_hash"
	//角色列表
	CacheRoleListHash = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + "role_list_hash"
)
