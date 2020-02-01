package dingtalk

import (
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/date"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/sdk/dingtalk"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
	"strconv"
	"time"
)

var log = logger.GetDefaultLogger()


func GetSuiteTicket() (string, error){
	cacheJson, err := cache.Get(consts.CacheDingTalkSuiteTicket)
	log.Infof("飞书AppTicket: %s", cacheJson)
	if err != nil{
		log.Error(err)
		return "", errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if cacheJson == ""{
		log.Error("app ticket为空")
		return "", errs.BuildSystemErrorInfo(errs.FeiShuAppTicketNotExistError)
	}
	cacheBo := &bo.FeiShuAppTicketCacheBo{}
	_ = json.FromJson(cacheJson, cacheBo)
	return cacheBo.AppTicket, nil
}

func SetSuiteTicket(appTicket string) error{
	dingConfig := config.GetDingTalkSdkConfig()
	if dingConfig == nil{
		log.Info("dingtalk config is nil")
		return errs.DingTalkConfigError
	}
	cacheJson := json.ToJsonIgnoreError(bo.FeiShuAppTicketCacheBo{
		AppId: strconv.FormatInt(dingConfig.AppId, 10),
		AppTicket: appTicket,
		LastUpdateTime: date.Format(time.Now()),
	})

	err := cache.Set(consts.CacheDingTalkSuiteTicket, cacheJson)
	if err != nil{
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func GetDingTalkClientRest(corpId string) (*sdk.DingTalkClient, error) {
	suiteTicket, err := GetSuiteTicket()
	if err != nil {
		return nil, err
	}
	return dingtalk.GetDingTalkClient(corpId, suiteTicket)
}

func GetDingTalkUserRoleBos(corpId string) ([]*bo.DingTalkUserRoleBo, errs.SystemErrorInfo) {
	client, err1 := GetDingTalkClientRest(corpId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.DingTalkClientError, err1)
	}

	fetchChild := false
	resp, err2 := client.GetDeptList(nil, &fetchChild, "1")
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.DingTalkClientError, err2)
	}
	if resp.ErrCode != 0 {
		log.Error(resp.ErrMsg)
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.DingTalkOpenApiCallError, resp.ErrMsg)
	}

	dingTalkUserRoleBos := &[]*bo.DingTalkUserRoleBo{}

	if len(resp.Department) > 0 {
		userIds := &[]string{}

		isError3, err3 := getDepartmentUserIds(client, resp, userIds)

		if isError3 {
			return nil, err3
		}

		//去重
		*userIds = slice.SliceUniqueString(*userIds)

		//赋值用户角色
		for _, userId := range *userIds {
			*dingTalkUserRoleBos = append(*dingTalkUserRoleBos, &bo.DingTalkUserRoleBo{
				UserId: userId,
			})
		}

		isError4, err4 := getAdminUserIds(client, dingTalkUserRoleBos)

		if isError4 {
			return nil, err4
		}
	}

	return *dingTalkUserRoleBos, nil
}

//获取指定部门id下的userIds集合
func getDepartmentUserIds(client *sdk.DingTalkClient, resp sdk.GetDeptListResp, userIds *[]string) (isNil bool, error errs.SystemErrorInfo) {

	deptStringIds := make([]string, 0)

	//获取钉钉部门下的部门id
	for _, dept := range resp.Department {
		deptStringIds = append(deptStringIds, strconv.FormatInt(dept.Id, 10))
	}

	//拼接root部门
	deptStringIds = append(deptStringIds, strconv.FormatInt(1, 10))

	for _, deptId := range deptStringIds {
		userListResp, err3 := client.GetDepMemberIds(deptId)
		if err3 != nil {
			log.Error(err3)
			return true, errs.BuildSystemErrorInfo(errs.DingTalkClientError, err3)
		}
		if userListResp.ErrCode != 0 {
			log.Error(userListResp.ErrMsg)
			return true, errs.BuildSystemErrorInfoWithMessage(errs.DingTalkOpenApiCallError, userListResp.ErrMsg)
		}

		if len(userListResp.UserIds) > 0 {
			*userIds = append(*userIds, userListResp.UserIds...)
		}
	}

	//不需要返回错误 错误内容返回nil
	return false, nil

}

//获取admin用户的userIds
func getAdminUserIds(client *sdk.DingTalkClient, dingTalkUserRoleBos *[]*bo.DingTalkUserRoleBo) (isNil bool, error errs.SystemErrorInfo) {

	//判断是否是root或者admin角色
	adminListResp, err4 := client.GetAdminList()
	if err4 != nil {
		log.Error(err4)
		return true, errs.BuildSystemErrorInfo(errs.DingTalkClientError, err4)
	}
	if adminListResp.ErrCode != 0 {
		log.Error(adminListResp.ErrMsg)
		return true, errs.BuildSystemErrorInfoWithMessage(errs.DingTalkOpenApiCallError, adminListResp.ErrMsg)
	}

	for _, usr := range adminListResp.AdminList {
		for _, dingTalkUserRoleBo := range *dingTalkUserRoleBos {
			if dingTalkUserRoleBo.UserId == usr.UserId {
				if usr.SysLevel == 1 {
					(*dingTalkUserRoleBo).IsRoot = true
				} else {
					(*dingTalkUserRoleBo).IsAdmin = true
				}
				break
			}
		}
	}

	//不需要返回错误 错误内容返回nil
	return false, nil

}

