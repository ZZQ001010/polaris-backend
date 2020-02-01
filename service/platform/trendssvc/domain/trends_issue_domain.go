package domain

import (
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/processvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"strconv"
)

func AddIssueTrends(issueTrendsBo bo.IssueTrendsBo) {
	var trendsBos []bo.TrendsBo = nil

	assemblyError := assemblyTrendsBos(issueTrendsBo, &trendsBos)

	if assemblyError != nil || trendsBos == nil {
		return
	}

	for i, _ := range trendsBos {
		trendsBos[i].CreateTime = types.Time(issueTrendsBo.OperateTime)
	}

	err := CreateTrendsBatch(trendsBos)
	if err != nil{
		log.Error(err)
	}
}

func assemblyTrendsBos(issueTrendsBo bo.IssueTrendsBo, trendsBos *[]bo.TrendsBo) errs.SystemErrorInfo {
	trendsType := issueTrendsBo.PushType
	var err1 errs.SystemErrorInfo = nil

	if trendsType == consts.PushTypeCreateIssue {
		*trendsBos, err1 = assemblyCreateIssueTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeUpdateIssue {
		*trendsBos, err1 = assemblyUpdateIssueTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeUpdateIssueMembers {
		*trendsBos, err1 = assemblyUpdateIssueMemberTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeDeleteIssue {
		*trendsBos, err1 = assemblyDeleteIssueTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeUpdateIssueStatus {
		*trendsBos, err1 = assemblyUpdateIssueStatusTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeUpdateRelationIssue {
		*trendsBos, err1 = assemblyUpdateRelationIssueTrends(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeIssueComment {
		*trendsBos, err1 = assemblyCreateIssueComment(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeUploadResource {
		*trendsBos, err1 = assemblyUploadResource(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeUpdateIssueProjectObjectType {
		*trendsBos, err1 = assemblyUpdateIssueProjectObjectType(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	} else if trendsType == consts.PushTypeDeleteResource {
		*trendsBos, err1 = assemblyDeleteResource(issueTrendsBo)
		if err1 != nil {
			log.Errorf(consts.Assembly_CreateIssueTrends_eorror_printf, err1)
			return err1
		}
	}

	return nil
}

//创建评论
func assemblyCreateIssueComment(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := make([]bo.TrendsBo, 0)
	ext := issueTrendsBo.Ext
	ext.IssueType = "T"
	ext.ObjName = issueTrendsBo.IssueTitle

	operatorId := issueTrendsBo.OperatorId

	commentBo := ext.CommentBo
	commentId := commentBo.Id

	//拼装动态(这里要求实时)
	trendsBo := bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationCreate,
		OperObjId:       commentId,
		OperObjType:     consts.TrendsOperObjectTypeComment,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeCreateIssueComment,
		RelationObjType: consts.TrendsOperObjectTypeIssue,
		NewValue:        &commentBo.Content,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(ext),
		Creator:         operatorId,
	}
	trendsBos = append(trendsBos, trendsBo)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("捕获到的错误：%s", r)
			}
		}()
		noticeErr := AddNoticeByCreateComment(issueTrendsBo)
		if noticeErr != nil {
			log.Error(noticeErr)
		}
	}()

	return trendsBos, nil
}

//创建任务
func assemblyCreateIssueTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := []bo.TrendsBo{}
	ext := bo.TrendExtensionBo{}
	ext.IssueType = "T"
	ext.ObjName = issueTrendsBo.IssueTitle

	newValue := issueTrendsBo.NewValue
	trendsBo := bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationCreate,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.ParentIssueId,
		RelationType:    consts.TrendsRelationTypeCreateIssue,
		RelationObjType: consts.TrendsOperObjectTypeIssue,
		NewValue:        &newValue,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(ext),
		Creator:         issueTrendsBo.OperatorId,
	}

	trendsBos = append(trendsBos, trendsBo)

	//if issueTrendsBo.ParentIssueId > 0 {
	//	parentIssueInfoResp := projectfacade.IssueInfo(projectvo.IssueInfoReqVo{
	//		OrgId:   issueTrendsBo.OrgId,
	//		UserId:  issueTrendsBo.OperatorId,
	//		IssueID: issueTrendsBo.ParentIssueId,
	//	})
	//	if parentIssueInfoResp.Failure() {
	//		log.Error(parentIssueInfoResp.Error())
	//		return nil, parentIssueInfoResp.Error()
	//	}
	//	ext.ObjName = parentIssueInfoResp.IssueInfo.Issue.Title
	//	ext.RelationIssue = bo.RelationIssue{ID: issueTrendsBo.IssueId, Title: issueTrendsBo.IssueTitle}
	//	createChildIssueTrendsBo := bo.TrendsBo{
	//		OrgId:           issueTrendsBo.OrgId,
	//		Module1:         consts.TrendsModuleOrg,
	//		Module2Id:       issueTrendsBo.ProjectId,
	//		Module2:         consts.TrendsModuleProject,
	//		Module3Id:       issueTrendsBo.ParentIssueId,
	//		Module3:         consts.TrendsModuleIssue,
	//		OperCode:        consts.RoleOperationCreate,
	//		OperObjId:       issueTrendsBo.IssueId,
	//		OperObjType:     consts.TrendsOperObjectTypeIssue,
	//		OperObjProperty: consts.BlankString,
	//		RelationObjId:   issueTrendsBo.IssueId,
	//		RelationObjType: consts.TrendsOperObjectTypeIssue,
	//		RelationType:    consts.TrendsRelationTypeCreateChildIssue,
	//		NewValue:        &newValue,
	//		OldValue:        nil,
	//		Ext:             json.ToJsonIgnoreError(ext),
	//		Creator:         issueTrendsBo.OperatorId,
	//	}
	//	trendsBos = append(trendsBos, createChildIssueTrendsBo)
	//}

	return trendsBos, nil
}

//更新任务
func assemblyUpdateIssueTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	newValue := issueTrendsBo.NewValue
	oldValue := issueTrendsBo.OldValue
	operObjProperty := issueTrendsBo.OperateObjProperty
	trendsBo := &bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: operObjProperty,
		RelationObjId:   0,
		RelationType:    consts.TrendsRelationTypeUpdateIssue,
		RelationObjType: consts.BlankString,
		NewValue:        &newValue,
		OldValue:        &oldValue,
		Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator:         issueTrendsBo.OperatorId,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("捕获到的错误：%s", r)
			}
		}()
		noticeErr := AddNoticeByChangeUpdateIssue(issueTrendsBo)
		if noticeErr != nil {
			log.Error(noticeErr)
		}
	}()


	return []bo.TrendsBo{*trendsBo}, nil
}

//更新成员
func assemblyUpdateIssueMemberTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	deletedFollowerIds, addedFollowerIds := util.GetDifMemberIds(issueTrendsBo.BeforeChangeFollowers, issueTrendsBo.AfterChangeFollowers)
	deletedParticipantIds, addedParticipantIds := util.GetDifMemberIds(issueTrendsBo.BeforeChangeParticipants, issueTrendsBo.AfterChangeParticipants)

	memberChangeInfos := [][]int64{deletedFollowerIds, addedFollowerIds, deletedParticipantIds, addedParticipantIds}
	trendsBos := &[]bo.TrendsBo{}
	//获取所有相关用户
	allRelationIds := slice.SliceUniqueInt64(append(append(append(append(append(deletedFollowerIds, addedFollowerIds...), deletedParticipantIds...), addedParticipantIds...), issueTrendsBo.BeforeOwner, issueTrendsBo.AfterOwner)))
	userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(consts.AppSourceChannelDingTalk, issueTrendsBo.OrgId, allRelationIds)
	if err != nil {
		log.Error("获取用户相关信息失败" + strs.ObjectToString(err))
		return nil, err
	}
	userInfo := map[int64]bo.SimpleUserInfoBo{}
	for _, v := range userInfos {
		userInfo[v.UserId] = bo.SimpleUserInfoBo{
			Id:     v.UserId,
			Name:   v.Name,
			Avatar: v.Avatar,
		}
	}
	for i, changeInfos := range memberChangeInfos {
		if len(changeInfos) > 0 {
			ext := bo.TrendExtensionBo{}
			ext.IssueType = "T"
			ext.ObjName = issueTrendsBo.IssueTitle
			for _, v := range changeInfos {
				ext.MemberInfo = append(ext.MemberInfo, userInfo[v])
			}
			beforeMap := map[string]interface{}{}
			afterMap := map[string]interface{}{}
			operCode := ""
			relationType := ""
			operObjProperty := ""

			assemblyPartInfo(i, issueTrendsBo, &operObjProperty, &operCode, &relationType, &beforeMap, &afterMap)

			newValue := json.ToJsonIgnoreError(afterMap)
			oldValue := json.ToJsonIgnoreError(beforeMap)
			trendsBo := bo.TrendsBo{
				OrgId:           issueTrendsBo.OrgId,
				Module1:         consts.TrendsModuleOrg,
				Module2Id:       issueTrendsBo.ProjectId,
				Module2:         consts.TrendsModuleProject,
				Module3Id:       issueTrendsBo.IssueId,
				Module3:         consts.TrendsModuleIssue,
				OperCode:        operCode,
				OperObjId:       issueTrendsBo.IssueId,
				OperObjType:     consts.TrendsOperObjectTypeIssue,
				OperObjProperty: operObjProperty,
				RelationObjId:   0,
				RelationType:    relationType,
				RelationObjType: consts.BlankString,
				NewValue:        &newValue,
				OldValue:        &oldValue,
				Ext:             json.ToJsonIgnoreError(ext),
				Creator:         issueTrendsBo.OperatorId,
			}
			*trendsBos = append(*trendsBos, trendsBo)
		}
	}

	if issueTrendsBo.BeforeOwner != issueTrendsBo.AfterOwner {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("捕获到的错误：%s", r)
				}
			}()
			noticeErr := AddNoticeByChangeIssueOwner(issueTrendsBo)
			if noticeErr != nil {
				log.Error(noticeErr)
			}
		}()

		//负责人动态
		ext := bo.TrendExtensionBo{}
		ext.IssueType = "T"
		ext.ObjName = issueTrendsBo.IssueTitle

		//var newValue, oldValue string
		//if _, ok := userInfo[issueTrendsBo.AfterOwner];ok {
		//	newValue = json.ToJsonIgnoreError(userInfo[issueTrendsBo.AfterOwner])
		//}
		//if _, ok := userInfo[issueTrendsBo.BeforeOwner];ok {
		//	oldValue = json.ToJsonIgnoreError(userInfo[issueTrendsBo.BeforeOwner])
		//}

		oldValue := strconv.FormatInt(issueTrendsBo.BeforeOwner, 10)
		newValue := strconv.FormatInt(issueTrendsBo.AfterOwner, 10)
		if _, ok := userInfo[issueTrendsBo.BeforeOwner]; ok {
			ext.MemberInfo = append(ext.MemberInfo, userInfo[issueTrendsBo.BeforeOwner])
		}
		if _, ok := userInfo[issueTrendsBo.AfterOwner]; ok {
			ext.MemberInfo = append(ext.MemberInfo, userInfo[issueTrendsBo.AfterOwner])
		}

		*trendsBos = append(*trendsBos, bo.TrendsBo{
			OrgId:           issueTrendsBo.OrgId,
			Module1:         consts.TrendsModuleOrg,
			Module2Id:       issueTrendsBo.ProjectId,
			Module2:         consts.TrendsModuleProject,
			Module3Id:       issueTrendsBo.IssueId,
			Module3:         consts.TrendsModuleIssue,
			OperCode:        consts.RoleOperationModify,
			OperObjId:       issueTrendsBo.IssueId,
			OperObjType:     consts.TrendsOperObjectTypeIssue,
			OperObjProperty: consts.TrendsOperObjPropertyNameOwner,
			RelationObjId:   0,
			RelationType:    consts.TrendsRelationTypeUpdateIssueOwner,
			RelationObjType: consts.BlankString,
			NewValue:        &newValue,
			OldValue:        &oldValue,
			Ext:             json.ToJsonIgnoreError(ext),
			Creator:         issueTrendsBo.OperatorId,
		})
	}
	return *trendsBos, nil
}

