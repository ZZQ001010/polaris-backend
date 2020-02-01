package bo

import "time"

type ProcessStepBo struct {
	Id          int64     `db:"id,omitempty" json:"id"`
	OrgId       int64     `db:"org_id,omitempty" json:"orgId"`
	ProcessId   int64     `db:"process_id,omitempty" json:"processId"`
	LangCode    string    `db:"lang_code,omitempty" json:"langCode"`
	Name        string    `db:"name,omitempty" json:"name"`
	StartStatus int64     `db:"start_status,omitempty" json:"startStatus"`
	EndStatus   int64     `db:"end_status,omitempty" json:"endStatus"`
	Sort        int       `db:"sort,omitempty" json:"sort"`
	IsDefault   int       `db:"is_default,omitempty" json:"isDefault"`
	Remark      string    `db:"remark,omitempty" json:"remark"`
	Status      int       `db:"status,omitempty" json:"status"`
	Creator     int64     `db:"creator,omitempty" json:"creator"`
	CreateTime  time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator     int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime  time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version     int       `db:"version,omitempty" json:"version"`
	IsDelete    int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*ProcessStepBo) TableName() string {
	return "ppm_prs_process_step"
}
