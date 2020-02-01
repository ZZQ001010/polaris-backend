package notice

import (
	"container/list"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"gopkg.in/fatih/set.v0"
)

var supportChannel = []string{consts.AppSourceChannelDingTalk, consts.AppSourceChannelFeiShu}

func PushIssue(issueNoticeBo bo.IssueNoticeBo) {
	orgId := issueNoticeBo.OrgId

	//pc推送
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			log.Errorf("捕获到的错误：%s", r)
	//		}
	//	}()
	//	PushIssuePC(issueNoticeBo)
	//}()

	//求他平台推送
	for _, channel := range supportChannel{
		if CheckOrgOutInfo(channel, orgId){
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Errorf("捕获到的错误：%s", r)
					}
				}()
				PushIssueByChannel(issueNoticeBo, channel)
			}()
		}
	}
}

func PushIssueComment(issueNoticeBo bo.IssueNoticeBo, content string, mentionedUserIds []int64){
	if CheckOrgOutInfo(consts.AppSourceChannelFeiShu, issueNoticeBo.OrgId){
		PushFsIssueComment(issueNoticeBo, content, mentionedUserIds)
	}
}

func CheckOrgOutInfo(sourceChannel string, orgId int64) bool{
	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(sourceChannel, orgId)
	if err != nil{
		log.Error(err)
		return false
	}
	if baseOrgInfo.OutOrgId == ""{
		return false
	}
	return true
}

func PushIssueByChannel(issueNoticeBo bo.IssueNoticeBo, sourceChannel string) {

	pushType := issueNoticeBo.PushType
	orgId := issueNoticeBo.OrgId

	//任务正常更新
	if pushType == consts.PushTypeCreateIssue || pushType == consts.PushTypeUpdateIssue || pushType == consts.PushTypeDeleteIssue || pushType == consts.PushTypeUpdateIssueStatus {
		bePushedUserInfos := GetNormalUserIds(issueNoticeBo, sourceChannel, PushNoticeTargetTypeModifyMsg)
		log.Infof("任务推送，需要推送的用户信息列表为 %s", json.ToJsonIgnoreError(bePushedUserInfos))

		if len(bePushedUserInfos) == 0{
			log.Infof("需要推送的人员数为0，不需要推送，消息%s", json.ToJsonIgnoreError(issueNoticeBo))
		}

		if sourceChannel == consts.AppSourceChannelDingTalk{
			bePushedMsg, err := GetDingTalkNormalMsg(issueNoticeBo)
			log.Infof("Dingtalk任务推送，需要推送的消息为 %s", json.ToJsonIgnoreError(bePushedMsg))
			if err != nil {
				log.Error(err)
				return
			}
			IssueNoticeDingTalkPush(orgId, bePushedUserInfos, *bePushedMsg)
		}else if sourceChannel == consts.AppSourceChannelFeiShu{
			if pushType == consts.PushTypeCreateIssue{
				//创建事件，推送负责人专属消息
				ownerId := issueNoticeBo.AfterOwner
				if ownerId != issueNoticeBo.OperatorId{
					//如果变动后的负责人和操作人不是同一个人，对该负责人推送负责人消息
					//处理负责人信息
					err := dealOwnerInfo(issueNoticeBo, 0, issueNoticeBo.AfterOwner, orgId, sourceChannel)
					if err != nil{
						log.Error(err)
						return
					}
					//去掉负责人
					bePushedNormalUerInfos := make([]bo.UserNoticeInfoBo, 0)
					for _, bePushedUserInfo := range bePushedUserInfos{
						if bePushedUserInfo.UserId != ownerId{
							bePushedNormalUerInfos = append(bePushedNormalUerInfos, bePushedUserInfo)
						}
					}
					bePushedUserInfos = bePushedNormalUerInfos
				}
			}

			if len(bePushedUserInfos) > 0{
				bePushedMsg, err := GetFeiShuNormalMsg(issueNoticeBo)
				log.Infof("飞书任务推送，需要推送的消息为 %s", json.ToJsonIgnoreError(bePushedMsg))
				if err != nil {
					log.Error(err)
					return
				}
				IssueNoticeFeiShuPush(orgId, bePushedUserInfos, bePushedMsg)
			}
		}

	//任务人员变动
	} else if pushType == consts.PushTypeUpdateIssueMembers {
		beDeletedUserInfos, beAddedUserInfos, err := GetDifNoticeUserInfos(issueNoticeBo.OrgId, issueNoticeBo.OperatorId, issueNoticeBo.BeforeChangeParticipants, issueNoticeBo.AfterChangeParticipants, sourceChannel, PushNoticeRangeTypeParticipant, PushNoticeNotLimit)

		if err != nil {
			log.Error(err)
			return
		}
		//处理参与人信息
		err = dealParticipantInfo(issueNoticeBo, beDeletedUserInfos, beAddedUserInfos, orgId, sourceChannel)

		if err != nil {
			log.Error(err)
			return
		}

		beDeletedFollowerInfos, beAddedFollowerInfos, err := GetDifNoticeUserInfos(issueNoticeBo.OrgId, issueNoticeBo.OperatorId, issueNoticeBo.BeforeChangeFollowers, issueNoticeBo.AfterChangeFollowers, sourceChannel, PushNoticeRangeTypeAttention, PushNoticeNotLimit)

		if err != nil {
			log.Error(err)
			return
		}

		//处理关注人信息
		err = dealFollowerInfo(issueNoticeBo, beDeletedFollowerInfos, beAddedFollowerInfos, orgId, sourceChannel)

		if err != nil {
			log.Error(err)
			return
		}

		//处理负责人信息
		err = dealOwnerInfo(issueNoticeBo, issueNoticeBo.BeforeOwner, issueNoticeBo.AfterOwner, orgId, sourceChannel)
		if err != nil{
			log.Error(err)
			return
		}
	}
}

