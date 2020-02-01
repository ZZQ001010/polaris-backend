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

func GetIssueObjectTypeById(orgId, sourceId int64) (*bo.IssueObjectTypeBo, errs.SystemErrorInfo) {
	list, err1 := GetIssueObjectTypeList(orgId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}
	for _, source := range list {
		if source.Id == sourceId {
			return &source, nil
		}
	}
	return nil, errs.BuildSystemErrorInfo(errs.SourceNotExist)
}

func GetIssueObjectTypeListByProjectObjectTypeId(orgId, projectTypeId int64) ([]bo.IssueObjectTypeBo, errs.SystemErrorInfo) {
	list, err1 := GetIssueObjectTypeList(orgId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}
	result := make([]bo.IssueObjectTypeBo, 0)
	for _, source := range list {
		if source.ProjectObjectTypeId == projectTypeId {
			result = append(result, source)
		}
	}
	return result, nil
}

func GetIssueObjectTypeList(orgId int64) ([]bo.IssueObjectTypeBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheIssueObjectTypeList, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	issueObjectTypeListBo := &[]bo.IssueObjectTypeBo{}
	issueObjectTypeListPo := &[]po.PpmPrsIssueObjectType{}
	issueObjectTypeListJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if issueObjectTypeListJson != "" {

		err = json.FromJson(issueObjectTypeListJson, issueObjectTypeListBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return *issueObjectTypeListBo, nil
	} else {
		err := mysql.SelectAllByCond(consts.TableIssueObjectType, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcStatus:   consts.AppStatusEnable,
		}, issueObjectTypeListPo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}

		_ = copyer.Copy(issueObjectTypeListPo, issueObjectTypeListBo)

		issueObjectTypeListJson, err := json.ToJson(issueObjectTypeListBo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, issueObjectTypeListJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		return *issueObjectTypeListBo, nil
	}
}

func DeleteIssueObjectTypeListCache(orgId int64) errs.SystemErrorInfo {

	key, err := util.ParseCacheKey(sconsts.CacheIssueObjectTypeList, map[string]interface{}{
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

	return nil
}
