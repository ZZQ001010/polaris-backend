package service

import (
	"fmt"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/date"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"strconv"
	"time"
	"upper.io/db.v3"
)

func JudgeProjectFiling(orgId, projectId int64) errs.SystemErrorInfo {
	projectInfo, err := domain.LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return err
	}
	if projectInfo.IsFilling == consts.ProjectIsFiling {
		return errs.ProjectIsFilingYet
	}

	return nil
}

func dealAuthAndInputType(orgId, currentUserId int64, input *vo.CreateIssueReq) errs.SystemErrorInfo {
	err := domain.AuthProject(orgId, currentUserId, input.ProjectID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationCreate)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	pass := orgfacade.VerifyOrgRelaxed(orgId, input.OwnerID)
	if !pass {
		return errs.BuildSystemErrorInfo(errs.IllegalityOwner)
	}

	pass, err = VerifyPriority(orgId, consts.PriorityTypeIssue, input.PriorityID)
	if err != nil || !pass {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.IllegalityPriority)
	}

	//这里做一个小逻辑，默认的查普通项目的项目类型
	if input.TypeID == nil {
		list, err := domain.GetProjectObjectTypeList(orgId, input.ProjectID)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		var defaultType *bo.ProjectObjectTypeBo = nil
		for _, v := range *list {
			if defaultType == nil {
				defaultType = &v
			}
			if v.LangCode == consts.ProjectObjectTypeLangCodeTask {
				input.TypeID = &v.Id
				break
			}
		}
		if input.TypeID == nil {
			defaultTypeId := defaultType.Id
			input.TypeID = &defaultTypeId
		}
	}
	return nil
}

func assemblyIssueBoParam(orgId int64, input *vo.CreateIssueReq) (initStatus int64, projectBo *bo.ProjectBo, issueCode *bo.IdCodes, err errs.SystemErrorInfo) {

	if input.TypeID == nil {
		log.Error("任务类型不存在")
		return 0, nil,nil, errs.BuildSystemErrorInfo(errs.ProjectObjectTypeNotExist)
	}

	initStatus, err1 := processfacade.GetProcessInitStatusIdRelaxed(orgId, input.ProjectID, *input.TypeID, consts.ProcessStatusCategoryIssue)
	if err1 != nil {
		log.Error(err1)
		return 0, nil,nil, err1
	}

	projectBo, err2 := domain.GetProject(orgId, input.ProjectID)
	if err2 != nil {
		log.Error(err2)
		return 0, nil,nil, err2
	}

	//取消对象类型编号
	issueCode, err4 := idfacade.ApplyIdRelaxed(orgId, projectBo.PreCode, "")
	if err4 != nil {
		return 0, nil, nil, err4
	}

	var zero int64 = 0
	if input.IterationID == nil {
		input.IterationID = &zero
	}

	return initStatus, projectBo, issueCode, nil
}

func changeIssueValue(entity *bo.IssueBo, input *vo.CreateIssueReq) errs.SystemErrorInfo{

	if input.ParentID != nil {
		//判断主任务是否存在
		parentIssueBo, err := domain.GetIssueBo(entity.OrgId, *input.ParentID)
		if err != nil{
			log.Error(err)
			return errs.ParentIssueNotExist
		}
		if parentIssueBo.ParentId != 0 {
			return errs.ParentIssueHasParent
		}
		entity.ParentId = *input.ParentID
		entity.ProjectObjectTypeId = parentIssueBo.ProjectObjectTypeId
	}
	if input.PlanEndTime != nil && input.PlanEndTime.IsNotNull() {
		entity.PlanEndTime = *input.PlanEndTime
	} else {
		entity.PlanEndTime = types.Time(consts.BlankTimeObject)
	}

	if input.PlanStartTime != nil && input.PlanStartTime.IsNotNull() {
		entity.PlanStartTime = *input.PlanStartTime
	} else {
		entity.PlanStartTime = types.Time(consts.BlankTimeObject)
	}

	if input.ModuleID != nil {
		entity.ModuleId = *input.ModuleID
	}
	if input.PlanWorkHour != nil {
		entity.PlanWorkHour = *input.PlanWorkHour
	}
	if input.VersionID != nil {
		entity.VersionId = *input.VersionID
	}
	return nil
}

