package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"upper.io/db.v3"
)

func ProjectDayStats(orgId int64, page uint, size uint, params *vo.ProjectDayStatReq) (*vo.ProjectDayStatList, errs.SystemErrorInfo) {

	projectId := params.ProjectID

	projectBo, err1 := domain.GetProject(orgId, projectId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityProject)
	}

	//默认查询十五天
	startTime, endTime, err := domain.CalDateRangeCond(params.StartDate, params.EndDate, 15)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//如果and后的日期是到天的，需要加一天
	condEndTime := endTime.AddDate(0, 0, 1)

	cond := db.Cond{
		consts.TcProjectId: projectBo.Id,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcStatDate:  db.Between(startTime.Format(consts.AppDateFormat), condEndTime.Format(consts.AppDateFormat)),
	}

	bos, total, err1 := domain.GetProjectDayStatBoList(page, size, cond)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	resultList := &[]*vo.ProjectDayStat{}
	err3 := copyer.Copy(bos, resultList)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return &vo.ProjectDayStatList{
		Total: total,
		List:  *resultList,
	}, nil
}

//date: yyyy-MM-dd
func AppendProjectDayStat(projectBo bo.ProjectBo, date string) errs.SystemErrorInfo {
	return domain.AppendProjectDayStat(projectBo, date)
}
