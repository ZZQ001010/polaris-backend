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
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/projectsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func ClearProjectObjectTypeList(orgId int64, projectId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheProjectObjectTypeList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	_, err := cache.Del(key)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func GetProjectObjectTypeList(orgId int64, projectId int64) (*[]bo.ProjectObjectTypeBo, errs.SystemErrorInfo) {
	return ProjectObjectTypesWithProjectByOrder(orgId, projectId, "")
}

func ProjectObjectTypesWithProjectByOrder(orgId int64, projectId int64, order string) (*[]bo.ProjectObjectTypeBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProjectObjectTypeList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	projectObjectTypeListPo := &[]po.PpmPrsProjectObjectType{}
	projectObjectTypeListBo := &[]bo.ProjectObjectTypeBo{}
	projectObjectTypeListJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if projectObjectTypeListJson != "" {
		err = json.FromJson(projectObjectTypeListJson, projectObjectTypeListBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return projectObjectTypeListBo, nil
	} else {
		if orgId == 0 {
			err = mysql.SelectAllByCondWithNumAndOrder(consts.TableProjectObjectType, db.Cond{
				consts.TcOrgId:    orgId,
				consts.TcIsDelete: consts.AppIsNoDelete,
				consts.TcStatus:   consts.AppStatusEnable,
			}, nil, 0, 0, order, projectObjectTypeListPo)
			if err != nil {
				return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
		} else {
			projectObjectTypeProcesses := &[]po.PpmPrsProjectObjectTypeProcess{}
			err := mysql.SelectAllByCond(consts.TableProjectObjectTypeProcess, db.Cond{
				consts.TcOrgId:     orgId,
				consts.TcProjectId: projectId,
				consts.TcIsDelete:  consts.AppIsNoDelete,
			}, projectObjectTypeProcesses)
			if err != nil {
				return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
			projectObjectTypeIds := make([]int64, len(*projectObjectTypeProcesses))
			for i, p := range *projectObjectTypeProcesses {
				projectObjectTypeIds[i] = p.ProjectObjectTypeId
			}
			err = mysql.SelectAllByCondWithNumAndOrder(consts.TableProjectObjectType, db.Cond{
				consts.TcOrgId:    orgId,
				consts.TcId:       db.In(projectObjectTypeIds),
				consts.TcIsDelete: consts.AppIsNoDelete,
				consts.TcStatus:   consts.AppStatusEnable,
			}, nil, 0, 0, order, projectObjectTypeListPo)
			if err != nil {
				return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
		}
		_ = copyer.Copy(projectObjectTypeListPo, projectObjectTypeListBo)

		projectObjectTypeListJson, err := json.ToJson(projectObjectTypeListBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, projectObjectTypeListJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return projectObjectTypeListBo, nil
	}
}

func GetProjectObjectTypeByLangCodeAndObjectType(orgId int64, projectId int64, langCode string, typ int) (*bo.ProjectObjectTypeBo, errs.SystemErrorInfo) {
	list, err := GetProjectObjectTypeList(orgId, projectId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	for _, objectType := range *list {
		if objectType.LangCode == langCode && objectType.ObjectType == typ {
			return &objectType, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.ProjectObjectTypeNotExist)
}

func GetProjectObjectTypeById(orgId int64, projectId int64, projectObjectTypeId int64) (*bo.ProjectObjectTypeBo, errs.SystemErrorInfo) {
	list, err := GetProjectObjectTypeList(orgId, projectId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	for _, objectType := range *list {
		if objectType.Id == projectObjectTypeId {
			return &objectType, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.ProjectObjectTypeNotExist)
}

func GetProjectObjectTypeOfIteration(orgId int64, projectId int64) (*bo.ProjectObjectTypeBo, errs.SystemErrorInfo) {
	objectType, err := GetProjectObjectTypeByLangCodeAndObjectType(orgId, projectId, consts.ProjectObjectTypeLangCodeIteration, consts.ProjectObjectTypeIteration)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return objectType, nil
}
