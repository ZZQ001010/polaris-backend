package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"sync"
	"upper.io/db.v3"
)

func IssueInfo(orgId, currentUserId, issueID int64, sourceChannel string) (*vo.IssueInfo, errs.SystemErrorInfo) {
	issueBo, err1 := domain.GetIssueBo(orgId, issueID)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	//err := domain.AuthIssue(orgId, currentUserId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationView)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	//}

	issueInfoResp := &vo.IssueInfo{}
	handlerFuncList := make([]func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup), 0)

	//issue信息和详情
	handlerFuncList = append(handlerFuncList, func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup) {
		defer wg.Add(-1)
		issueInfo, errorInfo := dealIssueInfo(orgId, issueID, issueBo)
		if errorInfo != nil {
			log.Error(errorInfo)
			return
		}
		issueInfoResp.Issue = issueInfo
	})

	//任务详情
	handlerFuncList = append(handlerFuncList, func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup) {
		defer wg.Add(-1)
		projectObjectTypeBo, err := domain.GetProjectObjectTypeById(orgId, issueBo.ProjectId, issueBo.ProjectObjectTypeId)
		if err != nil {
			log.Error(err)
			return
		}
		issueInfoResp.ProjectObjectTypeName = projectObjectTypeBo.Name
	})

	//状态信息
	handlerFuncList = append(handlerFuncList, func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup) {
		defer wg.Add(-1)
		statusInfo, errStatus := dealStatus(orgId, issueBo)
		if errStatus != nil {
			log.Error(errStatus)
			return
		}
		issueInfoResp.Status = statusInfo
		//流程 和流程步骤-接下来的状态
		nextStatusInfoList, errorProcess := dealProcessAndNextStatusList(orgId, issueBo, statusInfo)
		if errorProcess != nil {
			log.Error(errorProcess)
			return
		}
		issueInfoResp.NextStatus = nextStatusInfoList
	})

	//优先级
	handlerFuncList = append(handlerFuncList, func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup) {
		defer wg.Add(-1)
		priorityInfo, errPriority := dealPriorityInfo(orgId, issueBo)
		if errPriority != nil {
			log.Error(errPriority)
			return
		}
		issueInfoResp.Priority = priorityInfo
	})

	//项目信息
	handlerFuncList = append(handlerFuncList, func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup) {
		defer wg.Add(-1)
		projectInfoBo, err2 := domain.GetHomeProjectInfoBo(orgId, issueBo.ProjectId)
		if err2 != nil {
			log.Error(err2)
			return
		}
		projectInfo := &vo.HomeIssueProjectInfo{}
		copyErr := copyer.Copy(projectInfoBo, projectInfo)
		if copyErr != nil {
			log.Error(copyErr)
			return
		}
		issueInfoResp.Project = projectInfo
	})

	//负责人信息
	handlerFuncList = append(handlerFuncList, func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup) {
		defer wg.Add(-1)
		ownerId := issueBo.Owner
		ownerBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(sourceChannel, orgId, ownerId)
		if err != nil {
			log.Error(err)
			return
		}
		ownerInfo := AssemblyUserIdInfo(ownerBaseInfo)
		issueInfoResp.Owner = ownerInfo
	})

	//创建人信息
	handlerFuncList = append(handlerFuncList, func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup) {
		defer wg.Add(-1)
		creatorId := issueBo.Creator
		creatorBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(sourceChannel, orgId, creatorId)
		if err != nil {
			log.Error(err)
			return
		}
		creatorInfo := AssemblyUserIdInfo(creatorBaseInfo)
		issueInfoResp.CreatorInfo = creatorInfo
	})

	//参与人信息
	handlerFuncList = append(handlerFuncList, func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup) {
		defer wg.Add(-1)
		participantInfos, errParticipant := dealParticipant(orgId, issueBo, sourceChannel)
		if errParticipant != nil {
			log.Error(errParticipant)
			return
		}
		issueInfoResp.ParticipantInfos = participantInfos
	})

	//关注人信息
	handlerFuncList = append(handlerFuncList, func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup) {
		defer wg.Add(-1)
		followerInfos, errFollower := dealFollowerInfos(orgId, issueBo, sourceChannel)
		if errFollower != nil {
			log.Error(errFollower)
			return
		}
		issueInfoResp.FollowerInfos = followerInfos
	})

	//任务来源和类型
	handlerFuncList = append(handlerFuncList, func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup) {
		defer wg.Add(-1)
		//任务来源信息
		issueInfoResp.SourceInfo = dealIssueSourceInfo(orgId, issueBo)
		//任务类型
		issueInfoResp.TypeInfo = dealIssueObjectTypeInfo(orgId, issueBo)
	})

	//任务标签
	handlerFuncList = append(handlerFuncList, func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup) {
		defer wg.Add(-1)
		homeIssueTagInfos, errIssueTags := dealIssueTags(orgId, issueBo)
		if errIssueTags != nil {
			log.Error(errIssueTags)
			return
		}
		//任务标签
		issueInfoResp.Tags = homeIssueTagInfos
	})

	//子任务数量以及已完成的子任务数量
	handlerFuncList = append(handlerFuncList, func(issueInfoResp *vo.IssueInfo, wg *sync.WaitGroup) {
		defer wg.Add(-1)
		childNum, childFinishedNum, err := dealIssueChildNum(issueBo.Id, orgId)
		if err != nil {
			log.Error(err)
			return
		}
		issueInfoResp.ChildsNum = childNum
		issueInfoResp.ChildsFinishedNum = childFinishedNum
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
			currentFunc(issueInfoResp, &wg)
		}()
	}

	wg.Wait()
	return issueInfoResp, nil
}

