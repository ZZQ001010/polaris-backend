package domain

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/maps"
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

func StatisticIterationCountGroupByStatus(orgId, projectId int64) (*bo.IterationStatusTypeCountBo, errs.SystemErrorInfo) {
	var notStartTotal, processingTotal, finishedTotal int64 = 0, 0, 0
	iterationQueryCond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if projectId > 0 {
		iterationQueryCond[consts.TcProjectId] = projectId
	}
	iterationPos, err := dao.SelectIteration(iterationQueryCond)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	cacheStatusBos, err2 := processfacade.GetProcessStatusListByCategoryRelaxed(orgId, consts.ProcessStatusCategoryIteration)
	if err2 != nil {
		log.Error(err2)
		return nil, err2
	}
	statusCacheMap := maps.NewMap("StatusId", cacheStatusBos)

	for _, iterationPo := range *iterationPos {

		value, ok := statusCacheMap[iterationPo.Status]

		if !ok {
			continue
		}
		statusBo := value.(bo.CacheProcessStatusBo)
		switch statusBo.StatusType {
		case consts.ProcessStatusTypeNotStarted:
			notStartTotal++
		case consts.ProcessStatusTypeProcessing:
			processingTotal++
		case consts.ProcessStatusTypeCompleted:
			finishedTotal++
		}
	}

	return &bo.IterationStatusTypeCountBo{
		NotStartTotal:   notStartTotal,
		ProcessingTotal: processingTotal,
		FinishedTotal:   finishedTotal,
	}, nil
}

//date: yyyy-MM-dd
func AppendIterationStat(iterationBo bo.IterationBo, date string, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	statDate, timeParseError := time.Parse(consts.AppDateFormat, date)
	if timeParseError != nil {
		log.Error(timeParseError)
		return errs.BuildSystemErrorInfo(errs.SystemError)
	}

	orgId := iterationBo.OrgId
	projectId := iterationBo.ProjectId
	iterationId := iterationBo.Id

	id, err3 := idfacade.ApplyPrimaryIdRelaxed(consts.TableIterationStat)
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
		OrgId:       orgId,
		ProjectId:   &projectId,
		IterationId: &iterationId,
	})
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	iterationStatPo := &po.PpmStaIterationStat{}
	iterationStatPo.Id = id
	iterationStatPo.OrgId = orgId
	iterationStatPo.ProjectId = projectId
	iterationStatPo.IterationId = iterationId
	iterationStatPo.StatDate = statDate
	iterationStatPo.Status = projectCacheBo.Status

	issueStatusStatMap := maps.NewMap("ProjectTypeId", issueStatusStatBos)
	ext := bo.StatExtBo{}
	ext.Issue = bo.StatIssueExtBo{
		Data: issueStatusStatMap,
	}
	iterationStatPo.Ext = json.ToJsonIgnoreError(ext)
	//封装状态
	assemblyIterationStat(issueStatusStatBos, iterationStatPo)

	err5 := dao.InsertIterationStat(*iterationStatPo, tx...)
	if err5 != nil {
		log.Error(err5)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func assemblyIterationStat(issueStatusStatBos []bo.IssueStatusStatBo, iterationStatPo *po.PpmStaIterationStat) {
	for _, statBo := range issueStatusStatBos {
		switch statBo.ProjectTypeLangCode {
		case consts.ProjectObjectTypeLangCodeDemand:
			iterationStatPo.DemandCount += statBo.IssueCount
			iterationStatPo.DemandWaitCount += statBo.IssueWaitCount
			iterationStatPo.DemandRunningCount += statBo.IssueRunningCount
			iterationStatPo.DemandEndCount += statBo.IssueEndCount
		case consts.ProjectObjectTypeLangCodeTask:
			iterationStatPo.TaskCount += statBo.IssueCount
			iterationStatPo.TaskWaitCount += statBo.IssueWaitCount
			iterationStatPo.TaskRunningCount += statBo.IssueRunningCount
			iterationStatPo.TaskEndCount += statBo.IssueEndCount
		case consts.ProjectObjectTypeLangCodeBug:
			iterationStatPo.BugCount += statBo.IssueCount
			iterationStatPo.BugWaitCount += statBo.IssueWaitCount
			iterationStatPo.BugRunningCount += statBo.IssueRunningCount
			iterationStatPo.BugEndCount += statBo.IssueEndCount
		case consts.ProjectObjectTypeLangCodeTestTask:
			iterationStatPo.TesttaskCount += statBo.IssueCount
			iterationStatPo.TesttaskWaitCount += statBo.IssueWaitCount
			iterationStatPo.TesttaskRunningCount += statBo.IssueRunningCount
			iterationStatPo.TesttaskEndCount += statBo.IssueEndCount
		}
		iterationStatPo.IssueCount += statBo.IssueCount
		iterationStatPo.IssueWaitCount += statBo.IssueWaitCount
		iterationStatPo.IssueRunningCount += statBo.IssueRunningCount
		iterationStatPo.IssueEndCount += statBo.IssueEndCount
		iterationStatPo.StoryPointCount += statBo.StoryPointCount
		iterationStatPo.StoryPointWaitCount += statBo.StoryPointWaitCount
		iterationStatPo.StoryPointRunningCount += statBo.StoryPointRunningCount
		iterationStatPo.StoryPointEndCount += statBo.StoryPointEndCount
	}
}
