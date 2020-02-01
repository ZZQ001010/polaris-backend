package service

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/trendssvc/domain"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

type Ext struct {
	ObjName string
}

func TrendList(orgId, currentUserId int64, input *vo.TrendReq) (*vo.TrendsList, errs.SystemErrorInfo) {

	log.Infof(consts.UserLoginSentence, currentUserId, orgId)

	trendBo := &bo.TrendsQueryCondBo{
		LastTrendID: input.LastTrendID,
		OrgId:       orgId,
		ObjId:       input.ObjID,
		ObjType:     input.ObjType,
		OperId:      input.OperID,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Type:        input.Type,
		Page:        input.Page,
		Size:        input.Size,
		OrderType:   input.OrderType,
	}
	result, err := domain.QueryTrends(trendBo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TrendDomainError, err)
	}

	resultModel := &vo.TrendsList{}
	copyer.Copy(result, resultModel)
	creatorIds := []int64{}
	for _, v := range resultModel.List {
		creatorIds = append(creatorIds, v.Creator)
	}
	userInfo, err := orgfacade.GetBaseUserInfoBatchRelaxed("", orgId, creatorIds)
	if err != nil {
		return nil, err
	}
	userInfoById := map[int64]bo.BaseUserInfoBo{}
	for _, v := range userInfo {
		userInfoById[v.UserId] = v
	}
	issueIds := []int64{}
	projectIds := []int64{}
	projectObjectTypeIds := []int64{}
	for _, v := range resultModel.List {
		if v.Module3 == consts.TrendsModuleIssue && v.Module3Id != 0 {
			issueIds = append(issueIds, v.Module3Id)
		}
		if v.Module2 == consts.TrendsModuleProject && v.Module2Id != 0 && v.Module3 == consts.BlankString {
			projectIds = append(projectIds, v.Module2Id)
		}
		//变更任务栏获取最新信息
		if v.RelationType == consts.TrendsRelationTypeUpdateIssueProjectObjectType {
			if v.OldValue != nil {
				oldBo := &bo.ProjectObjectTypeAndProjectIdBo{}
				_ = json.FromJson(*v.OldValue, oldBo)
				if oldBo.ProjectId != 0 {
					projectIds = append(projectIds, oldBo.ProjectId)
				}
				if oldBo.ProjectObjectTypeId != 0 {
					projectObjectTypeIds = append(projectObjectTypeIds, oldBo.ProjectObjectTypeId)
				}
			}

			if v.NewValue != nil {
				newBo := &bo.ProjectObjectTypeAndProjectIdBo{}
				_ = json.FromJson(*v.NewValue, newBo)
				if newBo.ProjectId != 0 {
					projectIds = append(projectIds, newBo.ProjectId)
				}
				if newBo.ProjectObjectTypeId != 0 {
					projectObjectTypeIds = append(projectObjectTypeIds, newBo.ProjectObjectTypeId)
				}
			}

		}
	}
	//获取最新任务信息
	issueInfos := projectfacade.GetSimpleIssueInfoBatch(projectvo.GetSimpleIssueInfoBatchReqVo{OrgId: orgId, Ids: issueIds})
	if issueInfos.Failure() {
		log.Error(issueInfos.Error())
		return nil, issueInfos.Error()
	}
	issueInfosMap := maps.NewMap("ID", issueInfos.Data)

	//获取最新项目信息
	projectResp := projectfacade.GetSimpleProjectInfo(projectvo.GetSimpleProjectInfoReqVo{OrgId: orgId, Ids: slice.SliceUniqueInt64(projectIds)})
	if projectResp.Failure() {
		log.Error(projectResp.Error())
		return nil, projectResp.Error()
	}
	projectInfo := map[int64]string{}
	for _, v := range *projectResp.Data {
		projectInfo[v.ID] = v.Name
	}

	//获取任务栏信息
	projectObjectTypeResp := projectfacade.ProjectObjectTypeList(projectvo.ProjectObjectTypesReqVo{
		OrgId:  orgId,
		Params: &vo.ProjectObjectTypesReq{ObjectType: 2, Ids: projectObjectTypeIds},
	})
	if projectObjectTypeResp.Failure() {
		log.Error(projectObjectTypeResp.Error())
		return nil, projectObjectTypeResp.Error()
	}
	projectObjectTypeInfo := map[int64]string{}
	for _, v := range projectObjectTypeResp.ProjectObjectTypeList.List {
		projectObjectTypeInfo[v.ID] = v.Name
	}

	var lastId int64
	for k, v := range resultModel.List {
		if _, ok := userInfoById[v.Creator]; ok {
			(resultModel.List)[k].CreatorInfo = assemblyUserIdInfo(userInfoById[v.Creator])
		}
		ext := &vo.TrendExtension{}
		err := json.FromJson(v.Ext, ext)
		if err == nil {
			if (*ext).ObjName == nil {
				(resultModel.List)[k].OperObjName = consts.BlankString
			} else {
				(resultModel.List)[k].OperObjName = *ext.ObjName
			}

			(resultModel.List)[k].Extension = ext
		}

		//如果是任务，获取最新任务名称
		if v.Module3 == consts.TrendsModuleIssue && v.Module3Id != 0 {
			if _, ok := issueInfosMap[v.Module3Id]; ok {
				issueTemp := issueInfosMap[v.Module3Id].(vo.Issue)
				(resultModel.List)[k].OperObjName = issueTemp.Title
			}
		}

		//如果是项目，获取最新项目名称
		if v.Module2 == consts.TrendsModuleProject && v.Module2Id != 0 && v.Module3 == consts.BlankString {
			if _, ok := projectInfo[v.Module2Id]; ok {
				(resultModel.List)[k].OperObjName = projectInfo[v.Module2Id]
			}
		}

		if v.RelationType == consts.TrendsRelationTypeCreateIssueComment {
			if *v.NewValue != "" {
				(resultModel.List)[k].Comment = v.NewValue
			}
		} else if v.RelationType == consts.TrendsRelationTypeUpdateIssueProjectObjectType && v.OldValue != nil && v.NewValue != nil {
			//变更任务栏 重新拼接字符串
			oldBo := &bo.ProjectObjectTypeAndProjectIdBo{}
			_ = json.FromJson(*v.OldValue, oldBo)
			oldValue := ""
			if _, ok := projectObjectTypeInfo[oldBo.ProjectObjectTypeId]; ok {
				oldValue += projectObjectTypeInfo[oldBo.ProjectObjectTypeId]
			} else if v.Extension.ChangeList[0].OldValue != nil {
				oldValue = *v.Extension.ChangeList[0].OldValue
			}

			newBo := &bo.ProjectObjectTypeAndProjectIdBo{}
			_ = json.FromJson(*v.NewValue, newBo)
			newValue := ""
			if _, ok := projectObjectTypeInfo[newBo.ProjectObjectTypeId]; ok {
				newValue += projectObjectTypeInfo[newBo.ProjectObjectTypeId]
			} else if v.Extension.ChangeList[0].NewValue != nil {
				newValue = *v.Extension.ChangeList[0].NewValue
			}

			if oldBo.ProjectId != newBo.ProjectId {
				if _, ok := projectInfo[oldBo.ProjectId]; ok {
					oldValue += "(" + projectInfo[oldBo.ProjectId] + ")"
				}
				if _, ok := projectInfo[newBo.ProjectId]; ok {
					newValue += "(" + projectInfo[newBo.ProjectId] + ")"
				}
			}

			(resultModel.List)[k].Extension.ChangeList[0].OldValue = &oldValue
			(resultModel.List)[k].Extension.ChangeList[0].NewValue = &newValue
		}
		if len(v.Extension.ChangeList) > 0 {
			for i, list := range v.Extension.ChangeList {
				if *list.Field == "planEndTime" {
					name := consts.PlanEndTime
					(resultModel.List)[k].Extension.ChangeList[i].FieldName = &name
				}
			}
		}

		//如果查询的是任务（多加个判断是否是子任务处理）
		if input.ObjType != nil && *input.ObjType == consts.TrendsOperObjectTypeIssue && *input.ObjID == v.RelationObjID {
			if v.RelationType == consts.TrendsRelationTypeCreateIssue {
				(resultModel.List)[k].RelationType = consts.TrendsRelationTypeCreateChildIssue
			} else if v.RelationType == consts.TrendsRelationTypeDeleteIssue {
				(resultModel.List)[k].RelationType = consts.TrendsRelationTypeDeleteChildIssue
			}
		}
		lastId = v.ID
	}
	resultModel.LastTrendID = lastId
	return resultModel, nil
}

