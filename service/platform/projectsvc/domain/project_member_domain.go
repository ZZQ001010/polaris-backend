package domain

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"gopkg.in/fatih/set.v0"
	"strconv"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func HandleProjectMember(orgId int64, currentUserId int64, owner int64, projectId int64, memberIds []int64, followerIds []int64) ([]interface{}, []int64, errs.SystemErrorInfo) {
	//插入项目成员
	memberEntities := []interface{}{}
	addedMemberIds := []int64{}
	memberEntity := po.PpmProProjectRelation{}

	//1.负责人
	memberId, err := idfacade.ApplyPrimaryIdRelaxed(memberEntity.TableName())
	if err != nil {
		return memberEntities, nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	pass := orgfacade.VerifyOrgRelaxed(orgId, owner)

	if !pass {
		log.Info("owner " + strconv.FormatInt(owner, 10))
		return memberEntities, nil, errs.BuildSystemErrorInfo(errs.VerifyOrgError)
	}

	memberEntities = append(memberEntities, po.PpmProProjectRelation{
		Id:           memberId,
		OrgId:        orgId,
		ProjectId:    projectId,
		RelationId:   owner,
		RelationType: consts.IssueRelationTypeOwner,
		Creator:      currentUserId,
		CreateTime:   time.Now(),
		IsDelete:     consts.AppIsNoDelete,
		Status:       consts.ProjectMemberEffective,
		Updator:      currentUserId,
		UpdateTime:   time.Now(),
		Version:      1,
	})
	addedMemberIds = append(addedMemberIds, owner)

	//2.项目成员
	//默认创建者也是项目成员
	if owner != currentUserId {
		if bool, _ := slice.Contain(memberIds, currentUserId); !bool {
			memberIds = append(memberIds, currentUserId)
		}
	}

	memberIds = slice.SliceUniqueInt64(memberIds)
	if len(memberIds) != 0 {
		for _, v := range memberIds {
			pass := orgfacade.VerifyOrgRelaxed(orgId, v)
			if !pass {
				return memberEntities, nil, errs.BuildSystemErrorInfo(errs.VerifyOrgError)
			}
			memberId, err := idfacade.ApplyPrimaryIdRelaxed(memberEntity.TableName())
			if err != nil {
				return memberEntities, nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
			}
			memberEntities = append(memberEntities, po.PpmProProjectRelation{
				Id:           memberId,
				OrgId:        orgId,
				ProjectId:    projectId,
				RelationId:   v,
				RelationType: consts.IssueRelationTypeParticipant,
				Creator:      currentUserId,
				CreateTime:   time.Now(),
				IsDelete:     consts.AppIsNoDelete,
				Status:       consts.ProjectMemberEffective,
				Updator:      currentUserId,
				UpdateTime:   time.Now(),
				Version:      1,
			})
			addedMemberIds = append(addedMemberIds, v)
		}
	}

	//3.项目关注人
	followerIds = slice.SliceUniqueInt64(followerIds)
	if len(followerIds) != 0 {
		for _, v := range followerIds {
			pass := orgfacade.VerifyOrgRelaxed(orgId, v)
			if !pass {
				return memberEntities, nil, errs.BuildSystemErrorInfo(errs.VerifyOrgError)
			}
			followerId, err := idfacade.ApplyPrimaryIdRelaxed(memberEntity.TableName())
			if err != nil {
				return memberEntities, nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
			}
			memberEntities = append(memberEntities, po.PpmProProjectRelation{
				Id:           followerId,
				OrgId:        orgId,
				ProjectId:    projectId,
				RelationId:   v,
				RelationType: consts.IssueRelationTypeFollower,
				Creator:      currentUserId,
				CreateTime:   time.Now(),
				IsDelete:     consts.AppIsNoDelete,
				Status:       consts.ProjectMemberEffective,
				Updator:      currentUserId,
				UpdateTime:   time.Now(),
				Version:      1,
			})
			addedMemberIds = append(addedMemberIds, v)
		}
	}

	return memberEntities, addedMemberIds, nil
}

//我参与的
func GetParticipantMembers(orgId, currentUserId int64) ([]int64, errs.SystemErrorInfo) {
	projectIdsNeed := []int64{}
	memberEntities := &[]*po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond((&po.PpmProProjectRelation{}).TableName(), db.Cond{
		consts.TcIsDelete:     db.Eq(consts.AppIsNoDelete),
		consts.TcRelationType: db.Eq(consts.IssueRelationTypeParticipant),
		consts.TcRelationId:   db.Eq(currentUserId),
		consts.TcOrgId:        orgId,
	}, memberEntities)
	if err != nil {
		return projectIdsNeed, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	for _, v := range *memberEntities {
		projectIdsNeed = append(projectIdsNeed, v.ProjectId)
	}

	return projectIdsNeed, nil
}

//我参与的和我负责的
func GetParticipantMembersAndOwner(orgId, currentUserId int64) ([]int64, errs.SystemErrorInfo) {
	projectIdsNeed := []int64{}
	memberEntities := &[]*po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond((&po.PpmProProjectRelation{}).TableName(), db.Cond{
		consts.TcIsDelete:     db.Eq(consts.AppIsNoDelete),
		consts.TcRelationType: db.In([]int{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeOwner}),
		consts.TcRelationId:   db.Eq(currentUserId),
		consts.TcOrgId:        orgId,
	}, memberEntities)
	if err != nil {
		return projectIdsNeed, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	for _, v := range *memberEntities {
		projectIdsNeed = append(projectIdsNeed, v.ProjectId)
	}

	return projectIdsNeed, nil
}

func GetProjectMemberInfo(projectIds []int64, orgId int64, creatorIds []int64, sourceChannel string) (map[int64]bo.UserIDInfoBo, map[int64][]bo.UserIDInfoBo, map[int64][]bo.UserIDInfoBo, map[int64]bo.UserIDInfoBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	relatedInfo := &[]bo.RelationInfoTypeBo{}
	err1 := conn.Select("relation_id", "relation_type", "project_id").From("ppm_pro_project_relation").
		Where(db.Cond{
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcProjectId:    db.In(projectIds),
			consts.TcStatus:       1,
			consts.TcOrgId:        orgId,
			consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeOwner, consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower}),
		}).All(relatedInfo)
	if err1 != nil {
		log.Error(err1)
		return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	creatorInfo := map[int64]bo.UserIDInfoBo{}

	ownerInfo := map[int64]bo.UserIDInfoBo{}
	participantInfo := map[int64][]bo.UserIDInfoBo{}
	followerInfo := map[int64][]bo.UserIDInfoBo{}
	relatedIds := map[int64][]int64{}
	allRelationIds := []int64{}
	for _, v := range *relatedInfo {
		allRelationIds = append(allRelationIds, v.RelationId)
	}
	allRelationIds = append(allRelationIds, creatorIds...)
	allUserInfo, err := orgfacade.GetBaseUserInfoBatchRelaxed(sourceChannel, orgId, allRelationIds)
	if err != nil {
		return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	userInfoById := map[int64]bo.BaseUserInfoBo{}
	for _, v := range allUserInfo {
		userInfoById[v.UserId] = v
	}
	for _, v := range creatorIds {
		userInfo, ok := userInfoById[v]
		if !ok {
			continue
		}
		temp := bo.UserIDInfoBo{}

		temp.Name = userInfo.Name
		temp.Avatar = userInfo.Avatar
		temp.UserID = userInfo.UserId
		temp.EmplID = userInfo.OutUserId
		temp.IsDeleted = userInfo.OrgUserIsDelete == consts.AppIsDeleted
		temp.IsDisabled = userInfo.OrgUserStatus == consts.AppStatusDisabled
		creatorInfo[v] = temp
	}
	for _, v := range *relatedInfo {
		//userInfo, err := orgfacade.GetUserInfoRelaxed(orgId, v.RelationId)
		//if err != nil {
		//	return nil, nil, nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
		//}
		userInfo, ok := userInfoById[v.RelationId]
		if !ok {
			continue
		}
		temp := bo.UserIDInfoBo{}

		temp.Name = userInfo.Name
		temp.Avatar = userInfo.Avatar
		temp.UserID = userInfo.UserId
		temp.EmplID = userInfo.OutUserId
		temp.IsDeleted = userInfo.OrgUserIsDelete == consts.AppIsDeleted
		temp.IsDisabled = userInfo.OrgUserStatus == consts.AppStatusDisabled
		if v.RelationType == consts.IssueRelationTypeOwner {
			ownerInfo[v.ProjectId] = temp
		} else if v.RelationType == consts.IssueRelationTypeParticipant {
			participantInfo[v.ProjectId] = append(participantInfo[v.ProjectId], temp)
		} else if v.RelationType == consts.IssueRelationTypeFollower {
			followerInfo[v.ProjectId] = append(followerInfo[v.ProjectId], temp)
		}
		relatedIds[v.ProjectId] = append(relatedIds[v.ProjectId], v.UserId)
	}

	return ownerInfo, participantInfo, followerInfo, creatorInfo, nil
}

func JudgeIsProjectMember(currentUserId, orgId, projectId int64) (po.PpmProProjectRelation, errs.SystemErrorInfo) {
	member := &po.PpmProProjectRelation{}
	err := mysql.SelectOneByCond(member.TableName(), db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationId:   currentUserId,
		consts.TcProjectId:    projectId,
		consts.TcOrgId:        orgId,
		consts.TcRelationType: consts.IssueRelationTypeParticipant,
	}, member)
	if err != nil {
		return *member, errs.BuildSystemErrorInfo(errs.NotProjectParticipant)
	}

	return *member, nil
}

func GetChangeMembersAndDeleteOld(tx sqlbuilder.Tx, input bo.UpdateProjectBo, orgId int64, oldOwner int64, updPoint *mysql.Upd) (set.Interface, set.Interface, errs.SystemErrorInfo) {
	upd := *updPoint
	oldMembers := set.New(set.ThreadSafe)
	thisMembers := set.New(set.ThreadSafe)

	//成员更新
	if err := assemblyMembers(input, orgId, &oldMembers, &thisMembers, tx); err != nil {
		return nil, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	//负责人更新
	if util.FieldInUpdate(input.UpdateFields, "owner") {
		if input.Owner != nil && *input.Owner != oldOwner {
			upd[consts.TcOwner] = *input.Owner
			//oldMembers.Add(oldOwner)
			//删除旧有负责人
			err := tx.Collection((&po.PpmProProjectRelation{}).TableName()).Find(db.Cond{
				consts.TcRelationId:   oldOwner,
				consts.TcOrgId:        db.Eq(orgId),
				consts.TcProjectId:    db.Eq(input.ID),
				consts.TcRelationType: consts.IssueRelationTypeOwner,
			}).Update(map[string]int{
				consts.TcIsDelete: consts.AppIsDeleted,
			})
			if err != nil {
				return nil, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
			//(删除新负责人可能是旧参与人)
			err = tx.Collection((&po.PpmProProjectRelation{}).TableName()).Find(db.Cond{
				consts.TcRelationId:   input.Owner,
				consts.TcOrgId:        db.Eq(orgId),
				consts.TcProjectId:    db.Eq(input.ID),
				consts.TcRelationType: consts.IssueRelationTypeParticipant,
			}).Update(map[string]int{
				consts.TcIsDelete: consts.AppIsDeleted,
			})
			if err != nil {
				return nil, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
		}
	}

	return oldMembers, thisMembers, nil
}

func assemblyMembers(input bo.UpdateProjectBo, orgId int64, oldMembers, thisMembers *set.Interface, tx sqlbuilder.Tx) error {

	if util.FieldInUpdate(input.UpdateFields, "memberIds") {
		memberEntities := &[]po.PpmProProjectRelation{}
		err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcProjectId:    input.ID,
			consts.TcRelationType: consts.IssueRelationTypeParticipant,
			consts.TcIsDelete:     consts.AppIsNoDelete,
		}, memberEntities)
		if err != nil {
			return err
		}
		for _, v := range *memberEntities {
			(*oldMembers).Add(v.RelationId)
		}

		input.MemberIds = slice.SliceUniqueInt64(input.MemberIds)
		for _, v := range input.MemberIds {
			(*thisMembers).Add(v)
		}

		delMembers := set.Difference(*oldMembers, *thisMembers)
		if len(delMembers.List()) != 0 {
			err := tx.Collection((&po.PpmProProjectRelation{}).TableName()).Find(db.Cond{
				consts.TcRelationId:   db.In(delMembers.List()),
				consts.TcOrgId:        db.Eq(orgId),
				consts.TcProjectId:    db.Eq(input.ID),
				consts.TcRelationType: consts.IssueRelationTypeParticipant,
			}).Update(map[string]int{
				consts.TcIsDelete: consts.AppIsDeleted,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ChangeMembers(tx sqlbuilder.Tx, input bo.UpdateProjectBo, orgId int64, currentUserId int64, oldOwner int64, oldMembers, thisMembers set.Interface) errs.SystemErrorInfo {
	//需要删除的
	delMembers := set.Difference(oldMembers, thisMembers)
	memberAdd := []interface{}{}

	projectOwner := oldOwner

	ownerError := updateOwn(input, oldOwner, orgId, currentUserId, &projectOwner, &memberAdd)
	if ownerError != nil {
		return ownerError
	}

	if len(delMembers.List()) != 0 {
		err := DeleteRelationByDeleteMember(tx, delMembers.List(), projectOwner, input.ID, orgId, currentUserId)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
		}
	}

	addMembers := set.Difference(thisMembers, oldMembers)

	for _, v := range addMembers.List() {
		val, ok := v.(int64)
		if !ok {
			continue
		}
		pass := orgfacade.VerifyOrgRelaxed(orgId, val)
		if !pass {
			return errs.BuildSystemErrorInfo(errs.VerifyOrgError)
		}
		memberId, err := idfacade.ApplyPrimaryIdRelaxed((&po.PpmProProjectRelation{}).TableName())
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
		}
		memberAdd = append(memberAdd, po.PpmProProjectRelation{
			Id:           memberId,
			OrgId:        orgId,
			ProjectId:    input.ID,
			RelationId:   val,
			RelationType: consts.IssueRelationTypeParticipant,
			Creator:      currentUserId,
			CreateTime:   time.Now(),
			IsDelete:     consts.AppIsNoDelete,
			Status:       consts.ProjectMemberEffective,
			Updator:      currentUserId,
			UpdateTime:   time.Now(),
			Version:      1,
		})
	}

	err := mysql.TransBatchInsert(tx, &po.PpmProProjectRelation{}, memberAdd)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

func updateOwn(input bo.UpdateProjectBo, oldOwner int64, orgId int64, currentUserId int64, projectOwner *int64, memberAdd *[]interface{}) errs.SystemErrorInfo {

	if util.FieldInUpdate(input.UpdateFields, "owner") {
		if input.Owner != nil && *input.Owner != oldOwner {
			pass := orgfacade.VerifyOrgRelaxed(orgId, *input.Owner)
			if !pass {
				return errs.BuildSystemErrorInfo(errs.VerifyOrgError)
			}
			memberId, err := idfacade.ApplyPrimaryIdRelaxed((&po.PpmProProjectRelation{}).TableName())
			if err != nil {
				return errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
			}
			*memberAdd = append(*memberAdd, po.PpmProProjectRelation{
				Id:           memberId,
				OrgId:        orgId,
				ProjectId:    input.ID,
				RelationId:   *input.Owner,
				RelationType: consts.IssueRelationTypeOwner,
				Creator:      currentUserId,
				CreateTime:   time.Now(),
				IsDelete:     consts.AppIsNoDelete,
				Status:       consts.ProjectMemberEffective,
				Updator:      currentUserId,
				UpdateTime:   time.Now(),
				Version:      1,
			})
			*projectOwner = *input.Owner
		}
	}
	return nil
}

func JudgeIsFollower(projectId, currentUserId, orgId int64) (bool, errs.SystemErrorInfo) {
	isExist, err := mysql.IsExistByCond(consts.TableProjectRelation, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcRelationId:   currentUserId,
		consts.TcStatus:       1,
		consts.TcRelationType: consts.IssueRelationTypeFollower,
		consts.TcProjectId:    projectId,
	})
	if err != nil {
		return isExist, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return isExist, nil
}

func JudgeIsMember(projectId, currentUserId, orgId int64) (bool, errs.SystemErrorInfo) {
	isExist, err := mysql.IsExistByCond(consts.TableProjectRelation, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcRelationId:   currentUserId,
		consts.TcStatus:       1,
		consts.TcRelationType: db.In(consts.MemberRelationTypeList),
		consts.TcProjectId:    projectId,
	})
	if err != nil {
		return isExist, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return isExist, nil
}

func AddMember(projectId, orgId, userId, currentUserId int64, relationType int) (bool, errs.SystemErrorInfo) {
	memberId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectRelation)
	if err != nil {
		return false, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	err1 := mysql.Insert(&po.PpmProProjectRelation{
		Id:           memberId,
		IsDelete:     consts.AppIsNoDelete,
		OrgId:        orgId,
		ProjectId:    projectId,
		Status:       1,
		RelationType: relationType,
		RelationId:   userId,
		Creator:      currentUserId,
		CreateTime:   time.Now(),
		Updator:      currentUserId,
		UpdateTime:   time.Now(),
	})
	if err1 != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return true, nil
}

func DeleteMember(projectId, orgId, userId, currentUserId int64, relationType int) (bool, errs.SystemErrorInfo) {
	_, err := mysql.UpdateSmartWithCond(consts.TableProjectRelation, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcRelationId:   userId,
		consts.TcStatus:       1,
		consts.TcRelationType: relationType,
		consts.TcProjectId:    projectId,
	}, mysql.Upd{
		consts.TcIsDelete:   consts.AppIsDeleted,
		consts.TcUpdator:    currentUserId,
		consts.TcUpdateTime: time.Now(),
	})
	if err != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return true, nil
}

func RemoveProjectMember(orgId, userId int64, input vo.RemoveProjectMemberReq) errs.SystemErrorInfo {
	if len(input.MemberIds) == 0 {
		return errs.UpdateMemberIdsIsEmptyError
	}

	projectId := input.ProjectID

	projectInfo, infoErr := GetProjectInfo(projectId, orgId)
	if infoErr != nil {
		log.Error(infoErr)
		return infoErr
	}
	//负责人不能被移除
	if ok, _ := slice.Contain(input.MemberIds, projectInfo.Owner); ok {
		return errs.CannotRemoveProjectOwner
	}
	//删除项目关联
	_, updateErr := mysql.UpdateSmartWithCond(consts.TableProjectRelation, db.Cond{
		consts.TcProjectId:    projectId,
		consts.TcRelationId:   db.In(input.MemberIds),
		consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower}),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcUpdator:  userId,
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	if updateErr != nil {
		log.Error(updateErr)
		return errs.MysqlOperateError
	}

	//清掉用户的角色
	roleErr := rolefacade.RemoveRoleUserRelation(rolevo.RemoveRoleUserRelationReqVo{
		OrgId:      orgId,
		UserIds:    input.MemberIds,
		OperatorId: userId,
	})
	if roleErr.Failure() {
		log.Error(roleErr.Message)
	}

	//最后将用户信息缓存清掉
	clearErr := rolefacade.ClearUserRoleList(rolevo.ClearUserRoleReqVo{
		ProjectId: projectId,
		OrgId:     orgId,
		UserIds:   input.MemberIds,
	})
	if clearErr.Failure() {
		log.Error(clearErr.Error())
		return clearErr.Error()
	}

	refreshProjectAuthErr := RefreshProjectAuthBo(orgId, projectId)
	if refreshProjectAuthErr != nil {
		log.Error(refreshProjectAuthErr)
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = projectInfo.Name
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:      consts.PushTypeUpdateProjectMembers,
			OrgId:         orgId,
			ProjectId:     projectId,
			OperatorId:    userId,
			BeforeChangeMembers: input.MemberIds,
			AfterChangeMembers:  []int64{},
			Ext:           ext,
		})
	})
	return nil
}

func AddProjectMember(orgId, userId int64, input vo.RemoveProjectMemberReq) errs.SystemErrorInfo {
	if len(input.MemberIds) == 0 {
		return errs.UpdateMemberIdsIsEmptyError
	}

	projectId := input.ProjectID
	//判断项目是否存在
	projectInfo, infoErr := GetProjectInfo(projectId, orgId)
	if infoErr != nil {
		log.Error(infoErr)
		return infoErr
	}

	verifyOrgUserFlag := orgfacade.VerifyOrgUsersRelaxed(orgId, input.MemberIds)
	if ! verifyOrgUserFlag{
		log.Error("存在用户组织校验失败")
		return errs.VerifyOrgError
	}
	addIds, updateProjectRelationErr := UpdateProjectRelationWithRelationTypes(userId, orgId, projectId, consts.MemberRelationTypeList, consts.IssueRelationTypeParticipant, input.MemberIds)
	if updateProjectRelationErr != nil {
		log.Error(updateProjectRelationErr)
		return updateProjectRelationErr
	}

	refreshProjectAuthErr := RefreshProjectAuthBo(orgId, projectId)
	if refreshProjectAuthErr != nil {
		log.Error(refreshProjectAuthErr)
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = projectInfo.Name
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:      consts.PushTypeUpdateProjectMembers,
			OrgId:         orgId,
			ProjectId:     projectId,
			OperatorId:    userId,
			BeforeChangeMembers: []int64{},
			AfterChangeMembers:  addIds,
			Ext:           ext,
		})
	})
	return nil
}