//参与人
func assemblyParticipantIds(input *vo.CreateIssueReq, orgId int64, entity *bo.IssueBo, currentUserId int64) errs.SystemErrorInfo {

	if input.ParticipantIds != nil && len(input.ParticipantIds) > 0 {
		input.ParticipantIds = slice.SliceUniqueInt64(input.ParticipantIds)
		pcount := len(input.ParticipantIds)
		participantInfos := make([]bo.IssueUserBo, 0, pcount)

		pids, err3 := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueRelation, pcount)
		if err3 != nil {
			return err3
		}

		for i, participantId := range input.ParticipantIds {
			pass := orgfacade.VerifyOrgRelaxed(orgId, participantId)
			if !pass {
				return errs.BuildSystemErrorInfo(errs.VerifyOrgError)
			}

			participantInfos = append(participantInfos, bo.IssueUserBo{
				IssueRelationBo: bo.IssueRelationBo{
					Id:           pids.Ids[i].Id,
					OrgId:        entity.OrgId,
					ProjectId:    entity.ProjectId,
					IssueId:      entity.Id,
					RelationId:   participantId,
					RelationType: consts.IssueRelationTypeParticipant,
					Creator:      currentUserId,
					CreateTime:   types.NowTime(),
					Updator:      currentUserId,
					UpdateTime:   types.NowTime(),
					Version:      1,
				},
			})
		}
		entity.ParticipantInfos = &participantInfos
	}
	return nil
}

//关注者
func assembleyFollower(input *vo.CreateIssueReq, orgId int64, entity *bo.IssueBo, currentUserId int64) errs.SystemErrorInfo {
	if input.FollowerIds != nil && len(input.FollowerIds) > 0 {
		input.FollowerIds = slice.SliceUniqueInt64(input.FollowerIds)

		fcount := len(input.FollowerIds)
		followerInfos := make([]bo.IssueUserBo, 0, fcount)

		fids, err3 := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueRelation, fcount)
		if err3 != nil {
			return err3
		}

		for i, followerId := range input.FollowerIds {
			pass := orgfacade.VerifyOrgRelaxed(orgId, followerId)
			if !pass {
				return errs.BuildSystemErrorInfo(errs.VerifyOrgError)
			}

			followerInfos = append(followerInfos, bo.IssueUserBo{
				IssueRelationBo: bo.IssueRelationBo{
					Id:           fids.Ids[i].Id,
					OrgId:        entity.OrgId,
					ProjectId:    entity.ProjectId,
					IssueId:      entity.Id,
					RelationId:   followerId,
					RelationType: consts.IssueRelationTypeFollower,
					Creator:      currentUserId,
					CreateTime:   types.NowTime(),
					Updator:      currentUserId,
					UpdateTime:   types.NowTime(),
					Version:      1,
				},
			})
		}
		entity.FollowerInfos = &followerInfos
	}
	return nil
}

func CreateIssue(reqVo projectvo.CreateIssueReqVo) (*vo.Issue, errs.SystemErrorInfo) {
	return CreateIssueWithId(reqVo, 0)
}

