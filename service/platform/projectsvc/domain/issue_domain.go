package domain

import (
	"fmt"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/core/util/times"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/core/util/str"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"strconv"
	"sync"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func CreateIssue(issueBo *bo.IssueBo, sourceChannel string) errs.SystemErrorInfo {
	newUUID := uuid.NewUuid()
	lockKey := fmt.Sprintf("%s%d", consts.IssueRelateOperationLock, issueBo.Id)
	if issueBo.ParentId != 0 {
		lockKey = fmt.Sprintf("%s%d", consts.IssueRelateOperationLock, issueBo.ParentId)
	}
	suc, lockErr := cache.TryGetDistributedLock(lockKey, newUUID)
	if lockErr != nil{
		log.Error(lockErr)
		return errs.TryDistributedLockError
	}
	if suc{
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, newUUID); err != nil{
				log.Error(err)
			}
		}()
	}else{
		//未获取到锁，直接响应错误信息
		return errs.CreateIssueFail
	}
	issuePo, issueRelationPos, issueDetailPo, err1 := initIssueInfo(issueBo)

	if err1 != nil {
		return err1
	}

	//先插入任务关联
	err4 := dao.InsertIssueRelationBatch(*issueRelationPos)
	if err4 != nil {
		log.Error(err4)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
	}

	//任务和任务明细需要事务
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		issuePo.Sort = issuePo.Id

		err4 = mysql.TransInsert(tx, issuePo)
		if err4 != nil {
			log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
		}
		err4 = mysql.TransInsert(tx, issueDetailPo)
		if err4 != nil {
			log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
		}
		return nil
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	insertProjectMembersInputBo := &bo.InsertProjectMembersInputBo{
		OrgId:      issueBo.OrgId,
		ProjectId:  issueBo.ProjectId,
		OwnerInfo:  issueBo.OwnerInfo,
		OperatorId: issueBo.Creator,
	}
	if issueBo.ParticipantInfos != nil {
		insertProjectMembersInputBo.ParticipantInfos = *issueBo.ParticipantInfos
	}
	if issueBo.FollowerInfos != nil {
		insertProjectMembersInputBo.FollowerInfos = *issueBo.FollowerInfos
	}

	//再插入项目成员关联，任务关联成功项目关联失败影响不大，不需要事务
	err4 = InsertProjectMembers(insertProjectMembersInputBo)
	if err4 != nil {
		log.Error(err4)
	}

	if issueBo.Tags != nil {
		//插入任务和标签关联影响也不大
		err4 = IssueRelateTags(issuePo.OrgId, issuePo.ProjectId, issuePo.Id, issuePo.Creator, issueBo.Tags, nil)
		if err4 != nil {
			log.Error(err4)
		}
	}

	asyn.Execute(func() {
		beforeChangeFollowers := []int64{}
		beforeChangeParticipants := []int64{}
		if issueBo.FollowerInfos != nil {
			for _, follower := range *issueBo.FollowerInfos {
				beforeChangeFollowers = append(beforeChangeFollowers, follower.RelationId)
			}
		}
		if issueBo.ParticipantInfos != nil {
			for _, participant := range *issueBo.ParticipantInfos {
				beforeChangeParticipants = append(beforeChangeParticipants, participant.RelationId)
			}
		}

		var planStartTime *types.Time = nil
		var planEndTime *types.Time = nil
		if issueBo.PlanStartTime.IsNotNull() {
			planStartTime = &issueBo.PlanStartTime
		}
		if issueBo.PlanEndTime.IsNotNull() {
			planEndTime = &issueBo.PlanEndTime
		}

		blankUserIds := []int64{}
		issueTrendsBo := bo.IssueTrendsBo{
			PushType:      consts.PushTypeCreateIssue,
			OrgId:         issueBo.OrgId,
			OperatorId:    issueBo.Creator,
			IssueId:       issuePo.Id,
			ParentIssueId: issuePo.ParentId,
			ProjectId:     issuePo.ProjectId,
			PriorityId:    issuePo.PriorityId,
			ParentId:      issuePo.ParentId,

			IssueTitle:               issuePo.Title,
			IssueStatusId:            issuePo.Status,
			BeforeOwner:              issuePo.Owner,
			AfterOwner:               issuePo.Owner,
			BeforeChangeFollowers:    beforeChangeFollowers,
			AfterChangeFollowers:     blankUserIds,
			BeforeChangeParticipants: beforeChangeParticipants,
			AfterChangeParticipants:  blankUserIds,

			IssuePlanStartTime: planStartTime,
			IssuePlanEndTime:   planEndTime,

			NewValue: json.ToJsonIgnoreError(issueBo),
		}

		asyn.Execute(func() {
			PushIssueTrends(issueTrendsBo)
		})
		asyn.Execute(func() {
			PushIssueThirdPlatformNotice(issueTrendsBo)
		})
		asyn.Execute(func() {
			if sourceChannel == consts.AppSourceChannelFeiShu {
				CreateCalendarEvent(issueBo, issueBo.Creator, beforeChangeParticipants)
			}
		})
	})
	return nil
}

func GetCalendarInfo(orgId, projectId int64) (bool, string, string) {
	projectCalendarInfo, err := GetProjectCalendarInfo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return false, "", ""
	}
	if projectCalendarInfo.IsSyncOutCalendar != consts.IsSyncOutCalendar || projectCalendarInfo.CalendarId == consts.BlankString {
		log.Error("无对应项目日历或未设置导出日历")
		return false, "", ""
	}
	orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(consts.AppSourceChannelFeiShu, orgId)
	if err != nil {
		log.Errorf("组织外部信息不存在 %v", err)
		return false, "", ""
	}

	return true, orgBaseInfo.OutOrgId, projectCalendarInfo.CalendarId
}

func initIssueInfo(issueBo *bo.IssueBo) (*po.PpmPriIssue, *[]po.PpmPriIssueRelation, *po.PpmPriIssueDetail, errs.SystemErrorInfo) {
	issuePo := &po.PpmPriIssue{}

	err1 := util.ConvertObject(&issueBo, &issuePo)
	if err1 != nil {
		return nil, nil, nil, err1
	}
	issuePo.IsDelete = consts.AppIsNoDelete

	issueRelationPos, err2 := buildIssueRelationBos(issueBo)
	if err2 != nil {
		return nil, nil, nil, err2
	}

	issueDetailPo := &po.PpmPriIssueDetail{}
	err3 := util.ConvertObject(&issueBo.IssueDetailBo, &issueDetailPo)
	if err3 != nil {
		return nil, nil, nil, err3
	}
	issueDetailPo.IsDelete = consts.AppIsNoDelete

	return issuePo, issueRelationPos, issueDetailPo, nil

}

func buildIssueRelationBos(issueBo *bo.IssueBo) (*[]po.PpmPriIssueRelation, errs.SystemErrorInfo) {
	issueRelationOwnerBo := issueBo.OwnerInfo.IssueRelationBo
	issueRelationFollowerBos := bo.BuildIssueRelationBosFromUserBos(issueBo.FollowerInfos)
	issueRelationParticipantBos := bo.BuildIssueRelationBosFromUserBos(issueBo.ParticipantInfos)

	ppo := &po.PpmPriIssueRelation{}
	err2 := util.ConvertObject(&issueRelationOwnerBo, &ppo)
	if err2 != nil {
		return nil, err2
	}

	fpos := &[]po.PpmPriIssueRelation{}
	err3 := util.ConvertObject(&issueRelationFollowerBos, &fpos)
	if err3 != nil {
		return nil, err3
	}
	ppos := &[]po.PpmPriIssueRelation{}
	err4 := util.ConvertObject(&issueRelationParticipantBos, &ppos)
	if err4 != nil {
		return nil, err4
	}

	relationCount := 1 + len(*fpos) + len(*ppos)

	issueRelationPos := make([]po.PpmPriIssueRelation, 0, relationCount)

	issueRelationPos = append(issueRelationPos, *ppo)
	issueRelationPos = append(issueRelationPos, *fpos...)
	issueRelationPos = append(issueRelationPos, *ppos...)

	return &issueRelationPos, nil
}

func InsertProjectMembers(insertProjectMembersInputBo *bo.InsertProjectMembersInputBo) errs.SystemErrorInfo {
	projectId := insertProjectMembersInputBo.ProjectId
	orgId := insertProjectMembersInputBo.OrgId

	//防止项目成员重复插入
	uid := uuid.NewUuid()
	projectIdStr := strconv.FormatInt(projectId, 10)
	lockKey := consts.AddIssueScheduleProjectMemberLock + projectIdStr
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("获取%s锁时异常 %v", lockKey, err)
		return errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	}
	projectRelations, err1 := GetBeInsertedProjectMembers(insertProjectMembersInputBo)
	if err1 != nil {
		log.Error(err1)
		return err1
	}

	//再插入项目成员关联，任务关联成功项目关联失败影响不大，不需要事务
	err4 := dao.InsertProjectRelationBatch(projectRelations)
	if err4 != nil {
		log.Error(err4)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
	}

	refreshErr := RefreshProjectAuthBo(orgId, projectId)
	if refreshErr != nil {
		log.Error(refreshErr)
	}
	return nil
}

