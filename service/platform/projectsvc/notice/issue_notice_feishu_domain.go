package notice

import (
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/extra/feishu"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
	"time"
)

func GetFeiShuNormalMsg(issueNoticeBo bo.IssueNoticeBo) (*vo.Card, error) {
	orgId := issueNoticeBo.OrgId
	operatorId := issueNoticeBo.OperatorId
	pushType := issueNoticeBo.PushType

	operatorBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(consts.AppSourceChannelFeiShu, orgId, operatorId)
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
		title = operatorBaseInfo.Name + " 更新了任务"
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

	title = "⏰ " + title

	issueInfoAppLink := feishu.GetIssueInfoAppLink(issueNoticeBo.IssueId, issueNoticeBo.ParentId)

	msg := &vo.Card{
		Header: &vo.CardHeader{
			Title: &vo.CardHeaderTitle{
				Tag:     "plain_text",
				Content: title,
			},
		},
		Elements: []interface{}{
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "**所属项目**\n" + projectInfo.Name,
						},
					},
				},
			},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "**任务标题**\n" + issueNoticeBo.IssueTitle,
						},
					},
				},
			},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "**操作人: ** " + operatorBaseInfo.Name,
						},
					},
				},
			},
			vo.CardElementActionModule{
				Tag: "action",
				Actions: []interface{}{
					vo.ActionButton{
						Tag: "button",
						Text: vo.CardElementText{
							Tag:     "plain_text",
							Content: "🔍 查看详情",
						},
						Url:  issueInfoAppLink,
						Type: "default",
					},
				},
			},
		},
	}
	return msg, nil
}

func GetFeiShuMembersChangeMsg(issueNoticeBo bo.IssueNoticeBo, userInfos []bo.UserNoticeInfoBo, actionType int, domainType int) (*vo.Card, errs.SystemErrorInfo) {
	orgId := issueNoticeBo.OrgId
	operatorId := issueNoticeBo.OperatorId

	operatorBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(consts.AppSourceChannelFeiShu, orgId, operatorId)
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
	noticeTitle = "⏰ " + noticeTitle

	issueInfoAppLink := feishu.GetIssueInfoAppLink(issueNoticeBo.IssueId, issueNoticeBo.ParentId)

	msg := &vo.Card{
		Header: &vo.CardHeader{
			Title: &vo.CardHeaderTitle{
				Tag:     "plain_text",
				Content: noticeTitle,
			},
		},
		Elements: []interface{}{
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "**所属项目**\n" + projectInfo.Name,
						},
					},
				},
			},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "**任务标题**\n" + issueNoticeBo.IssueTitle,
						},
					},
				},
			},
			vo.CardElementContentModule{
				Tag: "div",
				Fields: []vo.CardElementField{
					{
						Text: vo.CardElementText{
							Tag:     "lark_md",
							Content: "**操作人: ** " + operatorBaseInfo.Name,
						},
					},
				},
			},
			vo.CardElementActionModule{
				Tag: "action",
				Actions: []interface{}{
					vo.ActionButton{
						Tag: "button",
						Text: vo.CardElementText{
							Tag:     "plain_text",
							Content: "🔍 查看详情",
						},
						Url:  issueInfoAppLink,
						Type: "default",
					},
				},
			},
		},
	}
	return msg, nil
}