func CreateIssueWithId(reqVo projectvo.CreateIssueReqVo, issueId int64) (*vo.Issue, errs.SystemErrorInfo) {
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	input := reqVo.CreateIssue
	sourceChannel := reqVo.SourceChannel

	err := dealAuthAndInputType(orgId, currentUserId, &input)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	initStatus, projectBo, issueCode, err := assemblyIssueBoParam(orgId, &input)

	//如果外部没有传id，则申请一个新的id
	if issueId <= 0{
		id, err3 := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssue)
		if err3 != nil {
			log.Error(err3)
			return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError)
		}
		issueId = id
	}
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	nowTime := types.NowTime()
	entity := &bo.IssueBo{
		Id:                  issueId,
		OrgId:               orgId,
		Code:                issueCode.Ids[0].Code,
		ProjectId:           input.ProjectID,
		ProjectObjectTypeId: *input.TypeID,
		Title:               input.Title,
		Owner:               input.OwnerID,
		OwnerChangeTime:     nowTime,
		PriorityId:          input.PriorityID,
		SourceId:            0,
		IssueObjectTypeId:   0,
		IterationId:         *input.IterationID,
		ModuleId:            0,
		Status:              initStatus,
		Creator:             currentUserId,
		CreateTime:          nowTime,
		Updator:             currentUserId,
		UpdateTime:          nowTime,
		Version:             1,
	}
	//tags
	entity.Tags = ConvertIssueTagReqVoToBo(input.Tags)

	issueErr := changeIssueValue(entity, &input)
	if issueErr != nil{
		log.Error(issueErr)
		return nil, issueErr
	}

	issueDetailId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueDetail)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}

	issueDetail := &bo.IssueDetailBo{
		Id:         issueDetailId,
		OrgId:      orgId,
		IssueId:    entity.Id,
		ProjectId:  entity.ProjectId,
		StoryPoint: 0,
		Tags:       "",
		Status:     entity.Status,
		Creator:    entity.Creator,
		CreateTime: entity.CreateTime,
		Updator:    entity.Updator,
		UpdateTime: entity.UpdateTime,
		Version:    1,
	}
	if input.Remark != nil {
		issueDetail.Remark = *input.Remark
	}

	entity.IssueDetailBo = *issueDetail

	err = assemblyParticipantIds(&input, orgId, entity, currentUserId)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	err = assembleyFollower(&input, orgId, entity, currentUserId)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	oid, err3 := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueRelation)
	if err3 != nil {
		return nil, err3
	}

	//负责人
	entity.OwnerInfo = &bo.IssueUserBo{
		IssueRelationBo: bo.IssueRelationBo{
			Id:           oid,
			OrgId:        entity.OrgId,
			ProjectId:    entity.ProjectId,
			IssueId:      entity.Id,
			RelationId:   input.OwnerID,
			RelationType: consts.IssueRelationTypeOwner,
			Creator:      currentUserId,
			CreateTime:   types.NowTime(),
			Updator:      currentUserId,
			UpdateTime:   types.NowTime(),
			Version:      1,
		},
	}

	err4 := domain.CreateIssue(entity, sourceChannel)
	if err4 != nil {
		return nil, err4
	}
	//不允许给子任务创建子任务
	if input.Children != nil && len(input.Children) > 0 && input.ParentID == nil {
		err = CreateChildIssue(input, orgId, currentUserId, initStatus, projectBo.PreCode, issueId, sourceChannel)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
		}
	}

	asyn.Execute(func() {
		PushAddIssueNotice(orgId, projectBo.Id, entity.Id, currentUserId)
	})

	issueResp := &vo.Issue{}
	err1 := copyer.Copy(entity, issueResp)
	if err1 != nil {
		log.Errorf("copyer.Copy: %q\n", err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	return issueResp, nil
}

func ConvertIssueTagReqVoToBo(issueTagReqVos []*vo.IssueTagReqInfo) []bo.IssueTagReqBo {
	if issueTagReqVos != nil && len(issueTagReqVos) > 0 {
		issueTagReqBos := make([]bo.IssueTagReqBo, 0)
		for _, tag := range issueTagReqVos {
			if tag != nil {
				issueTagReqBos = append(issueTagReqBos, bo.IssueTagReqBo{
					Id:   tag.ID,
					Name: tag.Name,
				})
			}
		}
		return issueTagReqBos
	}
	return nil
}

func CreateChildIssue(input vo.CreateIssueReq, orgId, currentUserId, initStatus int64, projectPreCode string, parentId int64, sourceChannel string) errs.SystemErrorInfo {

	nowTime := types.NowTime()
	for _, v := range input.Children {

		issueId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssue)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
		}
		issueCode, err2 := idfacade.ApplyIdRelaxed(orgId, projectPreCode, "")
		if err2 != nil {
			return errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
		}
		entity := &bo.IssueBo{
			Id:                  issueId,
			OrgId:               orgId,
			Code:                issueCode.Ids[0].Code,
			ProjectId:           input.ProjectID,
			ProjectObjectTypeId: *input.TypeID,
			Title:               v.Title,
			Owner:               v.OwnerID,
			PriorityId:          v.PriorityID,
			OwnerChangeTime:     nowTime,
			SourceId:            0,
			IssueObjectTypeId:   0,
			IterationId:         0,
			ModuleId:            0,
			Status:              initStatus,
			Creator:             currentUserId,
			CreateTime:          types.NowTime(),
			Updator:             currentUserId,
			UpdateTime:          types.NowTime(),
			Version:             1,
			ParentId:            parentId,
		}
		if v.TypeID != nil {
			entity.ProjectObjectTypeId = *v.TypeID
		}

		assemblyIssueTime(entity, v)

		issueDetailId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueDetail)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
		}
		if v.Remark == nil {
			blank := consts.BlankString
			v.Remark = &blank
		}
		issueDetail := &bo.IssueDetailBo{
			Id:         issueDetailId,
			OrgId:      orgId,
			IssueId:    entity.Id,
			ProjectId:  entity.ProjectId,
			StoryPoint: 0,
			Tags:       "",
			Status:     entity.Status,
			Creator:    entity.Creator,
			CreateTime: entity.CreateTime,
			Updator:    entity.Updator,
			UpdateTime: entity.UpdateTime,
			Version:    1,
			Remark:     *v.Remark,
		}

		entity.IssueDetailBo = *issueDetail

		//负责人
		oid, err3 := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueRelation)
		if err3 != nil {
			return err3
		}
		entity.OwnerInfo = &bo.IssueUserBo{
			IssueRelationBo: bo.IssueRelationBo{
				Id:           oid,
				OrgId:        entity.OrgId,
				ProjectId:    entity.ProjectId,
				IssueId:      entity.Id,
				RelationId:   v.OwnerID,
				RelationType: consts.IssueRelationTypeOwner,
				Creator:      currentUserId,
				CreateTime:   types.NowTime(),
				Updator:      currentUserId,
				UpdateTime:   types.NowTime(),
				Version:      1,
			},
		}

		err4 := domain.CreateIssue(entity, sourceChannel)
		if err4 != nil {
			return err4
		}

		asyn.Execute(func() {
			PushAddIssueNotice(orgId, entity.ProjectId, entity.Id, currentUserId)
		})
	}

	return nil
}

