package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) ProjectDayStats(params *projectvo.ProjectDayStatsReqVo) projectvo.ProjectDayStatsRespVo {
	list, err := service.ProjectDayStats(params.OrgId, params.Page, params.Size, params.Params)
	return projectvo.ProjectDayStatsRespVo{Err: vo.NewErr(err), ProjectDayStatList: list}
}

//date: yyyy-MM-dd
func (PostGreeter) AppendProjectDayStat(req projectvo.AppendProjectDayStatReqVo) vo.VoidErr {
	err := service.AppendProjectDayStat(req.ProjectBo, req.Date)
	return vo.VoidErr{Err: vo.NewErr(err)}
}
