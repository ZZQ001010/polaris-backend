package domain

import (
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/resourcefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"strconv"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func DeleteAllIssueRelation(tx sqlbuilder.Tx, operatorId, orgId, issueId int64) errs.SystemErrorInfo {
	//删除之前的关联
	_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIssueId:  issueId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operatorId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func DeleteAllIssueRelationByIds(operatorId, orgId int64, relationIds []int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	//删除之前的关联
	_, err := mysql.TransUpdateSmartWithCond(tx[0], consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       db.In(relationIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete:   consts.AppIsDeleted,
		consts.TcUpdator:    operatorId,
		consts.TcUpdateTime: time.Now(),
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func DeleteIssueRelation(operatorId int64, issueBo bo.IssueBo, relationType int) errs.SystemErrorInfo {
	orgId := issueBo.OrgId
	issueId := issueBo.Id
	//删除之前的关联
	_, err := mysql.UpdateSmartWithCond(consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      issueId,
		consts.TcRelationType: relationType,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operatorId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func DeleteIssueRelationByIds(operatorId int64, issueBo bo.IssueBo, relationType int, relationIds []int64) errs.SystemErrorInfo {
	if relationIds == nil || len(relationIds) == 0 {
		return nil
	}

	orgId := issueBo.OrgId
	issueId := issueBo.Id

	err := mysql.TransX(func(tx sqlbuilder.Tx) error{
		//删除之前的关联
		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcIssueId:      issueId,
			consts.TcRelationId:   db.In(relationIds),
			consts.TcRelationType: relationType,
			consts.TcIsDelete:     consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  operatorId,
		})
		if err != nil {
			log.Error(err)
			return err
		}

		//删除文件
		resp := resourcefacade.DeleteResource(resourcevo.DeleteResourceReqVo{
			Input: bo.DeleteResourceBo{
				ResourceIds: relationIds,
				UserId: operatorId,
				OrgId: orgId,
				ProjectId: issueBo.ProjectId,
			},
		})
		if resp.Failure(){
			log.Error(resp.Message)
			return resp.Error()
		}
		return nil
	})

	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

func UpdateIssueRelationSingle(operatorId int64, issueBo bo.IssueBo, relationType int, newUserId int64) (*bo.IssueRelationBo, errs.SystemErrorInfo) {
	bos, err := UpdateIssueRelation(operatorId, issueBo, relationType, []int64{newUserId})
	if err != nil {
		return nil, err
	}
	if len(bos) == 0 {
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError)
	}
	return &bos[0], nil
}

func UpdateIssueRelation(operatorId int64, issueBo bo.IssueBo, relationType int, newUserIds []int64) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	orgId := issueBo.OrgId
	issueId := issueBo.Id

	//防止项目成员重复插入
	uid := uuid.NewUuid()
	issueIdStr := strconv.FormatInt(issueId, 10)
	relationTypeStr := strconv.Itoa(relationType)
	lockKey := consts.AddIssueRelationLock + issueIdStr + "#" + relationTypeStr
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
	issueRelations := &[]po.PpmPriIssueRelation{}
	err5 := mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
		consts.TcIssueId:      issueBo.Id,
		consts.TcRelationId:   db.In(newUserIds),
		consts.TcRelationType: relationType,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, issueRelations)
	if err5 != nil {
		log.Error(err5)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	//check，去掉已有的关联
	if len(*issueRelations) > 0 {
		alreadyExistRelationIdMap := map[int64]bool{}
		for _, issueRelation := range *issueRelations {
			alreadyExistRelationIdMap[issueRelation.RelationId] = true
		}
		notRelationUserIds := make([]int64, 0)
		for _, newUserId := range newUserIds {
			if _, ok := alreadyExistRelationIdMap[newUserId]; !ok {
				notRelationUserIds = append(notRelationUserIds, newUserId)
			}
		}
		newUserIds = notRelationUserIds
	}
	newUserIds = slice.SliceUniqueInt64(newUserIds)

	issueRelationBos := make([]bo.IssueRelationBo, len(newUserIds))
	if len(newUserIds) == 0 {
		return issueRelationBos, nil
	}

	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueRelation, len(newUserIds))
	if err != nil {
		log.Errorf("id generate: %q\n", err)
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}

	issueRelationPos := make([]po.PpmPriIssueRelation, len(newUserIds))
	for i, newUserId := range newUserIds {
		id := ids.Ids[i].Id
		issueRelation := &po.PpmPriIssueRelation{}
		issueRelation.Id = id
		issueRelation.OrgId = orgId
		issueRelation.ProjectId = issueBo.ProjectId
		issueRelation.IssueId = issueBo.Id
		issueRelation.RelationId = newUserId
		issueRelation.RelationType = relationType
		issueRelation.Creator = operatorId
		issueRelation.Updator = operatorId
		issueRelation.IsDelete = consts.AppIsNoDelete
		issueRelationPos[i] = *issueRelation

		issueRelationBos[i] = bo.IssueRelationBo{
			Id:           id,
			OrgId:        issueBo.OrgId,
			IssueId:      issueBo.Id,
			RelationId:   newUserId,
			RelationType: consts.IssueRelationTypeOwner,
			Creator:      operatorId,
			CreateTime:   types.NowTime(),
			Updator:      operatorId,
			UpdateTime:   types.NowTime(),
			Version:      1,
		}
	}

	err2 := mysql.BatchInsert(&po.PpmPriIssueRelation{}, slice.ToSlice(issueRelationPos))
	if err2 != nil {
		log.Errorf("mysql.BatchInsert(): %q\n", err2)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}
	return issueRelationBos, nil
}

func GetIssueRelationIdsByRelateType(orgId int64, issueId int64, relationType int) (*[]int64, errs.SystemErrorInfo) {
	issueParticipantRelations, _, err := dao.SelectIssueRelationByPage(db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      issueId,
		consts.TcRelationType: relationType,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, bo.PageBo{
		Order: consts.TcCreateTime + " desc",
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	relationIds := make([]int64, len(*issueParticipantRelations))
	for i, participantRelation := range *issueParticipantRelations {
		relationIds[i] = participantRelation.RelationId
	}
	return &relationIds, nil
}

func GetIssueRelationByRelateTypeList(orgId int64, issueId int64, relationTypes []int) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	issueParticipantRelations, _, err := dao.SelectIssueRelationByPage(db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      issueId,
		consts.TcRelationType: db.In(relationTypes),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, bo.PageBo{
		Order: consts.TcCreateTime + " desc",
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	relationBos := &[]bo.IssueRelationBo{}
	_ = copyer.Copy(issueParticipantRelations, relationBos)
	return *relationBos, nil
}

//创建者也可以删除资源
func GetIssueResourceIdsByCreator(orgId int64, issueId int64, ids []int64, creatorId int64) (*[]int64, errs.SystemErrorInfo) {
	issueRelations, _, err := dao.SelectIssueRelationByPage(db.Cond{
		consts.TcId:           db.In(ids),
		consts.TcOrgId:        orgId,
		consts.TcIssueId:      issueId,
		consts.TcRelationType: consts.IssueRelationTypeResource,
		consts.TcCreator:      creatorId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, bo.PageBo{
		Order: consts.TcCreateTime + " desc",
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	relationIds := make([]int64, len(*issueRelations))
	for i, participantRelation := range *issueRelations {
		relationIds[i] = participantRelation.Id
	}
	return &relationIds, nil
}

func VerifyRelationIssue(issueIds []int64, projectObjectTypeId int64, orgId int64) errs.SystemErrorInfo {
	issueList := &[]po.PpmPriIssue{}
	err := mysql.SelectAllByCond(consts.TableIssue, db.Cond{
		consts.TcOrgId:               orgId,
		consts.TcId:                  db.In(issueIds),
		consts.TcProjectObjectTypeId: projectObjectTypeId,
	}, issueList)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	issueIds = slice.SliceUniqueInt64(issueIds)
	trueIssueIds := []int64{}
	for _, v := range *issueList {
		trueIssueIds = append(trueIssueIds, v.Id)
	}
	//去重后的数组长度相同则表示传递的任务id都有效
	if len(issueIds) == len(trueIssueIds) {
		return nil
	}

	return errs.BuildSystemErrorInfo(errs.RelationIssueError)
}

func RelationIssueList(orgId, issueId int64) ([]po.PpmPriIssue, errs.SystemErrorInfo) {
	issueRelationList := &[]po.PpmPriIssueRelation{}
	issueList := &[]po.PpmPriIssue{}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcIssueId:  issueId,
	}, issueRelationList)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	if len(*issueRelationList) == 0 {
		return *issueList, nil
	}

	issueIds := []int64{}
	for _, v := range *issueRelationList {
		issueIds = append(issueIds, v.RelationId)
	}

	err = mysql.SelectAllByCond(consts.TableIssue, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(issueIds),
	}, issueList)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	return *issueList, nil
}

func GetIssueRelationByResource(orgId int64, projectId int64, resourceIds []int64) (*[]po.PpmPriIssueRelation, errs.SystemErrorInfo) {
	issueParticipantRelations, _, err := dao.SelectIssueRelationByPage(db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: consts.IssueRelationTypeResource,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationId:   db.In(resourceIds),
	}, bo.PageBo{
		Order: consts.TcId + " desc",
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return issueParticipantRelations, nil
}

func GetTotalResourceByRelationCond(cond db.Cond) (*[]po.PpmPriIssueRelation, errs.SystemErrorInfo) {
	pos := &[]po.PpmPriIssueRelation{}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, cond, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	resourceIds := []int64{}
	for _, value := range *pos {
		if isContain, _ := slice.Contain(resourceIds, value.RelationId); !isContain {
			resourceIds = append(resourceIds, value.RelationId)
		}
	}
	return pos, nil
}

func DeleteProjectAttachment(orgId, operatorId, projectId int64, resourceIds []int64) errs.SystemErrorInfo {
	issueRelationPos, err := GetIssueRelationByResource(orgId, projectId, resourceIds)
	if err != nil {
		log.Error(err)
		return err
	}
	realResourceIds := make([]int64, 0)
	relationIds := make([]int64, len(*issueRelationPos))
	realResourceMap := make(map[int64]bool)
	for index, value := range *issueRelationPos {
		realResourceMap[value.RelationId] = true
		relationIds[index] = value.Id
	}
	for key, _ := range realResourceMap {
		realResourceIds = append(realResourceIds, key)
	}
	if len(realResourceIds) != len(resourceIds) {
		return errs.InvalidResourceIdsError
	}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := DeleteAllIssueRelationByIds(orgId, operatorId, relationIds, tx)
		if err != nil {
			log.Error(err)
			return nil
		}
		deleteInput := bo.DeleteResourceBo{
			ResourceIds: resourceIds,
			UserId:      operatorId,
			OrgId:       orgId,
			ProjectId:   projectId,
		}
		resp := resourcefacade.DeleteResource(resourcevo.DeleteResourceReqVo{Input: deleteInput})
		if resp.Failure() {
			log.Error(resp.Error())
			return resp.Error()
		}
		return nil
	})

	asyn.Execute(func() {
		reqVo := resourcevo.GetResourceByIdReqBody{
			ResourceIds: resourceIds,
		}
		req := resourcevo.GetResourceByIdReqVo{GetResourceByIdReqBody: reqVo}
		resp := resourcefacade.GetResourceById(req)
		resourceBos := resp.ResourceBos
		resourceNames := make([]string, len(resourceBos))
		for index, value := range resourceBos {
			resourceNames[index] = value.Name
		}

		trendBo := bo.ProjectTrendsBo{
			PushType:   consts.PushTypeDeleteResource,
			OrgId:      orgId,
			ProjectId:  projectId,
			OperatorId: operatorId,
			NewValue:   json.ToJsonIgnoreError(resourceNames),
		}

		asyn.Execute(func() {
			PushProjectTrends(trendBo)
		})
		asyn.Execute(func() {
			PushProjectThirdPlatformNotice(trendBo)
		})
	})

	return nil
}

func GetRelationInfoByIssueIds(issueIds []int64, relationTypes []int) ([]bo.IssueRelationBo, errs.SystemErrorInfo) {
	relationInfos := &[]po.PpmPriIssueRelation{}
	cond := db.Cond{
		consts.TcIsDelete:consts.AppIsNoDelete,
		consts.TcIssueId:db.In(issueIds),
	}
	if len(relationTypes) != 0 {
		cond[consts.TcRelationType] = db.In(relationTypes)
	}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, cond, relationInfos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.IssueRelationBo{}
	copyErr := copyer.Copy(relationInfos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.ObjectCopyError
	}

	return *bos, nil
}
func GetIssueMembers(orgId int64, issueId int64) (*bo.IssueMembersBo, errs.SystemErrorInfo){
	relationBos, err := GetIssueRelationByRelateTypeList(orgId, issueId, consts.MemberRelationTypeList)
	if err != nil{
		log.Error(err)
		return nil, err
	}

	var owner int64
	participantMap := map[int64]bool{}
	followerMap := map[int64]bool{}
	memberMap := map[int64]bool{}
	for _, v := range relationBos {
		if v.RelationType == consts.IssueRelationTypeFollower {
			followerMap[v.RelationId] = true
		} else if v.RelationType == consts.IssueRelationTypeParticipant {
			participantMap[v.RelationId] = true
		} else if v.RelationType == consts.IssueRelationTypeOwner {
			owner = v.RelationId
		}
		memberMap[v.RelationId] = true
	}

	var followerIds, participantIds, memberIds []int64
	for k, _ := range followerMap{
		followerIds = append(followerIds, k)
	}

	for k, _ := range participantMap{
		participantIds = append(participantIds, k)
	}

	for k, _ := range memberMap{
		memberIds = append(memberIds, k)
	}


	return &bo.IssueMembersBo{
		MemberIds: memberIds,
		OwnerId: owner,
		ParticipantIds: participantIds,
		FollowerIds: followerIds,
	}, nil
}
