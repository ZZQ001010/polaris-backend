package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/md5"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/pinyin"
	"github.com/galaxy-book/polaris-backend/common/extra/feishu"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
	"upper.io/db.v3/lib/sqlbuilder"
)

/**
orgId: 组织id
fsDepIdMap: 飞书部门id-内部部门id
*/
func InitFsUserList(orgId int64, tenantKey string, fsDepIdMap map[string]int64, superAdminRoleId, normalAdminRoleId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo{
	tenant, err := feishu.GetTenant(tenantKey)
	if err != nil{
		log.Error(err)
		return err
	}

	userIdCache := map[string]bool{}
	userPoList := make([]po.PpmOrgUser, 0)
	userOutPoList := make([]po.PpmOrgUserOutInfo, 0)
	userConfigList := make([]po.PpmOrgUserConfig, 0)
	userOrgList := make([]po.PpmOrgUserOrganization, 0)
	userDepList := make([]po.PpmOrgUserDepartment, 0)

	//根部门id
	rootDepId := fsDepIdMap["0"]
	//先排除根部门, 因为飞书不支持查询根部门下的用户
	delete(fsDepIdMap, "0")


	surplusOpenIds, err := feishu.GetScopeOpenIds(tenantKey)
	if err != nil{
		log.Error(err)
		return err
	}

	hasSuperAdmin := false

	batch := 30
	offset := 0
	surplusSize := len(surplusOpenIds)
	if surplusSize > 0{
		for ;; {
			limit := offset + batch
			if surplusSize < limit{
				limit = surplusSize
			}
			openIds := surplusOpenIds[offset: limit]

			userBatchResp, err := tenant.GetUserBatchGetV2(nil, openIds)
			if err != nil{
				log.Error(err)
				return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}
			if userBatchResp.Code != 0{
				log.Error(userBatchResp.Msg)
				return errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}

			userPoIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUser, len(openIds))
			if idErr != nil{
				log.Error(idErr)
				return idErr
			}
			for i, fsUser := range userBatchResp.Data.Users{
				userNativeId := userPoIds.Ids[i].Id
				if _, ok := userIdCache[fsUser.OpenId]; ! ok{
					userPo := assemblyFeiShuUserInfo(orgId, fsUser)
					userPo.Id = userNativeId
					userPoList = append(userPoList, userPo)
					userOutPoList = append(userOutPoList, assemblyFeiShuUserOutInfo(orgId, userNativeId, fsUser))
					userConfigList = append(userConfigList, assemblyFeiShuUserConfigInfo(orgId, userNativeId))
					userOrgList = append(userOrgList, assemblyFeiShuUserOrgRelationInfo(orgId, userNativeId))

					if fsUser.IsTenantManager{
						roleId := superAdminRoleId
						if hasSuperAdmin{
							roleId = normalAdminRoleId
						}

						err2 := InitFsManager(orgId, userNativeId, roleId, tx)
						if err2 != nil {
							log.Error(err2)
							return err2
						}

						hasSuperAdmin = true
					}
				}
				userDeps := fsUser.Departments
				if userDeps != nil && len(userDeps) > 0{
					hasDep := false
					for _, userDep := range userDeps{
						if depId, ok := fsDepIdMap[userDep]; ok{
							//这些用户归入部门
							userDepList = append(userDepList, assemblyFeiShuUserDepRelationInfo(orgId, userNativeId, depId))
							hasDep = true
						}
					}
					if !hasDep{
						userDepList = append(userDepList, assemblyFeiShuUserDepRelationInfo(orgId, userNativeId, rootDepId))
					}
				}
				userIdCache[fsUser.OpenId] = true
			}

			if surplusSize <= limit{
				break
			}
			offset += batch
		}
	}

	if len(userPoList) > 0{
		userOutPoIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserOutInfo, len(userOutPoList))
		if idErr != nil{
			log.Error(idErr)
			return idErr
		}
		for i, _ := range userOutPoList{
			userOutPoList[i].Id = userOutPoIds.Ids[i].Id
		}

		userConfigPoIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserConfig, len(userConfigList))
		if idErr != nil{
			log.Error(idErr)
			return idErr
		}
		for i, _ := range userConfigList{
			userConfigList[i].Id = userConfigPoIds.Ids[i].Id
		}

		userOrgPoIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserOrganization, len(userOrgList))
		if idErr != nil{
			log.Error(idErr)
			return idErr
		}
		for i, _ := range userOrgList{
			userOrgList[i].Id = userOrgPoIds.Ids[i].Id
		}

		batchInsert := mysql.TransBatchInsert(tx, &po.PpmOrgUser{}, slice.ToSlice(userPoList))
		if batchInsert != nil {
			log.Error(batchInsert)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, batchInsert)
		}

		log.Info(json.ToJsonIgnoreError(userOutPoList))

		batchInsert = mysql.TransBatchInsert(tx, &po.PpmOrgUserOutInfo{}, slice.ToSlice(userOutPoList))
		if batchInsert != nil {
			log.Error(batchInsert)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, batchInsert)
		}

		batchInsert = mysql.TransBatchInsert(tx, &po.PpmOrgUserOrganization{}, slice.ToSlice(userOrgList))
		if batchInsert != nil {
			log.Error(batchInsert)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, batchInsert)
		}

		batchInsert = mysql.TransBatchInsert(tx, &po.PpmOrgUserConfig{}, slice.ToSlice(userConfigList))
		if batchInsert != nil {
			log.Error(batchInsert)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, batchInsert)
		}
	}

	if len(userDepList) > 0{
		userDepPoIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(userDepList))
		if idErr != nil{
			log.Error(idErr)
			return idErr
		}
		for i, _ := range userDepList{
			userDepList[i].Id = userDepPoIds.Ids[i].Id
		}
		batchInsert := mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, slice.ToSlice(userDepList))
		if batchInsert != nil {
			log.Error(batchInsert)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, batchInsert)
		}
	}

	return nil
}

