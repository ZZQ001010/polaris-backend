package domain

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/facade/resourcefacade"
	"time"
)

func CreateIssueResource(issueBo bo.IssueBo, createResourceReqBo bo.IssueCreateResourceReqBo) (int64, errs.SystemErrorInfo) {
	var resourceId int64 = 0
	operatorId := createResourceReqBo.OperatorId

	bucketName := ""
	md5 := ""
	if createResourceReqBo.BucketName != nil {
		bucketName = *createResourceReqBo.BucketName
	}
	if createResourceReqBo.Md5 != nil {
		md5 = *createResourceReqBo.Md5
	}

	resourceType := consts.OssResource
	//这边认为，bucketName不传就是本地上传
	if bucketName == "" {
		resourceType = consts.LocalResource
	}
	//20191231 修改为sourceType 原filetype代表文件类型
	//fileType := consts.OssPolicyTypeIssueResource
	sourceType := consts.OssPolicyTypeIssueResource
	createResourceBo := bo.CreateResourceBo{
		ProjectId:  issueBo.ProjectId,
		OrgId:      issueBo.OrgId,
		Path:       createResourceReqBo.ResourcePath,
		Name:       createResourceReqBo.FileName,
		Suffix:     createResourceReqBo.FileSuffix,
		Bucket:     bucketName,
		Type:       resourceType,
		Size:       createResourceReqBo.ResourceSize,
		Md5:        md5,
		OperatorId: createResourceReqBo.OperatorId,
		SourceType: &sourceType,
	}

	respVo := resourcefacade.CreateResource(resourcevo.CreateResourceReqVo{CreateResourceBo: createResourceBo})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return 0, respVo.Error()
	}
	resId := respVo.ResourceId

	resourceId = resId

	_, err2 := UpdateIssueRelationSingle(operatorId, issueBo, consts.IssueRelationTypeResource, resId)
	if err2 != nil {
		log.Error(err2)
		return 0, err2
	}

	asyn.Execute(func() {
		resourceTrend := []bo.ResourceInfoBo{}
		resourceTrend = append(resourceTrend, bo.ResourceInfoBo{
			Name:       createResourceReqBo.FileName,
			Url:        createResourceReqBo.ResourcePath,
			Size:       createResourceReqBo.ResourceSize,
			UploadTime: time.Now(),
			Suffix:     createResourceReqBo.FileSuffix,
		})
		issueTrendsBo := bo.IssueTrendsBo{
			PushType:      consts.PushTypeUploadResource,
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
			PushIssueTrends(issueTrendsBo)
		})
		asyn.Execute(func() {
			PushIssueThirdPlatformNotice(issueTrendsBo)
		})
	})

	return resourceId, nil
}

func DeleteIssueResource(issueBo bo.IssueBo, deletedResourceIds []int64, operatorId int64) errs.SystemErrorInfo {
	return DeleteIssueRelationByIds(operatorId, issueBo, consts.IssueRelationTypeResource, deletedResourceIds)
}
