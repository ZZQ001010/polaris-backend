package po

import "time"

type PpmOrgDepartmentOutInfo struct {
	Id                       int64     `db:"id,omitempty" json:"id"`
	OrgId                    int64     `db:"org_id,omitempty" json:"orgId"`
	DepartmentId             int64     `db:"department_id,omitempty" json:"departmentId"`
	SourcePlatform           string    `db:"source_platform,omitempty" json:"sourcePlatform"`
	SourceChannel            string    `db:"source_channel,omitempty" json:"sourceChannel"`
	OutOrgDepartmentId       string    `db:"out_org_department_id,omitempty" json:"outOrgDepartmentId"`
	OutOrgDepartmentCode     string    `db:"out_org_department_code,omitempty" json:"outOrgDepartmentCode"`
	Name                     string    `db:"name,omitempty" json:"name"`
	OutOrgDepartmentParentId string    `db:"out_org_department_parent_id,omitempty" json:"outOrgDepartmentParentId"`
	Status                   int       `db:"status,omitempty" json:"status"`
	Creator                  int64     `db:"creator,omitempty" json:"creator"`
	CreateTime               time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator                  int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime               time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version                  int       `db:"version,omitempty" json:"version"`
	IsDelete                 int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrgDepartmentOutInfo) TableName() string {
	return "ppm_org_department_out_info"
}
