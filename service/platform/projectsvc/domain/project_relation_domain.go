package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"strconv"
	"upper.io/db.v3"
)

func GetProjectRelationBo(projectBo bo.ProjectBo, relationType int, relationId int64) (*bo.ProjectRelationBo, errs.SystemErrorInfo) {
	projectRelation, err := dao.SelectOneProjectRelation(db.Cond{
		consts.TcOrgId:        projectBo.OrgId,
		consts.TcProjectId:    projectBo.Id,
		consts.TcRelationId:   relationId,
		consts.TcRelationType: relationType,
		consts.TcStatus:       consts.AppStatusEnable,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotRelatedError)
	}

	bo := &bo.ProjectRelationBo{}
	err1 := copyer.Copy(projectRelation, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func GetProjectRelationByType(projectId int64, relationTypes []int64) (*[]bo.ProjectRelationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcProjectId:    projectId,
		consts.TcRelationType: db.In(relationTypes),
		consts.TcStatus:       consts.AppStatusEnable,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotRelatedError)
	}
	bos := &[]bo.ProjectRelationBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bos, nil
}

func GetProjectRelationByTypeAndUserId(userId int64, relationTypes []int64) (*[]bo.ProjectRelationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcRelationId:   userId,
		consts.TcRelationType: db.In(relationTypes),
		consts.TcStatus:       consts.AppStatusEnable,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotRelatedError)
	}
	bos := &[]bo.ProjectRelationBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bos, nil
}

//更新项目关联，多关联类型check
//relationTypes: 使用该类型范围内的数据与relationIds做过滤，得到真正需要关联的id
//targetRelationType: 目标关联类型，新增的关联以此类型为准
func UpdateProjectRelationWithRelationTypes(operatorId, orgId, projectId int64, relationTypes []int, targetRelationType int, relationIds []int64) ([]int64, errs.SystemErrorInfo) {
	//防止项目成员重复插入
	uid := uuid.NewUuid()
	projectIdStr := strconv.FormatInt(projectId, 10)
	lockKey := consts.AddProjectRelationLock + projectIdStr
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("获取%s锁时异常 %v", lockKey, err)
		return nil, errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	} else {
		return nil, errs.BuildSystemErrorInfo(errs.GetDistributedLockError)
	}

	//预先查询已有的关联
	projectRelations := &[]po.PpmProProjectRelation{}
	err5 := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcProjectId:    projectId,
		consts.TcRelationId:   db.In(relationIds),
		consts.TcRelationType: db.In(relationTypes),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, projectRelations)
	if err5 != nil {
		log.Error(err5)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	//check，去掉已有的关联
	if len(*projectRelations) > 0 {
		notRelationUserIds := make([]int64, 0)
		allExistIds := []int64{}
		for _, issueRelation := range *projectRelations {
			allExistIds = append(allExistIds, issueRelation.RelationId)
		}
		for _, id := range relationIds {
			exist, err := slice.Contain(allExistIds, id)
			if err != nil {
				log.Error(err)
				continue
			}
			if !exist {
				notRelationUserIds = append(notRelationUserIds, id)
			}
		}

		relationIds = notRelationUserIds
	}
	relationIds = slice.SliceUniqueInt64(relationIds)

	relationIdsSize := len(relationIds)
	if relationIdsSize == 0 {
		return relationIds, nil
	}

	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, relationIdsSize)
	if err != nil {
		log.Errorf("id generate: %q\n", err)
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}

	projectRelationPos := make([]po.PpmProProjectRelation, relationIdsSize)
	for i, relationId := range relationIds {
		id := ids.Ids[i].Id
		issueRelation := &po.PpmProProjectRelation{}
		issueRelation.Id = id
		issueRelation.OrgId = orgId
		issueRelation.ProjectId = projectId
		issueRelation.RelationId = relationId
		issueRelation.RelationType = targetRelationType
		issueRelation.Creator = operatorId
		issueRelation.Updator = operatorId
		issueRelation.IsDelete = consts.AppIsNoDelete
		projectRelationPos[i] = *issueRelation
	}

	err2 := mysql.BatchInsert(&po.PpmProProjectRelation{}, slice.ToSlice(projectRelationPos))
	if err2 != nil {
		log.Errorf("mysql.BatchInsert(): %q\n", err2)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}

	return relationIds, nil
}

//更新项目关联，带分布式锁
func UpdateProjectRelation(operatorId, orgId, projectId int64, relationType int, relationIds []int64) errs.SystemErrorInfo {
	_, err := UpdateProjectRelationWithRelationTypes(operatorId, orgId, projectId, []int{relationType}, relationType, relationIds)
	return err
}
