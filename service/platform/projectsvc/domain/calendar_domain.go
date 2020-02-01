package domain

import (
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/extra/feishu"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	consts2 "github.com/galaxy-book/feishu-sdk-golang/core/consts"
	vo2 "github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
	"strconv"
	"time"
	"upper.io/db.v3"
)

//更新日历订阅者
func UpdateCalendarAttendees(orgId int64, calendarId string, addMembers []int64, delMembers []int64, projectId int64) {
	if calendarId == "" {
		projectCalendarInfo, err := GetProjectCalendarInfo(orgId, projectId)
		if err != nil {
			log.Error(err)
			return
		}
		if projectCalendarInfo.IsSyncOutCalendar != consts.IsSyncOutCalendar || projectCalendarInfo.CalendarId == consts.BlankString {
			log.Error("无对应项目日历或未设置导出日历")
			return
		}
		calendarId = projectCalendarInfo.CalendarId
	}
	orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(consts.AppSourceChannelFeiShu, orgId)
	if err != nil {
		log.Errorf("组织外部信息不存在 %v", err)
		return
	}
	tenant, err := feishu.GetTenant(orgBaseInfo.OutOrgId)
	if err != nil {
		log.Error(err)
		return
	}
	if len(addMembers) > 0 {
		resp, err := orgfacade.GetBaseUserInfoBatchRelaxed(consts.AppSourceChannelFeiShu, orgId, addMembers)
		if err != nil {
			log.Error(err)
			return
		}
		for _, v := range resp {
			_, fsErr := tenant.AddCalendarAttendeesAcl(calendarId, vo2.AddCalendarAttendeesAclReq{
				//目前相关人员都给最大权限
				Role: consts2.AccessRoleReader,
				Scope: vo2.CalendarScope{
					Type:   "user",
					OpenId: v.OutUserId,
				},
			})
			if fsErr != nil {
				log.Error(fsErr)
				return
			}
		}
	}

	if len(delMembers) > 0 {
		resp, err := orgfacade.GetBaseUserInfoBatchRelaxed(consts.AppSourceChannelFeiShu, orgId, delMembers)
		if err != nil {
			log.Error(err)
			return
		}
		for _, v := range resp {
			_, fsErr := tenant.DeleteCalendarAttendeesAcl(calendarId, v.OutUserId)
			if fsErr != nil {
				log.Error(fsErr)
				return
			}
		}
	}

}

//创建日历
func CreateCalendar(isSyncOutCalendar *int, orgId int64, projectId int64, userId int64, addMemberIds []int64) {
	if isSyncOutCalendar != nil && *isSyncOutCalendar == consts.IsSyncOutCalendar {
		//防止重复插入
		uid := uuid.NewUuid()
		projectIdStr := strconv.FormatInt(projectId, 10)
		lockKey := consts.CreateCalendarLock + projectIdStr
		suc, err := cache.TryGetDistributedLock(lockKey, uid)
		if err != nil {
			log.Errorf("获取%s锁时异常 %v", lockKey, err)
			return
		}
		if suc {
			defer func() {
				if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
					log.Error(err)
				}
			}()
		}
		projectCalendarInfo, err := GetProjectCalendarInfo(orgId, projectId)
		if err != nil {
			log.Error(err)
			return
		}
		if projectCalendarInfo.CalendarId != "" {
			log.Info("日历已插入" + projectIdStr)
			return
		}
		orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(consts.AppSourceChannelFeiShu, orgId)
		if err != nil {
			log.Errorf("组织外部信息不存在 %v", err)
		}
		projectInfo, err := GetProject(orgId, projectId)
		if err != nil {
			log.Error(err)
			return
		}
		//创建日历
		tenant, err := feishu.GetTenant(orgBaseInfo.OutOrgId)
		if err != nil {
			log.Error(err)
			return
		}
		resp, fsErr := tenant.CreateCalendar(vo2.CreateCalendarReq{
			Summary:           projectInfo.Name,
			Description:       projectInfo.Remark,
			DefaultAccessRole: consts2.AccessRoleReader,
		})
		if fsErr != nil {
			log.Error(fsErr)
			return
		}

		if resp.Code != 200000 {
			log.Error("创建日历失败" + resp.Msg)
			return
		}
		//创建关联关系
		memberId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectRelation)
		if err != nil {
			log.Error(err)
			return
		}
		err1 := dao.InsertProjectRelation(po.PpmProProjectRelation{
			Id:           memberId,
			OrgId:        orgId,
			ProjectId:    projectId,
			RelationType: consts.IssueRelationTypeCalendar,
			RelationCode: resp.Data.Id,
			Creator:      userId,
			CreateTime:   time.Now(),
			IsDelete:     consts.AppIsNoDelete,
			Status:       consts.ProjectMemberEffective,
			Updator:      userId,
			UpdateTime:   time.Now(),
			Version:      1,
		})
		if err1 != nil {
			log.Error(err1)
		}
		_, err2 := GetProjectCalendarInfo(orgId, projectId)
		if err2 != nil {
			log.Error("缓存信息失败" + err2.Message())
		}
		//访问控制
		UpdateCalendarAttendees(orgId, resp.Data.Id, addMemberIds, []int64{}, projectId)

		//创建日程（可能创建日历之前就有任务）
		createCalendarEventsBefore(projectId, userId)
	}
}

