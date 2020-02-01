package domain

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/trendsvo"
	"github.com/galaxy-book/polaris-backend/facade/trendsfacade"
)

func CreateIssueComment(issueBo bo.IssueBo, comment string, mentionedUserIds []int64, operatorId int64) (int64, errs.SystemErrorInfo) {
	orgId := issueBo.OrgId
	projectId := issueBo.ProjectId

	//拼装评论
	commentBo := bo.CommentBo{
		OrgId:      orgId,
		ProjectId:  projectId,
		ObjectId:   issueBo.Id,
		ObjectType: consts.TrendsOperObjectTypeIssue,
		Content:    comment,
		Creator:    operatorId,
		Updator:    operatorId,
		IsDelete:   consts.AppIsNoDelete,
	}

	respVo := trendsfacade.CreateComment(trendsvo.CreateCommentReqVo{CommentBo: commentBo})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return 0, respVo.Error()
	}

	commentId := respVo.CommentId
	commentBo.Id = commentId




	asyn.Execute(func() {
		issueMembersBo, err := GetIssueMembers(orgId, issueBo.Id)
		if err != nil{
			log.Error(err)
			return
		}

		beforeParticipantIds := issueMembersBo.ParticipantIds
		beforeFollowerIds := issueMembersBo.FollowerIds

		issueTrendsBo := bo.IssueTrendsBo{
			PushType:                 consts.PushTypeIssueComment,
			OrgId:                    issueBo.OrgId,
			IssueId:                  issueBo.Id,
			ParentIssueId:            issueBo.ParentId,
			ProjectId:                issueBo.ProjectId,
			PriorityId:				  issueBo.PriorityId,
			ParentId:				  issueBo.ParentId,

			OperatorId:               operatorId,

			IssueTitle:               issueBo.Title,
			IssueStatusId:            issueBo.Status,
			BeforeOwner:              issueBo.Owner,
			AfterOwner:               issueBo.Owner,
			BeforeChangeFollowers: beforeFollowerIds,
			AfterChangeFollowers: beforeFollowerIds,
			BeforeChangeParticipants: beforeParticipantIds,
			AfterChangeParticipants: beforeParticipantIds,

			Ext: bo.TrendExtensionBo{
				MentionedUserIds: mentionedUserIds,
				CommentBo: commentBo,
			},
		}

		asyn.Execute(func() {
			PushIssueTrends(issueTrendsBo)
		})

		asyn.Execute(func() {
			PushIssueThirdPlatformNotice(issueTrendsBo)
		})
	})

	return commentId, nil
}
