package feishu

import (
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/temp"
)

const(
	AppLinkParamNameAppId = "AppId"
	AppLinkParamNameIssueId = "IssueId"
	AppLinkParamNameParentId = "ParentId"
	AppLinkParamNameProjectId = "ProjectId"

	AppLinkDefault = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=sidebar-semi"
	//任务详情
	AppLinkIssueInfo = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=sidebar-semi&path=pages%2fPC%2fAddTask%2findex%3fid%3d{{.IssueId}}%26parentId%3d{{.ParentId}}"

	AppLinkMobileOpenQRCode = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=sidebar-semi&path=pages/PC/Home/index"
	AppLinkDefaultProjectConfigure = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=sidebar-semi&path=pages/PC/Configure/index"

	//项目统计
	AppLinkProjectStatistic = "https://applink.feishu.cn/client/mini_program/open?appId={{.AppId}}&mode=sidebar-semi&path=pages/PC/ProjectStatistics/index?projectId={{.ProjectId}}"

	AppGuide = "https://polaris.feishu.cn/docs/doccn5PXJnUDDhLee2Uksn337mb"
	AppConfigDoc = "https://polaris.feishu.cn/docs/doccn56SrnQKluE4VnLnL4MTMKe"
)

func GetDefaultAppLink() string{
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil{
		log.Error("飞书配置为空")
		return ""
	}
	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
	}
	link, err := temp.Render(AppLinkDefault, params)
	if err != nil{
		log.Error(err)
		return ""
	}
	return link
}

func GetMobileOpenQRCodeAppLink() string{
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil{
		log.Error("飞书配置为空")
		return ""
	}
	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
	}
	link, err := temp.Render(AppLinkMobileOpenQRCode, params)
	if err != nil{
		log.Error(err)
		return ""
	}
	return link
}

func GetDefaultProjectConfigureAppLink() string{
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil{
		log.Error("飞书配置为空")
		return ""
	}
	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
	}
	link, err := temp.Render(AppLinkDefaultProjectConfigure, params)
	if err != nil{
		log.Error(err)
		return ""
	}
	return link
}

func GetIssueInfoAppLink(issueId int64, parentId int64) string{
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil{
		log.Error("飞书配置为空")
		return ""
	}

	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
		AppLinkParamNameIssueId: issueId,
		AppLinkParamNameParentId: parentId,
	}

	link, err := temp.Render(AppLinkIssueInfo, params)
	if err != nil{
		log.Error(err)
		return ""
	}
	log.Infof("获取任务详情appLink: %s", link)
	return link
}

func GetProjectStatisticAppLink(projectId int64) string{
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil{
		log.Error("飞书配置为空")
		return ""
	}

	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
		AppLinkParamNameProjectId: projectId,
	}

	link, err := temp.Render(AppLinkProjectStatistic, params)
	if err != nil{
		log.Error(err)
		return ""
	}
	return link
}

func GetProjectStatisticPcLink(projectId int64) string{
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil{
		log.Error("飞书配置为空")
		return ""
	}

	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
		AppLinkParamNameProjectId: projectId,
	}

	link, err := temp.Render(fsConfig.CardJumpLink.ProjectDailyPcUrl, params)
	if err != nil{
		log.Error(err)
		return ""
	}
	return link
}

func GetPersonalStatisticPcLink() string{
	fsConfig := config.GetConfig().FeiShu
	if fsConfig == nil{
		log.Error("飞书配置为空")
		return ""
	}

	params := map[string]interface{}{
		AppLinkParamNameAppId: fsConfig.AppId,
	}

	link, err := temp.Render(fsConfig.CardJumpLink.PersonalDailyPcUrl, params)
	if err != nil{
		log.Error(err)
		return ""
	}
	log.Infof("个人统计页面地址 %s", link)
	return link
}