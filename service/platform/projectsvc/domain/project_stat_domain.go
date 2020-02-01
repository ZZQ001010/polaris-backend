package domain

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//date: yyyy-MM-dd
func AppendProjectDayStat(projectBo bo.ProjectBo, date string, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	statDate, timeParseError := time.Parse(consts.AppDateFormat, date)
	if timeParseError != nil {
		log.Error(timeParseError)
		return errs.BuildSystemErrorInfo(errs.SystemError)
	}

	orgId := projectBo.OrgId
	projectId := projectBo.Id

	id, err3 := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectDayStat)
	if err3 != nil {
		log.Error(err3)
		return errs.BuildSystemErrorInfo(errs.ApplyIdError, err3)
	}

	projectCacheBo, err := LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	issueStatusStatBos, err1 := GetIssueStatusStat(bo.IssueStatusStatCondBo{
		OrgId:     orgId,
		ProjectId: &projectId,
	})
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	projectDayStatPo := &po.PpmStaProjectDayStat{}
	projectDayStatPo.Id = id
	projectDayStatPo.OrgId = orgId
	projectDayStatPo.ProjectId = projectId
	projectDayStatPo.StatDate = statDate
	projectDayStatPo.Status = projectCacheBo.Status

	issueStatusStatMap := maps.NewMap("ProjectTypeId", issueStatusStatBos)
	ext := bo.StatExtBo{}
	ext.Issue = bo.StatIssueExtBo{
		Data: issueStatusStatMap,
	}
	projectDayStatPo.Ext = json.ToJsonIgnoreError(ext)
	//封装状态
	assemblyProjectStat(issueStatusStatBos, projectDayStatPo)

	err5 := dao.InsertProjectDayStat(*projectDayStatPo, tx...)
	if err5 != nil {
		log.Error(err5)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func assemblyProjectStat(issueStatusStatBos []bo.IssueStatusStatBo, projectDayStatPo *po.PpmStaProjectDayStat) {
	for _, statBo := range issueStatusStatBos {
		switch statBo.ProjectTypeLangCode {
		case consts.ProjectObjectTypeLangCodeDemand:
			projectDayStatPo.DemandCount += statBo.IssueCount
			projectDayStatPo.DemandWaitCount += statBo.IssueWaitCount
			projectDayStatPo.DemandRunningCount += statBo.IssueRunningCount
			projectDayStatPo.DemandEndCount += statBo.IssueEndCount
			projectDayStatPo.DemandOverdueCount += statBo.IssueOverdueCount
		case consts.ProjectObjectTypeLangCodeTask:
			projectDayStatPo.TaskCount += statBo.IssueCount
			projectDayStatPo.TaskWaitCount += statBo.IssueWaitCount
			projectDayStatPo.TaskRunningCount += statBo.IssueRunningCount
			projectDayStatPo.TaskEndCount += statBo.IssueEndCount
			projectDayStatPo.TaskOverdueCount += statBo.IssueOverdueCount
		case consts.ProjectObjectTypeLangCodeBug:
			projectDayStatPo.BugCount += statBo.IssueCount
			projectDayStatPo.BugWaitCount += statBo.IssueWaitCount
			projectDayStatPo.BugRunningCount += statBo.IssueRunningCount
			projectDayStatPo.BugEndCount += statBo.IssueEndCount
			projectDayStatPo.BugOverdueCount += statBo.IssueOverdueCount
		case consts.ProjectObjectTypeLangCodeTestTask:
			projectDayStatPo.TesttaskCount += statBo.IssueCount
			projectDayStatPo.TesttaskWaitCount += statBo.IssueWaitCount
			projectDayStatPo.TesttaskRunningCount += statBo.IssueRunningCount
			projectDayStatPo.TesttaskEndCount += statBo.IssueEndCount
			projectDayStatPo.TesttaskOverdueCount += statBo.IssueOverdueCount
		}
		projectDayStatPo.IssueCount += statBo.IssueCount
		projectDayStatPo.IssueWaitCount += statBo.IssueWaitCount
		projectDayStatPo.IssueRunningCount += statBo.IssueRunningCount
		projectDayStatPo.IssueEndCount += statBo.IssueEndCount
		projectDayStatPo.IssueOverdueCount += statBo.IssueOverdueCount
		projectDayStatPo.StoryPointCount += statBo.StoryPointCount
		projectDayStatPo.StoryPointWaitCount += statBo.StoryPointWaitCount
		projectDayStatPo.StoryPointRunningCount += statBo.StoryPointRunningCount
		projectDayStatPo.StoryPointEndCount += statBo.StoryPointEndCount
	}
}

func GetProjectCountByOwnerId(orgId, ownerId int64) (int64, errs.SystemErrorInfo){
	finishedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryProject, consts.ProcessStatusTypeCompleted)
	if err != nil || len(*finishedIds) == 0 {
		log.Errorf("proxies.GetProcessStatusId: %q\n", err)
		return 0, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	total, err1 := mysql.SelectCountByCond(consts.TableProject, db.Cond{
		consts.TcOrgId: orgId,
		consts.TcOwner: ownerId,
		consts.TcStatus:    db.NotIn(finishedIds),
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcIsFiling: consts.AppIsNotFilling,
	})
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return int64(total), nil
}
