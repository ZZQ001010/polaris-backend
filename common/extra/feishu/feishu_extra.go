package feishu

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
	"github.com/galaxy-book/feishu-sdk-golang/sdk"
)

var log = logger.GetDefaultLogger()

func GetTenant(tenantKey string) (*sdk.Tenant, errs.SystemErrorInfo){
	cacheTenantInfo, err := GetTenantAccessToken(tenantKey)
	if err != nil{
		log.Error(err)
		return nil, err
	}
	return &sdk.Tenant{
		TenantAccessToken: cacheTenantInfo.TenantAccessToken,
	}, nil
}

/**
获取部门列表
*/
func GetDeptList(tenantKey string) ([]vo.DepartmentRestInfoVo, errs.SystemErrorInfo) {
	return GetDeptListWithDepId(tenantKey, "0")
}

func GetDeptListWithDepId(tenantKey string, depId string) ([]vo.DepartmentRestInfoVo, errs.SystemErrorInfo) {
	tenant, err := GetTenant(tenantKey)
	if err != nil{
		log.Error(err)
		return nil, err
	}

	var respVo *vo.GetDepartmentSimpleListV2RespVo = nil
	var err1 error = nil
	var hasMore = true

	var pageSize = 100
	var pageToken = ""

	depList := make([]vo.DepartmentRestInfoVo, 0)
	for ;;{
		if !hasMore{
			break
		}
		respVo, err1 = tenant.GetDepartmentSimpleListV2(depId, pageToken, pageSize, true)
		if err1 != nil{
			log.Error(err1)
			return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
		}
		if respVo.Code != 0{
			log.Error(respVo.Msg)
			return nil, errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, respVo.Msg)
		}

		hasMore = respVo.Data.HasMore
		pageToken = respVo.Data.PageToken

		depList = append(depList, respVo.Data.DepartmentInfos...)
	}

	return depList, nil
}


func GetScopeDeps(tenantKey string) ([]vo.DepartmentRestInfoVo, errs.SystemErrorInfo){
	tenant, err := GetTenant(tenantKey)
	if err != nil{
		log.Error(err)
		return nil, err
	}
	resp, scopeErr := tenant.GetScopeV2()
	if scopeErr != nil{
		log.Error(scopeErr)
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}
	if resp.Code != 0{
		log.Error(resp.Msg)
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}

	depList := make([]vo.DepartmentRestInfoVo, 0)

	deps := resp.Data.AuthedDepartments
	if deps != nil && len(deps) > 0{
		if deps[0] == "0"{
			depList, err = GetDeptList(tenantKey)
			if err != nil{
				log.Error(err)
				return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}
			if resp.Code != 0{
				log.Error(resp.Msg)
				return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}
		}else{
			for _, dep := range deps{
				depDetails, err := GetDeptListWithDepId(tenantKey, dep)
				if err != nil{
					log.Error(err)
					return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
				}
				if resp.Code != 0{
					log.Error(resp.Msg)
					return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
				}
				if depDetails != nil && len(depDetails) > 0{
					for _, depDetail := range depDetails{
						depList = append(depList, vo.DepartmentRestInfoVo{
							Id:  depDetail.Id,
							Name: depDetail.Name,
							ParentId: depDetail.ParentId,
						})
					}
				}
			}
		}
	}

	//做去重
	depMap := map[string]vo.DepartmentRestInfoVo{}
	for _, dep := range depList{
		depMap[dep.Id] = dep
	}
	depList = make([]vo.DepartmentRestInfoVo, 0)
	for _, dep := range depMap{
		depList = append(depList, dep)
	}

	//增加根部门
	//depList = append(depList, vo.DepartmentRestInfoVo{
	//	Id: "0",
	//	Name: "飞书平台组织",
	//	ParentId: "Root_Department_Identification",
	//})
	return depList, nil
}

func GetDeptUserInfosByDeptIds(tenantKey string, deptIds []string) ([]vo.UserDetailInfoV2, errs.SystemErrorInfo){
	tenant, err := GetTenant(tenantKey)
	if err != nil{
		log.Error(err)
		return nil, err
	}
	userInfos := make([]vo.UserDetailInfoV2, 0)
	if deptIds == nil || len(deptIds) == 0{
		return userInfos, nil
	}

	for _, outDepId := range deptIds{
		var respVo *vo.GetUserBatchGetV2Resp = nil
		var err1 error = nil
		var hasMore = true

		var pageSize = 100
		var pageToken = ""

		for ;; {
			if !hasMore {
				break
			}

			respVo, err1 = tenant.GetDepartmentUserDetailListV2(outDepId, pageToken, pageSize, true)
			if err1 != nil{
				log.Error(err1)
				return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
			}
			if respVo.Code != 0{
				log.Error(respVo.Msg)
				return nil, errs.BuildSystemErrorInfoWithMessage(errs.FeiShuOpenApiCallError, respVo.Msg)
			}
			hasMore = respVo.Data.HasMore
			pageToken = respVo.Data.PageToken

			userList := respVo.Data.Users

			if userList != nil && len(userList) > 0{
				userInfos = append(userInfos, userList...)
			}
		}
	}
	return userInfos, nil
}

func GetScopeOpenIds(tenantKey string) ([]string, errs.SystemErrorInfo){
	tenant, err := GetTenant(tenantKey)
	if err != nil{
		log.Error(err)
		return nil, err
	}

	resp, scopeErr := tenant.GetScopeV2()
	if scopeErr != nil{
		log.Error(scopeErr)
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}
	if resp.Code != 0{
		log.Error(resp.Msg)
		return nil, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError)
	}

	//获取当前授权范围内的所有用户
	scopeUsers := resp.Data.AuthedUsers
	scopeOpenIds := make([]string, 0)
	if scopeUsers != nil && len(scopeUsers) > 0{
		for _, scopeUser := range scopeUsers{
			scopeOpenIds = append(scopeOpenIds, scopeUser.OpenId)
		}
	}

	//获取当前授权范围内所有的部门
	scopeDeptIds := resp.Data.AuthedDepartments
	fsUserInfos, err := GetDeptUserInfosByDeptIds(tenantKey, scopeDeptIds)
	if err != nil{
		log.Error(err)
		return nil, err
	}
	for _, fsUserInfo := range fsUserInfos{
		scopeOpenIds = append(scopeOpenIds, fsUserInfo.OpenId)
	}

	scopeOpenIds = slice.SliceUniqueString(scopeOpenIds)

	return scopeOpenIds, nil
}