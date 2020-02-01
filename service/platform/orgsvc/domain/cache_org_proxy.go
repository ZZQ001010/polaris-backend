package domain

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/orgsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetOrgIdByOutOrgId(sourceChannel string, outOrgId string) (int64, errs.SystemErrorInfo){
	key, err5 := util.ParseCacheKey(sconsts.CacheOutOrgIdRelationId, map[string]interface{}{
		consts.CacheKeyOutOrgIdConstName:      outOrgId,
		consts.CacheKeySourceChannelConstName: sourceChannel,
	})
	if err5 != nil {
		log.Error(err5)
		return 0, err5
	}
	orgIdInfoJson, err := cache.Get(key)
	if err != nil {
		return 0, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if orgIdInfoJson != ""{
		orgIdInfo := &bo.OrgIdInfo{}
		err := json.FromJson(orgIdInfoJson, orgIdInfo)
		if err != nil {
			return 0, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return orgIdInfo.OrgId, nil
	}else{
		orgBo, err := GetOrgByOutOrgId(sourceChannel, outOrgId)
		if err != nil{
			log.Error(err)
			return 0, err
		}
		orgIdInfo := bo.OrgIdInfo{
			OutOrgId: outOrgId,
			OrgId: orgBo.Id,
		}
		orgIdInfoJson = json.ToJsonIgnoreError(orgIdInfo)
		err1 := cache.SetEx(key, orgIdInfoJson, consts.GetCacheBaseExpire())
		if err1 != nil {
			return 0, errs.BuildSystemErrorInfo(errs.RedisOperateError, err1)
		}
		return orgIdInfo.OrgId, nil
	}
}

func GetBaseOrgOutInfo(sourceChannel string, orgId int64) (*bo.BaseOrgOutInfoBo, errs.SystemErrorInfo){
	key, err5 := util.ParseCacheKey(sconsts.CacheBaseOrgOutInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:      orgId,
		consts.CacheKeySourceChannelConstName: sourceChannel,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}
	outOrgInfoJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	if outOrgInfoJson != ""{
		orgOutInfoBo := &bo.BaseOrgOutInfoBo{}
		err := json.FromJson(outOrgInfoJson, orgOutInfoBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return orgOutInfoBo, nil
	}else{
		orgOutInfo := &po.PpmOrgOrganizationOutInfo{}
		err = mysql.SelectOneByCond(orgOutInfo.TableName(), db.Cond{
			consts.TcOrgId:         orgId,
			consts.TcSourceChannel: sourceChannel,
			consts.TcIsDelete:      consts.AppIsNoDelete,
		}, orgOutInfo)
		if err != nil {
			log.Info(strs.ObjectToString(err))
			return nil, errs.OrgOutInfoNotExist
		}
		orgOutInfoBo := &bo.BaseOrgOutInfoBo{
			OrgId: orgId,
			OutOrgId: orgOutInfo.OutOrgId,
			SourceChannel: sourceChannel,
		}
		outOrgInfoJson = json.ToJsonIgnoreError(orgOutInfoBo)
		err1 := cache.SetEx(key, outOrgInfoJson, consts.GetCacheBaseExpire())
		if err1 != nil {
			log.Error(err1)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err1)
		}
		return orgOutInfoBo, nil
	}
}

func GetBaseOrgInfoByOutOrgId(sourceChannel string, outOrgId string) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo){
	orgId, err := GetOrgIdByOutOrgId(sourceChannel, outOrgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return GetBaseOrgInfo(sourceChannel, orgId)
}

func ClearCacheBaseOrgInfo(sourceChannel string, orgId int64) errs.SystemErrorInfo {
	key, err5 := util.ParseCacheKey(sconsts.CacheBaseOrgInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:         orgId,
		consts.CacheKeySourceChannelConstName: sourceChannel,
	})

	if err5 != nil {
		log.Error(err5)
		return err5
	}

	_, err := cache.Del(key)

	if err != nil {
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}

func GetBaseOrgInfo(sourceChannel string, orgId int64) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheBaseOrgInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:         orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	baseOrgInfoJson, err := cache.Get(key)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	baseOrgInfo := &bo.BaseOrgInfoBo{}
	if baseOrgInfoJson != "" {
		err := json.FromJson(baseOrgInfoJson, baseOrgInfo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}

	} else {
		orgInfo := &po.PpmOrgOrganization{}
		err := mysql.SelectOneByCond(orgInfo.TableName(), db.Cond{
			consts.TcId:       orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, orgInfo)
		if err != nil {
			log.Info(strs.ObjectToString(err))
			return nil, errs.BuildSystemErrorInfo(errs.OrgNotExist)
		}

		baseOrgInfo.OrgId = orgId
		baseOrgInfo.OrgName = orgInfo.Name
		baseOrgInfo.OrgOwnerId = orgInfo.Owner
		baseOrgInfoJson, err = json.ToJson(baseOrgInfo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, baseOrgInfoJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
	}
	if sourceChannel != ""{
		baseOrgOutInfo, err := GetBaseOrgOutInfo(sourceChannel, orgId)
		if err != nil{
			log.Error(err)
		}else{
			baseOrgInfo.OutOrgId = baseOrgOutInfo.OutOrgId
			baseOrgInfo.SourceChannel = baseOrgOutInfo.SourceChannel
			baseOrgInfo.OrgId = orgId
		}
	}

	return baseOrgInfo, nil
}

func GetOutDeptAndInnerDept(orgId int64, tx *sqlbuilder.Tx) (map[string]int64, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheDeptRelation, map[string]interface{}{
		consts.CacheKeyOrgIdConstName: orgId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	deptRelationListJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	deptRelationList := &map[string]int64{}
	if deptRelationListJson != "" {
		err = json.FromJson(deptRelationListJson, deptRelationList)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return *deptRelationList, nil
	} else {
		deptOurInfoList := &[]po.PpmOrgDepartmentOutInfo{}

		selectErr := queryDepartmentOutInfWithTx(tx, deptOurInfoList, orgId)

		log.Info("部门关联关系: " + strs.ObjectToString(deptOurInfoList))
		if selectErr != nil {
			return *deptRelationList, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
		}
		for _, v := range *deptOurInfoList {
			(*deptRelationList)[v.OutOrgDepartmentId] = v.DepartmentId
		}
		deptRelationListJson, err := json.ToJson(deptRelationList)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, deptRelationListJson, consts.GetCacheBaseExpire())
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}

		return *deptRelationList, nil
	}
}

func queryDepartmentOutInfWithTx(tx *sqlbuilder.Tx, deptOurInfoList *[]po.PpmOrgDepartmentOutInfo, orgId int64) error {
	var selectErr error
	if tx != nil {
		//TODO 未定义TransSelectAllByCond，先使用SelectAllByCond(不使用事务不会有问题)
		selectErr = mysql.TransSelectAllByCond(*tx, consts.TableDepartmentOutInfo, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, deptOurInfoList)
	} else {
		selectErr = mysql.SelectAllByCond(consts.TableDepartmentOutInfo, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, deptOurInfoList)
	}
	return selectErr
}