func assemblyIssueTime(entity *bo.IssueBo, v *vo.IssueChildren) {
	if v.PlanEndTime != nil && v.PlanEndTime.IsNotNull() {
		entity.PlanEndTime = *v.PlanEndTime
	} else {
		entity.PlanEndTime = types.Time(consts.BlankTimeObject)
	}

	if v.PlanStartTime != nil && v.PlanStartTime.IsNotNull() {
		entity.PlanStartTime = *v.PlanStartTime
	} else {
		entity.PlanStartTime = types.Time(consts.BlankTimeObject)
	}
	if v.PlanWorkHour != nil {
		entity.PlanWorkHour = *v.PlanWorkHour
	}
}

func UpdateIssue(reqVo projectvo.UpdateIssueReqVo) (*vo.UpdateIssueResp, errs.SystemErrorInfo) {
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	input := reqVo.Input
	sourceChannel := reqVo.SourceChannel

	issueBo, err1 := domain.GetIssueBo(orgId, input.ID)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	err := domain.AuthIssue(orgId, currentUserId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	newIssueBo := &bo.IssueBo{}
	_ = copyer.Copy(issueBo, newIssueBo)

	issueUpdateBo := bo.IssueUpdateBo{
		IssueBo: *issueBo,
	}

	issueUpdateBo.OperatorId = currentUserId

	upd := mysql.Upd{}
	if input.UpdateFields == nil || len(input.UpdateFields) == 0 {
		return nil, errs.BuildSystemErrorInfo(errs.UpdateFiledIsEmpty)
	}

	changeList := []bo.TrendChangeListBo{}
	err1 = UpdateIssueCondAssembly(&upd, newIssueBo, orgId, input, issueBo, &changeList)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}
	//备注
	needUpdateRemark(input, &issueUpdateBo, newIssueBo, issueBo, &changeList)

	//修改负责人
	ownError := needUpdateOwn(input, orgId, issueBo, &upd, &issueUpdateBo, newIssueBo)

	if ownError != nil {
		return nil, ownError
	}
	//参与人
	needUpdateParticipant(input, &issueUpdateBo)
	//关注人
	needUpdateFollower(input, &issueUpdateBo)

	issueUpdateBo.IssueUpdateCond = upd
	issueUpdateBo.NewIssueBo = *newIssueBo
	err5 := domain.UpdateIssue(issueUpdateBo, changeList, sourceChannel)
	if err5 != nil {
		log.Error(err5)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err5)
	}

	asyn.Execute(func() {
		PushModifyIssueNotice(issueBo.OrgId, issueBo.ProjectId, issueBo.Id, currentUserId)
	})

	return &vo.UpdateIssueResp{
		ID: issueBo.Id,
	}, nil
}

