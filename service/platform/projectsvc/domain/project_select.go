package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetProject(orgId int64, projectId int64) (*bo.ProjectBo, errs.SystemErrorInfo) {
	project := &po.PpmProProject{}
	err := mysql.SelectOneByCond(project.TableName(), db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       projectId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, project)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}
	projectBo := &bo.ProjectBo{}
	err1 := copyer.Copy(project, projectBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return projectBo, nil
}

func GetProjectBoList(orgId int64, ids []int64) ([]bo.ProjectBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmProProject{}
	err := mysql.SelectAllByCond(consts.TableProject, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(ids),
	}, pos)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bos := &[]bo.ProjectBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return *bos, nil
}

//通过项目类型langCode获取项目列表
func GetProjectBoListByProjectTypeLangCode(orgId int64, projectTypeLangCode *string) ([]bo.ProjectBo, errs.SystemErrorInfo) {
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if projectTypeLangCode != nil {
		projectType, err3 := GetProjectTypeByLangCode(orgId, *projectTypeLangCode)
		if err3 != nil {
			log.Error(err3)
			return nil, err3
		}
		cond[consts.TcProjectTypeId] = projectType.Id
	}

	pos := &[]po.PpmProProject{}
	err := mysql.SelectAllByCond(consts.TableProject, cond, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bos := &[]bo.ProjectBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return *bos, nil
}

func GetProjectInfoByOrgIds(orgIds []int64) ([]bo.ProjectBo, errs.SystemErrorInfo) {

	cond := db.Cond{
		consts.TcOrgId:    db.In(orgIds),
		consts.TcIsFiling: consts.AppIsNotFilling,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}

	pos := &[]po.PpmProProject{}
	err := mysql.SelectAllByCond(consts.TableProject, cond, pos)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bos := &[]bo.ProjectBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return *bos, nil

}
