package domain

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	dingtalk2 "github.com/galaxy-book/polaris-backend/common/extra/dingtalk"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"strconv"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//通用版
func GeneralInitOrg(initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {

	isDingTalk := initOrgBo.SourceChannel == consts.AppSourceChannelDingTalk

	//判断这个组织是否已经存在 内含存在的时候更新
	_, isUpdateFail, err := GetGeneralOrgInfoByOutOrgId(initOrgBo.OutOrgId, initOrgBo.PermanentCode, initOrgBo.SourceChannel, isDingTalk, tx)

	//不需要初始化并且已经更新过了
	if err == nil {
		log.Errorf("组织已经存在，不需要初始化，初始化信息为：%s", json.ToJsonIgnoreError(initOrgBo))
		return 0, errs.BuildSystemErrorInfo(errs.OrgNotNeedInitError)
	}
	//走这里说明组织存在 且更新时候异常 直接返回出去
	if isUpdateFail != nil && *isUpdateFail {
		logger.GetDefaultLogger().Error("组织初始化，更新组织时异常" + strs.ObjectToString(err))
		return 0, err
	}

	//组织信息初始化
	orgId, err := GeneralOrgInfoInit(initOrgBo, isDingTalk, tx)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	//处理外部信息
	_, err = GeneralOrgOutInfoInit(initOrgBo, orgId, isDingTalk, tx)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	_, err = GeneralOrgConfigInfoInit(orgId, tx)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	//权限、角色初始化
	roleInitResp := rolefacade.RoleInit(rolevo.RoleInitReqVo{
		OrgId: orgId,
	})
	if roleInitResp.Failure() {
		log.Error(roleInitResp.Message)
		return 0, roleInitResp.Error()
	}
	log.Info("权限、角色初始化成功")

	//优先级，任务类型，任务来源初始化
	priorityInfo := projectfacade.ProjectInit(projectvo.ProjectInitReqVo{OrgId: orgId})
	if priorityInfo.Failure() {
		log.Error(priorityInfo.Message)
		return 0, priorityInfo.Error()
	}
	log.Info("优先级初始化成功")

	err = InitDepartment(
		orgId,
		initOrgBo.OutOrgId,
		initOrgBo.SourceChannel,
		roleInitResp.RoleInitResp.OrgSuperAdminRoleId,
		roleInitResp.RoleInitResp.OrgNormalAdminRoleId,
		tx)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	//钉钉单独初始化的内容
	dErr := DingDingInit(initOrgBo.OutOrgId, orgId, roleInitResp.RoleInitResp, tx)

	if dErr != nil {
		return 0, errs.BuildSystemErrorInfo(errs.OrgInitError, dErr)
	}

	return orgId, nil
}

func GetGeneralOrgInfoByOutOrgId(outOrgId string, permanentCode string, sourceChannel string, isDingTalk bool, tx sqlbuilder.Tx) (baseOrgInfo *bo.BaseOrgInfoBo, updateFlag *bool, returnErr errs.SystemErrorInfo) {
	outOrgInfo := &po.PpmOrgOrganizationOutInfo{}

	conds := db.Cond{
		consts.TcOutOrgId:      outOrgId,
		consts.TcSourceChannel: sourceChannel,
	}

	//飞书的拼接未删除的条件
	if !isDingTalk {
		conds[consts.TcIsDelete] = consts.AppIsNoDelete
	}

	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, conds, outOrgInfo)
	//err 不为空说明组织不存在 查不到
	if err != nil {
		log.Error(err)
		return nil, nil, errs.BuildSystemErrorInfo(errs.OrgOutInfoNotExist)
	}

	//这里需要先判断 如果是钉钉的已经存在的情况下需要更新一下,然后后面的不进行初始化操作了
	if isDingTalk {
		isUpdateFail := true
		//钉钉的单独重置对应的数据
		orgAuthTicketInfo := &bo.OrgAuthTicketInfoBo{}

		//组织已存在
		logger.GetDefaultLogger().Infof("钉钉组织 %s 已经存在, 删除状态 %d", outOrgId, outOrgInfo.IsDelete)

		if outOrgInfo.IsDelete == consts.AppIsDeleted {
			outOrgInfo.IsDelete = consts.AppIsNoDelete

		}
		if outOrgInfo.AuthTicket != "" {
			json.FromJson(outOrgInfo.AuthTicket, orgAuthTicketInfo)
		}
		orgAuthTicketInfo.PermanentCode = permanentCode
		newAuthTicket, _ := json.ToJson(orgAuthTicketInfo)

		outOrgInfo.AuthTicket = newAuthTicket

		AssemblyOrgOutInfo(outOrgId, outOrgInfo)

		err = mysql.TransUpdate(tx, outOrgInfo)
		//这边翻过来 update 失败就不处理了
		if err != nil {
			logger.GetDefaultLogger().Error("组织初始化，更新组织时异常" + strs.ObjectToString(err))
			return nil, &isUpdateFail, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}

	//获取原本的组织信息 不存在就出去初始化
	orgInfo, err1 := GetOrgBoById(outOrgInfo.OrgId)
	if err1 != nil {
		log.Error(err1)
		return nil, nil, errs.BuildSystemErrorInfo(errs.OrgNotExist)
	}

	//组装数据返回
	return &bo.BaseOrgInfoBo{
		OrgId:         orgInfo.Id,
		OrgName:       orgInfo.Name,
		OutOrgId:      outOrgId,
		SourceChannel: sourceChannel,
	}, nil, nil
}

//组织信息初始化
func GeneralOrgInfoInit(initOrgBo bo.InitOrgBo, isDingTalk bool, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {

	isAuth := 0
	//如果是dingtalk的,部分信息来源于钉钉openApi
	if isDingTalk {

		authInfo, err := GetCorpAuthInfo(initOrgBo.OutOrgId)
		if err != nil {
			return 0, err
		}
		//赋值
		initOrgBo.OrgName = authInfo.AuthCorpInfo.CorpName
		initOrgBo.OrgLogo = authInfo.AuthCorpInfo.CorpLogoUrl
		initOrgBo.OrgProvince = authInfo.AuthCorpInfo.CorpProvince
		initOrgBo.OrgCity = authInfo.AuthCorpInfo.CorpCity
		//赋值给initOrgBo 西面判断的时候用
		initOrgBo.IsAuthenticated = authInfo.AuthCorpInfo.IsAuthenticated

	}

	if initOrgBo.IsAuthenticated {
		isAuth = 1
	}

	orgId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrganization)
	if err != nil {
		return 0, err
	}
	org := &po.PpmOrgOrganization{}
	//统一处理的id,状态,名字,sourceChannel
	org.Id = orgId
	org.Status = consts.AppStatusEnable
	org.IsDelete = consts.AppIsNoDelete
	org.SourceChannel = initOrgBo.SourceChannel

	org.Name = initOrgBo.OrgName
	org.LogoUrl = initOrgBo.OrgLogo
	org.Address = initOrgBo.OrgProvince + initOrgBo.OrgCity
	org.IsAuthenticated = isAuth
	//插入org 信息
	err2 := mysql.TransInsert(tx, org)
	if err2 != nil {
		logger.GetDefaultLogger().Error("组织初始化，添加组织时异常:" + strs.ObjectToString(err2))
		return 0, err
	}
	return orgId, nil
}

//组织外部信息初始化
func GeneralOrgOutInfoInit(initOrgBo bo.InitOrgBo, orgId int64, isDingTalk bool, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	isAuth := 0
	if initOrgBo.IsAuthenticated {
		isAuth = 1
	}

	orgOutInfoId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrganizationOutInfo)
	if err != nil {
		return 0, err
	}
	orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
	orgOutInfo.Id = orgOutInfoId
	orgOutInfo.OrgId = orgId
	orgOutInfo.IsDelete = consts.AppIsNoDelete
	orgOutInfo.Status = consts.AppStatusEnable
	orgOutInfo.SourceChannel = initOrgBo.SourceChannel

	if isDingTalk {
		authInfo, err := GetCorpAuthInfo(initOrgBo.OutOrgId)
		if err != nil {
			return 0, err
		}
		//赋值
		initOrgBo.OrgName = authInfo.AuthCorpInfo.CorpName
		initOrgBo.OrgLogo = authInfo.AuthCorpInfo.CorpLogoUrl
		initOrgBo.OrgProvince = authInfo.AuthCorpInfo.CorpProvince
		initOrgBo.OrgCity = authInfo.AuthCorpInfo.CorpCity
		//赋值给initOrgBo 西面判断的时候用
		initOrgBo.IsAuthenticated = authInfo.AuthCorpInfo.IsAuthenticated
	}

	orgOutInfo.Name = initOrgBo.OrgName
	orgOutInfo.OutOrgId = initOrgBo.OutOrgId
	orgOutInfo.Industry = initOrgBo.Industry
	orgOutInfo.IsAuthenticated = isAuth
	orgOutInfo.AuthLevel = strconv.Itoa(initOrgBo.AuthLevel)

	err2 := mysql.TransInsert(tx, orgOutInfo)
	if err2 != nil {
		logger.GetDefaultLogger().Error("组织初始化，添加外部组织信息时异常: " + strs.ObjectToString(err2))
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}
	return orgOutInfoId, nil
}

