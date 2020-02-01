package service

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/core/util/pinyin"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/domain"
	"strings"
	"time"
)

var userErrorTemplate = " GetCurrentUser : %v\n"

func PersonalInfo(orgId, userId int64, sourceChannel string) (*vo.PersonalInfo, errs.SystemErrorInfo) {

	userInfoBo, passwordSet, err1 := domain.GetUserInfo(orgId, userId, sourceChannel)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.GetUserInfoError, err1)
	}

	personalInfo := &vo.PersonalInfo{}
	copyErr := copyer.Copy(userInfoBo, personalInfo)
	if copyErr != nil {
		logger.GetDefaultLogger().Error(strs.ObjectToString(copyErr))
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	if passwordSet {
		personalInfo.PasswordSet = 1
	}

	return personalInfo, nil
}

func GetUserIds(orgId int64, corpId, sourceChannel string, empIds []string) ([]*vo.UserIDInfo, errs.SystemErrorInfo) {
	resultIds := make([]*vo.UserIDInfo, len(empIds))
	for i, empId := range empIds {
		baseUserInfo, err := domain.GetBaseUserInfoByEmpId(sourceChannel, orgId, empId)
		if err != nil {
			baseUserInfo, err = domain.UserInit(orgId, corpId, empId, sourceChannel)
			if err != nil {
				log.Error(err)
				return nil, errs.BuildSystemErrorInfo(errs.UserNotInitError, err)
			}
		}
		resultIds[i] = &vo.UserIDInfo{
			UserID:     baseUserInfo.UserId,
			Name:       baseUserInfo.Name,
			Avatar:     baseUserInfo.Avatar,
			EmplID:     baseUserInfo.OutUserId,
			IsDeleted:  baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted,
			IsDisabled: baseUserInfo.OrgUserStatus == consts.AppStatusDisabled,
		}
	}
	return resultIds, nil
}

func GetUserId(orgId int64, corpId, sourceChannel, empId string) (*vo.UserIDInfo, errs.SystemErrorInfo) {
	userIdInfos, err := GetUserIds(orgId, corpId, sourceChannel, []string{empId})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if userIdInfos == nil || len(userIdInfos) == 0 {
		log.Errorf("GetUserIds 获取到空 %d %s %s %s", orgId, corpId, sourceChannel, empId)
		return nil, errs.BuildSystemErrorInfo(errs.UserNotExist)
	}
	return userIdInfos[0], nil
}

func UserConfigInfo(orgId, userId int64) (*vo.UserConfig, errs.SystemErrorInfo) {
	//cacheUserInfo, err := GetCurrentUser(ctx)
	//if err != nil {
	//	log.Errorf(userErrorTemplate, err)
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//orgId := cacheUserInfo.OrgId
	//userId := cacheUserInfo.UserId

	userConfig, err := domain.GetUserConfigInfo(orgId, userId)
	if err != nil {
		userConfig = &bo.UserConfigBo{}
		userConfigBo, err2 := domain.InsertUserConfig(orgId, userId)
		if err2 != nil {
			log.Error(err2)
			return nil, errs.BuildSystemErrorInfo(errs.UserConfigUpdateError, err2)
		}
		err3 := copyer.Copy(userConfigBo, userConfig)
		if err3 != nil {
			log.Error(err3)
			return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err3)
		}
	}

	configInfo := &vo.UserConfig{}
	copyErr := copyer.Copy(userConfig, configInfo)
	if copyErr != nil {
		logger.GetDefaultLogger().Error(strs.ObjectToString(copyErr))
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return configInfo, nil
}

func UpdateUserConfig(orgId, userId int64, input vo.UpdateUserConfigReq) (*vo.UpdateUserConfigResp, errs.SystemErrorInfo) {
	userConfigBo := &bo.UserConfigBo{}
	copyErr := copyer.Copy(input, userConfigBo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err1 := domain.UpdateUserConfig(orgId, userId, *userConfigBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.UserConfigUpdateError, err1)
	}

	//Nico: 暂时先这样做双写，之后优化
	err1 = domain.DeleteUserConfigInfo(orgId, userId)
	if err1 != nil {
		log.Error(err1)
	}

	return &vo.UpdateUserConfigResp{
		ID: input.ID,
	}, nil
}


func UpdateUserPcConfig(orgId, userId int64, input vo.UpdateUserPcConfigReq) (*vo.UpdateUserConfigResp, errs.SystemErrorInfo) {
	userConfigBo := &bo.UserConfigBo{}

	if util.FieldInUpdate(input.UpdateFields, "pcNoticeOpenStatus"){
		userConfigBo.PcNoticeOpenStatus = *input.PcNoticeOpenStatus
	}
	if util.FieldInUpdate(input.UpdateFields, "pcIssueRemindMessageStatus"){
		userConfigBo.PcIssueRemindMessageStatus = *input.PcIssueRemindMessageStatus
	}
	if util.FieldInUpdate(input.UpdateFields, "pcOrgMessageStatus"){
		userConfigBo.PcOrgMessageStatus = *input.PcOrgMessageStatus
	}
	if util.FieldInUpdate(input.UpdateFields, "pcProjectMessageStatus"){
		userConfigBo.PcProjectMessageStatus = *input.PcProjectMessageStatus
	}
	if util.FieldInUpdate(input.UpdateFields, "pcCommentAtMessageStatus"){
		userConfigBo.PcCommentAtMessageStatus = *input.PcCommentAtMessageStatus
	}

	err1 := domain.UpdateUserPcConfig(orgId, userId, *userConfigBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.UserConfigUpdateError, err1)
	}

	//Nico: 暂时先这样做双写，之后优化
	err1 = domain.DeleteUserConfigInfo(orgId, userId)
	if err1 != nil {
		log.Error(err1)
	}

	return &vo.UpdateUserConfigResp{
		ID: 0,
	}, nil
}

