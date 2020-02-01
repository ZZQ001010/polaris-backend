package service

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/core/util/pinyin"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/commonvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/facade/commonfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/domain"
	"strings"
)

func GetOrgBoList() ([]bo.OrganizationBo, errs.SystemErrorInfo) {
	return domain.GetOrgBoList()
}

func GetBaseOrgInfoByOutOrgId(sourceChannel string, outOrgId string) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseOrgInfoByOutOrgId(sourceChannel, outOrgId)
}

func CreateOrg(req orgvo.CreateOrgReqVo, sourceChannel, sourcePlatform string) (int64, errs.SystemErrorInfo) {
	creatorId := req.Data.CreatorId
	createReqInfo := req.Data.CreateOrgReq
	userInfoBo, _, err := domain.GetUserBo(creatorId)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return 0, err
	}
	orgId, err := domain.CreateOrg(bo.CreateOrgBo{OrgName: createReqInfo.OrgName}, creatorId, sourceChannel, sourcePlatform, "")
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return 0, err
	}

	//初始化组织相关资源
	err = CreateOrgRelationResource(orgId, creatorId, sourceChannel, sourcePlatform, req.Data.CreateOrgReq.OrgName)
	if err != nil {
		log.Error(err)
		return orgId, err
	}

	//用户和组织关联
	err = domain.AddUserOrgRelation(orgId, creatorId, true, false, false)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return orgId, err
	}
	upd := mysql.Upd{}
	if userInfoBo.OrgID == 0 {
		//更新用户的orgId
		upd[consts.TcOrgId] = orgId
	}
	if createReqInfo.CreatorName != nil {
		creatorName := strings.Trim(*createReqInfo.CreatorName, " ")
		//creatorNameLen := str.CountStrByGBK(creatorName)
		//if creatorNameLen == 0 || creatorNameLen > 20 {
		//	log.Error("姓名长度错误")
		//	return orgId, errs.BuildSystemErrorInfo(errs.UserNameLenError)
		//}
		isNameRight := format.VerifyUserNameFormat(creatorName)
		if !isNameRight {
			return orgId, errs.BuildSystemErrorInfo(errs.UserNameLenError)
		}

		//更新用户的名称
		upd[consts.TcName] = creatorName
		upd[consts.TcNamePinyin] = pinyin.ConvertToPinyin(creatorName)
	}
	if len(upd) > 0 {
		err = domain.UpdateUserInfo(creatorId, upd)
		if err != nil {
			log.Info(strs.ObjectToString(err))
			return orgId, err
		}
	}

	//刷新用户缓存
	userToken := req.Data.UserToken
	err = UpdateCacheUserInfoOrgId(userToken, orgId)
	if err != nil {
		log.Info(strs.ObjectToString(err))
	}

	err = domain.ClearBaseUserInfo(orgId, creatorId)
	if err != nil {
		log.Error(err)
	}

	return orgId, nil
}

func bindingUserAndRole(isRoot, isAdmin bool, orgId, userId int64, roleInitResp *bo.RoleInitResp) errs.SystemErrorInfo {
	if isRoot {
		err2 := rolefacade.RoleUserRelationRelaxed(orgId, userId, roleInitResp.OrgSuperAdminRoleId)
		if err2 != nil {
			log.Error(err2)
			return err2
		}
	}
	if isAdmin {
		err2 := rolefacade.RoleUserRelationRelaxed(orgId, userId, roleInitResp.OrgNormalAdminRoleId)
		if err2 != nil {
			log.Error(err2)
			return err2
		}
	}
	return nil
}

