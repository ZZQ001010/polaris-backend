package trendsfacade

import (
	"errors"
	"fmt"
	
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/http"
	
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/trendsvo"
)


func AddIssueTrends(req trendsvo.AddIssueTrendsReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/trendssvc/addIssueTrends", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	queryParams["key"] = req.Key
	requestBody := json.ToJsonIgnoreError(req.IssueTrendsBo)
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("trendssvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func AddOrgTrends(req trendsvo.AddOrgTrendsReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/trendssvc/addOrgTrends", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.OrgTrendsBo)
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("trendssvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func AddProjectTrends(req trendsvo.AddProjectTrendsReqVo) vo.VoidErr {
	respVo := &vo.VoidErr{}
	
	reqUrl := fmt.Sprintf("%s/api/trendssvc/addProjectTrends", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.ProjectTrendsBo)
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("trendssvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func CreateComment(req trendsvo.CreateCommentReqVo) trendsvo.CreateCommentRespVo {
	respVo := &trendsvo.CreateCommentRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/trendssvc/createComment", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.CommentBo)
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("trendssvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func CreateTrends(req trendsvo.CreateTrendsReqVo) trendsvo.CreateTrendsRespVo {
	respVo := &trendsvo.CreateTrendsRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/trendssvc/createTrends", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	requestBody := json.ToJsonIgnoreError(req.TrendsBo)
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("trendssvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func GetMQTTChannelKey(req trendsvo.GetMQTTChannelKeyReqVo) trendsvo.GetMQTTChannelKeyRespVo {
	respVo := &trendsvo.GetMQTTChannelKeyRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/trendssvc/getMQTTChannelKey", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("trendssvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func NoticeList(req trendsvo.NoticeListReqVo) trendsvo.NoticeListRespVo {
	respVo := &trendsvo.NoticeListRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/trendssvc/noticeList", config.GetPreUrl("trendssvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("trendssvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func TrendList(req trendsvo.TrendListReqVo) trendsvo.TrendListRespVo {
	respVo := &trendsvo.TrendListRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/trendssvc/trendList", config.GetPreUrl("trendssvc"))
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("trendssvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func UnreadNoticeCount(req trendsvo.UnreadNoticeCountReqVo) trendsvo.UnreadNoticeCountRespVo {
	respVo := &trendsvo.UnreadNoticeCountRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/trendssvc/unreadNoticeCount", config.GetPreUrl("trendssvc"))
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("trendssvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


