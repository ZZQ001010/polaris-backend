package domain

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func SelectList(cond db.Cond, union *db.Union, page int, size int, order interface{}) (*[]bo.IssueBo, int64, errs.SystemErrorInfo) {
	issues := &[]*po.PpmPriIssue{}
	total, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableIssue, cond, union, page, size, order, issues)
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	issueBos := &[]bo.IssueBo{}
	err2 := copyer.Copy(*issues, issueBos)
	if err2 != nil {
		log.Error(err2)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err2)
	}

	return issueBos, int64(total), nil
}

func SelectIssueRemindInfoList(issueIdsCondBo bo.SelectIssueIdsCondBo, page int, size int) ([]bo.IssueRemindInfoBo, int64, errs.SystemErrorInfo) {
	issuePrefix := "i."
	projectPrefix := "p."
	//规则，未完成
	cond := db.Cond{
		issuePrefix + consts.TcIsDelete: consts.AppIsNoDelete,
		//通过end_time来筛选未完成的项目
		issuePrefix + consts.TcEndTime:    db.Lt(consts.BlankElasticityTime),
		projectPrefix + consts.TcIsDelete: consts.AppIsNoDelete,
		projectPrefix + consts.TcIsFiling: consts.AppIsNotFilling,
	}

	//条件组装
	if issueIdsCondBo.AfterPlanEndTime != nil && issueIdsCondBo.BeforePlanEndTime != nil {
		cond[issuePrefix+consts.TcPlanEndTime] = db.Between(*issueIdsCondBo.BeforePlanEndTime, *issueIdsCondBo.AfterPlanEndTime)
	}

	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	mid := conn.Select(
		issuePrefix+consts.TcId,
		issuePrefix+consts.TcPlanEndTime,
		issuePrefix+consts.TcOwner,
		issuePrefix+consts.TcOrgId,
		issuePrefix+consts.TcProjectId,
		issuePrefix+consts.TcTitle,
		issuePrefix+consts.TcParentId,
	).From(consts.TableIssue + " i").LeftJoin(consts.TableProject + " p").On("i.project_id = p.id").Where(cond).OrderBy("i.id asc").Paginate(uint(size)).Page(uint(page))

	//查询总数
	total, err := mid.TotalEntries()
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	issueIdBos := &[]bo.IssueRemindInfoBo{}
	//总数大于0的话才需要去取数据
	if total > 0 {
		err = mid.All(issueIdBos)
		if err != nil {
			log.Error(err)
			return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
	}
	return *issueIdBos, int64(total), nil
}

func AllIssueForProject(orgId, projectId int64, isParent bool) ([]bo.IssueAndDetailInfoBo, errs.SystemErrorInfo) {
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
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	resBo := &[]bo.IssueAndDetailInfoBo{}
	cond := db.Cond{
		"i." + consts.TcIsDelete:  consts.AppIsNoDelete,
		"i." + consts.TcId:        db.Raw("d." + consts.TcIssueId),
		"i." + consts.TcOrgId:     orgId,
		"i." + consts.TcProjectId: projectId,
	}
	if !isParent {
		cond["i."+consts.TcParentId] = db.NotEq(0)
	} else {
		cond["i."+consts.TcParentId] = 0
	}
	err = conn.Select(db.Raw("i.id, i.project_object_type_id, i.title, i.priority_id, i.plan_start_time, i.plan_end_time, i.parent_id, i.owner, i.status, i.creator, i.create_time, d.remark")).
		From("ppm_pri_issue i", "ppm_pri_issue_detail d").Where(cond).OrderBy("i.project_object_type_id").All(resBo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	return *resBo, nil
}
