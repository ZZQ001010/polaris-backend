package domain

import (
	"fmt"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/date"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/md5"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/pinyin"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

func UserInit(orgId int64, corpId string, outUserId string, sourceChannel string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	log.Infof("用户初始化操作, orgId: %d, corpId %s, outUserId %s", orgId, corpId, outUserId)
	baseUserInfo, err := GetBaseUserInfoByEmpId(sourceChannel, orgId, outUserId)
	if err != nil {
		//这里做用户初始化的兜底
		lockKey := consts.InitUserLock + sourceChannel + outUserId
		suc, err := cache.TryGetDistributedLock(lockKey, outUserId)
		log.Infof("准备获取分布式锁 %v", suc)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		if suc {
			log.Infof("获取分布式锁成功 %v", suc)
			defer func() {
				if _, err := cache.ReleaseDistributedLock(lockKey, outUserId); err != nil{
					log.Error(err)
				}
			}()

			var err3 errs.SystemErrorInfo = nil
			baseUserInfo, err3 = GetBaseUserInfoByEmpId(sourceChannel, orgId, outUserId)
			if err3 != nil {
				log.Error(err3)
			}
			if baseUserInfo != nil {
				return baseUserInfo, nil
			}

			if sourceChannel == consts.AppSourceChannelDingTalk {
				//钉钉走自己的逻辑
				err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
					_, err3 = InitDingTalkUser(orgId, corpId, outUserId, tx)
					if err3 != nil {
						log.Error(err3)
						return err3
					}
					return nil
				})
				if err1 != nil {
					log.Error(err1)
					return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
				}
			} else if sourceChannel == consts.AppSourceChannelFeiShu {
				//飞书走飞书逻辑
				err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
					_, err := InitFsUser(orgId, corpId, outUserId, tx)
					if err != nil {
						log.Error(err)
						return err
					}
					return nil
				})
				if err1 != nil {
					log.Error(err1)
					return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
				}
			} else {
				return nil, errs.BuildSystemErrorInfo(errs.SourceChannelNotDefinedError)
			}

			baseUserInfo, err3 = GetBaseUserInfoByEmpId(sourceChannel, orgId, outUserId)
			if err3 != nil {
				log.Error(err3)
				return nil, err3
			}
			return baseUserInfo, nil
		} else {
			baseUserInfo, err = GetBaseUserInfoByEmpId(sourceChannel, orgId, outUserId)
			if err != nil {
				log.Error(err)
				return nil, errs.UserInitError
			}
		}
	}
	return baseUserInfo, nil
}

