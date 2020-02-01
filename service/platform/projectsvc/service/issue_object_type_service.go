package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
)

func IssueObjectTypeList(orgId int64, page uint, size uint, params *vo.IssueObjectTypesReq) (*vo.IssueObjectTypeList, errs.SystemErrorInfo) {
	var typeList interface{} = nil
	var err1 error = nil
	if params != nil && params.ProjectObjectTypeID != nil {
		typeList, err1 = domain.GetIssueObjectTypeListByProjectObjectTypeId(orgId, *params.ProjectObjectTypeID)
	} else {
		typeList, err1 = domain.GetIssueObjectTypeList(orgId)
	}
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err1)
	}
	resultList := &[]*vo.IssueObjectType{}
	copyErr := copyer.Copy(typeList, resultList)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return &vo.IssueObjectTypeList{
		Total: int64(len(*resultList)),
		List:  *resultList,
	}, nil
}

func CreateIssueObjectType(currentUserId int64, input vo.CreateIssueObjectTypeReq) (*vo.Void, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//currentUserId := cacheUserInfo.UserId

	//TODO 权限
	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	//}

	entity := &bo.IssueObjectTypeBo{}
	copyErr := copyer.Copy(input, entity)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueObjectType)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	entity.Id = id
	entity.Creator = currentUserId
	entity.Updator = currentUserId

	//删除缓存
	err = domain.DeleteIssueObjectTypeListCache(input.OrgID)

	if err != nil {
		return nil, err
	}

	err1 := domain.CreateIssueObjectType(entity)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	return &vo.Void{
		ID: id,
	}, nil
}

func UpdateIssueObjectType(currentUserId int64, input vo.UpdateIssueObjectTypeReq) (*vo.Void, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//currentUserId := cacheUserInfo.UserId

	//TODO 权限
	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	//}

	entity := &bo.IssueObjectTypeBo{}
	copyErr := copyer.Copy(input, entity)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	entity.Updator = currentUserId

	//是否存在
	_, err2 := domain.GetIssueObjectTypeBo(entity.Id)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	}

	//删除缓存
	err := domain.DeleteIssueObjectTypeListCache(input.OrgID)

	if err != nil {
		return nil, err
	}

	err1 := domain.UpdateIssueObjectType(entity)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	return &vo.Void{
		ID: input.ID,
	}, nil
}

func DeleteIssueObjectType(orgId, currentUserId int64, input vo.DeleteIssueObjectTypeReq) (*vo.Void, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//currentUserId := cacheUserInfo.UserId
	targetId := input.ID

	//TODO 权限
	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	//}

	bo, err1 := domain.GetIssueObjectTypeBo(targetId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	//删除缓存 暂时先用orgId 后面用 input中的orgId orgId的包含校验放在权限中去处理
	err := domain.DeleteIssueObjectTypeListCache(orgId)

	if err != nil {
		return nil, err
	}

	err2 := domain.DeleteIssueObjectType(bo, currentUserId)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	}

	return &vo.Void{
		ID: targetId,
	}, nil
}
