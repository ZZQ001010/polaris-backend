package domain

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//批量获取多个任务的tags
//入参：组织id，任务id列表
func GetIssueTagsByIssueIds(orgId int64, issueIds []int64) ([]bo.IssueTagBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return nil, errs.MysqlOperateError
	}

	issueTagInfo := &[]po.TagInfoWithIssue{}
	err1 := conn.Select(db.Raw("i.issue_id, i.tag_id, t.name, t.bg_style, t.font_style")).From("ppm_pri_issue_tag i", "ppm_pri_tag t").Where(db.Cond{
		"i.issue_id":  db.In(issueIds),
		"i.is_delete": consts.AppIsNoDelete,
		"t.is_delete": consts.AppIsNoDelete,
		"i.tag_id":    db.Raw("t.id"),
		"t.org_id":    orgId,
	}).All(issueTagInfo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}
	issueTagBos := make([]bo.IssueTagBo, 0)
	for _, v := range *issueTagInfo {
		issueTagBos = append(issueTagBos, bo.IssueTagBo{
			TagId:     v.TagId,
			IssueId:   v.IssueId,
			TagName:   v.Name,
			FontStyle: v.FontStyle,
			BgStyle:   v.BgStyle,
		})
	}
	return issueTagBos, nil
}

//任务关联tags
func IssueRelateTags(orgId, projectId, issueId, operatorId int64, addTags []bo.IssueTagReqBo, delTags []bo.IssueTagReqBo) errs.SystemErrorInfo {
	err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if addTags != nil && len(addTags) > 0 {
			err := IssueAddTags(orgId, projectId, issueId, operatorId, addTags, tx)
			if err != nil {
				log.Error(err)
				return err
			}
		}
		if delTags != nil && len(delTags) > 0 {
			delTagIds := make([]int64, 0)
			for _, tagInfo := range delTags {
				delTagIds = append(delTagIds, tagInfo.Id)
			}
			_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueTag, db.Cond{
				consts.TcOrgId:    orgId,
				consts.TcIssueId:  issueId,
				consts.TcTagId:    db.In(delTagIds),
				consts.TcIsDelete: consts.AppIsNoDelete,
			}, mysql.Upd{
				consts.TcIsDelete: consts.AppIsDeleted,
				consts.TcUpdator:  operatorId,
			})
			if err1 != nil {
				log.Error(err1)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
			}
		}
		return nil
	})
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}
	return nil
}

func IssueAddTags(orgId, projectId, issueId, operatorId int64, addTags []bo.IssueTagReqBo, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	//防止项目成员重复插入
	//uid := uuid.NewUuid()
	//issueIdStr := strconv.FormatInt(issueId, 10)
	//lockKey := consts.AddIssueTagsLock + issueIdStr
	//suc, err := cache.TryGetDistributedLock(lockKey, uid)
	//if err != nil {
	//	log.Errorf("获取%s锁时异常 %v", lockKey, err)
	//	return errs.TryDistributedLockError
	//}
	//if suc {
	//	defer cache.ReleaseDistributedLock(lockKey, uid)
	//} else {
	//	return errs.BuildSystemErrorInfo(errs.GetDistributedLockError)
	//}

	addTagIds := make([]int64, 0)
	for _, tagInfo := range addTags {
		addTagIds = append(addTagIds, tagInfo.Id)
	}

	//预先查询已有的关联
	issueTags := &[]po.PpmPriIssueTag{}
	err5 := mysql.SelectAllByCond(consts.TableIssueTag, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIssueId:  issueId,
		consts.TcTagId:    db.In(addTagIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, issueTags)
	if err5 != nil {
		log.Error(err5)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	//删除已存在的关联
	notRelationTags := make([]bo.IssueTagReqBo, 0)
	alreadyExistTagIds := make([]int64, 0)
	for _, issueTag := range *issueTags {
		alreadyExistTagIds = append(alreadyExistTagIds, issueTag.TagId)
	}

	//去重ids
	repetIds := make([]int64, 0)
	for _, issueTag := range addTags {
		exist, err := slice.Contain(alreadyExistTagIds, issueTag.Id)
		if err != nil {
			log.Error(err)
			continue
		}
		if !exist {
			if exist, _ = slice.Contain(repetIds, issueTag.Id); !exist {
				notRelationTags = append(notRelationTags, issueTag)
				repetIds = append(repetIds, issueTag.Id)
			}
		}
	}
	addTags = notRelationTags

	if len(addTags) == 0 {
		return nil
	}
	ids, err1 := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueTag, len(addTags))
	if err1 != nil {
		log.Error(err1)
		return err1
	}
	pos := make([]po.PpmPriIssueTag, len(addTags))
	for i, addTag := range addTags {
		pos[i] = po.PpmPriIssueTag{
			Id:        ids.Ids[i].Id,
			OrgId:     orgId,
			ProjectId: projectId,
			IssueId:   issueId,
			TagId:     addTag.Id,
			TagName:   addTag.Name,
			Creator:   operatorId,
			Updator:   operatorId,
		}
	}

	err := mysql.TransBatchInsert(tx, &po.PpmPriIssueTag{}, slice.ToSlice(pos))
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

//获取任务的标签列表
//入参：组织id，任务id
func GetIssueTags(orgId, issueId int64) ([]bo.IssueTagBo, errs.SystemErrorInfo) {
	return GetIssueTagsByIssueIds(orgId, []int64{issueId})
}

func IssueTagStat(tagIds []int64) ([]bo.IssueTagStatBo, errs.SystemErrorInfo) {
	resPo := &[]po.IssueTagStat{}
	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	selectErr := conn.Select(db.Raw("count(distinct(issue_id)) as total, tag_id")).From(consts.TableIssueTag).Where(db.Cond{
		consts.TcIsDelete:consts.AppIsNoDelete,
		consts.TcTagId:db.In(tagIds),
	}).GroupBy(consts.TcTagId).All(resPo)
	if selectErr != nil {
		log.Error(selectErr)
		return nil, errs.MysqlOperateError
	}

	resBo := &[]bo.IssueTagStatBo{}
	_ = copyer.Copy(resPo, resBo)

	return *resBo, nil
}

func DeleteIssueTags(orgId int64, issueIds []int64, operatorId int64) errs.SystemErrorInfo {
	_, err := mysql.UpdateSmartWithCond(consts.TableIssueTag, db.Cond{
		consts.TcOrgId:orgId,
		consts.TcIssueId:db.In(issueIds),
		consts.TcIsDelete:consts.AppIsNoDelete,
	}, mysql.Upd{
			consts.TcUpdator:operatorId,
			consts.TcIsDelete:consts.AppIsDeleted,
	})

	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}

	return nil
}