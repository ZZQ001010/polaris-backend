package po

import "time"

type PpmPrsProcess struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	LangCode   string    `db:"lang_code,omitempty" json:"langCode"`
	Name       string    `db:"name,omitempty" json:"name"`
	IsDefault  int       `db:"is_default,omitempty" json:"isDefault"`
	Type       int       `db:"type,omitempty" json:"type"`
	Sort       int       `db:"sort,omitempty" json:"sort"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmPrsProcess) TableName() string {
	return "ppm_prs_process"
}