//组织配置初始化
func GeneralOrgConfigInfoInit(orgId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	sysConfig := &po.PpmOrcConfig{}

	payLevel := &po.PpmBasPayLevel{}
	err := mysql.SelectById(payLevel.TableName(), 1, payLevel)
	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	orgConfigId, err1 := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrgConfig)
	if err1 != nil {
		log.Error(err1)
		return 0, err1
	}

	sysConfig.Id = orgConfigId
	sysConfig.OrgId = orgId
	sysConfig.TimeZone = "Asia/Shanghai"
	sysConfig.TimeDifference = "+08:00"
	sysConfig.PayLevel = 1
	sysConfig.PayStartTime = time.Now()
	sysConfig.PayEndTime = time.Now().Add(time.Duration(payLevel.Duration) * time.Second)
	sysConfig.Language = "zh-CN"
	sysConfig.RemindSendTime = "09:00:00"
	sysConfig.DatetimeFormat = "yyyy-MM-dd HH:mm:ss"
	sysConfig.PasswordLength = 6
	sysConfig.PasswordRule = 1
	sysConfig.MaxLoginFailCount = 0
	sysConfig.Status = consts.AppStatusEnable
	err2 := mysql.TransInsert(tx, sysConfig)
	if err2 != nil {
		logger.GetDefaultLogger().Error("组织初始化，添加组织配置信息时异常: " + strs.ObjectToString(err2))
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}

	return orgConfigId, nil
}

