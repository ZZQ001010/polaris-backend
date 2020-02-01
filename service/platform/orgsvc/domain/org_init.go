package domain

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/common/sdk/dingtalk"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/pinyin"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/idvo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"github.com/pkg/errors"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"

	"strconv"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

const larkDepartmentInitSql = consts.TemplateDirPrefix + "lark_department_init.template"
const larkUserInitSql = consts.TemplateDirPrefix + "lark_user_init.template"

func OrgInit(corpId string, permanentCode string, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	org := &po.PpmOrgOrganization{}
	orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
	orgAuthTicketInfo := &bo.OrgAuthTicketInfoBo{}

	err := tx.Collection(orgOutInfo.TableName()).Find(db.Cond{consts.TcOutOrgId: corpId, consts.TcSourceChannel: consts.AppSourceChannelDingTalk}).One(orgOutInfo)

	//有组织的时候就重置一下 没有的时候就初始化
	if err == nil {
		//组织已存在
		logger.GetDefaultLogger().Infof("组织 %s 已经存在, 删除状态 %d", corpId, orgOutInfo.IsDelete)

		if orgOutInfo.IsDelete == consts.AppIsDeleted {
			orgOutInfo.IsDelete = consts.AppIsNoDelete
		}
		if orgOutInfo.AuthTicket != "" {
			json.FromJson(orgOutInfo.AuthTicket, orgAuthTicketInfo)
		}
		orgAuthTicketInfo.PermanentCode = permanentCode
		newAuthTicket, _ := json.ToJson(orgAuthTicketInfo)

		orgOutInfo.AuthTicket = newAuthTicket
		AssemblyOrgOutInfo(corpId, orgOutInfo)

		err = mysql.TransUpdate(tx, orgOutInfo)
		if err != nil {
			logger.GetDefaultLogger().Error("组织初始化，更新组织时异常" + strs.ObjectToString(err))
			return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		return orgOutInfo.OrgId, nil
	} else {
		logger.GetDefaultLogger().Infof("组织 %s 准备创建", corpId)
		//申请Id
		orgOutInfoId, orgId, err := dealOrgOutInfoIdAndorgId(orgOutInfo, org)

		if err != nil {
			return 0, err
		}

		orgOutInfo.Id = orgOutInfoId
		orgOutInfo.OrgId = orgId
		orgOutInfo.OutOrgId = corpId
		orgOutInfo.IsDelete = consts.AppIsNoDelete
		orgOutInfo.Status = consts.AppStatusEnable
		orgOutInfo.SourceChannel = consts.AppSourceChannelDingTalk

		orgAuthTicketInfo.PermanentCode = permanentCode
		newAuthTicket, _ := json.ToJson(orgAuthTicketInfo)

		orgOutInfo.AuthTicket = newAuthTicket

		//调用钉钉获取一些企业外部的信息
		AssemblyOrgOutInfo(corpId, orgOutInfo)

		err2 := mysql.TransInsert(tx, orgOutInfo)
		if err2 != nil {
			logger.GetDefaultLogger().Error("组织初始化，添加外部组织信息时异常: " + strs.ObjectToString(err2))
			return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
		}

		org.Id = orgId
		org.Status = consts.AppStatusEnable
		org.IsDelete = consts.AppIsNoDelete
		org.SourceChannel = consts.AppSourceChannelDingTalk

		AssemblyOrg(corpId, org)

		err2 = mysql.TransInsert(tx, org)
		if err2 != nil {
			logger.GetDefaultLogger().Error("组织初始化，添加组织时异常:" + strs.ObjectToString(err2))
			return 0, err
		}
		return orgId, nil
	}
}

func dealOrgOutInfoIdAndorgId(orgOutInfo *po.PpmOrgOrganizationOutInfo, org *po.PpmOrgOrganization) (int64, int64, errs.SystemErrorInfo) {
	orgOutInfoIdVo := idfacade.ApplyPrimaryId(idvo.ApplyPrimaryIdReqVo{Code: orgOutInfo.TableName()})
	if orgOutInfoIdVo.Failure() {
		return int64(0), int64(0), orgOutInfoIdVo.Error()
	}

	orgIdVo := idfacade.ApplyPrimaryId(idvo.ApplyPrimaryIdReqVo{Code: org.TableName()})
	if orgIdVo.Failure() {
		return int64(0), int64(0), orgIdVo.Error()
	}

	return orgOutInfoIdVo.Id, orgIdVo.Id, nil
}

func OrgOwnerInit(orgId int64, owner int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	org := &po.PpmOrgOrganization{}
	org.Id = orgId
	org.Owner = owner
	err := mysql.TransUpdate(tx, org)
	if err != nil {
		logger.GetDefaultLogger().Error(strs.ObjectToString(err))
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

//组装外部信息
func AssemblyOrgOutInfo(corpId string, orgOutInfo *po.PpmOrgOrganizationOutInfo) errs.SystemErrorInfo {
	//获取企业授权信息
	authInfo, err := GetCorpAuthInfo(corpId)
	if err != nil {
		return err
	}

	authCorpInfo := authInfo.AuthCorpInfo

	orgOutInfo.Name = authCorpInfo.CorpName
	orgOutInfo.Industry = authCorpInfo.Industry

	isAuth := 0
	if authCorpInfo.IsAuthenticated {
		isAuth = 1
	}
	orgOutInfo.IsAuthenticated = isAuth
	orgOutInfo.AuthLevel = strconv.FormatInt(authCorpInfo.AuthLevel, 10)

	return nil
}

func AssemblyOrg(corpId string, org *po.PpmOrgOrganization) errs.SystemErrorInfo {
	authInfo, err := GetCorpAuthInfo(corpId)
	if err != nil {
		return err
	}

	authCorpInfo := authInfo.AuthCorpInfo

	org.Name = authCorpInfo.CorpName
	org.LogoUrl = authCorpInfo.CorpLogoUrl
	org.Address = authCorpInfo.CorpProvince + authCorpInfo.CorpCity

	isAuth := 0
	if authCorpInfo.IsAuthenticated {
		isAuth = 1
	}
	org.IsAuthenticated = isAuth

	return nil
}

func GetCorpAuthInfo(corpId string) (sdk.GetAuthInfoResp, errs.SystemErrorInfo) {
	suiteTicket, err := GetSuiteTicket()
	if err != nil {
		return sdk.GetAuthInfoResp{}, err
	}
	//创建企业对象
	corpProxy := dingtalk.GetSDKProxy().CreateCorp(corpId, suiteTicket)
	resp, err2 := corpProxy.GetAuthInfo()
	if err2 != nil {
		return sdk.GetAuthInfoResp{}, errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, err2)
	}

	if resp.ErrCode != 0 {
		return sdk.GetAuthInfoResp{}, errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, errors.New(resp.ErrMsg))
	}
	resp, err2 = corpProxy.GetAuthInfo()
	if err2 != nil {
		return resp, errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, err2)
	}
	return resp, nil
}

