package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/processvo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/domain"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func ProcessStatusList(page uint, size uint) (*vo.ProcessStatusList, errs.SystemErrorInfo) {
	cond := db.Cond{}
	cond[consts.TcIsDelete] = consts.AppIsNoDelete

	bos, total, err := domain.GetProcessStatusBoList(page, size, cond)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}

	resultList := &[]*vo.ProcessStatus{}
	copyErr := copyer.Copy(bos, resultList)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &vo.ProcessStatusList{
		Total: total,
		List:  *resultList,
	}, nil
}

func CreateProcessStatus(req processvo.CreateProcessStatusReqVo) (*vo.Void, errs.SystemErrorInfo) {

	input := req.CreateProcessStatusReq

	currentUserId := req.UserId

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

	entity := &bo.ProcessStatusBo{}
	copyErr := copyer.Copy(input, entity)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProcessStatus)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	entity.Id = id
	entity.Creator = currentUserId
	entity.Updator = currentUserId

	//清除缓存 包含关联关系缓存
	err = domain.DeleteCacheProcessStatusList(req.OrgId, entity.Id)

	if err != nil {
		return nil, err
	}

	err1 := domain.CreateProcessStatus(entity)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	return &vo.Void{
		ID: id,
	}, nil
}

func UpdateProcessStatus(req processvo.UpdateProcessStatusReqVo) (*vo.Void, errs.SystemErrorInfo) {
	//TODO 权限
	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	//}

	input := req.UpdateProcessStatusReq

	currentUserId := req.UserId

	entity := &bo.ProcessStatusBo{}
	copyErr := copyer.Copy(input, entity)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	entity.Updator = currentUserId

	//是否存在
	_, err2 := domain.GetProcessStatusBo(entity.Id)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	}

	//清除缓存 包含关联关系缓存
	err := domain.DeleteCacheProcessStatusList(req.OrgId, entity.Id)

	if err != nil {
		return nil, err
	}

	err1 := domain.UpdateProcessStatus(entity)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	return &vo.Void{
		ID: input.ID,
	}, nil
}

func DeleteProcessStatus(req processvo.DeleteProcessStatusReq) (*vo.Void, errs.SystemErrorInfo) {

	input := req.DeleteProcessStatusReq

	currentUserId := req.UserId

	targetId := input.ID

	//TODO 权限
	//err = AuthIssue(orgId, currentUserId, input.ID, consts.RoleOperationPathOrgProIssueT, consts.RoleOperationModify)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	//}

	bo, err1 := domain.GetProcessStatusBo(targetId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}
	//清除缓存 包含关联关系缓存
	err := domain.DeleteCacheProcessStatusList(req.OrgId, targetId)

	if err != nil {
		return nil, err
	}

	err2 := domain.DeleteProcessStatus(bo, currentUserId)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	}

	return &vo.Void{
		ID: targetId,
	}, nil
}

func GetProcessStatus(orgId int64, id int64) (*bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	return domain.GetProcessStatus(orgId, id)
}

func ProcessStatusInit(orgId int64, contextMap map[string]interface{}) errs.SystemErrorInfo {
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		return domain.ProcessStatusInit(orgId, contextMap, tx)
	})
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func GetProcessStatusByCategory(orgId int64, statusId int64, category int) (*bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	return domain.GetProcessStatusByCategory(orgId, statusId, category)
}

func GetProcessStatusListByCategory(orgId int64, category int) ([]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	return domain.GetProcessStatusListByCategory(orgId, category)
}

func GetProcessStatusIds(orgId int64, category int, typ int) (*[]int64, errs.SystemErrorInfo) {
	return domain.GetProcessStatusIds(orgId, category, typ)
}

func GetProcessStatusList(orgId, processId int64) (*[]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	return domain.GetProcessStatusList(orgId, processId)
}

func GetProcessInitStatusId(orgId, projectId, projectObjectTypeId int64, category int) (int64, errs.SystemErrorInfo) {
	return domain.GetProcessInitStatusId(orgId, projectId, projectObjectTypeId, category)
}

func GetDefaultProcessStatusId(orgId int64, processId int64, category int) (int64, errs.SystemErrorInfo) {
	return domain.GetDefaultProcessStatusId(orgId, processId, category)
}