func GetBeInsertedProjectMembers(insertProjectMembersInputBo *bo.InsertProjectMembersInputBo) ([]po.PpmProProjectRelation, errs.SystemErrorInfo) {
	projectId := insertProjectMembersInputBo.ProjectId
	orgId := insertProjectMembersInputBo.OrgId
	projectMembers := &[]po.PpmProProjectRelation{}

	strconv.FormatInt(projectId, 10)

	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: db.In([]int{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower}),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, projectMembers)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	relationUsers := &[]bo.IssueUserBo{}
	if insertProjectMembersInputBo.ParticipantInfos != nil {
		*relationUsers = append(*relationUsers, insertProjectMembersInputBo.ParticipantInfos...)
	}
	if insertProjectMembersInputBo.FollowerInfos != nil {
		*relationUsers = append(*relationUsers, insertProjectMembersInputBo.FollowerInfos...)
	}
	if insertProjectMembersInputBo.OwnerInfo != nil {
		*relationUsers = append(*relationUsers, *insertProjectMembersInputBo.OwnerInfo)
	}
	beInsertProjectParticipantMemberIds := &[]int64{}
	beInsertProjectFollowerMemberIds := &[]int64{}

	dealInsertMemberIds(relationUsers, projectMembers,
		beInsertProjectParticipantMemberIds, beInsertProjectFollowerMemberIds)

	beInsertProjectMembers := &[]po.PpmProProjectRelation{}

	projectMembersError := dealBeInserProjectMembers(beInsertProjectParticipantMemberIds, beInsertProjectFollowerMemberIds, beInsertProjectMembers,
		orgId, projectId, insertProjectMembersInputBo.OperatorId)

	if projectMembersError != nil {
		return nil, projectMembersError
	}

	return *beInsertProjectMembers, nil
}

func dealBeInserProjectMembers(beInsertProjectParticipantMemberIds, beInsertProjectFollowerMemberIds *[]int64,
	beInsertProjectMembers *[]po.PpmProProjectRelation, orgId, projectId,
	operatorId int64) errs.SystemErrorInfo {

	if len(*beInsertProjectParticipantMemberIds) > 0 {
		for _, memberId := range *beInsertProjectParticipantMemberIds {
			id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectRelation)
			if err != nil {
				return errs.BuildSystemErrorInfo(errs.ApplyIdError)
			}
			projectMember := po.PpmProProjectRelation{
				Id:           id,
				OrgId:        orgId,
				ProjectId:    projectId,
				TeamId:       0,
				RelationId:   memberId,
				RelationType: consts.IssueRelationTypeParticipant,
				Status:       consts.AppStatusEnable,
				Creator:      operatorId,
				IsDelete:     consts.AppIsNoDelete,
			}
			*beInsertProjectMembers = append(*beInsertProjectMembers, projectMember)
		}
	}

	if len(*beInsertProjectFollowerMemberIds) > 0 {
		for _, memberId := range *beInsertProjectFollowerMemberIds {
			id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectRelation)
			if err != nil {
				return errs.BuildSystemErrorInfo(errs.ApplyIdError)
			}
			projectMember := po.PpmProProjectRelation{
				Id:           id,
				OrgId:        orgId,
				ProjectId:    projectId,
				TeamId:       0,
				RelationId:   memberId,
				RelationType: consts.IssueRelationTypeFollower,
				Status:       consts.AppStatusEnable,
				Creator:      operatorId,
				IsDelete:     consts.AppIsNoDelete,
			}
			*beInsertProjectMembers = append(*beInsertProjectMembers, projectMember)
		}
	}
	return nil

}

//判断是否包含参与者和追随者
func judgecontainsParticipantAndFollower(projectMember po.PpmProProjectRelation, relationUser bo.IssueUserBo,
	containsParticipant *bool, containsFollower *bool) {

	if projectMember.RelationId == relationUser.RelationId {
		if projectMember.RelationType == consts.IssueRelationTypeParticipant {
			*containsParticipant = true
		} else {
			*containsFollower = true
		}
	}
}

//处理MermberIds
func dealInsertMemberIds(relationUsers *[]bo.IssueUserBo, projectMembers *[]po.PpmProProjectRelation,
	beInsertProjectParticipantMemberIds *[]int64, beInsertProjectFollowerMemberIds *[]int64) {

	for _, relationUser := range *relationUsers {
		containsParticipant := false
		containsFollower := false
		for _, projectMember := range *projectMembers {

			judgecontainsParticipantAndFollower(projectMember, relationUser, &containsParticipant, &containsFollower)
		}
		if !containsParticipant {
			if relationUser.RelationType == consts.IssueRelationTypeParticipant || relationUser.RelationType == consts.IssueRelationTypeOwner {
				*beInsertProjectParticipantMemberIds = append(*beInsertProjectParticipantMemberIds, relationUser.RelationId)
			}
		}
		if !containsFollower {
			if relationUser.RelationType == consts.IssueRelationTypeFollower {
				*beInsertProjectFollowerMemberIds = append(*beInsertProjectFollowerMemberIds, relationUser.RelationId)
			}
		}
	}
}

func DeleteRelationByDeleteMember(tx sqlbuilder.Tx, delMembers []interface{}, projectOwner int64, projectId int64, orgId int64, currentUserId int64) errs.SystemErrorInfo {
	issueRelation := &[]po.PpmPriIssueRelation{}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: consts.IssueRelationTypeOwner,
		consts.TcRelationId:   db.In(delMembers),
	}, issueRelation)

	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	//原有负责人、参与人直接从任务中移除
	_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeOwner}),
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationId:   db.In(delMembers),
	}, mysql.Upd{
		consts.TcIsDelete:   consts.AppIsDeleted,
		consts.TcUpdator:    currentUserId,
		consts.TcUpdateTime: time.Now(),
	})

	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	//任务负责人则移除并调整为项目负责人
	length := len(*issueRelation)
	if length > 0 {
		issueIds := make([]int64, length)
		issuePoInfo := make([]interface{}, length)
		for k, v := range *issueRelation {
			issueIds = append(issueIds, v.IssueId)
			id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueRelation)
			if err != nil {
				mysql.Rollback(tx)
				return errs.BuildSystemErrorInfo(errs.ApplyIdError)
			}
			issuePoInfo[k] = po.PpmPriIssueRelation{
				Id:           id,
				OrgId:        orgId,
				IssueId:      v.IssueId,
				RelationType: consts.IssueRelationTypeOwner,
				RelationId:   projectOwner,
				Creator:      currentUserId,
				CreateTime:   time.Now(),
				Updator:      currentUserId,
				UpdateTime:   time.Now(),
				IsDelete:     consts.AppIsNoDelete,
			}
		}

		err = mysql.TransBatchInsert(tx, &po.PpmPriIssueRelation{}, issuePoInfo)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
			consts.TcId: db.In(issueIds),
		}, mysql.Upd{
			consts.TcOwner: projectOwner,
		})
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}

	return nil
}

