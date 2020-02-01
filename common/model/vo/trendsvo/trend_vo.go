package trendsvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type TrendListRespVo struct {
	vo.Err
	TrendList *vo.TrendsList `json:"trendList"`
}

type TrendListReqVo struct {
	Input  *vo.TrendReq
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type AddIssueTrendsReqVo struct {
	IssueTrendsBo bo.IssueTrendsBo `json:"issueTrendsBo"`
	Key           string           `json:"key"`
}

type AddProjectTrendsReqVo struct {
	ProjectTrendsBo bo.ProjectTrendsBo `json:"projectTrendsBo"`
}

type AddOrgTrendsReqVo struct {
	OrgTrendsBo bo.OrgTrendsBo `json:"orgTrendsBo"`
}

type CreateTrendsReqVo struct {
	TrendsBo *bo.TrendsBo `json:"trendsBo"`
}

type CreateTrendsRespVo struct {
	TrendsId *int64 `json:"data"`

	vo.Err
}