func dealFollowerInfo(issueNoticeBo bo.IssueNoticeBo, beDeletedFollowerInfos *[]bo.UserNoticeInfoBo, beAddedFollowerInfos *[]bo.UserNoticeInfoBo,
	orgId int64, sourceChannel string) errs.SystemErrorInfo {
	if sourceChannel == consts.AppSourceChannelDingTalk{
		if beDeletedFollowerInfos != nil {
			msg, err := GetDingTalkMembersChangeMsg(issueNoticeBo, *beDeletedFollowerInfos, 1, 2)
			if err != nil {
				log.Error(err)
				return err
			}
			IssueNoticeDingTalkPush(orgId, *beDeletedFollowerInfos, *msg)
		}
		if beAddedFollowerInfos != nil {
			msg, err := GetDingTalkMembersChangeMsg(issueNoticeBo, *beAddedFollowerInfos, 2, 2)
			if err != nil {
				log.Error(err)
				return err
			}
			IssueNoticeDingTalkPush(orgId, *beAddedFollowerInfos, *msg)
		}
	} else if sourceChannel == consts.AppSourceChannelFeiShu{
		if beDeletedFollowerInfos != nil {
			msg, err := GetFeiShuMembersChangeMsg(issueNoticeBo, *beDeletedFollowerInfos, 1, 2)
			if err != nil {
				log.Error(err)
				return err
			}
			IssueNoticeFeiShuPush(orgId, *beDeletedFollowerInfos, msg)
		}
		if beAddedFollowerInfos != nil {
			msg, err := GetFeiShuMembersChangeMsg(issueNoticeBo, *beAddedFollowerInfos, 2, 2)
			if err != nil {
				log.Error(err)
				return err
			}
			IssueNoticeFeiShuPush(orgId, *beAddedFollowerInfos, msg)
		}
	}


	return nil
}

func dealOwnerInfo(issueNoticeBo bo.IssueNoticeBo, beforeOwnerId int64, afterOwnerId int64, orgId int64, sourceChannel string) errs.SystemErrorInfo{
	if beforeOwnerId == afterOwnerId{
		return nil
	}
	afterOwnerInfo, err := orgfacade.GetBaseUserInfoRelaxed(sourceChannel, orgId, afterOwnerId)
	if err != nil{
		log.Error(err)
		return err
	}

	userConfig, err := orgfacade.GetUserConfigInfoRelaxed(orgId, afterOwnerId)
	if err != nil {
		log.Errorf("获取%d用户配置失败, %v", afterOwnerId, err)
		return err
	}

	if userConfig.OwnerRangeStatus != 1{
		return nil
	}

	if sourceChannel == consts.AppSourceChannelDingTalk{
		//待实现
	} else if sourceChannel == consts.AppSourceChannelFeiShu{
		msg, err := GetFeiShuOwnerChangeMsg(issueNoticeBo)
		if err != nil {
			log.Error(err)
			return err
		}
		IssueNoticeFeiShuPush(orgId, []bo.UserNoticeInfoBo{
			{
				UserId: afterOwnerInfo.UserId,
				OutUserId: afterOwnerInfo.OutUserId,
				Name: afterOwnerInfo.Name,
			},
		}, msg)
	}
	return nil
}