func GetUserBo(userId int64) (*bo.UserInfoBo, bool, errs.SystemErrorInfo) {
	user := &po.PpmOrgUser{}
	err1 := mysql.SelectOneByCond(user.TableName(), db.Cond{
		consts.TcId:       userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, user)
	if err1 != nil {
		logger.GetDefaultLogger().Error(strs.ObjectToString(err1))
		return nil, false, errs.BuildSystemErrorInfo(errs.UserNotFoundError, err1)
	}

	userInfo := &bo.UserInfoBo{}
	err1 = copyer.Copy(user, userInfo)
	if err1 != nil {
		return nil, false, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	passwordSet := false
	if user.Password != consts.BlankString {
		passwordSet = true
	}
	return userInfo, passwordSet, nil
}

func GetUserInfo(orgId int64, userId int64, sourceChannel string) (*bo.UserInfoBo, bool, errs.SystemErrorInfo) {
	baseUserInfo, err := GetBaseUserInfo(sourceChannel, orgId, userId)
	if err != nil {
		log.Error(err)
		return nil, false, errs.BuildSystemErrorInfo(errs.UserNotInitError)
	}

	baseOrgInfo, err := GetBaseOrgInfo(sourceChannel, orgId)
	if err != nil {
		log.Error(err)
		return nil, false, errs.BuildSystemErrorInfo(errs.OrgNotInitError)
	}

	ownerInfo, passwordSet, err := GetUserBo(userId)
	if err != nil {
		log.Error(err)
		return nil, false, err
	}
	//部分属性覆盖
	ownerInfo.EmplID = &baseUserInfo.OutUserId
	ownerInfo.OrgID = orgId
	ownerInfo.OrgName = baseOrgInfo.OrgName

	//这里先默认写死
	ownerInfo.Rimanente = 10
	ownerInfo.Level = 1
	ownerInfo.LevelName = "试用级别"

	return ownerInfo, passwordSet, nil
}

func GetOutUserInfoListBySourceChannel(sourceChannel string, page int, size int) ([]bo.UserOutInfoBo, errs.SystemErrorInfo) {
	userOutInfos := &[]po.PpmOrgUserOutInfo{}
	_, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableUserOutInfo, db.Cond{
		consts.TcSourceChannel: sourceChannel,
		consts.TcStatus:        consts.AppStatusEnable,
		consts.TcIsDelete:      consts.AppIsNoDelete,
	}, nil, page, size, "org_id asc", userOutInfos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	userOutInfoBos := &[]bo.UserOutInfoBo{}
	err = copyer.Copy(userOutInfos, userOutInfoBos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *userOutInfoBos, nil
}

func GetOutUserInfo(outUserId, sourceChannel string) (*bo.UserOutInfoBo, errs.SystemErrorInfo) {
	userOutInfo := &po.PpmOrgUserOutInfo{}
	err := mysql.SelectOneByCond(userOutInfo.TableName(), db.Cond{
		consts.TcOutUserId:     outUserId,
		consts.TcSourceChannel: sourceChannel,
		consts.TcIsDelete:      consts.AppIsNoDelete}, userOutInfo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.UserOutInfoNotExist)
	}
	userOutInfoBo := &bo.UserOutInfoBo{}
	_ = copyer.Copy(userOutInfo, userOutInfoBo)
	return userOutInfoBo, nil
}

func GetUserInfoListByOrg(orgId int64) ([]bo.SimpleUserInfoBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	userInfos := &[]po.PpmOrgUser{}

	userAlias := "user."
	orgAlias := "org."
	selectErr := conn.Select("user.name", "user.id").
		From(consts.TableUser+" user", consts.TableUserOrganization+" org").
		Where(db.Cond{
			userAlias + consts.TcId:       db.Raw(orgAlias + consts.TcUserId),
			orgAlias + consts.TcIsDelete:  consts.AppIsNoDelete,
			userAlias + consts.TcIsDelete: consts.AppIsNoDelete,
			orgAlias + consts.TcStatus:    consts.AppIsInitStatus,
			orgAlias + consts.TcOrgId:     orgId,
		}).
		All(userInfos)
	if selectErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
	}
	userInfoBos := &[]bo.SimpleUserInfoBo{}
	copyErr := copyer.Copy(userInfos, userInfoBos)

	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *userInfoBos, nil
}

//如果err不等于空，说明用户未注册
func GetUserInfoByMobile(phoneNumber string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	userPo := &po.PpmOrgUser{}
	err := mysql.SelectOneByCond(consts.TableUser, db.Cond{
		consts.TcMobile:   phoneNumber,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, userPo)
	if err != nil {
		log.Error(err)
		return nil, errs.UserNotExist
	}
	userBo := &bo.UserInfoBo{}
	copyErr := copyer.Copy(userPo, userBo)
	if copyErr != nil {
		log.Error(copyErr)
	}
	return userBo, nil
}

func GetUserInfoByEmail(email string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	userPo := &po.PpmOrgUser{}
	err := mysql.SelectOneByCond(consts.TableUser, db.Cond{
		consts.TcEmail:   email,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, userPo)
	if err != nil {
		log.Error(err)
		return nil, errs.UserNotExist
	}
	userBo := &bo.UserInfoBo{}
	copyErr := copyer.Copy(userPo, userBo)
	if copyErr != nil {
		log.Error(copyErr)
	}
	return userBo, nil
}

//如果err不等于空，说明用户未注册
func GetUserInfoByLoginNameAndPwd(loginName string, pwd string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	userPo := &po.PpmOrgUser{}
	err := mysql.SelectOneByCond(consts.TableUser, db.Cond{
		consts.TcLoginName: loginName,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}, userPo)
	if err != nil {
		log.Error(err)
		return nil, errs.UserNotExist
	}

	salt := userPo.PasswordSalt
	pwd = md5.Md5V(salt + pwd)

	if userPo.Password != pwd {
		return nil, errs.BuildSystemErrorInfo(errs.PwdLoginUsrOrPwdNotMatch)
	}
	userBo := &bo.UserInfoBo{}
	copyErr := copyer.Copy(userPo, userBo)
	if copyErr != nil {
		log.Error(copyErr)
	}
	return userBo, nil
}

//loginName为允许为账号，邮箱，手机号
func GetUserInfoByPwd(loginName string, pwd string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return nil, errs.MysqlOperateError
	}

	userPo := &po.PpmOrgUser{}
	err = conn.Collection(consts.TableUser).Find(db.And(db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}, db.Or(db.Cond{
		consts.TcMobile: loginName,
	},db.Cond{
		consts.TcEmail: loginName,
	}))).One(userPo)
	if err != nil {
		if err == db.ErrNoMoreRows{
			log.Error(err)
			return nil, errs.UserNotExist
		}else{
			return nil, errs.MysqlOperateError
		}
	}
	salt := userPo.PasswordSalt
	pwd = util.PwdEncrypt(pwd, salt)

	if userPo.Password != pwd {
		return nil, errs.BuildSystemErrorInfo(errs.PwdLoginUsrOrPwdNotMatch)
	}
	userBo := &bo.UserInfoBo{}
	copyErr := copyer.Copy(userPo, userBo)
	if copyErr != nil {
		log.Error(copyErr)
	}
	return userBo, nil
}

//loginName为允许为账号，邮箱，手机号
func GetUserInfoByLoginName(loginName string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return nil, errs.MysqlOperateError
	}

	userPo := &po.PpmOrgUser{}
	err = conn.Collection(consts.TableUser).Find(db.And(db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}, db.Or(db.Cond{
		consts.TcMobile: loginName,
	},db.Cond{
		consts.TcEmail: loginName,
	}))).One(userPo)
	if err != nil {
		if err == db.ErrNoMoreRows{
			return nil, errs.UserNotExist
		}else{
			log.Error(err)
			return nil, errs.MysqlOperateError
		}
	}
	userBo := &bo.UserInfoBo{}
	copyErr := copyer.Copy(userPo, userBo)
	if copyErr != nil {
		log.Error(copyErr)
	}
	return userBo, nil
}