func InitFsManager(orgId int64, userId int64, roleId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo{
	err2 := rolefacade.RoleUserRelationRelaxed(orgId, userId, roleId)
	if err2 != nil {
		log.Error(err2)
		return err2
	}
	err2 = OrgOwnerInit(orgId, userId, tx)
	if err2 != nil {
		log.Error(err2)
		return err2
	}
	return nil
}

func InitFsUser(orgId int64, tenantKey string, openUserId string, tx sqlbuilder.Tx) (*bo.UserInfoBo, errs.SystemErrorInfo){
	topDep, err1 := GetTopDepartmentInfo(orgId)
	if err1 != nil{
		log.Error(err1)
		return nil, err1
	}

	tenant, err1 := feishu.GetTenant(tenantKey)
	if err1 != nil{
		log.Error(err1)
		return nil, err1
	}
	userBatchGetResp, err := tenant.GetUserBatchGetV2(nil, []string{openUserId})
	if err != nil{
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, err)
	}
	if userBatchGetResp.Code != 0{
		log.Errorf("err %s", json.ToJsonIgnoreError(userBatchGetResp))
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, userBatchGetResp.Msg)
	}
	userInfos := userBatchGetResp.Data.Users
	if len(userInfos) == 0{
		log.Errorf("user不存在")
		return nil, errs.BuildSystemErrorInfo(errs.UserNotExist)
	}
	fsUser := userInfos[0]

	userNativeId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUser)
	if idErr != nil{
		log.Error(idErr)
		return nil, idErr
	}
	userOutInfoId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOutInfo)
	if idErr != nil{
		log.Error(idErr)
		return nil, idErr
	}
	userConfigId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserConfig)
	if idErr != nil{
		log.Error(idErr)
		return nil, idErr
	}
	userOrgId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOrganization)
	if idErr != nil{
		log.Error(idErr)
		return nil, idErr
	}
	userDepId, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserDepartment)
	if idErr != nil{
		log.Error(idErr)
		return nil, idErr
	}
	userPo := assemblyFeiShuUserInfo(orgId, fsUser)
	userPo.Id = userNativeId
	userOutPo := assemblyFeiShuUserOutInfo(orgId, userNativeId, fsUser)
	userOutPo.Id = userOutInfoId
	userConfigPo := assemblyFeiShuUserConfigInfo(orgId, userNativeId)
	userConfigPo.Id = userConfigId
	userOrgRelationPo := assemblyFeiShuUserOrgRelationInfo(orgId, userNativeId)
	userOrgRelationPo.Id = userOrgId
	userDepRelationPo := assemblyFeiShuUserDepRelationInfo(orgId, userNativeId, topDep.Id)
	userDepRelationPo.Id = userDepId

	dbErr := mysql.TransInsert(tx, &userPo)
	if dbErr != nil{
		log.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	dbErr = mysql.TransInsert(tx, &userOutPo)
	if dbErr != nil{
		log.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	dbErr = mysql.TransInsert(tx, &userConfigPo)
	if dbErr != nil{
		log.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	dbErr = mysql.TransInsert(tx, &userOrgRelationPo)
	if dbErr != nil{
		log.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	dbErr = mysql.TransInsert(tx, &userDepRelationPo)
	if dbErr != nil{
		log.Error(dbErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	userBo := &bo.UserInfoBo{}
	_ = copyer.Copy(userPo, userBo)
	return userBo, nil
}


func assemblyFeiShuUserInfo(orgId int64, fsUserDetailInfo vo.UserDetailInfoV2) po.PpmOrgUser{
	phoneNumber := fsUserDetailInfo.Mobile
	sourceChannel := consts.AppSourceChannelFeiShu
	sourcePlatform := consts.AppSourceChannelFeiShu
	name := fsUserDetailInfo.Name

	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userPo := &po.PpmOrgUser{
		OrgId:              orgId,
		Name:               name,
		NamePinyin:         pinyin.ConvertToPinyin(name),
		Avatar:             fsUserDetailInfo.Avatar.AvatarOrigin,
		LoginName:          phoneNumber, //
		LoginNameEditCount: 0,
		Email:              fsUserDetailInfo.Email,
		Mobile:             phoneNumber,
		Password:           pwd,
		PasswordSalt:       salt,
		SourceChannel:      sourceChannel,
		SourcePlatform:     sourcePlatform,
	}
	return *userPo
}

func assemblyFeiShuUserOutInfo(orgId int64, userId int64, fsUserDetailInfo vo.UserDetailInfoV2) po.PpmOrgUserOutInfo{
	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userOutInfo := &po.PpmOrgUserOutInfo{}
	userOutInfo.UserId = userId
	userOutInfo.OrgId = orgId
	userOutInfo.OutOrgUserId = fsUserDetailInfo.OpenId
	userOutInfo.OutUserId = fsUserDetailInfo.OpenId
	userOutInfo.IsDelete = consts.AppIsNoDelete
	userOutInfo.Status = consts.AppStatusEnable
	userOutInfo.SourceChannel = consts.AppSourceChannelFeiShu
	userOutInfo.Name = fsUserDetailInfo.Name
	userOutInfo.Avatar = fsUserDetailInfo.Avatar.AvatarOrigin
	userOutInfo.JobNumber = fsUserDetailInfo.EmployeeNo

	return *userOutInfo
}

func assemblyFeiShuUserConfigInfo(orgId int64, userId int64) po.PpmOrgUserConfig{
	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userConfigInfo := &po.PpmOrgUserConfig{}
	userConfigInfo.UserId = userId
	userConfigInfo.OrgId = orgId

	return *userConfigInfo
}

func assemblyFeiShuUserOrgRelationInfo(orgId int64, userId int64) po.PpmOrgUserOrganization{
	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userOrgRelationInfo := &po.PpmOrgUserOrganization{}
	userOrgRelationInfo.UserId = userId
	userOrgRelationInfo.OrgId = orgId
	userOrgRelationInfo.Status = consts.AppStatusEnable
	userOrgRelationInfo.UseStatus = consts.AppStatusDisabled
	userOrgRelationInfo.CheckStatus = consts.AppCheckStatusSuccess

	return *userOrgRelationInfo
}

func assemblyFeiShuUserDepRelationInfo(orgId int64, userId int64, depId int64) po.PpmOrgUserDepartment{
	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userDepRelationInfo := &po.PpmOrgUserDepartment{}
	userDepRelationInfo.UserId = userId
	userDepRelationInfo.OrgId = orgId
	userDepRelationInfo.DepartmentId = depId

	return *userDepRelationInfo
}