func dealParticipantInfo(issueNoticeBo bo.IssueNoticeBo, beDeletedUserInfos *[]bo.UserNoticeInfoBo, beAddedUserInfos *[]bo.UserNoticeInfoBo,
	orgId int64, sourceChannel string) errs.SystemErrorInfo {

	if sourceChannel == consts.AppSourceChannelDingTalk{
		if beDeletedUserInfos != nil {
			msg, err := GetDingTalkMembersChangeMsg(issueNoticeBo, *beDeletedUserInfos, 1, 1)
			if err != nil {
				log.Error(err)
				return err
			}
			IssueNoticeDingTalkPush(orgId, *beDeletedUserInfos, *msg)
		}
		if beAddedUserInfos != nil {
			msg, err := GetDingTalkMembersChangeMsg(issueNoticeBo, *beAddedUserInfos, 2, 1)
			if err != nil {
				log.Error(err)
				return err
			}
			IssueNoticeDingTalkPush(orgId, *beAddedUserInfos, *msg)
		}
	} else if sourceChannel == consts.AppSourceChannelFeiShu{
		if beDeletedUserInfos != nil {
			msg, err := GetFeiShuMembersChangeMsg(issueNoticeBo, *beDeletedUserInfos, 1, 1)
			if err != nil {
				log.Error(err)
				return err
			}
			IssueNoticeFeiShuPush(orgId, *beDeletedUserInfos, msg)
		}
		if beAddedUserInfos != nil {
			msg, err := GetFeiShuMembersChangeMsg(issueNoticeBo, *beAddedUserInfos, 2, 1)
			if err != nil {
				log.Error(err)
				return err
			}
			IssueNoticeFeiShuPush(orgId, *beAddedUserInfos, msg)
		}
	}
	return nil
}

//filter: 需要推送的人，限定推送范围
func GetNormalUserIdsWithFilter(issueNoticeBo bo.IssueNoticeBo, sourceChannel string, targetType int, filter map[int64]bool) []bo.UserNoticeInfoBo {
	orgId := issueNoticeBo.OrgId
	operatorId := issueNoticeBo.OperatorId

	ownerArray := []int64{issueNoticeBo.BeforeOwner}
	userIdsList := [][]int64{ownerArray, issueNoticeBo.BeforeChangeParticipants, issueNoticeBo.BeforeChangeFollowers}

	//返回链表
	noticeUserList := dealIssueNoticeUserIdsList(userIdsList, operatorId, orgId, targetType)

	//转换array并去重
	bePushedUserIds := make([]int64, noticeUserList.Len())
	i := 0
	for e := noticeUserList.Front(); e != nil; e = e.Next() {
		bePushedUserIds[i] = e.Value.(int64)
		i++
	}
	bePushedUserIds = slice.SliceUniqueInt64(bePushedUserIds)

	if filter != nil && len(filter) > 0{
		userIds := make([]int64, 0)
		for _, bePushedUserId := range bePushedUserIds{
			if _, ok := filter[bePushedUserId]; ok{
				userIds = append(userIds, bePushedUserId)
			}
		}
		bePushedUserIds = userIds
	}

	userNoticeInfos := make([]bo.UserNoticeInfoBo, len(bePushedUserIds))
	if len(bePushedUserIds) > 0 {
		baseUserInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(sourceChannel, orgId, bePushedUserIds)
		if err != nil {
			log.Error(err)
			return userNoticeInfos
		}
		userMap := maps.NewMap("UserId", baseUserInfos)
		for i, userId := range bePushedUserIds {
			if baseUserInfoInterface, ok := userMap[userId]; ok{
				if baseUserInfo, ok := baseUserInfoInterface.(bo.BaseUserInfoBo); ok{
					userNoticeInfos[i] = bo.UserNoticeInfoBo{
						UserId:    baseUserInfo.UserId,
						OutUserId: baseUserInfo.OutUserId,
						Name:      baseUserInfo.Name,
					}
				}
			}
		}
	}

	return userNoticeInfos
}

func GetNormalUserIds(issueNoticeBo bo.IssueNoticeBo, sourceChannel string, targetType int) []bo.UserNoticeInfoBo {
	return GetNormalUserIdsWithFilter(issueNoticeBo, sourceChannel, targetType, nil)
}