/**
获取部门列表
*/
func GetScopeDeps(corpId string) ([]sdk.DepartmentInfo, errs.SystemErrorInfo) {
	departmentInfos := make([]sdk.DepartmentInfo, 0)

	client, err1 := GetDingTalkClientRest(corpId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.DingTalkClientError, err1)
	}

	scopes, err := client.GetAuthScopes()
	if err != nil{
		log.Error(err)
		return nil, errs.DingTalkClientError
	}
	if scopes.ErrCode != 0{
		log.Error(scopes.ErrMsg)
		return nil, errs.DingTalkClientError
	}

	authDeps := scopes.AuthOrgScopes.AuthedDept
	if authDeps != nil && len(authDeps) > 0{
		fetchChild := true
		for _, depId := range authDeps{

			depDetailResp, err := client.GetDeptDetail(strconv.FormatInt(depId, 10), nil)
			if err != nil{
				log.Error(err)
				return nil, errs.DingTalkClientError
			}
			if depDetailResp.ErrCode != 0{
				log.Error(depDetailResp.ErrMsg)
				return nil, errs.DingTalkClientError
			}
			departmentInfos = append(departmentInfos, sdk.DepartmentInfo{
				Id: depDetailResp.Id,
				Name: depDetailResp.Name,
				ParentId: depDetailResp.ParentId,
				CreateDeptGroup: depDetailResp.CreateDeptGroup,
				AutoAddUser: depDetailResp.AutoAddUser,
			})

			resp, err2 := client.GetDeptList(nil, &fetchChild, strconv.FormatInt(depId, 10))
			if err2 != nil{
				log.Error(err2)
				return nil, errs.DingTalkClientError
			}
			if resp.ErrCode != 0{
				log.Error(resp.ErrMsg)
				return nil, errs.DingTalkClientError
			}
			if resp.Department != nil{
				departmentInfos = append(departmentInfos, resp.Department...)
			}
		}
	}

	//在最后做过滤
	depMap := map[int64]sdk.DepartmentInfo{}
	for _, dep := range departmentInfos{
		depMap[dep.Id] = dep
	}

	departmentInfos = make([]sdk.DepartmentInfo, 0)
	for _, dep := range depMap{
		departmentInfos = append(departmentInfos, dep)
	}
	return departmentInfos, nil
}

func GetScopeUsers(corpId string) ([]sdk.UserList, errs.SystemErrorInfo){
	client, err1 := GetDingTalkClientRest(corpId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.DingTalkClientError, err1)
	}

	scopes, err := client.GetAuthScopes()
	if err != nil{
		log.Error(err)
		return nil, errs.DingTalkClientError
	}
	if scopes.ErrCode != 0{
		log.Error(scopes.ErrMsg)
		return nil, errs.DingTalkClientError
	}

	scopeUserIds := scopes.AuthOrgScopes.AuthedUser

	//获取其余部门下的用户
	deps, depsErr := GetScopeDeps(corpId)
	if depsErr != nil{
		log.Error(depsErr)
		return nil, depsErr
	}

	//unionId - userInfo
	userInfoMap := map[string]sdk.UserList{}
	//userId - bool
	userInfoExistMap := map[string]bool{}
	for _, dep := range deps{

		offset := int64(0)
		size := int64(50)

		for ;;{
			//获取对应部门下的用户
			resp, err := client.GetDepMemberDetailList(strconv.FormatInt(dep.Id, 10), consts.AppSourceChannelDingTalkDefaultLang, offset, size, "")
			if err != nil{
				log.Error(err)
				return nil, errs.DingTalkClientError
			}
			if resp.ErrCode != 0{
				log.Error(resp.ErrMsg)
				return nil, errs.DingTalkClientError
			}

			userList := resp.UserList
			if userList != nil && len(userList) > 0{
				for _, user := range userList{
					userInfoMap[user.UnionId] = user
					userInfoExistMap[user.UserId] = true
				}
			}

			if !resp.HasMore{
				break
			}

			offset += size
		}
	}

	//获取授权范围内的用户信息
	if scopeUserIds != nil && len(scopeUserIds) > 0{
		for _, userId := range scopeUserIds{
			if _, ok := userInfoExistMap[userId]; ! ok{
				resp, err := client.GetUserDetail(userId, nil)
				if err != nil{
					log.Error(err)
					return nil, errs.DingTalkClientError
				}
				if resp.ErrCode != 0{
					log.Error(resp.ErrMsg)
					return nil, errs.DingTalkClientError
				}

				userInfoMap[resp.UnionId] = resp.UserList
			}
		}
	}

	userList := make([]sdk.UserList, 0)
	for _, user := range userInfoMap{
		userList = append(userList, user)
	}

	return userList, nil
}