func GetIssueInfoByProject(projectIds []int64, orgId int64) (map[int64]bo.IssueStatistic, error) {
	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	issueStat := map[int64]bo.IssueStatistic{}
	//获取任务信息
	statistic := &[]*bo.IssueStatistic{}

	finishedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
	if err != nil || len(*finishedIds) == 0 {
		log.Errorf("proxies.GetProcessStatusId: %q\n", err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	var finishStr = consts.BlankString
	for _, v := range *finishedIds {
		finishStr += strconv.FormatInt(v, 10) + ","
	}
	finishStr = str.Substr(finishStr, 0, -1)

	_ = conn.Select(db.Raw("project_id,count(id) AS `all`,sum(CASE WHEN `status` not in (" + finishStr + ") and `plan_end_time` > '1970-01-01 00:00:00' and `plan_end_time` < now() THEN 1 ELSE 0 END) AS `overdue`,sum(CASE WHEN `status` in (" + finishStr + ") THEN 1 ELSE 0 END) AS `finish`")).From((&po.PpmPriIssue{}).TableName()).Where(db.Cond{
		consts.TcIsDelete:  db.Eq(consts.AppIsNoDelete),
		consts.TcProjectId: db.In(projectIds),
	}).GroupBy(consts.TcProjectId).All(statistic)
	for _, v := range *statistic {
		issueStat[v.ProjectId] = *v
	}

	return issueStat, nil
}

func UpdateIssue(issueUpdateBo bo.IssueUpdateBo, changeList []bo.TrendChangeListBo, sourceChannel string) errs.SystemErrorInfo {
	log.Info(strs.ObjectToString(issueUpdateBo))

	issueBo := issueUpdateBo.IssueBo

	newIssueBo := issueUpdateBo.NewIssueBo
	orgId := issueBo.OrgId
	operatorId := issueUpdateBo.OperatorId
	issueId := issueBo.Id

	beforeParticipantIds, err3 := GetIssueRelationIdsByRelateType(orgId, issueId, consts.IssueRelationTypeParticipant)
	if err3 != nil {
		log.Error(err3)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}
	afterParticipantIds := *beforeParticipantIds

	beforeFollowerIds, err3 := GetIssueRelationIdsByRelateType(orgId, issueId, consts.IssueRelationTypeFollower)
	if err3 != nil {
		log.Error(err3)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}
	afterFollowerIds := *beforeFollowerIds

	//更新负责人，不需要做diff
	err2 := assemblyUpdateOwnerInfo(&issueBo, issueUpdateBo, operatorId)
	if err2 != nil {
		return err2
	}

	//组装更新参与者信息
	updateParticipantErr := assemblyUpdateParticipant(&issueBo, &issueUpdateBo, operatorId, beforeParticipantIds, &afterParticipantIds)

	if updateParticipantErr != nil {
		return updateParticipantErr
	}

	followError := assemblyUpdateFollow(&issueBo, &issueUpdateBo, operatorId, beforeFollowerIds, &afterFollowerIds)

	if followError != nil {
		return followError
	}

	//更新备注
	if issueUpdateBo.IssueDetailRemark != nil {
		err5 := UpdateIssueDetailRemark(issueBo, operatorId, *issueUpdateBo.IssueDetailRemark)
		if err5 != nil {
			log.Error(err5)
			return errs.BuildSystemErrorInfo(errs.IssueDomainError, err5)
		}
	}

	//更新任务
	upd := issueUpdateBo.IssueUpdateCond
	//只是单纯的做check，upd正常情况下不会为空的
	if upd == nil {
		upd = mysql.Upd{}
	}
	upd[consts.TcUpdator] = operatorId
	err := mysql.UpdateSmart(consts.TableIssue, issueBo.Id, upd)
	if err != nil {
		log.Errorf("mysql.TransUpdateSmart: %q\n", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	insertProjectMembersInputBo := &bo.InsertProjectMembersInputBo{
		OrgId:      issueBo.OrgId,
		ProjectId:  issueBo.ProjectId,
		OwnerInfo:  issueBo.OwnerInfo,
		OperatorId: issueBo.Creator,
	}
	if issueBo.ParticipantInfos != nil {
		insertProjectMembersInputBo.ParticipantInfos = *issueBo.ParticipantInfos
	}
	if issueBo.FollowerInfos != nil {
		insertProjectMembersInputBo.FollowerInfos = *issueBo.FollowerInfos
	}

	//最后更新项目成员，可避免事务
	err = InsertProjectMembers(insertProjectMembersInputBo)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	asyn.Execute(func() {
		afterOwnerId := issueBo.Owner
		if issueUpdateBo.OwnerId != nil {
			afterOwnerId = *issueUpdateBo.OwnerId
		}
		pushType := consts.PushTypeUpdateIssue
		if issueUpdateBo.UpdateParticipant || issueUpdateBo.UpdateFollower || issueUpdateBo.NewIssueBo.Owner != issueUpdateBo.IssueBo.Owner {
			pushType = consts.PushTypeUpdateIssueMembers
		}
		if pushType == consts.PushTypeUpdateIssue && len(changeList) == 0 {
			return
		}
		ext := bo.TrendExtensionBo{}
		ext.IssueType = "T"
		ext.ObjName = issueBo.Title
		ext.ChangeList = changeList

		//最新的计划时间
		var planStartTime *types.Time = nil
		var planEndTime *types.Time = nil
		if newIssueBo.PlanStartTime.IsNotNull() {
			planStartTime = &newIssueBo.PlanStartTime
		}
		if newIssueBo.PlanEndTime.IsNotNull() {
			planEndTime = &newIssueBo.PlanEndTime
		}

		issueTrendsBo := bo.IssueTrendsBo{
			PushType:      pushType,
			OrgId:         orgId,
			OperatorId:    operatorId,
			IssueId:       issueBo.Id,
			ParentIssueId: issueBo.ParentId,
			ProjectId:     issueBo.ProjectId,
			PriorityId:    issueBo.PriorityId,
			ParentId:      issueBo.ParentId,

			IssueTitle:               newIssueBo.Title,
			IssueStatusId:            issueBo.Status,
			BeforeOwner:              issueBo.Owner,
			AfterOwner:               afterOwnerId,
			BeforeChangeFollowers:    *beforeFollowerIds,
			AfterChangeFollowers:     afterFollowerIds,
			BeforeChangeParticipants: *beforeParticipantIds,
			AfterChangeParticipants:  afterParticipantIds,

			IssuePlanStartTime: planStartTime,
			IssuePlanEndTime:   planEndTime,

			SourceChannel: sourceChannel,

			NewValue: json.ToJsonIgnoreError(newIssueBo),
			OldValue: json.ToJsonIgnoreError(issueBo),
			Ext:      ext,
		}

		asyn.Execute(func() {
			PushIssueTrends(issueTrendsBo)
		})
		asyn.Execute(func() {
			PushIssueThirdPlatformNotice(issueTrendsBo)
		})
		asyn.Execute(func() {
			if sourceChannel == consts.AppSourceChannelFeiShu {
				UpdateCalendarEvent(issueUpdateBo, issueId, orgId, *beforeParticipantIds, afterParticipantIds)
			}
		})
	})
	return nil
}

//更新followers
func assemblyUpdateFollow(issueBo *bo.IssueBo, issueUpdateBo *bo.IssueUpdateBo, operatorId int64, beforeFollowerIds, afterFollowerIds *[]int64) errs.SystemErrorInfo {
	if issueUpdateBo.UpdateFollower {

		if issueUpdateBo.Followers == nil || len(issueUpdateBo.Followers) == 0 {
			//delete
			err5 := DeleteIssueRelation(operatorId, *issueBo, consts.IssueRelationTypeFollower)
			if err5 != nil {
				log.Error(err5)
				return err5
			}
		} else {
			updateFollowerError := updateFollow(issueBo, issueUpdateBo, operatorId, beforeFollowerIds, afterFollowerIds)

			if updateFollowerError != nil {
				return updateFollowerError
			}
		}
	}

	return nil
}

func updateFollow(issueBo *bo.IssueBo, issueUpdateBo *bo.IssueUpdateBo, operatorId int64, beforeFollowerIds, afterFollowerIds *[]int64) errs.SystemErrorInfo {
	orgId := issueBo.OrgId

	followerInfos := &[]bo.IssueUserBo{}

	issueUpdateBo.Followers = slice.SliceUniqueInt64(issueUpdateBo.Followers)
	*afterFollowerIds = issueUpdateBo.Followers
	//dif
	deletedFollowerIds, addedFollowerIds := util.GetDifMemberIds(*beforeFollowerIds, *afterFollowerIds)

	verifyOrgUserFlag := orgfacade.VerifyOrgUsersRelaxed(orgId, addedFollowerIds)
	if !verifyOrgUserFlag {
		log.Error("存在用户组织校验失败")
		return errs.VerifyOrgError
	}

	//删除
	err1 := DeleteIssueRelationByIds(operatorId, *issueBo, consts.IssueRelationTypeFollower, deletedFollowerIds)
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	//新增
	issueRelationBos, err4 := UpdateIssueRelation(operatorId, *issueBo, consts.IssueRelationTypeFollower, addedFollowerIds)
	if err4 != nil {
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err4)
	}

	for _, issueRelationBo := range issueRelationBos {
		*followerInfos = append(*followerInfos, bo.IssueUserBo{
			IssueRelationBo: issueRelationBo,
		})
	}

	issueBo.FollowerInfos = followerInfos

	return nil
}

//afterParticipantIds,beforeParticipantIds   issueBo.ParticipantInfos 需要赋值
func assemblyUpdateParticipant(issueBo *bo.IssueBo, issueUpdateBo *bo.IssueUpdateBo, operatorId int64, beforeParticipantIds, afterParticipantIds *[]int64) errs.SystemErrorInfo {

	if issueUpdateBo.UpdateParticipant {

		if issueUpdateBo.Participants == nil || len(issueUpdateBo.Participants) == 0 {
			//delete
			err5 := DeleteIssueRelation(operatorId, *issueBo, consts.IssueRelationTypeParticipant)
			*afterParticipantIds = issueUpdateBo.Participants
			if err5 != nil {
				log.Error(err5)
				return err5
			}
		} else {
			participantError := updateParticipant(issueBo, issueUpdateBo, operatorId, beforeParticipantIds, afterParticipantIds)

			if participantError != nil {
				return participantError
			}
		}
	}

	return nil
}

//更新参与者
func updateParticipant(issueBo *bo.IssueBo, issueUpdateBo *bo.IssueUpdateBo, operatorId int64, beforeParticipantIds, afterParticipantIds *[]int64) errs.SystemErrorInfo {
	orgId := issueBo.OrgId

	participantInfos := &[]bo.IssueUserBo{}
	issueUpdateBo.Participants = slice.SliceUniqueInt64(issueUpdateBo.Participants)
	*afterParticipantIds = issueUpdateBo.Participants

	//dif
	deletedParticipantIds, addedParticipantIds := util.GetDifMemberIds(*beforeParticipantIds, *afterParticipantIds)

	verifyOrgUserFlag := orgfacade.VerifyOrgUsersRelaxed(orgId, addedParticipantIds)
	if !verifyOrgUserFlag {
		log.Error("存在用户组织校验失败")
		return errs.VerifyOrgError
	}

	//删除
	err1 := DeleteIssueRelationByIds(operatorId, *issueBo, consts.IssueRelationTypeParticipant, deletedParticipantIds)
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	//新增
	issueRelationBos, err4 := UpdateIssueRelation(operatorId, *issueBo, consts.IssueRelationTypeParticipant, addedParticipantIds)
	if err4 != nil {
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err4)
	}
	for _, issueRelationBo := range issueRelationBos {
		*participantInfos = append(*participantInfos, bo.IssueUserBo{
			IssueRelationBo: issueRelationBo,
		})
	}
	issueBo.ParticipantInfos = participantInfos

	return nil
}

//更新负责人，不需要做diff
func assemblyUpdateOwnerInfo(issueBo *bo.IssueBo, issueUpdateBo bo.IssueUpdateBo, operatorId int64) errs.SystemErrorInfo {

	if issueUpdateBo.OwnerId != nil {
		err5 := DeleteIssueRelation(operatorId, *issueBo, consts.IssueRelationTypeOwner)
		if err5 != nil {
			log.Error(err5)
			return err5
		}
		issueRelationBo, err2 := UpdateIssueRelationSingle(operatorId, *issueBo, consts.IssueRelationTypeOwner, *issueUpdateBo.OwnerId)
		if err2 != nil {
			log.Error(err2)
			return errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
		}
		issueBo.OwnerInfo = &bo.IssueUserBo{
			IssueRelationBo: *issueRelationBo,
		}
	}
	return nil
}

func DeleteIssue(issueBo *bo.IssueBo, operatorId int64, sourceChannel string) errs.SystemErrorInfo {
	orgId := issueBo.OrgId
	issueId := issueBo.Id

	beforeParticipantIds, err1 := GetIssueRelationIdsByRelateType(orgId, issueId, consts.IssueRelationTypeParticipant)
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}
	beforeFollowerIds, err1 := GetIssueRelationIdsByRelateType(orgId, issueId, consts.IssueRelationTypeFollower)
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	count, err := mysql.SelectCountByCond(consts.TableIssue, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcParentId: issueId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if count > 0 {
		return errs.BuildSystemErrorInfo(errs.ExistingSubTask)
	}

	conn, err := mysql.GetConnect()
	if err != nil {
		log.Errorf(consts.DBOpenErrorSentence, err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	tx, err := conn.NewTx(nil)
	if err != nil {
		log.Errorf(consts.TxOpenErrorSentence, err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	defer mysql.Close(conn, tx)

	err = mysql.TransUpdateSmart(tx, consts.TableIssue, issueId, mysql.Upd{
		consts.TcUpdator:  operatorId,
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	issueDetail := &po.PpmPriIssueDetail{}
	_, err = mysql.TransUpdateSmartWithCond(tx, issueDetail.TableName(), db.Cond{
		consts.TcOrgId:    operatorId,
		consts.TcIssueId:  issueId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcUpdator:  operatorId,
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	err3 := DeleteAllIssueRelation(tx, operatorId, orgId, issueId)
	if err3 != nil {
		log.Error(err3)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}

	err = tx.Commit()
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	asyn.Execute(func() {
		blank := []int64{}
		issueTrendsBo := bo.IssueTrendsBo{
			PushType:      consts.PushTypeDeleteIssue,
			OrgId:         orgId,
			OperatorId:    operatorId,
			IssueId:       issueBo.Id,
			ParentIssueId: issueBo.ParentId,
			ProjectId:     issueBo.ProjectId,
			PriorityId:    issueBo.PriorityId,
			ParentId:      issueBo.ParentId,

			IssueTitle:               issueBo.Title,
			IssueStatusId:            issueBo.Status,
			BeforeOwner:              issueBo.Owner,
			AfterOwner:               0,
			BeforeChangeFollowers:    *beforeFollowerIds,
			AfterChangeFollowers:     blank,
			BeforeChangeParticipants: *beforeParticipantIds,
			AfterChangeParticipants:  blank,

			SourceChannel: sourceChannel,
		}
		asyn.Execute(func() {
			PushIssueTrends(issueTrendsBo)
		})
		asyn.Execute(func() {
			PushIssueThirdPlatformNotice(issueTrendsBo)
		})
	})

	return nil
}

func ConvertIssueBoToHomeIssueInfo(issueBo bo.IssueBo) (*bo.HomeIssueInfoBo, error) {
	homeIssueInfos, err := ConvertIssueBosToHomeIssueInfos(issueBo.OrgId, []bo.IssueBo{issueBo})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}
	if len(homeIssueInfos) == 0 {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError)
	}
	return &homeIssueInfos[0], nil
}

func GetHomeIssuePriorityInfoBo(orgId, priorityId int64) (*bo.HomeIssuePriorityInfoBo, errs.SystemErrorInfo) {
	priority, err := GetPriorityById(orgId, priorityId)
	if err != nil {
		log.Error(err)
		//return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	} else {
		priorityInfo := ConvertPriorityCacheInfoToHomeIssuePriorityInfo(*priority)
		return priorityInfo, nil
	}
	return &bo.HomeIssuePriorityInfoBo{}, nil
}

func ConvertPriorityCacheInfoToHomeIssuePriorityInfo(priorityCacheInfo bo.PriorityBo) *bo.HomeIssuePriorityInfoBo {
	priorityInfo := &bo.HomeIssuePriorityInfoBo{}
	priorityInfo.ID = priorityCacheInfo.Id
	priorityInfo.Name = priorityCacheInfo.Name
	priorityInfo.FontStyle = priorityCacheInfo.FontStyle
	priorityInfo.BgStyle = priorityCacheInfo.BgStyle
	return priorityInfo
}

func GetHomeIssueStatusInfoBo(orgId, statusId int64) (*bo.HomeIssueStatusInfoBo, errs.SystemErrorInfo) {
	status, err := processfacade.GetProcessStatusRelaxed(orgId, statusId)
	statusInfo := &bo.HomeIssueStatusInfoBo{}
	if err != nil {
		log.Errorf("status %d, err %v", statusId, err)
		//return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	} else {
		return ConvertStatusInfoToHomeIssueStatusInfo(*status), nil
	}
	return statusInfo, nil
}

func ConvertStatusInfoToHomeIssueStatusInfo(statusInfo bo.CacheProcessStatusBo) *bo.HomeIssueStatusInfoBo {
	homeStatusInfo := &bo.HomeIssueStatusInfoBo{}
	homeStatusInfo.ID = statusInfo.StatusId
	homeStatusInfo.Name = statusInfo.Name
	homeStatusInfo.Type = statusInfo.StatusType
	homeStatusInfo.BgStyle = statusInfo.BgStyle
	homeStatusInfo.FontStyle = statusInfo.FontStyle
	return homeStatusInfo
}

func GetHomeProjectInfoBo(orgId, projectId int64) (*bo.HomeIssueProjectInfoBo, errs.SystemErrorInfo) {
	projectCacheInfo, err := LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	projectInfo := ConvertProjectCacheInfoToHomeIssueProjectInfo(*projectCacheInfo)
	return projectInfo, nil
}

func ConvertProjectCacheInfoToHomeIssueProjectInfo(projectCacheInfo bo.ProjectAuthBo) *bo.HomeIssueProjectInfoBo {
	projectInfo := &bo.HomeIssueProjectInfoBo{}
	projectInfo.Name = projectCacheInfo.Name
	projectInfo.ID = projectCacheInfo.Id
	projectInfo.IsFilling = projectCacheInfo.IsFilling
	return projectInfo
}

func ConvertBaseUserInfoToHomeIssueOwnerInfo(baseUserInfo bo.BaseUserInfoBo) *bo.HomeIssueOwnerInfoBo {
	ownerInfo := &bo.HomeIssueOwnerInfoBo{}
	ownerInfo.ID = baseUserInfo.UserId
	ownerInfo.Name = baseUserInfo.Name
	ownerInfo.Avatar = &baseUserInfo.Avatar
	ownerInfo.IsDeleted = baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted
	ownerInfo.IsDisabled = baseUserInfo.OrgUserStatus == consts.AppStatusDisabled
	return ownerInfo
}

func ConvertIssueTagBosToMapGroupByIssueId(issueTagBos []bo.IssueTagBo) maps.LocalMap {
	issueTagBoMap := maps.LocalMap{}
	if issueTagBos != nil && len(issueTagBos) > 0 {
		for _, issueTagBo := range issueTagBos {
			issueId := issueTagBo.IssueId
			if issueTagsInterface, ok := issueTagBoMap[issueId]; ok {
				if issueTags, ok := issueTagsInterface.(*[]bo.IssueTagBo); ok {
					*issueTags = append(*issueTags, issueTagBo)
				}
			} else {
				issueTagBoMap[issueId] = &[]bo.IssueTagBo{issueTagBo}
			}
		}
	}
	return issueTagBoMap
}

func ConvertIssueTagBosToHomeIssueTagBos(issueTagBos []bo.IssueTagBo) []bo.HomeIssueTagInfoBo {
	homeIssueTagBos := make([]bo.HomeIssueTagInfoBo, 0)
	if issueTagBos != nil && len(issueTagBos) > 0 {
		for _, issueTagBo := range issueTagBos {
			homeIssueTagBos = append(homeIssueTagBos, bo.HomeIssueTagInfoBo{
				ID:        issueTagBo.TagId,
				Name:      issueTagBo.TagName,
				FontStyle: issueTagBo.FontStyle,
				BgStyle:   issueTagBo.BgStyle,
			})
		}
	}
	return homeIssueTagBos
}

func UpdateIssueStatusByStatusType(issueBo bo.IssueBo, operatorId int64, nextStatusType int, needModifyChildStatus int, sourceChannel string) errs.SystemErrorInfo {
	orgId := issueBo.OrgId
	projectId := issueBo.ProjectId
	projectObjectTypeId := issueBo.ProjectObjectTypeId
	if nextStatusType > 3 || nextStatusType < 1 {
		log.Error("要更新的状态类型不明确")
		return errs.BuildSystemErrorInfo(errs.IssueStatusUpdateError)
	}
	processId, err := GetProjectProcessId(orgId, projectId, projectObjectTypeId)
	if err != nil {
		log.Error(err)
		return err
	}
	updateStatusInfos, err := processfacade.GetProcessStatusListRelaxed(orgId, processId)
	//updateStatusIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, nextStatusType)
	if err != nil {
		log.Errorf("proxies.GetProcessStatusIdsRelaxed: %c\n", err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}
	if len(*updateStatusInfos) == 0 {
		log.Errorf("组织%d，不存在状态类型为%d的状态", orgId, nextStatusType)
		return errs.BuildSystemErrorInfo(errs.ProcessStatusNotExist)
	}
	//先赋值，防止为空
	updateFirstStatusId := (*updateStatusInfos)[0].StatusId
	for _, statusInfo := range *updateStatusInfos {
		if statusInfo.Category == consts.ProcessStatusCategoryIssue && statusInfo.StatusType == nextStatusType {
			updateFirstStatusId = statusInfo.StatusId
		}
	}
	return UpdateIssueStatus(issueBo, operatorId, updateFirstStatusId, needModifyChildStatus, sourceChannel)
}

func UpdateIssueStatus(issueBo bo.IssueBo, operatorId int64, nextStatusId int64, needModifyChildStatus int, sourceChannel string) errs.SystemErrorInfo {
	orgId := issueBo.OrgId

	if issueBo.Status == nextStatusId {
		log.Error("任务状态更新-要更新的状态和当前状态一致，更新失败")
		return errs.BuildSystemErrorInfo(errs.IssueStatusUpdateError)
	}

	nextStatus, err1 := processfacade.GetProcessStatusByCategoryRelaxed(orgId, nextStatusId, consts.ProcessStatusCategoryIssue)
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.ProcessStatusNotExist, err1)
	}

	upd := mysql.Upd{
		consts.TcStatus:  nextStatus.StatusId,
		consts.TcUpdator: operatorId,
	}
	if nextStatus.StatusType == consts.ProcessStatusTypeCompleted {
		upd[consts.TcEndTime] = times.GetBeiJingTime()
	} else if nextStatus.StatusType == consts.ProcessStatusTypeNotStarted {
		upd[consts.TcStartTime] = consts.BlankTime
		upd[consts.TcEndTime] = consts.BlankTime
	} else if nextStatus.StatusType == consts.ProcessStatusTypeProcessing {
		upd[consts.TcStartTime] = times.GetBeiJingTime()
		upd[consts.TcEndTime] = consts.BlankTime
	}

	//检查是否存在未完成的子任务
	if nextStatus.StatusType == consts.ProcessStatusTypeCompleted {
		finishedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
		if err != nil {
			log.Errorf("proxies.GetProcessStatusIdsRelaxed: %c\n", err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}

		cond := db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcStatus:   db.NotIn(finishedIds),
			consts.TcParentId: issueBo.Id,
		}

		count, err2 := mysql.SelectCountByCond(consts.TableIssue, cond)
		if err2 != nil {
			log.Errorf("mysql.SelectAllByCond: %v\n", err2)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
		}

		if count > 0 {
			if needModifyChildStatus == 1 {
				//如果需要同步更新子任务状态
				//Nico: fix 子任务状态的开始和结束时间同时也更新掉
				_, err2 := mysql.UpdateSmartWithCond(consts.TableIssue, cond, upd)
				if err2 != nil {
					log.Errorf("mysql.UpdateSmartWithCond: %v\n", err2)
					return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
				}
			} else {
				return errs.BuildSystemErrorInfo(errs.ExistingNotFinishedSubTask)
			}
		}
	}

	err := mysql.UpdateSmart(consts.TableIssue, issueBo.Id, upd)
	if err != nil {
		log.Errorf("mysql.UpdateSmart: %c\n", err)
		return errs.BuildSystemErrorInfo(errs.IssueStatusUpdateError, err)
	}

	asyn.Execute(func() {
		issueMembersBo, err := GetIssueMembers(orgId, issueBo.Id)
		if err != nil{
			log.Error(err)
			return
		}

		beforeParticipantIds := issueMembersBo.ParticipantIds
		beforeFollowerIds := issueMembersBo.FollowerIds

		operateObjProperty := consts.TrendsOperObjPropertyNameStatus
		oldValueMap := map[string]interface{}{
			operateObjProperty: issueBo.Status,
		}
		newValueMap := map[string]interface{}{
			operateObjProperty: nextStatusId,
		}

		//状态列表
		statusList, err := processfacade.GetProcessStatusListByCategoryRelaxed(orgId, consts.ProcessStatusCategoryIssue)
		if err != nil {
			log.Error(strs.ObjectToString(err))
			return
		}

		change := bo.TrendChangeListBo{
			Field:     "status",
			FieldName: consts.Status,
		}
		for _, v := range statusList {
			if v.StatusId == issueBo.Status {
				change.OldValue = v.Name
			} else if v.StatusId == nextStatusId {
				change.NewValue = v.Name
			}
		}
		changeList := []bo.TrendChangeListBo{}
		changeList = append(changeList, change)
		ext := bo.TrendExtensionBo{
			IssueType:  "T",
			ObjName:    issueBo.Title,
			ChangeList: changeList,
		}

		issueTrendsBo := bo.IssueTrendsBo{
			PushType:      consts.PushTypeUpdateIssueStatus,
			OrgId:         orgId,
			OperatorId:    operatorId,
			IssueId:       issueBo.Id,
			ParentIssueId: issueBo.ParentId,
			ProjectId:     issueBo.ProjectId,
			PriorityId:    issueBo.PriorityId,
			ParentId:      issueBo.ParentId,

			IssueTitle:               issueBo.Title,
			IssueStatusId:            nextStatusId, //更新后的状态id
			BeforeOwner:              issueBo.Owner,
			AfterOwner:               0,
			BeforeChangeFollowers:    beforeFollowerIds,
			AfterChangeFollowers:     beforeFollowerIds,
			BeforeChangeParticipants: beforeParticipantIds,
			AfterChangeParticipants:  beforeParticipantIds,

			SourceChannel: sourceChannel,

			NewValue:           json.ToJsonIgnoreError(newValueMap),
			OldValue:           json.ToJsonIgnoreError(oldValueMap),
			OperateObjProperty: operateObjProperty,
			Ext:                ext,
		}

		asyn.Execute(func() {
			PushIssueTrends(issueTrendsBo)
		})
		asyn.Execute(func() {
			PushIssueThirdPlatformNotice(issueTrendsBo)
		})
	})
	return nil
}

//返回要切换的项目id， 如果是0，表示不需要切换
//第二个参数是项目对象类型的项目id
func CheckProjectObjectTypeSwitchProject(projectId, projectObjectTypeId, orgId int64) (int64, int64, errs.SystemErrorInfo) {
	//获取项目和项目对象类型关联对象
	processBos, err := GetProjectObjectTypeProcessByCond(projectObjectTypeId, orgId)
	if err != nil {
		log.Error(err)
		return 0, 0, err
	}

	//验证要移动的任务栏是否在当前项目下
	switchProjectId := int64(0)
	projectObjectTypeProjectId := int64(0)
	for _, bo := range *processBos {
		projectObjectTypeProjectId = bo.ProjectId
		if bo.ProjectId != projectId {
			switchProjectId = bo.ProjectId
			break
		}
	}

	return switchProjectId, projectObjectTypeProjectId, nil
}

func UpdateIssueProjectObjectType(orgId, operatorId int64, issueBo bo.IssueBo, projectObjectTypeId int64) errs.SystemErrorInfo {
	newUUID := uuid.NewUuid()
	lockKey := fmt.Sprintf("%s%d", consts.IssueRelateOperationLock, issueBo.Id)
	suc, lockErr := cache.TryGetDistributedLock(lockKey, newUUID)
	if lockErr != nil{
		log.Error(lockErr)
		return errs.TryDistributedLockError
	}
	if suc{
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, newUUID); err != nil{
				log.Error(err)
			}
		}()
	}else{
		//未获取到锁，直接响应错误信息
		return errs.UpdateIssueProjectObjectTypeFail
	}

	if issueBo.ParentId != 0 {
		return errs.CannotMoveChildIssue
	}
	issueId := issueBo.Id
	oldProjectId := issueBo.ProjectId
	//获取要切换的项目id
	switchProjectId, projectObjectTypeProjectId, err1 := CheckProjectObjectTypeSwitchProject(issueBo.ProjectId, projectObjectTypeId, orgId)
	if err1 != nil {
		log.Error(err1)
		return err1
	}
	if projectObjectTypeProjectId == 0 {
		log.Error("项目对象类型的项目id不存在")
		return errs.ProjectTypeProjectObjectTypeNotExist
	}
	childPos := &[]po.PpmPriIssue{}
	//查询是否存在子任务
	err := mysql.SelectAllByCond(consts.TableIssue, db.Cond{
		consts.TcParentId: issueId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, childPos)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	//默认的
	ids := []int64{issueId}
	if len(*childPos) > 0 {
		for _, data := range *childPos {
			ids = append(ids, data.Id)
		}
	}

	//更新带上项目id，防止并发导致脏数据
	upd := mysql.Upd{
		consts.TcUpdator:             operatorId,
		consts.TcProjectObjectTypeId: projectObjectTypeId,
		consts.TcProjectId:           projectObjectTypeProjectId,
	}

	if switchProjectId > 0 {
		log.Infof("任务 %d 准备切换至项目 %d 下", issueId, switchProjectId)
		//校验当前用户有没有该项目的创建权限
		issueBo.ProjectId = switchProjectId
		authErr := AuthProject(orgId, operatorId, switchProjectId, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationCreate)
		if authErr != nil {
			log.Error(authErr)
			return authErr
		}
		//设置要切换的项目id
		upd[consts.TcProjectId] = switchProjectId

	}

	//更新父子任务的工作栏以及项目id
	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return errs.MysqlOperateError
	}
	_, err = conn.Update(consts.TableIssue).Set(upd).Where(db.And(
		db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
		},
		db.Or(
			db.Cond{
				consts.TcId: issueBo.Id,
			},
			db.Cond{
				consts.TcParentId: issueBo.Id,
			},
		),
	)).Exec()
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return errs.IssueRelationUpdateError
	}

	if switchProjectId > 0 {
		//文件/附件/动态皆不转移到新项目
		//删除任务成员和关注人
		//_, updateErr := mysql.UpdateSmartWithCond(consts.TableIssueRelation, db.Cond{
		//	consts.TcOrgId:orgId,
		//	consts.TcIssueId:db.In(ids),
		//	consts.TcIsDelete:consts.AppIsNoDelete,
		//	consts.TcRelationType:db.In([]int64{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower}),
		//}, mysql.Upd{
		//	consts.TcUpdator:operatorId,
		//	consts.TcIsDelete:consts.AppIsDeleted,
		//})
		//if updateErr != nil {
		//	log.Error(updateErr)
		//	return errs.MysqlOperateError
		//}
		memberErr := switchIssueMember(orgId, ids, operatorId, switchProjectId, conn)
		if memberErr != nil {
			log.Error(memberErr)
			return memberErr
		}

		//标签新项目没有则创建，有则使用新项目的
		tagErr := switchIssueTag(ids, switchProjectId, orgId, operatorId)
		if tagErr != nil {
			log.Error(tagErr)
			return tagErr
		}
		//更新飞书日历
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("捕获到的错误：%s", r)
				}
			}()
			calendarErr := SwitchCalendar(orgId, oldProjectId, ids, operatorId, switchProjectId)
			if calendarErr != nil {
				log.Error(calendarErr)
				return
			}
		}()
	}

	asyn.Execute(func() {
		if projectObjectTypeId != 0 && projectObjectTypeId != issueBo.ProjectObjectTypeId {
			projectObjInfo, err := dao.SelectProjectObjectType(db.Cond{consts.TcId: db.In([]int64{issueBo.ProjectObjectTypeId, projectObjectTypeId})})
			if err != nil {
				log.Error(err)
				return
			}
			var oldName, newName string
			for _, objectType := range *projectObjInfo {
				if objectType.Id == issueBo.ProjectObjectTypeId {
					oldName = objectType.Name
				} else if objectType.Id == projectObjectTypeId {
					newName = objectType.Name
				}
			}
			changeList := []bo.TrendChangeListBo{}
			changeList = append(changeList, bo.TrendChangeListBo{
				Field:     "projectObjectType",
				FieldName: consts.ProjectObjectType,
				OldValue:  oldName,
				NewValue:  newName,
			})
			issueTrendsBo := bo.IssueTrendsBo{
				PushType:      consts.PushTypeUpdateIssueProjectObjectType,
				OrgId:         issueBo.OrgId,
				OperatorId:    operatorId,
				IssueId:       issueBo.Id,
				ParentIssueId: issueBo.ParentId,
				ProjectId:     issueBo.ProjectId,
				PriorityId:    issueBo.PriorityId,
				IssueTitle:    issueBo.Title,
				ParentId:      issueBo.ParentId,
				OldValue:      json.ToJsonIgnoreError(bo.ProjectObjectTypeAndProjectIdBo{ProjectObjectTypeId: issueBo.ProjectObjectTypeId, ProjectId: oldProjectId}),
				NewValue:      json.ToJsonIgnoreError(bo.ProjectObjectTypeAndProjectIdBo{ProjectObjectTypeId: projectObjectTypeId, ProjectId: projectObjectTypeProjectId}),

				Ext: bo.TrendExtensionBo{
					ObjName:    issueBo.Title,
					ChangeList: changeList,
				},
			}
			asyn.Execute(func() {
				PushIssueTrends(issueTrendsBo)
			})
			asyn.Execute(func() {
				PushIssueThirdPlatformNotice(issueTrendsBo)
			})
		}
	})
	return nil
}

