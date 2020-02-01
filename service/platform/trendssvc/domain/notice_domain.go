package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/extra/mqtt"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/trendssvc/po"
	"upper.io/db.v3"
)

func UnreadNoticeCount(orgId, userId int64) (uint64, errs.SystemErrorInfo) {
	count, err := mysql.SelectCountByCond(consts.TableNotice, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcNoticer:  userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.NoticeUnReadStatus,
	})
	if err != nil {
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return count, nil
}

func GetNoticeList(orgId, userId int64, page, size int, input *vo.NoticeListReq) (uint64, *[]bo.NoticeBo, errs.SystemErrorInfo) {
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcNoticer:  userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if input != nil && input.Type != nil {
		cond[consts.TcType] = input.Type
	}
	count, err := mysql.SelectCountByCond(consts.TableNotice, cond)
	if err != nil {
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	noticePo := &[]po.PpmTreNotice{}
	err = mysql.SelectAllByCondWithNumAndOrder(consts.TableNotice, cond, nil, page, size, "create_time desc", noticePo)
	if err != nil {
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	noticeBo := &[]bo.NoticeBo{}
	copyErr := copyer.Copy(noticePo, noticeBo)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return count, noticeBo, nil
}

func AddNoticeByChangeProjectMember(delMembers, addMembers []int64, projectTrendsBo bo.ProjectTrendsBo) errs.SystemErrorInfo {
	noticePo := []po.PpmTreNotice{}
	if len(delMembers) >= 0 {
		for _, v := range delMembers {
			noticePo = append(noticePo, po.PpmTreNotice{
				OrgId:     projectTrendsBo.OrgId,
				ProjectId: projectTrendsBo.ProjectId,
				Type:      1, //项目通知
				Content:   "将您移出了项目【" + projectTrendsBo.Ext.ObjName + "】",
				RelationType:consts.TrendsRelationTypeDeletedProjectParticipant,
				Noticer:   v,
				Creator:   projectTrendsBo.OperatorId,
			})
		}
	}
	if len(addMembers) >= 0 {
		for _, v := range addMembers {
			noticePo = append(noticePo, po.PpmTreNotice{
				OrgId:     projectTrendsBo.OrgId,
				ProjectId: projectTrendsBo.ProjectId,
				Type:      1, //项目通知
				Content:   "将您加入了项目【" + projectTrendsBo.Ext.ObjName + "】",
				RelationType:consts.TrendsRelationTypeAddedProjectParticipant,
				Noticer:   v,
				Creator:   projectTrendsBo.OperatorId,
			})
		}
	}
	err := InsertNotice(noticePo)
	if err != nil {
		return err
	}

	return nil
}

func InsertNotice(noticePo []po.PpmTreNotice) errs.SystemErrorInfo {
	if len(noticePo) == 0 {
		return nil
	}
	resp, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableNotice, len(noticePo))
	if err != nil {
		log.Error(err)
		return err
	}
	for k, _ := range noticePo {
		if len(resp.Ids) >= k {
			noticePo[k].Id = resp.Ids[k].Id
		}
	}

	creator := noticePo[0].Creator
	orgId := noticePo[0].OrgId
	creatorInfo, err := orgfacade.GetBaseUserInfoRelaxed("", orgId, creator)
	if err != nil{
		log.Error(err)
		return err
	}

	insertErr := mysql.BatchInsert(&po.PpmTreNotice{}, slice.ToSlice(noticePo))
	if insertErr != nil {
		log.Error(insertErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertErr)
	}

	asyn.Execute(func() {

		for _, noticeInfo := range noticePo {
			notice := noticeInfo
			asyn.Execute(func() {
				channel := util.GetMQTTUserChannel(notice.OrgId, notice.Noticer)
				_ = mqtt.Publish(channel, json.ToJsonIgnoreError(bo.MQTTNoticeBo{
					Type: consts.MQTTNoticeTypeRemind,
					Body: bo.MQTTRemindNotice{
						OperatorId: creator,
						OperatorName: creatorInfo.Name,
						OrgID: orgId,
						Content: notice.Content,
						Data: notice,
					},
				}))
			})
		}
	})

	return nil
}

func AddNoticeByChangeProjectNotice(projectTrendsBo bo.ProjectTrendsBo) errs.SystemErrorInfo {
	//获取项目相关人员
	relationResp := projectfacade.GetProjectRelation(projectvo.GetProjectRelationReqVo{
		ProjectId:    projectTrendsBo.ProjectId,
		RelationType: []int64{consts.IssueRelationTypeOwner, consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower},
	})
	if relationResp.Failure() {
		log.Error(relationResp.Error())
		return relationResp.Error()
	}
	if len(relationResp.Data) == 0 {
		return nil
	}
	noticePo := []po.PpmTreNotice{}
	for _, v := range relationResp.Data {
		noticePo = append(noticePo, po.PpmTreNotice{
			OrgId:     projectTrendsBo.OrgId,
			ProjectId: projectTrendsBo.ProjectId,
			Type:      1, //项目通知
			RelationType:    consts.TrendsRelationTypeUpdateProject,
			Content:   "更新了项目公告【" + projectTrendsBo.Ext.ObjName + "】",
			Noticer:   v.RelationId,
			Creator:   projectTrendsBo.OperatorId,
		})
	}

	err := InsertNotice(noticePo)
	if err != nil {
		return err
	}

	return nil
}

func AddNoticeByChangeIssueOwner(issueTrendsBo bo.IssueTrendsBo) errs.SystemErrorInfo {
	noticePo := []po.PpmTreNotice{}
	noticePo = append(noticePo, po.PpmTreNotice{
		OrgId:     issueTrendsBo.OrgId,
		ProjectId: issueTrendsBo.ProjectId,
		IssueId:   issueTrendsBo.IssueId,
		Type:      1, //项目通知
		RelationType:    consts.TrendsRelationTypeUpdateIssueOwner,
		Content:   "将任务指派给你【" + issueTrendsBo.IssueTitle + "】",
		Noticer:   issueTrendsBo.AfterOwner,
		Creator:   issueTrendsBo.OperatorId,
	})

	err := InsertNotice(noticePo)
	if err != nil {
		return err
	}

	return nil
}

func issueRelateIds(issueTrendsBo bo.IssueTrendsBo) ([]int64, errs.SystemErrorInfo) {
	relatedIds := make([]int64, 0)
	//获取任务相关人员
	resp := projectfacade.GetIssueMembers(projectvo.GetIssueMembersReqVo{
		OrgId:   issueTrendsBo.OrgId,
		IssueId: issueTrendsBo.IssueId,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return relatedIds, resp.Error()
	}
	respData := resp.Data
	return respData.MemberIds, nil
}

func AddNoticeByChangeUpdateIssue(issueTrendsBo bo.IssueTrendsBo) errs.SystemErrorInfo {
	relatedIds, err1 := issueRelateIds(issueTrendsBo)
	if err1 != nil {
		return err1
	}
	noticePo := []po.PpmTreNotice{}
	for _, v := range relatedIds {
		noticePo = append(noticePo, po.PpmTreNotice{
			OrgId:     issueTrendsBo.OrgId,
			ProjectId: issueTrendsBo.ProjectId,
			IssueId:   issueTrendsBo.IssueId,
			Type:      1, //项目通知
			RelationType:    consts.TrendsRelationTypeUpdateIssue,
			Content:   "更新了任务详情【" + issueTrendsBo.IssueTitle + "】",
			Noticer:   v,
			Creator:   issueTrendsBo.OperatorId,
		})
	}

	err := InsertNotice(noticePo)
	if err != nil {
		return err
	}

	return nil
}

func AddNoticeByCreateComment(issueTrendsBo bo.IssueTrendsBo) errs.SystemErrorInfo {
	trendsExt := issueTrendsBo.Ext
	mentionedUserIds := trendsExt.MentionedUserIds
	ext := issueTrendsBo.Ext
	commentBo := ext.CommentBo
	operatorId := issueTrendsBo.OperatorId

	mentionedUserMap := map[int64]bool{}
	if mentionedUserIds != nil{
		for _, mentionUserId := range mentionedUserIds{
			mentionedUserMap[mentionUserId] = true
		}
	}

	relatedIds, err1 := issueRelateIds(issueTrendsBo)
	if err1 != nil {
		log.Error(err1)
		return err1
	}
	noticePo := []po.PpmTreNotice{}
	for _, v := range relatedIds {
		if v == operatorId{
			continue
		}
		if _, ok := mentionedUserMap[v]; ok{
			noticePo = append(noticePo, po.PpmTreNotice{
				OrgId:     issueTrendsBo.OrgId,
				ProjectId: issueTrendsBo.ProjectId,
				IssueId:   issueTrendsBo.IssueId,
				Type:      1, //项目通知
				RelationType:    consts.TrendsRelationTypeCreateIssueComment,
				Content:   "评论了任务【" + commentBo.Content + "】",
				Noticer:   v,
				Creator:   operatorId,
			})
		}else{
			noticePo = append(noticePo, po.PpmTreNotice{
				OrgId:     issueTrendsBo.OrgId,
				ProjectId: issueTrendsBo.ProjectId,
				IssueId:   issueTrendsBo.IssueId,
				Type:      1, //项目通知
				RelationType:    consts.TrendsRelationTypeCreateIssueComment,
				Content:   "评论了任务【" + issueTrendsBo.IssueTitle + "】",
				Noticer:   v,
				Creator:   operatorId,
			})
		}
	}

	err := InsertNotice(noticePo)
	if err != nil {
		return err
	}

	return nil
}

func AddNoticeByChangeUpdateIssueStatus(issueTrendsBo bo.IssueTrendsBo) errs.SystemErrorInfo {
	relatedIds, err1 := issueRelateIds(issueTrendsBo)
	if err1 != nil {
		log.Error(err1)
		return err1
	}
	noticePo := []po.PpmTreNotice{}
	for _, v := range relatedIds {
		noticePo = append(noticePo, po.PpmTreNotice{
			OrgId:     issueTrendsBo.OrgId,
			ProjectId: issueTrendsBo.ProjectId,
			IssueId:   issueTrendsBo.IssueId,
			Type:      1, //项目通知
			RelationType:    consts.TrendsRelationTypeUpdateIssueStatus,
			Content:   "完成了任务【" + issueTrendsBo.IssueTitle + "】",
			Noticer:   v,
			Creator:   issueTrendsBo.OperatorId,
		})
	}

	err := InsertNotice(noticePo)
	if err != nil {
		return err
	}

	return nil
}

func AddNoticeByApplyJoinOrg(orgTrendsBo bo.OrgTrendsBo, userInfos []bo.BaseUserInfoBo) errs.SystemErrorInfo{
	orgId := orgTrendsBo.OrgId
	operatorId := orgTrendsBo.OperatorId
	respVo := rolefacade.GetOrgAdminUser(rolevo.GetOrgAdminUserReqVo{
		OrgId: orgId,
	})
	if respVo.Failure(){
		log.Error(respVo.Message)
		return respVo.Error()
	}

	adminUserIds := respVo.Data
	if len(adminUserIds) == 0{
		return nil
	}

	userNames := ""
	//整理被加入的成员列表
	for _, userInfo := range userInfos{
		userNames += userInfo.Name + ", "
	}
	userNames = userNames[0 : len(userNames) - 1]

	noticePos := make([]po.PpmTreNotice, 0)
	for _, userId := range adminUserIds{
		noticePos = append(noticePos, po.PpmTreNotice{
			OrgId:     orgId,
			Type:      2, //项目通知
			RelationType:    consts.TrendsRelationTypeApplyJoinOrg,
			Content:   userNames + "申请加入组织",
			Noticer:   userId,
			Creator:   operatorId,
		})
	}

	err := InsertNotice(noticePos)
	if err != nil {
		return err
	}

	return nil
}