//创建日程（可能创建日历之前就有任务）
func createCalendarEventsBefore(projectId, userId int64)  {
	issuesInfo, count,err := SelectList(db.Cond{
		consts.TcProjectId:projectId,
		consts.TcIsDelete:consts.AppIsNoDelete,
	}, nil, 0, 0, nil)
	if err != nil {
		log.Error(err)
		return
	}
	if count == 0 {
		return
	}
	issueIds := []int64{}
	for _, issueBo := range *issuesInfo {
		issueIds = append(issueIds, issueBo.Id)
	}

	//获取已经创建过日程的任务
	relationInfo, err1 := GetRelationInfoByIssueIds(issueIds, []int{consts.IssueRelationTypeCalendar})
	if err1 != nil {
		log.Error(err1)
		return
	}

	alreadyCreateIds := []int64{}
	for _, relationBo := range relationInfo {
		alreadyCreateIds = append(alreadyCreateIds, relationBo.IssueId)
	}

	needCreateIds := []int64{}
	for _, id := range issueIds {
		if ok, _ := slice.Contain(alreadyCreateIds, id); !ok {
			needCreateIds = append(needCreateIds, id)
		}
	}
	if len(needCreateIds) == 0 {
		return
	}

	//获取参与人
	participants, err1 := GetRelationInfoByIssueIds(issueIds, []int{consts.IssueRelationTypeParticipant})
	if err1 != nil {
		log.Error(err1)
		return
	}

	participantMap := map[int64][]int64{}
	for _, participant := range participants {
		participantMap[participant.IssueId] = append(participantMap[participant.IssueId], participant.RelationId)
	}

	for _, issueBo := range *issuesInfo {
		if ok, _ := slice.Contain(needCreateIds, issueBo.Id); ok {
			temp := []int64{}
			if _, ok := participantMap[issueBo.Id]; ok {
				temp = participantMap[issueBo.Id]
			}
			CreateCalendarEvent(&issueBo, userId, temp)
		}
	}
}
//更新日历
func UpdateCalendar(input vo.UpdateProjectReq, orgId int64, oldMembers, newMembers []int64, userId int64) {
	projectCalendarInfo, err := GetProjectCalendarInfo(orgId, input.ID)
	if err != nil {
		log.Error(err)
		return
	}
	if projectCalendarInfo.IsSyncOutCalendar != consts.IsSyncOutCalendar {
		log.Error("未设置导出日历")
		return
	}
	if projectCalendarInfo.CalendarId == consts.BlankString {
		log.Info("无对应项目日历,重新生成日历")
		CreateCalendar(&projectCalendarInfo.IsSyncOutCalendar, orgId, input.ID, userId, newMembers)
		return
	}

	update := vo2.UpdateCalendarReq{}
	needUpdate := 0
	for _, v := range input.UpdateFields {
		if v == "name" && input.Name != nil {
			update.Summary = *input.Name
			needUpdate = 1
		}
		if v == "remark" && input.Remark != nil {
			update.Description = *input.Remark
			needUpdate = 1
		}
	}

	orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(consts.AppSourceChannelFeiShu, orgId)
	if err != nil {
		log.Errorf("组织外部信息不存在 %v", err)
		return
	}
	if needUpdate != 0 {
		//更新日历
		tenant, err := feishu.GetTenant(orgBaseInfo.OutOrgId)
		if err != nil {
			log.Error(err)
		}
		resp, fsErr := tenant.UpdateCalendar(projectCalendarInfo.CalendarId, update)
		if fsErr != nil {
			log.Error(fsErr)
		}

		if resp.Code != 200000 {
			log.Error("更新日历失败" + resp.Msg)
			return
		}
		log.Info("更新日历成功")
		return
	} else {
		log.Info("日历无需更新")
	}

	//创建成员订阅日历
	deleted, added := util.GetDifMemberIds(oldMembers, newMembers)
	UpdateCalendarAttendees(orgId, projectCalendarInfo.CalendarId, added, deleted, input.ID)
}