//组装更新的issue条件参数
func UpdateIssueCondAssembly(updPoint *mysql.Upd, newIssueBo *bo.IssueBo, orgId int64, input vo.UpdateIssueReq, issueBo *bo.IssueBo, changeList *[]bo.TrendChangeListBo) errs.SystemErrorInfo {

	//判断名字是否需要更新
	needUpdateTitle(input, updPoint, newIssueBo, issueBo, changeList)

	//判断PlanTime是否需要更新
	planTimeError := needUpdatePlanTime(input, updPoint, newIssueBo, issueBo, changeList)

	if planTimeError != nil {
		return planTimeError
	}
	//计划工作小时
	needUpdatePlanWorkHour(input, updPoint, newIssueBo, issueBo, changeList)

	//判断是否更新优先级
	priorityError := needUpdatePriority(input, orgId, updPoint, newIssueBo, issueBo, changeList)
	if priorityError != nil {
		return priorityError
	}

	//处理迭代信息
	iterationError := needUpdateIteration(input, updPoint, newIssueBo)
	if iterationError != nil {
		return iterationError
	}

	//来源
	sourceError := needUpdateSource(input, updPoint, newIssueBo, orgId, issueBo, changeList)
	if sourceError != nil {
		return sourceError
	}

	//类型
	issueObjectTypeError := needUpdateObjectType(input, updPoint, newIssueBo, orgId, issueBo, changeList)
	if issueObjectTypeError != nil {
		return issueObjectTypeError
	}

	return nil
}

func needUpdateObjectType(input vo.UpdateIssueReq, updPoint *mysql.Upd, newIssueBo *bo.IssueBo, orgId int64, issueBo *bo.IssueBo, changeList *[]bo.TrendChangeListBo) errs.SystemErrorInfo {
	if NeedUpdate(input.UpdateFields, "issueObjectTypeId") {
		if input.IssueObjectTypeID == nil || *input.IssueObjectTypeID == 0 {
			return nil
		}
		if domain.IssueObjectTypeExist(orgId, *input.IssueObjectTypeID) {
			(*updPoint)[consts.TcIssueObjectTypeId] = input.IssueObjectTypeID
			newIssueBo.IssueObjectTypeId = *input.IssueObjectTypeID
		} else {
			return errs.BuildSystemErrorInfo(errs.SourceNotExist)
		}

		old, err := domain.GetIssueObjectTypeBo(issueBo.IssueObjectTypeId)
		if err != nil {
			return err
		}
		new, err := domain.GetIssueObjectTypeBo(*input.IssueObjectTypeID)
		if err != nil {
			return err
		}

		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "issueObjectTypeId",
			FieldName: consts.IssueObjectTypeId,
			OldValue:  old.Name,
			NewValue:  new.Name,
		})
	}
	return nil
}

func needUpdateSource(input vo.UpdateIssueReq, updPoint *mysql.Upd, newIssueBo *bo.IssueBo, orgId int64, issueBo *bo.IssueBo, changeList *[]bo.TrendChangeListBo) errs.SystemErrorInfo {
	if NeedUpdate(input.UpdateFields, "sourceId") {
		if input.SourceID == nil || *input.SourceID == 0 {
			return nil
		}
		if domain.SourceExist(orgId, *input.SourceID) {
			(*updPoint)[consts.TcSourceId] = input.SourceID
			newIssueBo.SourceId = *input.SourceID
		} else {
			return errs.BuildSystemErrorInfo(errs.SourceNotExist)
		}
		sourceInfo, err := domain.GetIssueSourceInfo(orgId, []int64{issueBo.SourceId, *input.SourceID})
		if err != nil {
			return err
		}
		var old, new string
		for _, v := range sourceInfo {
			if v.Id == issueBo.SourceId {
				old = v.Name
			} else if v.Id == *input.SourceID {
				new = v.Name
			}
		}
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "sourceId",
			FieldName: consts.Source,
			OldValue:  old,
			NewValue:  new,
		})
	}
	return nil
}

//关注者
func needUpdateFollower(input vo.UpdateIssueReq, issueUpdateBo *bo.IssueUpdateBo) {
	if NeedUpdate(input.UpdateFields, "followerIds") {
		issueUpdateBo.UpdateFollower = true
		//关注者
		if input.FollowerIds != nil {
			issueUpdateBo.Followers = input.FollowerIds
		}
	}
}

//参与人
func needUpdateParticipant(input vo.UpdateIssueReq, issueUpdateBo *bo.IssueUpdateBo) {

	if NeedUpdate(input.UpdateFields, "participantIds") {
		issueUpdateBo.UpdateParticipant = true
		//参与人
		if input.ParticipantIds != nil {
			issueUpdateBo.Participants = input.ParticipantIds
		}
	}
}

