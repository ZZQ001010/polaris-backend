package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/trendsvo"
	"github.com/galaxy-book/polaris-backend/service/platform/trendssvc/service"
)

func (PostGreeter) TrendList(reqVo trendsvo.TrendListReqVo) trendsvo.TrendListRespVo {
	trendList, err := service.TrendList(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return trendsvo.TrendListRespVo{Err: vo.NewErr(err), TrendList: trendList}
}

func (PostGreeter) AddIssueTrends(req trendsvo.AddIssueTrendsReqVo) vo.VoidErr {
	service.AddIssueTrends(req.IssueTrendsBo)
	return vo.VoidErr{Err: vo.NewErr(nil)}
}

//添加项目趋势
func (PostGreeter) AddProjectTrends(req trendsvo.AddProjectTrendsReqVo) vo.VoidErr {
	service.AddProjectTrends(req.ProjectTrendsBo)
	return vo.VoidErr{Err: vo.NewErr(nil)}
}

func (PostGreeter) AddOrgTrends(req trendsvo.AddOrgTrendsReqVo) vo.VoidErr {
	service.AddOrgTrends(req.OrgTrendsBo)
	return vo.VoidErr{Err: vo.NewErr(nil)}
}

func (PostGreeter) CreateTrends(req trendsvo.CreateTrendsReqVo) trendsvo.CreateTrendsRespVo {
	trendsId, err := service.CreateTrends(req.TrendsBo)
	return trendsvo.CreateTrendsRespVo{TrendsId: trendsId, Err: vo.NewErr(err)}
}
