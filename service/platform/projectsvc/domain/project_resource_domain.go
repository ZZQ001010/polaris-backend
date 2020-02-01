package domain

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/facade/resourcefacade"
)

func CreateProjectResource(createResourceReqBo bo.ProjectCreateResourceReqBo) (int64, errs.SystemErrorInfo) {
	var resourceId int64 = 0

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
	sourcetype := consts.OssPolicyTypeProjectResource
	createResourceBo := bo.CreateResourceBo{
		ProjectId:  createResourceReqBo.ProjectId,
		OrgId:      createResourceReqBo.OrgId,
		Path:       createResourceReqBo.ResourcePath,
		Name:       createResourceReqBo.FileName,
		Suffix:     createResourceReqBo.FileSuffix,
		Bucket:     bucketName,
		Type:       resourceType,
		Size:       createResourceReqBo.ResourceSize,
		Md5:        md5,
		OperatorId: createResourceReqBo.OperatorId,
		FolderId:   &createResourceReqBo.FolderId,
		SourceType: &sourcetype,
	}

	respVo := resourcefacade.CreateResource(resourcevo.CreateResourceReqVo{CreateResourceBo: createResourceBo})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return 0, respVo.Error()
	}

	asyn.Execute(func() {
		//新增动态
		ext := bo.TrendExtensionBo{ObjName: createResourceBo.Name + createResourceBo.Suffix}
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeCreateProjectFile,
			OrgId:      createResourceBo.OrgId,
			ProjectId:  createResourceBo.ProjectId,
			OperatorId: createResourceBo.OperatorId,
			Ext:        ext,
		})
	})

	resourceId = respVo.ResourceId
	return resourceId, nil
}

func UpdateProjectResourceName(updateResourceBo bo.UpdateResourceInfoBo) (int64, errs.SystemErrorInfo) {
	respVo := resourcefacade.UpdateResourceInfo(resourcevo.UpdateResourceInfoReqVo{
		Input: updateResourceBo,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return 0, respVo.Error()
	}
	//新增动态
	changes := bo.TrendChangeListBo{
		Field:     "resourceName",
		FieldName: consts.ProjectResourceName,
		OldValue:  respVo.OldBo[0].Name,
		NewValue:  respVo.NewBo[0].Name,
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{ChangeList: []bo.TrendChangeListBo{
			changes,
		},
			ObjName: json.ToJsonIgnoreError(respVo.OldBo[0].Name),
		}
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeUpdateProjectFile,
			OrgId:      updateResourceBo.OrgId,
			ProjectId:  updateResourceBo.ProjectId,
			OperatorId: updateResourceBo.UserId,
			Ext:        ext,
		})
	})

	resourceId := updateResourceBo.ResourceId
	//time.Sleep(10 * time.Second)
	return resourceId, nil
}

func UpdateProjectResourceFolder(updateResourceBo bo.UpdateResourceFolderBo) ([]int64, errs.SystemErrorInfo) {
	respVo := resourcefacade.UpdateResourceFolder(resourcevo.UpdateResourceFolderReqVo{
		Input: updateResourceBo,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	//新增动态
	changes := bo.TrendChangeListBo{
		Field:     "resourceFolder",
		FieldName: consts.ProjectResourceFolder,
		OldValue:  *respVo.CurrentFolderName,
		NewValue:  *respVo.TargetFolderName,
	}
	resourceNames := []string{}
	resourceIds := []int64{}
	for _, value := range respVo.OldBo {
		resourceNames = append(resourceNames, value.Name)
		resourceIds = append(resourceIds, value.Id)
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{
			ChangeList: []bo.TrendChangeListBo{
				changes,
			},
			ObjName: json.ToJsonIgnoreError(resourceNames),
		}
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeUpdateProjectFile,
			OrgId:      updateResourceBo.OrgId,
			ProjectId:  updateResourceBo.ProjectId,
			OperatorId: updateResourceBo.UserId,
			Ext:        ext,
		})
	})
	return resourceIds, nil
}

func DeleteProjectResource(deleteBo bo.DeleteResourceBo) ([]int64, errs.SystemErrorInfo) {
	respVo := resourcefacade.DeleteResource(resourcevo.DeleteResourceReqVo{
		Input: deleteBo,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}

	resourceNames := []string{}
	resourceIds := []int64{}
	for _, value := range respVo.OldBo {
		resourceNames = append(resourceNames, value.Name)
		resourceIds = append(resourceIds, value.Id)
	}
	asyn.Execute(func(){
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:   consts.PushTypeDeleteProjectFile,
			OrgId:      deleteBo.OrgId,
			ProjectId:  deleteBo.ProjectId,
			OperatorId: deleteBo.UserId,
			NewValue:   json.ToJsonIgnoreError(resourceNames),
		})
	})
	return resourceIds, nil
}

func GetProjectResource(bo bo.GetResourceBo) (*vo.ResourceList, errs.SystemErrorInfo) {
	respVo := resourcefacade.GetResource(resourcevo.GetResourceReqVo{
		Input: bo,
	})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.ResourceList, nil
}