func UserRegister(regInfo bo.UserSMSRegisterInfo) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	loginName := regInfo.PhoneNumber
	if loginName == ""{
		loginName = regInfo.Email
	}
	if loginName == ""{
		return nil, errs.UserRegisterError
	}

	//注册时对手机号加锁
	lockKey := consts.UserBindLoginNameLock + loginName
	uid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uid)
	if lockErr != nil{
		log.Error(lockErr)
		return nil, errs.UserRegisterError
	}
	if suc{
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil{
				log.Error(err)
			}
		}()
	}else{
		log.Error("注册失败")
		return nil, errs.UserRegisterError
	}

	userPo, err := assemblyUserRegisterUserInfo(regInfo)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//插入用户
	err1 := mysql.Insert(userPo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.UserRegisterError)
	}

	//即时注册时插入失败，查看时也会做二次check并插入
	err = insertUserConfig(userPo.OrgId, userPo.Id)
	if err != nil {
		log.Error(err)
	}

	userBo := &bo.UserInfoBo{}
	copyErr := copyer.Copy(userPo, userBo)
	if copyErr != nil {
		log.Error(copyErr)
	}
	return userBo, nil
}

func assemblyUserRegisterUserInfo(regInfo bo.UserSMSRegisterInfo) (*po.PpmOrgUser, errs.SystemErrorInfo) {
	phoneNumber := regInfo.PhoneNumber
	email := regInfo.Email
	sourceChannel := regInfo.SourceChannel
	sourcePlatform := regInfo.SourcePlatform
	name := regInfo.Name

	userId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableUser)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError)
	}
	userPo := &po.PpmOrgUser{
		Id:                 userId,
		OrgId:              0,
		Name:               name,
		NamePinyin:         pinyin.ConvertToPinyin(name),
		Avatar:             "",
		LoginName:          phoneNumber, //
		LoginNameEditCount: 0,
		Email:              email,
		Mobile:             phoneNumber,
		SourceChannel:      sourceChannel,
		SourcePlatform:     sourcePlatform,
		//SourceObjId:,
		Creator: userId,
		Updator: userId,
	}
	return userPo, nil
}

