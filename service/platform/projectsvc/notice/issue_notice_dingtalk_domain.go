package notice

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/extra/dingtalk"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
)

func GetDingTalkNormalMsg(issueNoticeBo bo.IssueNoticeBo) (*sdk.WorkNoticeMsg, error) {
	orgId := issueNoticeBo.OrgId
	operatorId := issueNoticeBo.OperatorId
	pushType := issueNoticeBo.PushType

	operatorBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(consts.AppSourceChannelDingTalk, orgId, operatorId)
	if err != nil {
		log.Errorf("查询组织 %d 用户 %d 信息出现异常 %v", orgId, operatorId, err)
		return nil, err
	}

	projectInfo, err1 := domain.LoadProjectAuthBo(orgId, issueNoticeBo.ProjectId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}

	title := operatorBaseInfo.Name + " 创建了新的任务"
	if pushType == consts.PushTypeUpdateIssue {
		title = operatorBaseInfo.Name + " 更新了任务内容"
	} else if pushType == consts.PushTypeDeleteIssue {
		title = operatorBaseInfo.Name + " 删除了任务"
	} else if pushType == consts.PushTypeUpdateIssueStatus {
		status, err := processfacade.GetProcessStatusRelaxed(orgId, issueNoticeBo.IssueStatusId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		title = operatorBaseInfo.Name + " 更新了任务状态为 " + status.Name
	}
	author := " "

	msg := &sdk.WorkNoticeMsg{
		MsgType: "oa",
		OA: &sdk.OANotice{
			MsgUrl: "http://study.ikuvn.com",
			Head: sdk.OANoticeHead{
				BgColor: "00CCFF",
				Text:    "Polaris",
			},
			Body: sdk.OANoticeBody{
				Title: &title,
				Form: &[]sdk.OANoticeBodyForm{
					{
						Key:   "任务标题: ",
						Value: issueNoticeBo.IssueTitle,
					}, {
						Key:   "操作人: ",
						Value: operatorBaseInfo.Name,
					}, {
						Key:   "所属项目：",
						Value: projectInfo.Name,
					},
				},
				Author: &author,
			},
		},
	}
	return msg, nil
}


func GetDingTalkMembersChangeMsg(issueNoticeBo bo.IssueNoticeBo, userInfos []bo.UserNoticeInfoBo, actionType int, domainType int) (*sdk.WorkNoticeMsg, errs.SystemErrorInfo) {
	orgId := issueNoticeBo.OrgId
	operatorId := issueNoticeBo.OperatorId

	operatorBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(consts.AppSourceChannelDingTalk, orgId, operatorId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	projectInfo, err1 := domain.LoadProjectAuthBo(orgId, issueNoticeBo.ProjectId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}

	action := ""
	if actionType == 1 {
		action = "移除了"
	} else {
		action = "添加了"
	}
	domain := ""
	if domainType == 1 {
		domain = "参与者"
	} else {
		domain = "关注者"
	}

	outUserNameStr := ""
	for _, userInfo := range userInfos {
		outUserNameStr += userInfo.Name + ","
	}
	if len(outUserNameStr) > 0 {
		outUserNameStr = outUserNameStr[0 : len(outUserNameStr)-1]
	}

	noticeTitle := operatorBaseInfo.Name + " " + action + domain + " " + outUserNameStr

	author := " "

	msg := &sdk.WorkNoticeMsg{
		MsgType: "oa",
		OA: &sdk.OANotice{
			MsgUrl: "http://study.ikuvn.com",
			Head: sdk.OANoticeHead{
				BgColor: "00CCFF",
				Text:    "Polaris",
			},
			Body: sdk.OANoticeBody{
				Title: &noticeTitle,
				Form: &[]sdk.OANoticeBodyForm{
					{
						Key:   "任务标题: ",
						Value: issueNoticeBo.IssueTitle,
					}, {
						Key:   "操作人: ",
						Value: operatorBaseInfo.Name,
					}, {
						Key:   "所属项目：",
						Value: projectInfo.Name,
					},
				},
				Author: &author,
			},
		},
	}
	return msg, nil
}


func IssueNoticeDingTalkPush(orgId int64, userInfos []bo.UserNoticeInfoBo, msg sdk.WorkNoticeMsg) {
	outUserIdsStr := ""
	for _, userInfo := range userInfos {
		if userInfo.OutUserId != ""{
			outUserIdsStr += userInfo.OutUserId + ","
		}
	}
	if len(outUserIdsStr) > 0 {
		outUserIdsStr = outUserIdsStr[0 : len(outUserIdsStr)-1]
	}
	orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(consts.AppSourceChannelDingTalk, orgId)
	if err != nil {
		log.Errorf("组织外部信息不存在 %v", err)
		return
	}
	client, err1 := dingtalk.GetDingTalkClientRest(orgBaseInfo.OutOrgId)
	if err1 != nil {
		log.Errorf("获取dingtalk client时发生异常 %v", err)
		return
	}
	resp, err1 := client.SendWorkNotice(&outUserIdsStr, nil, false, msg)
	if err1 != nil {
		log.Error("发送ding talk 工作通知时发生异常" + strs.ObjectToString(err))
		return
	}
	if resp.ErrCode != 0 {
		log.Error("发送ding talk 失败" + resp.ErrMsg)
		return
	}
	str, _ := json.ToJson(resp)
	log.Info("发送成功: " + str)
}