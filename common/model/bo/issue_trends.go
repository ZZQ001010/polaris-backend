package bo

import (
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"time"
)

type IssueTrendsBo struct {
	PushType      consts.IssueNoticePushType `json:"pushType"`      //推送类型
	OrgId         int64                      `json:"orgId"`         //组织id
	OperatorId    int64                      `json:"operatorId"`    //操作人id
	IssueId       int64                      `json:"issueId"`       //任务id
	ParentIssueId int64                      `json:"parentIssueId"` //父任务id
	ProjectId     int64                      `json:"projectId"`     //项目id
	PriorityId    int64                      `json:"priorityId"`    //优先级id
	ParentId      int64                      `json:"parentId"`      //父任务id

	IssueTitle         string      `json:"issueTitle"`         //更新后的任务标题
	IssueRemark        string      `json:"issueRemark"`        //任务描述
	IssueStatusId      int64       `json:"issueStatusId"`      //更新后的任务状态
	IssuePlanStartTime *types.Time `json:"issuePlanStartTime"` //更新后的任务开始时间
	IssuePlanEndTime   *types.Time `json:"issuePlanEndTime"`   //更新后的任务结束时间

	SourceChannel string `json:"sourceChannel"` //来源通道

	BeforeOwner              int64   `json:"beforeOwner"`              //之前的负责人
	AfterOwner               int64   `json:"afterOwner"`               //新的负责人
	BeforeChangeFollowers    []int64 `json:"beforeChangeFollowers"`    //之前的关注人
	AfterChangeFollowers     []int64 `json:"afterChangeFollowers"`     //之后的负责人
	BeforeChangeParticipants []int64 `json:"beforeChangeParticipants"` //之前的参与人
	AfterChangeParticipants  []int64 `json:"afterChangeParticipants"`  //之后的参与人
	BindIssues               []int64 `json:"bindIssues"`               //关联任务
	UnbindIssues             []int64 `json:"unbindIssues"`             //取消关联的任务

	OnlyNotice bool `json:"onlyNotice"` //如果为true，表示只处理通知

	OperateObjProperty string           `json:"operateObjProperty"` //操作属性
	NewValue           string           `json:"newValue"`           //新值
	OldValue           string           `json:"oldValue"`           //老值
	Ext                TrendExtensionBo `json:"ext"`
	OperateTime        time.Time        `json:"operateTime"` //操作时间
}

type TrendExtensionBo struct {
	IssueType     string              `json:"issueType"`
	ObjName       string              `json:"objName"`
	ChangeList    []TrendChangeListBo `json:"changeList"`
	MemberInfo    []SimpleUserInfoBo  `json:"memberInfo"`
	RelationIssue RelationIssue       `json:"relationIssue"`

	MentionedUserIds []int64          `json:"mentionedUserIds"` //提及的用户列表
	CommentBo        CommentBo        `json:"commentBo"`
	ResourceInfo     []ResourceInfoBo `json:"resourceInfo"`
}

type ResourceInfoBo struct {
	Url        string    `json:"url"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	UploadTime time.Time `json:"uploadTime"`
	Suffix     string    `json:"suffix"`
}

type TrendChangeListBo struct {
	Field     string `json:"field"`
	FieldName string `json:"fieldName"`
	OldValue  string `json:"oldValue"`
	NewValue  string `json:"newValue"`
}

type RelationIssue struct {
	// 关联信息id
	ID int64 `json:"id"`
	// 关联信息名称
	Title string `json:"title"`
}