func GetFeiShuOwnerChangeMsg(issueNoticeBo bo.IssueNoticeBo) (*vo.Card, errs.SystemErrorInfo) {
	planStartTime := dealPlanTime(issueNoticeBo.IssuePlanStartTime)
	planEndTime := dealPlanTime(issueNoticeBo.IssuePlanEndTime)
	issueId := issueNoticeBo.IssueId
	orgId := issueNoticeBo.OrgId
	ownerId := issueNoticeBo.AfterOwner
	operatorId := issueNoticeBo.OperatorId
	projectId := issueNoticeBo.ProjectId
	priorityId := issueNoticeBo.PriorityId

	ownerBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(consts.AppSourceChannelFeiShu, orgId, ownerId)
	if err != nil {
		log.Errorf("查询组织 %d 用户 %d 信息出现异常 %v", orgId, operatorId, err)
		return nil, err
	}

	projectInfo, err := domain.LoadProjectAuthBo(orgId, projectId)
	if err != nil{
		log.Error(err)
		return nil, err
	}

	priorityInfo, err := domain.GetPriorityById(orgId, priorityId)
	if err != nil{
		log.Error(err)
		return nil, err
	}

	noticeTitle := "⏰ 新任务"
	issueInfoAppLink := feishu.GetIssueInfoAppLink(issueNoticeBo.IssueId, issueNoticeBo.ParentId)

	elements := []interface{}{
		vo.CardElementContentModule{
			Tag: "div",
			Fields: []vo.CardElementField{
				{
					Text: vo.CardElementText{
						Tag:     "lark_md",
						Content: "**所属项目**\n" + projectInfo.Name,
					},
				},
			},
		},
		vo.CardElementContentModule{
			Tag: "div",
			Fields: []vo.CardElementField{
				{
					Text: vo.CardElementText{
						Tag:     "lark_md",
						Content: "**任务标题**\n" + issueNoticeBo.IssueTitle,
					},
				},
			},
		},
	}

	if issueNoticeBo.IssueRemark != "" {
		elements = append(elements, vo.CardElementContentModule{
			Tag: "div",
			Fields: []vo.CardElementField{
				{
					Text: vo.CardElementText{
						Tag:     "lark_md",
						Content: "**任务内容**\n" + issueNoticeBo.IssueRemark,
					},
				},
			},
		})
	}
	elements = append(elements, vo.CardElementContentModule{
		Tag: "div",
		Fields: []vo.CardElementField{
				{
					Text: vo.CardElementText{
						Tag:     "lark_md",
						Content: "**负责人: ** " + ownerBaseInfo.Name + "\t\t **优先级: **" + priorityInfo.Name,
					},
				},
			},
		},
		vo.CardElementContentModule{
			Tag: "div",
			Fields: []vo.CardElementField{
				{
					Text: vo.CardElementText{
						Tag:     "lark_md",
						Content: "**计划开始时间:**",
					},
				},
			},
		},
		vo.CardElementActionModule{
			Tag: "action",
			Actions: []interface{}{
				vo.ActionDatePicker{
					Tag:             "picker_datetime",
					InitialDatetime: planStartTime,
					Placeholder: &vo.CardElementText{
						Tag:     "plain_text",
						Content: "修改计划开始时间",
					},
					Confirm: &vo.CardElementConfirm{
						Title: &vo.CardHeaderTitle{
							Tag:     "plain_text",
							Content: "确认要修改这个任务的计划开始时间吗?",
						},
						Text: &vo.CardElementText{
							Tag:     "plain_text",
							Content: "",
						},
					},
					Value: map[string]interface{}{
						consts.FsCardValueIssueId:  issueId,
						consts.FsCardValueCardType: consts.FsCardTypeIssueRemind,
						consts.FsCardValueAction:   consts.FsCardActionUpdatePlanStartTime,
						consts.FsCardValueOrgId:    orgId,
						consts.FsCardValueUserId:   ownerId,
					},
				},
			},
		},
		vo.CardElementContentModule{
			Tag: "div",
			Fields: []vo.CardElementField{
				{
					Text: vo.CardElementText{
						Tag:     "lark_md",
						Content: "**计划截止时间:**",
					},
				},
			},
		},
		vo.CardElementActionModule{
			Tag: "action",
			Actions: []interface{}{
				vo.ActionDatePicker{
					Tag:             "picker_datetime",
					InitialDatetime: planEndTime,
					Placeholder: &vo.CardElementText{
						Tag:     "plain_text",
						Content: "修改截止时间",
					},
					Confirm: &vo.CardElementConfirm{
						Title: &vo.CardHeaderTitle{
							Tag:     "plain_text",
							Content: "确认要修改这个任务的截止时间吗?",
						},
						Text: &vo.CardElementText{
							Tag:     "plain_text",
							Content: "",
						},
					},
					Value: map[string]interface{}{
						consts.FsCardValueIssueId:  issueId,
						consts.FsCardValueCardType: consts.FsCardTypeIssueRemind,
						consts.FsCardValueAction:   consts.FsCardActionUpdatePlanEndTime,
						consts.FsCardValueOrgId:    orgId,
						consts.FsCardValueUserId:   ownerId,
					},
				},
			},
		},
		vo.CardElementActionModule{
			Tag: "action",
			Actions: []interface{}{
				vo.ActionButton{
					Tag: "button",
					Text: vo.CardElementText{
						Tag:     "lark_md",
						Content: "🔍 查看详情",
					},
					Type: "default",
					Url:  issueInfoAppLink,
				},
			},
		})

	msg := &vo.Card{
		Header: &vo.CardHeader{
			Title: &vo.CardHeaderTitle{
				Tag:     "plain_text",
				Content: noticeTitle,
			},
		},
		Elements: elements,
	}
	return msg, nil
}

func dealPlanTime(t *types.Time) string {
	if t == nil {
		return ""
	}
	return time.Time(*t).Format(consts.AppTimeFormatYYYYMMDDHHmm)
}

func IssueNoticeFeiShuPush(orgId int64, userInfos []bo.UserNoticeInfoBo, msg *vo.Card) {
	openIds := make([]string, len(userInfos))
	for i, userInfo := range userInfos {
		if userInfo.OutUserId != "" {
			openIds[i] = userInfo.OutUserId
		}
	}
	orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(consts.AppSourceChannelFeiShu, orgId)
	if err != nil {
		log.Errorf("组织外部信息不存在 %v", err)
		return
	}

	tenant, err1 := feishu.GetTenant(orgBaseInfo.OutOrgId)
	if err1 != nil {
		log.Error(err)
		return
	}

	resp, err2 := tenant.SendMessageBatch(vo.BatchMsgVo{
		OpenIds: openIds,
		MsgType: "interactive",
		Card:    msg,
	})
	if err2 != nil {
		log.Error(err2)
		return
	}

	if resp.Code != 0 {
		log.Error("发送飞书卡片通知失败" + resp.Msg)
		return
	}
	str, _ := json.ToJson(resp)
	log.Info("发送成功: " + str)
}
