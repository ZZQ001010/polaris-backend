package msgfacade

import (
	"errors"
	"fmt"
	
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/http"
	
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/msgvo"
)


func InsertMqConsumeFailMsg(req msgvo.InsertMqConsumeFailMsgReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/msgsvc/insertMqConsumeFailMsg", config.GetPreUrl("msgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["msgType"] = req.MsgType
	queryParams["orgId"] = req.OrgId
	requestBody := json.ToJsonIgnoreError(req.Msg)
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("msgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func PushMsgToMq(req msgvo.PushMsgToMqReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/msgsvc/pushMsgToMq", config.GetPreUrl("msgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["msgType"] = req.MsgType
	queryParams["orgId"] = req.OrgId
	requestBody := json.ToJsonIgnoreError(req.Msg)
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("msgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func SendLoginSMS(req msgvo.SendLoginSMSReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/msgsvc/sendLoginSMS", config.GetPreUrl("msgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["phoneNumber"] = req.PhoneNumber
	queryParams["code"] = req.Code
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("msgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func SendMail(req msgvo.SendMailReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/msgsvc/sendMail", config.GetPreUrl("msgsvc"))
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("msgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func SendSMS(req msgvo.SendSMSReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/msgsvc/sendSMS", config.GetPreUrl("msgsvc"))
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("msgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UpdateMsgStatus(req msgvo.UpdateMsgStatusReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/msgsvc/updateMsgStatus", config.GetPreUrl("msgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["msgId"] = req.MsgId
	queryParams["newStatus"] = req.NewStatus
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("msgsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


