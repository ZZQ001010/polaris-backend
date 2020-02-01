package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/date"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"time"
	"upper.io/db.v3"
)

func GetIssueStatusStat(issueStatusStatCondBo bo.IssueStatusStatCondBo) ([]bo.IssueStatusStatBo, errs.SystemErrorInfo) {
	orgId := issueStatusStatCondBo.OrgId

	//TODO 查询时将项目对象类型的name和langCode也查出来
	//SLOW
	issueAndDetailUnionBos, err1 := GetIssueAndDetailUnionBoList(bo.IssueBoListCond{
		OrgId:        orgId,
		ProjectId:    issueStatusStatCondBo.ProjectId,
		IterationId:  issueStatusStatCondBo.IterationId,
		RelationType: issueStatusStatCondBo.RelationType,
		UserId:       &issueStatusStatCondBo.UserId,
	})

	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	statusBoList, err2 := processfacade.GetProcessStatusListByCategoryRelaxed(orgId, consts.ProcessStatusCategoryIssue)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err2)
	}

	projectObjectTypeIds := make([]int64, 0)
	for _, issueUnionInfo := range issueAndDetailUnionBos {
		exist, err := slice.Contain(projectObjectTypeIds, issueUnionInfo.IssueProjectObjectTypeId)
		if err != nil {
			log.Error(err)
			continue
		}
		if !exist {
			projectObjectTypeIds = append(projectObjectTypeIds, issueUnionInfo.IssueProjectObjectTypeId)
		}
	}

	projectObjectTypeListPos := &[]po.PpmPrsProjectObjectType{}
	err := mysql.SelectAllByCond(consts.TableProjectObjectType, db.Cond{
		consts.TcOrgId: orgId,
		consts.TcId:    db.In(projectObjectTypeIds),
		//consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus: consts.AppStatusEnable,
	}, projectObjectTypeListPos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	projectObjectTypeBos := &[]bo.ProjectObjectTypeBo{}
	_ = copyer.Copy(projectObjectTypeListPos, projectObjectTypeBos)

	statusBoLocalCache := maps.NewMap("StatusId", statusBoList)
	projectObjectTypeLocalCache := maps.NewMap("Id", projectObjectTypeBos)
	return dealStat(issueAndDetailUnionBos, statusBoLocalCache, projectObjectTypeLocalCache)
}

func dealStat(issueAndDetailUnionBos []bo.IssueAndDetailUnionBo, statusBoLocalCache, projectObjectTypeLocalCache maps.LocalMap) ([]bo.IssueStatusStatBo, errs.SystemErrorInfo) {
	extData := map[int64]bo.IssueStatusStatBo{}
	for _, targetBo := range issueAndDetailUnionBos {
		projectObjectTypeBo, statusBo := assemblyObject(targetBo, statusBoLocalCache, projectObjectTypeLocalCache)

		if projectObjectTypeBo == nil || statusBo == nil {
			continue
		}

		key := projectObjectTypeBo.Id
		data, ok := extData[key]
		if !ok {
			data = bo.IssueStatusStatBo{
				ProjectTypeId:       key,
				ProjectTypeName:     projectObjectTypeBo.Name,
				ProjectTypeLangCode: projectObjectTypeBo.LangCode,
			}
		}

		data.IssueCount++
		data.StoryPointCount++

		todayZero := date.GetZeroTime(time.Now())
		now := time.Now()
		tomorrowZero := todayZero.AddDate(0, 0, 1)
		theDayAfterTomorrowZero := tomorrowZero.AddDate(0, 0, 1)

		if statusBo.StatusType == consts.ProcessStatusTypeNotStarted {
			data.IssueWaitCount++
			data.StoryPointWaitCount++
		} else if statusBo.StatusType == consts.ProcessStatusTypeProcessing {
			data.IssueRunningCount++
			data.StoryPointRunningCount++
		} else if statusBo.StatusType == consts.ProcessStatusTypeCompleted {
			data.IssueEndCount++
			data.StoryPointEndCount++

			//今日完成
			if targetBo.EndTime.After(todayZero) && targetBo.EndTime.Before(tomorrowZero) {
				data.IssueEndTodayCount++
			}
		}

		if targetBo.CreateTime.After(todayZero) && targetBo.CreateTime.Before(tomorrowZero){
			//今日创建
			data.TodayCreateCount ++
		}

		if statusBo.StatusType != consts.ProcessStatusTypeCompleted {
			if targetBo.PlanEndTime.After(consts.BlankTimeObject) && targetBo.PlanEndTime.Before(time.Now()) {
				//已逾期
				data.IssueOverdueCount++
			}
			if targetBo.PlanEndTime.Equal(now) || (targetBo.PlanEndTime.After(now) && targetBo.PlanEndTime.Before(tomorrowZero)) {
				//今日到期
				data.IssueOverdueTodayCount++
			}
			//明日逾期
			if targetBo.PlanEndTime.Equal(tomorrowZero) || (targetBo.PlanEndTime.After(tomorrowZero) && targetBo.PlanEndTime.Before(theDayAfterTomorrowZero)) {
				//即将逾期
				data.IssueOverdueTomorrowCount++
			}

			if targetBo.OwnerChangeTime.After(todayZero) && targetBo.OwnerChangeTime.Before(tomorrowZero) {
				//今日指派给我
				data.TodayCount++
			}

		}
		//逾期完成
		if statusBo.StatusType == consts.ProcessStatusTypeCompleted && targetBo.PlanEndTime.Before(targetBo.EndTime) && targetBo.PlanEndTime.After(consts.BlankTimeObject) {
			data.IssueOverdueEndCount++
		}

		extData[projectObjectTypeBo.Id] = data
	}

	statBos := make([]bo.IssueStatusStatBo, 0)
	for _, v := range extData {
		statBos = append(statBos, v)
	}
	return statBos, nil
}

func assemblyObject(targetBo bo.IssueAndDetailUnionBo, statusBoLocalCache, projectObjectTypeLocalCache maps.LocalMap) (
	*bo.ProjectObjectTypeBo, *bo.CacheProcessStatusBo) {
	var statusResultBo *bo.CacheProcessStatusBo = nil
	var projectObjectTypeResultBo *bo.ProjectObjectTypeBo = nil
	if statusBoCacheObj, ok := statusBoLocalCache[targetBo.IssueStatusId]; ok {
		if statusBo, ok := statusBoCacheObj.(bo.CacheProcessStatusBo); ok {
			statusResultBo = &statusBo
		}
	}

	if projectObjectTypeCacheObj, ok := projectObjectTypeLocalCache[targetBo.IssueProjectObjectTypeId]; ok {
		if projectObjectTypeBo, ok := projectObjectTypeCacheObj.(bo.ProjectObjectTypeBo); ok {
			projectObjectTypeResultBo = &projectObjectTypeBo
		}
	}
	return projectObjectTypeResultBo, statusResultBo
}
