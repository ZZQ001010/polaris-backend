package consts

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
)

var (
	//流程缓存
	CacheProcessList = consts.CacheKeyPrefix + consts.ProcesssvcApplicationName + consts.CacheKeyOfOrg + "process_list"
	//流程状态列表
	CacheProcessStatusList = consts.CacheKeyPrefix + consts.ProcesssvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProcess + "process_status_list"
	//状态列表
	CacheStatusList = consts.CacheKeyPrefix + consts.ProcesssvcApplicationName + consts.CacheKeyOfOrg + "process_status_list"
	//流程步骤列表
	CacheProcessStepList = consts.CacheKeyPrefix + consts.ProcesssvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProcess + "process_step_list"
)
