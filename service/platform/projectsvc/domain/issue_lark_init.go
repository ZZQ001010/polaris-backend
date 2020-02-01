package domain

import (
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/vo/processvo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"strconv"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
)

const larkIssueInitSql = consts.TemplateDirPrefix + "lark_issue_init.template"

func LarkIssueInit(orgId int64, zhangsanId, lisiId, projectId, operatorId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	contextMap := map[string]interface{}{}
	contextMap["OrgId"] = orgId
	contextMap["OperatorId"] = operatorId
	contextMap["ZhangSanId"] = zhangsanId
	contextMap["LiSiId"] = lisiId
	contextMap["ProjectId"] = projectId
	contextMap["NowTime"] = types.NowTime()

	statCount := 4
	//stat id 申请
	statIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectDayStat, statCount)
	if err != nil {
		return err
	}
	statIdCount := 1
	now := time.Now()
	for _, v := range statIds.Ids {
		contextMap["StatId"+strconv.Itoa(statIdCount)] = v.Id
		contextMap["Day"+strconv.Itoa(statIdCount)] = now.AddDate(0, 0, -statIdCount).Format(consts.AppDateFormat)
		statIdCount++
	}

	count := 8
	//issue id 申请
	issueIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssue, count)
	if err != nil {
		return err
	}
	issueIdCount := 1
	for _, v := range issueIds.Ids {
		contextMap["IssueId"+strconv.Itoa(issueIdCount)] = v.Id
		issueIdCount++
	}
	//issuedetail id 申请
	issueDetailIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueDetail, count)
	if err != nil {
		return err
	}
	issueDetailIdCount := 1
	for _, v := range issueDetailIds.Ids {
		contextMap["IssueDetailId"+strconv.Itoa(issueDetailIdCount)] = v.Id
		issueDetailIdCount++
	}

	//issuerelation id 申请
	relationIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueRelation, count)
	if err != nil {
		return err
	}
	relationIdCount := 1
	for _, v := range relationIds.Ids {
		contextMap["RelationId"+strconv.Itoa(relationIdCount)] = v.Id
		relationIdCount++
	}

	objectTypes, err := GetProjectObjectTypeList(orgId, projectId)
	if err != nil {
		return err
	}
	contextMap["ObjectTypeDemand"] = 0
	for _, v := range *objectTypes {
		if v.Name == "需求" {
			contextMap["ObjectTypeDemand"] = v.Id
			break
		}
	}

	//优先级
	allPriority, err := GetPriorityList(orgId)
	if err != nil {
		return err
	}
	for _, v := range *allPriority {
		if v.LangCode == "Priority.Issue.P0" {
			contextMap["PriorityP0"] = v.Id
		} else if v.LangCode == "Priority.Issue.P1" {
			contextMap["PriorityP1"] = v.Id
		} else if v.LangCode == "Priority.Issue.P3" {
			contextMap["PriorityP3"] = v.Id
		}
	}

	//状态
	processStatusResp := processfacade.GetProcessStatusListByCategory(processvo.GetProcessStatusListByCategoryReqVo{
		OrgId:    orgId,
		Category: 3,
	})
	if processStatusResp.Failure() {
		return processStatusResp.Error()
	}
	for _, v := range processStatusResp.CacheProcessStatusBoList {
		if v.Name == "未完成" {
			contextMap["NotStartStatus"] = v.StatusId
		} else if v.Name == "处理中" {
			contextMap["RunningStatus"] = v.StatusId
		}
	}

	insertErr := util.ReadAndWrite(larkIssueInitSql, contextMap, tx)
	if insertErr != nil {
		log.Error(insertErr)
		return errs.BuildSystemErrorInfo(errs.BaseDomainError, insertErr)
	}

	return nil
}
