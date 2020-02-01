package bo

import (
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
)

type ProjectMemberChangeBo struct {
	PushType            consts.IssueNoticePushType //推送类型
	OrgId               int64
	ProjectId           int64
	OperatorId          int64
	BeforeChangeMembers []int64
	AfterChangeMembers  []int64

	OperateObjProperty string
	NewValue           string
	OldValue           string
}

type ProjectAuthBo struct {
	Id           int64   `json:"id"`
	Name         string  `json:"name"`
	Creator      int64   `json:"creator"`
	Owner        int64   `json:"owner"`
	Status       int64   `json:"status"`
	PublicStatus int     `json:"publicStatus"` //公共状态
	IsFilling    int     `json:"isFilling"`
	Participants []int64 `json:"participants"`
	Followers    []int64 `json:"followers"`
	ProjectType  int64   `json:"projectType"`
}

type ProjectBo struct {
	Id            int64      `db:"id,omitempty" json:"id"`
	OrgId         int64      `db:"org_id,omitempty" json:"orgId"`
	Code          string     `db:"code,omitempty" json:"code"`
	Name          string     `db:"name,omitempty" json:"name"`
	PreCode       string     `db:"pre_code,omitempty" json:"preCode"`
	Owner         int64      `db:"owner,omitempty" json:"owner"`
	ProjectTypeId int64      `db:"project_type_id,omitempty" json:"projectTypeId"`
	PriorityId    int64      `db:"priority_id,omitempty" json:"priorityId"`
	PlanStartTime types.Time `db:"plan_start_time,omitempty" json:"planStartTime"`
	PlanEndTime   types.Time `db:"plan_end_time,omitempty" json:"planEndTime"`
	PublicStatus  int        `db:"public_status,omitempty" json:"publicStatus"`
	ResourceId    int64      `db:"resource_id,omitempty" json:"resourceId"`
	IsFiling      int        `db:"is_filing,omitempty" json:"isFiling"`
	Remark        string     `db:"remark,omitempty" json:"remark"`
	Status        int64      `db:"status,omitempty" json:"status"`
	Creator       int64      `db:"creator,omitempty" json:"creator"`
	CreateTime    types.Time `db:"create_time,omitempty" json:"createTime"`
	Updator       int64      `db:"updator,omitempty" json:"updator"`
	UpdateTime    types.Time `db:"update_time,omitempty" json:"updateTime"`
	Version       int        `db:"version,omitempty" json:"version"`
	IsDelete      int        `db:"is_delete,omitempty" json:"isDelete"`
}

type UpdateProjectBo struct {
	ID            int64       `json:"id"`
	Code          *string     `json:"code"`
	Name          *string     `json:"name"`
	PreCode       *string     `json:"preCode"`
	Owner         *int64      `json:"owner"`
	PriorityID    *int64      `json:"priorityId"`
	PlanStartTime *types.Time `json:"planStartTime"`
	PlanEndTime   *types.Time `json:"planEndTime"`
	PublicStatus  *int        `json:"publicStatus"`
	ResourceID    *int64      `json:"resourceId"`
	IsFiling      *int        `json:"isFiling"`
	Remark        *string     `json:"remark"`
	Status        *int64      `json:"status"`
	ResourcePath  *string     `json:"resourcePath"`
	ResourceType  *int        `json:"resourceType"`
	MemberIds     []int64     `json:"memberIds"`
	FollowerIds   []int64     `json:"followerIds"`
	UpdateFields  []string    `json:"updateFields"`
}

type RelationInfoTypeBo struct {
	UserId       int64  `db:"user_id,omitempty" json:"userId"`
	RelationType int    `db:"relation_type,omitempty" json:"relationType"`
	RelationId   int64  `db:"relation_id,omitempty" json:"relationId"`
	ProjectId    int64  `db:"project_id,omitempty" json:"projectId"`
	OutOrgUserId string `db:"out_org_user_id,omitempty" json:"outOrgUserId"`
	OutUserId    string `db:"out_user_id,omitempty" json:"outUserId"`
	Name         string `db:"name,omitempty" json:"name"`
	Avatar       string `db:"avatar,omitempty" json:"avatar"`
}

type ProjectStatBo struct {
	IterationTotal int64 `json:"iterationTotal"`
	TaskTotal      int64 `json:"taskTotal"`
	MemberTotal    int64 `json:"memberTotal"`
}

type InsertProjectMembersInputBo struct {
	OrgId int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
	OperatorId int64 `json:"operatorId"`

	OwnerInfo        *IssueUserBo   `json:"ownerInfo"`
	ParticipantInfos []IssueUserBo `json:"participantInfos"`
	FollowerInfos    []IssueUserBo `json:"followerInfos"`
}

type IssueRelationUserInfosBo struct {
	OwnerInfo        *IssueUserBo   `json:"ownerInfo"`
	ParticipantInfos []IssueUserBo `json:"participantInfos"`
	FollowerInfos    []IssueUserBo `json:"followerInfos"`
}