func assemblyPartInfo(i int, issueTrendsBo bo.IssueTrendsBo, operObjProperty, operCode, relationType *string, beforeMap, afterMap *map[string]interface{}) {
	if i == 0 || i == 1 {
		*operObjProperty = consts.TrendsOperObjPropertyNameFollower
		(*beforeMap)[*operObjProperty] = issueTrendsBo.BeforeChangeFollowers
		(*afterMap)[*operObjProperty] = issueTrendsBo.AfterChangeFollowers
		if i == 0 {
			*operCode = consts.RoleOperationUnbind
			*relationType = consts.TrendsRelationTypeDeleteIssueFollower
		} else if i == 1 {
			*operCode = consts.RoleOperationBind
			*relationType = consts.TrendsRelationTypeAddedIssueFollower
		}
	} else if i == 2 || i == 3 {
		*operObjProperty = consts.TrendsOperObjPropertyNameParticipant
		(*beforeMap)[*operObjProperty] = issueTrendsBo.BeforeChangeParticipants
		(*afterMap)[*operObjProperty] = issueTrendsBo.AfterChangeParticipants
		if i == 2 {
			*operCode = consts.RoleOperationUnbind
			*relationType = consts.TrendsRelationTypeDeletedIssueParticipant
		} else if i == 3 {
			*operCode = consts.RoleOperationBind
			*relationType = consts.TrendsRelationTypeAddedIssueParticipant
		}
	}

}