func GetProjectAllMember(orgId, projectId int64, page, size int) (int64, []bo.ProjectRelationBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	defer func() {
		if err := conn.Close(); err != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	}()
	if err != nil {
		log.Error(err)
		return 0, nil, errs.MysqlOperateError
	}

	cond := db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeOwner, consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower}),
	}
	countPo := &po.PpmProProjectRelation{}
	countErr := conn.Select(db.Raw("count(distinct relation_id) as id")).From(consts.TableProjectRelation).Where(cond).One(countPo)
	if countErr != nil {
		log.Error(countErr)
		return 0, nil, errs.MysqlOperateError
	}
	count := countPo.Id
	pos := &[]po.PpmProProjectRelation{}
	//获取所有成员（最小的relation_id代表最高的用户角色（负责人），最早的创建时间表示加入时间，随机挑选一名操作人）
	mid := conn.Select(db.Raw("relation_id, min(relation_type) as relation_type, min(create_time) as create_time, min(creator) as creator")).
		From(consts.TableProjectRelation).
		Where(cond).GroupBy(consts.TcRelationId)
	if page > 0 && size > 0 {
		mid.Offset((page - 1) * size).Limit(size)
	}
	selectErr := mid.All(pos)
	if selectErr != nil {
		log.Error(selectErr)
		return 0, nil, errs.MysqlOperateError
	}

	bos := &[]bo.ProjectRelationBo{}
	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}

	return count, *bos, nil
}
