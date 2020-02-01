package websitefacade

import (
	"errors"
	"fmt"
	
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/http"
	
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/websitevo"
)


func RegisterWebSiteContact(req websitevo.RegisterWebSiteContactReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/websitesvc/registerWebSiteContact", config.GetPreUrl("websitesvc"))
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
		respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.ServerError, errors.New(fmt.Sprintf("websitesvc response code %d", respStatusCode))))
		log.Errorf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return *respVo
	}
	jsonConvertErr := json.FromJson(respBody, respVo)
	if jsonConvertErr != nil{
	respVo.Err = vo.NewErr(errs.BuildSystemErrorInfo(errs.JSONConvertError, jsonConvertErr))
	}
	return *respVo
}


