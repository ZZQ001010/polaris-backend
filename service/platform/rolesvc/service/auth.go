package service

import (
	"fmt"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/rolesvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/domain"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/po"
	"strconv"
	"strings"
)

func Authenticate(orgId int64, userId int64, projectAuthInfo *bo.ProjectAuthBo, issueAuthInfo *bo.IssueAuthBo, path string, operation string) errs.SystemErrorInfo {
	//如果是issue负责人，拥有所有权限
	if issueAuthInfo != nil && userId == issueAuthInfo.Owner {
		log.Infof("权限校验成功，用户 %d 是任务 %d 的负责人", userId, issueAuthInfo.Id)
		return nil
	}
	projectId := int64(0)

	//如果是项目负责人，拥有所有权限
	if projectAuth(&projectId, userId, projectAuthInfo) {
		return nil
	}

	//验证当前人是否是超管
	userAdminFlag, err := GetUserAdminFlag(orgId, userId)
	if err != nil {
		log.Error(err)
		return err
	}
	if userAdminFlag.IsAdmin {
		log.Infof("权限校验成功，用户 %d 是组织 %d 的超级管理员", userId, orgId)
		return nil
	}

	//判断操作编号是否合法
	_, err = GetRoleOperationByCode(operation)
	if err != nil {
		log.Errorf("操作%s不合法", operation)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	authPath := path
	checkCompensationErr := checkCompensation(orgId, authPath)
	if checkCompensationErr != nil {
		log.Error(checkCompensationErr)
	}

	//path替换
	path = strings.ReplaceAll(path, "{org_id}", strconv.FormatInt(orgId, 10))
	//specificPah用来判断是否存在对项目做了特殊权限处理, 优先查找specificPah对应的权限项
	specificPah := strings.ReplaceAll(path, "{pro_id}", strconv.FormatInt(projectId, 10))
	//normalPath用来查询默认的权限
	normalPath := strings.ReplaceAll(path, "{pro_id}", "0")

	//获取用户所有角色
	roleIds, err1 := GetUserRoleIds(orgId, userId, projectId, projectAuthInfo, issueAuthInfo)
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.GetUserRoleError, err1)
	}

	//开始查找角色是否有相应的权限
	rolePermissionOperations, err := GetRolePermissionOperationListByPath(orgId, roleIds, specificPah)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if len(*rolePermissionOperations) == 0 && projectId != 0 {
		//寻找全局项目权限
		rolePermissionOperations, err = GetRolePermissionOperationListByPath(orgId, roleIds, normalPath)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}

	//校验权限
	//for _, permission := range *rolePermissionOperations {
	//	operationCodes := strings.Trim(permission.OperationCodes, " ")
	//	if operationCodes == "*" {
	//		log.Infof("权限校验通过 userId %d, roleId %d, operationCodes %s, path %s", userId, permission.RoleId, permission.OperationCodes, permission.PermissionPath)
	//		return nil
	//	} else {
	//		mat, err1 := regexp.MatchString(operation, permission.OperationCodes)
	//		if err1 != nil {
	//			log.Errorf("权限匹配，正则校验异常 %v", err1)
	//			continue
	//		}
	//		if mat {
	//			log.Infof("权限校验通过 userId %d, roleId %d, operationCodes %s, path %s", userId, permission.RoleId, permission.OperationCodes, permission.PermissionPath)
	//			return nil
	//		}
	//	}
	//}

	//如果为true 权限校验通过
	if checkPermission(rolePermissionOperations, userId, operation) {
		return nil
	}

	return errs.BuildSystemErrorInfo(errs.NoOperationPermissions)
}

