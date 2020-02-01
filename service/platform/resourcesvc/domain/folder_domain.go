package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/resourcesvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/resourcesvc/po"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func CreateFolder(input bo.CreateFolderBo) (int64, errs.SystemErrorInfo) {
	folderId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableFolder)
	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	folderPo := po.PpmResFolder{
		Id:        folderId,
		OrgId:     input.OrgId,
		ProjectId: input.ProjectId,
		Name:      input.Name,
		ParentId:  input.ParentId,
		FileType:  input.FileType,
		Creator:   input.UserId,
		Updator:   input.UserId,
	}
	err0 := dao.InsertFolder(folderPo)
	if err0 != nil {
		log.Error(err0)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err0)
	}
	return folderId, nil
}

//func UpdateFolderParentId(folderId, parentId, userId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
//	folderPo, err := dao.SelectFolderById(folderId, tx...)
//	if err != nil {
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	folderPo.Updator = userId
//	folderPo.UpdateTime = time.Now()
//	folderPo.ParentId = parentId
//	err = dao.UpdateFolder(*folderPo, tx...)
//	if err != nil {
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	return nil
//}

//func UpdateFolderName(folderId, userId int64, newName string, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
//	folderPo, err := dao.SelectFolderById(folderId, tx...)
//	if err != nil {
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	folderPo.Updator = userId
//	folderPo.UpdateTime = time.Now()
//	folderPo.Name = newName
//	err = dao.UpdateFolder(*folderPo, tx...)
//	if err != nil {
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	return nil
//}

func UpdateFolder(folderId int64, input bo.UpdateFolderBo, tx ...sqlbuilder.Tx) (mysql.Upd, errs.SystemErrorInfo) {
	//folderPos, err := dao.SelectFolderByIds([]int64{folderId}, tx...)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}
	//folderPo := (*folderPos)[0]
	//if util.FieldInUpdate(input.UpdateFields, "parentId") && input.ParentID != nil {
	//	folderPo.ParentId = *input.ParentID
	//}
	//if util.FieldInUpdate(input.UpdateFields, "name") && input.Name != nil {
	//	folderPo.Name = *input.Name
	//}
	//folderPo.Updator = input.UserId
	//folderPo.UpdateTime = time.Now()
	//err = dao.UpdateFolder(folderPo, tx...)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}
	//return &folderPo, nil
	upd := mysql.Upd{}
	if util.FieldInUpdate(input.UpdateFields, "parentId") && input.ParentID != nil {
		err := CheckFolderIds([]int64{*input.ParentID}, input.ProjectID, input.OrgId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		upd[consts.TcParentId] = *input.ParentID
	} else if util.FieldInUpdate(input.UpdateFields, "name") && input.Name != nil {
		isNameRight := format.VerifyFolderNameFormat(*input.Name)
		if !isNameRight {
			return nil,errs.InvalidFolderNameError
		}
		upd[consts.TcName] = *input.Name
	}
	if len(upd) != 0 {
		upd[consts.TcUpdator] = input.UserId
		upd[consts.TcUpdateTime] = time.Now()
		cond := db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcId:       folderId,
		}
		err := dao.UpdateFolderByCond(cond, upd)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	return upd, nil
}

func CheckFolderIds(folderIds []int64, projectId, orgId int64) errs.SystemErrorInfo {
	isExist, err := dao.FolderIdIsExist(folderIds, projectId, orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	if !isExist {
		log.Error(errs.InvalidFolderIdsError)
		return errs.InvalidFolderIdsError
	}
	return nil
}
func DeleteFolder(folderIds []int64, userId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	upd := mysql.Upd{}
	upd[consts.TcIsDelete] = consts.AppIsDeleted
	upd[consts.TcUpdator] = userId
	upd[consts.TcUpdateTime] = time.Now()
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(folderIds),
	}
	err := dao.UpdateFolderByCond(cond, upd, tx...)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func GetFolder(parentId *int64, projectId int64, page bo.PageBo) (*[]po.PpmResFolder, uint64, errs.SystemErrorInfo) {
	cond := db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
	}
	if parentId != nil {
		cond[consts.TcParentId] = *parentId
	}
	folderPos, total, err := dao.SelectFolderByPage(cond, page)
	if err != nil {
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return folderPos, total, nil
}

func GetFolderById(folderIds []int64) ([]bo.FolderBo, errs.SystemErrorInfo) {
	resourceEntities := &[]po.PpmResFolder{}
	err := mysql.SelectAllByCond((&po.PpmResFolder{}).TableName(), db.Cond{
		consts.TcIsDelete: db.Eq(consts.AppIsNoDelete),
		consts.TcId:       db.In(folderIds),
	}, resourceEntities)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.FolderBo{}
	_ = copyer.Copy(resourceEntities, bos)

	return *bos, nil
}
