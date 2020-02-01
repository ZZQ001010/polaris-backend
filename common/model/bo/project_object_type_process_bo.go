package bo

import "time"

type PpmPrsProjectObjectTypeProcessBo struct {
	Id                  int64     `db:"id,omitempty" json:"id"`
	OrgId               int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId           int64     `db:"project_id,omitempty" json:"projectId"`
	ProjectObjectTypeId int64     `db:"project_object_type_id,omitempty" json:"projectObjectTypeId"`
	ProcessId           int64     `db:"process_id,omitempty" json:"processId"`
	Sort                int       `db:"sort,omitempty" json:"sort"`
	Creator             int64     `db:"creator,omitempty" json:"creator"`
	CreateTime          time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator             int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime          time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version             int       `db:"version,omitempty" json:"version"`
	IsDelete            int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmPrsProjectObjectTypeProcessBo) TableName() string {
	return "ppm_prs_project_object_type_process"
}
