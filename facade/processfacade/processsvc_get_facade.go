package processfacade

import (
	"errors"
	"fmt"
	
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/http"
	
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/processvo"
)


func GetDefaultProcessStatusId(req processvo.GetDefaultProcessIdReqVo) processvo.GetDefaultProcessIdRespVo {
	respVo := &processvo.GetDefaultProcessIdRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/processsvc/getDefaultProcessStatusId", config.GetPreUrl("processsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["processId"] = req.ProcessId
	queryParams["category"] = req.Category
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("processsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetNextProcessStepStatusList(req processvo.GetNextProcessStepStatusListReqVo) processvo.GetNextProcessStepStatusListRespVo {
	respVo := &processvo.GetNextProcessStepStatusListRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/processsvc/getNextProcessStepStatusList", config.GetPreUrl("processsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["processId"] = req.ProcessId
	queryParams["startStatusId"] = req.StartStatusId
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("processsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetProcessById(req processvo.GetProcessByIdReqVo) processvo.GetProcessByIdRespVo {
	respVo := &processvo.GetProcessByIdRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/processsvc/getProcessById", config.GetPreUrl("processsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["id"] = req.Id
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("processsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetProcessInitStatusId(req processvo.GetProcessInitStatusIdReqVo) processvo.GetProcessInitStatusIdRespVo {
	respVo := &processvo.GetProcessInitStatusIdRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/processsvc/getProcessInitStatusId", config.GetPreUrl("processsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["projectId"] = req.ProjectId
	queryParams["projectObjectTypeId"] = req.ProjectObjectTypeId
	queryParams["category"] = req.Category
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("processsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetProcessStatus(req processvo.GetProcessStatusReqVo) processvo.GetProcessStatusRespVo {
	respVo := &processvo.GetProcessStatusRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/processsvc/getProcessStatus", config.GetPreUrl("processsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["id"] = req.Id
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("processsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetProcessStatusByCategory(req processvo.GetProcessStatusByCategoryReqVo) processvo.GetProcessStatusByCategoryRespVo {
	respVo := &processvo.GetProcessStatusByCategoryRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/processsvc/getProcessStatusByCategory", config.GetPreUrl("processsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["statusId"] = req.StatusId
	queryParams["category"] = req.Category
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("processsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetProcessStatusIds(req processvo.GetProcessStatusIdsReqVo) processvo.GetProcessStatusIdsRespVo {
	respVo := &processvo.GetProcessStatusIdsRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/processsvc/getProcessStatusIds", config.GetPreUrl("processsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["category"] = req.Category
	queryParams["typ"] = req.Typ
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("processsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetProcessStatusList(req processvo.GetProcessStatusListReqVo) processvo.GetProcessStatusListRespVo {
	respVo := &processvo.GetProcessStatusListRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/processsvc/getProcessStatusList", config.GetPreUrl("processsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["processId"] = req.ProcessId
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("processsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetProcessStatusListByCategory(req processvo.GetProcessStatusListByCategoryReqVo) processvo.GetProcessStatusListByCategoryRespVo {
	respVo := &processvo.GetProcessStatusListByCategoryRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/processsvc/getProcessStatusListByCategory", config.GetPreUrl("processsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["category"] = req.Category
	queryParams["statusType"] = req.StatusType
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("processsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


