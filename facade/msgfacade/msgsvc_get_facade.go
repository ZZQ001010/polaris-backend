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


func GetFailMsgList(req msgvo.GetFailMsgListReqVo) msgvo.GetFailMsgListRespVo {
	respVo := &msgvo.GetFailMsgListRespVo{}
	
	reqUrl := fmt.Sprintf("%s/api/msgsvc/getFailMsgList", config.GetPreUrl("msgsvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["msgType"] = req.MsgType
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
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