//负责人直接带过去，其余的如果是目标项目的人则带过去
func switchIssueMember(orgId int64, issueIds []int64, operatorId int64, projectId int64, conn sqlbuilder.Database) errs.SystemErrorInfo {
	//获取目标项目的所有成员
	newMemberList := &[]po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:orgId,
		consts.TcProjectId:projectId,
		consts.TcIsDelete:consts.AppIsNoDelete,
		consts.TcRelationType:db.In(consts.MemberRelationTypeList),
	}, newMemberList)
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}
	newMemberIds := []int64{}
	for _, relation := range *newMemberList {
		newMemberIds = append(newMemberIds, relation.RelationId)
	}

	//获取目标任务的所有负责人
	ownerList := &[]po.PpmPriIssueRelation{}
	err = mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:orgId,
		consts.TcIssueId:db.In(issueIds),
		consts.TcIsDelete:consts.AppIsNoDelete,
		consts.TcRelationType:db.In([]int64{consts.IssueRelationTypeOwner}),
	}, ownerList)
	//需要新增到目标项目的人
	ownerIds := []int64{}
	for _, relation := range *ownerList {
		if ok, _ := slice.Contain(newMemberIds, relation.RelationId); !ok {
			ownerIds = append(ownerIds, relation.RelationId)
		}
	}

	//不属于新项目成员的关注人和参与人直接去除(和负责人相关的不移出，因为会加入新的项目)
	allRelateIds := append(newMemberIds, ownerIds...)
	_, updateErr := mysql.UpdateSmartWithCond(consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:orgId,
		consts.TcIssueId:db.In(issueIds),
		consts.TcIsDelete:consts.AppIsNoDelete,
		consts.TcRelationType:db.In([]int64{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower}),
		consts.TcRelationId:db.NotIn(allRelateIds),
	}, mysql.Upd{
		consts.TcUpdator:operatorId,
		consts.TcIsDelete:consts.AppIsDeleted,
	})
	if updateErr != nil {
		log.Error(updateErr)
		return errs.MysqlOperateError
	}
	//负责人移动到新项目
	moveErr := UpdateProjectRelation(operatorId, orgId, projectId, consts.IssueRelationTypeFollower, ownerIds)
	if moveErr != nil {
		log.Error(moveErr)
		return moveErr
	}
	//属于新项目的关注人和参与人带过去，负责人直接带过去
	upd := mysql.Upd{
		consts.TcUpdator:operatorId,
		consts.TcProjectId:projectId,
	}
	_, err = conn.Update(consts.TableIssueRelation).Set(upd).Where(db.And(
		db.Cond{
			consts.TcOrgId:orgId,
			consts.TcIssueId:db.In(issueIds),
			consts.TcIsDelete:consts.AppIsNoDelete,
		},
		db.Or(
			db.Cond{
				consts.TcRelationType:db.In([]int64{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower}),
				consts.TcRelationId:db.In(allRelateIds),
			},
			db.Cond{
				consts.TcRelationType:db.In([]int64{consts.IssueRelationTypeOwner}),
			},
		),
	)).Exec()
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return errs.IssueRelationUpdateError
	}
	return nil
}

