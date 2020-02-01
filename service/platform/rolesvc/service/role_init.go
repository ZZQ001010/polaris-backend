package service

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/po"
	"strconv"
	"upper.io/db.v3/lib/sqlbuilder"
)

const RoleInitSql = consts.TemplateDirPrefix + "role_init.template"

func RoleInit(orgId int64, tx sqlbuilder.Tx) (*bo.RoleInitResp, errs.SystemErrorInfo) {
	maps := map[string]interface{}{}
	maps["OrgId"] = orgId

	//ppm_rol_role_group主键封装
	roleGroup := &po.PpmRolRoleGroup{}
	roleGroupIds, err := idfacade.ApplyMultipleIdRelaxed(0, roleGroup.TableName(), "", 3)
	if err != nil {
		return nil, err
	}
	for k, v := range roleGroupIds.Ids {
		maps["RoleGroupId"+strconv.Itoa(k+1)] = v.Id
	}
	logger.GetDefaultLogger().Infof("ppm_rol_role_group主键分配完成")

	//ppm_rol_role主键封装
	role := &po.PpmRolRole{}
	roleIds, err := idfacade.ApplyMultipleIdRelaxed(0, role.TableName(), "", 14)
	if err != nil {
		return nil, err
	}
	for k, v := range roleIds.Ids {
		maps["RoleId"+strconv.Itoa(k+1)] = v.Id
	}
	//获取组织超级管理员角色id和普通管理员角色id
	orgSuperAdminRoleId := maps["RoleId7"].(int64)
	orgNormalAdminRoleId := maps["RoleId8"].(int64)
	logger.GetDefaultLogger().Infof("ppm_rol_role主键分配完成")

	//ppm_rol_role_permission_operation主键封装
	rolePermissionOperation := &po.PpmRolRolePermissionOperation{}
	rolePermissionOperationIds, err := idfacade.ApplyMultipleIdRelaxed(0, rolePermissionOperation.TableName(), "", 163)
	if err != nil {
		return nil, err
	}
	for k, v := range rolePermissionOperationIds.Ids {
		maps["PermissionOperation"+strconv.Itoa(k+1)] = v.Id
	}
	logger.GetDefaultLogger().Infof("ppm_rol_role_permission_operation主键分配完成")

	//ppm_rol_permission主键封装
	permission := &po.PpmRolPermission{}
	permissionIds, err := idfacade.ApplyMultipleIdRelaxed(0, permission.TableName(), "", 29)
	if err != nil {
		return nil, err
	}
	for k, v := range permissionIds.Ids {
		maps["PermissionId"+strconv.Itoa(k+1)] = v.Id
	}
	logger.GetDefaultLogger().Infof("ppm_rol_permission主键分配完成")

	//ppm_rol_permission_operation主键封装
	rolPermissionOperation := &po.PpmRolPermission{}
	rolPermissionOperationIds, err := idfacade.ApplyMultipleIdRelaxed(0, rolPermissionOperation.TableName(), "", 119)
	if err != nil {
		return nil, err
	}
	for k, v := range rolPermissionOperationIds.Ids {
		maps["RolPermissionOperationId"+strconv.Itoa(k+1)] = v.Id
	}
	logger.GetDefaultLogger().Infof("ppm_rol_permission_operation主键分配完成")

	err = util.ReadAndWrite(RoleInitSql, maps, tx)
	if err != nil {
		return nil, err
	}

	return &bo.RoleInitResp{
		OrgSuperAdminRoleId:  orgSuperAdminRoleId,
		OrgNormalAdminRoleId: orgNormalAdminRoleId,
	}, nil
}
