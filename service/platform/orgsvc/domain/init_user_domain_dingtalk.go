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
	"github.com/galaxy-book/polaris-backend/common/extra/dingtalk"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InitDingTalkUserList(orgId int64, tenantKey string, fsDepIdMap map[int64]int64, superAdminRoleId, normalAdminRoleId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo{
	userIdCache := map[string]bool{}
	userPoList := make([]po.PpmOrgUser, 0)
	userOutPoList := make([]po.PpmOrgUserOutInfo, 0)
	userConfigList := make([]po.PpmOrgUserConfig, 0)
	userOrgList := make([]po.PpmOrgUserOrganization, 0)
	userDepList := make([]po.PpmOrgUserDepartment, 0)

	//根部门id
	rootDepId := fsDepIdMap[1]

	users, err := dingtalk.GetScopeUsers(tenantKey)
	if err != nil{
		log.Error(err)
		return err
	}

	surplusSize := len(users)
	if surplusSize == 0{
		return nil
	}

	userPoIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUser, surplusSize)
	if idErr != nil{
		log.Error(idErr)
		return idErr
	}

	hasSuperAdmin := false

	for i, fsUser := range users{
		userNativeId := userPoIds.Ids[i].Id
		if _, ok := userIdCache[fsUser.UnionId]; ! ok{
			userPo := assemblyDingTalkUserInfo(orgId, fsUser)
			userPo.Id = userNativeId
			userPoList = append(userPoList, userPo)
			userOutPoList = append(userOutPoList, assemblyDingTalkUserOutInfo(orgId, userNativeId, fsUser))
			userConfigList = append(userConfigList, assemblyFeiShuUserConfigInfo(orgId, userNativeId))
			userOrgList = append(userOrgList, assemblyFeiShuUserOrgRelationInfo(orgId, userNativeId))

			if fsUser.IsAdmin{
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
		userDeps := fsUser.Department
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
		userIdCache[fsUser.UnionId] = true
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

func InitDingTalkUser(orgId int64, corpId string, unionId string, tx sqlbuilder.Tx) (*bo.UserInfoBo, errs.SystemErrorInfo){
	topDep, err1 := GetTopDepartmentInfo(orgId)
	if err1 != nil{
		log.Error(err1)
		return nil, err1
	}

	client, err2 := dingtalk.GetDingTalkClientRest(corpId)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.DingTalkClientError, err2)
	}

	userIdResp, err2 := client.GetUserIdByUnionId(unionId)
	if err2 != nil{
		log.Error(err2)
		return nil, errs.DingTalkOpenApiCallError
	}
	if userIdResp.ErrCode != 0{
		log.Error(userIdResp.ErrMsg)
		return nil, errs.DingTalkOpenApiCallError
	}

	userId := userIdResp.UserId

	userDetailResp, err2 := client.GetUserDetail(userId, nil)
	if err2 != nil{
		log.Error(err2)
		return nil, errs.DingTalkOpenApiCallError
	}
	if userDetailResp.ErrCode != 0{
		log.Error(userDetailResp.ErrMsg)
		return nil, errs.DingTalkOpenApiCallError
	}

	user := userDetailResp.UserList

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
	userPo := assemblyDingTalkUserInfo(orgId, user)
	userPo.Id = userNativeId
	userOutPo := assemblyDingTalkUserOutInfo(orgId, userNativeId, user)
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

func assemblyDingTalkUserInfo(orgId int64, duser sdk.UserList) po.PpmOrgUser{
	phoneNumber := ""
	sourceChannel := consts.AppSourceChannelDingTalk
	sourcePlatform := consts.AppSourceChannelDingTalk
	name := duser.Name

	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userPo := &po.PpmOrgUser{
		OrgId:              orgId,
		Name:               name,
		NamePinyin:         pinyin.ConvertToPinyin(name),
		Avatar:             duser.Avatar,
		LoginName:          phoneNumber, //
		LoginNameEditCount: 0,
		Email:              "",
		Mobile:             phoneNumber,
		Password:           pwd,
		PasswordSalt:       salt,
		SourceChannel:      sourceChannel,
		SourcePlatform:     sourcePlatform,
	}
	return *userPo
}

func assemblyDingTalkUserOutInfo(orgId int64, userId int64, duser sdk.UserList) po.PpmOrgUserOutInfo{
	pwd := uuid.NewUuid()
	salt := uuid.NewUuid()
	pwd = md5.Md5V(salt + pwd)
	userOutInfo := &po.PpmOrgUserOutInfo{}
	userOutInfo.UserId = userId
	userOutInfo.OrgId = orgId
	userOutInfo.OutOrgUserId = duser.UnionId
	userOutInfo.OutUserId = duser.UnionId
	userOutInfo.IsDelete = consts.AppIsNoDelete
	userOutInfo.Status = consts.AppStatusEnable
	userOutInfo.SourceChannel = consts.AppSourceChannelDingTalk
	userOutInfo.Name = duser.Name
	userOutInfo.Avatar = duser.Avatar
	userOutInfo.JobNumber = duser.JobNumber

	return *userOutInfo
}