//增加组织成员，添加关联以及加入顶级部门
//inCheck：是否需要被审核
func AddOrgMember(orgId, userId int64, operatorId int64, inCheck bool, inDisabled bool) errs.SystemErrorInfo {
	userOrgRelation, err := GetUserOrganizationNewestRelation(orgId, userId)
	//关联不存在或者已删除，或者审核未通过，此时允许新增关联
	if (err != nil && err.Code() == errs.UserOrgNotRelation.Code()) || userOrgRelation.IsDelete == consts.AppIsDeleted || userOrgRelation.CheckStatus == consts.AppCheckStatusFail {
		log.Errorf("用户%d和组织%d需要做关联", userId, orgId)
		//上锁
		lockKey := fmt.Sprintf("%s%d:%d", consts.UserAndOrgRelationLockKey, orgId, userId)
		lockUuid := uuid.NewUuid()

		suc, lockErr := cache.TryGetDistributedLock(lockKey, lockUuid)
		if lockErr != nil {
			log.Error(lockErr)
			return errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
		}
		if suc {
			defer func() {
				if _, err := cache.ReleaseDistributedLock(lockKey, lockUuid); err != nil {
					log.Error(err)
				}
			}()
			//二次check
			userOrgRelation, err := GetUserOrganizationNewestRelation(orgId, userId)
			if (err != nil && err.Code() == errs.UserOrgNotRelation.Code()) || userOrgRelation.IsDelete == consts.AppIsDeleted || userOrgRelation.CheckStatus == consts.AppCheckStatusFail {
				//组织用户做关联
				log.Errorf("用户%d和组织%d开始关联", userId, orgId)
				err = AddUserOrgRelation(orgId, userId, false, inCheck, inDisabled)
				//判断关联是否失败
				if err != nil {
					log.Error(err)
					return err
				}
				if !inCheck {
					//获取顶级部门
					topDep, err := GetTopDepartmentInfo(orgId)
					if err != nil {
						log.Error(err)
						return err
					}
					log.Infof("获取顶级部门成功 %s", json.ToJsonIgnoreError(topDep))

					//将用户加到顶级部门中
					err = BoundDepartmentUser(orgId, []int64{userId}, topDep.Id, operatorId, false)
					if err != nil {
						log.Error(err)
					}
					log.Infof("用户%d和部门%d关联成功", userId, topDep.Id)
				}
			}
		}
	}
	clearErr := ClearBaseUserInfo(orgId, userId)
	if err != nil {
		log.Error(clearErr)
		return clearErr
	}
	return nil
}

//添加用户和组织关联
func AddUserOrgRelation(orgId, userId int64, inUsed bool, inCheck bool, inDisabled bool) errs.SystemErrorInfo {
	userOrgId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOrganization)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.ApplyIdError)
	}

	useStatus := consts.AppStatusDisabled
	if inUsed {
		useStatus = consts.AppStatusEnable
	}

	checkStatus := consts.AppCheckStatusSuccess
	status := consts.AppStatusEnable
	if inCheck {
		checkStatus = consts.AppCheckStatusWait
		status = consts.AppStatusDisabled
	}

	if inDisabled{
		status = consts.AppStatusDisabled
	}

	userOrgPo := po.PpmOrgUserOrganization{
		Id:          userOrgId,
		OrgId:       orgId,
		UserId:      userId,
		CheckStatus: checkStatus,
		UseStatus:   useStatus,
		Status:      status,
		Creator:     userId,
		Updator:     userId,
	}

	err1 := mysql.Insert(&userOrgPo)
	if err != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func UpdateUserInfo(userId int64, upd mysql.Upd) errs.SystemErrorInfo {
	err := mysql.UpdateSmart(consts.TableUser, userId, upd)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func UpdateUserDefaultOrg(userId, orgId int64) errs.SystemErrorInfo {
	updateUserInfoErr := UpdateUserInfo(userId, mysql.Upd{
		consts.TcOrgId: orgId,
	})
	if updateUserInfoErr != nil {
		log.Error(updateUserInfoErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

//用户登录时回调
func UserLoginHook(userId int64, orgId int64) errs.SystemErrorInfo {
	//更新用户最后登录时间
	err := UpdateUserInfo(userId, mysql.Upd{
		consts.TcLastLoginTime: date.FormatTime(types.NowTime()),
	})
	if err != nil {
		log.Error(err)
	}
	if orgId != 0 {
		//更新使用状态
		_, err1 := mysql.UpdateSmartWithCond(consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcUserId:   userId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcUseStatus: consts.AppStatusEnable,
		})
		if err1 != nil {
			log.Error(err1)
			err = errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
	}
	return err
}

func BatchGetUserDetailInfo(userIds []int64) ([]bo.UserInfoBo, errs.SystemErrorInfo) {
	userIds = slice.SliceUniqueInt64(userIds)
	po := &[]po.PpmOrgUser{}
	err := mysql.SelectAllByCond(consts.TableUser, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		//consts.TcOrgId:    orgId,
		consts.TcId: db.In(userIds),
	}, po)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	bo := &[]bo.UserInfoBo{}
	copyErr := copyer.Copy(po, bo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return *bo, nil
}