func dealIssueNoticeUserIdsList(userIdsList [][]int64, operatorId int64, orgId int64, targetType int) *list.List {

	noticeUserList := list.New()

	for i, userIds := range userIdsList {
		if userIds != nil {
			dealIssueNoticeUserIds(i, userIds, operatorId, orgId, targetType, noticeUserList)
		}
	}
	return noticeUserList
}

func dealIssueNoticeUserIds(i int, userIds []int64, operatorId int64, orgId int64, targetType int, noticeUserList *list.List) {
	for _, userId := range userIds {
		if userId == operatorId {
			continue
		}
		userConfig, err := orgfacade.GetUserConfigInfoRelaxed(orgId, userId)
		if err != nil {
			log.Errorf("获取%d用户配置失败, %v", userId, err)
			continue
		}
		if userPushRangeConfigContinueFlag(i, userConfig) {
			continue
		}
		if userPushTargetConfigContinueFlag(targetType, userConfig){
			continue
		}
		noticeUserList.PushBack(userId)
	}
}

func GetDifNoticeUserInfos(orgId int64, operatorId int64, beforeUserIds []int64, afterUserIds []int64, sourceChannel string, rangeType int, targetType int) (*[]bo.UserNoticeInfoBo, *[]bo.UserNoticeInfoBo, errs.SystemErrorInfo) {
	beforeChangeMembersSet := set.New(set.ThreadSafe)
	if beforeUserIds != nil && len(beforeUserIds) > 0{
		for _, member := range beforeUserIds {
			beforeChangeMembersSet.Add(member)
		}
	}

	afterChangeMembersSet := set.New(set.ThreadSafe)
	if afterUserIds != nil && len(afterUserIds) > 0{
		for _, member := range afterUserIds {
			afterChangeMembersSet.Add(member)
		}
	}

	deletedMembersSet := set.Difference(beforeChangeMembersSet, afterChangeMembersSet)
	addedMembersSet := set.Difference(afterChangeMembersSet, beforeChangeMembersSet)

	deletedUserInfos := ConvertUserIdsToUserInfoWithSet(orgId, operatorId, deletedMembersSet, sourceChannel, rangeType, targetType)
	addedUserInfos := ConvertUserIdsToUserInfoWithSet(orgId, operatorId, addedMembersSet, sourceChannel, rangeType, targetType)

	return deletedUserInfos, addedUserInfos, nil
}

//rangeType：范围类型
//targetType: 目标类型
func ConvertUserIdsToUserInfo(orgId int64, operatorId int64, convertUserIds []int64, sourceChannel string, rangeType int, targetType int) *[]bo.UserNoticeInfoBo {
	userNoticeInfos := &[]bo.UserNoticeInfoBo{}

	userIds := make([]int64, 0)
	for _, userId := range convertUserIds {
		if userId == operatorId {
			continue
		}
		userIds = append(userIds, userId)
	}
	userIds = slice.SliceUniqueInt64(userIds)

	baseUserInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(sourceChannel, orgId, userIds)
	if err != nil {
		log.Error(err)
		return userNoticeInfos
	}
	userMap := maps.NewMap("UserId", baseUserInfos)
	for _, userId := range userIds {
		userConfig, err := orgfacade.GetUserConfigInfoRelaxed(orgId, userId)
		if err != nil {
			log.Errorf("获取%d用户配置失败, %v", userId, err)
			continue
		}

		if userPushRangeConfigContinueFlag(rangeType, userConfig){
			continue
		}
		if userPushRangeConfigContinueFlag(targetType, userConfig){
			continue
		}

		if baseUserInfoInterface, ok := userMap[userId]; ok{
			if baseUserInfo, ok := baseUserInfoInterface.(bo.BaseUserInfoBo); ok{
				*userNoticeInfos = append(*userNoticeInfos, bo.UserNoticeInfoBo{
					UserId:    baseUserInfo.UserId,
					OutUserId: baseUserInfo.OutUserId,
					Name:      baseUserInfo.Name,
				})
			}
		}
	}
	return userNoticeInfos
}


func ConvertUserIdsToUserInfoWithSet(orgId int64, operatorId int64, memberSet set.Interface, sourceChannel string, rangeType int, targetType int) *[]bo.UserNoticeInfoBo {
	userIds := make([]int64, 0)
	for _, member := range memberSet.List() {
		if userId, ok := member.(int64); ok{
			userIds = append(userIds, userId)
		}
	}
	return ConvertUserIdsToUserInfo(orgId, operatorId, userIds, sourceChannel, rangeType, targetType)
}

