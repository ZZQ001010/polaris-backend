package service

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/service/platform/resourcesvc/domain"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

func CreateResource(createResourceBo bo.CreateResourceBo, tx ...sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	return domain.CreateResource(createResourceBo, tx...)
}

func UpdateResourceInfo(input bo.UpdateResourceInfoBo) (*resourcevo.UpdateResourceData, errs.SystemErrorInfo) {
	orgId := input.OrgId
	resourceId := input.ResourceId
	projectId := input.ProjectId
	updateFields := input.UpdateFields
	resp := &resourcevo.UpdateResourceData{}
	if updateFields == nil || len(updateFields) == 0 {
		return nil, errs.UpdateFiledIsEmpty
	}
	err := domain.CheckResourceIds([]int64{resourceId}, projectId, orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	bos, err := domain.GetResourceByIds([]int64{resourceId})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	oldBo := bos[0]
	resp.OldBo = append(resp.OldBo, oldBo)
	newPo, err := domain.UpdateResourceInfo(resourceId, input)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	newBo := &bo.ResourceBo{}
	err1 := copyer.Copy(newPo, newBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	resp.NewBo = append(resp.NewBo, *newBo)
	return resp, nil
}

func UpdateResourceFolder(input bo.UpdateResourceFolderBo) (*resourcevo.UpdateResourceData, errs.SystemErrorInfo) {
	orgId := input.OrgId
	userId := input.UserId
	resourceIds := input.ResourceIds
	projectId := input.ProjectId
	currentFolderId := input.CurrentFolderId
	targetFolderId := input.TargetFolderID
	resp := &resourcevo.UpdateResourceData{}
	err := domain.CheckFolderIds([]int64{currentFolderId,targetFolderId}, projectId, orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = domain.UpdateResourceFolderId(resourceIds, currentFolderId, targetFolderId, userId, orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	bos, err := domain.GetResourceByIds(resourceIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	blank:=""
	if currentFolderId!=0 {
		currentBos, err := domain.GetFolderById([]int64{currentFolderId})
		if err != nil {
			log.Error(err)
			return nil, err
		}
		resp.CurrentFolderName = &currentBos[0].Name
	}else {
		resp.CurrentFolderName = &blank
	}
	if targetFolderId!=0 {
		targetBos, err := domain.GetFolderById([]int64{targetFolderId})
		if err != nil {
			log.Error(err)
			return nil, err
		}
		resp.TargetFolderName = &targetBos[0].Name
	}else {
		resp.TargetFolderName = &blank
	}
	resp.OldBo = bos
	return resp, nil
}

func DeleteResource(deleteBo bo.DeleteResourceBo) (*resourcevo.UpdateResourceData, errs.SystemErrorInfo) {
	orgId := deleteBo.OrgId
	userId := deleteBo.UserId
	resourceIds := deleteBo.ResourceIds
	folderId := deleteBo.FolderId
	//仅文件做文件夹校验
	//if folderId != nil {
	//	err := domain.CheckRelation(resourceIds, *folderId, orgId)
	//	if err != nil {
	//		log.Error(err)
	//		return nil, err
	//	}
	//}
	resp := &resourcevo.UpdateResourceData{}
	bos, err := domain.GetResourceByIds(resourceIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = domain.DeleteResource(resourceIds, folderId, orgId, userId)
	if err != nil {
		return nil, err
	}
	resp.OldBo = bos
	return resp, nil
}
func GetResource(input bo.GetResourceBo) (*vo.ResourceList, errs.SystemErrorInfo) {
	folderId := input.FolderId
	orgId := input.OrgId
	projectId := input.ProjectId
	cond := db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
		consts.TcOrgId:     orgId,
		//新增获取的文件类型 2019/12/27
		consts.TcSourceType: input.SourceType,
	}
	if input.FileType != nil {
		cond[consts.TcFileType] = *input.FileType
	}
	if input.KeyWord != nil {
		cond[consts.TcName] = db.Like("%" + *input.KeyWord + "%")
	}
	pageBo := bo.PageBo{Page: input.Page, Size: input.Size, Order: "id desc"}
	if folderId != nil {
		err := domain.CheckFolderIds([]int64{*folderId}, projectId, orgId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		resourceIds, err := domain.GetResourceIdsByFolderId(*folderId, orgId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		cond[consts.TcId] = db.In(resourceIds)
	}

	resourceBos, total, err := domain.GetResourceBoListByPage(cond, pageBo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resourceVos := &[]*vo.Resource{}
	copyErr := copyer.Copy(resourceBos, resourceVos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	creatorIds := make([]int64, 0)
	for _, value := range *resourceVos {
		creatorIds = append(creatorIds, value.Creator)
		value.PathCompressed = util.GetCompressedPath(value.Host+value.Path, value.Type)
	}
	ownerMap, err := domain.GetBaseUserInfoMap(orgId, creatorIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for i, resourceVo := range *resourceVos {
		if ownerInfoInterface, ok := ownerMap[resourceVo.Creator]; ok {
			ownerInfo := ownerInfoInterface.(bo.BaseUserInfoBo)
			(*resourceVos)[i].CreatorName = ownerInfo.Name
		} else {
			log.Errorf("用户 %d 信息不存在，组织id %d", resourceVo.Creator, orgId)
		}
	}
	return &vo.ResourceList{
		List:  *resourceVos,
		Total: total,
	}, nil
}

//func InsertResource(tx sqlbuilder.Tx, resourcePath string, orgId int64, currentUserId int64, resourceType int, fileName string) (int64, errs.SystemErrorInfo) {
//	return domain.InsertResource(tx, resourcePath, orgId, currentUserId, resourceType, fileName)
//}

//获取资源信息
func GetResourceById(resourceIds []int64) ([]bo.ResourceBo, errs.SystemErrorInfo) {
	return domain.GetResourceByIds(resourceIds)
}

func GetIdByPath(orgId int64, resourcePath string, resourceType int) (int64, errs.SystemErrorInfo) {
	return domain.GetIdByPath(orgId, resourcePath, resourceType)
}