//检测补偿
func checkCompensation(orgId int64, authPath string) errs.SystemErrorInfo {
	if _, ok := sconsts.RolePermissionOperationDefineMap[authPath]; !ok {
		return nil
	}

	paths, err := GetCompensatoryRolePermissionPaths(orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	if exist, _ := slice.Contain(paths, authPath); exist {
		return nil
	}
	lockKey := fmt.Sprintf("%s%d", consts.RolePermissionPathCompensatoryLockKey, orgId)
	uUid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uUid)
	if lockErr != nil {
		log.Error(lockErr)
		return errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uUid); err != nil {
				log.Error(err)
			}
		}()

		//二次校验
		paths, err := GetCompensatoryRolePermissionPaths(orgId)
		if err != nil {
			log.Error(err)
			return err
		}
		if exist, _ := slice.Contain(paths, authPath); !exist {
			//这时要开始做补偿
			//roleList, err := GetRoleList(orgId, 0)
			roleList, err := domain.GetSysDefaultRoles(orgId)
			if err != nil {
				log.Error(err)
				return err
			}
			roleMap := maps.NewMap("LangCode", roleList)

			//组装po对象
			insertPos := make([]po.PpmRolRolePermissionOperation, 0)
			for k, rp := range sconsts.RolePermissionOperationDefineMap {
				if exist, _ := slice.Contain(paths, k); !exist {
					paths = append(paths, k)

					len := len(rp.RoleLangCodes)
					if len > 0 {
						permissionPath := strings.ReplaceAll(k, "{org_id}", strconv.FormatInt(orgId, 10))
						for _, roleLangCode := range rp.RoleLangCodes {
							if roleInfo, ok := roleMap[roleLangCode]; ok {
								roleCacheBo := roleInfo.(bo.RoleBo)
								insertPos = append(insertPos, po.PpmRolRolePermissionOperation{
									OrgId:          orgId,
									RoleId:         roleCacheBo.Id,
									ProjectId:      0,
									PermissionId:   rp.RolePermissionId,
									PermissionPath: permissionPath,
									OperationCodes: rp.OperationCodes,
								})
							}
						}
					}
				}
			}
			//插入操作
			if len(insertPos) > 0 {
				len := len(insertPos)
				ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableRolePermissionOperation, len)
				if err != nil {
					log.Error(err)
					return err
				}
				for i, _ := range insertPos {
					insertPos[i].Id = ids.Ids[i].Id
				}
				//设置缓存
				cacheErr := SetCompensatoryRolePermissionPaths(orgId, paths)
				if cacheErr != nil {
					log.Error(cacheErr)
					return cacheErr
				}

				dbErr := mysql.BatchInsert(&po.PpmRolRolePermissionOperation{}, slice.ToSlice(insertPos))
				if dbErr != nil {
					log.Error(dbErr)
					return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
				}
			}
		}
	}

	return nil
}

//传递指针引用更改值
func projectAuth(projectId *int64, userId int64, projectAuthInfo *bo.ProjectAuthBo) bool {

	if projectAuthInfo != nil {
		*projectId = projectAuthInfo.Id
		if userId == projectAuthInfo.Owner {
			log.Infof("权限校验成功，用户 %d 是项目 %d 的负责人", userId, projectAuthInfo.Id)
			return true
		}
	}
	return false
}

//校验权限
func checkPermission(rolePermissionOperations *[]bo.RolePermissionOperationBo, userId int64, operation string) bool {

	for _, permission := range *rolePermissionOperations {
		operationCodes := strings.Trim(permission.OperationCodes, " ")
		if operationCodes == "*" {
			log.Infof("权限校验通过 userId %d, roleId %d, operationCodes %s, path %s", userId, permission.RoleId, permission.OperationCodes, permission.PermissionPath)
			return true
		} else {
			//mat, err1 := regexp.MatchString(operation, permission.OperationCodes)
			//if err1 != nil {
			//	log.Errorf("权限匹配，正则校验异常 %v", err1)
			//	continue
			//}
			if util.RoleOperationCodesMatch(operation, permission.OperationCodes) {
				log.Infof("权限校验通过 userId %d, roleId %d, operationCodes %s, path %s", userId, permission.RoleId, permission.OperationCodes, permission.PermissionPath)
				return true
			}
		}
	}
	return false
}

func GetUserRoleIds(orgId, userId, projectId int64, projectAuthInfo *bo.ProjectAuthBo, issueAuthInfo *bo.IssueAuthBo) ([]int64, error) {
	//开始加载用户所有的角色
	roleUsers, err1 := GetUserRoleListByProjectId(orgId, userId, projectId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
	}
	noSpecific := len(*roleUsers) == 0 && projectId != 0

	//如果没有对项目增加特殊权限，查询通用权限
	if noSpecific {
		roleUsers, err1 = GetUserRoleListByProjectId(orgId, userId, 0)
		if err1 != nil {
			log.Error(err1)
			return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
		}
	}

	roleIds := make([]int64, 0)
	for _, roleUser := range *roleUsers {
		roleIds = append(roleIds, roleUser.RoleId)
	}

	//如果没有对项目做单独的权限设定或者查询的是通用权限，那么增加特殊角色
	err := assemblyRoleIds(&roleIds, projectId, noSpecific, orgId, userId, projectAuthInfo, issueAuthInfo)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return roleIds, nil
}

