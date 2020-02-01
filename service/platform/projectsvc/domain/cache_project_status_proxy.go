package domain

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/processvo"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/projectsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetProjectStatus(orgId, projectId int64) (*[]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProjectProcessStatus, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	processStatusJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	processStatusList := &[]bo.CacheProcessStatusBo{}
	if processStatusJson != "" {
		err = json.FromJson(processStatusJson, processStatusList)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return processStatusList, nil
	} else {
		projectRelationList := &[]po.PpmProProjectRelation{}
		err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcProjectId:    projectId,
			consts.TcStatus:       consts.AppStatusEnable,
			consts.TcRelationType: consts.IssueRelationTypeStatus,
		}, projectRelationList)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}

		if len(*projectRelationList) == 0 {
			return nil, errs.BuildSystemErrorInfo(errs.ProcessStatusNotExist)
		}
		//for _, v := range *projectRelationList {
		//	respVo := processfacade.GetProcessStatus(processvo.GetProcessStatusReqVo{OrgId: orgId, Id: v.Id})
		//	if respVo.Successful() {
		//		s := respVo.CacheProcessStatusBo
		//		*processStatusList = append(*processStatusList, *s)
		//	}
		//}
		assemblyProcessStatusList(projectRelationList, processStatusList, orgId)

		processListJson, err := json.ToJson(processStatusList)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, processListJson, consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}

		return processStatusList, nil
	}
}

func assemblyProcessStatusList(projectRelationList *[]po.PpmProProjectRelation, processStatusList *[]bo.CacheProcessStatusBo, orgId int64) {
	for _, v := range *projectRelationList {
		respVo := processfacade.GetProcessStatus(processvo.GetProcessStatusReqVo{OrgId: orgId, Id: v.RelationId})
		if respVo.Successful() {
			s := respVo.CacheProcessStatusBo
			*processStatusList = append(*processStatusList, *s)
		}
	}
}
