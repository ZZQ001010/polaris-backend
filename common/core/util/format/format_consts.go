package format

const (
	EmailPattern               = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //用户邮箱
	PasswordPattern            = `^[a-zA-Z]\w*$`                               //用户密码
	UserNamePattern            = "^[ 0-9A-Za-z]{1,20}$"                        //用户名
	OrgNamePattern             = "^[\u4e00-\u9fa5|0-9|a-zA-Z]{1,20}$"          //组织名
	OrgCodePattern             = "^[0-9|a-zA-Z]{1,20}$"                        //网址后缀编号
	OrgAdressPattern           = `^.{0,100}$`                                  //组织地址
	ProjectNamePattern         = `^.{1,30}$`                                   //项目名
	ProjectPreviousCodePattern = `^[a-zA-Z|0-9]{0,10}$`                        //项目前缀编号
	ProjectRemarkPattern       = `^[\s\S]{0,500}$`                             //项目描述(简介)
	//ProjectNoticePattern       = `^.{0,2000}$`                                 //项目公告
	IssueNamePattern = `^.{1,200}$` //任务名
	//IssueRemarkPattern           = `^.{0,10000}$`                                //任务描述(详情)
	IssueCommenPattern           = `^.{1,500}$`                                 //任务评论
	ProjectObjectTypeNamePattern = `^.{1,30}$`                                  //标题栏名
	ResourceNamePattern          = `^[^\\\\/:*?\"<>|]{1,300}(\.[a-zA-Z0-9]+)?$` //资源名
	FolderNamePattern            = `^[^\\\\/:*?\"<>|]{1,30}$`                   //文件夹名
	RoleNamePattern              = "^[a-zA-Z|0-9]{1,20}$"                       //角色名
)

const (
	ChinesePattern  = "[\u4e00-\u9fa5]+?"
	AllBlankPattern = "^ +$"
)
