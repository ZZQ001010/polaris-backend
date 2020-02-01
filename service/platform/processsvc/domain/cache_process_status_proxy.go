package domain

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/processsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/po"
	"strconv"
	"upper.io/db.v3"
)

var log = logger.GetDefaultLogger()

func GetNextProcessStepStatusList(orgId, processId, startStatusId int64) (*[]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	nextStatusList := &[]bo.CacheProcessStatusBo{}
	processStepList, err := GetProcessStepList(orgId, processId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, processStep := range *processStepList {
		if processStep.StartStatus == startStatusId {
			nextStatusId := processStep.EndStatus
			status, err := GetProcessStatus(orgId, nextStatusId)
			status.DisplayName = processStep.Name
			if err != nil {
				log.Error(err)
				return nil, err
			}
			*nextStatusList = append(*nextStatusList, *status)
		}
	}
	return nextStatusList, nil
}

func GetProcessStatus(orgId int64, id int64) (*bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	list, err := GetStatusList(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, status := range *list {
		if status.StatusId == id {
			return &status, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.ProcessStatusNotExist)
}

func GetProcessStatusListByIds(orgId int64, ids []int64) ([]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	list, err := GetStatusList(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	processStatusBos := make([]bo.CacheProcessStatusBo, 0, len(*list))
	for _, status := range *list {
		for _, statusId := range ids {
			if status.StatusId == statusId {
				processStatusBos = append(processStatusBos, status)
				break
			}
		}
	}
	return processStatusBos, nil
}

func GetStatusList(orgId int64) (*[]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheStatusList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	processStatusListJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if processStatusListJson != "" {
		cacheProcessStatusList := &[]bo.CacheProcessStatusBo{}
		err = json.FromJson(processStatusListJson, cacheProcessStatusList)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return cacheProcessStatusList, nil
	} else {
		processStatusListPo := &[]po.PpmPrsProcessStatus{}
		err = mysql.SelectAllByCond(consts.TableProcessStatus, db.Cond{
			consts.TcOrgId:    db.In([]int64{orgId, 0}),
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcStatus:   consts.AppStatusEnable,
		}, processStatusListPo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		cacheProcessStatusList := make([]bo.CacheProcessStatusBo, len(*processStatusListPo))
		for i, status := range *processStatusListPo {
			cacheProcessStatusList[i] = bo.CacheProcessStatusBo{
				StatusId:   status.Id,
				StatusType: status.Type,
				Category:   status.Category,
				Name:       status.Name,
				BgStyle:    status.BgStyle,
				FontStyle:  status.FontStyle,
			}
		}
		processStatusListJson, err = json.ToJson(cacheProcessStatusList)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, processStatusListJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return &cacheProcessStatusList, nil
	}
}

func GetProcessStatusList(orgId, processId int64) (*[]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProcessStatusList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProcessIdConstName: processId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	processProcessStatusListPo := &[]po.PpmPrsProcessProcessStatus{}
	processStatusListJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if processStatusListJson != "" {
		cacheProcessStatusList := &[]bo.CacheProcessStatusBo{}
		err = json.FromJson(processStatusListJson, cacheProcessStatusList)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return cacheProcessStatusList, nil
	} else {
		err := mysql.SelectAllByCond(consts.TableProcessProcessStatus, db.Cond{
			consts.TcOrgId:     db.In([]int64{orgId, 0}),
			consts.TcProcessId: processId,
			consts.TcIsDelete:  consts.AppIsNoDelete,
		}, processProcessStatusListPo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.ProcessProcessStatusRelationError)
		}
		cacheProcessStatusList := make([]bo.CacheProcessStatusBo, len(*processProcessStatusListPo))
		statusError := dealProcessInitStatus(processProcessStatusListPo, orgId, &cacheProcessStatusList)

		if statusError != nil {

			return nil, statusError
		}

		processStatusListJson, err = json.ToJson(cacheProcessStatusList)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, processStatusListJson, consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return &cacheProcessStatusList, nil
	}
}

func dealProcessInitStatus(processProcessStatusList *[]po.PpmPrsProcessProcessStatus, orgId int64, cacheProcessStatusList *[]bo.CacheProcessStatusBo) errs.SystemErrorInfo {
	for i, processProcessStatus := range *processProcessStatusList {
		processStatus, err := GetProcessStatus(orgId, processProcessStatus.ProcessStatusId)
		if err != nil {
			log.Error(err)
			return err
		}
		processStatus.IsInit = processProcessStatus.IsInitStatus == 1
		(*cacheProcessStatusList)[i] = *processStatus
	}
	return nil
}

func GetProcessInitStatusId(orgId, projectId, projectObjectTypeId int64, category int) (int64, errs.SystemErrorInfo) {
	respVo := projectfacade.GetProjectProcessId(
		projectvo.GetProjectProcessIdReqVo{
			OrgId:               orgId,
			ProjectId:           projectId,
			ProjectObjectTypeId: projectObjectTypeId,
		})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return 0, respVo.Error()
	}
	processId := respVo.ProcessId

	log.Info("processId: " + strconv.FormatInt(processId, 10))

	process, err := GetProcess(orgId, processId)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	statusList, err := GetProcessStatusList(orgId, process.Id)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	for _, status := range *statusList {
		if status.Category == category {
			if status.IsInit {
				return status.StatusId, nil
			}
		}
	}
	return 0, errs.BuildSystemErrorInfo(errs.ProcessProcessStatusInitStatueNotExist)
}

//typ为-1表示匹配所有
func GetProcessStatusIds(orgId int64, category int, typ int) (*[]int64, errs.SystemErrorInfo) {
	statusList, err := GetStatusList(orgId)
	if err != nil {
		return nil, err
	}

	res := &[]int64{}
	for _, status := range *statusList {
		if status.Category == category {
			if typ == -1 || status.StatusType == typ {
				*res = append(*res, status.StatusId)
			}
		}
	}
	return res, nil
}

func GetProcessStatusByCategory(orgId int64, statusId int64, category int) (*bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	statusList, err := GetStatusList(orgId)
	if err != nil {
		return nil, err
	}
	for _, status := range *statusList {
		if status.StatusId == statusId && status.Category == category {
			return &status, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.ProcessStatusNotExist)
}

func GetProcessStatusListByCategory(orgId int64, category int) ([]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	statusList, err := GetStatusList(orgId)
	if err != nil {
		return nil, err
	}
	cacheProcessStatusBoList := make([]bo.CacheProcessStatusBo, 0, 10)
	for _, status := range *statusList {
		if status.Category == category {
			cacheProcessStatusBoList = append(cacheProcessStatusBoList, status)
		}
	}
	return cacheProcessStatusBoList, nil
}

func GetDefaultProcessStatusId(orgId int64, processId int64, category int) (int64, errs.SystemErrorInfo) {
	statusList, err := GetProcessStatusList(orgId, processId)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	for _, status := range *statusList {
		if status.Category == category {
			if status.IsInit {
				return status.StatusId, nil
			}
		}
	}
	return 0, errs.BuildSystemErrorInfo(errs.ProcessProcessStatusInitStatueNotExist)
}

//删除关联缓存
func DeleteCacheProcessStatusList(orgId, processId int64) errs.SystemErrorInfo {

	key, err := util.ParseCacheKey(sconsts.CacheStatusList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})

	if err != nil {
		log.Error(err)
		return err
	}

	_, err1 := cache.Del(key)

	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	//删除关联缓存
	key, err = util.ParseCacheKey(sconsts.CacheProcessStatusList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProcessIdConstName: processId,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	_, err1 = cache.Del(key)

	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	return nil
}
