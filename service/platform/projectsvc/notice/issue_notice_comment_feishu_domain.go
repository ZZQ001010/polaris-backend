package notice

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
)

func GetFeiShuCommentMsg(issueNoticeBo bo.IssueNoticeBo, content string) (*vo.Card, error) {
	orgId := issueNoticeBo.OrgId
	operatorId := issueNoticeBo.OperatorId

	operatorBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(consts.AppSourceChannelFeiShu, orgId, operatorId)
	if err != nil {
		log.Errorf("查询组织 %d 用户 %d 信息出现异常 %v", orgId, operatorId, err)
		return nil, err
	}

	msg := &vo.Card{
		Header: &vo.CardHeader{
			Title: &vo.CardHeaderTitle{
				Tag:     "plain_text",
				Content: "@我",
			},
		},
		Elements: []interface{}{
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
							Content: "**" + operatorBaseInfo.Name + ": ** " + util.RenderCommentContentToMarkDown(content),
						},
					},
				},
			},
		},
	}
	return msg, nil
}

func PushFsIssueComment(issueNoticeBo bo.IssueNoticeBo, content string, mentionedUserIds []int64) {
	userMap := map[int64]bool{}
	if mentionedUserIds == nil || len(mentionedUserIds) == 0 {
		return
	}
	for _, userId := range mentionedUserIds {
		userMap[userId] = true
	}

	//目前只有飞书
	bePushedUserInfos := GetNormalUserIdsWithFilter(issueNoticeBo, consts.AppSourceChannelFeiShu, PushNoticeTargetTypeCommentAtMsg, userMap)

	if bePushedUserInfos != nil && len(bePushedUserInfos) > 0 {
		card, err := GetFeiShuCommentMsg(issueNoticeBo, content)
		if err != nil {
			log.Error(err)
			return
		}
		IssueNoticeFeiShuPush(issueNoticeBo.OrgId, bePushedUserInfos, card)
	}
}