func GetSuiteTicket() (string, errs.SystemErrorInfo) {
	val, err := cache.Get(consts.CacheDingTalkSuiteTicket)
	if err != nil {
		return val, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	return val, nil
}

//
//func InitDepartment(orgId int64, corpId string, tx sqlbuilder.Tx) errs.SystemErrorInfo {
//	deptList, err := dingtalk2.GetDeptList(corpId)
//	if err != nil {
//		return err
//	}
//	//把组织作为根部门插入
//	client, err1 := dingtalk2.GetDingTalkClientRest(corpId)
//	if err1 != nil {
//		log.Error(err1)
//		return errs.BuildSystemErrorInfo(errs.DingTalkClientError, err1)
//	}
//	rootDept, rootErr := client.GetDeptDetail("1", nil)
//	if rootErr != nil {
//		log.Error(rootErr)
//		return errs.BuildSystemErrorInfo(errs.DingTalkClientError, rootErr)
//	}
//	deptList.Department = append(deptList.Department, struct {
//		Id              int64  `json:"id"`
//		Name            string `json:"name"`
//		ParentId        int64  `json:"parentid"`
//		CreateDeptGroup bool   `json:"createDeptGroup"`
//		AutoAddUser     bool   `json:"autoAddUser"`
//	}{
//		Id:              rootDept.Id,
//		Name:            rootDept.Name,
//		ParentId:        rootDept.ParentId,
//		CreateDeptGroup: rootDept.CreateDeptGroup,
//		AutoAddUser:     rootDept.AutoAddUser,
//	})
//
//	departmentInfo := make([]interface{}, len(deptList.Department))
//	outDepartmentInfo := make([]interface{}, len(deptList.Department))
//	outIds := map[int64]int64{}
//
//	for k, v := range deptList.Department {
//		id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableDepartment)
//		if err != nil {
//			logger.GetDefaultLogger().Error(strs.ObjectToString(err))
//			return err
//		}
//		outId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableDepartmentOutInfo)
//		if err != nil {
//			logger.GetDefaultLogger().Error(strs.ObjectToString(err))
//			return err
//		}
//
//		outIds[v.Id] = id
//		departmentInfo[k] = &po.PpmOrgDepartment{
//			Id:            id,
//			OrgId:         orgId,
//			Name:          v.Name,
//			ParentId:      v.ParentId,
//			SourceChannel: consts.AppSourceChannelDingTalk,
//		}
//
//		outDepartmentInfo[k] = po.PpmOrgDepartmentOutInfo{
//			Id:                       outId,
//			OrgId:                    orgId,
//			DepartmentId:             id,
//			SourceChannel:            consts.AppSourceChannelDingTalk,
//			OutOrgDepartmentId:       strconv.FormatInt(v.Id, 10),
//			Name:                     v.Name,
//			OutOrgDepartmentParentId: strconv.FormatInt(v.ParentId, 10),
//		}
//	}
//
//	for k, v := range departmentInfo {
//		if a, ok := v.(*po.PpmOrgDepartment); ok {
//			departmentInfo[k].(*po.PpmOrgDepartment).ParentId = outIds[a.ParentId]
//		}
//	}
//
//	departErr := mysql.TransBatchInsert(tx, &po.PpmOrgDepartment{}, departmentInfo)
//	if departErr != nil {
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, departErr)
//	}
//	outDepartErr := mysql.TransBatchInsert(tx, &po.PpmOrgDepartmentOutInfo{}, outDepartmentInfo)
//	if outDepartErr != nil {
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, outDepartErr)
//	}
//	return nil
//}