//删除任务
func assemblyDeleteIssueTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBos := []bo.TrendsBo{}
	ext := bo.TrendExtensionBo{}
	ext.IssueType = "T"
	ext.ObjName = issueTrendsBo.IssueTitle
	trendsBos = append(trendsBos, bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationDelete,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.ParentIssueId,
		RelationType:    consts.TrendsRelationTypeDeleteIssue,
		RelationObjType: consts.TrendsOperObjectTypeIssue,
		NewValue:        nil,
		OldValue:        nil,
		Ext:             json.ToJsonIgnoreError(ext),
		Creator:         issueTrendsBo.OperatorId,
	})

	//if issueTrendsBo.ParentIssueId != 0 {
	//	//获取父任务的标题
	//	issueInfo := projectfacade.IssueInfo(projectvo.IssueInfoReqVo{UserId: issueTrendsBo.OperatorId, OrgId: issueTrendsBo.OrgId, IssueID: issueTrendsBo.ParentIssueId})
	//	if issueInfo.Failure() {
	//		log.Error(issueInfo.Error())
	//		return trendsBos, issueInfo.Error()
	//	}
	//	ext.ObjName = issueInfo.IssueInfo.Issue.Title
	//	ext.RelationIssue = bo.RelationIssue{Title: issueTrendsBo.IssueTitle, ID: issueTrendsBo.IssueId}
	//	trendsBos = append(trendsBos, bo.TrendsBo{
	//		OrgId:           issueTrendsBo.OrgId,
	//		Module1:         consts.TrendsModuleOrg,
	//		Module2Id:       issueTrendsBo.ProjectId,
	//		Module2:         consts.TrendsModuleProject,
	//		Module3Id:       issueTrendsBo.ParentIssueId,
	//		Module3:         consts.TrendsModuleIssue,
	//		OperCode:        consts.RoleOperationDelete,
	//		OperObjId:       issueTrendsBo.IssueId,
	//		OperObjType:     consts.TrendsOperObjectTypeIssue,
	//		OperObjProperty: consts.BlankString,
	//		RelationObjId:   issueTrendsBo.IssueId,
	//		RelationType:    consts.TrendsRelationTypeDeleteChildIssue,
	//		RelationObjType: consts.TrendsOperObjectTypeIssue,
	//		NewValue:        nil,
	//		OldValue:        nil,
	//		Ext:             json.ToJsonIgnoreError(ext),
	//		Creator:         issueTrendsBo.OperatorId,
	//	})
	//}

	return trendsBos, nil
}

