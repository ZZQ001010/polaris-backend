package domain

import (
	"fmt"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/core/util/image"
	"github.com/galaxy-book/polaris-backend/common/core/util/str"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/resourcesvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/resourcesvc/po"
	"github.com/disintegration/imaging"
	"image/jpeg"
	"os"
	"strings"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

func CreateResource(createResourceBo bo.CreateResourceBo, tx ...sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	if ok, _ := slice.Contain([]int{consts.LocalResource, consts.OssResource, consts.DingDiskResource}, createResourceBo.Type); !ok {
		return 0, errs.BuildSystemErrorInfo(errs.InvalidResourceType)
	}
	isNameRight := format.VerifyResourceNameFormat(createResourceBo.Name)
	if !isNameRight {
		return 0, errs.InvalidResourceNameError
	}
	//新增folderId相关逻辑,为了保持原有逻辑,所以这里添加if条件分支 2019/12/12
	if createResourceBo.FolderId != nil {
		//if len(createResourceBo.Name) > 15 || createResourceBo.Name == "" {
		//	return 0, errs.InvalidResourceNameError
		//}
		//判断folderId是否存在
		folderIsExist, err := dao.FolderIdIsExist([]int64{*createResourceBo.FolderId}, createResourceBo.ProjectId, createResourceBo.OrgId)
		if err != nil {
			log.Error(err)
			return 0, err
		}
		if !folderIsExist {
			log.Error(errs.FolderIdNotExistError)
			return 0, errs.FolderIdNotExistError
		}
	}
	resourceId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableResource)
	if err != nil {
		return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	host, path := str.UrlParse(createResourceBo.Path)

	suffix := createResourceBo.Suffix
	if suffix == ""{
		suffix = util.ParseFileSuffix(createResourceBo.Name)
	}

	resourceEntity := po.PpmResResource{
		Id:        resourceId,
		OrgId:     createResourceBo.OrgId,
		ProjectId: createResourceBo.ProjectId,
		Host:      host,
		Path:      path,
		Name:      createResourceBo.Name,
		Type:      createResourceBo.Type,
		Suffix:    suffix,
		Bucket:    createResourceBo.Bucket,
		Size:      createResourceBo.Size,
		Md5:       createResourceBo.Md5,
		Creator:   createResourceBo.OperatorId,
		Updator:   createResourceBo.OperatorId,
		IsDelete:  consts.AppIsNoDelete,
	}
	//新增filetype逻辑   2019/12/20
	if createResourceBo.SourceType != nil {
		resourceEntity.SourceType = *createResourceBo.SourceType
	}
	//新增自动检测fileType逻辑 2019/12/30
	suffStr := strings.ToUpper(strings.TrimSpace(suffix))
	if value, ok := consts.FileTypes[suffStr]; ok {
		resourceEntity.FileType = value
	} else {
		resourceEntity.FileType = consts.FileTypeOthers
	}
	compressErr := compressImage(createResourceBo)
	if compressErr != nil {
		log.Error(compressErr)
	}

	err2 := dao.InsertResource(resourceEntity, tx...)
	if err2 != nil {
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}
	//新增folderId相关逻辑,为了保持原有逻辑,所以这里添加if条件分支 2019/12/12
	if createResourceBo.FolderId != nil {
		midtableId, err1 := idfacade.ApplyPrimaryIdRelaxed(consts.TableFolderResource)
		if err1 != nil {
			log.Error(err1)
			return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError, err1)
		}
		//插入中间表数据
		midtableEntity := po.PpmResFolderResource{
			Id:         midtableId,
			OrgId:      createResourceBo.OrgId,
			ResourceId: resourceId,
			FolderId:   *createResourceBo.FolderId,
			Creator:    createResourceBo.OperatorId,
			Updator:    createResourceBo.OperatorId,
			IsDelete:   consts.AppIsNoDelete,
		}
		err2 := dao.InsertMidTable(midtableEntity, tx...)
		if err2 != nil {
			log.Error(err2)
			return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
		}
	}
	return resourceId, nil
}

func CheckResourceIds(resourceIds []int64, projectId, orgId int64) errs.SystemErrorInfo {
	isExist, err := dao.ResourceIdIsExist(resourceIds, orgId, projectId)
	if err != nil {
		log.Error(err)
		return err
	}
	if !isExist {
		log.Error(errs.InvalidResourceIdsError)
		return errs.InvalidResourceIdsError
	}
	return nil
}

