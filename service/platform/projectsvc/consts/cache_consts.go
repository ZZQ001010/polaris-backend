package consts

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
)

var (
	//用户配置缓存
	CacheBaseProjectInfo = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "baseinfo"
	//用户基础信息缓存key
	CacheProjectObjectTypeList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "project_object_type"
	//优先级列表
	CachePriorityList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + "priority_list"
	//项目类型
	CacheProjectTypeList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + "project_type_list"
	//项目类型与项目对象类型关联缓存
	CacheProjectTypeProjectObjectTypeList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + "project_type_project_object_type"
	//项目状态缓存
	CacheProjectProcessStatus = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "process_status"
	//任务来源列表
	CacheIssueSourceList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + "issue_source_list"
	//任务类型列表
	CacheIssueObjectTypeList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + "issue_object_type_list"
	//项目日历信息缓存
	CacheProjectCalendarInfo = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "project_calendar_info"
	//项目标签缓存
	CacheProjectTagInfo = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "tag_stat_info"
)