//更新任务状态
func assemblyUpdateIssueStatusTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	newValue := issueTrendsBo.NewValue
	oldValue := issueTrendsBo.OldValue
	operObjProperty := issueTrendsBo.OperateObjProperty

	trendsBo := &bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModifyStatus,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: operObjProperty,
		RelationObjId:   0,
		RelationType:    consts.TrendsRelationTypeUpdateIssueStatus,
		RelationObjType: consts.BlankString,
		NewValue:        &newValue,
		OldValue:        &oldValue,
		Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator:         issueTrendsBo.OperatorId,
	}

	resp := processfacade.GetProcessStatusIds(processvo.GetProcessStatusIdsReqVo{
		OrgId:    issueTrendsBo.OrgId,
		Category: consts.ProcessStatusCategoryIssue,
		Typ:      consts.ProcessStatusTypeCompleted,
	})
	if resp.Failure() {
		//不影响主流程
		log.Error(resp.Error())
		return []bo.TrendsBo{*trendsBo}, nil
	}
	if resp.ProcessStatusIds == nil {
		return []bo.TrendsBo{*trendsBo}, nil
	}
	for _, v := range *resp.ProcessStatusIds {
		if v == issueTrendsBo.IssueStatusId {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Errorf("捕获到的错误：%s", r)
					}
				}()
				noticeErr := AddNoticeByChangeUpdateIssueStatus(issueTrendsBo)
				if noticeErr != nil {
					log.Error(noticeErr)
				}
			}()

			break
		}
	}

	return []bo.TrendsBo{*trendsBo}, nil
}