func switchIssueTag(issueIds []int64, projectId int64, orgId int64, operatorId int64) errs.SystemErrorInfo {
	newUUID := uuid.NewUuid()
	lockKey := fmt.Sprintf("%s%d", consts.CreateProjectTagLock, projectId)
	suc, lockErr := cache.TryGetDistributedLock(lockKey, newUUID)
	if lockErr != nil{
		log.Error(lockErr)
		return errs.TryDistributedLockError
	}
	if suc{
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, newUUID); err != nil{
				log.Error(err)
			}
		}()
	}else{
		//未获取到锁，直接响应错误信息
		return errs.CreateTagFail
	}
	//获取目标项目的所有标签
	_, tags, err := GetTagList(db.Cond{
		consts.TcIsDelete:consts.AppIsNoDelete,
		consts.TcProjectId:projectId,
		consts.TcOrgId:orgId,
	}, 0, 0)
	if err != nil {
		log.Error(err)
		return err
	}

	//获取涉及到的任务的标签
	issueTags, err := GetIssueTagsByIssueIds(orgId, issueIds)
	if err != nil {
		log.Error(err)
		return err
	}
	tagMap := maps.NewMap("Name", *tags)

	//需要更新的标签id->任务id
	upd := map[int64][]int64{}
	createMap := []bo.TagBo{}
	createTagByIssueIds := map[string][]int64{}
	for _, tag := range issueTags {
		if _, ok := tagMap[tag.TagName]; ok {
			//如果存在则使用新项目的标签
			temp := tagMap[tag.TagName].(bo.TagBo)
			upd[temp.Id] = append(upd[temp.Id], tag.IssueId)
		} else {
			//如果不存在则新生成
			if _, ok := createTagByIssueIds[tag.TagName]; !ok {
				createMap = append(createMap, bo.TagBo{BgStyle:tag.BgStyle, FontStyle:tag.FontStyle, Name:tag.TagName})
				createTagByIssueIds[tag.TagName] = []int64{tag.IssueId}
			} else {
				createTagByIssueIds[tag.TagName] = append(createTagByIssueIds[tag.TagName], tag.IssueId)
			}
		}
	}

	if len(createMap) > 0 {
		tagInfos, err := InsertTag(orgId, operatorId, projectId, createMap)
		if err != nil {
			log.Error(err)
			return err
		}
		for _, tag := range tagInfos {
			upd[tag.Id] = createTagByIssueIds[tag.Name]
		}
	}

	if len(upd) > 0 {
		//删除旧有标签关联
		err := DeleteIssueTags(orgId, issueIds, operatorId)
		if err != nil {
			log.Error(err)
			return err
		}
		//插入新关联
		var allCount int
		for _, ids := range upd {
			allCount += len(ids)
		}
		if allCount == 0 {
			return nil
		}

		ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueTag, allCount)
		if err != nil {
			log.Errorf("id generate: %q\n", err)
			return errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
		}

		inserts := []interface{}{}
		var i int64 = 0
		for tagId, idsForIssue := range upd {
			for _, id := range idsForIssue {
				inserts = append(inserts, po.PpmPriIssueTag{
					Id:ids.Ids[i].Id,
					IssueId:id,
					OrgId:orgId,
					ProjectId:projectId,
					TagId:tagId,
				})
				i++
			}
		}
		insertErr := mysql.BatchInsert(&po.PpmPriIssueTag{}, inserts)
		if insertErr != nil {
			log.Error(insertErr)
			return errs.MysqlOperateError
		}
	}

	return nil
}