//负责人
func needUpdateOwn(input vo.UpdateIssueReq, orgId int64, issueBo *bo.IssueBo, upd *mysql.Upd, issueUpdateBo *bo.IssueUpdateBo, newIssueBo *bo.IssueBo) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "ownerId") {
		if input.OwnerID == nil {
			return errs.BuildSystemErrorInfo(errs.IssueOwnerCantBeNull)
		}
		if !orgfacade.VerifyOrgRelaxed(orgId, *input.OwnerID) {
			return errs.BuildSystemErrorInfoWithMessage(errs.IllegalityOwner, "负责人不存在")
		}
		//如果所有者变动， 则更新
		if issueBo.Owner != *input.OwnerID {
			ownerChangeTime := types.NowTime()
			(*upd)[consts.TcOwner] = *input.OwnerID
			(*upd)[consts.TcOwnerChangeTime] = date.FormatTime(ownerChangeTime)
			issueUpdateBo.OwnerId = input.OwnerID
			newIssueBo.Owner = *input.OwnerID
			newIssueBo.OwnerChangeTime = ownerChangeTime
		}
	}

	return nil
}

//备注
func needUpdateRemark(input vo.UpdateIssueReq, issueUpdateBo *bo.IssueUpdateBo, newIssueBo *bo.IssueBo, issueBo *bo.IssueBo, changeList *[]bo.TrendChangeListBo) {

	if NeedUpdate(input.UpdateFields, "remark") {
		remark := ""
		if input.Remark != nil {
			remark = *input.Remark
		}
		(*issueUpdateBo).IssueDetailRemark = &remark
		newIssueBo.Remark = remark
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "remark",
			FieldName: consts.Remark,
			OldValue:  issueBo.Remark,
			NewValue:  remark,
		})
	}
}

//优先级
func needUpdatePriority(input vo.UpdateIssueReq, orgId int64, updPoint *mysql.Upd, newIssueBo *bo.IssueBo, issueBo *bo.IssueBo, changeList *[]bo.TrendChangeListBo) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "priorityId") {
		if input.PriorityID != nil {
			suc, err := VerifyPriority(orgId, consts.PriorityTypeIssue, *input.PriorityID)
			if err != nil {
				return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
			}
			if !suc {
				return errs.BuildSystemErrorInfo(errs.IllegalityPriority, err)
			}
			(*updPoint)[consts.TcPriorityId] = input.PriorityID
			newIssueBo.PriorityId = *input.PriorityID
			old, err := domain.GetPriorityBo(issueBo.PriorityId)
			if err != nil {
				return err
			}
			new, err := domain.GetPriorityBo(*input.PriorityID)
			if err != nil {
				return err
			}
			*changeList = append(*changeList, bo.TrendChangeListBo{
				Field:     "priorityId",
				FieldName: consts.Priority,
				OldValue:  old.Name,
				NewValue:  new.Name,
			})
		}
	}

	return nil

}

//迭代信息
func needUpdateIteration(input vo.UpdateIssueReq, updPoint *mysql.Upd, newIssueBo *bo.IssueBo) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "iterationId") {
		if input.IterationID != nil && *input.IterationID != 0 {
			_, err := domain.GetIterationBo(*input.IterationID)
			if err != nil {
				return errs.BuildSystemErrorInfo(errs.IterationNotExist)
			}
		}
		(*updPoint)[consts.TcIterationId] = input.IterationID
		newIssueBo.IterationId = *input.IterationID
	}
	return nil
}

//计划工作时间
func needUpdatePlanWorkHour(input vo.UpdateIssueReq, updPoint *mysql.Upd, newIssueBo *bo.IssueBo, issueBo *bo.IssueBo, changeList *[]bo.TrendChangeListBo) {
	upd := *updPoint

	if NeedUpdate(input.UpdateFields, "planWorkHour") {
		updatePlanWorkHour := 0
		if input.PlanWorkHour != nil {
			updatePlanWorkHour = *input.PlanWorkHour
		}
		upd[consts.TcPlanWorkHour] = updatePlanWorkHour
		newIssueBo.PlanWorkHour = updatePlanWorkHour
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "planWorkHour",
			FieldName: consts.PlanWorkHour,
			OldValue:  strconv.Itoa(issueBo.PlanWorkHour),
			NewValue:  strconv.Itoa(updatePlanWorkHour),
		})
	}
}

