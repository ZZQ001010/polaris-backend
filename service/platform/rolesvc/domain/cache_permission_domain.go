package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	bo2 "github.com/galaxy-book/polaris-backend/common/model/bo"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/rolesvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/po"
	"github.com/prometheus/common/log"
	"upper.io/db.v3"
)

func GetPermissionList() (*[]bo2.PermissionBo, errs.SystemErrorInfo) {
	key := sconsts.CachePermissionList
	listJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	bo := &[]bo2.PermissionBo{}
	if listJson != "" {
		err := json.FromJson(listJson, bo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	} else {
		po := &[]po.PpmRolPermission{}
		selectErr := mysql.SelectAllByCond(consts.TablePermission, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcStatus:   consts.AppStatusEnable,
		}, po)
		if selectErr != nil {
			log.Error(selectErr)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
		}
		_ = copyer.Copy(po, bo)
		listJson, err = json.ToJson(bo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, listJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
	}

	return bo, nil
}

func GetPermissionByType(permissionType int) ([]bo2.PermissionBo, errs.SystemErrorInfo) {
	res := []bo2.PermissionBo{}
	list, err := GetPermissionList()
	if err != nil {
		log.Error(err)
		return res, err
	}
	for _, v := range *list {
		if v.IsShow == consts.AppShowEnable && v.Type == permissionType {
			res = append(res, v)
		}
	}

	return res, nil
}

//获取项目权限项（只取部分）
func GetProjectPermission() ([]bo2.PermissionBo, errs.SystemErrorInfo) {
	res := []bo2.PermissionBo{}
	list, err := GetPermissionList()
	if err != nil {
		log.Error(err)
		return res, err
	}
	selectedPermission := []string{
		//任务
		consts.PermissionProIssue4,
		//任务栏
		consts.PermissionOrgProjectObjectType,
		//文件管理
		consts.PermissionProFile,
		//标签管理
		consts.PermissionProTag,
		//附件管理
		consts.PermissionProAttachment,
		//项目成员管理
		consts.PermissionProMember,
		//项目相关
		consts.PermissionProConfig,
		//项目角色管理
		consts.PermissionProRole,
	}
	for _, v := range *list {
		if v.IsShow == consts.AppShowEnable && v.Type == consts.PermissionTypePro {
			if ok, _ := slice.Contain(selectedPermission, v.LangCode); ok {
				res = append(res, v)
			}
		}
	}

	return res, nil
}

func GetPermissionById(permissionId int64) (bo2.PermissionBo, errs.SystemErrorInfo) {
	res := bo2.PermissionBo{}
	list, err := GetPermissionList()
	if err != nil {
		log.Error(err)
		return res, err
	}

	for _, bo := range *list {
		if bo.Id == permissionId {
			return bo, err
		}
	}

	return res, errs.PermissionNotExist
}
