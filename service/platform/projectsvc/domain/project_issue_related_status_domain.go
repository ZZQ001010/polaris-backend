package domain

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
)

func GetProjectRelatedStatus(orgId int64, projectId int64, projectObjectTypeId int64) ([]bo.HomeIssueStatusInfoBo, errs.SystemErrorInfo) {
	process, err1 := GetProjectProcessBo(orgId, projectId, projectObjectTypeId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err1)
	}

	cacheStatusBos, err1 := processfacade.GetProcessStatusListRelaxed(orgId, process.Id)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err1)
	}
	//err := mysql.SelectByQuery(consts.SelectProjectRelatedStatus, statusIdBos, orgId, projectId, projectObjectTypeId)
	//if err != nil{
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	//}
	//statusIds := make([]int64, len(*statusIdBos))
	//for i, statusIdBo := range *statusIdBos{
	//	statusIds[i] = statusIdBo.StatusId
	//}
	//cacheStatusList, err1 := proxies.GetProcessStatusListByIds(orgId, statusIds)
	//if err1 != nil{
	//	log.Error(err1)
	//	return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err1)
	//}
	statusInfoBos := make([]bo.HomeIssueStatusInfoBo, len(*cacheStatusBos))
	for i, cacheStatus := range *cacheStatusBos {
		statusInfoBos[i] = bo.HomeIssueStatusInfoBo{
			ID:        cacheStatus.StatusId,
			Name:      cacheStatus.Name,
			BgStyle:   cacheStatus.BgStyle,
			FontStyle: cacheStatus.FontStyle,
			Type:      cacheStatus.StatusType,
		}
	}
	return statusInfoBos, nil
}
