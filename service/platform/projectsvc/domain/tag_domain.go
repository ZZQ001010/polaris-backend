package domain

import (
	"fmt"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/pinyin"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetTagCount(cond db.Cond) (uint64, errs.SystemErrorInfo) {
	count, err := mysql.SelectCountByCond(consts.TableTag, cond)
	if err != nil {
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return count, nil
}

//func InsertTag(orgId, userId int64, name string, bgStyle string, fontStyle string, projectId int64, tagsInfo []bo.TagBo) (int64, errs.SystemErrorInfo) {
func InsertTag(orgId, userId int64, projectId int64, tagsInfo []bo.TagBo) ([]bo.TagBo, errs.SystemErrorInfo) {

	newUUID := uuid.NewUuid()
	lockKey := fmt.Sprintf("%s%d", consts.CreateTagLock, projectId)
	suc, lockErr := cache.TryGetDistributedLock(lockKey, newUUID)
	if lockErr != nil{
		log.Error(lockErr)
		return nil, errs.TryDistributedLockError
	}
	if suc{
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, newUUID); err != nil{
				log.Error(err)
			}
		}()
	}else{
		//未获取到锁，直接响应错误信息
		return nil, errs.RepeatTag
	}

	nameArr := []string{}
	newTagBo := []bo.TagBo{}
	//去重
	for _, tagBo := range tagsInfo {
		if ok, _ := slice.Contain(nameArr, tagBo.Name); !ok {
			newTagBo = append(newTagBo, tagBo)
			nameArr = append(nameArr, tagBo.Name)
		}
	}
	//查看是否重复（如果已有，则直接返回现有的标签）
	tagPos := &[]po.PpmPriTag{}
	err := mysql.SelectAllByCond(consts.TableTag,db.Cond{
		consts.TcOrgId: orgId,
		consts.TcName:  db.In(nameArr),
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcProjectId:projectId,
	}, tagPos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	} else if len(*tagPos) < len(newTagBo) {
		existTags := []string{}
		for _, tag := range *tagPos {
			existTags = append(existTags, tag.Name)
		}
		ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableTag, len(newTagBo) - len(*tagPos))
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
		}
		i := 0
		insertPos := []po.PpmPriTag{}
		for _, tagBo := range newTagBo {
			if ok, _ := slice.Contain(existTags, tagBo.Name); ok {
				continue
			}
			spellName := pinyin.ConvertToPinyin(tagBo.Name)
			insertPos = append(insertPos, po.PpmPriTag{
				Id:         ids.Ids[i].Id,
				OrgId:      orgId,
				ProjectId:  projectId,
				Name:       tagBo.Name,
				NamePinyin: spellName,
				Creator:    userId,
				BgStyle:    tagBo.BgStyle,
				FontStyle:  tagBo.FontStyle,
			})
			i++
		}
		insertErr := mysql.BatchInsert(&po.PpmPriTag{}, slice.ToSlice(insertPos))
		if insertErr != nil {
			log.Error(insertErr)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertErr)
		}
		*tagPos = append(*tagPos, insertPos...)
	}

	tagBos := &[]bo.TagBo{}
	copyErr := copyer.Copy(tagPos, tagBos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.ObjectCopyError
	}

	return *tagBos, nil
}

func GetTagList(cond db.Cond, page int, size int) (uint64, *[]bo.TagBo, errs.SystemErrorInfo) {
	count, err := GetTagCount(cond)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	tagPo := &[]po.PpmPriTag{}
	listErr := mysql.SelectAllByCondWithNumAndOrder(consts.TableTag, cond, nil, page, size, "create_time desc", tagPo)
	if listErr != nil {
		log.Error(listErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	tagBo := &[]bo.TagBo{}
	copyErr := copyer.Copy(tagPo, tagBo)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return count, tagBo, nil
}

func GetTagInfo(tagId int64) (*bo.TagBo, errs.SystemErrorInfo) {
	info := &po.PpmPriTag{}
	err := mysql.SelectOneByCond(consts.TableTag, db.Cond{
		consts.TcIsDelete:consts.AppIsNoDelete,
		consts.TcId:tagId,
	}, info)
	if err != nil {
		log.Error(err)
		if err == db.ErrNoMoreRows {
			return nil, errs.TagNotExist
		} else {
			return nil, errs.MysqlOperateError
		}
	}

	infoBo := &bo.TagBo{}
	_ = copyer.Copy(info, infoBo)

	return infoBo, nil
}

func DeleteTag(orgId, userId int64, projectId int64, tagIds []int64) (int64, errs.SystemErrorInfo) {
	var count int64
	var err error
	err = mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除标签
		count, err = mysql.UpdateSmartWithCond(consts.TableTag, db.Cond{
			consts.TcId:        db.In(tagIds),
			consts.TcIsDelete:  consts.AppIsNoDelete,
			consts.TcProjectId: projectId,
			consts.TcOrgId:     orgId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  userId,
		})
		if err != nil {
			log.Error(err)
			return err
		}

		//删除标签关联
		_, err1 := mysql.UpdateSmartWithCond(consts.TableIssueTag, db.Cond{
			consts.TcOrgId:orgId,
			consts.TcIsDelete:consts.AppIsNoDelete,
			consts.TcTagId:db.In(tagIds),
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  userId,
		})
		if err1 != nil {
			log.Error(err1)
			return err1
		}

		return nil
	})
	if err != nil {
		log.Error(err)
		return 0, errs.MysqlOperateError
	}

	return count, nil
}

func UpdateTag(orgId, userId, tagId int64, upd mysql.Upd) errs.SystemErrorInfo {
	upd[consts.TcUpdator] = userId
	_, err := mysql.UpdateSmartWithCond(consts.TableTag, db.Cond{
		consts.TcId:        tagId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcOrgId:     orgId,
	}, upd)
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}

	return nil
}