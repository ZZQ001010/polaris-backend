package bo

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/library/db/mysql"
	"time"
)

var log = logger.GetDefaultLogger()

type IssueRelationBo struct {
	Id           int64      `db:"id,omitempty" json:"id"`
	OrgId        int64      `db:"org_id,omitempty" json:"orgId"`
	ProjectId    int64      `db:"project_id,omitempty" json:"projectId"`
	IssueId      int64      `db:"issue_id,omitempty" json:"issueId"`
	RelationId   int64      `db:"relation_id,omitempty" json:"relationId"`
	RelationCode string     `db:"relation_code,omitempty" json:"relationCode"`
	RelationType int        `db:"relation_type,omitempty" json:"relationType"`
	Creator      int64      `db:"creator,omitempty" json:"creator"`
	CreateTime   types.Time `db:"create_time,omitempty" json:"createTime"`
	Updator      int64      `db:"updator,omitempty" json:"updator"`
	UpdateTime   types.Time `db:"update_time,omitempty" json:"updateTime"`
	Version      int        `db:"version,omitempty" json:"version"`
	IsDelete     int        `db:"is_delete,omitempty" json:"isDelete"`
}

type IssueUserBo struct {
	IssueRelationBo
}

type IssueBo struct {
	Id                  int64      `json:"id"`
	OrgId               int64      `json:"orgId"`
	Code                string     `json:"code"`
	ProjectId           int64      `json:"projectId"`
	ProjectObjectTypeId int64      `json:"projectObjectTypeId"`
	Title               string     `json:"title"`
	Owner               int64      `json:"owner"`
	OwnerChangeTime     types.Time `json:"ownerChangeTime"`
	PriorityId          int64      `json:"priorityId"`
	SourceId            int64      `json:"sourceId"`
	IssueObjectTypeId   int64      `json:"issueObjectTypeId"`
	PlanStartTime       types.Time `json:"planStartTime"`
	PlanEndTime         types.Time `json:"planEndTime"`
	StartTime           types.Time `json:"startTime"`
	EndTime             types.Time `json:"endTime"`
	PlanWorkHour        int        `json:"planWorkHour"`
	IterationId         int64      `json:"iterationId"`
	VersionId           int64      `json:"versionId"`
	ModuleId            int64      `json:"moduleId"`
	ParentId            int64      `json:"parentId"`
	Status              int64      `json:"status"`
	Creator             int64      `json:"creator"`
	Sort                int64      `json:"sort"`
	CreateTime          types.Time `json:"createTime"`
	Updator             int64      `json:"updator"`
	UpdateTime          types.Time `json:"updateTime"`
	Version             int        `json:"version"`

	IssueDetailBo `json:"-"`

	OwnerInfo        *IssueUserBo    `json:"ownerInfo"`
	ParticipantInfos *[]IssueUserBo  `json:"participantInfos"`
	FollowerInfos    *[]IssueUserBo  `json:"followerInfos"`
	Tags             []IssueTagReqBo `json:"tags"`
}

//任务和明细联合的Bo
type IssueAndDetailUnionBo struct {
	IssueId                  int64     `db:"issueId,omitempty"`
	IssueStatusId            int64     `db:"issueStatusId,omitempty"`
	IssueProjectObjectTypeId int64     `db:"issueProjectObjectTypeId,omitempty"`
	StoryPoint               int       `db:"storyPoint,omitempty"`
	PlanEndTime              time.Time `db:"planEndTime,omitempty"`
	EndTime                  time.Time `db:"endTime,omitempty"`
	OwnerChangeTime          time.Time `db:"ownerChangeTime,omitempty"`
	CreateTime 				 time.Time `db:"createTime,omitempty"`
}

type IssueUpdateBo struct {
	IssueBo    IssueBo
	NewIssueBo IssueBo

	IssueUpdateCond   mysql.Upd
	IssueDetailRemark *string
	OwnerId           *int64
	UpdateParticipant bool
	UpdateFollower    bool
	Participants      []int64
	Followers         []int64

	OperatorId int64
}

type IssueDailyNoticeBo struct {
	//所有未完成的本人负责的任务数量
	PendingSum uint64
	//截止时间为今天的本人负责的任务数量
	DueOfTodaySum uint64
	//已逾期的本人负责的任务数量
	OverdueSum uint64
	//即将逾期的本人负责的任务数量
	BeAboutToOverdueSum uint64
}

func BuildIssueRelationBosFromUserBos(bos *[]IssueUserBo) *[]IssueRelationBo {
	if bos == nil {
		return &([]IssueRelationBo{})
	}
	relationBos := make([]IssueRelationBo, 0, len(*bos))
	for _, v := range *bos {
		relationBos = append(relationBos, v.IssueRelationBo)
	}
	return &relationBos
}

//任务验证bo
type IssueAuthBo struct {
	Id           int64
	Owner        int64
	ProjectId    int64
	Creator      int64
	Status       int64
	Participants []int64
	Followers    []int64
}

type HomeIssueOwnerInfoBo struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar"`
	// 是否已被删除，为true则代表被组织移除
	IsDeleted bool `json:"isDeleted"`
	// 是否已被禁用, 为true则代表被组织禁用
	IsDisabled bool `json:"isDisabled"`
}

type HomeIssuePriorityInfoBo struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	BgStyle   string `json:"bgStyle"`
	FontStyle string `json:"fontStyle"`
}