func CreateOrgRelationResource(orgId int64, creatorId int64, sourceChannel, sourcePlatform string, orgName string) errs.SystemErrorInfo {
	//初始化角色、优先级
	//权限、角色初始化
	roleInitResp := rolefacade.RoleInit(rolevo.RoleInitReqVo{
		OrgId: orgId,
	})
	if roleInitResp.Failure() {
		logger.GetDefaultLogger().Error("Rollback" + roleInitResp.Message)
		return roleInitResp.Error()
	}
	log.Info("权限、角色初始化成功")

	//为用户绑定超级管理员
	err := bindingUserAndRole(true, false, orgId, creatorId, roleInitResp.RoleInitResp)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("组织超级管理员角色赋予成功")

	//优先级初始化
	priorityInfo := projectfacade.ProjectInit(projectvo.ProjectInitReqVo{OrgId: orgId})
	if priorityInfo.Failure() {
		return errs.BuildSystemErrorInfo(errs.BaseDomainError, priorityInfo.Error())
	}
	log.Info("优先级初始化成功")

	//部门初始化
	departmentId, departmentErr := domain.LarkDepartmentInit(orgId, sourceChannel, sourcePlatform, orgName, creatorId)
	if departmentErr != nil {
		return errs.BuildSystemErrorInfo(errs.BaseDomainError, departmentErr)
	}
	log.Info("部门初始化成功")

	//用户和顶级部门绑定
	err = domain.BoundDepartmentUser(orgId, []int64{creatorId}, departmentId, creatorId, true)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func UserOrganizationList(userId int64) (*vo.UserOrganizationListResp, errs.SystemErrorInfo) {

	organizationBo, err := domain.GetUserOrganizationIdList(userId)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	orgIds := []int64{}
	enableOrgIds := []int64{}
	enableStatus := consts.AppStatusEnable
	disabledStatus := consts.AppStatusDisabled
	for _, value := range *organizationBo {
		orgIds = append(orgIds, value.OrgId)
		if value.Status == enableStatus {
			enableOrgIds = append(enableOrgIds, value.OrgId)
		}
	}

	bos, err := domain.GetOrgBoListByIds(orgIds)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	resultList := &[]*vo.UserOrganization{}
	copyErr := copyer.Copy(bos, resultList)

	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	for k, v := range *resultList {
		if ok, _ := slice.Contain(enableOrgIds, v.ID); ok {
			(*resultList)[k].OrgIsEnabled = &enableStatus
		} else {
			(*resultList)[k].OrgIsEnabled = &disabledStatus
		}
	}

	return &vo.UserOrganizationListResp{
		List: *resultList,
	}, nil
}

func SwitchUserOrganization(orgId, userId int64, token string) errs.SystemErrorInfo {
	//监测可用性
	baseUserInfo, err := GetBaseUserInfo("", orgId, userId)
	if err != nil {
		log.Error(err)
		return err
	}

	err = baseUserInfoOrgStatusCheck(*baseUserInfo)
	if err != nil {
		log.Error(err)
		return err
	}

	//更改用户缓存的orgId
	err = UpdateCacheUserInfoOrgId(token, orgId)
	if err != nil {
		log.Error(err)
		return err
	}

	//修改用户默认组织
	updateUserInfoErr := domain.UpdateUserDefaultOrg(userId, orgId)
	if updateUserInfoErr != nil {
		log.Error(updateUserInfoErr)
	}

	return nil
}

//获取组织信息
func OrganizationInfo(req orgvo.OrganizationInfoReqVo) (*vo.OrganizationInfoResp, errs.SystemErrorInfo) {

	bo, err := domain.GetOrgBoById(req.OrgId)

	if err != nil {
		return nil, err
	}

	//跨服务查询
	resp := commonfacade.AreaInfo(commonvo.AreaInfoReqVo{
		IndustryID: bo.IndustryId,
		CountryID:  bo.CountryId,
		ProvinceID: bo.ProvinceId,
		CityID:     bo.CityId,
	})

	if resp.Failure() {
		log.Error(resp.Message)
		return nil, resp.Error()
	}

	infoResp := vo.OrganizationInfoResp{
		OrgID:         bo.Id,
		OrgName:       bo.Name,
		Code:          bo.Code,
		WebSite:       bo.WebSite,
		IndustryID:    bo.IndustryId,
		IndustryName:  resp.AreaInfoResp.IndustryName,
		Scale:         bo.Scale,
		CountryID:     bo.CountryId,
		CountryCname:  resp.AreaInfoResp.CountryCname,
		ProvinceID:    bo.ProvinceId,
		ProvinceCname: resp.AreaInfoResp.ProvinceCname,
		CityID:        bo.CityId,
		CityCname:     resp.AreaInfoResp.CityCname,
		Address:       bo.Address,
		LogoURL:       bo.LogoUrl,
		Owner:         bo.Owner,
	}

	return &infoResp, nil
}

//对于自己创建的组织，暂时不支持转让
//
//对于加入的企业只有查看全，无操作权
//
//暂时只做基本设置
func UpdateOrganizationSetting(req orgvo.UpdateOrganizationSettingReqVo) (int64, errs.SystemErrorInfo) {

	input := req.Input
	// Owns转让的成员需要判断是否在这个组织里面 暂定
	organizationBo, err := domain.GetOrgBoById(input.OrgID)

	if err != nil {
		return 0, err
	}
	//不是所有者不可以更改信息
	if organizationBo.Owner != req.UserId {
		return 0, errs.BuildSystemErrorInfo(errs.OrgOwnTransferError)
	}

	//更改Own的接口拆开来,  这一期暂时也不做
	updateOrgBo, err := assemblyOrganizationBo(input, req.UserId, organizationBo)

	if err != nil {
		return 0, err
	}

	err = domain.UpdateOrg(*updateOrgBo)

	if err != nil {
		return 0, err
	}
	return input.OrgID, nil
}

func assemblyOrganizationBo(input vo.UpdateOrganizationSettingsReq, userId int64, orgOrganization *bo.OrganizationBo) (*bo.UpdateOrganizationBo, errs.SystemErrorInfo) {
	//公用初始化
	orgBo := bo.OrganizationBo{Id: input.OrgID}

	upd := &mysql.Upd{}
	//名字
	err := assemblyOrgName(input, upd)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//网址
	err = assemblyCode(input, upd)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//行业
	assemblyIndustryID(input, upd)
	//组织规模
	assemblyScaleID(input, upd)
	//所在国家
	assemblyCountryID(input, upd)
	// 所在省份
	assemblyProvince(input, upd)
	// 所在城市
	assemblyCity(input, upd)
	// 组织地址
	err = assemblyAddress(input, upd)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// 组织logo地址
	err = assemblyLogoUrl(input, upd)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//sourceChannel
	orgBo.SourceChannel = orgOrganization.SourceChannel

	return &bo.UpdateOrganizationBo{
		Bo:                     orgBo,
		OrganizationUpdateCond: *upd,
	}, nil
}

func assemblyOrgName(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {
	if NeedUpdate(input.UpdateFields, "orgName") {
		orgName := strings.TrimSpace(input.OrgName)
		//orgNameLen := len(orgName)
		//if orgNameLen == 0 || orgNameLen > 256 {
		//	return errs.BuildSystemErrorInfo(errs.OrgNameLenError)
		//}
		isOrgNameRight := format.VerifyOrgNameFormat(orgName)
		if !isOrgNameRight {
			return errs.OrgNameLenError
		}

		(*upd)[consts.TcName] = input.OrgName
	}
	return nil
}

//网址
func assemblyCode(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "code") {
		if input.Code != nil {
			orgCode := *input.Code
			orgCode = strings.TrimSpace(orgCode)

			//判断当前组织有没有设置过code
			organizationBo, err := domain.GetOrgBoById(input.OrgID)
			if err != nil {
				log.Error(err)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
			if organizationBo.Code != consts.BlankString {
				return errs.BuildSystemErrorInfo(errs.OrgCodeAlreadySetError)
			}

			//orgCodeLen := strs.Len(orgCode)
			////判断长度
			//if orgCodeLen > sconsts.OrgCodeLength || orgCodeLen < 1 {
			//	return errs.BuildSystemErrorInfo(errs.OrgCodeLenError)
			//}
			isOrgCodeRight := format.VerifyOrgCodeFormat(orgCode)
			if !isOrgCodeRight {
				return errs.OrgCodeLenError
			}

			_, err = domain.GetOrgBoByCode(orgCode)
			//查不到才能更改
			if err != nil {
				(*upd)[consts.TcCode] = orgCode
			} else {
				return errs.BuildSystemErrorInfo(errs.OrgCodeExistError)
			}
		}
	}

	return nil
}

//组织行业Id
func assemblyIndustryID(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "industryId") {

		if input.IndustryID != nil {
			(*upd)[consts.TcIndustryId] = *input.IndustryID
		} else {
			(*upd)[consts.TcIndustryId] = 0
		}
	}
}

//组织规模
func assemblyScaleID(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "scale") {

		if input.Scale != nil {
			(*upd)[consts.TcScale] = *input.Scale
		} else {
			(*upd)[consts.TcScale] = 0
		}
	}
}

