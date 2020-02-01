package domain

import (
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func RelateIteration(orgId, projectId, iterationId, operatorId int64, addIssueIds []int64, delIssueIds []int64) errs.SystemErrorInfo {
	projectObjectType, err := GetProjectObjectTypeByLangCodeAndObjectType(orgId, projectId, consts.ProjectObjectTypeLangCodeFeature, consts.ProjectObjectTypeTask)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(addIssueIds) > 0 {
			_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
				consts.TcOrgId:               orgId,
				consts.TcId:                  db.In(addIssueIds),
				consts.TcProjectId:           projectId,
				consts.TcIterationId:         db.Eq(0),
				consts.TcProjectObjectTypeId: db.NotEq(projectObjectType.Id),
				consts.TcIsDelete:            consts.AppIsNoDelete,
			}, mysql.Upd{
				consts.TcIterationId: iterationId,
				consts.TcUpdator:     operatorId,
			})
			if err1 != nil {
				log.Error(err1)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
			}
		}

		if len(delIssueIds) > 0 {
			_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
				consts.TcOrgId:               orgId,
				consts.TcId:                  db.In(delIssueIds),
				consts.TcProjectId:           projectId,
				consts.TcIterationId:         db.NotEq(0),
				consts.TcProjectObjectTypeId: db.NotEq(projectObjectType.Id),
				consts.TcIsDelete:            consts.AppIsNoDelete,
			}, mysql.Upd{
				consts.TcIterationId: 0,
				consts.TcUpdator:     operatorId,
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

func JudgeIterationIsExist(orgId, id int64) bool {
	return dao.JudgeIterationIsExist(orgId, id)
}
