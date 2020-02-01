package domain

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/trendsvo"
	"github.com/galaxy-book/polaris-backend/facade/trendsfacade"
	"time"
)

func PushOrgTrends(orgTrendsBo bo.OrgTrendsBo) {
	orgTrendsBo.OperateTime = time.Now()
	//动态改成同步的
	resp := trendsfacade.AddOrgTrends(trendsvo.AddOrgTrendsReqVo{OrgTrendsBo: orgTrendsBo})
	if resp.Failure(){
		log.Error(resp.Message)
	}
}