func LarkDepartmentInit(orgId int64, sourceChannel, sourcePlatform string, orgName string, creator int64) (int64, errs.SystemErrorInfo) {
	deparmentVo := idfacade.ApplyPrimaryId(idvo.ApplyPrimaryIdReqVo{Code: consts.TableDepartment})
	if deparmentVo.Failure() {
		log.Error(deparmentVo.Message)
		return 0, deparmentVo.Error()
	}
	//outDeparmentVo := idfacade.ApplyPrimaryId(idvo.ApplyPrimaryIdReqVo{Code: consts.TableDepartmentOutInfo})
	//if outDeparmentVo.Failure() {
	//	log.Error(outDeparmentVo.Message)
	//	return 0, outDeparmentVo.Error()
	//}
	contextMap := map[string]interface{}{}
	contextMap["OrgId"] = orgId
	contextMap["OrgName"] = orgName
	contextMap["DepartmentId"] = deparmentVo.Id
	//contextMap["OutDepartmentId"] = outDeparmentVo.Id
	contextMap["SourceChannel"] = sourceChannel
	contextMap["SourcePlatform"] = sourcePlatform
	contextMap["Creator"] = creator
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		insertErr := util.ReadAndWrite(larkDepartmentInitSql, contextMap, tx)
		if insertErr != nil {
			return errs.BuildSystemErrorInfo(errs.BaseDomainError, insertErr)
		}

		return nil
	})
	if err != nil {
		return 0, errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}

	return deparmentVo.Id, nil
}

//初始化张三和李四
func LarkUserInit(orgId int64, sourceChannel, sourcePlatform string, departmentId int64) (int64, int64, errs.SystemErrorInfo) {
	count := 2
	contextMap := map[string]interface{}{}
	contextMap["OrgId"] = orgId
	contextMap["UserName1"] = "张三"
	contextMap["UserNamePy1"] = pinyin.ConvertToPinyin("张三")
	contextMap["UserName2"] = "李四"
	contextMap["UserNamePy2"] = pinyin.ConvertToPinyin("张三")
	contextMap["SourceChannel"] = sourceChannel
	contextMap["SourcePlatform"] = sourcePlatform
	contextMap["DepartmentId"] = departmentId

	//user id 申请
	userIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUser, count)
	if err != nil {
		return 0, 0, err
	}
	userIdsCount := 1
	var zhangsan, lisi int64
	for _, v := range userIds.Ids {
		contextMap["UserId"+strconv.Itoa(userIdsCount)] = v.Id
		if userIdsCount == 1 {
			zhangsan = v.Id
		} else if userIdsCount == 2 {
			lisi = v.Id
		}
		userIdsCount++
	}

	//userconfig id 申请
	userconfigIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserConfig, count)
	if err != nil {
		return 0, 0, err
	}
	userconfigIdsCount := 1
	for _, v := range userconfigIds.Ids {
		contextMap["UserConfigId"+strconv.Itoa(userconfigIdsCount)] = v.Id
		userconfigIdsCount++
	}

	//UserOutInfo id 申请
	userOutInfoIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserOutInfo, count)
	if err != nil {
		return 0, 0, err
	}
	userOutInfoIdsCount := 1
	for _, v := range userOutInfoIds.Ids {
		contextMap["UserOutInfoId"+strconv.Itoa(userOutInfoIdsCount)] = v.Id
		userOutInfoIdsCount++
	}

	//UserOrg id 申请
	userOrgIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserOrganization, count)
	if err != nil {
		return 0, 0, err
	}
	userOrgIdsCount := 1
	for _, v := range userOrgIds.Ids {
		contextMap["UserOrgId"+strconv.Itoa(userOrgIdsCount)] = v.Id
		userOrgIdsCount++
	}

	//UserDept id 申请
	userDeptIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, count)
	if err != nil {
		return 0, 0, err
	}
	userDeptIdsCount := 1
	for _, v := range userDeptIds.Ids {
		contextMap["UserDeptId"+strconv.Itoa(userDeptIdsCount)] = v.Id
		userDeptIdsCount++
	}

	//RoleUser id 申请
	roleIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableRoleUser, count)
	if err != nil {
		return 0, 0, err
	}
	roleIdsCount := 1
	for _, v := range roleIds.Ids {
		contextMap["RoleUserId"+strconv.Itoa(roleIdsCount)] = v.Id
		roleIdsCount++
	}
	err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		insertErr := util.ReadAndWrite(larkUserInitSql, contextMap, tx)
		if insertErr != nil {
			return insertErr
		}

		return nil
	})

	if err1 != nil {
		return 0, 0, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	return zhangsan, lisi, nil
}