type HomeIssueProjectInfoBo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	IsFilling int `json:"isFilling"`
}

type HomeIssueStatusInfoBo struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	DisplayName *string `json:"displayName"`
	BgStyle     string  `json:"bgStyle"`
	FontStyle   string  `json:"fontStyle"`
	Type        int     `json:"type"`
}

type HomeIssueInfoBo struct {
	Issue             IssueBo                  `json:"issue"`
	Project           *HomeIssueProjectInfoBo  `json:"project"`
	Owner             *HomeIssueOwnerInfoBo    `json:"owner"`
	Status            *HomeIssueStatusInfoBo   `json:"status"`
	Priority          *HomeIssuePriorityInfoBo `json:"priority"`
	Tags              []HomeIssueTagInfoBo     `json:"tags"`
	ChildsNum         int64                    `json:"childsNum"`
	ChildsFinishedNum int64                    `json:"childsFinishedNum"`
}

type HomeIssueTagInfoBo struct {
	// 标签id
	ID int64 `json:"id"`
	// 标签名
	Name string `json:"name"`
	// 背景颜色
	BgStyle string `json:"bgStyle"`
	// 字体颜色
	FontStyle string `json:"fontStyle"`
}

type UserIDInfoBo struct {
	//用户id
	UserID int64 `json:"userId"`
	//姓名
	Name string `json:"name"`
	//头像
	Avatar string `json:"avatar"`
	//外部用户id
	EmplID string `json:"emplId"`
	//唯一id(弃用，统一emplId，也就是外部用户id)
	UnionID string `json:"unionId"`
	// 是否已被删除，为true则代表被组织移除
	IsDeleted bool `json:"isDeleted"`
	// 是否已被禁用, 为true则代表被组织禁用
	IsDisabled bool `json:"isDisabled"`
}

type DepartmentMemberInfoBo struct {
	// 用户id
	UserID int64 `json:"userId"`
	// 用户名称
	Name string `json:"name"`
	//姓名拼音
	NamePy string `json:"namePy"`
	// 用户头像
	Avatar string `json:"avatar"`
	// 工号：企业下唯一
	EmplID string `json:"emplId"`
	// unionId： 开发者账号下唯一
	UnionID string `json:"unionId"`
	// 用户部门id
	DepartmentID int64 `json:"departmentId"`
}

type IssueInfoBo struct {
	Issue    IssueBo                  `json:"issue"`
	Project  *HomeIssueProjectInfoBo  `json:"project"`
	Owner    *HomeIssueOwnerInfoBo    `json:"owner"`
	Status   *HomeIssueStatusInfoBo   `json:"status"`
	Priority *HomeIssuePriorityInfoBo `json:"priority"`

	ParticipantInfos []*UserIDInfoBo          `json:"participantInfos"`
	FollowerInfos    []*UserIDInfoBo          `json:"followerInfos"`
	NextStatus       []*HomeIssueStatusInfoBo `json:"nextStatus"`
}

type IssueReportDetailBo struct {
	Total          int64
	ReportUserName string
	StartTime      string
	EndTime        string
	List           []HomeIssueInfoBo
	ShareID        string
}

type TaskStatBo struct {
	NotStart   int64 `json:"notStart"`
	Processing int64 `json:"processing"`
	Completed  int64 `json:"completed"`
}

type IssueBoListCond struct {
	OrgId        int64
	ProjectId    *int64
	IterationId  *int64
	RelationType *int
	UserId       *int64
	Ids          []int64
}

type IssueChildCountBo struct {
	ParentIssueId int64 `db:"parentIssueId"`
	Count         int64 `db:"count"`
}

type SelectIssueIdsCondBo struct {
	//计划开始时间开始范围
	BeforePlanEndTime *string `json:"beforePlanEndTime"`
	//计划开始时间结束范围
	AfterPlanEndTime *string `json:"afterPlanEndTime"`
}

type IssueIdBo struct {
	Id int64 `json:"id" db:"id"`
}

type IssueAndDetailInfoBo struct {
	Id                  int64      `db:"id,omitempty" json:"id"`
	ProjectObjectTypeId int64      `db:"project_object_type_id,omitempty" json:"projectObjectTypeId"`
	Title               string     `db:"title,omitempty" json:"title"`
	PriorityId          int64      `db:"priority_id,omitempty" json:"priorityId"`
	Owner               int64      `db:"owner,omitempty" json:"owner"`
	PlanStartTime       types.Time `db:"plan_start_time,omitempty" json:"planStartTime"`
	PlanEndTime         types.Time `db:"plan_end_time,omitempty" json:"planEndTime"`
	ParentId            int64      `db:"parent_id,omitempty" json:"parentId"`
	Status              int64      `db:"status,omitempty" json:"status"`
	Creator             int64      `db:"creator,omitempty" json:"creator"`
	CreateTime          types.Time `db:"create_time,omitempty" json:"createTime"`
	Remark              string     `db:"remark,omitempty" json:"remark"`
}

type IssueCombineBo struct {
	IssueAndDetailInfoBo
	Children []IssueAndDetailInfoBo `json:"children"`
}

type IssueMembersBo struct {
	MemberIds []int64 `json:"memberIds"`
	OwnerId int64 `json:"ownerId"`
	ParticipantIds []int64 `json:"participantIds"`
	FollowerIds []int64 `json:"followerIds"`
}