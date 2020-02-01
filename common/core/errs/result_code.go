package errs

import "github.com/galaxy-book/common/core/errors"

var (
	//成功
	OK = errors.OK

	//token错误
	RequestError = errors.RequestError

	//认证错误
	Unauthorized = errors.Unauthorized

	//禁止访问
	ForbiddenAccess = errors.ForbiddenAccess

	//请求地址不存在
	PathNotFound = errors.PathNotFound

	//不支持该方法
	MethodNotAllowed = errors.MethodNotAllowed

	//Token过期
	TokenExpires = errors.TokenExpires

	//请求参数错误
	ServerError = errors.ServerError

	//过载保护,服务暂不可用
	ServiceUnavailable = errors.ServiceUnavailable

	//服务调用超时
	Deadline = errors.Deadline

	//超出限制
	LimitExceed = errors.LimitExceed

	//参数错误
	ParamError = errors.ParamError

	//文件过大
	FileTooLarge = errors.FileTooLarge

	//文件类型错误
	FileTypeError = errors.FileTypeError

	//文件或目录不存在
	FileNotExist = errors.FileNotExist

	//文件路径为空
	FilePathIsNull = errors.FilePathIsNull

	//读取文件失败
	FileReadFail = errors.FileReadFail

	//错误未定义
	ErrorUndefined = errors.ErrorUndefined

	//业务失败
	BusinessFail = errors.BusinessFail

	//系统异常
	SystemError = errors.SystemError

	//未知错误
	UnknownError = errors.UnknownError

	//业务中详细异常定义
	//>>>工具异常
	MysqlOperateError          = errors.MysqlOperateError
	RedisOperateError          = errors.RedisOperateError
	GetDistributedLockError    = errors.GetDistributedLockError
	OssError                   = errors.AddResultCodeInfo(300103, "Oss异常", "ResultCode.OssError")
	RocketMQProduceInitError   = errors.RocketMQProduceInitError
	RocketMQSendMsgError       = errors.RocketMQSendMsgError
	RocketMQConsumerInitError  = errors.RocketMQConsumerInitError
	RocketMQConsumerStartError = errors.RocketMQConsumerStartError
	RocketMQConsumerStopError  = errors.RocketMQConsumerStopError

	DbMQSendMsgError         = errors.DbMQSendMsgError
	DbMQCreateConsumerError  = errors.DbMQCreateConsumerError
	DbMQConsumerStartedError = errors.DbMQConsumerStartedError

	KafkaMqSendMsgError           = errors.KafkaMqSendMsgError
	KafkaMqSendMsgCantBeNullError = errors.KafkaMqSendMsgCantBeNullError
	KafkaMqConsumeMsgError        = errors.KafkaMqConsumeMsgError
	KafkaMqConsumeStartError      = errors.KafkaMqConsumeStartError

	RunModeUnsupportUpload = errors.AddResultCodeInfo(300104, "该部署模式不支持本地上传", "ResultCode.RunModeUnsupportUpload")

	SystemBusy = errors.AddResultCodeInfo(300105, "系统繁忙，请稍后重试", "ResultCode.SystemBusy")

	JSONConvertError = errors.AddResultCodeInfo(300201, "Json转换出现异常", "ResultCode.JSONConvertError")
	ObjectCopyError  = errors.AddResultCodeInfo(300202, "对象copy出现异常", "ResultCode.ObjectCopyError")
	CacheProxyError  = errors.AddResultCodeInfo(300203, "缓存代理出现异常", "ResultCode.CacheProxyError")
	ObjectTypeError  = errors.AddResultCodeInfo(300204, "对象类型错误", "ResultCode.ObjectTypeError")

	ApplyIdError        = errors.AddResultCodeInfo(300205, "ID申请异常", "ResultCode.ApplyIdError")
	ApplyIdCountTooMany = errors.AddResultCodeInfo(300206, "申请id数量过多", "ResultCode.ApplyIdCountTooMany")
	TypeConvertError    = errors.AddResultCodeInfo(300207, "类型转换出现异常", "ResultCode.TypeConvertError")
	UpdateFiledIsEmpty  = errors.AddResultCodeInfo(300208, "未更新任何信息", "ResultCode.UpdateFiledIsEmpty")

	TokenAuthError      = errors.AddResultCodeInfo(300301, "身份认证异常，请重新登录", "ResultCode.TokenAuthError")
	TokenNotExist       = errors.AddResultCodeInfo(300302, "身份认证失败，请重新登录", "ResultCode.TokenNotExist")
	SuiteTicketError    = errors.AddResultCodeInfo(300303, "获取SuiteTicket异常", "ResultCode.SuiteTicketError")
	GetContextError     = errors.AddResultCodeInfo(300304, "获取请求上下文异常", "ResultCode.GetContextError")
	TemplateRenderError = errors.AddResultCodeInfo(300305, "模板解析失败", "ResultCode.TemplateRenderError")
	DecryptError = errors.AddResultCodeInfo(300401, "参数解密异常", "ResultCode.DecryptError")
	CaptchaError = errors.AddResultCodeInfo(300402, "验证码错误", "ResultCode.CaptchaError")

	MQTTKeyGenError = errors.AddResultCodeInfo(300501, "生成key发生异常", "ResultCode.MQTTKeyGenError")
	MQTTPublishError = errors.AddResultCodeInfo(300502, "MQTT推送消息发生异常", "ResultCode.MQTTPublishError")
	MQTTConnectError = errors.AddResultCodeInfo(300503, "MQTT连接时发生异常", "ResultCode.MQTTConnectError")
	MQTTMissingConfigError = errors.AddResultCodeInfo(300504, "MQTT缺少配置", "ResultCode.MQTTMissingConfigError")


	TryDistributedLockError = errors.TryDistributedLockError

	//>>>业务异常
	InitDbFail                     = errors.AddResultCodeInfo(400000, "初始化db失败", "ResultCode.InitDbFail")
	ObjectRecordNotFoundError      = errors.AddResultCodeInfo(400001, "对象记录不存在", "ResultCode.ObjectRecordNotFoundError")
	DingTalkUserInfoNotInitedError = errors.AddResultCodeInfo(400002, "钉钉用户没有初始化", "ResultCode.DingTalkUserInfoNotInitedError")
	UserNotFoundError              = errors.AddResultCodeInfo(400003, "用户信息不存在或已经删除", "ResultCode.UserNotFoundError")
	CacheUserInfoNotExistError     = errors.AddResultCodeInfo(400004, "令牌对应的用户信息不存在", "ResultCode.CacheUserInfoNotExistError")

	PageSizeOverflowMaxSizeError = errors.AddResultCodeInfo(400005, "请求页长超出最大页长限制", "ResultCode.PageSizeOverflowMaxSizeError")
	OutOfConditionError          = errors.AddResultCodeInfo(400006, "请求条件超出限制", "ResultCode.OutOfConditionError")
	ConditionHandleError         = errors.AddResultCodeInfo(400007, "条件处理异常", "ResultCode.ConditionHandleError")

	ReqParamsValidateError = errors.AddResultCodeInfo(400008, "请求参数校验异常", "ResultCode.ReqParamsValidateError")

	OrgNotInitError        = errors.AddResultCodeInfo(400009, "组织未初始化", "ResultCode.OrgNotInitError")
	UserConfigNotExist     = errors.AddResultCodeInfo(400010, "用户配置不存在", "Result.UserConfigNotExist")
	OrgNotExist            = errors.AddResultCodeInfo(400011, "组织不存在", "ResultCode.OrgNotExist")
	OrgInitError           = errors.AddResultCodeInfo(400012, "组织初始化异常", "ResultCode.OrgInitError")
	OrgOwnTransferError    = errors.AddResultCodeInfo(400013, "非组织创建者不能更改信息", "ResultCode.OrgOwnTransferError")
	OrgOutInfoNotExist     = errors.AddResultCodeInfo(400014, "组织外部信息不存在", "ResultCode.OrgOutInfoNotExist")
	UserOutInfoNotExist    = errors.AddResultCodeInfo(400015, "用户外部信息不存在", "ResultCode.UserOutInfoNotExist")
	UserOutInfoNotError    = errors.AddResultCodeInfo(400016, "用户外部信息错误", "ResultCode.UserOutInfoNotError")
	OrgCodeAlreadySetError = errors.AddResultCodeInfo(400017, "组织网址不能二次修改", "ResultCode.OrgCodeAlreadySetError")
	//OrgCodeLenError           = errors.AddResultCodeInfo(400018, "组织网址后缀长度请控制在64字以内", "ResultCode.OrgWebSiteSettingLenError")
	OrgCodeLenError   = errors.AddResultCodeInfo(400018, "组织网址后缀只能输入20个字符,包含数字和英文", "ResultCode.OrgWebSiteSettingLenError")
	OrgCodeExistError = errors.AddResultCodeInfo(400019, "组织网址后缀已被占用，请重新输入", "ResultCode.OrgCodeExistError")
	//OrgAddressLenError        = errors.AddResultCodeInfo(400020, "详情地址不得超过256字", "ResultCode.OrgAddressLenError")
	OrgAddressLenError        = errors.AddResultCodeInfo(400020, "详情地址不得超过100字", "ResultCode.OrgAddressLenError")
	OrgLogoLenError           = errors.AddResultCodeInfo(400021, "组织logo路径长度不能超过512字", "ResultCode.OrgLogoLenError")
	OrgUserRoleModifyError    = errors.AddResultCodeInfo(400022, "无权修改当前角色", "ResultCode.OrgUserRoleModifyError")
	OrgRoleGroupNotExist      = errors.AddResultCodeInfo(400023, "角色分组不存在", "ResultCode.OrgRoleGroupNotExist")
	OrgUserUnabled            = errors.AddResultCodeInfo(400024, "您已被当前组织禁止访问，请联系管理员解除限制", "ResultCode.OrgUserUnabled")
	OrgRoleNoExist            = errors.AddResultCodeInfo(400025, "角色不存在", "ResultCode.OrgRoleNoExist")
	OrgUserDeleted            = errors.AddResultCodeInfo(400026, "您已被该组织移除", "ResultCode.OrgUserDeleted")
	OrgUserCheckStatusUnabled = errors.AddResultCodeInfo(400027, "您未通过该组织的审核", "ResultCode.OrgUserCheckStatusUnabled")

	RepeatProjectName = errors.AddResultCodeInfo(400101, "项目名重复", "ResultCode.RepeatProjectName")
	NotAllInitOrgUser = errors.AddResultCodeInfo(400102, "当前成员或负责人不属于组织成员", "Result.NotAllInitOrgUser")

	ExistingNotFinishedSubTask = errors.AddResultCodeInfo(400103, "当前任务下还有未完成的子任务", "Result.ExistingSubTask")
	VerifyOrgError             = errors.AddResultCodeInfo(400104, "存在无效用户，请刷新重试", "Result.VerifyOrgError")
	ProcessNotExist            = errors.AddResultCodeInfo(400105, "流程不存在", "Result.ProcessNotExist")
	IssueNotExist              = errors.AddResultCodeInfo(400106, "问题不存在", "Result.IssueNotExist")
	ProcessStatusNotExist      = errors.AddResultCodeInfo(400107, "流程状态不存在", "Result.ProcessStatusNotExist")
	NotAllowQuitProject        = errors.AddResultCodeInfo(400108, "负责人不允许退出项目", "Result.NotAllowQuitProject")
	NotProjectParticipant      = errors.AddResultCodeInfo(400109, "抱歉，您不是当前项目成员", "Result.NotProjectParticipant")
	PriorityNotExist           = errors.AddResultCodeInfo(400110, "优先级不存在", "Result.PriorityNotExist")
	ProjectPreCodeExist        = errors.AddResultCodeInfo(400111, "项目前缀编号已存在，请手动输入", "Result.ProjectPreCodeExist")

	RepeatProjectPrecode                   = errors.AddResultCodeInfo(400112, "项目前缀编号重复", "ResultCode.RepeatProjectPrecode")
	CreateProjectTimeError                 = errors.AddResultCodeInfo(400113, "项目截至时间必须大于开始时间", "ResultCode.CreateProjectTimeError")
	ParentIssueNotExist                    = errors.AddResultCodeInfo(400114, "父任务不存在", "ResultCode.ParentIssueNotExist")
	ExistingSubTask                        = errors.AddResultCodeInfo(400115, "删除失败，当前任务下还有未删除的子任务", "Result.ExistingSubTask")
	IssueAlreadyBeDeleted                  = errors.AddResultCodeInfo(400116, "任务不存在或已被删除", "Result.IssueAlreadyBeDeleted")
	ProcessProcessStatusRelationError      = errors.AddResultCodeInfo(400117, "流程状态关联异常", "Result.ProcessProcessStatusRelationError")
	ProcessProcessStatusInitStatueNotExist = errors.AddResultCodeInfo(400118, "流程初始状态不存在", "Result.ProcessProcessStatusInitStatueNotExist")
	ProjectNotExist                        = errors.AddResultCodeInfo(400119, "项目不存在", "Result.ProjectNotExist")
	RoleNotExist                           = errors.AddResultCodeInfo(400120, "角色不存在", "Result.RoleNotExist")
	RoleOperationNotExist                  = errors.AddResultCodeInfo(400121, "角色操作不存在", "Result.RoleOperationNotExist")
	GetUserRoleError                       = errors.AddResultCodeInfo(400122, "获取用户角色时发生异常", "Result.GetUserRoleError")
	ProjectNotInit                         = errors.AddResultCodeInfo(400123, "项目尚未初始化", "Result.ProjectNotInit")
	GetUserInfoError                       = errors.AddResultCodeInfo(400124, "获取用户信息异常", "Result.GetUserInfoError")
	IssueCondAssemblyError                 = errors.AddResultCodeInfo(400125, "任务查询条件封装异常", "Result.IssueCondAssemblyError")
	IssueDetailNotExist                    = errors.AddResultCodeInfo(400126, "任务详情不存在", "Result.IssueDetailNotExist")
	AlreadyStarProject                     = errors.AddResultCodeInfo(400127, "项目已关注", "Result.AlreadyStarProject")
	NotYetStarProject                      = errors.AddResultCodeInfo(400128, "项目尚未关注", "Result.NotYetStarProject")
	TargetNotExist                         = errors.AddResultCodeInfo(400129, "操作对象不存在", "Result.TargetNotExist")
	InvalidResourceType                    = errors.AddResultCodeInfo(400130, "资源类型有误", "Result.InvalidResourceType")
	ProjectObjectTypeProcessNotExist       = errors.AddResultCodeInfo(400131, "项目对象类型对应的流程不存在", "Result.ProjectObjectTypeProcessNotExist")
	IterationExistingNotFinishedTask       = errors.AddResultCodeInfo(400132, "当前迭代存在未完成的任务", "Result.IterationExistingNotFinishedTask")
	ProjectTypeNotExist                    = errors.AddResultCodeInfo(400133, "项目类型不存在", "Result.ProjectTypeNotExist")
	OssPolicyTypeError                     = errors.AddResultCodeInfo(400134, "错误的策略类型", "Result.OssPolicyTypeError")
	RelationIssueError                     = errors.AddResultCodeInfo(400135, "关联的任务有误", "Result.RelationIssueError")
	ParentIssueRelationChildIssueError     = errors.AddResultCodeInfo(400136, "父子任务不能关联", "Result.ParentIssueRelationChildIssueError")
	IterationNotExist                      = errors.AddResultCodeInfo(400137, "迭代不存在", "Result.IterationNotExist")
	ProjectNotRelatedError                 = errors.AddResultCodeInfo(400138, "项目未关联对应的资源", "Result.ProjectNotRelatedError")
	SourceNotExist                         = errors.AddResultCodeInfo(400139, "来源不存在", "Result.SourceNotExist")
	IssueObjectTypeNotExist                = errors.AddResultCodeInfo(400140, "任务类型不存在", "Result.IssueObjectTypeNotExist")
	ResourceNotExist                       = errors.AddResultCodeInfo(400141, "资源不存在", "Result.ResourceNotExist")
	ProjectTypeNormalError                 = errors.AddResultCodeInfo(400142, "项目不是普通任务", "Result.ProjectTypeNormalError")
	InviteCodeInvalid                      = errors.AddResultCodeInfo(400143, "邀请链接失效", "Result.InviteCodeInvalid")
	UnSupportLoginType                     = errors.AddResultCodeInfo(400144, "不支持的登录方式", "Result.UnSupportLoginType")
	ProjectIsFilingYet                     = errors.AddResultCodeInfo(400145, "项目已归档", "Result.ProjectIsFilingYet")
	LastProjectObjectType                  = errors.AddResultCodeInfo(400146, "最后一个泳道无法删除", "Result.LastProjectObjectType")
	PasswordEmptyError                     = errors.AddResultCodeInfo(400147, "请输入密码", "Result.PasswordEmptyError")
	PasswordNotSetError                    = errors.AddResultCodeInfo(400148, "密码未设置", "Result.PasswordNotSetError")
	PasswordNotMatchError                  = errors.AddResultCodeInfo(400149, "密码验证错误", "Result.PasswordNotMatchError")
	ParentIssueHasParent                  = errors.AddResultCodeInfo(400150, "子任务不允许创建子任务", "Result.ParentIssueHasParent")
	CreateIssueFail                  = errors.AddResultCodeInfo(400151, "创建任务失败", "Result.CreateIssueFail")

	IssueStatusUpdateError                = errors.AddResultCodeInfo(400202, "任务状态更新失败", "Result.IssueStatusUpdateError")
	UserConfigUpdateError                 = errors.AddResultCodeInfo(400203, "用户设置更新失败", "Result.UserConfigUpdateError")
	UserConfigInsertError                 = errors.AddResultCodeInfo(400204, "用户设置失败", "Result.UserConfigUpdateError")
	IterationIssueRelateError             = errors.AddResultCodeInfo(400205, "迭代和任务关联失败", "Result.IterationIssueRelateError")
	IterationStatusUpdateError            = errors.AddResultCodeInfo(400206, "迭代状态更新失败", "Result.IterationStatusUpdateError")
	IssueRelationUpdateError              = errors.AddResultCodeInfo(400207, "任务关联更新失败", "Result.RelationUpdateError")
	ProjectStatusUpdateError              = errors.AddResultCodeInfo(400208, "项目状态更新失败", "Result.ProjectStatusUpdateError")
	IssueProjectObjectTypeNotParttenError = errors.AddResultCodeInfo(400209, "任务项目对象类型不匹配", "Result.IssueProjectObjectTypeNotParttenError")
	IssueOwnerCantBeNull                  = errors.AddResultCodeInfo(400301, "任务负责人不能为空", "Result.IssueOwnerCantBeNull")
	DepartmentNotExist                    = errors.AddResultCodeInfo(400302, "部门不存在", "Result.DepartmentNotExist")
	ParentDepartmentNotExist              = errors.AddResultCodeInfo(400303, "父部门不存在", "Result.ParentDepartmentNotExist")
	TopDepartmentNotExist                 = errors.AddResultCodeInfo(400304, "顶级部门不存在", "Result.TopDepartmentNotExist")
	ProjectObjectTypeCantBeNullError      = errors.AddResultCodeInfo(400305, "项目对象类型不能为空", "Result.ProjectObjectTypeCantBeNullError")
	PlanEndTimeInvalidError               = errors.AddResultCodeInfo(400306, "计划结束时间需要大于开始时间", "Result.PlanEndTimeInvalidError")
	//OrgNameLenError                            = errors.AddResultCodeInfo(400307, "组织名称长度请控制在256字以内", "Result.OrgNameLenError")
	OrgNameLenError                            = errors.AddResultCodeInfo(400307, "组织名称包含非法字符或超出20个字符", "Result.OrgNameLenError")
	UserNameLenError                           = errors.AddResultCodeInfo(400308, "姓名包含非法字符或超出20个字符", "Result.UserNameLenError")
	RepeatTag                                  = errors.AddResultCodeInfo(400309, "标签已存在", "Result.RepeatTag")
	IssueSortReferenceError                    = errors.AddResultCodeInfo(400310, "任务排序参照物不能为空", "Result.IssueSortReferenceError")
	IssueSortReferenceInvalidError             = errors.AddResultCodeInfo(400311, "任务排序参照物无效", "Result.IssueSortReferenceInvalidError")
	DateRangeError                             = errors.AddResultCodeInfo(400312, "时间范围错误", "Result.DateRangeError")
	ImportDataEmpty                            = errors.AddResultCodeInfo(400313, "导入数据为空", "Result.ImportDataEmpty")
	ImportFileNotExist                         = errors.AddResultCodeInfo(400314, "未上传数据文件", "Result.ImportFileNotExist")
	NotDefaultStyle                            = errors.AddResultCodeInfo(400315, "样式无效", "Result.NotDefaultStyle")
	LengthOutOfLimit                           = errors.AddResultCodeInfo(400316, "标签为空或长度超出限制", "Result.LengthOutOfLimit")
	DateParseError                             = errors.AddResultCodeInfo(400317, "时间解析异常", "Result.DateParseError")
	DailyProjectReportError                    = errors.AddResultCodeInfo(400318, "当日项目已发送", "Result.DailyProjectReportError")
	PageInvalidError                           = errors.AddResultCodeInfo(400319, "页码无效", "Result.PageInvalidError")
	PageSizeInvalidError                       = errors.AddResultCodeInfo(400320, "页长无效", "Result.PageSizeInvalidError")
	UserOrgNotRelation                         = errors.AddResultCodeInfo(400321, "用户不是该组织成员", "Result.UserOrgNotRelation")
	UserDisabledError                          = errors.AddResultCodeInfo(400322, "已经被组织禁用", "Result.UserDisabledError")
	InvalidImportFile                          = errors.AddResultCodeInfo(400323, "文件格式有误，请上传xls、xlsx格式的文件", "Result.InvalidImportFile")
	FileParseFail                              = errors.AddResultCodeInfo(400324, "文件解析失败,请下载最新文件模板或检查文件内容", "Result.FileParseFail")
	TooLargeImportData                         = errors.AddResultCodeInfo(400325, "导入任务数据过大", "Result.TooLargeImportData")
	TooLongProjectRemark                       = errors.AddResultCodeInfo(400326, "项目简介应少于500字", "Result.TooLongProjectRemark")
	ProjectCodeLenError                        = errors.AddResultCodeInfo(400327, "项目编号长度不得超过64个字", "Result.ProjectCodeLenError")
	ProjectNameLenError                        = errors.AddResultCodeInfo(400328, "项目名称长度不得超过256个字", "Result.ProjectNameLenError")
	ProjectPreCodeLenError                     = errors.AddResultCodeInfo(400329, "项目前缀编号长度不得超过16个字", "Result.ProjectPreCodeLenError")
	ProjectRemarkLenError                      = errors.AddResultCodeInfo(400330, "项目描述长度不得超过512个字", "Result.ProjectRemarkLenError")
	ProjectIsArchivedWhenModifyIssue           = errors.AddResultCodeInfo(400331, "不允许操作归档项目下的任务", "Result.ProjectIsArchivedWhenModifyIssue")
	NoPrivateProjectPermissions                = errors.AddResultCodeInfo(400332, "没有私有项目操作权限", "Result.NoPrivateProjectPermissions")
	ChildIssueForFirst                         = errors.AddResultCodeInfo(400333, "第一条任务不能是子任务", "Result.ChildIssueForFirst")
	ProjectObjectTypeSameName                  = errors.AddResultCodeInfo(400334, "泳道名字重复", "Result.ProjectObjectTypeSameName")
	ProjectNameEmpty                           = errors.AddResultCodeInfo(400335, "项目名称不能为空", "Result.ProjectNameEmpty")
	UpdateMemberIdsIsEmptyError                = errors.AddResultCodeInfo(400336, "变动的成员列表为空", "Result.UpdateMemberIdsIsEmptyError")
	UpdateMemberStatusFail                     = errors.AddResultCodeInfo(400337, "修改成员状态失败", "Result.UpdateMemberStatusFail")
	CantUpdateStatusWhenParentIssueIsCompleted = errors.AddResultCodeInfo(400338, "父任务已完成，无法修改子任务状态", "Result.CantUpdateStatusWhenParentIssueIsCompleted")
	RoleNameLenErr                             = errors.AddResultCodeInfo(400339, "角色名称不能为空且得超过10个字符", "Result.RoleNameLenErr")
	DefaultRoleCantModify                      = errors.AddResultCodeInfo(400340, "默认角色不允许编辑", "Result.DefaultRoleCantModify")
	RoleModifyBusy                             = errors.AddResultCodeInfo(400341, "角色更新繁忙", "Result.RoleEditBusy")
	RoleNameRepeatErr                          = errors.AddResultCodeInfo(400342, "角色名称重复", "Result.RoleNameRepeatErr")
	CannotRemoveProjectOwner                   = errors.AddResultCodeInfo(400343, "项目负责人不能被移除", "Result.CannotRemoveProjectOwner")
	DefaultRoleNameErr                         = errors.AddResultCodeInfo(400344, "与系统角色名称冲突", "Result.DefaultRoleNameErr")
	SourceChannelNotDefinedError               = errors.AddResultCodeInfo(400355, "来源通道未定义", "Result.SourceChannelNotDefinedError")
	OrgNotNeedInitError                        = errors.AddResultCodeInfo(400356, "组织已存在，不需要初始化", "Result.OrgNotNeedInitError")
	IssueCommentLenError                       = errors.AddResultCodeInfo(400357, "评论不得为空且不能超过500字", "Result.IssueCommentLenError")
	//IssueRemarkLenError                        = errors.AddResultCodeInfo(400358, "描述不得为空且不能超过500字", "Result.IssueRemarkLenError")
	IssueRemarkLenError        = errors.AddResultCodeInfo(400358, "描述不能超过10000字", "Result.IssueRemarkLenError")
	AuthCodeIsNull             = errors.AddResultCodeInfo(400359, "验证码不得为空", "Result.AuthCodeIsNull")
	ContactRemarkLenErr        = errors.AddResultCodeInfo(400360, "问题反馈描述不得超过512字", "Result.ContactRemarkLenErr")
	ContactResourceInfoLenErr  = errors.AddResultCodeInfo(400361, "问题反馈资源信息不得超过2048字", "Result.ContactResourceInfoLenErr")
	ContactResourceSizeErr     = errors.AddResultCodeInfo(400362, "问题反馈图片数量不能超过5个", "Result.ContactResourceSizeErr")
	PwdAlreadySettingsErr      = errors.AddResultCodeInfo(400363, "密码已设置过", "Result.PwdAlreadySettingsErr")
	PwdFormatError             = errors.AddResultCodeInfo(400364, "密码需要以字母开头，长度在6~18之间，只能包含字母、数字和下划线", "Result.PwdLengthError")
	TagNotExist                = errors.AddResultCodeInfo(400365, "标签不存在", "Result.TagNotExist")
	InvalidProjectNameError    = errors.AddResultCodeInfo(400366, "项目名不能超出20个字符", "Result.InvalidProjectNameError")
	InvalidProjectPreCodeError = errors.AddResultCodeInfo(400367, "项目前缀编号只能输入10个字符,包含数字和英文", "Result.InvalidProjectPreCodeError")
	InvalidProjectRemarkError  = errors.AddResultCodeInfo(400368, "项目简介不能超出500个字符", "Result.InvalidProjectRemarkError")
	IssueRelateTagFail = errors.AddResultCodeInfo(400369, "任务关联标签失败", "Result.IssueRelateTagFail")
	CreateTagFail 			   = errors.AddResultCodeInfo(400370, "创建标签失败", "Result.CreateTagFail")
	ProjectNoticeLenError      = errors.AddResultCodeInfo(400371, "项目公告不能超出2000字", "Result.ProjectNoticeLenError")

	InvalidSex               = errors.AddResultCodeInfo(400401, "性别不在正常范围内", "Result.InvalidSex")
	IssueTitleError          = errors.AddResultCodeInfo(400402, "任务标题包含非法字符或超出200个字符", "Result.IssueTitleError")
	FolderIdNotExistError    = errors.AddResultCodeInfo(400403, "文件夹不存在", "Result.FolderIdNotExistError")
	InvalidResourceNameError = errors.AddResultCodeInfo(400404, "文件名包含非法字符或超出300个字符", "Result.InvalidResourceNameError")
	InvalidFolderNameError   = errors.AddResultCodeInfo(400405, "文件夹名包含非法字符或超出30个字符", "Result.InvalidResourceNameError")
	InvalidFolderIdsError    = errors.AddResultCodeInfo(400406, "无效的文件夹ids", "Result.InvalidFolderIdsError")
	InvalidResourceIdsError  = errors.AddResultCodeInfo(400407, "无效的文件ids", "Result.InvalidResourceIdsError")
	ParentIdIsItselfError    = errors.AddResultCodeInfo(400409, "目标文件夹是自己本身,无需移动", "Result.ParentIdIsItselfError")
	ResouceNotInFolderError  = errors.AddResultCodeInfo(400410, "文件不在该文件夹下", "Result.ResouceNotInFolderError")
	ReourceTypeMismatchType  = errors.AddResultCodeInfo(400411, "文件类型不匹配", "Result.ReourceTypeMismatchType")
	EncodeNotSupport         = errors.AddResultCodeInfo(400412, "不支持的编码类型", "Result.EncodeNotSupport")
	SetUserPasswordError     = errors.AddResultCodeInfo(400413, "设置密码失败", "Result.SetUserPasswordError")
	UnBindLoginNameFail      = errors.AddResultCodeInfo(400414, "解绑登录方式失败", "Result.UnBindLoginNameFail")
	BindLoginNameFail        = errors.AddResultCodeInfo(400415, "绑定登录方式失败", "Result.BindLoginNameFail")
	NotBindAccountError      = errors.AddResultCodeInfo(400416, "该登录方式未绑定任何账号", "Result.NotBindAccountError")
	AccountAlreadyBindError      = errors.AddResultCodeInfo(400417, "该登录方式已绑定其它账号", "Result.AccountAlreadyBindError")
	EmailNotBindAccountError      = errors.AddResultCodeInfo(400418, "该邮箱未绑定任何账户，请重新输入或使用手机验证码登录", "Result.EmailNotBindAccountError")
	MobileNotBindAccountError      = errors.AddResultCodeInfo(400419, "该手机号未绑定任何账号", "Result.MobileNotBindAccountError")

	//User
	UserInitError                       = errors.AddResultCodeInfo(400501, "用户初始化失败", "Result.UserInitError")
	UserNotInitError                    = errors.AddResultCodeInfo(400502, "用户未初始化", "Result.UserNotInitError")
	UserNotExist                        = errors.AddResultCodeInfo(400503, "用户不存在", "Result.UserNotExist")
	UserInfoGetFail                     = errors.AddResultCodeInfo(400504, "用户信息获取失败", "Result.UserInfoGetFail")
	UserRegisterError                   = errors.AddResultCodeInfo(400505, "用户注册失败", "Result.UserRegisterError")
	LarkInitError                       = errors.AddResultCodeInfo(400506, "示例数据已初始化", "Result.LarkInitError")
	UserSexFail                         = errors.AddResultCodeInfo(400507, "用户性别错误", "Result.UserSexFail")
	UserNameEmpty                       = errors.AddResultCodeInfo(400508, "用户姓名不能为空串", "Result.UserNameEmpty")
	EmailNotRegisterError               = errors.AddResultCodeInfo(400509, "当前邮箱未注册", "Result.EmailNotRegisterError")
	EmailNotBindError                   = errors.AddResultCodeInfo(400510, "邮箱未绑定", "Result.EmailNotBindError")
	MobileNotBindError                  = errors.AddResultCodeInfo(400511, "手机号未绑定", "Result.MobileNotBindError")
	EmailAlreadyBindError               = errors.AddResultCodeInfo(400512, "邮箱已绑定, 请先解绑", "Result.EmailAlreadyBindError")
	MobileAlreadyBindError              = errors.AddResultCodeInfo(400513, "手机号已绑定， 请先解绑", "Result.MobileAlreadyBindError")
	EmailAlreadyBindByOtherAccountError = errors.AddResultCodeInfo(400514, "该邮箱已被其他账户绑定", "Result.EmailAlreadyBindByOtherAccountError")
	MobileAlreadyBindOtherAccountError  = errors.AddResultCodeInfo(400515, "该手机号已被其他账户绑定", "Result.MobileAlreadyBindOtherAccountError")

	//动态
	TrendsCreateError     = errors.AddResultCodeInfo(401001, "动态创建失败", "Result.TrendsCreateError")
	TrendsObjTypeNilError = errors.AddResultCodeInfo(401002, "对象id有值的情况下对象类型不能为空", "Result.TrendsObjTypeNilError")
	TrendsObjIdNilError   = errors.AddResultCodeInfo(401003, "对象类型有值的情况下对象id不能为空", "Result.TrendsObjIdNilError")

	// 项目对象类型不存在
	ProjectObjectTypeNotExist            = errors.AddResultCodeInfo(402001, "任务栏不存在", "Result.ProjectObjectTypeNotExist")
	ProjectTypeProjectObjectTypeNotExist = errors.AddResultCodeInfo(402002, "项目类型与项目对象类型关联不存在", "Result.ProjectTypeProjectObjectTypeNotExist")
	ProjectObjectTypeDeleteFailExistIssue      = errors.AddResultCodeInfo(402003, "任务栏中存在任务不可删除，请先将该任务栏中的任务移出", "Result.ProjectTypeDeleteFailExistIssue")
	InvalidProjectObjectTypeName      = errors.AddResultCodeInfo(402004, "任务栏名称不能为空且不能超过30字", "Result.InvalidProjectObjectTypeName")
	CannotMoveChildIssue      = errors.AddResultCodeInfo(402005, "子任务不可单独移动任务栏", "Result.CannotMoveChildIssue")
	UpdateIssueProjectObjectTypeFail      = errors.AddResultCodeInfo(402006, "移动任务栏失败", "Result.UpdateIssueProjectObjectTypeFail")

	//domain
	ProjectDomainError    = errors.AddResultCodeInfo(405001, "项目领域出错", "Result.ProjectDomainError")
	IssueDomainError      = errors.AddResultCodeInfo(405002, "任务领域出错", "Result.IssueDomainError")
	UserDomainError       = errors.AddResultCodeInfo(405003, "用户领域出错", "Result.UserDomainError")
	BaseDomainError       = errors.AddResultCodeInfo(405004, "领域出错", "Result.BaseDomainError")
	TrendDomainError      = errors.AddResultCodeInfo(405005, "动态领域出错", "Result.TrendDomainError")
	IterationDomainError  = errors.AddResultCodeInfo(405006, "迭代领域出错", "Result.IterationDomainError")
	ObjectTypeDomainError = errors.AddResultCodeInfo(405007, "对象类型领域出错", "Result.ObjectTypeDomainError")
	ResourceDomainError   = errors.AddResultCodeInfo(405008, "资源领域出错", "Result.ResourceDomainError")
	ProcessDomainError    = errors.AddResultCodeInfo(405009, "流程领域出错", "Result.ProcessDomainError")
	DepartmentDomainError = errors.AddResultCodeInfo(405010, "部门领域出错", "Result.DepartmentDomainError")

	//权限验证领域
	IllegalityRoleOperation = errors.AddResultCodeInfo(407001, "非法的操作code", "Result.IllegalityRoleOperation")
	UserRoleNotDefinition   = errors.AddResultCodeInfo(407002, "用户角色未定义", "Result.UserRoleNotDefinition")
	NoOperationPermissions  = errors.AddResultCodeInfo(407003, "没有操作权限", "Result.NoOperationPermissions")
	PermissionNotExist      = errors.AddResultCodeInfo(407004, "权限项不存在", "Result.PermissionNotExist")

	//>>>dingtalk open api error
	SuiteTicketNotExistError      = errors.AddResultCodeInfo(600001, "suiteTicket失效或不存在", "ResultCode.SuiteTicketNotExistError")
	DingTalkOpenApiCallError      = errors.AddResultCodeInfo(600002, "钉钉OpenApi调用异常", "ResultCode.DingTalkOpenApiCallError")
	DingTalkAvoidCodeInvalidError = errors.AddResultCodeInfo(600003, "钉钉免登code失效", "ResultCode.DingTalkAvoidCodeInvalidError")
	DingTalkClientError           = errors.AddResultCodeInfo(600004, "钉钉Client获取失败", "ResultCode.DingTalkClientError")
	DingTalkGetUserInfoError      = errors.AddResultCodeInfo(600005, "钉钉获取用户信息失败", "ResultCode.DingTalkGetUserInfoError")
	DingTalkOrgInitError = errors.AddResultCodeInfo(600006, "钉钉企业初始化失败", "ResultCode.DingTalkOrgInitError")
	DingTalkConfigError = errors.AddResultCodeInfo(600007, "钉钉配置错误", "ResultCode.DingTalkConfigError")

	//>>> 飞书 open api err
	FeiShuOpenApiCallError                = errors.AddResultCodeInfo(606001, "飞书OpenApi调用异常", "ResultCode.FeiShuOpenApiCallError")
	FeiShuAppTicketNotExistError          = errors.AddResultCodeInfo(606002, "飞书AppTicket不存在", "ResultCode.FeiShuAppTicketNotExistError")
	FeiShuConfigNotExistError             = errors.AddResultCodeInfo(606003, "飞书配置不存在", "ResultCode.FeiShuConfigNotExistError")
	FeiShuClientTenantError               = errors.AddResultCodeInfo(606004, "飞书客户端获取失败", "ResultCode.FeiShuClientTenantError")
	FeiShuGetAppAccessTokenError          = errors.AddResultCodeInfo(606005, "飞书获取AppAccessToken失败", "ResultCode.FeiShuGetAppAccessTokenError")
	FeiShuGetTenantAccessTokenError       = errors.AddResultCodeInfo(606006, "飞书获取TenantAccessToken失败", "ResultCode.FeiShuGetTenantAccessTokenError")
	FeiShuAuthCodeInvalid                 = errors.AddResultCodeInfo(606007, "飞书用户授权失败", "ResultCode.FeiShuAuthCodeInvalid")
	FeiShuCardCallSignVerifyError         = errors.AddResultCodeInfo(606008, "飞书卡片回调签名校验失败", "ResultCode.FeiShuCardCallSigVerifyError")
	FeiShuCardCallMsgRepetError           = errors.AddResultCodeInfo(606009, "飞书卡片消息重复推送", "ResultCode.FeiShuCardCallMsgRepetError")
	FeiShuUserNotInAppUseScopeOfAuthority = errors.AddResultCodeInfo(606010, "不在应用使用授权范围内，请联系组织管理员解决~", "ResultCode.FeiShuUserNotInAppUseScopeOfAuthority")

	//Login Error
	SMSLoginCodeSendError                = errors.AddResultCodeInfo(601001, "登录验证码发送失败", "ResultCode.SMSLoginCodeSendError")
	SMSPhoneNumberFormatError            = errors.AddResultCodeInfo(601002, "手机号格式错误，请重新输入", "ResultCode.SMSPhoneNumberFormatError")
	SMSSendLimitError                    = errors.AddResultCodeInfo(601003, "发送过于频繁（服务商）", "ResultCode.SMSSendLimitError")
	SMSSendTimeLimitError                = errors.AddResultCodeInfo(601004, "发送过于频繁", "ResultCode.SMSSendTimeLimitError")
	SMSLoginCodeInvalid                  = errors.AddResultCodeInfo(601005, "验证码已失效，请重新获取", "ResultCode.SMSLoginCodeInvalid")
	SMSLoginCodeNotMatch                 = errors.AddResultCodeInfo(601006, "验证码错误，请重新获取", "ResultCode.SMSLoginCodeNotMatch")
	SMSLoginCodeVerifyFailTimesOverLimit = errors.AddResultCodeInfo(601007, "验证码错误，失败次数过多，请重新发送", "ResultCode.SMSLoginCodeVerifyFailTimesOverLimit")
	PwdLoginCodeNotMatch                 = errors.AddResultCodeInfo(601008, "图形验证码错误", "ResultCode.PwdLoginCodeNotMatch")
	PwdLoginUsrOrPwdNotMatch             = errors.AddResultCodeInfo(601009, "用户名或密码错误", "ResultCode.PwdLoginUsrOrPwdNotMatch")

	EmailFormatErr       = errors.AddResultCodeInfo(602001, "邮箱格式错误", "ResultCode.EmailFormatErr")
	EmailSubjectEmptyErr = errors.AddResultCodeInfo(602002, "邮箱标题不能为空", "ResultCode.EmailSubjectEmptyErr")
	EmailSendErr         = errors.AddResultCodeInfo(602003, "邮件发送失败", "ResultCode.EmailSendErr")

	NotSupportedContactAddressType = errors.AddResultCodeInfo(603001, "不支持的联系方式类型", "ResultCode.NotSupportedContactAddressType")
	NotSupportedAuthCodeType       = errors.AddResultCodeInfo(603002, "不支持的验证码类型", "ResultCode.NotSupportedAuthCodeType")
	NotSupportedRegisterType       = errors.AddResultCodeInfo(603003, "暂时不支持该注册方式", "ResultCode.NotSupportedRegisterType")

	//非法对象
	IllegalityPriority    = errors.AddResultCodeInfo(800001, "非法优先级", "Result.IllegalityPriority")
	IllegalityOwner       = errors.AddResultCodeInfo(800002, "非法负责人", "Result.IllegalityOwner")
	IllegalityFollower    = errors.AddResultCodeInfo(800003, "非法关注人", "Result.IllegalityFollower")
	IllegalityParticipant = errors.AddResultCodeInfo(800004, "非法参与人", "Result.IllegalityParticipant")
	IllegalityProject     = errors.AddResultCodeInfo(800005, "非法项目", "Result.IllegalityProject")
	IllegalityOrg         = errors.AddResultCodeInfo(800006, "非法组织", "Result.IllegalityOrg")
	IllegalityIteration   = errors.AddResultCodeInfo(800007, "非法的迭代", "Result.IllegalityIteration")
	IllegalityIssue       = errors.AddResultCodeInfo(800008, "任务不存在或已被删除", "Result.IllegalityIssue")
	IllegalityMQTTChannelType       = errors.AddResultCodeInfo(800009, "非法的通道类型", "Result.IllegalityMQTTChannelType")

	//项目关联对象
	ProjectRelationNotExist = errors.AddResultCodeInfo(900001, "项目关联对象不存在", "Result.ProjectRelationNotExist")
)
