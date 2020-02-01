package bo

import "time"

type ProjectObjectTypeBo struct {
	Id                               int64                            `db:"id,omitempty" json:"id"`
	OrgId                            int64                            `db:"org_id,omitempty" json:"orgId"`
	LangCode                         string                           `db:"lang_code,omitempty" json:"langCode"`
	PreCode                          string                           `db:"pre_code,omitempty" json:"preCode"`
	Name                             string                           `db:"name,omitempty" json:"name"`
	ObjectType                       int                              `db:"object_type,omitempty" json:"objectType"`
	BgStyle                          string                           `db:"bg_style,omitempty" json:"bgStyle"`
	FontStyle                        string                           `db:"font_style,omitempty" json:"fontStyle"`
	Icon                             string                           `db:"icon,omitempty" json:"icon"`
	Sort                             int                              `db:"sort,omitempty" json:"sort"`
	OrginalSort                      int                              `db:"sort,omitempty" json:"orginalSort"`
	Remark                           string                           `db:"remark,omitempty" json:"remark"`
	IsReadonly                       int                              `db:"is_readonly,omitempty" json:"isReadonly"`
	Status                           int                              `db:"status,omitempty" json:"status"`
	Creator                          int64                            `db:"creator,omitempty" json:"creator"`
	CreateTime                       time.Time                        `db:"create_time,omitempty" json:"createTime"`
	Updator                          int64                            `db:"updator,omitempty" json:"updator"`
	UpdateTime                       time.Time                        `db:"update_time,omitempty" json:"updateTime"`
	Version                          int                              `db:"version,omitempty" json:"version"`
	IsDelete                         int                              `db:"is_delete,omitempty" json:"isDelete"`
	PpmPrsProjectObjectTypeProcessBo PpmPrsProjectObjectTypeProcessBo `json:"-"`
	AfterID                          *int64                           `json:"afterId"`
	BeforeID                         *int64                           `json:"beforeId"`
}

type ProjectObjectTypeBoSorter []ProjectObjectTypeBo

func (ms ProjectObjectTypeBoSorter) Len() int {
	return len(ms)
}
func (ms ProjectObjectTypeBoSorter) Less(i, j int) bool {
	return ms[i].Sort < ms[j].Sort // 按值排序
}
func (ms ProjectObjectTypeBoSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (*ProjectObjectTypeBo) TableName() string {
	return "ppm_prs_project_object_type"
}

type ProjectObjectTypeRestInfoBo struct {
	// 主键
	ID int64 `json:"id"`
	// 语言编号
	LangCode string `json:"langCode"`
	// 名称
	Name string `json:"name"`
	// 类型,1迭代，2问题
	ObjectType int `json:"objectType"`
}

type ProjectSupportObjectTypeListBo struct {
	// 项目支持的对象类型
	ProjectSupportList []*ProjectObjectTypeRestInfoBo `json:"projectSupportList"`
	// 迭代支持的对象类型
	IterationSupportList []*ProjectObjectTypeRestInfoBo `json:"iterationSupportList"`
}

type ProjectObjectTypeAndProjectIdBo struct {
	ProjectObjectTypeId int64 `json:"projectObjectTypeId"`
	ProjectId           int64 `json:"projectId"`
}
