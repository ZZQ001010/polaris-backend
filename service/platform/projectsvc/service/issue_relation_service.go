package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/facade/resourcefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"time"
)

func CreateIssueComment(orgId, currentUserId int64, input vo.CreateIssueCommentReq) (*vo.Void, errs.SystemErrorInfo) {
	issueId := input.IssueID

	issueBo, err1 := domain.GetIssueBo(orgId, issueId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	err := domain.AuthIssue(orgId, currentUserId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationComment)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	commentId, err2 := domain.CreateIssueComment(*issueBo, input.Comment, input.MentionedUserIds, currentUserId)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
	}

	return &vo.Void{
		ID: commentId,
	}, nil
}

func CreateIssueResource(orgId, operatorId int64, input vo.CreateIssueResourceReq) (*vo.Void, errs.SystemErrorInfo) {
	issueId := input.IssueID

	issueBo, err := domain.GetIssueBo(orgId, issueId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIssue)
	}

	err = domain.AuthIssue(orgId, operatorId, *issueBo, consts.RoleOperationPathOrgProAttachment, consts.RoleOperationUpload)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	resourceId, err2 := domain.CreateIssueResource(*issueBo, bo.IssueCreateResourceReqBo{
		ResourcePath: input.ResourcePath,
		ResourceSize: input.ResourceSize,
		FileName:     input.FileName,
		FileSuffix:   input.FileSuffix,
		Md5:          input.Md5,
		BucketName:   input.BucketName,
		OperatorId:   operatorId,
	})

	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
	}

	return &vo.Void{
		ID: resourceId,
	}, nil
}

func CreateIssueRelationIssue(orgId, operatorId int64, input vo.UpdateIssueAndIssueRelateReq) (*vo.Void, errs.SystemErrorInfo) {
	issueId := input.IssueID

	issueBo, err := domain.GetIssueBo(orgId, issueId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIssue)
	}

	err = domain.AuthIssue(orgId, operatorId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	//真正增加的关联问题id
	realAddRelationIds := []int64{}
	if len(input.AddRelateIssueIds) > 0 {
		err = domain.VerifyRelationIssue(input.AddRelateIssueIds, issueBo.ProjectObjectTypeId, orgId)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
		}
		relationInfo, err1 := domain.UpdateIssueRelation(operatorId, *issueBo, consts.IssueRelationTypeIssue, input.AddRelateIssueIds)
		if err1 != nil {
			log.Error(err1)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
		}

		for _, v := range relationInfo {
			realAddRelationIds = append(realAddRelationIds, v.RelationId)
		}
	}

	if len(input.DelRelateIssueIds) > 0 {
		err1 := domain.DeleteIssueRelationByIds(operatorId, *issueBo, consts.IssueRelationTypeIssue, input.DelRelateIssueIds)
		if err1 != nil {
			log.Error(err1)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
		}
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.IssueType = "T"
		ext.ObjName = issueBo.Title

		issueTrendsBo := bo.IssueTrendsBo{
			PushType:      consts.PushTypeUpdateRelationIssue,
			OrgId:         orgId,
			OperatorId:    operatorId,
			IssueId:       issueBo.Id,
			ParentIssueId: issueBo.ParentId,
			ProjectId:     issueBo.ProjectId,
			IssueTitle:    issueBo.Title,
			IssueStatusId: issueBo.Status,
			BeforeOwner:   issueBo.Owner,
			ParentId:      issueBo.ParentId,

			BindIssues:   realAddRelationIds,
			UnbindIssues: input.DelRelateIssueIds,
			Ext:          ext,
		}

		asyn.Execute(func() {
			domain.PushIssueTrends(issueTrendsBo)
		})
		asyn.Execute(func() {
			domain.PushIssueThirdPlatformNotice(issueTrendsBo)
		})
	})

	return &vo.Void{
		ID: input.IssueID,
	}, nil
}

func DeleteIssueResource(orgId, operatorId int64, input vo.DeleteIssueResourceReq) (*vo.Void, errs.SystemErrorInfo) {

	issueId := input.IssueID

	issueBo, err := domain.GetIssueBo(orgId, issueId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIssue)
	}

	//projectInfo, err := domain.LoadProjectAuthBo(orgId, issueBo.ProjectId)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}

	//去重
	input.DeletedResourceIds = slice.SliceUniqueInt64(input.DeletedResourceIds)

	//这里权限写死:已上传文件，可删除其文件，拥有删除权限的人员为，任务负责人、附件上传者，项目负责人，超级管理员
	hasPermission := true

	authErr := domain.AuthIssue(orgId, operatorId, *issueBo, consts.RoleOperationPathOrgProAttachment, consts.RoleOperationDelete)
	if authErr != nil {
		log.Error(authErr)
		hasPermission = false
		//return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, authErr)
	}

	//err = AuthIssue(orgId, operatorId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationUnbind)
	if !hasPermission {
		log.Error(err)
		//小逻辑：允许创建者删除文件资源
		relationIds, err := domain.GetIssueResourceIdsByCreator(orgId, issueId, input.DeletedResourceIds, operatorId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		//判断要删除的资源是否是这个人创建的
		if relationIds != nil && len(*relationIds) == len(input.DeletedResourceIds) {
			hasPermission = true
		}
	}

	if !hasPermission {
		return nil, errs.BuildSystemErrorInfo(errs.NoOperationPermissions)
	}

	judgeErr := JudgeProjectFiling(orgId, issueBo.ProjectId)
	if judgeErr != nil {
		log.Error(judgeErr)
		return nil, judgeErr
	}

	err1 := domain.DeleteIssueResource(*issueBo, input.DeletedResourceIds, operatorId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
	}

	asyn.Execute(func() {
		bos, total, err1 := resourcefacade.GetResourceBoListRelaxed(0, 0, resourcevo.GetResourceBoListCond{
			ResourceIds: &input.DeletedResourceIds,
			OrgId:       orgId,
		})
		if err1 != nil {
			log.Error(err1)
			return
		}
		if total == 0 {
			log.Error("无更新")
			return
		}
		resourceTrend := []bo.ResourceInfoBo{}
		for _, resourceBo := range *bos {
			resourceTrend = append(resourceTrend, bo.ResourceInfoBo{
				Name:       resourceBo.Name,
				Url:        resourceBo.Host + resourceBo.Path,
				Size:       resourceBo.Size,
				UploadTime: time.Now(),
				Suffix:     resourceBo.Suffix,
			})
		}

		issueTrendsBo := bo.IssueTrendsBo{
			PushType:      consts.PushTypeDeleteResource,
			OrgId:         issueBo.OrgId,
			OperatorId:    operatorId,
			IssueId:       issueBo.Id,
			ParentIssueId: issueBo.ParentId,
			ProjectId:     issueBo.ProjectId,
			PriorityId:    issueBo.PriorityId,
			ParentId:      issueBo.ParentId,

			IssueTitle:    issueBo.Title,
			IssueStatusId: issueBo.Status,

			Ext: bo.TrendExtensionBo{
				ObjName:      issueBo.Title,
				ResourceInfo: resourceTrend,
			},
		}
		asyn.Execute(func() {
			domain.PushIssueTrends(issueTrendsBo)
		})
		asyn.Execute(func() {
			domain.PushIssueThirdPlatformNotice(issueTrendsBo)
		})
	})

	return &vo.Void{
		ID: issueId,
	}, nil
}