func CheckRelation(resourceIds []int64, folderId int64, orgId int64) errs.SystemErrorInfo {
	isExist, err := dao.RelationIsExist(resourceIds, folderId, orgId)
	if err != nil {
		return err
	}
	if !isExist {
		return errs.ResouceNotInFolderError
	}
	return nil
}
func DeleteResource(resourceIds []int64, folderId *int64, orgId, userId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	//仅文件删除关联关系
	if folderId != nil {
		upd := mysql.Upd{}
		upd[consts.TcIsDelete] = consts.AppIsDeleted
		upd[consts.TcUpdator] = userId
		upd[consts.TcUpdateTime] = time.Now()
		cond := db.Cond{
			consts.TcIsDelete:   consts.AppIsNoDelete,
			consts.TcFolderId:   folderId,
			consts.TcResourceId: db.In(resourceIds),
			consts.TcOrgId:      orgId,
		}
		err := dao.UpdateMidTableByCond(cond, upd, tx...)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	upd := mysql.Upd{}
	upd[consts.TcIsDelete] = consts.AppIsDeleted
	upd[consts.TcUpdator] = userId
	upd[consts.TcUpdateTime] = time.Now()
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(resourceIds),
	}
	_, err := dao.UpdateResourceByCond(cond, upd, tx...)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func compressImage(createResourceBo bo.CreateResourceBo) errs.SystemErrorInfo {
	if createResourceBo.Type == consts.LocalResource && createResourceBo.DistPath != "" {
		distPath := createResourceBo.DistPath
		suffix := createResourceBo.Suffix
		if _, ok := consts.ImgTypeMap[strings.ToUpper(suffix)]; ok {
			//固定高120
			newImg, err := image.ResizeAuto(distPath, 120, imaging.Lanczos)
			if err != nil {
				log.Error(err)
			} else {
				afterPath := util.GetCompressedPath(distPath, createResourceBo.Type)
				f, err := os.Create(afterPath)
				if err != nil {
					log.Error(err)
				} else {
					defer func() {
						if err := f.Close(); err != nil {
							log.Error(err)
						}
					}()
					imgErr := jpeg.Encode(f, newImg, nil)
					if imgErr != nil {
						log.Error(imgErr)
					}
				}
			}
		}
	}
	return nil
}

func InsertResource(tx sqlbuilder.Tx, resourcePath string, orgId int64, currentUserId int64, resourceType int, name string) (int64, errs.SystemErrorInfo) {
	resourceEntity := &po.PpmResResource{}
	if ok, _ := slice.Contain([]int{consts.LocalResource, consts.OssResource, consts.DingDiskResource}, resourceType); !ok {
		return 0, errs.BuildSystemErrorInfo(errs.InvalidResourceType)
	}
	nameSplit := strings.Split(resourcePath, "/")
	fileName := nameSplit[len(nameSplit)-1]
	suffix := ""
	if strings.Index(fileName, ".") != -1 {
		suffixSplit := strings.Split(resourcePath, ".")
		suffix = suffixSplit[len(suffixSplit)-1]
	}

	resourceId, err := idfacade.ApplyPrimaryIdRelaxed(resourceEntity.TableName())
	if err != nil {
		return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	if name != "" {
		fileName = name
	}
	host, path := str.UrlParse(resourcePath)
	*resourceEntity = po.PpmResResource{
		OrgId:  orgId,
		Path:   path,
		Name:   fileName,
		Suffix: suffix,
		//Md5:     md5.Md5V(nameSplit[len(nameSplit)-1]),
		Host:       host,
		Creator:    currentUserId,
		CreateTime: time.Now(),
		Id:         resourceId,
		Type:       resourceType,
		Updator:    currentUserId,
		UpdateTime: time.Now(),
		Version:    1,
	}
	_, err2 := tx.Collection(resourceEntity.TableName()).Insert(resourceEntity)
	if err2 != nil {
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}

	return resourceId, nil
}

//获取资源信息
func GetResourceByIds(resourceIds []int64) ([]bo.ResourceBo, errs.SystemErrorInfo) {
	resourceEntities := &[]po.PpmResResource{}
	err := mysql.SelectAllByCond((&po.PpmResResource{}).TableName(), db.Cond{
		consts.TcIsDelete: db.Eq(consts.AppIsNoDelete),
		consts.TcId:       db.In(resourceIds),
	}, resourceEntities)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.ResourceBo{}
	_ = copyer.Copy(resourceEntities, bos)

	return *bos, nil
}

func GetIdByPath(orgId int64, resourcePath string, resourceType int) (int64, errs.SystemErrorInfo) {
	resourceInfo := &bo.ResourceTypeBo{}
	host, path := str.UrlParse(resourcePath)
	err := mysql.SelectOneByCond(consts.TableResource, db.Cond{
		consts.TcIsDelete: db.Eq(consts.AppIsNoDelete),
		consts.TcPath:     db.Eq(path),
		consts.TcHost:     db.Eq(host),
		consts.TcOrgId:    orgId,
		consts.TcType:     resourceType,
	}, resourceInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return 0, errs.BuildSystemErrorInfo(errs.ResourceNotExist)
		} else {
			return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
	}
	return resourceInfo.ID, nil
}

func GetResourceBoList(page uint, size uint, input resourcevo.GetResourceBoListCond) (*[]bo.ResourceBo, int64, errs.SystemErrorInfo) {
	cond := db.Cond{}
	cond[consts.TcOrgId] = input.OrgId
	if input.ResourceIds != nil {
		cond[consts.TcId] = db.In(*input.ResourceIds)
	}
	if input.IsDelete != nil {
		cond[consts.TcIsDelete] = *input.IsDelete
	}
	pos, total, err := dao.SelectResourceByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "id desc",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.ResourceBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}
func GetResourceBoListByPage(cond db.Cond, pageBo bo.PageBo) (*[]bo.ResourceBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectResourceByPage(cond, pageBo)
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.ResourceBo{}
	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func UpdateResourceInfo(resourceId int64, input bo.UpdateResourceInfoBo, tx ...sqlbuilder.Tx) (*po.PpmResResource, errs.SystemErrorInfo) {
	resourcePo, err := dao.SelectResourceById(resourceId, tx...)
	if err != nil {
		return nil,errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if util.FieldInUpdate(input.UpdateFields, "fileName") && input.FileName != nil {
		isNameRight := format.VerifyResourceNameFormat(*input.FileName)
		if !isNameRight {
			return nil,errs.InvalidResourceNameError
		}
		resourcePo.Name = *input.FileName
	}
	if util.FieldInUpdate(input.UpdateFields, "fileSuffix") && input.FileSuffix != nil {
		resourcePo.Name = *input.FileSuffix
	}
	resourcePo.Updator = input.UserId
	resourcePo.UpdateTime = time.Now()
	err = dao.UpdateResource(*resourcePo, tx...)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return resourcePo, nil
}

func UpdateResourceFolderId(resourceIds []int64, currentFolderId, targetFolderId, userId, orgId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	//更新关联锁
	lockKey := fmt.Sprintf("%s%d", consts.UpdateResourceFolderLock, targetFolderId)
	uid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uid)
	if lockErr != nil{
		log.Error(lockErr)
		return errs.TryDistributedLockError
	}
	if suc{
		defer func() {
			if _, e := cache.ReleaseDistributedLock(lockKey, uid); e != nil{
				log.Error(e)
			}
		}()
		//查询有没有关联
		hasRelation, err := dao.ResourceFolderHasRelation(resourceIds, targetFolderId, orgId)
		if err != nil{
			log.Error(err)
			return err
		}

		//和目标文件夹已经有关联了就不要再移动了
		if hasRelation{
			return errs.SystemBusy
		}

		//修改老的关联
		upd := mysql.Upd{}
		upd[consts.TcFolderId] = targetFolderId
		upd[consts.TcUpdator] = userId
		upd[consts.TcUpdateTime] = time.Now()
		cond := db.Cond{
			consts.TcIsDelete:   consts.AppIsNoDelete,
			consts.TcFolderId:   currentFolderId,
			consts.TcResourceId: db.In(resourceIds),
			consts.TcOrgId:      orgId,
		}
		err = dao.UpdateMidTableByCond(cond, upd, tx...)
		if err != nil {
			log.Error(err)
			return err
		}
	}else{
		return errs.SystemBusy
	}
	return nil
}

func GetResourceIdsByFolderId(folderId, orgId int64) (*[]int64, errs.SystemErrorInfo) {
	midtablePos, err := dao.SelectMidTablePoByFolderId(folderId, orgId)
	if err != nil {
		return nil, err
	}
	var resourceIds []int64
	for _, value := range *midtablePos {
		resourceIds = append(resourceIds, value.ResourceId)
	}
	return &resourceIds, nil
}

func GetBaseUserInfoMap(orgId int64, userIds []int64) (map[interface{}]interface{}, errs.SystemErrorInfo) {
	ownerInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed("", orgId, userIds)
	if err != nil {
		return nil, err
	}
	ownerMap := maps.NewMap("UserId", ownerInfos)
	return ownerMap, nil
}
