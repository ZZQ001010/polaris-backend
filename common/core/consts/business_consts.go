package consts

//Issue 关联类型
const (
	//负责人
	IssueRelationTypeOwner = 1
	//参与人
	IssueRelationTypeParticipant = 2
	//关注人
	IssueRelationTypeFollower = 3
	//关联任务
	IssueRelationTypeIssue = 4
	//资源
	IssueRelationTypeResource = 5

	//状态
	IssueRelationTypeStatus = 20
	//飞书日历日程
	IssueRelationTypeCalendar = 21
)

var MemberRelationTypeList = []int{IssueRelationTypeOwner, IssueRelationTypeParticipant, IssueRelationTypeFollower}

//项目对象类型LangCode
const (
	ProjectObjectTypeLangCodeIteration = "Project.ObjectType.Iteration"
	ProjectObjectTypeLangCodeBug       = "Project.ObjectType.Bug"
	ProjectObjectTypeLangCodeTestTask  = "Project.ObjectType.TestTask"
	ProjectObjectTypeLangCodeFeature   = "Project.ObjectType.Feature"
	ProjectObjectTypeLangCodeDemand    = "Project.ObjectType.Demand"
	ProjectObjectTypeLangCodeTask      = "Project.ObjectType.Task"
)

//项目类型LangCode
const (
	//普通
	ProjectTypeLangCodeNormalTask = "ProjectType.NormalTask"
	//敏捷
	ProjectTypeLangCodeAgile = "ProjectType.Agile"
)

//流程langCode
const (
	ProcessLangCodeDefaultTask      = "Process.Issue.DefaultTask"      //默认任务流程
	ProcessLangCodeDefaultAgileTask = "Process.Issue.DefaultAgileTask" //默认敏捷项目任务流程
	ProcessLangCodeDefaultBug       = "Process.Issue.DefaultBug"       //默认缺陷流程
	ProcessLangCodeDefaultTestTask  = "Process.Issue.DefaultTestTask"  //默认测试任务流程
	ProcessLangCodeDefaultProject   = "Process.DefaultProject"         //默认项目流程
	ProcessLangCodeDefaultIteration = "Process.DefaultIteration"       //默认迭代流程
	ProcessLangCodeDefaultFeature   = "Process.Issue.DefaultFeature"   //默认特性流程
	ProcessLangCodeDefaultDemand    = "Process.Issue.DefaultDemand"    //默认需求流程
)

//项目对象类型ObjectType
const (
	ProjectObjectTypeIteration = 1
	ProjectObjectTypeTask      = 2
)

//Process Status Type
const (
	//未开始
	ProcessStatusTypeNotStarted = 1
	//进行中
	ProcessStatusTypeProcessing = 2
	//已完成
	ProcessStatusTypeCompleted = 3
)

const (
	//项目流程
	ProcessPrject = 1
	//迭代流程
	ProcessIteration = 2
	//问题流程
	ProcessIssue = 3
)

const (
	//优先级类型-项目
	PriorityTypeProject = 1
	//优先级类型-需求/任务等优先级
	PriorityTypeIssue = 2
)

const (
	//项目状态
	ProcessStatusCategoryProject = 1
	//迭代状态
	ProcessStatusCategoryIteration = 2
	//任务状态
	ProcessStatusCategoryIssue = 3
)

//user team relation Type
const (
	UserTeamRelationTypeLeader = 1
	UserTeamRelationTypeMember = 2
)

const (
	DailyReport   = 1
	WeeklyReport  = 2
	MonthlyReport = 3
)

const TemplateDirPrefix = "resources/template/"

//资源存储方式
const (
	//本地
	LocalResource = 1
	//oss
	OssResource = 2
	//钉盘
	DingDiskResource = 3
	//集群
	PrivateResource = 4
)

var ImgTypeMap = map[string]bool{
	"GIF":  true,
	"JPG":  true,
	"JPEG": true,
	"PNG":  true,
	"TIFF": true,
	"WEBP": true,
}

