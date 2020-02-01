package consts

const (
	TrendsIssueExt = "{\"issueType\":\"T\"}"
)

const (
	TrendsModuleOrg     = "Org"
	TrendsModuleProject = "Project"
	TrendsModuleIssue   = "Issue"
	TrendsModuleRole    = "Role"
)

const (
	TrendsOperObjectTypeIssue   = "Issue"
	TrendsOperObjectTypeComment = "Comment"
	TrendsOperObjectTypeProject = "Project"
	TrendsOperObjectTypeRole    = "Role"
	TrendsOperObjectTypeUser = "User"
	TrendsOperObjectTypeOrg = "Org"
)

//动态事件，放于relationType字段，长度限制32（任务）
const (
	//新增任务
	TrendsRelationTypeCreateIssue = "AddIssue"
	//添加评论
	TrendsRelationTypeCreateIssueComment = "AddIssueComment"
	//新增任务（子）
	TrendsRelationTypeCreateChildIssue = "AddChildIssue"
	//更新任务
	TrendsRelationTypeUpdateIssue = "UpdIssue"
	//更新任务状态
	TrendsRelationTypeUpdateIssueStatus = "UpdIssueStatus"
	//删除任务
	TrendsRelationTypeDeleteIssue = "DelIssue"
	//删除任务（子）
	TrendsRelationTypeDeleteChildIssue = "DelChildIssue"
	//添加关注人
	TrendsRelationTypeAddedIssueFollower = "AddIssueFollower"
	//去除关注人
	TrendsRelationTypeDeleteIssueFollower = "DelIssueFollower"
	//添加成员
	TrendsRelationTypeAddedIssueParticipant = "AddIssueParticipant"
	//修改任务负责人
	TrendsRelationTypeUpdateIssueOwner = "UpdateIssueOwner"
	//删除成员
	TrendsRelationTypeDeletedIssueParticipant = "DelIssueParticipant"
	//添加关联任务
	TrendsRelationTypeAddRelationIssue = "AddRelationIssue"
	//删除关联任务
	TrendsRelationTypeDeleteRelationIssue = "DelRelationIssue"
	//上传附件
	TrendsRelationTypeUploadResource = "UploadResource"
	//删除附件
	TrendsRelationTypeDeleteResource = "DeleteResource"
	//变更任务栏
	TrendsRelationTypeUpdateIssueProjectObjectType = "UpdateIssueProjectObjectType"
)

//动态事件，放于relationType字段，长度限制32（项目）
const (
	//创建项目
	TrendsRelationTypeCreateProject = "AddProject"
	//更新项目
	TrendsRelationTypeUpdateProject = "UpdProject"
	//关注项目
	TrendsRelationTypeStarProject = "StarProject"
	//取关项目
	TrendsRelationTypeUnstarProject = "UnstarProject"
	//退出项目
	TrendsRelationTypeUnbindProject = "UnbindProject"
	//更新项目状态
	TrendsRelationTypeUpdateProjectStatus = "UpdProjectStatus"
	//修改项目负责人
	TrendsRelationTypeUpdateProjectOwner = "UpdateProjectOwner"
	//添加关注人
	TrendsRelationTypeAddedProjectFollower = "AddProjectFollower"
	//去除关注人
	TrendsRelationTypeDeleteProjectFollower = "DelProjectFollower"
	//添加成员
	TrendsRelationTypeAddedProjectParticipant = "AddProjectParticipant"
	//删除成员
	TrendsRelationTypeDeletedProjectParticipant = "DelProjectParticipant"
	//批量插入任务
	TrendsRelationTypeCreateIssueBatch = "CreateIssueBatch"
	//上传文件
	TrendsRelationTypeUploadProjectFile = "UploadProjectFile"
	//更新文件
	TrendsRelationTypeUpDateProjectFile = "UpdateProjectFile"
	//删除文件
	TrendsRelationTypeDeleteProjectFile = "DeleteProjectFile"
	//上传文件
	TrendsRelationTypeCreateProjectFolder = "CreateProjectFolder"
	//更新文件
	TrendsRelationTypeUpdateProjectFolder = "UpdateProjectFolder"
	//删除文件
	TrendsRelationTypeDeleteProjectFolder = "DeleteProjectFolder"
)

//组织动态事件
const(
	TrendsRelationTypeApplyJoinOrg = "ApplyJoinOrg"
)

//项目动态查询有效的动态类型
var ValidRelationTypesOfProject = []string{
	TrendsRelationTypeCreateProject,
	TrendsRelationTypeUpdateProject,
	TrendsRelationTypeStarProject,
	TrendsRelationTypeUnstarProject,
	TrendsRelationTypeUnbindProject,
	TrendsRelationTypeUpdateProjectStatus,
	TrendsRelationTypeAddedProjectFollower,
	TrendsRelationTypeDeleteProjectFollower,
	TrendsRelationTypeAddedProjectParticipant,
	TrendsRelationTypeDeletedProjectParticipant,
	TrendsRelationTypeCreateIssue,
	TrendsRelationTypeDeleteIssue,
	TrendsRelationTypeCreateIssueBatch,
	TrendsRelationTypeUpdateProjectOwner,
}


var TrendsOperObjPropertyMap = map[string]string{}

//动态修改字段名称定义
const (
	TrendsOperObjPropertyNameStatus      = "status"
	TrendsOperObjPropertyNameFollower    = "follower"
	TrendsOperObjPropertyNameParticipant = "participant"
	TrendsOperObjPropertyNameOwner       = "owner"
)

const (
	TrendsRelationTypeUser = "User"
)

//别名
func init() {
	TrendsOperObjPropertyMap[TrendsOperObjPropertyNameStatus] = "状态"
}
