package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetProjectObjectTypeBoList(page uint, size uint, cond db.Cond) (*[]bo.ProjectObjectTypeBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectProjectObjectTypeByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.ProjectObjectTypeBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetProjectObjectTypeBo(id int64) (*bo.ProjectObjectTypeBo, errs.SystemErrorInfo) {
	po, err := dao.SelectProjectObjectTypeById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.ProjectObjectTypeBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func CreateProjectObjectType(bo *bo.ProjectObjectTypeBo) errs.SystemErrorInfo {

	pbo := bo.PpmPrsProjectObjectTypeProcessBo

	ppo := &po.PpmPrsProjectObjectTypeProcess{}

	err := copyer.Copy(pbo, ppo)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}

	po := &po.PpmPrsProjectObjectType{}
	err = copyer.Copy(bo, po)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}

	//任务和任务明细需要事务
	err2 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		err = mysql.TransInsert(tx, po)
		if err != nil {
			log.Errorf(consts.Mysql_TransInsert_error_printf, err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		err = mysql.TransInsert(tx, ppo)
		if err != nil {
			log.Errorf(consts.Mysql_TransInsert_error_printf, err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		return nil
	})

	if err2 != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}

	//
	//err2 := dao.InsertProjectObjectType(*po)
	//if err2 != nil {
	//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	//}
	return nil
}

func UpdateProjectObjectType(bo *bo.ProjectObjectTypeBo, projectId int64) errs.SystemErrorInfo {

	beforeProjectObjectType := &po.PpmPrsProjectObjectType{}

	afterProjectObjectType := &po.PpmPrsProjectObjectType{}

	projectObjectTypeProcesses := &[]po.PpmPrsProjectObjectTypeProcess{}
	err := mysql.SelectAllByCond(consts.TableProjectObjectTypeProcess, db.Cond{
		consts.TcOrgId:     bo.OrgId,
		consts.TcProjectId: projectId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}, projectObjectTypeProcesses)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	projectObjectTypeIds := make([]int64, len(*projectObjectTypeProcesses))
	for i, p := range *projectObjectTypeProcesses {
		projectObjectTypeIds[i] = p.ProjectObjectTypeId
	}

	projectObjectTypePo := &po.PpmPrsProjectObjectType{}
	copyErr := copyer.Copy(bo, projectObjectTypePo)
	if copyErr != nil {
		log.Error(copyErr)
		return errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	err2 := mysql.TransX(func(tx sqlbuilder.Tx) error {

		if bo.BeforeID != nil {

			err := mysql.SelectOneByCond(consts.TableProjectObjectType, db.Cond{
				consts.TcIsDelete: consts.AppIsNoDelete,
				consts.TcStatus:   consts.AppIsInitStatus,
				consts.TcId:       bo.BeforeID,
			}, beforeProjectObjectType)

			if err != nil {
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}

			//添加更亲其他的排序
			_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableProjectObjectType, db.Cond{
				consts.TcSort + " ": db.Lte(beforeProjectObjectType.Sort),
				consts.TcOrgId:      bo.OrgId,
				consts.TcId:         db.In(projectObjectTypeIds),
			}, mysql.Upd{
				consts.TcSort: db.Raw("sort -1"),
			})

			if err2 != nil {
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}

			projectObjectTypePo.Sort = beforeProjectObjectType.Sort

		}
		//更新后面的
		if bo.AfterID != nil {

			err := mysql.SelectOneByCond(consts.TableProjectObjectType, db.Cond{
				consts.TcIsDelete: consts.AppIsNoDelete,
				consts.TcStatus:   consts.AppIsInitStatus,
				consts.TcId:       bo.AfterID,
			}, afterProjectObjectType)

			if err != nil {
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}

			_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableProjectObjectType, db.Cond{
				consts.TcSort:       db.Gte(afterProjectObjectType.Sort),
				consts.TcOrgId:      bo.OrgId,
				consts.TcId:         db.In(projectObjectTypeIds),
			}, mysql.Upd{
				consts.TcSort: db.Raw("sort +1"),
			})

			if err2 != nil {
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
			}

			projectObjectTypePo.Sort = afterProjectObjectType.Sort
		}

		err := mysql.TransUpdate(tx, projectObjectTypePo)
		if err != nil {
			log.Errorf(consts.Mysql_TransInsert_error_printf, err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		return nil
	})

	if err2 != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}
	return nil
}

//projectObjectTypeId 有的时候代表是更新
func CheckSameProjectObjectTypeName(orgId, projectId int64, name string, projectObjectTypeId *int64) (map[int64]interface{}, errs.SystemErrorInfo) {
	//从缓存中获取所有的project_object_type
	bos, err := ProjectObjectTypesWithProjectByOrder(orgId, projectId, "")

	if err != nil {
		log.Error(err)
		return nil, err
	}

	projectObjectTypeMap := make(map[int64]interface{}, len(*bos))

	for _, value := range *bos {
		//如果是更新的情况 不等于本身 相同的名字返回错误
		if projectObjectTypeId != nil && name == value.Name && *projectObjectTypeId != value.Id {
			return nil, errs.BuildSystemErrorInfo(errs.ProjectObjectTypeSameName)
		}

		//走这里对比就是新增的情况如果名字相等)那么就返回相同
		if projectObjectTypeId == nil && name == value.Name {
			return nil, errs.BuildSystemErrorInfo(errs.ProjectObjectTypeSameName)
		}

		projectObjectTypeMap[value.Id] = value
	}

	return projectObjectTypeMap, nil
}

func JudgeLastProjectObjectType(orgId, projectId, projectObjectTpyeId int64) errs.SystemErrorInfo {

	count, err := mysql.SelectCountByCond(consts.TableProjectObjectTypeProcess, db.Cond{
		consts.TcOrgId:               orgId,
		consts.TcProjectId:           projectId,
		consts.TcIsDelete:            consts.AppIsNoDelete,
		consts.TcProjectObjectTypeId: db.NotEq(projectObjectTpyeId),
	})
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	if count < 1 {
		return errs.BuildSystemErrorInfo(errs.LastProjectObjectType)
	}
	return nil
}

func DeleteProjectObjectType(bo *bo.ProjectObjectTypeBo, projectId, projectObjectTpyeId, operatorId int64) errs.SystemErrorInfo {

	process := &po.PpmPrsProjectObjectTypeProcess{}
	err1 := mysql.SelectOneByCond(consts.TableProjectObjectTypeProcess, db.Cond{
		consts.TcOrgId:               bo.OrgId,
		consts.TcProjectId:           projectId,
		consts.TcIsDelete:            consts.AppIsNoDelete,
		consts.TcProjectObjectTypeId: projectObjectTpyeId,
	}, process)

	if err1 != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}

	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除时候删除对应的关联
		err := mysql.UpdateSmart(consts.TableProjectObjectTypeProcess, process.Id, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})

		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		//删除更新projectObjectType表
		_, err = dao.DeleteProjectObjectTypeById(bo.Id, operatorId, tx)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		return nil
	})
	return nil
}
