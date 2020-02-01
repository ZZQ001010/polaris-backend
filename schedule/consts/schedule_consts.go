package consts

import "github.com/galaxy-book/polaris-backend/common/core/consts"

const (
	//项目每日日报是否已发送的缓存key
	ScheduleDailyProjectSendCacheKey = consts.CacheKeyPrefix + consts.SchedulesvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "schedule"
	//用户邀请code有效时间为: 6小时
	CacheDailyProjectSendExpire = 60 * 60 * 12

	//重试retry发送项目每日日报的key
	CacheDailyProjectPlanEndTimeLastScanTime = consts.CacheKeyPrefix + consts.SchedulesvcApplicationName + "daily_project_plan_end_time_last_scan_time"

	//任务截止时间通知-上次扫描的时间
	CacheIssuePlanEndTimeLastScanTime = consts.CacheKeyPrefix + consts.SchedulesvcApplicationName + "issue_plan_end_time_last_scan_time"
)

const (
	IntervalSecondUnit = "s"
	IntervalMinuteUnit = "m"
	IntervalHourUnit   = "h"
)

const (
	Plus  = "+"
	Minus = "-"
)

const DefaultSourceChannel = consts.AppSourceChannelDingTalk

const (
	DailyProjectReportTitle  = "项目日报"
	DailyFinishCountTitle    = "**今日完成: **"
	DailyRemainingCountTitle = "**剩余未完成: **"
	DailyOverdueCountTitle   = "**已逾期: **"
	ViewDetail               = "🔍 查看详情"
)

const (
	DefaultRetryPushNum = 1
	RetryUpperLimit     = 3
)