func UpdateUserDefaultProjectIdConfig(orgId, userId int64, input vo.UpdateUserDefaultProjectConfigReq) (*vo.UpdateUserConfigResp, errs.SystemErrorInfo) {
	userConfigBo, err := domain.GetUserConfigInfo(orgId, userId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defaultProjectId := input.DefaultProjectID

	cacheProjectInfoResp := projectfacade.GetCacheProjectInfo(projectvo.GetCacheProjectInfoReqVo{
		OrgId:     orgId,
		ProjectId: defaultProjectId,
	})
	if cacheProjectInfoResp.Failure() {
		log.Error(cacheProjectInfoResp.Message)
		return nil, cacheProjectInfoResp.Error()
	}

	err1 := domain.UpdateUserDefaultProjectIdConfig(orgId, userId, *userConfigBo, defaultProjectId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.UserConfigUpdateError, err1)
	}

	err1 = domain.DeleteUserConfigInfo(orgId, userId)
	if err1 != nil {
		log.Error(err1)
	}

	return &vo.UpdateUserConfigResp{
		ID: userConfigBo.ID,
	}, nil
}

func UpdateUserInfo(orgId, userId int64, input vo.UpdateUserInfoReq) (*vo.Void, errs.SystemErrorInfo) {

	upd := &mysql.Upd{}
	//头像
	assemblyAvatar(input, upd)
	//姓名
	nameErr := assemblyName(input, upd)

	if nameErr != nil {
		log.Error(nameErr)
		return nil, nameErr
	}

	//出生日期
	assemblyBirthday(input, upd)
	//性别
	sexErr := assemblySex(input, upd)

	if sexErr != nil {
		log.Error(sexErr)
		return nil, sexErr
	}

	err := domain.UpdateUserInfo(userId, *upd)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = domain.ClearBaseUserInfo(orgId, userId)
	if err != nil {
		log.Error(err)
	}

	return &vo.Void{
		ID: userId,
	}, nil
}

func assemblySex(input vo.UpdateUserInfoReq, upd *mysql.Upd) errs.SystemErrorInfo {
	if NeedUpdate(input.UpdateFields, "sex") {

		if input.Sex != nil {

			if *input.Sex != consts.Male && *input.Sex != consts.Female {
				return errs.BuildSystemErrorInfo(errs.UserSexFail)
			}
			(*upd)[consts.TcSex] = *input.Sex
		}
	}
	return nil
}

func assemblyBirthday(input vo.UpdateUserInfoReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "birthday") {

		if input.Birthday != nil {
			birthday := time.Time(*input.Birthday)
			(*upd)[consts.TcBirthday] = birthday
		}
	}
}

//组装个人头像信息
func assemblyAvatar(input vo.UpdateUserInfoReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "avatar") {

		if input.Avatar != nil {
			(*upd)[consts.TcAvatar] = *input.Avatar
		}
	}
}

//组装名字
func assemblyName(input vo.UpdateUserInfoReq, upd *mysql.Upd) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "name") {

		if input.Name != nil {

			name := strings.Trim(*input.Name, " ")
			//nameLen := str.CountStrByGBK(name)
			//
			//if nameLen == 0 || nameLen > 20 {
			//	log.Error("姓名长度错误")
			//	return errs.BuildSystemErrorInfo(errs.UserNameLenError)
			//}
			isNameRight := format.VerifyUserNameFormat(name)
			if !isNameRight {
				return errs.BuildSystemErrorInfo(errs.UserNameLenError)
			}
			(*upd)[consts.TcName] = name
			(*upd)[consts.TcNamePinyin] = pinyin.ConvertToPinyin(name)
		}
	}
	return nil
}

func GetUserInfoByUserIds(input orgvo.GetUserInfoByUserIdsReqVo) (*[]orgvo.GetUserInfoByUserIdsRespVo, errs.SystemErrorInfo) {

	bos, err := domain.GetBaseUserInfoBatch("", input.OrgId, input.UserIds)

	if err != nil {
		return nil, err
	}

	vos := &[]orgvo.GetUserInfoByUserIdsRespVo{}

	copyError := copyer.Copy(bos, vos)

	if copyError != nil {
		log.Error(copyError)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyError)
	}
	return vos, nil
}

func VerifyOrg(orgId int64, userId int64) bool {
	return domain.VerifyOrg(orgId, userId)
}

func VerifyOrgUsers(orgId int64, userIds []int64) bool {
	return domain.VerifyOrgUsers(orgId, userIds)
}

func GetUserInfo(orgId int64, userId int64, sourceChannel string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	res, _, err := domain.GetUserInfo(orgId, userId, sourceChannel)
	return res, err
}

func GetOutUserInfoListBySourceChannel(sourceChannel string, page int, size int) ([]bo.UserOutInfoBo, errs.SystemErrorInfo) {
	return domain.GetOutUserInfoListBySourceChannel(sourceChannel, page, size)
}

func GetUserInfoListByOrg(orgId int64) ([]bo.SimpleUserInfoBo, errs.SystemErrorInfo) {
	return domain.GetUserInfoListByOrg(orgId)
}
