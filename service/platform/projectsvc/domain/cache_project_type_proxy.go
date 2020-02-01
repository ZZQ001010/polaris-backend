package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/projectsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetProjectTypeList(orgId int64) (*[]bo.ProjectTypeBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProjectTypeList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	projectTypeListPo := &[]po.PpmPrsProjectType{}
	projectTypeListBo := &[]bo.ProjectTypeBo{}
	projectTypeListJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if projectTypeListJson != "" {

		err = json.FromJson(projectTypeListJson, projectTypeListBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return projectTypeListBo, nil
	} else {
		err := mysql.SelectAllByCond(consts.TableProjectType, db.Cond{
			consts.TcOrgId:    db.In([]int64{orgId, 0}),
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcStatus:   consts.AppStatusEnable,
		}, projectTypeListPo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_ = copyer.Copy(projectTypeListPo, projectTypeListBo)
		projectTypeListJson, err := json.ToJson(projectTypeListBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, projectTypeListJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return projectTypeListBo, nil
	}
}

func GetProjectTypeByLangCode(orgId int64, langCode string) (*bo.ProjectTypeBo, errs.SystemErrorInfo) {
	list, err := GetProjectTypeList(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, projectType := range *list {
		if projectType.LangCode == langCode {
			return &projectType, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.ProjectTypeNotExist)
}

func GetProjectTypeById(orgId int64, id int64) (*bo.ProjectTypeBo, errs.SystemErrorInfo) {
	list, err := GetProjectTypeList(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, projectType := range *list {
		if projectType.Id == id {
			return &projectType, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.ProjectTypeNotExist)
}

func GetProjectProcessBo(orgId int64, projectId int64, projectObjectTypeId int64) (*bo.ProcessBo, errs.SystemErrorInfo) {
	processId, err := GetProjectProcessId(orgId, projectId, projectObjectTypeId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectObjectTypeProcessNotExist)
	}
	return processfacade.GetProcessByIdRelaxed(orgId, processId)
}

func GetProjectProcessId(orgId int64, projectId int64, projectObjectTypeId int64) (int64, errs.SystemErrorInfo) {
	projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{}
	err := mysql.SelectOneByCond(projectObjectTypeProcess.TableName(), db.Cond{
		consts.TcOrgId:               orgId,
		consts.TcProjectId:           projectId,
		consts.TcProjectObjectTypeId: projectObjectTypeId,
		consts.TcIsDelete:            consts.AppIsNoDelete,
	}, projectObjectTypeProcess)
	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.ProjectObjectTypeProcessNotExist)
	}
	return projectObjectTypeProcess.ProcessId, nil
}
