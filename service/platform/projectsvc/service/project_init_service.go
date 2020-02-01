package service

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"strconv"
	"upper.io/db.v3/lib/sqlbuilder"
)

//问题来源
const IssueSourcePrimaryKey = 20
const IssueSourceSql = consts.TemplateDirPrefix + "ppm_prs_issue_source.template"

//问题类型
const IssueObjectTypePrimaryKey = 21
const IssueObjectTypeSql = consts.TemplateDirPrefix + "ppm_prs_issue_object_type.template"

const MaxApply = 20

const ProjectObjectFeatureId = 2
const ProjectObjectDemandId = 3
const ProjectObjectBugId = 4
const ProjectObjectTaskId = 5

func ProjectInit(orgId int64, tx sqlbuilder.Tx) (map[string]interface{}, errs.SystemErrorInfo) {
	//逻辑保留，暂时不需要传递上下文
	contextMap := map[string]interface{}{}

	contextMap["OrgId"] = orgId
	//初始化sql写死
	contextMap["ProjectObjectFeatureId"] = ProjectObjectFeatureId
	contextMap["ProjectObjectDemandId"] = ProjectObjectDemandId
	contextMap["ProjectObjectBugId"] = ProjectObjectBugId
	contextMap["ProjectObjectTaskId"] = ProjectObjectTaskId

	//初始化优先级
	err := domain.PriorityInit(orgId, tx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//封装任务来源id
	err1 := assemblyIdForIssueSource(orgId, contextMap)
	if err != nil {
		log.Error(err)
		return nil, err1
	}

	//封装任务类型id
	err1 = assemblyIdForIssueType(orgId, contextMap)
	if err != nil {
		log.Error(err)
		return nil, err1
	}

	//执行初始化sql
	err2 := util.ReadAndWrite(IssueSourceSql, contextMap, tx)
	if err2 != nil {
		return nil, err2
	}
	err2 = util.ReadAndWrite(IssueObjectTypeSql, contextMap, tx)
	if err2 != nil {
		return nil, err2
	}
	return contextMap, nil
}

//注册ppm_prs_issue_source主键id
func assemblyIdForIssueSource(orgId int64, contextMap map[string]interface{}) errs.SystemErrorInfo {

	idForIssueSource := 1
	for i := IssueSourcePrimaryKey; i >= 0; i -= MaxApply {
		applyThis := 0
		if i > MaxApply {
			applyThis = MaxApply
		} else {
			applyThis = i
		}
		ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueSource, applyThis)
		if err != nil {
			return err
		}
		for _, v := range ids.Ids {
			contextMap["IssueSource"+strconv.Itoa(idForIssueSource)] = v.Id
			idForIssueSource++
		}
	}
	return nil
}

//注册ppm_prs_issue_object_type主键id
func assemblyIdForIssueType(orgId int64, contextMap map[string]interface{}) errs.SystemErrorInfo {
	idForIssueType := 1
	for i := IssueObjectTypePrimaryKey; i >= 0; i -= MaxApply {
		applyThis := 0
		if i > MaxApply {
			applyThis = MaxApply
		} else {
			applyThis = i
		}
		ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueObjectType, applyThis)
		if err != nil {
			return err
		}
		for _, v := range ids.Ids {
			contextMap["IssueObjectType"+strconv.Itoa(idForIssueType)] = v.Id
			idForIssueType++
		}
	}
	return nil
}