//获取任务关联用户信息
func GetIssueRelationUserInfos(orgId int64, issueIds []int64) (*bo.IssueRelationUserInfosBo, errs.SystemErrorInfo) {
	issueRelation := &[]po.PpmPriIssueRelation{}
	relationErr := mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: db.In([]int{consts.IssueRelationTypeOwner, consts.IssueRelationTypeFollower, consts.IssueRelationTypeParticipant}),
		consts.TcIssueId:      db.In(issueIds),
	}, issueRelation)
	if relationErr != nil {
		log.Error(relationErr)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, relationErr)
	}
	var followers, participants []bo.IssueUserBo
	var owner *bo.IssueUserBo
	participantMap := map[int64]bool{}
	followerMap := map[int64]bool{}
	for _, v := range *issueRelation {
		relationBo := &bo.IssueRelationBo{}
		_ = copyer.Copy(v, relationBo)
		issueUserBo := bo.IssueUserBo{IssueRelationBo: *relationBo}
		if v.RelationType == consts.IssueRelationTypeFollower {
			if _, ok := followerMap[v.RelationId]; !ok {
				followers = append(followers, issueUserBo)
				followerMap[v.RelationId] = true
			}
		} else if v.RelationType == consts.IssueRelationTypeParticipant {
			if _, ok := participantMap[v.RelationId]; !ok {
				participants = append(participants, issueUserBo)
				participantMap[v.RelationId] = true
			}
		} else if v.RelationType == consts.IssueRelationTypeOwner {
			owner = &issueUserBo
		}
	}
	return &bo.IssueRelationUserInfosBo{
		OwnerInfo:        owner,
		ParticipantInfos: followers,
		FollowerInfos:    participants,
	}, nil
}