func assemblyRoleIds(roleIds *[]int64, projectId int64, noSpecific bool, orgId, userId int64, projectAuthInfo *bo.ProjectAuthBo, issueAuthInfo *bo.IssueAuthBo) (error error) {

	if projectId == 0 || noSpecific {
		//如果任务信息不为空，则获取任务相关角色，否则获取项目角色
		//if issueAuthInfo != nil {
		//	err := loadRoleByProjectAuthInfo(roleIds, orgId, userId, issueAuthInfo.Creator, issueAuthInfo.Followers, issueAuthInfo.Participants)
		//	if err != nil {
		//		log.Error(err)
		//		return err
		//	}
		//}
		if projectAuthInfo != nil {
			//加载创建者相关信息,参与者 ,关注者
			err := loadRoleByProjectAuthInfo(roleIds, orgId, userId, projectAuthInfo.Creator, projectAuthInfo.Followers, projectAuthInfo.Participants)
			if err != nil {
				log.Error(err)
				return err
			}
		}
		//组织成员
		if len(*roleIds) == 0 {
			orgMemberRole, err := GetRoleByLangCode(orgId, consts.RoleGroupSpecialMember)
			if err != nil {
				log.Error(err)
				return errs.BuildSystemErrorInfo(errs.CacheProxyError)
			}
			*roleIds = append(*roleIds, orgMemberRole.Id)
		}

		//访客
		visitorRole, err := GetRoleByLangCode(orgId, consts.RoleGroupSpecialVisitor)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError)
		}
		*roleIds = append(*roleIds, visitorRole.Id)
	}
	return nil
}

func loadRoleByProjectAuthInfo(roleIds *[]int64, orgId, userId int64, creator int64, followers []int64, participants []int64) (error error) {

	//加载创建者相关信息
	creatorIsError, creatorError := loadCreateRole(roleIds, orgId, userId, creator)
	if creatorIsError {
		log.Error(creatorIsError)
		return creatorError
	}

	if participants != nil && len(participants) > 0 {
		//参与者
		participantIsError, ParticipantError := loadParticipantRole(roleIds, orgId, userId, participants)

		if participantIsError {
			log.Error(participantIsError)
			return ParticipantError
		}
	}

	if followers != nil && len(followers) > 0 {
		//关注者
		followIsError, followError := loadFollowRole(roleIds, orgId, userId, followers)
		if followIsError {
			log.Error(followIsError)
			return followError
		}
	}

	return nil
}

//加载创建者相关角色
func loadCreateRole(roleIds *[]int64, orgId, userId int64, creator int64) (isNil bool, error error) {

	if userId == creator {
		createRole, err := GetRoleByLangCode(orgId, consts.RoleGroupSpecialCreator)
		if err != nil {
			log.Error(err)
			return true, errs.BuildSystemErrorInfo(errs.CacheProxyError)
		}
		*roleIds = append(*roleIds, createRole.Id)
	}
	//不需要返回 异常为空
	return false, nil
}

//判断是否是参与者 和加载参与者权限
func loadParticipantRole(roleIds *[]int64, orgId, userId int64, participants []int64) (isNil bool, error error) {

	isParticipant, err2 := slice.Contain(participants, userId)
	if err2 != nil {
		log.Error(err2)
		return true, errs.BuildSystemErrorInfo(errs.SystemError)
	}
	if isParticipant {
		//workerRole, err := GetRoleByLangCode(orgId, consts.RoleGroupSpecialWorker)
		workerRole, err := GetRoleByLangCode(orgId, consts.RoleGroupProMember)
		if err != nil {
			log.Error(err)
			return true, errs.BuildSystemErrorInfo(errs.CacheProxyError)
		}
		*roleIds = append(*roleIds, workerRole.Id)
	}
	return false, nil
}

//关注着
func loadFollowRole(roleIds *[]int64, orgId, userId int64, followers []int64) (isNil bool, error error) {

	isFollower, err2 := slice.Contain(followers, userId)
	if err2 != nil {
		log.Error(err2)
		return true, errs.BuildSystemErrorInfo(errs.SystemError)
	}
	if isFollower {
		//attentionRole, err := GetRoleByLangCode(orgId, consts.RoleGroupSpecialAttention)
		attentionRole, err := GetRoleByLangCode(orgId, consts.RoleGroupProMember)
		if err != nil {
			log.Error(err)
			return true, errs.BuildSystemErrorInfo(errs.CacheProxyError)
		}
		*roleIds = append(*roleIds, attentionRole.Id)
	}
	return false, nil
}