//所在国家
func assemblyCountryID(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "countryId") {

		if input.CountryID != nil {
			(*upd)[consts.TcCountryId] = *input.CountryID
		} else {
			(*upd)[consts.TcCountryId] = 0
		}
	}
}

//省份
func assemblyProvince(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "provinceId") {

		if input.ProvinceID != nil {
			(*upd)[consts.TcProvinceId] = *input.ProvinceID
		} else {
			(*upd)[consts.TcProvinceId] = 0
		}
	}
}

//城市
func assemblyCity(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "cityId") {

		if input.CityID != nil {
			(*upd)[consts.TcCityId] = *input.CityID
		} else {
			(*upd)[consts.TcCityId] = 0
		}
	}
}

//地址
func assemblyAddress(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "address") {

		if input.Address != nil {
			//len := strs.Len(*input.Address)
			//if len > 256 {
			//	return errs.BuildSystemErrorInfo(errs.OrgAddressLenError)
			//}
			isAdressRight := format.VerifyOrgAdressFormat(*input.Address)
			if !isAdressRight {
				return errs.OrgAddressLenError
			}

			(*upd)[consts.TcAddress] = *input.Address
		}
	}
	return nil
}

//Logo
func assemblyLogoUrl(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "logoUrl") {
		if input.LogoURL != nil {
			logoLen := strs.Len(*input.LogoURL)
			if logoLen > 512 {
				return errs.BuildSystemErrorInfo(errs.OrgLogoLenError)
			}

			(*upd)[consts.TcLogoUrl] = *input.LogoURL
		}
	}
	return nil
}

func ScheduleOrganizationPageList(reqVo orgvo.ScheduleOrganizationPageListReqVo) (*orgvo.ScheduleOrganizationPageListResp, errs.SystemErrorInfo) {

	page := reqVo.Page
	size := reqVo.Size

	bos, count, err := domain.ScheduleOrganizationPageList(page, size)

	if err != nil {
		return nil, err
	}

	list := &[]*orgvo.ScheduleOrganizationListResp{}

	copyErr := copyer.Copy(bos, list)

	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &orgvo.ScheduleOrganizationPageListResp{
		Total:                        count,
		ScheduleOrganizationListResp: list,
	}, nil

}

//通过来源获取组织id列表
func GetOrgIdListBySourceChannel(sourceChannel string, page int, size int) ([]int64, errs.SystemErrorInfo) {
	return domain.GetOrgIdListBySourceChannel(sourceChannel, page, size)
}
