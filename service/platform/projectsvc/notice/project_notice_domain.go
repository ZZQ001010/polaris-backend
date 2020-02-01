package notice

import (
	"container/list"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/extra/dingtalk"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/processvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
	"gopkg.in/fatih/set.v0"
	"strings"
)

func ProjectMemberChangeNotice(projectMemberChangeBo bo.ProjectMemberChangeBo) {
	pushType := projectMemberChangeBo.PushType
	switch pushType {
	case consts.PushTypeUpdateProjectMembers:
		pushProjectMemberChangeNotice(
			projectMemberChangeBo.OrgId,
			projectMemberChangeBo.ProjectId,
			projectMemberChangeBo.OperatorId,
			projectMemberChangeBo.BeforeChangeMembers,
			projectMemberChangeBo.AfterChangeMembers)
	case consts.PushTypeUpdateProjectStatus:
		pushProjectStatusModifyNotice(projectMemberChangeBo)
	}
}

func parserProjectStatusName(projectMemberChangeBo bo.ProjectMemberChangeBo) (newStatusName string, oldStatusName string, ok bool) {
	orgId := projectMemberChangeBo.OrgId

	oldValue := &map[string]interface{}{}
	newValue := &map[string]interface{}{}
	err := json.FromJson(projectMemberChangeBo.OldValue, oldValue)
	if err != nil {
		log.Error(err)
		return
	}
	err = json.FromJson(projectMemberChangeBo.NewValue, newValue)
	if err != nil {
		log.Error(err)
		return
	}
	propertyName := projectMemberChangeBo.OperateObjProperty
	oldStatusId, ok := (*oldValue)[propertyName].(int64)
	if !ok {
		log.Error("parser old status err")
	}
	newStatusId, ok := (*newValue)[propertyName].(int64)
	if !ok {
		log.Error("parser new status err")
	}
	oldStatusInfoRespVo := processfacade.GetProcessStatus(processvo.GetProcessStatusReqVo{
		OrgId: orgId,
		Id:    oldStatusId,
	})
	if oldStatusInfoRespVo.Failure() {
		log.Error(oldStatusInfoRespVo.Message)
		return
	}
	oldStatusInfo := oldStatusInfoRespVo.CacheProcessStatusBo

	newStatusInfoRespVo := processfacade.GetProcessStatus(processvo.GetProcessStatusReqVo{
		OrgId: orgId,
		Id:    newStatusId,
	})
	if newStatusInfoRespVo.Failure() {
		log.Error(newStatusInfoRespVo.Message)
		return
	}
	newStatusInfo := newStatusInfoRespVo.CacheProcessStatusBo
	return oldStatusInfo.Name, newStatusInfo.Name, true
}

func pushProjectStatusModifyNotice(projectMemberChangeBo bo.ProjectMemberChangeBo) {
	orgId := projectMemberChangeBo.OrgId
	projectId := projectMemberChangeBo.ProjectId
	operatorId := projectMemberChangeBo.OperatorId

	oldStatusName, newStatusName, ok := parserProjectStatusName(projectMemberChangeBo)
	if !ok {
		return
	}

	projectAuthBo, err := domain.LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return
	}
	userInfos := GetProjectNormalUserIds(*projectAuthBo, orgId, operatorId, projectMemberChangeBo.PushType)

	if len(userInfos) == 0 {
		log.Info("push user count is zero...")
		return
	}

	pushUserIds := make([]string, 0)
	for _, userInfo := range userInfos {
		pushUserIds = append(pushUserIds, userInfo.OutUserId)
	}

	operatorBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(consts.AppSourceChannelDingTalk, orgId, operatorId)
	if err != nil {
		log.Error(err)
		return
	}
	author := " "
	noticeTitle := "项目状态由 " + oldStatusName + " 变更为 " + newStatusName
	msg := sdk.WorkNoticeMsg{
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
						Key:   "项目名称: ",
						Value: projectAuthBo.Name,
					}, {
						Key:   "操作人: ",
						Value: operatorBaseInfo.Name,
					},
				},
				Author: &author,
			},
		},
	}

	pushDingTalkMsg(orgId, pushUserIds, msg)
}

func pushProjectMemberChangeNotice(orgId, projectId int64, operatorId int64, beforeChangeMembers []int64, afterChangeMembers []int64) {

	beforeChangeMembersSet := set.New(set.ThreadSafe)
	for _, member := range beforeChangeMembers {
		beforeChangeMembersSet.Add(member)
	}
	afterChangeMembersSet := set.New(set.ThreadSafe)
	for _, member := range afterChangeMembers {
		afterChangeMembersSet.Add(member)
	}

	deletedMembersSet := set.Difference(beforeChangeMembersSet, afterChangeMembersSet)
	addedMembersSet := set.Difference(afterChangeMembersSet, beforeChangeMembersSet)

	if deletedMembersSet.Size() > 0 {
		pushMemberMsg(orgId, projectId, operatorId, deletedMembersSet, 1)
	}
	if addedMembersSet.Size() > 0 {
		pushMemberMsg(orgId, projectId, operatorId, addedMembersSet, 2)
	}
}

