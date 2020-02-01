package domain

import (
	"fmt"
	"github.com/galaxy-book/common/core/util/date"
	"github.com/galaxy-book/common/core/util/times"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"time"
	"upper.io/db.v3"
)

var getProcessError = "proxies.GetProcessStatusId: %v\n"

func IssueCondRelatedTypeAssembly(issueCond db.Cond, relatedType int, relatedUserId int64) {
	switch relatedType {
	case 1:
		issueCond["creator"] = relatedUserId
	case 2:
		issueCond[consts.TcOwner] = relatedUserId
	case 3, 4:
		rt := consts.IssueRelationTypeParticipant
		if relatedType == 4 {
			rt = consts.IssueRelationTypeFollower
		}
		issueCond[consts.TcId] = db.In(db.Raw("select ir.issue_id as id from ppm_pri_issue_relation ir where ir.is_delete = 2 and ir.relation_id = ? and ir.relation_type = ?", relatedUserId, rt))
	}
}

func IssueCondNoRelatedTypeAssembly(issueCond db.Cond, relatedUserId int64, orgId int64, isAdmin bool) {
	args := []interface{}{orgId}
	sql := "select p.id as id from ppm_pro_project p where p.org_id = ? and p.is_delete = 2"
	//增加项目限制
	if ! isAdmin{
		sql += " and (p.public_status = 1 or p.id in (SELECT DISTINCT pr.project_id FROM ppm_pro_project_relation pr WHERE pr.relation_id = ? AND relation_type in (1,2,3) AND pr.is_delete = 2))"
		args = append(args, relatedUserId)
	}
	issueCond[consts.TcProjectId] = db.In(db.Raw(sql, args...))
}

func IssueCondTagId(issueCond db.Cond, orgId int64, tagIds []int64) {
	//查询tag
	issueCond[consts.TcId+" IN"] = db.In(db.Raw("select it.issue_id from ppm_pri_issue_tag it where it.org_id = ? and it.tag_id in ? and it.is_delete = 2", orgId, tagIds))
}

func IssueCondResourceId(issueCond db.Cond, orgId, resourceId int64) {
	//查询tag
	issueCond[consts.TcId+" IN "] = db.In(db.Raw("select ir.issue_id from ppm_pri_issue_relation ir where ir.org_id = ? and ir.relation_id = ? and ir.relation_type = ? and ir.is_delete = 2", orgId, resourceId, consts.IssueRelationTypeResource))
}

func IssueCondFiling(issueCond db.Cond, orgId int64, isFiling int) {
	if isFiling == 3 {
		return
	} else {
		issueCond[consts.TcProjectId+" IN"] = db.In(db.Raw("select id from ppm_pro_project where is_filing = ? and org_id = ? and is_delete = 2", isFiling, orgId))
	}
}

func IssueCondOrderBy(orgId int64, orderByType int) interface{} {
	var orderBy interface{} = nil
	switch orderByType {
	case 1:
		//项目分组
		orderBy = db.Raw("project_id desc, id desc")
	case 2:
		priorities, err := GetPriorityListByType(orgId, consts.PriorityTypeIssue)
		if err != nil {
			log.Error(err)
			orderBy = db.Raw("(select sort from ppm_prs_priority p where p.id = priority_id) asc, plan_end_time asc, id desc")
		} else {
			bo.SortPriorityBo(*priorities)
			orderBySort := ""
			for _, priority := range *priorities {
				orderBySort += fmt.Sprintf(",%d", priority.Id)
			}
			orderBy = db.Raw("FIELD(priority_id" + orderBySort + ")")
		}
	case 3:
		//创建时间降序
		orderBy = db.Raw("create_time desc, id desc")
	case 4:
		//更新时间降序
		orderBy = db.Raw("update_time desc, id desc")
	case 5:
		//按开始时间最早
		orderBy = db.Raw("if (plan_start_time > '1970-02-01 00:00:00',1,0) desc, plan_start_time asc, id desc")
	case 6:
		//按开始时间最晚
		orderBy = db.Raw("plan_start_time desc, id desc")
	case 8:
		//按截止时间最近
		orderBy = db.Raw("plan_end_time desc, id desc")
	case 9:
		//按创建时间最早
		orderBy = db.Raw("create_time asc, id desc")
	case 10:
		//按照sort排序
		orderBy = db.Raw("sort asc, id desc")
	}
	return orderBy
}