func assemblyUpdateRelationIssueTrends(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	allRelateIssues := append(issueTrendsBo.BindIssues, issueTrendsBo.UnbindIssues...)
	if len(allRelateIssues) <= 0 {
		return nil, nil
	}
	//获取问题信息
	issueInfo := projectfacade.GetIssueInfoList(projectvo.IssueInfoListReqVo{
		IssueIds: allRelateIssues,
	})
	if issueInfo.Failure() {
		log.Error(issueInfo.Error())
		return nil, issueInfo.Error()
	}
	issueInfoById := map[int64]vo.Issue{}
	for _, v := range issueInfo.IssueInfos {
		issueInfoById[v.ID] = v
	}
	trendsBos := []bo.TrendsBo{}
	if len(issueTrendsBo.UnbindIssues) > 0 {
		deleteIds := slice.SliceUniqueInt64(issueTrendsBo.UnbindIssues)
		for _, v := range deleteIds {
			if _, ok := issueInfoById[v]; !ok {
				continue
			}
			issueTrendsBo.Ext.RelationIssue = bo.RelationIssue{
				ID:    v,
				Title: issueInfoById[v].Title,
			}

			trendsBos = append(trendsBos, bo.TrendsBo{
				OrgId:       issueTrendsBo.OrgId,
				Module1:     consts.TrendsModuleOrg,
				Module2Id:   issueTrendsBo.ProjectId,
				Module2:     consts.TrendsModuleProject,
				Module3Id:   issueTrendsBo.IssueId,
				Module3:     consts.TrendsModuleIssue,
				OperCode:    consts.RoleOperationUnbind,
				OperObjId:   issueTrendsBo.IssueId,
				OperObjType: consts.TrendsOperObjectTypeIssue,
				//OperObjProperty: operObjProperty,
				RelationObjId:   v,
				RelationType:    consts.TrendsRelationTypeDeleteRelationIssue,
				RelationObjType: consts.BlankString,
				Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
				Creator:         issueTrendsBo.OperatorId,
			})
		}
	}

	if len(issueTrendsBo.BindIssues) > 0 {
		addIds := slice.SliceUniqueInt64(issueTrendsBo.BindIssues)
		for _, v := range addIds {
			if _, ok := issueInfoById[v]; !ok {
				continue
			}
			issueTrendsBo.Ext.RelationIssue = bo.RelationIssue{
				ID:    v,
				Title: issueInfoById[v].Title,
			}

			trendsBos = append(trendsBos, bo.TrendsBo{
				OrgId:       issueTrendsBo.OrgId,
				Module1:     consts.TrendsModuleOrg,
				Module2Id:   issueTrendsBo.ProjectId,
				Module2:     consts.TrendsModuleProject,
				Module3Id:   issueTrendsBo.IssueId,
				Module3:     consts.TrendsModuleIssue,
				OperCode:    consts.RoleOperationBind,
				OperObjId:   issueTrendsBo.IssueId,
				OperObjType: consts.TrendsOperObjectTypeIssue,
				//OperObjProperty: operObjProperty,
				RelationObjId:   v,
				RelationType:    consts.TrendsRelationTypeAddRelationIssue,
				RelationObjType: consts.BlankString,
				Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
				Creator:         issueTrendsBo.OperatorId,
			})
		}
	}
	return trendsBos, nil
}