func pushMemberMsg(orgId int64, projectId int64, operatorId int64, membersSet set.Interface, typ int) {
	pushUserIds := &[]string{}
	pushUserNames := &[]string{}

	projectBo, err := domain.LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return
	}

	operatorBaseInfo, err := orgfacade.GetBaseUserInfoRelaxed(consts.AppSourceChannelDingTalk, orgId, operatorId)
	if err != nil {
		log.Error(err)
		return
	}

	//处理pushUserId和pushUserName
	appendUserInfo(orgId, operatorId, membersSet, pushUserIds, pushUserNames)

	if len(*pushUserIds) == 0 {
		log.Info("推送用户数组长度为1， 不需要分享")
		return
	}

	action := ""
	if typ == 1 {
		action = "移除了"
	} else {
		action = "添加了"
	}

	noticeTitle := operatorBaseInfo.Name + " " + action + "参与者 " + strings.Join(*pushUserNames, ",")
	author := " "

	msg := sdk.WorkNoticeMsg{
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
						Key:   "项目名称: ",
						Value: projectBo.Name,
					}, {
						Key:   "操作人: ",
						Value: operatorBaseInfo.Name,
					},
				},
				Author: &author,
			},
		},
	}
	pushDingTalkMsg(orgId, *pushUserIds, msg)
}

func appendUserInfo(orgId, operatorId int64, membersSet set.Interface, pushUserIds, pushUserNames *[]string) {

	for _, member := range membersSet.List() {
		userId := member.(int64)

		if userId == operatorId {
			continue
		}

		userConfig, err := orgfacade.GetUserConfigInfoRelaxed(orgId, userId)
		if err != nil {
			log.Errorf("获取%d用户配置失败, %v", userId, err)
			continue
		}
		if userConfig.ParticipantRangeStatus != 1 {
			continue
		}
		if userConfig.RelationMessageStatus != 1 {
			continue
		}

		baseUserInfo, baseUserInfoErr := orgfacade.GetBaseUserInfoRelaxed(consts.AppSourceChannelDingTalk, orgId, userId)
		if baseUserInfoErr != nil {
			log.Error(baseUserInfoErr)
			continue
		}

		*pushUserIds = append(*pushUserIds, baseUserInfo.OutUserId)
		*pushUserNames = append(*pushUserNames, baseUserInfo.Name)
	}
}

func pushDingTalkMsg(orgId int64, pushUserIds []string, msg sdk.WorkNoticeMsg) {
	pushUserIdsStr := strings.Join(pushUserIds, ",")

	orgBaseInfo, err := orgfacade.GetBaseOrgInfoRelaxed(consts.AppSourceChannelDingTalk, orgId)
	if err != nil {
		log.Errorf("组织外部信息不存在 %v", err)
		return
	}
	client, err2 := dingtalk.GetDingTalkClientRest(orgBaseInfo.OutOrgId)
	if err2 != nil {
		log.Errorf("获取dingtalk client时发生异常 %v", err2)
		return
	}
	resp, err1 := client.SendWorkNotice(&pushUserIdsStr, nil, false, msg)
	if err1 != nil {
		log.Error("发送ding talk 工作通知时发生异常" + strs.ObjectToString(err1))
		return
	}
	if resp.ErrCode != 0 {
		log.Error("发送ding talk 失败" + resp.ErrMsg)
		return
	}
	log.Infof("推送项目通知成功，code %d, msg: %s", resp.ErrCode, resp.ErrMsg)
}

func GetProjectNormalUserIds(projectAuthBo bo.ProjectAuthBo, orgId, operatorId int64, pushType consts.IssueNoticePushType) []bo.UserNoticeInfoBo {
	ownerArray := []int64{projectAuthBo.Owner}
	userIdsList := [][]int64{ownerArray, projectAuthBo.Participants, projectAuthBo.Followers}

	//返回链表
	noticeUserList := dealProjectNoticeUserIdsList(userIdsList, operatorId, orgId, pushType)

	//转换array并去重
	bePushedUserIds := make([]int64, noticeUserList.Len())
	i := 0
	for e := noticeUserList.Front(); e != nil; e = e.Next() {
		bePushedUserIds[i] = e.Value.(int64)
		i++
	}
	bePushedUserIds = slice.SliceUniqueInt64(bePushedUserIds)

	userNoticeInfos := make([]bo.UserNoticeInfoBo, len(bePushedUserIds))

	if len(bePushedUserIds) > 0 {
		for i, userId := range bePushedUserIds {
			baseUserInfo, err := orgfacade.GetBaseUserInfoRelaxed(consts.AppSourceChannelDingTalk, orgId, userId)
			if err != nil {
				log.Error(err)
				continue
			}
			userNoticeInfos[i] = bo.UserNoticeInfoBo{
				UserId:    baseUserInfo.UserId,
				OutUserId: baseUserInfo.OutUserId,
				Name:      baseUserInfo.Name,
			}
		}
	}
	return userNoticeInfos
}

func dealProjectNoticeUserIdsList(userIdsList [][]int64, operatorId int64, orgId int64, pushType consts.IssueNoticePushType) *list.List {

	noticeUserList := list.New()

	for i, userIds := range userIdsList {
		if userIds != nil {
			dealProjectNoticeUserIds(i, userIds, operatorId, orgId, pushType, noticeUserList)
		}
	}
	return noticeUserList
}

func dealProjectNoticeUserIds(i int, userIds []int64, operatorId int64, orgId int64, pushType consts.IssueNoticePushType, noticeUserList *list.List) {

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
		if (pushType == consts.PushTypeCreateProject || pushType == consts.PushTypeUpdateProject || pushType == consts.PushTypeUpdateProjectStatus) && userConfig.ModifyMessageStatus == 1 {
			noticeUserList.PushBack(userId)
		}
	}
}

