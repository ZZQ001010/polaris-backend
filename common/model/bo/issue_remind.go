package bo

import (
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
)

type IssueRemindMqBo struct {
	//要处理的任务id
	IssueRemindInfoList []IssueRemindInfoBo `json:"issueRemindInfoList"`
	//traceId
	TraceId string `json:"traceId"`
	//推送类型
	PushType  consts.IssueNoticePushType `json:"pushType"` //推送类型
}

type IssueRemindInfoBo struct {
	//任务id
	Id int64 `json:"id" db:"id"`
	//计划结束时间
	PlanEndTime types.Time `json:"planEndTime" db:"plan_end_time"`
	//负责人id
	OwnerId int64 `json:"ownerId" db:"owner"`
	//组织id
	OrgId int64 `json:"orgId" db:"org_id"`
	//项目id
	ProjectId int64 `json:"projectId" db:"project_id"`
	//标题
	Title string `json:"title" db:"title"`
	//父任务id
	ParentId int64 `json:"parentId" db:"parent_id"`
}