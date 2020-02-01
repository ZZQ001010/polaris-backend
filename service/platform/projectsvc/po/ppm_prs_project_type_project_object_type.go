package po

import "time"

type PpmPrsProjectTypeProjectObjectType struct {
	Id                  int64     `db:"id,omitempty" json:"id"`
	OrgId               int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectTypeId       int64     `db:"project_type_id,omitempty" json:"projectTypeId"`
	ProjectObjectTypeId int64     `db:"project_object_type_id,omitempty" json:"projectObjectTypeId"`
	Remark              string    `db:"remark,omitempty" json:"remark"`
	DefaultProcessId    int64     `db:"default_process_id,omitempty" json:"defaultProcessId"`
	IsReadonly          int       `db:"is_readonly,omitempty" json:"isReadonly"`
	Status              int       `db:"status,omitempty" json:"status"`
	Creator             int64     `db:"creator,omitempty" json:"creator"`
	CreateTime          time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator             int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime          time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version             int       `db:"version,omitempty" json:"version"`
	IsDelete            int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmPrsProjectTypeProjectObjectType) TableName() string {
	return "ppm_prs_project_type_project_object_type"
}
