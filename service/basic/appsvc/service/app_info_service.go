package services

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/appvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/idvo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/basic/appsvc/domain"
)

//func AppInfoList(ctx context.Context, page uint, size uint, cond db.Cond) (*vo.AppInfoList, error) {
//	cond["is_delete"] = consts.AppIsNoDelete
//
//	bos, total, err := basedomain.GetAppInfoBoList(page, size, cond)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
//	}
//
//	resultList := &[]*vo.AppInfo{}
//	copyErr := copyer.Copy(bos, resultList)
//	if copyErr != nil {
//		log.Errorf("对象copy异常: %v", copyErr)
//		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
//	}
//
//	return &vo.AppInfoList{
//		Total: total,
//		List:  *resultList,
//	}, nil
//}
//

func GetAppInfoByActive(code string) (*vo.AppInfo, errs.SystemErrorInfo) {
	bo, err := domain.GetAppInfoBoByCode(code)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if bo.CheckStatus != consts.AppCheckStatusSuccess || bo.Status != consts.AppStatusEnable {
		//状态不可用，不输出
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}

	appInfo := &vo.AppInfo{}
	copyErr := copyer.Copy(bo, appInfo)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return appInfo, nil
}

func CreateAppInfo(input appvo.CreateAppInfoReqVo) (*vo.Void, errs.SystemErrorInfo) {
	entity := &bo.AppInfoBo{}
	copyErr := copyer.Copy(input.CreateAppInfo, entity)

	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	entity.Creator = input.UserId
	entity.Updator = input.UserId

	idvo := idfacade.ApplyPrimaryId(idvo.ApplyPrimaryIdReqVo{Code: consts.TableAppInfo})
	if idvo.Code != errs.OK.Code() {
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError)
	}
	entity.Id = idvo.Id

	err1 := domain.CreateAppInfo(entity)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}

	return &vo.Void{
		ID: idvo.Id,
	}, nil
}

func UpdateAppInfo(input appvo.UpdateAppInfoReqVo) (*vo.Void, errs.SystemErrorInfo) {

	entity := &bo.AppInfoBo{}
	copyErr := copyer.Copy(input.Input, entity)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	entity.Updator = input.UserId

	//是否存在
	_, err2 := domain.GetAppInfoBoNoCache(entity.Id)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	}

	err1 := domain.UpdateAppInfo(entity)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	return &vo.Void{
		ID: input.Input.ID,
	}, nil
}

func DeleteAppInfo(input appvo.DeleteAppInfoReqVo) (*vo.Void, errs.SystemErrorInfo) {

	bo, err1 := domain.GetAppInfoBoNoCache(input.Input.ID)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}

	err2 := domain.DeleteAppInfo(bo, input.UserId)
	if err2 != nil {
		log.Error(err2)
		return nil, err2
	}

	return &vo.Void{
		ID: input.Input.ID,
	}, nil
}