//url类型
const (
	//http路径
	UrlTypeHttpPath = 1
	//dist路径
	UrlTypeDistPath = 2
)

const (
	//公用项目
	PublicProject = 1
	//私有项目
	PrivateProject = 2
)

const (
	ProjectMemberEffective = 1
	ProjectMemberDisabled  = 2
)

const (
	ProjectIsFiling    = 1
	ProjectIsNotFiling = 2
)

//消息类型定义
const (
	MessageTypeIssueTrends   = 20
	MessageTypeProjectTrends = 21
)

//消息状态
const (
	MessageStatusWait      = 1
	MessageStatusDoing     = 2
	MessageStatusCompleted = 3
	MessageStatusFail      = 4
)

//是否导出飞书日历
const (
	IsSyncOutCalendar    = 1
	IsNotSyncOutCalendar = 2
)

//性别 1:男性 2:女性
const (
	Male   = 1
	Female = 2
)

//联系地址类型
const (
	ContactAddressTypeMobile = 1
	ContactAddressTypeEmail = 2
)

//验证码验证类型1: 登录验证码，2：注册验证码，3：修改密码验证码，4：找回密码验证码，5：绑定，6，解绑
const(
	AuthCodeTypeLogin = 1
	AuthCodeTypeRegister = 2
	AuthCodeTypeResetPwd = 3
	AuthCodeTypeRetrievePwd = 4
	AuthCodeTypeBind = 5
	AuthCodeTypeUnBind = 6
)

//登录类型
const(
	//短信验证码登录
	LoginTypeSMSCode = 1
	//账号密码登录
	LoginTypePwd = 2
	//邮箱登录
	LoginTypeMail = 3
)

//注册类型
const(
	//短信注册
	RegisterTypeSMSCode = 1
	//账号密码注册
	RegisterTypePwd = 2
	//邮箱注册
	RegisterTypeMail = 3
)

//mqtt channelType
const (
	//项目
	MQTTChannelTypeProject = 1
	//组织
	MQTTChannelTypeOrg = 2
	//个人
	MQTTChannelTypeUser = 3
)

//mqtt channel
const (
	MQTTChannelPrefix = "MQTT/"
	MQTTChannelRoot = MQTTChannelPrefix + "#/"
	MQTTChannelOrg = MQTTChannelPrefix + "org/{{.OrgId}}/channel/"
	MQTTChannelProject = MQTTChannelPrefix + "org/{{.OrgId}}/project/{{.ProjectId}}/channel/"
	MQTTChannelUser = MQTTChannelPrefix + "org/{{.OrgId}}/user/{{.UserId}}/channel/"

	MQTTChannelKeyOrg = "OrgId"
	MQTTChannelKeyProject = "ProjectId"
	MQTTChannelKeyUser = "UserId"
)

//mqtt默认ttl(单位：秒)
const MQTTDefaultTTL = 0
//mqtt默认权限
const MQTTDefaultPermissions = "r"
//mqtt默认全局权限
const MQTTDefaultRootPermissions = "rwlsp"

//mqtt 通知类型
const (
	//提醒
	MQTTNoticeTypeRemind = 1
	//数据刷新
	MQTTNoticeTypeDataRefresh = 2
)

//mqtt 数据刷新通知类型
const (
	//新增
	MQTTDataRefreshActionAdd = "ADD"
	//删除
	MQTTDataRefreshActionDel = "DEL"
	//更新
	MQTTDataRefreshActionModify = "MODIFY"
)

//mqtt 数据刷新对象类型
const(
	//项目
	MQTTDataRefreshTypePro = "PRO"
	//任务
	MQTTDataRefreshTypeIssue = "ISSUE"
	//标签
	MQTTDataRefreshTypeTag = "TAG"
	//人员
	MQTTDataRefreshTypeMember = "MEMBER"
)