//创建日程
func CreateCalendarEvent(issueBo *bo.IssueBo, userId int64, participant []int64) {
	//防止重复插入
	uid := uuid.NewUuid()
	issueIdStr := strconv.FormatInt(issueBo.Id, 10)
	lockKey := consts.CreateCalendarEventLock + issueIdStr
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("获取%s锁时异常 %v", lockKey, err)
		return
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	}
	//如果原来没有则新建
	_, err = dao.SelectOneIssueRelation(db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: consts.IssueRelationTypeCalendar,
		consts.TcIssueId:      issueBo.Id,
		consts.TcOrgId:        issueBo.OrgId,
	})
	if err != nil {
		if err != db.ErrNoMoreRows {
			log.Error(err)
			return
		}
	} else {
		log.Info("日程已创建" + issueIdStr)
		return
	}
	defaultTime := types.Time(consts.BlankTimeObject)
	if issueBo.PlanStartTime == defaultTime || issueBo.PlanEndTime == defaultTime {
		log.Error("日程创建失败：开始时间和结束时间必填")
		return
	}
	ok, outOrgId, calendarId := GetCalendarInfo(issueBo.OrgId, issueBo.ProjectId)
	if !ok {
		return
	}
	attendeesId := []int64{}
	attendeesId = append(append(attendeesId, issueBo.Owner), participant...)
	//if issueBo.ParticipantInfos != nil {
	//	for _, v := range *issueBo.ParticipantInfos {
	//		attendeesId = append(attendeesId, v.RelationId)
	//	}
	//}
	//关注人不加入日程参与人
	//if issueBo.FollowerInfos != nil {
	//	for _, v := range *issueBo.FollowerInfos {
	//		attendeesId = append(attendeesId, v.RelationId)
	//	}
	//}
	attendeesId = slice.SliceUniqueInt64(attendeesId)
	if len(attendeesId) == 0 {
		log.Info("日程创建失败：无相关人员")
		return
	}
	resp, err := orgfacade.GetBaseUserInfoBatchRelaxed(consts.AppSourceChannelFeiShu, issueBo.OrgId, attendeesId)
	if err != nil {
		log.Error(err)
		return
	}
	attendees := []vo2.Attendees{}
	for _, v := range resp {
		attendees = append(attendees, vo2.Attendees{
			OpenId:      v.OutUserId,
			DisplayName: v.Name,
		})
	}

	//创建日程
	tenant, err := feishu.GetTenant(outOrgId)
	if err != nil {
		log.Error(err)
	}
	start, _ := time.ParseInLocation(consts.AppTimeFormat, issueBo.PlanStartTime.String(), time.Local)
	end, _ := time.ParseInLocation(consts.AppTimeFormat, issueBo.PlanEndTime.String(), time.Local)
	createReq := vo2.CreateCalendarEventReq{
		Summary: issueBo.Title,
		Start: vo2.TimeFormat{
			TimeStamp: start.Unix(),
		},
		End: vo2.TimeFormat{
			TimeStamp: end.Unix(),
		},
		Attendees: &attendees,
	}
	if issueBo.IssueDetailBo.Remark != consts.BlankString {
		createReq.Description = issueBo.IssueDetailBo.Remark
	}

	resp1, fsErr := tenant.CreateCalendarEvent(calendarId, createReq)
	if fsErr != nil {
		log.Error(fsErr)
	}

	if resp1.Code != 200000 {
		log.Error("日程创建失败：" + resp1.Msg)
	}
	log.Info("创建日程成功")

	//创建关联关系
	relationId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueRelation)
	if err != nil {
		log.Error(err)
		return
	}
	err1 := dao.InsertIssueRelation(po.PpmPriIssueRelation{
		Id:           relationId,
		OrgId:        issueBo.OrgId,
		IssueId:      issueBo.Id,
		RelationType: consts.IssueRelationTypeCalendar,
		RelationCode: resp1.Data.Id,
		Creator:      userId,
		CreateTime:   time.Now(),
		Updator:      userId,
		UpdateTime:   time.Now(),
		IsDelete:     consts.AppIsNoDelete,
	})
	if err1 != nil {
		log.Error(err1)
		return
	}
}