//更新时间校验
func needUpdatePlanTime(input vo.UpdateIssueReq, updPoint *mysql.Upd, newIssueBo *bo.IssueBo, issueBo *bo.IssueBo, changeList *[]bo.TrendChangeListBo) errs.SystemErrorInfo {
	upd := *updPoint

	if NeedUpdate(input.UpdateFields, "planEndTime") {
		var updatePlanEndTime = consts.BlankTimeObject
		if input.PlanEndTime != nil && input.PlanEndTime.IsNotNull() {
			updatePlanEndTime = time.Time(*input.PlanEndTime)
		}
		upd[consts.TcPlanEndTime] = date.Format(updatePlanEndTime)
		newIssueBo.PlanEndTime = types.Time(updatePlanEndTime)
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "planEndTime",
			FieldName: consts.PlanEndTime,
			OldValue:  issueBo.PlanEndTime.String(),
			NewValue:  date.Format(updatePlanEndTime),
		})
	}
	if NeedUpdate(input.UpdateFields, "planStartTime") {
		var updatePlanStartTime = consts.BlankTimeObject
		if input.PlanStartTime != nil && input.PlanStartTime.IsNotNull() {
			updatePlanStartTime = time.Time(*input.PlanStartTime)
		}
		upd[consts.TcPlanStartTime] = date.Format(updatePlanStartTime)
		newIssueBo.PlanStartTime = types.Time(updatePlanStartTime)
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "planStartTime",
			FieldName: consts.PlanStartTime,
			OldValue:  issueBo.PlanStartTime.String(),
			NewValue:  date.Format(updatePlanStartTime),
		})
	}
	//更新时间校验
	planStartTime := newIssueBo.PlanStartTime
	planEndTime := newIssueBo.PlanEndTime

	fmt.Println(planStartTime, planEndTime)
	if planStartTime.IsNotNull() && planEndTime.IsNotNull() {
		if time.Time(planEndTime).Before(time.Time(planStartTime)) {
			return errs.BuildSystemErrorInfo(errs.PlanEndTimeInvalidError)
		}
	}

	return nil
}

//需要更新标题
func needUpdateTitle(input vo.UpdateIssueReq, updPoint *mysql.Upd, newIssueBo *bo.IssueBo, issueBo *bo.IssueBo, changeList *[]bo.TrendChangeListBo) {
	upd := *updPoint

	if NeedUpdate(input.UpdateFields, "title") {
		updateTitle := consts.BlankString
		if input.Title != nil {
			updateTitle = *input.Title
		}
		upd[consts.TcTitle] = updateTitle
		newIssueBo.Title = updateTitle
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "title",
			FieldName: consts.Title,
			OldValue:  issueBo.Title,
			NewValue:  updateTitle,
		})
	}
}

func DeleteIssue(reqVo projectvo.DeleteIssueReqVo) (*vo.Issue, errs.SystemErrorInfo) {
	input := reqVo.Input
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	sourceChannel := reqVo.SourceChannel

	issueBo, err2 := domain.GetIssueBo(orgId, input.ID)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
	}

	err := domain.AuthIssue(orgId, currentUserId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationDelete)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	err3 := domain.DeleteIssue(issueBo, currentUserId, sourceChannel)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}

	asyn.Execute(func() {
		PushDelIssueNotice(issueBo.OrgId, issueBo.ProjectId, []int64{issueBo.Id})
	})

	result := &vo.Issue{}
	copyErr := copyer.Copy(issueBo, result)
	if copyErr != nil {
		log.Errorf("copyer.Copy(): %q\n", copyErr)
	}
	return result, nil
}