func ConvertIssueBosToHomeIssueInfos(orgId int64, issueBos []bo.IssueBo) ([]bo.HomeIssueInfoBo, errs.SystemErrorInfo) {
	issuesSize := len(issueBos)
	list := make([]bo.HomeIssueInfoBo, 0)

	if issuesSize == 0 {
		return list, nil
	}

	//转换map
	var statusMap = maps.LocalMap{}
	var priorityMap = maps.LocalMap{}
	var issueChildCountMap = maps.LocalMap{}
	var issueChildFinishedCountMap = maps.LocalMap{}
	var projectMap = maps.LocalMap{}
	var ownerMap = maps.LocalMap{}
	//interface.(*[]bo.IssueTagBo)
	var issueTagsMap = maps.LocalMap{}

	ownerIds := make([]int64, 0)
	projectIds := make([]int64, 0)
	issueIds := make([]int64, 0)
	for _, issueBo := range issueBos {
		ownerId := issueBo.Owner
		projectId := issueBo.ProjectId
		ownerIdExist, _ := slice.Contain(ownerIds, ownerId)
		if !ownerIdExist {
			ownerIds = append(ownerIds, ownerId)
		}
		projectIdExist, _ := slice.Contain(projectIds, projectId)
		if !projectIdExist {
			projectIds = append(projectIds, projectId)
		}
		issueIds = append(issueIds, issueBo.Id)
	}

	handlerFuncList := make([]func(wg *sync.WaitGroup), 0)
	//状态
	handlerFuncList = append(handlerFuncList, func(wg *sync.WaitGroup) {
		defer wg.Add(-1)
		statusList, err := processfacade.GetProcessStatusListByCategoryRelaxed(orgId, consts.ProcessStatusCategoryIssue)
		if err != nil {
			log.Error(err)
			return
		}
		statusMap = maps.NewMap("StatusId", statusList)
	})

	//优先级
	handlerFuncList = append(handlerFuncList, func(wg *sync.WaitGroup) {
		defer wg.Add(-1)
		priorityList, err := GetPriorityList(orgId)
		if err != nil {
			log.Error(err)
			return
		}
		priorityMap = maps.NewMap("Id", priorityList)
	})

	//任务数量
	handlerFuncList = append(handlerFuncList, func(wg *sync.WaitGroup) {
		defer wg.Add(-1)
		//完成状态的ids
		finishedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
		if err != nil {
			logger.GetDefaultLogger().Error(strs.ObjectToString(err))
			return
		}

		conn, err1 := mysql.GetConnect()
		defer func() {
			if conn != nil {
				if err := conn.Close(); err != nil {
					logger.GetDefaultLogger().Info(strs.ObjectToString(err))
				}
			}
		}()
		if err1 != nil {
			log.Error(err1)
			return
		}

		//所有的任务
		issueChildCount := &[]bo.IssueChildCountBo{}

		dbErr := conn.Select(db.Raw("count(*) as count"), db.Raw("parent_id as parentIssueId")).From(consts.TableIssue).Where(db.Cond{
			consts.TcParentId: db.In(issueIds),
			consts.TcIsDelete: consts.AppIsNoDelete,
		}).GroupBy(consts.TcParentId).All(issueChildCount)
		if dbErr != nil {
			log.Error(dbErr)
			return
		}
		issueChildFinishedCount := &[]bo.IssueChildCountBo{}
		dbErr = conn.Select(db.Raw("count(*) as count"), db.Raw("parent_id as parentIssueId")).From(consts.TableIssue).Where(db.Cond{
			consts.TcStatus:   db.In(finishedIds),
			consts.TcParentId: db.In(issueIds),
			consts.TcIsDelete: consts.AppIsNoDelete,
		}).GroupBy(consts.TcParentId).All(issueChildFinishedCount)
		if dbErr != nil {
			log.Error(dbErr)
			return
		}
		issueChildCountMap = maps.NewMap("ParentIssueId", issueChildCount)
		issueChildFinishedCountMap = maps.NewMap("ParentIssueId", issueChildFinishedCount)
	})

	//项目
	handlerFuncList = append(handlerFuncList, func(wg *sync.WaitGroup) {
		defer wg.Add(-1)
		projectInfos, err := GetProjectAuthBoBatch(orgId, projectIds)
		if err != nil {
			log.Error(err)
			return
		}
		projectMap = maps.NewMap("Id", projectInfos)
	})

	//负责人
	handlerFuncList = append(handlerFuncList, func(wg *sync.WaitGroup) {
		defer wg.Add(-1)
		ownerInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(consts.AppSourceChannelDingTalk, orgId, ownerIds)
		if err != nil {
			log.Error(err)
			return
		}
		ownerMap = maps.NewMap("UserId", ownerInfos)
	})

	//任务tags
	handlerFuncList = append(handlerFuncList, func(wg *sync.WaitGroup) {
		defer wg.Add(-1)
		issueTagBos, err := GetIssueTagsByIssueIds(orgId, issueIds)
		if err != nil {
			log.Error(err)
			return
		}
		//interface.(*[]bo.IssueTagBo)
		issueTagsMap = ConvertIssueTagBosToMapGroupByIssueId(issueTagBos)
	})

	var wg sync.WaitGroup
	wg.Add(len(handlerFuncList))

	for _, handlerFunc := range handlerFuncList {
		currentFunc := handlerFunc
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("捕获到的错误：%s", r)
				}
			}()
			currentFunc(&wg)
		}()
	}

	wg.Wait()

	for _, issueBo := range issueBos {
		homeIssueInfoBo := bo.HomeIssueInfoBo{}
		//start
		if statusCacheInfo, ok := statusMap[issueBo.Status]; ok {
			homeIssueInfoBo.Status = ConvertStatusInfoToHomeIssueStatusInfo(statusCacheInfo.(bo.CacheProcessStatusBo))
		}
		if priorityCacheInfo, ok := priorityMap[issueBo.PriorityId]; ok {
			homeIssueInfoBo.Priority = ConvertPriorityCacheInfoToHomeIssuePriorityInfo(priorityCacheInfo.(bo.PriorityBo))
		}
		if projectCacheInfo, ok := projectMap[issueBo.ProjectId]; ok {
			homeIssueInfoBo.Project = ConvertProjectCacheInfoToHomeIssueProjectInfo(projectCacheInfo.(bo.ProjectAuthBo))
		}
		if userCacheInfo, ok := ownerMap[issueBo.Owner]; ok {
			homeIssueInfoBo.Owner = ConvertBaseUserInfoToHomeIssueOwnerInfo(userCacheInfo.(bo.BaseUserInfoBo))
		}

		if issueChildCountInfo, ok := issueChildCountMap[issueBo.Id]; ok {
			homeIssueInfoBo.ChildsNum = issueChildCountInfo.(bo.IssueChildCountBo).Count
		}
		if issueChildCountInfo, ok := issueChildFinishedCountMap[issueBo.Id]; ok {
			homeIssueInfoBo.ChildsFinishedNum = issueChildCountInfo.(bo.IssueChildCountBo).Count
		}

		if issueTagsInterface, ok := issueTagsMap[issueBo.Id]; ok {
			if issueTagBos, ok := issueTagsInterface.(*[]bo.IssueTagBo); ok {
				homeIssueInfoBo.Tags = ConvertIssueTagBosToHomeIssueTagBos(*issueTagBos)
			}
		}

		homeIssueInfoBo.Issue = issueBo
		list = append(list, homeIssueInfoBo)
	}
	return list, nil
}