//暂时钉钉专有的初始化流程
func DingDingInit(corpId string, orgId int64, roleInitResp *bo.RoleInitResp, tx sqlbuilder.Tx) error {

	//获取授权企业钉钉用户角色 单独处理
	dingtalkUserRoleBos, err3 := dingtalk2.GetDingTalkUserRoleBos(corpId)
	if err3 != nil {
		log.Error("Rollback" + strs.ObjectToString(err3))
	}

	teamId, initErr := TeamInit(orgId, tx)
	if initErr != nil {
		log.Error(initErr)
		return initErr
	}

	//用户角色init
	rootUserId, err := InitUserRoles(dingtalkUserRoleBos, orgId, teamId, corpId, roleInitResp, tx)
	if err != nil {
		return err
	}

	//拥有这初始化

	err1 := ownerinit(rootUserId, orgId, teamId, tx)

	if err1 != nil {
		return err1
	}

	return nil

}

func InitUserRoles(dingtalkUserRoleBos []*bo.DingTalkUserRoleBo, orgId, teamId int64, corpId string,
	roleInitResp *bo.RoleInitResp, tx sqlbuilder.Tx) (int64, error) {

	var rootId = ""
	var rootUserId int64 = -1

	for _, usr := range dingtalkUserRoleBos {
		isRoot := false
		isAdmin := usr.IsAdmin
		if usr.IsRoot && rootId == "" {
			rootId = usr.UserId
			isRoot = true
		}

		//初始化用户
		userId, err := UserInitByOrg(usr.UserId, corpId, orgId, tx)

		if err != nil {
			logger.GetDefaultLogger().Error("Rollback" + err.Message())
			return 0, err
		}

		if isRoot {
			rootUserId = userId
		}

		err1 := TeamUserInit(orgId, teamId, userId, isRoot, tx)
		if err1 != nil {
			logger.GetDefaultLogger().Error("Rollback" + err1.Message())
			return 0, err1
		}

		//角色和人绑定
		//if isRoot {
		//	err2 := rolefacade.RoleUserRelationRelaxed(orgId, userId, roleInitResp.OrgSuperAdminRoleId)
		//	if err2 != nil {
		//		log.Error(err2)
		//		return 0, err2
		//	}
		//} else if isAdmin {
		//	err2 := rolefacade.RoleUserRelationRelaxed(orgId, userId, roleInitResp.OrgNormalAdminRoleId)
		//	if err2 != nil {
		//		log.Error(err2)
		//		return 0, err2
		//	}
		//}
		if err2 := bindingUserAndRole(isRoot, isAdmin, orgId, userId, roleInitResp); err2 != nil {
			return 0, err2
		}

	}
	return rootUserId, nil
}

func bindingUserAndRole(isRoot, isAdmin bool, orgId, userId int64, roleInitResp *bo.RoleInitResp) errs.SystemErrorInfo {
	if isRoot {
		err2 := rolefacade.RoleUserRelationRelaxed(orgId, userId, roleInitResp.OrgSuperAdminRoleId)
		if err2 != nil {
			log.Error(err2)
			return err2
		}
	} else if isAdmin {
		err2 := rolefacade.RoleUserRelationRelaxed(orgId, userId, roleInitResp.OrgNormalAdminRoleId)
		if err2 != nil {
			log.Error(err2)
			return err2
		}
	}
	return nil
}

func ownerinit(rootUserId, orgId, teamId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	if rootUserId != -1 {
		//设置组织的owner
		//err = orgfacade.OrgOwnerInit(orgId, rootUserId, tx)

		err := OrgOwnerInit(orgId, rootUserId, tx)

		if err != nil {
			logger.GetDefaultLogger().Error("Rollback " + err.Message())
			return err
		}

		//设置团队的owner
		//err = orgfacade.TeamOwnerInit(teamId, rootUserId, tx)
		err1 := TeamOwnerInit(teamId, rootUserId, tx)

		if err1 != nil {
			logger.GetDefaultLogger().Error("Rollback " + err1.Message())
			return err
		}
	}
	return nil
}
