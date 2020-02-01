package po

import "time"

type PpmPrsProcessProcessStatus struct {
	Id              int64     `db:"id,omitempty" json:"id"`
	OrgId           int64     `db:"org_id,omitempty" json:"orgId"`
	ProcessId       int64     `db:"process_id,omitempty" json:"processId"`
	ProcessStatusId int64     `db:"process_status_id,omitempty" json:"processStatusId"`
	IsInitStatus    int       `db:"is_init_status,omitempty" json:"isInitStatus"`
	Sort            int       `db:"sort,omitempty" json:"sort"`
	Creator         int64     `db:"creator,omitempty" json:"creator"`
	CreateTime      time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator         int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime      time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version         int       `db:"version,omitempty" json:"version"`
	IsDelete        int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmPrsProcessProcessStatus) TableName() string {
	return "ppm_prs_process_process_status"
}
