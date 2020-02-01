package orgfacade

import (
	"errors"
	"fmt"
	
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/http"
	
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
)


func BatchGetUserDetailInfo(req orgvo.BatchGetUserInfoReq) orgvo.BatchGetUserInfoResp {
	respVo := &orgvo.BatchGetUserInfoResp{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/batchGetUserDetailInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.UserIds)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func BindLoginName(req orgvo.BindLoginNameReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/bindLoginName", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func CheckLoginName(req orgvo.CheckLoginNameReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/checkLoginName", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func CreateOrg(req orgvo.CreateOrgReqVo) orgvo.CreateOrgRespVo {
	respVo := &orgvo.CreateOrgRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/createOrg", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := json.ToJsonIgnoreError(req.Data)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func DepartmentInit(req orgvo.DepartmentInitReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/departmentInit", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["corpId"] = req.CorpId
	queryParams["sourceChannel"] = req.SourceChannel
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func DepartmentMembers(req orgvo.DepartmentMembersReqVo) orgvo.DepartmentMembersRespVo {
	respVo := &orgvo.DepartmentMembersRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/departmentMembers", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := json.ToJsonIgnoreError(req.Params)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func Departments(req orgvo.DepartmentsReqVo) orgvo.DepartmentsRespVo {
	respVo := &orgvo.DepartmentsRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/departments", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["currentUserId"] = req.CurrentUserId
	queryParams["orgId"] = req.OrgId
	requestBody := json.ToJsonIgnoreError(req.Params)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func FeiShuAuth(req vo.FeiShuAuthReq) orgvo.FeiShuAuthRespVo {
	respVo := &orgvo.FeiShuAuthRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/feiShuAuth", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["code"] = req.Code
	queryParams["codeType"] = req.CodeType
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GeneralInitOrg(req orgvo.InitOrgReqVo) orgvo.OrgInitRespVo {
	respVo := &orgvo.OrgInitRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/generalInitOrg", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.InitOrg)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetBaseUserInfoBatch(req orgvo.GetBaseUserInfoBatchReqVo) orgvo.GetBaseUserInfoBatchRespVo {
	respVo := &orgvo.GetBaseUserInfoBatchRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getBaseUserInfoBatch", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["orgId"] = req.OrgId
	requestBody := json.ToJsonIgnoreError(req.UserIds)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetOrgIdListBySourceChannel(req orgvo.GetOrgIdListBySourceChannelReqVo) orgvo.GetOrgIdListBySourceChannelRespVo {
	respVo := &orgvo.GetOrgIdListBySourceChannelRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgIdListBySourceChannel", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetOrgUserInfoListBySourceChannel(req orgvo.GetOrgUserInfoListBySourceChannelReq) orgvo.GetOrgUserInfoListBySourceChannelResp {
	respVo := &orgvo.GetOrgUserInfoListBySourceChannelResp{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getOrgUserInfoListBySourceChannel", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetPwdLoginCode(req orgvo.GetPwdLoginCodeReqVo) orgvo.GetPwdLoginCodeRespVo {
	respVo := &orgvo.GetPwdLoginCodeRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getPwdLoginCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["captchaId"] = req.CaptchaId
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetUserIds(req orgvo.GetUserIdsReqVo) orgvo.GetUserIdsRespVo {
	respVo := &orgvo.GetUserIdsRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUserIds", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["sourceChannel"] = req.SourceChannel
	queryParams["corpId"] = req.CorpId
	queryParams["orgId"] = req.OrgId
	requestBody := json.ToJsonIgnoreError(req.EmpIdsBody)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetUserInfoByUserIds(req orgvo.GetUserInfoByUserIdsReqVo) orgvo.GetUserInfoByUserIdsListRespVo {
	respVo := &orgvo.GetUserInfoByUserIdsListRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/getUserInfoByUserIds", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := json.ToJsonIgnoreError(req.UserIds)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func InitOrg(req orgvo.InitOrgReqVo) orgvo.OrgInitRespVo {
	respVo := &orgvo.OrgInitRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/initOrg", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.InitOrg)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func OrgInit(req orgvo.OrgInitReqVo) orgvo.OrgInitRespVo {
	respVo := &orgvo.OrgInitRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/orgInit", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["corpId"] = req.CorpId
	queryParams["permanentCode"] = req.PermanentCode
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func OrgOwnerInit(req orgvo.OrgOwnerInitReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/orgOwnerInit", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["owner"] = req.Owner
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func OrgSysConfigInit(req orgvo.OrgVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/orgSysConfigInit", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func OrgUserList(req orgvo.OrgUserListReq) orgvo.OrgUserListResp {
	respVo := &orgvo.OrgUserListResp{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/orgUserList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func OrganizationInfo(req orgvo.OrganizationInfoReqVo) orgvo.OrganizationInfoRespVo {
	respVo := &orgvo.OrganizationInfoRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/organizationInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func RemoveOrgMember(req orgvo.RemoveOrgMemberReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/removeOrgMember", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func ResetPassword(req orgvo.ResetPasswordReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/resetPassword", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func RetrievePassword(req orgvo.RetrievePasswordReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/retrievePassword", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func ScheduleOrganizationPageList(req orgvo.ScheduleOrganizationPageListReqVo) orgvo.ScheduleOrganizationPageListRespVo {
	respVo := &orgvo.ScheduleOrganizationPageListRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/scheduleOrganizationPageList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func SendAuthCode(req orgvo.SendAuthCodeReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/sendAuthCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func SendSMSLoginCode(req orgvo.SendSMSLoginCodeReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/sendSMSLoginCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func SetPassword(req orgvo.SetPasswordReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/setPassword", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func SetPwdLoginCode(req orgvo.SetPwdLoginCodeReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/setPwdLoginCode", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["captchaId"] = req.CaptchaId
	queryParams["captchaPassword"] = req.CaptchaPassword
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func SwitchUserOrganization(req orgvo.SwitchUserOrganizationReqVo) orgvo.SwitchUserOrganizationRespVo {
	respVo := &orgvo.SwitchUserOrganizationRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/switchUserOrganization", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["token"] = req.Token
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func TeamInit(req orgvo.OrgVo) orgvo.TeamInitRespVo {
	respVo := &orgvo.TeamInitRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/teamInit", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func TeamOwnerInit(req orgvo.TeamOwnerInitReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/teamOwnerInit", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["teamId"] = req.TeamId
	queryParams["owner"] = req.Owner
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func TeamUserInit(req orgvo.TeamUserInitReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/teamUserInit", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["teamId"] = req.TeamId
	queryParams["userId"] = req.UserId
	queryParams["isRoot"] = req.IsRoot
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UnbindLoginName(req orgvo.UnbindLoginNameReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/unbindLoginName", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UpdateOrgMemberCheckStatus(req orgvo.UpdateOrgMemberCheckStatusReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateOrgMemberCheckStatus", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UpdateOrgMemberStatus(req orgvo.UpdateOrgMemberStatusReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateOrgMemberStatus", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UpdateOrganizationSetting(req orgvo.UpdateOrganizationSettingReqVo) orgvo.UpdateOrganizationSettingRespVo {
	respVo := &orgvo.UpdateOrganizationSettingRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateOrganizationSetting", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UpdateUserConfig(req orgvo.UpdateUserConfigReqVo) orgvo.UpdateUserConfigRespVo {
	respVo := &orgvo.UpdateUserConfigRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateUserConfig", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := json.ToJsonIgnoreError(req.UpdateUserConfigReq)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UpdateUserDefaultProjectIdConfig(req orgvo.UpdateUserDefaultProjectIdConfigReqVo) orgvo.UpdateUserConfigRespVo {
	respVo := &orgvo.UpdateUserConfigRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateUserDefaultProjectIdConfig", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := json.ToJsonIgnoreError(req.UpdateUserDefaultProjectIdConfigReq)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UpdateUserInfo(req orgvo.UpdateUserInfoReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateUserInfo", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := json.ToJsonIgnoreError(req.UpdateUserInfoReq)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UpdateUserPcConfig(req orgvo.UpdateUserPcConfigReqVo) orgvo.UpdateUserConfigRespVo {
	respVo := &orgvo.UpdateUserConfigRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/updateUserPcConfig", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := json.ToJsonIgnoreError(req.UpdateUserPcConfigReq)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UserInitByOrg(req orgvo.UserInitByOrgReqVo) orgvo.UserInitByOrgRespVo {
	respVo := &orgvo.UserInitByOrgRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/userInitByOrg", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["corpId"] = req.CorpId
	queryParams["orgId"] = req.OrgId
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UserLogin(req orgvo.UserLoginReqVo) orgvo.UserSMSLoginRespVo {
	respVo := &orgvo.UserSMSLoginRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/userLogin", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.UserLoginReq)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UserOrganizationList(req orgvo.UserOrganizationListReqVo) orgvo.UserOrganizationListRespVo {
	respVo := &orgvo.UserOrganizationListRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/userOrganizationList", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UserQuit(req orgvo.UserQuitReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/userQuit", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["token"] = req.Token
requestBody := ""
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UserRegister(req orgvo.UserRegisterReqVo) orgvo.UserRegisterRespVo {
	respVo := &orgvo.UserRegisterRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/userRegister", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func VerifyOrgUsers(req orgvo.VerifyOrgUsersReqVo) vo.BoolRespVo {
	respVo := &vo.BoolRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/orgsvc/verifyOrgUsers", config.GetPreUrl("orgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := json.ToJsonIgnoreError(req.Input)
fullUrl := reqUrl + http.ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody

	respBody, respStatusCode, err := http.Post(reqUrl, queryParams, requestBody)

	
	//Process the response
	if err != nil {
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, err))
		log.Errorf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299{
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("orgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


