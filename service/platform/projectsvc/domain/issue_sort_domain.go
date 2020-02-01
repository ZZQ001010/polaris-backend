package domain

import (
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func UpdateIssueSort(issueBo bo.IssueBo, operatorId int64, beforeId, afterId *int64) errs.SystemErrorInfo{
	orgId := issueBo.OrgId

	isBefore := false
	refId := int64(0)

	if beforeId != nil{
		refId = *beforeId
		isBefore = true
	}else if afterId != nil{
		refId = *afterId
	}
	if refId == issueBo.Id{
		log.Error("要排序的任务和目标任务的id不能一致")
		return errs.BuildSystemErrorInfo(errs.IssueSortReferenceInvalidError)
	}

	refIssueBo, err := GetIssueBo(orgId, refId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error{
		targetSort := issueBo.Sort
		if isBefore{
			targetSort = refIssueBo.Sort + 1
			_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
				consts.TcOrgId: orgId,
				consts.TcProjectId: refIssueBo.ProjectId,
				consts.TcProjectObjectTypeId: refIssueBo.ProjectObjectTypeId,
				consts.TcSort: db.Gt(refIssueBo.Sort),
			}, mysql.Upd{
				consts.TcSort: db.Raw("sort + 1"),
				consts.TcUpdator: operatorId,
			})
			if err1 != nil{
				log.Error(err1)
				return err1
			}
		}else{
			targetSort = refIssueBo.Sort - 1
			_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
				consts.TcOrgId: orgId,
				consts.TcProjectId: refIssueBo.ProjectId,
				consts.TcProjectObjectTypeId: refIssueBo.ProjectObjectTypeId,
				consts.TcSort: db.Lt(refIssueBo.Sort),
			}, mysql.Upd{
				consts.TcSort: db.Raw("sort - 1"),
				consts.TcUpdator: operatorId,
			})
			if err2 != nil{
				log.Error(err2)
				return err2
			}
		}

		err3 := mysql.TransUpdateSmart(tx, consts.TableIssue, issueBo.Id, mysql.Upd{
			consts.TcSort: targetSort,
			consts.TcUpdator: operatorId,
		})
		if err3 != nil{
			log.Error(err3)
			return err3
		}
		return nil
	})
	if transErr != nil{
		log.Error(transErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}