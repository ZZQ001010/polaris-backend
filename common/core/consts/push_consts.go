package consts

//推送常量定义
const (
	PushIssueRemind            = "PushIssueRemind"
	PushIssueCommentAndAt      = "PushIssueCommentAndAt"
	PushIssueUpdate            = "PushIssueUpdate"
	PushRelatedContentDynamics = "PushRelatedContentDynamics"
)

type IssueNoticePushType int32

//任务动态推送类型
const (
	//创建任务
	PushTypeCreateIssue IssueNoticePushType = 0
	//更新任务
	PushTypeUpdateIssue IssueNoticePushType = 1
	//更新任务状态
	PushTypeUpdateIssueStatus IssueNoticePushType = 2
	//删除任务
	PushTypeDeleteIssue IssueNoticePushType = 3
	//更新任务成员
	PushTypeUpdateIssueMembers IssueNoticePushType = 4
	//更新任务关联
	PushTypeUpdateRelationIssue IssueNoticePushType = 5
	//每日项目日报推送- 通知项目推送
	PushTypeDailyProjectReportNoticProject IssueNoticePushType = 6
	//每日项目日报推送- 通知消息服务
	PushTypeDailyProjectReportNoticMsg IssueNoticePushType = 7
	//任务提醒
	PushTypeIssueRemind IssueNoticePushType = 8
	//任务评论
	PushTypeIssueComment IssueNoticePushType = 9
	//上传附件
	PushTypeUploadResource IssueNoticePushType = 10
	//更新任务泳道
	PushTypeUpdateIssueProjectObjectType IssueNoticePushType = 11
	//删除附件
	PushTypeDeleteResource IssueNoticePushType = 12
)

//项目动态推送类型
const (
	PushTypeCreateProject        IssueNoticePushType = 1000
	PushTypeUpdateProject        IssueNoticePushType = 1001
	PushTypeUpdateProjectMembers IssueNoticePushType = 1002
	PushTypeStarProject          IssueNoticePushType = 1003
	PushTypeUnstarProject        IssueNoticePushType = 1004
	PushTypeUnbindProject        IssueNoticePushType = 1005
	PushTypeUpdateProjectStatus  IssueNoticePushType = 1006
	PushTypeCreateIssueBatch     IssueNoticePushType = 1007
	PushTypeCreateProjectFile    IssueNoticePushType = 1008
	PushTypeUpdateProjectFile    IssueNoticePushType = 1009
	PushTypeDeleteProjectFile    IssueNoticePushType = 1010
	PushTypeCreateProjectFolder  IssueNoticePushType = 1011
	PushTypeUpdateProjectFolder  IssueNoticePushType = 1012
	PushTypeDeleteProjectFolder  IssueNoticePushType = 1013
)

//组织动态推送类型
const (
	//申请加入组织
	PushTypeApplyJoinOrg	IssueNoticePushType = 2001
	//通过申请
	PushTypeApplicationApproved	IssueNoticePushType = 2002
	//提升为管理员
	PushTypePromotionToOrgManager	IssueNoticePushType = 2003
)

const (
	Code                  = "编号"
	Title                 = "标题"
	Owner                 = "负责人"
	Status                = "状态"
	PlanStartTime         = "计划开始时间"
	PlanEndTime           = "截止时间"
	PlanWorkHour          = "计划工作时间"
	Priority              = "优先级"
	Source                = "来源"
	IssueObjectTypeId     = "类型"
	Remark                = "备注"
	PublicStatus          = "项目公开性"
	ProjectNotice         = "项目公告"
	ProjectResourcePath   = "项目封面"
	ProjectResourceName   = "项目文件名"
	ProjectResourceFolder = "项目文件父级文件夹名"
	ProjectFolderName     = "项目文件夹名"
	ProjectFolderParentId = "项目文件夹父文件夹"
	ProjectObjectType     = "任务栏"
)
