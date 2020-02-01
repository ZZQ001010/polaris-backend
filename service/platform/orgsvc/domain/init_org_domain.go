package domain

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"strconv"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
)

//留着做对比方便的注释
func InitOrg(initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	_, err := GetOrgInfoByOutOrgId(initOrgBo.OutOrgId, initOrgBo.SourceChannel)
	if err == nil {
		log.Errorf("组织已经存在，不需要初始化，初始化信息为：%s", json.ToJsonIgnoreError(initOrgBo))
		return 0, errs.BuildSystemErrorInfo(errs.OrgNotNeedInitError)
	}

	orgId, err := OrgInfoInit(initOrgBo, tx)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	_, err = OrgOutInfoInit(initOrgBo, orgId, tx)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	_, err = OrgConfigInfoInit(orgId, tx)
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

	return orgId, nil
}

//组织信息初始化
func OrgInfoInit(initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	isAuth := 0
	if initOrgBo.IsAuthenticated {
		isAuth = 1
	}

	orgId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrganization)
	if err != nil {
		return 0, err
	}
	org := &po.PpmOrgOrganization{}
	org.Id = orgId
	org.Status = consts.AppStatusEnable
	org.IsDelete = consts.AppIsNoDelete
	org.SourceChannel = initOrgBo.SourceChannel
	org.Name = initOrgBo.OrgName
	org.LogoUrl = initOrgBo.OrgLogo
	org.Address = initOrgBo.OrgProvince + initOrgBo.OrgCity
	org.IsAuthenticated = isAuth
	err2 := mysql.TransInsert(tx, org)
	if err2 != nil {
		logger.GetDefaultLogger().Error("组织初始化，添加组织时异常:" + strs.ObjectToString(err2))
		return 0, err
	}
	return orgId, nil
}

//组织外部信息初始化
func OrgOutInfoInit(initOrgBo bo.InitOrgBo, orgId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
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
func OrgConfigInfoInit(orgId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
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