//状态类型,1:未完成，2：已完成，3：未开始，4：进行中 5: 已逾期
func IssueCondStatusAssembly(issueCond db.Cond, orgId int64, status int) errs.SystemErrorInfo {
	var statusIds []int64 = nil

	if status == 1 {
		notStartedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeNotStarted)
		if err != nil {
			log.Errorf(getProcessError, err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		processingIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeProcessing)
		if err != nil {
			log.Errorf(getProcessError, err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		statusIds = append(*notStartedIds, *processingIds...)
	} else if status == 2 {
		finishedId, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
		if err != nil {
			log.Errorf(getProcessError, err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		statusIds = *finishedId
	} else if status == 3 {
		notStartedIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeNotStarted)
		if err != nil {
			log.Errorf(getProcessError, err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		statusIds = *notStartedIds
	} else if status == 4 {
		processingIds, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeProcessing)
		if err != nil {
			log.Errorf(getProcessError, err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		statusIds = *processingIds
	} else if status == 5 {
		//已经逾期筛选条件，并且未完成
		nowTime := time.Now()
		//逾期
		issueCond[consts.TcPlanEndTime] = db.Between(consts.BlankElasticityTime, date.Format(nowTime))
		//未完成
		finishedId, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
		if err != nil {
			log.Errorf(getProcessError, err)
			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
		}
		issueCond[consts.TcStatus] = db.NotIn(*finishedId)
	}
	if status != 5 {
		issueCond[consts.TcStatus] = db.In(statusIds)
	}
	return nil
}

func IssueCondRelationMemberAssembly(queryCond db.Cond, input *vo.HomeIssueInfoReq) {
	selectOthersIssueSql := "select ir.issue_id as id from ppm_pri_issue_relation ir where ir.is_delete = 2"
	memberIdsList := &[]interface{}{}
	selectOthersIssueSqlIsJoint := false
	selectOthersIssueSqlJointTag := " and "

	//可以在基本条件封装时做处理
	//if input.OwnerIds != nil{
	//	*memberIdsList = append(*memberIdsList, input.OwnerIds)
	//	if selectOthersIssueSqlIsJoint{
	//		selectOthersIssueSqlJointTag = " or "
	//	}else{
	//		selectOthersIssueSqlIsJoint = !selectOthersIssueSqlIsJoint
	//	}
	//	selectOthersIssueSql += selectOthersIssueSqlJointTag + "(ir.relation_id in ? and ir.relation_type = 1)"
	//}
	if input.ParticipantIds != nil {
		*memberIdsList = append(*memberIdsList, input.ParticipantIds)
		if selectOthersIssueSqlIsJoint {
			selectOthersIssueSqlJointTag = " or "
		} else {
			selectOthersIssueSqlIsJoint = !selectOthersIssueSqlIsJoint
		}
		selectOthersIssueSql += selectOthersIssueSqlJointTag + "(ir.relation_id in ? and ir.relation_type = 2)"
	}
	if input.FollowerIds != nil {
		*memberIdsList = append(*memberIdsList, input.FollowerIds)
		if selectOthersIssueSqlIsJoint {
			selectOthersIssueSqlJointTag = " or "
		} else {
			selectOthersIssueSqlIsJoint = !selectOthersIssueSqlIsJoint
		}
		selectOthersIssueSql += selectOthersIssueSqlJointTag + "(ir.relation_id in ? and ir.relation_type = 3)"
	}
	if selectOthersIssueSqlIsJoint {
		queryCond[consts.TcId] = db.In(db.Raw(selectOthersIssueSql, *memberIdsList...))
	}
}

func IssueCondCombinedCondAssembly(queryCond db.Cond, input *vo.HomeIssueInfoReq, currentUserId int64, orgId int64) errs.SystemErrorInfo {
	if input.CombinedType != nil {
		todayTimeQuantum := times.GetTodayTimeQuantum()
		switch *input.CombinedType {
		//1: 今日指派给我
		case 1:
			queryCond[consts.TcOwnerChangeTime] = db.Between(date.Format(todayTimeQuantum[0]), date.Format(todayTimeQuantum[1]))
			queryCond[consts.TcOwner] = currentUserId
		//2: 最近截止，展示已逾期任务和预计时间小于后天凌晨的任务
		case 2:
			tomorrowTime := todayTimeQuantum[1].Add(time.Duration(60*60*24) * time.Second)
			queryCond[consts.TcPlanEndTime] = db.Between(consts.BlankElasticityTime, date.Format(tomorrowTime))
			queryCond[consts.TcOwner] = currentUserId
		//3: 今日逾期
		case 3:
			nowTime := time.Now()
			todayBegin := date.GetZeroTime(nowTime)
			todayEnd := date.GetZeroTime(nowTime).Add((86400 - 1) * time.Second)
			//今日到期并且尚未完成
			queryCond[consts.TcPlanEndTime] = db.Between(date.Format(todayBegin), date.Format(todayEnd))
			//未完成
			finishedId, err := processfacade.GetProcessStatusIdsRelaxed(orgId, consts.ProcessStatusCategoryIssue, consts.ProcessStatusTypeCompleted)
			if err != nil {
				log.Errorf(getProcessError, err)
				return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
			}
			queryCond[consts.TcStatus] = db.NotIn(*finishedId)
		//4:逾期完成
		case 4:
			//逾期
			queryCond[consts.TcPlanEndTime+" "] = db.Gt(consts.BlankElasticityTime)
			queryCond[consts.TcEndTime] = db.Gt(db.Raw(consts.TcPlanEndTime))
		//5: 即将逾期:预计时间小于后天凌晨的任务
		case 5:
			tomorrowTime := todayTimeQuantum[1].Add(time.Duration(60*60*24) * time.Second)
			queryCond[consts.TcPlanEndTime] = db.Between(date.Format(time.Now()), date.Format(tomorrowTime))
		//6：今日创建的任务
		case 6:
			queryCond[consts.TcCreateTime] = db.Between(date.Format(todayTimeQuantum[0]), date.Format(todayTimeQuantum[1]))
		}
	}

	return nil
}

//增量查询条件封装
func IssueCondLastUpdateTimeCondAssembly(queryCond db.Cond, input *vo.HomeIssueInfoReq) {
	if input.LastUpdateTime != nil {
		//删除is_delete条件，因为增量查询要将删除的变动也查出来
		delete(queryCond, consts.TcIsDelete)
		queryCond[consts.TcUpdateTime] = db.Gte(date.FormatTime(*input.LastUpdateTime))
	}
}