func GetIssueBo(orgId, issueId int64) (*bo.IssueBo, errs.SystemErrorInfo) {
	issue := &po.PpmPriIssue{}
	err := mysql.SelectOneByCond(issue.TableName(), db.Cond{
		consts.TcId:       issueId,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, issue)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIssue, err)
	}
	issueBo := &bo.IssueBo{}
	err1 := copyer.Copy(issue, issueBo)
	if err1 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}
	return issueBo, nil
}

func GetCountIssueByProjectObjectTypeId(projectObjecTypeId, projectId, orgId int64) (bool, errs.SystemErrorInfo) {
	issue := &po.PpmPriIssue{}
	count, err := mysql.SelectCountByCond(issue.TableName(), db.Cond{
		consts.TcProjectObjectTypeId: projectObjecTypeId,
		consts.TcProjectId:           projectId,
		consts.TcOrgId:               orgId,
		consts.TcIsDelete:            consts.AppIsNoDelete,
	})

	if err != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	result := int(count) > 0
	return result, nil
}

func GetIssueAuthBo(issueBo bo.IssueBo, currentUserId int64) (*bo.IssueAuthBo, errs.SystemErrorInfo) {
	//orgId := issueBo.OrgId
	issueId := issueBo.Id

	//暂时去掉任务特殊角色
	//issueRelation := &[]po.PpmPriIssueRelation{}
	//relationErr := mysql.SelectAllByCond(consts.TableIssueRelation, db.Cond{
	//	consts.TcOrgId:        orgId,
	//	consts.TcIsDelete:     consts.AppIsNoDelete,
	//	consts.TcRelationType: db.In([]int{consts.IssueRelationTypeFollower, consts.IssueRelationTypeParticipant}),
	//	consts.TcIssueId:      issueId,
	//}, issueRelation)
	//if relationErr != nil {
	//	log.Error(relationErr)
	//	return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, relationErr)
	//}
	//var followers, participants []int64
	//for _, v := range *issueRelation {
	//	if v.RelationType == consts.IssueRelationTypeFollower {
	//		followers = append(followers, v.RelationId)
	//	} else {
	//		participants = append(participants, v.RelationId)
	//	}
	//}
	issueAuthBo := &bo.IssueAuthBo{
		Id:        issueId,
		Owner:     issueBo.Owner,
		Creator:   issueBo.Creator,
		ProjectId: issueBo.ProjectId,
		Status:    issueBo.Status,
		//Followers:    followers,
		//Participants: participants,
	}

	//2019-12-24-nico: 产品想要支持父任务负责人操作所有子任务
	if issueBo.ParentId != 0 {
		parentIssueBo, err := GetIssueBo(issueBo.OrgId, issueBo.ParentId)
		if err != nil {
			log.Error(err)
		} else {
			if parentIssueBo.Owner == currentUserId {
				issueAuthBo.Owner = currentUserId
			}
		}
	}
	//err1 := copyer.Copy(issue, issueAuthBo)
	//if err1 != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	//}
	return issueAuthBo, nil
}

func GetIssueBoList(issueListCond bo.IssueBoListCond) ([]bo.IssueBo, errs.SystemErrorInfo) {
	issueList := &po.PpmPriIssue{}

	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    issueListCond.OrgId,
	}
	if issueListCond.ProjectId != nil {
		cond[consts.TcProjectId] = *issueListCond.ProjectId
	}
	if issueListCond.IterationId != nil {
		cond[consts.TcIterationId] = *issueListCond.IterationId
	}
	if issueListCond.Ids != nil {
		cond[consts.TcId] = db.In(issueListCond.Ids)
	}
	err1 := mysql.SelectAllByCond(consts.TableIssue, cond, issueList)

	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bos := &[]bo.IssueBo{}
	err2 := copyer.Copy(issueList, bos)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *bos, nil
}

func GetIssueAndDetailUnionBoList(issueListCond bo.IssueBoListCond) ([]bo.IssueAndDetailUnionBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	conn.SetLogging(true)
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	orgId := issueListCond.OrgId

	issueAlias := "issue."
	issueRelationAlias := "issueRelation."
	cond := db.Cond{
		issueAlias + consts.TcId + " ":         db.Raw("issueDetail.issue_id"),
		issueAlias + consts.TcIsDelete:         consts.AppIsNoDelete,
		issueAlias + consts.TcId:               db.Raw("issueRelation.issue_id"),
		issueRelationAlias + consts.TcIsDelete: consts.AppIsNoDelete,
		issueAlias + consts.TcOrgId:            issueListCond.OrgId,
	}

	if issueListCond.ProjectId != nil {
		cond[issueAlias+consts.TcProjectId] = *issueListCond.ProjectId
	}
	if issueListCond.IterationId != nil {
		cond[issueAlias+consts.TcIterationId] = *issueListCond.IterationId
	}
	//如果没有项目和迭代这两个条件，只能统计未归档的项目下的任务
	if issueListCond.ProjectId == nil && issueListCond.IterationId == nil {
		//IssueCondFiling(cond, orgId, consts.AppIsNotFilling)
		cond[issueAlias+consts.TcProjectId] = db.In(db.Raw("select id from ppm_pro_project where is_filing = 2 and org_id = ? and is_delete = 2", orgId))
	}
	if issueListCond.RelationType != nil && issueListCond.UserId != nil {
		if *issueListCond.RelationType == 1 {
			//我负责的
			cond[issueRelationAlias+consts.TcRelationId] = issueListCond.UserId
			cond[issueRelationAlias+consts.TcRelationType] = consts.IssueRelationTypeOwner
		} else if *issueListCond.RelationType == 2 {
			//我参与的
			cond[issueRelationAlias+consts.TcRelationId] = issueListCond.UserId
			cond[issueRelationAlias+consts.TcRelationType] = consts.IssueRelationTypeParticipant
		} else if *issueListCond.RelationType == 3 {
			//我关注的
			cond[issueRelationAlias+consts.TcRelationId] = issueListCond.UserId
			cond[issueRelationAlias+consts.TcRelationType] = consts.IssueRelationTypeFollower
		} else if *issueListCond.RelationType == 4 {
			//我发起的
			cond[issueAlias+consts.TcCreator] = issueListCond.UserId
		}
	}

	unionBoList := &[]bo.IssueAndDetailUnionBo{}
	err1 := conn.Select(db.Raw("distinct(issue.id) as issueId,issue.status as issueStatusId,issueDetail.story_point as storyPoint,issue.project_object_type_id as issueProjectObjectTypeId,issue.plan_end_time as planEndTime,issue.end_time as endTime,issue.owner_change_time as ownerChangeTime,issue.create_time as createTime")).
		From(consts.TableIssue+" issue", consts.TableIssueDetail+" issueDetail", consts.TableIssueRelation+" issueRelation").
		Where(cond).
		All(unionBoList)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}

	return *unionBoList, nil
}

func GetIssueInfoList(issueIds []int64) ([]bo.IssueBo, errs.SystemErrorInfo) {
	issueList := &[]po.PpmPriIssue{}

	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(issueIds),
	}
	err1 := mysql.SelectAllByCond(consts.TableIssue, cond, issueList)

	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bos := &[]bo.IssueBo{}
	err2 := copyer.Copy(issueList, bos)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *bos, nil
}

func GetToNowCondIssueCount(projectId int64, timePoint time.Time, processStatus []int64) (int, errs.SystemErrorInfo) {
	issue := &po.PpmPriIssue{}

	//获取当天00:00
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	count, err := mysql.SelectCountByCond(issue.TableName(), db.Cond{
		consts.TcProjectId: projectId,
		consts.TcEndTime:   db.Between(startTime, timePoint),
		consts.TcStatus:    db.In(processStatus),
		consts.TcIsDelete:  consts.AppIsNoDelete,
	})

	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	return int(count), nil
}

func GetUnFinishIssueCount(projectId int64, processStatus []int64) (int, errs.SystemErrorInfo) {

	issue := &po.PpmPriIssue{}

	count, err := mysql.SelectCountByCond(issue.TableName(), db.Cond{
		consts.TcStatus:    db.In(processStatus),
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
	})

	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	return int(count), nil
}

func GetOverdueCondIssueCount(projectId int64, timePoint time.Time, processStatus []int64) (int, errs.SystemErrorInfo) {
	issue := &po.PpmPriIssue{}

	count, err := mysql.SelectCountByCond(issue.TableName(), db.Cond{
		consts.TcProjectId:   projectId,
		consts.TcStatus:      db.In(processStatus),
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcPlanEndTime: db.Before(timePoint),
	})

	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	return int(count), nil
}