func GetIssueRestInfos(orgId int64, page int, size int, input *vo.IssueRestInfoReq) (*vo.IssueRestInfoResp, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	log.Errorf(" GetCurrentUser: %q\n", err)
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//orgId := cacheUserInfo.OrgId

	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if input != nil {
		dealDbCond(input, &cond)
		if input.Status != nil {
			err1 := domain.IssueCondStatusAssembly(cond, orgId, *input.Status)
			if err1 != nil {
				log.Error(err1)
				return nil, errs.BuildSystemErrorInfo(errs.IssueCondAssemblyError, err1)
			}
		}
	}
	issueList, total, err := domain.SelectList(cond, nil, page, size, "")
	if err != nil {
		log.Errorf(" issuedomain.SelectList: %q\n", err)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	issueRestInfos := make([]*vo.IssueRestInfo, len(*issueList))

	finishedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
	if err != nil {
		log.Errorf("proxies.GetProcessStatusId: %q\n", err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	ownerIds := make([]int64, 0)
	for _, issueBo := range *issueList {
		ownerId := issueBo.Owner
		ownerIdExist, _ := slice.Contain(ownerIds, ownerId)
		if !ownerIdExist {
			ownerIds = append(ownerIds, ownerId)
		}
	}
	ownerInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(consts.AppSourceChannelDingTalk, orgId, ownerIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	statusList, err := processfacade.GetProcessStatusListByCategoryRelaxed(orgId, consts.ProcessStatusCategoryIssue)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	statusMap := maps.NewMap("StatusId", statusList)
	ownerMap := maps.NewMap("UserId", ownerInfos)

	priorityList, err := domain.GetPriorityList(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	priorityMap := maps.NewMap("Id", priorityList)
	for i, issueInfo := range *issueList {
		finished, err := slice.Contain(*finishedIds, issueInfo.Status)
		if err != nil {
			log.Error(err)
			continue
		}
		issueRestInfo := &vo.IssueRestInfo{
			ID:          issueInfo.Id,
			Title:       issueInfo.Title,
			OwnerID:     issueInfo.Owner,
			Finished:    finished,
			PlanEndTime: issueInfo.PlanEndTime,
			PlanStartTime:issueInfo.PlanStartTime,
			EndTime:     issueInfo.EndTime,
		}

		if userCacheInfo, ok := ownerMap[issueInfo.Owner]; ok {
			baseUserInfo := userCacheInfo.(bo.BaseUserInfoBo)
			issueRestInfo.OwnerName = baseUserInfo.Name
			issueRestInfo.OwnerAvatar = baseUserInfo.Avatar
			issueRestInfo.OwnerIsDeleted = baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted
			issueRestInfo.OwnerIsDisabled = baseUserInfo.OrgUserStatus == consts.AppStatusDisabled
		}

		if statusCacheInfoInterface, ok := statusMap[issueInfo.Status]; ok {
			if statusCacheInfo, ok := statusCacheInfoInterface.(bo.CacheProcessStatusBo); ok {
				issueRestInfo.StatusID = statusCacheInfo.StatusId
				issueRestInfo.StatusName = statusCacheInfo.Name
			}
		}

		//优先级
		if priorityCacheInfo, ok := priorityMap[issueInfo.PriorityId]; ok {
			priorityInfo := priorityCacheInfo.(bo.PriorityBo)
			issueRestInfo.PriorityInfo = &vo.HomeIssuePriorityInfo{
				ID:        priorityInfo.Id,
				Name:      priorityInfo.Name,
				BgStyle:   priorityInfo.BgStyle,
				FontStyle: priorityInfo.FontStyle,
			}
		}

		issueRestInfos[i] = issueRestInfo
	}
	return &vo.IssueRestInfoResp{
		Total: int64(total),
		List:  issueRestInfos,
	}, nil
}

//处理db的cond
func dealDbCond(input *vo.IssueRestInfoReq, cond *db.Cond) {
	if input.ProjectID != nil {
		(*cond)[consts.TcProjectId] = *input.ProjectID
	}
	if input.ParentID != nil {
		(*cond)[consts.TcParentId] = *input.ParentID
	}
	if input.IssueIds != nil && len(input.IssueIds) > 0 {
		(*cond)[consts.TcId] = db.In(input.IssueIds)
	}
}

func UpdateIssueProjectObjectType(orgId, currentUserId int64, input vo.UpdateIssueProjectObjectTypeReq) (*vo.Void, errs.SystemErrorInfo) {
	//获取原本的Issue
	issueBo, err := domain.GetIssueBo(orgId, input.ID)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IssueNotExist, err)
	}

	err = domain.AuthIssue(orgId, currentUserId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	//更新任务的项目对象类型
	err1 := domain.UpdateIssueProjectObjectType(orgId, currentUserId, *issueBo, input.ProjectObjectTypeID)

	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	asyn.Execute(func() {
		PushModifyIssueNotice(issueBo.OrgId, issueBo.ProjectId, issueBo.Id, currentUserId)
	})

	return &vo.Void{
		ID: input.ID,
	}, nil
}

func GetIssueInfoList(ids []int64) ([]vo.Issue, errs.SystemErrorInfo) {
	result, err := domain.GetIssueInfoList(ids)
	if err != nil {
		return nil, err
	}
	issueResp := &[]vo.Issue{}
	err1 := copyer.Copy(result, issueResp)
	if err1 != nil {
		log.Errorf("copyer.Copy: %q\n", err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}

	return *issueResp, nil
}
