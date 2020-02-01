package domain

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"time"
	"upper.io/db.v3"
)

func SelectIssueAssignCount(orgId int64, projectId int64, rankTop int) ([]bo.IssueAssignCountBo, errs.SystemErrorInfo) {
	finishedIds, err2 := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
	if err2 != nil || len(*finishedIds) == 0 {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
	}
	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if conn != nil {
				if err := conn.Close(); err != nil {
					logger.GetDefaultLogger().Info(strs.ObjectToString(err))
				}
			}
		}
	}()
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	issueAssignCountBos := &[]bo.IssueAssignCountBo{}
	err1 := conn.Select(db.Raw("count(*) as count"), consts.TcOwner).From(consts.TableIssue).Where(db.Cond{
		consts.TcProjectId: projectId,
		consts.TcStatus:    db.NotIn(finishedIds),
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}).GroupBy(consts.TcOwner).OrderBy("count desc").Limit(rankTop).All(issueAssignCountBos)

	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	return *issueAssignCountBos, nil
}

func SelectIssueAssignTotalCount(orgId int64, projectId int64) (int64, errs.SystemErrorInfo) {
	finishedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
	if err != nil || len(*finishedIds) == 0 {
		log.Errorf("proxies.GetProcessStatusId: %q\n", err)
		return 0, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	total, err1 := mysql.SelectCountByCond(consts.TableIssue, db.Cond{
		consts.TcProjectId: projectId,
		consts.TcStatus:    db.NotIn(finishedIds),
		consts.TcIsDelete:  consts.AppIsNoDelete,
	})
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return int64(total), nil
}

func GetIssueCountByOwnerId(orgId, ownerId int64) (int64, errs.SystemErrorInfo) {
	finishedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
	if err != nil || len(*finishedIds) == 0 {
		log.Errorf("proxies.GetProcessStatusId: %q\n", err)
		return 0, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	total, err1 := mysql.SelectCountByCond(consts.TableIssue, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcOwner:    ownerId,
		consts.TcStatus:   db.NotIn(finishedIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return int64(total), nil
}

//任务完成统计图，根据endTime做groupBy，统计当前负责人在指定时间段之内的每天任务完成数量
func IssueDailyPersonalWorkCompletionStat(orgId, ownerId int64, startDate *types.Time, endDate *types.Time) ([]bo.IssueDailyPersonalWorkCompletionStatBo, errs.SystemErrorInfo) {
	finishedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
	if err != nil || len(*finishedIds) == 0 {
		log.Errorf("proxies.GetProcessStatusId: %q\n", err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	//默认查询七天
	startTime, endTime, dateErr := CalDateRangeCond(startDate, endDate, 7)
	if dateErr != nil {
		log.Error(dateErr)
		return nil, dateErr
	}

	//如果and后的日期是到天的，需要加一天
	condEndTime := endTime.AddDate(0, 0, 1)

	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcOwner:    ownerId,
		consts.TcStatus:   db.In(finishedIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcEndTime:  db.Between(startTime.Format(consts.AppDateFormat), condEndTime.Format(consts.AppDateFormat)),
	}
	IssueCondFiling(cond, orgId, consts.AppIsNotFilling)

	conn, mysqlErr := mysql.GetConnect()
	conn.SetLogging(true)
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if mysqlErr != nil {
		log.Error(mysqlErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	resultBos := &[]bo.IssueDailyPersonalWorkCompletionStatBo{}
	mysqlErr = conn.Select(db.Raw("count(*) as count"), db.Raw("date(end_time) as date")).From(consts.TableIssue).Where(cond).GroupBy(db.Raw("date(end_time)")).All(resultBos)
	if mysqlErr != nil {
		log.Error(mysqlErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	resultBoMap := map[string]bo.IssueDailyPersonalWorkCompletionStatBo{}
	for _, resultBo := range *resultBos {
		date := resultBo.Date.Format(consts.AppDateFormat)
		resultBoMap[date] = resultBo
	}

	//处理结果集，补填空缺
	afterDealResultBos := make([]bo.IssueDailyPersonalWorkCompletionStatBo, 0)
	cursorTime := *startTime
	cursorTimeDateFormat := cursorTime.Format(consts.AppDateFormat)
	cursorTime, _ = time.Parse(consts.AppDateFormat, cursorTimeDateFormat)
	endTimeDateFormat := endTime.Format(consts.AppDateFormat)
	refEndTime, _ := time.Parse(consts.AppDateFormat, endTimeDateFormat)
	for {
		if cursorTime.After(refEndTime) {
			break
		}
		var targetBo *bo.IssueDailyPersonalWorkCompletionStatBo = nil
		cursorTimeDateFormat := cursorTime.Format(consts.AppDateFormat)
		if resultBo, ok := resultBoMap[cursorTimeDateFormat]; ok {
			targetBo = &resultBo
		} else {
			targetBo = &bo.IssueDailyPersonalWorkCompletionStatBo{
				Count: 0,
				Date:  cursorTime,
			}
		}
		afterDealResultBos = append(afterDealResultBos, *targetBo)
		cursorTime = cursorTime.AddDate(0, 0, 1)
	}
	return afterDealResultBos, nil
}

//计算日期范围条件
//startDate: 开始时间
//endDate: 结束时间
//cond: 条件
//defaultDay: 默认查询多少天前
func CalDateRangeCond(startDate *types.Time, endDate *types.Time, defaultDay int) (*time.Time, *time.Time, errs.SystemErrorInfo) {
	if startDate == nil && endDate == nil {
		currentTime := time.Now()
		sd := types.Time(currentTime.AddDate(0, 0, -(defaultDay - 1)))
		ed := types.Time(currentTime)
		startDate = &sd
		endDate = &ed
	}
	if startDate != nil && endDate != nil {
	} else {
		log.Error("没有提供明确的时间范围")
		return nil, nil, errs.BuildSystemErrorInfo(errs.DateRangeError)
	}
	startTime := time.Time(*startDate)
	endTime := time.Time(*endDate)

	return &startTime, &endTime, nil
}
