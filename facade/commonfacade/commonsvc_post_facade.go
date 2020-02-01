package commonfacade

import (
	"errors"
	"fmt"
	
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/http"
	
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/commonvo"
)


func AreaInfo(req commonvo.AreaInfoReqVo) commonvo.AreaInfoRespVo {
	respVo := &commonvo.AreaInfoRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/commonsvc/areaInfo", config.GetPreUrl("commonsvc"))
	queryParams := map[string]interface{}{}
	queryParams["industryID"] = req.IndustryID
	queryParams["countryID"] = req.CountryID
	queryParams["provinceID"] = req.ProvinceID
	queryParams["cityID"] = req.CityID
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("commonsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


func AreaLinkageList(req commonvo.AreaLinkageListReqVo) commonvo.AreaLinkageListRespVo {
	respVo := &commonvo.AreaLinkageListRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/commonsvc/areaLinkageList", config.GetPreUrl("commonsvc"))
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("commonsvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