func getCommentByValue(value string) string {
	if value == "" {
		return ""
	}
	commentBo := &bo.CommentBo{}
	err := json.FromJson(value, commentBo)
	if err != nil {
		log.Error("转化comment失败： " + strs.ObjectToString(err))
		return ""
	}

	return commentBo.Content
}

func assemblyUserIdInfo(baseUserInfo bo.BaseUserInfoBo) *vo.UserIDInfo {
	return &vo.UserIDInfo{
		UserID: baseUserInfo.UserId,
		Name:   baseUserInfo.Name,
		Avatar: baseUserInfo.Avatar,
		EmplID: baseUserInfo.OutUserId,
	}
}

func AddIssueTrends(issueTrendsBo bo.IssueTrendsBo) {
	domain.AddIssueTrends(issueTrendsBo)
}

//添加项目趋势
func AddProjectTrends(projectTrendsBo bo.ProjectTrendsBo) {
	domain.AddProjectTrends(projectTrendsBo)
}

func AddOrgTrends(orgTrendsBo bo.OrgTrendsBo) {
	domain.AddOrgTrends(orgTrendsBo)
}

func CreateTrends(trendsBo *bo.TrendsBo, tx ...sqlbuilder.Tx) (*int64, errs.SystemErrorInfo) {
	return domain.CreateTrends(trendsBo, tx...)
}