//上传附件
func assemblyUploadResource(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBo := &bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeUploadResource,
		RelationObjType: consts.BlankString,
		//NewValue:        &commentBo.Content,
		//OldValue:        nil,
		Ext:     json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator: issueTrendsBo.OperatorId,
	}

	return []bo.TrendsBo{*trendsBo}, nil
}

//变更任务栏
func assemblyUpdateIssueProjectObjectType(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBo := &bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeUpdateIssueProjectObjectType,
		RelationObjType: consts.BlankString,
		NewValue:        &issueTrendsBo.NewValue,
		OldValue:        &issueTrendsBo.OldValue,
		Ext:             json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator:         issueTrendsBo.OperatorId,
	}

	return []bo.TrendsBo{*trendsBo}, nil
}

//删除附件
func assemblyDeleteResource(issueTrendsBo bo.IssueTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	trendsBo := &bo.TrendsBo{
		OrgId:           issueTrendsBo.OrgId,
		Module1:         consts.TrendsModuleOrg,
		Module2Id:       issueTrendsBo.ProjectId,
		Module2:         consts.TrendsModuleProject,
		Module3Id:       issueTrendsBo.IssueId,
		Module3:         consts.TrendsModuleIssue,
		OperCode:        consts.RoleOperationModify,
		OperObjId:       issueTrendsBo.IssueId,
		OperObjType:     consts.TrendsOperObjectTypeIssue,
		OperObjProperty: consts.BlankString,
		RelationObjId:   issueTrendsBo.IssueId,
		RelationType:    consts.TrendsRelationTypeDeleteResource,
		RelationObjType: consts.BlankString,
		//NewValue:        &commentBo.Content,
		//OldValue:        nil,
		Ext:     json.ToJsonIgnoreError(issueTrendsBo.Ext),
		Creator: issueTrendsBo.OperatorId,
	}

	return []bo.TrendsBo{*trendsBo}, nil
}
