package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/times"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetIssueDetailBo(orgId, issueId int64) (*bo.IssueDetailBo, errs.SystemErrorInfo) {
	//获取issue详情
	issueDetail := &po.PpmPriIssueDetail{}
	err := mysql.SelectOneByCond(issueDetail.TableName(), db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIssueId:  issueId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, issueDetail)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDetailNotExist)
	}
	issueDetailBo := &bo.IssueDetailBo{}
	err1 := copyer.Copy(issueDetail, issueDetailBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return issueDetailBo, nil
}

func UpdateIssueDetailRemark(issueBo bo.IssueBo, operatorId int64, remark string) errs.SystemErrorInfo {
	//detail
	issueDetail := &po.PpmPriIssueDetail{}
	_, err := mysql.UpdateSmartWithCond(issueDetail.TableName(), db.Cond{
		consts.TcOrgId:    issueBo.OrgId,
		consts.TcIssueId:  issueBo.Id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcRemark:     remark,
		consts.TcUpdator:    operatorId,
		consts.TcUpdateTime: times.GetBeiJingTime(),
	})
	if err != nil {
		log.Errorf("mysql.UpdateSmartWithCond: %c\n", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}
