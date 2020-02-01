package rolefacade

import (
	"errors"
	"fmt"
	
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/http"
	
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
)


func GetOrgAdminUser(req rolevo.GetOrgAdminUserReqVo) rolevo.GetOrgAdminUserRespVo {
	respVo := &rolevo.GetOrgAdminUserRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getOrgAdminUser", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	respBody, respStatusCode, err := http.Get(reqUrl, queryParams)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("rolesvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetOrgRoleList(req rolevo.GetOrgRoleListReqVo) rolevo.GetOrgRoleListRespVo {
	respVo := &rolevo.GetOrgRoleListRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getOrgRoleList", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	respBody, respStatusCode, err := http.Get(reqUrl, queryParams)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("rolesvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetOrgRoleUser(req rolevo.GetOrgRoleUserReqVo) rolevo.GetOrgRoleUserRespVo {
	respVo := &rolevo.GetOrgRoleUserRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getOrgRoleUser", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	respBody, respStatusCode, err := http.Get(reqUrl, queryParams)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("rolesvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetPersonalPermissionInfo(req rolevo.GetPersonalPermissionInfoReqVo) rolevo.GetPersonalPermissionInfoRespVo {
	respVo := &rolevo.GetPersonalPermissionInfoRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getPersonalPermissionInfo", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["projectId"] = req.ProjectId
	queryParams["issueId"] = req.IssueId
	queryParams["sourceChannel"] = req.SourceChannel
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	respBody, respStatusCode, err := http.Get(reqUrl, queryParams)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("rolesvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetProjectRoleList(req rolevo.GetProjectRoleListReqVo) rolevo.GetOrgRoleListRespVo {
	respVo := &rolevo.GetOrgRoleListRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getProjectRoleList", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	respBody, respStatusCode, err := http.Get(reqUrl, queryParams)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("rolesvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetUserAdminFlag(req rolevo.GetUserAdminFlagReqVo) rolevo.GetUserAdminFlagRespVo {
	respVo := &rolevo.GetUserAdminFlagRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/rolesvc/getUserAdminFlag", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	respBody, respStatusCode, err := http.Get(reqUrl, queryParams)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("rolesvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func PermissionOperationList(req rolevo.PermissionOperationListReqVo) rolevo.PermissionOperationListRespVo {
	respVo := &rolevo.PermissionOperationListRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/rolesvc/permissionOperationList", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["roleId"] = req.RoleId
	queryParams["projectId"] = req.ProjectId
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	respBody, respStatusCode, err := http.Get(reqUrl, queryParams)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("rolesvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UpdateUserOrgRole(req rolevo.UpdateUserOrgRoleReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/rolesvc/updateUserOrgRole", config.GetPreUrl("rolesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["userId"] = req.UserId
	queryParams["roleId"] = req.RoleId
	queryParams["projectId"] = req.ProjectId
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	respBody, respStatusCode, err := http.Get(reqUrl, queryParams)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("rolesvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