//更新日程
func UpdateCalendarEvent(issueUpdateBo bo.IssueUpdateBo, issueId, orgId int64, beforeParticipants, afterParticipants []int64) {
	ok, outOrgId, calendarId := GetCalendarInfo(orgId, issueUpdateBo.IssueBo.ProjectId)
	if !ok {
		return
	}
	//如果原来没有则新建
	relation, err := dao.SelectOneIssueRelation(db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: consts.IssueRelationTypeCalendar,
		consts.TcIssueId:      issueId,
		consts.TcOrgId:        orgId,
	})
	if err != nil {
		log.Info("更新时生成日程")
		CreateCalendarEvent(&issueUpdateBo.NewIssueBo, issueUpdateBo.OperatorId, afterParticipants)
		return
	}
	//如果有则更新（基本信息，订阅者）
	//更新日程
	defaultTime := types.Time(consts.BlankTimeObject)
	if issueUpdateBo.NewIssueBo.PlanStartTime == defaultTime || issueUpdateBo.NewIssueBo.PlanEndTime == defaultTime {
		log.Error("时间不能为空")
		return
	}
	tenant, err := feishu.GetTenant(outOrgId)
	if err != nil {
		log.Error(err)
	}
	start, _ := time.ParseInLocation(consts.AppTimeFormat, issueUpdateBo.NewIssueBo.PlanStartTime.String(), time.Local)
	end, _ := time.ParseInLocation(consts.AppTimeFormat, issueUpdateBo.NewIssueBo.PlanEndTime.String(), time.Local)
	createReq := vo2.CreateCalendarEventReq{
		Summary: issueUpdateBo.NewIssueBo.Title,
		Start: vo2.TimeFormat{
			TimeStamp: start.Unix(),
		},
		End: vo2.TimeFormat{
			TimeStamp: end.Unix(),
		},
	}
	if issueUpdateBo.IssueDetailRemark != nil {
		createReq.Description = *issueUpdateBo.IssueDetailRemark
	}

	resp1, fsErr := tenant.UpdateCalendarEvent(calendarId, relation.RelationCode, createReq)
	if fsErr != nil {
		log.Error(fsErr)
	}

	if resp1.Code != 200000 {
		log.Error("更新日程失败" + resp1.Msg)
	}
	//更新订阅者
	deletedParticipantIds, addedParticipantIds := util.GetDifMemberIds(beforeParticipants, afterParticipants)
	if issueUpdateBo.IssueBo.Owner != issueUpdateBo.NewIssueBo.Owner {
		deletedParticipantIds = append(deletedParticipantIds, issueUpdateBo.IssueBo.Owner)
		addedParticipantIds = append(addedParticipantIds, issueUpdateBo.NewIssueBo.Owner)
	}
	attendees := []vo2.AttendeesResp{}
	attendeesId := slice.SliceUniqueInt64(append(deletedParticipantIds, addedParticipantIds...))
	if len(attendeesId) == 0 {
		log.Info("无相关人员")
		return
	}
	resp, err := orgfacade.GetBaseUserInfoBatchRelaxed(consts.AppSourceChannelFeiShu, orgId, attendeesId)
	if err != nil {
		log.Error(err)
		return
	}
	newMap := maps.NewMap("UserId", resp)
	for _, v := range deletedParticipantIds {
		info, ok := newMap[v]
		if !ok {
			continue
		}
		baseUserInfo := info.(bo.BaseUserInfoBo)
		attendees = append(attendees, vo2.AttendeesResp{
			Attendees: vo2.Attendees{
				OpenId:      baseUserInfo.OutUserId,
				DisplayName: baseUserInfo.Name,
			},
			Status: consts2.ActionRemove,
		})
	}
	for _, v := range addedParticipantIds {
		info, ok := newMap[v]
		if !ok {
			continue
		}
		baseUserInfo := info.(bo.BaseUserInfoBo)
		attendees = append(attendees, vo2.AttendeesResp{
			Attendees: vo2.Attendees{
				OpenId:      baseUserInfo.OutUserId,
				DisplayName: baseUserInfo.Name,
			},
			Status: consts2.ActionInvite,
		})
	}

	resp2, fsErr := tenant.UpdateCalendarEventAttendees(calendarId, relation.RelationCode, vo2.UpdateCalendarEventAtendeesReq{
		Attendees: attendees,
	})
	if fsErr != nil {
		log.Error(fsErr)
	}

	if resp2.Code != 200000 {
		log.Error("更新日程订阅者失败" + resp2.Msg)
	}
	log.Info("更新日程成功")
}