func RelatedIssueList(orgId int64, input vo.RelatedIssueListReq) (*vo.IssueRestInfoResp, errs.SystemErrorInfo) {
	issueList, err := domain.RelationIssueList(orgId, input.IssueID)
	if err != nil {
		log.Errorf(" issuedomain.RelationIssueList: %q\n", err)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	issueRestInfos := make([]*vo.IssueRestInfo, len(issueList))

	finishedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
	if err != nil {
		log.Errorf("proxies.GetProcessStatusId: %q\n", err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	for i, issueInfo := range issueList {
		finished, err := slice.Contain(*finishedIds, issueInfo.Status)
		baseUserInfo, err := orgfacade.GetDingTalkBaseUserInfoRelaxed(orgId, issueInfo.Owner)
		if err != nil {
			log.Errorf("proxies.GetDingTalkBaseUserInfo: %q\n", err)
			return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}

		statusInfo, err := domain.GetHomeIssueStatusInfoBo(orgId, issueInfo.Status)
		if err != nil {
			log.Errorf("proxies.GetHomeIssueStatusInfoBo: %q\n", err)
			return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}

		issueRestInfos[i] = &vo.IssueRestInfo{
			ID:          issueInfo.Id,
			Title:       issueInfo.Title,
			OwnerID:     issueInfo.Owner,
			OwnerName:   baseUserInfo.Name,
			OwnerAvatar: baseUserInfo.Avatar,
			Finished:    finished,
			StatusID:    issueInfo.Status,
			StatusName:  statusInfo.Name,
		}
	}
	return &vo.IssueRestInfoResp{
		Total: int64(len(issueList)),
		List:  issueRestInfos,
	}, nil
}

func IssueResources(orgId int64, page uint, size uint, input *vo.GetIssueResourcesReq) (*vo.ResourceList, errs.SystemErrorInfo) {
	issueId := input.IssueID

	_, err := domain.GetIssueBo(orgId, issueId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityIssue)
	}

	resourceIds, err1 := domain.GetIssueRelationIdsByRelateType(orgId, issueId, consts.IssueRelationTypeResource)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError)
	}

	isDelete := consts.AppIsNoDelete
	bos, total, err1 := resourcefacade.GetResourceBoListRelaxed(page, size, resourcevo.GetResourceBoListCond{
		ResourceIds: resourceIds,
		OrgId:       orgId,
		IsDelete:    &isDelete,
	})
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ResourceDomainError, err1)
	}
	if bos != nil && len(*bos) > 0 {
		for _, resBo := range *bos {
			resBo.Path = resBo.Host + resBo.Path
		}
	}

	resultList := &[]*vo.Resource{}
	copyErr := copyer.Copy(bos, resultList)

	creatorIds := make([]int64, 0)
	for _, res := range *resultList {
		creatorIds = append(creatorIds, res.Creator)
	}

	creatorInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed("", orgId, creatorIds)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return nil, err
	}

	creatorMap := maps.NewMap("UserId", creatorInfos)

	for _, res := range *resultList {
		if userInfo, ok := creatorMap[res.Creator]; ok {
			res.CreatorName = userInfo.(bo.BaseUserInfoBo).Name
		}
		//拼接host
		res.Path = util.JointUrl(res.Host, res.Path)
		res.PathCompressed = util.GetCompressedPath(res.Path, res.Type)
	}

	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &vo.ResourceList{
		Total: total,
		List:  *resultList,
	}, nil

}

func GetIssueMembers(orgId int64, issueId int64) (*projectvo.GetIssueMembersRespData, errs.SystemErrorInfo){
	issueMembersBo, err := domain.GetIssueMembers(orgId, issueId)
	if err != nil{
		log.Error(err)
		return nil, err
	}

	respVo := &projectvo.GetIssueMembersRespData{}
	_ = copyer.Copy(issueMembersBo, respVo)
	return respVo, nil
}