//
//
//func PushDingTalkNotice(typ string, corpId string, orgId, operatorId int64, ownerIds *[]int64, followerIds *[]int64, participantIds *[]int64, content interface{}) {
//	suiteTicket, err := dingtalk.GetSuiteTicket()
//	if err != nil {
//		log.Error("获取SuiteTicket时发生异常：", err)
//		return
//	}
//	client, err := dingtalk.GetDingTalkClient(corpId, suiteTicket)
//	if err != nil {
//		log.Error("获取dingtalk client时发生异常", err)
//		return
//	}
//
//	userIdsList := []*[]int64{ownerIds, followerIds, participantIds}
//
//	noticeUserList := list.New()
//	for i, userIds := range userIdsList {
//		if userIds != nil {
//			for _, userId := range *userIds {
//				if userId == operatorId {
//					continue
//				}
//				userConfig, err := proxies.GetUserConfigInfo(orgId, userId)
//				if err != nil {
//					log.Errorf("获取%d用户配置失败, %v", userId, err)
//					continue
//				}
//				if i == 1 && userConfig.OwnerRangeStatus != 1 {
//					continue
//				} else if i == 2 && userConfig.ParticipantRangeStatus != 1 {
//					continue
//				} else if i == 3 && userConfig.AttentionRangeStatus != 1 {
//					continue
//				}
//				if typ == consts.PushIssueUpdate && userConfig.ModifyMessageStatus == 1 {
//					noticeUserList.PushBack(userId)
//				} else if typ == consts.PushIssueRemind && userConfig.RemindMessageStatus == 1 {
//					noticeUserList.PushBack(userId)
//				} else if typ == consts.PushIssueCommentAndAt && userConfig.CommentAtMessageStatus == 1 {
//					noticeUserList.PushBack(userId)
//				} else if typ == consts.PushRelatedContentDynamics && userConfig.RelationMessageStatus == 1 {
//					noticeUserList.PushBack(userId)
//				}
//			}
//		}
//	}
//
//	bePushedUserIds := make([]int64, noticeUserList.Len())
//
//	i := 0
//	for e := noticeUserList.Front(); e != nil; e = e.Next() {
//		bePushedUserIds[i] = e.Value.(int64)
//		i++
//	}
//
//	bePushedUserIds = slice.SliceUniqueInt64(bePushedUserIds)
//
//	var msg *sdk.WorkNoticeMsg = nil
//
//	if typ == consts.PushIssueUpdate {
//		msg, err = GetIssueUpdateNoticeMsg(content)
//		if err != nil {
//			log.Error(err)
//			return
//		}
//	} else if typ == consts.PushIssueRemind {
//	} else if typ == consts.PushIssueCommentAndAt {
//	} else if typ == consts.PushRelatedContentDynamics {
//	}
//
//	if msg != nil {
//		NoticePushToDingTalk(*client, orgId, bePushedUserIds, *msg)
//	}
//}
//
//func GetIssueUpdateNoticeMsg(content interface{}) (*sdk.WorkNoticeMsg, error) {
//	issue, ok := content.(bo.NoticeUpdateIssueBo)
//
//	if !ok {
//		log.Error("NoticePushCreateIssue content 不是 NoticeUpdateIssueBo 类型")
//		return nil, errors.New("NoticePushCreateIssue content 不是 NoticeUpdateIssueBo 类型")
//	}
//
//	author := " "
//
//	msg := sdk.WorkNoticeMsg{
//		MsgType: "oa",
//		OA: &sdk.OANotice{
//			MsgUrl: "http://study.ikuvn.com",
//			Head: sdk.OANoticeHead{
//				BgColor: "00CCFF",
//				Text:    "Polaris",
//			},
//			Body: sdk.OANoticeBody{
//				Title: &issue.NoticeTitle,
//				Form: &[]sdk.OANoticeBodyForm{
//					{
//						Key:   "任务标题: ",
//						Value: issue.IssueTitle,
//					}, {
//						Key:   "操作人: ",
//						Value: issue.UpdateUserName,
//					}, {
//						Key:   "状态: ",
//						Value: issue.StatusName,
//					},
//				},
//				Content: &issue.NoticeContent,
//				Author:  &author,
//			},
//		},
//	}
//	log.Info(issue)
//	str, _ := json.ToJson(msg)
//	log.Infof("要推送的信息为：%s", str)
//	return &msg, nil
//}
//
