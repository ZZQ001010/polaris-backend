package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
)

func CreateProjectResource(orgId, operatorId int64, input vo.CreateProjectResourceReq) (*vo.Void, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProject(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.RoleOperationUpload)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resourceId, err := domain.CreateProjectResource(bo.ProjectCreateResourceReqBo{
		ProjectId:    input.ProjectID,
		OrgId:        orgId,
		ResourcePath: input.ResourcePath,
		ResourceSize: input.ResourceSize,
		FileName:     input.FileName,
		FileSuffix:   input.FileSuffix,
		Md5:          input.Md5,
		BucketName:   input.BucketName,
		OperatorId:   operatorId,
		FolderId:     input.FolderID,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.Void{
		ID: resourceId,
	}, nil
}

func UpdateProjectResourceName(orgId, operatorId int64, input vo.UpdateProjectResourceNameReq) (*vo.Void, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProject(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.RoleOperationModify)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	updateBo := &bo.UpdateResourceInfoBo{}
	err1 := copyer.Copy(&input, updateBo)
	if err1 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	updateBo.UserId = operatorId
	updateBo.OrgId = orgId
	resourceId, err := domain.UpdateProjectResourceName(*updateBo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.Void{
		ID: resourceId,
	}, nil
}

func UpdateProjectResourceFolder(orgId, operatorId int64, input vo.UpdateProjectResourceFolderReq) (*vo.UpdateProjectResourceFolderResp, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProject(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.RoleOperationModify)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	bo := &bo.UpdateResourceFolderBo{}
	err1 := copyer.Copy(&input, bo)
	if err1 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	bo.UserId = operatorId
	bo.OrgId = orgId
	resourceIds, err := domain.UpdateProjectResourceFolder(*bo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.UpdateProjectResourceFolderResp{
		ResourceIds: resourceIds,
	}, nil
}

func DeleteProjectResource(orgId, operatorId int64, input vo.DeleteProjectResourceReq) (*vo.DeleteProjectResourceResp, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProject(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.RoleOperationDelete)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resourceIds, err := domain.DeleteProjectResource(bo.DeleteResourceBo{
		ResourceIds: input.ResourceIds,
		FolderId:    &input.FolderID,
		OrgId:       orgId,
		UserId:      operatorId,
		ProjectId:   input.ProjectID,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.DeleteProjectResourceResp{
		ResourceIds: resourceIds,
	}, nil
}

func GetProjectResource(orgId, operatorId int64, page, size int, input vo.ProjectResourceReq) (*vo.ResourceList, errs.SystemErrorInfo) {
	//此处的path需要更改
	err := domain.AuthProjectWithOutPermission(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProFile, consts.RoleOperationView)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resourceList, err := domain.GetProjectResource(bo.GetResourceBo{
		FolderId:  &input.FolderID,
		OrgId:     orgId,
		UserId:    operatorId,
		ProjectId: input.ProjectID,
		Page:      page,
		Size:      size,
		//新增文件来源类型 2019/12/27
		SourceType: consts.OssPolicyTypeProjectResource,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return resourceList, nil
}
