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

func GetProjectTypeProjectObjectTypeListByProjectType(orgId, projectTypeId int64) (*[]bo.ProjectTypeProjectObjectTypeBo, errs.SystemErrorInfo) {
	list, err := GetProjectTypeProjectObjectTypeList(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resultList := make([]bo.ProjectTypeProjectObjectTypeBo, 0, 10)
	for _, obj := range *list {
		if obj.ProjectTypeId == projectTypeId {
			resultList = append(resultList, obj)
		}
	}
	return &resultList, nil
}

func GetProjectTypeProjectObjectTypeList(orgId int64) (*[]bo.ProjectTypeProjectObjectTypeBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProjectTypeProjectObjectTypeList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	projectTypeProjectObjectTypeListJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	projectTypeProjectObjectTypeListPo := &[]po.PpmPrsProjectTypeProjectObjectType{}
	projectTypeProjectObjectTypeListBo := &[]bo.ProjectTypeProjectObjectTypeBo{}
	if projectTypeProjectObjectTypeListJson != "" {
		err = json.FromJson(projectTypeProjectObjectTypeListJson, projectTypeProjectObjectTypeListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return projectTypeProjectObjectTypeListBo, nil
	} else {
		err := mysql.SelectAllByCond(consts.TableProjectTypeProjectObjectType, db.Cond{
			consts.TcOrgId:    db.In([]int64{orgId, 0}),
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, projectTypeProjectObjectTypeListPo)
		_ = copyer.Copy(projectTypeProjectObjectTypeListPo, projectTypeProjectObjectTypeListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.ProjectTypeProjectObjectTypeNotExist)
		}
		processListJson, err := json.ToJson(projectTypeProjectObjectTypeListBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, processListJson, consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return projectTypeProjectObjectTypeListBo, nil
	}
}