//中途同步日历
func SyncCalendarConfirm(orgId, userId, projectId int64) {
	projectCalendarInfo, err := GetProjectCalendarInfo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return
	}
	if projectCalendarInfo.CalendarId != "" {
		log.Info("项目已同步日历")
		createCalendarEventsBefore(projectId, userId)
		return
	}
	info, err := GetProjectRelationByType(projectId, []int64{consts.IssueRelationTypeOwner, consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower})
	if err != nil {
		log.Error(err)
		return
	}
	addIds := []int64{}
	for _, v := range *info {
		addIds = append(addIds, v.RelationId)
	}
	CreateCalendar(&projectCalendarInfo.IsSyncOutCalendar, orgId, projectId, userId, addIds)
}

func SwitchCalendar(orgId, oldProjectId int64, issueIds []int64, operatorId int64, newProjectId int64) errs.SystemErrorInfo {
	relationInfo, err := GetRelationInfoByIssueIds(issueIds, []int{consts.IssueRelationTypeCalendar})
	if err != nil {
		return err
	}
	//删除日程关联
	_, updateErr := mysql.UpdateSmartWithCond(consts.TableIssueRelation, db.Cond{
		consts.TcOrgId:orgId,
		consts.TcIssueId:db.In(issueIds),
		consts.TcIsDelete:consts.AppIsNoDelete,
		consts.TcRelationType:db.In([]int64{consts.IssueRelationTypeCalendar}),
	}, mysql.Upd{
		consts.TcUpdator:operatorId,
		consts.TcIsDelete:consts.AppIsDeleted,
	})
	if updateErr != nil {
		log.Error(updateErr)
		return errs.MysqlOperateError
	}
	eventIds := []string{}
	for _, relationBo := range relationInfo {
		eventIds = append(eventIds, relationBo.RelationCode)
	}

	if len(eventIds) > 0 {
		ok, outOrgId, calendarId := GetCalendarInfo(orgId, oldProjectId)
		if ok {
			tenant, err := feishu.GetTenant(outOrgId)
			if err != nil {
				log.Error(err)
				return err
			}

			//如果有，删除已有日程
			for _, id := range eventIds {
				_, deleteErr := tenant.DeleteCalendarEvent(calendarId, id)
				if deleteErr != nil {
					return errs.BuildSystemErrorInfo(errs.SystemError, deleteErr)
				}
			}
		}
	}

	ok, _, _ := GetCalendarInfo(orgId, newProjectId)
	if !ok {
		return nil
	}
	issuesInfo, issueErr := GetIssueInfoList(issueIds)
	if issueErr != nil {
		log.Error(issueErr)
		return issueErr
	}
	for _, issueBo := range issuesInfo {
		CreateCalendarEvent(&issueBo, operatorId, []int64{})
	}

	return nil
}