func dealIssueChildNum(issueId, orgId int64) (int64, int64, errs.SystemErrorInfo) {
	count, err1 := mysql.SelectCountByCond(consts.TableIssue, db.Cond{
		consts.TcParentId: issueId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err1 != nil {
		log.Error(err1)
		return 0, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	finishedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
	if err != nil {
		log.Error(err)
		return 0, 0, err
	}

	finishedCount, err1 := mysql.SelectCountByCond(consts.TableIssue, db.Cond{
		consts.TcParentId: issueId,
		consts.TcStatus:   db.In(finishedIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err1 != nil {
		log.Error(err1)
		return 0, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return int64(count), int64(finishedCount), nil
}

func dealIssueInfo(orgId int64, issueID int64, issueBo *bo.IssueBo) (*vo.Issue, errs.SystemErrorInfo) {

	//Issue信息
	issueInfo := &vo.Issue{}
	err := copyer.Copy(issueBo, issueInfo)
	if issueBo.ParentId != 0 {
		issueParentBo, err := domain.GetIssueBo(orgId, issueBo.ParentId)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.ParentIssueNotExist)
		}
		//指针变量.属性获取到对应的那个变量 赋值
		issueInfo.ParentTitle = issueParentBo.Title
	}

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}

	//获取issue详情
	issueDetail, err1 := domain.GetIssueDetailBo(orgId, issueID)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}
	issueInfo.Remark = &issueDetail.Remark

	return issueInfo, nil

}

func dealStatus(orgId int64, issueBo *bo.IssueBo) (*vo.HomeIssueStatusInfo, errs.SystemErrorInfo) {

	statusInfoBo, err1 := domain.GetHomeIssueStatusInfoBo(orgId, issueBo.Status)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}
	statusInfo := &vo.HomeIssueStatusInfo{}
	err := copyer.Copy(statusInfoBo, statusInfo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return statusInfo, nil
}

func dealProcessAndNextStatusList(orgId int64, issueBo *bo.IssueBo, statusInfo *vo.HomeIssueStatusInfo) ([]*vo.HomeIssueStatusInfo, errs.SystemErrorInfo) {

	process, err := domain.GetProjectProcessBo(orgId, issueBo.ProjectId, issueBo.ProjectObjectTypeId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	//流程步骤-接下来的状态
	nextStatusList, err := processfacade.GetNextProcessStepStatusListRelaxed(orgId, process.Id, statusInfo.ID)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	nextStatusInfoList := make([]*vo.HomeIssueStatusInfo, len(*nextStatusList))
	for i, nextStatus := range *nextStatusList {
		nextStatusName := nextStatus.DisplayName
		nextStatusInfoList[i] = &vo.HomeIssueStatusInfo{
			ID:          nextStatus.StatusId,
			Name:        nextStatus.Name,
			BgStyle:     nextStatus.BgStyle,
			FontStyle:   nextStatus.FontStyle,
			Type:        nextStatus.StatusType,
			DisplayName: &nextStatusName,
		}
	}

	return nextStatusInfoList, nil
}

func dealPriorityInfo(orgId int64, issueBo *bo.IssueBo) (*vo.HomeIssuePriorityInfo, errs.SystemErrorInfo) {

	priorityInfoBo, err2 := domain.GetHomeIssuePriorityInfoBo(orgId, issueBo.PriorityId)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
	}
	priorityInfo := &vo.HomeIssuePriorityInfo{}
	err := copyer.Copy(priorityInfoBo, priorityInfo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return priorityInfo, nil
}

func dealParticipant(orgId int64, issueBo *bo.IssueBo, sourceChannel string) ([]*vo.UserIDInfo, errs.SystemErrorInfo) {
	//参与人信息
	participantIds, err7 := domain.GetIssueRelationIdsByRelateType(orgId, issueBo.Id, consts.IssueRelationTypeParticipant)
	if err7 != nil {
		log.Error(err7)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err7)
	}

	participantInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(sourceChannel, orgId, *participantIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	participantIdInfos := make([]*vo.UserIDInfo, len(*participantIds))
	for i, participantInfo := range participantInfos {
		participantIdInfos[i] = AssemblyUserIdInfo(&participantInfo)
	}
	return participantIdInfos, nil
}

func dealFollowerInfos(orgId int64, issueBo *bo.IssueBo, sourceChannel string) ([]*vo.UserIDInfo, errs.SystemErrorInfo) {
	followerIds, err7 := domain.GetIssueRelationIdsByRelateType(orgId, issueBo.Id, consts.IssueRelationTypeFollower)
	if err7 != nil {
		log.Error(err7)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err7)
	}
	followerInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(sourceChannel, orgId, *followerIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	followerIdInfos := make([]*vo.UserIDInfo, len(*followerIds))
	for i, followerInfo := range followerInfos {
		followerIdInfos[i] = AssemblyUserIdInfo(&followerInfo)
	}
	return followerIdInfos, nil

}

func dealIssueTags(orgId int64, issueBo *bo.IssueBo) ([]*vo.HomeIssueTagInfo, errs.SystemErrorInfo) {
	issueTagBos, err := domain.GetIssueTags(orgId, issueBo.Id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	homeIssueTagInfo := make([]*vo.HomeIssueTagInfo, len(issueTagBos))
	for i, issueTagBo := range issueTagBos {
		homeIssueTagInfo[i] = &vo.HomeIssueTagInfo{
			ID:        issueTagBo.TagId,
			Name:      issueTagBo.TagName,
			FontStyle: issueTagBo.FontStyle,
			BgStyle:   issueTagBo.BgStyle,
		}
	}
	return homeIssueTagInfo, nil
}

func dealIssueSourceInfo(orgId int64, issueBo *bo.IssueBo) *vo.IssueSourceInfo {
	issueSource, err1 := domain.GetIssueSourceById(orgId, issueBo.SourceId)
	if err1 != nil {
		log.Error(err1)
		return nil
	}
	issueSourceInfo := &vo.IssueSourceInfo{
		ID:   issueSource.Id,
		Name: issueSource.Name,
	}
	return issueSourceInfo
}

func dealIssueObjectTypeInfo(orgId int64, issueBo *bo.IssueBo) *vo.IssueObjectTypeInfo {
	issueType, err1 := domain.GetIssueObjectTypeById(orgId, issueBo.IssueObjectTypeId)
	if err1 != nil {
		log.Error(err1)
		return nil
	}
	issueTypeInfo := &vo.IssueObjectTypeInfo{
		ID:   issueType.Id,
		Name: issueType.Name,
	}
	return issueTypeInfo
}
