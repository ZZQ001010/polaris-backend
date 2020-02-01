package bo

//缓存用户登录信息
type CacheUserInfoBo struct {
	OutUserId     string `json:"outUserId"`
	SourceChannel string `json:"sourceChannel"`
	UserId        int64  `json:"userId"`
	CorpId        string `json:"corpId"`
	OrgId         int64  `json:"orgId"`
}

type BaseUserInfoBo struct {
	UserId        int64  `json:"userId"`
	OutUserId     string `json:"outUserId"` //有可能为空
	OrgId         int64  `json:"orgId"`
	OutOrgId      string `json:"outOrgId"` //有可能为空
	Name          string `json:"name"`
	NamePy		  string `json:"namePy"`
	Avatar        string `json:"avatar"`
	HasOutInfo    bool   `json:"hasOutInfo"`
	HasOrgOutInfo bool   `json:"hasOrgOutInfo"`

	OrgUserIsDelete    int `json:"orgUserIsDelete"`           //是否被组织移除
	OrgUserStatus      int `json:"orgUserStatus"`      //用户组织状态
	OrgUserCheckStatus int `json:"orgUserCheckStatus"` //用户组织审核状态
}

//用户基本信息扩展
type BaseUserInfoExtBo struct {
	BaseUserInfoBo

	//部门id
	DepartmentId int64 `json:"departmentId"`
}

type BaseUserOutInfoBo struct {
	UserId    int64  `json:"userId"`
	OutUserId string `json:"outUserId"` //有可能为空
	OutOrgId  string `json:"outOrgId"`  //有可能为空
	OrgId     int64  `json:"orgId"`
}

type BaseOrgInfoBo struct {
	OrgId         int64  `json:"orgId"`
	OrgName       string `json:"orgName"`
	OrgOwnerId      int64 `json:"orgOwnerId"`
	OutOrgId      string `json:"outOrgId"`
	SourceChannel string `json:"sourceChannel"`
}

type BaseOrgOutInfoBo struct {
	OrgId         int64  `json:"orgId"`
	OutOrgId      string `json:"outOrgId"`
	SourceChannel string `json:"sourceChannel"`
}

type CacheProcessStatusBo struct {
	StatusId    int64  `json:"statusId"`
	StatusType  int    `json:"statusType"`
	Category    int    `json:"category"`
	IsInit      bool   `json:"isInit"`
	BgStyle     string `json:"bgStyle"`
	FontStyle   string `json:"fontStyle"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type CacheProjectCalendarInfoBo struct {
	IsSyncOutCalendar int    `json:"isSyncOutCalendar"`
	CalendarId        string `json:"calendarId"`
}
