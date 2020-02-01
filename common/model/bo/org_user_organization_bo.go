package bo

import "time"

type PpmOrgUserOrganizationBo struct {
	Id               int64     `db:"id,omitempty" json:"id"`
	OrgId            int64     `db:"org_id,omitempty" json:"orgId"`
	UserId           int64     `db:"user_id,omitempty" json:"userId"`
	CheckStatus      int       `db:"check_status,omitempty" json:"checkStatus"`
	UseStatus        int       `db:"use_status,omitempty" json:"useStatus"`
	Status           int       `db:"status,omitempty" json:"status"`
	StatusChangeTime time.Time `db:"status_change_time,omitempty" json:"statusChangeTime"`
	AuditorId        int64     `db:"auditor_id,omitempty" json:"auditorId"`
	AuditTime        time.Time `db:"audit_time,omitempty" json:"auditTime"`
	Creator          int64     `db:"creator,omitempty" json:"creator"`
	CreateTime       time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator          int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime       time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version          int       `db:"version,omitempty" json:"version"`
	IsDelete         int       `db:"is_delete,omitempty" json:"isDelete"`
}

type OrgMemberChangeBo struct {
	//变动类型：1 禁用，2 启用，3 新加入组织(正常)，4 从组织移除,5:新加入组织（禁用）
	ChangeType int `json:"changeType"`
	//组织id
	OrgId int64 `json:"orgId"`
	//变动人员id
	UserId int64 `json:"userId"`
	//变动OpenIds
	OpenId string `json:"openId"`
	//来源
	SourceChannel string `json:"sourceChannel"`
}