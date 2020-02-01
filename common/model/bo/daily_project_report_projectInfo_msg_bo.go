package bo

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"time"
)

type DailyProjectReportProjectInfoMsgBo struct {
	ScheduleTraceId string                     `json:"traceId"`
	PushType        consts.IssueNoticePushType `json:"pushType"` //推送类型
	OrgId           int64                      `json:"orgId"`
	ProjectId       int64                      `json:"projectId"`
	Owner           int64                      `json:"owner"`
}

type DailyProjectReportMsgBo struct {
	ScheduleTraceId string                     `json:"traceId"`
	PushType        consts.IssueNoticePushType `json:"pushType"` //推送类型
	OrgId           int64                      `json:"orgId"`
	Owner           int64                      `json:"owner"`
	//项目id
	ProjectId int64 `json:"projectId"`
	//项目名称
	ProjectName string `json:"projectName"`
	//今日完成
	DailyFinishCount int64 `json:"dailyFinishCount"`
	//剩余未完成
	RemainingCount int64 `json:"dailyRemainingCount"`
	//逾期任务数量
	OverdueCount int64 `json:"dailyOverdueCount"`
}

type RetryPushMsg struct {
	Num       int       `json:"num"`
	OrgId     int64     `json:"orgId"`
	RetryTime time.Time `json:"retryTime"`
}
