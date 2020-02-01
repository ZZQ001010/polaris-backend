package domain

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/po"
	"strconv"
	"strings"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//流程状态初始化条数
const PrsProcessStatusPrimaryKey = 27
const ProcessStatusSql = consts.TemplateDirPrefix + "ppm_prs_process_status.template"

//默认流程流程状态关联初始化条数
const PrsProcessProcessStatusPrimaryKey = 41
const ProcessProcessStatusSql = consts.TemplateDirPrefix + "ppm_prs_process_process_status.template"

//流程步骤初始化条数
const PrsProcessStepPrimaryKey = 148
const ProcessStepSql = consts.TemplateDirPrefix + "ppm_prs_process_step.template"

//问题来源
const IssueSourcePrimaryKey = 20
const IssueSourceSql = consts.TemplateDirPrefix + "ppm_prs_issue_source.template"

//问题类型
const IssueObjectTypePrimaryKey = 21
const IssueObjectTypeSql = consts.TemplateDirPrefix + "ppm_prs_issue_object_type.template"

const MaxApply = 20

//注册ppm_prs_process_status主键id
func assemblyIdForProcessStatus(orgId int64, contextMap map[string]interface{}) errs.SystemErrorInfo {
	//注册ppm_prs_process_status主键id
	idForProcessStatus := 1
	for i := PrsProcessStatusPrimaryKey; i >= 0; i -= MaxApply {
		applyThis := 0
		if i > MaxApply {
			applyThis = MaxApply
		} else {
			applyThis = i
		}
		ids, err := idfacade.ApplyMultipleIdRelaxed(orgId, (&po.PpmPrsProcessStatus{}).TableName(), "", int64(applyThis))
		if err != nil {
			return err
		}
		for _, v := range ids.Ids {
			contextMap["ProcessStatusId"+strconv.Itoa(idForProcessStatus)] = v.Id
			idForProcessStatus++
		}
	}
	return nil
}

//注册ppm_prs_process_process_status主键id
func assemblyIdForProcessProcessStatus(orgId int64, contextMap map[string]interface{}) errs.SystemErrorInfo {
	idForProcessProcessStatus := 1
	for i := PrsProcessProcessStatusPrimaryKey; i >= 0; i -= MaxApply {
		applyThis := 0
		if i > MaxApply {
			applyThis = MaxApply
		} else {
			applyThis = i
		}
		ids, err := idfacade.ApplyMultipleIdRelaxed(orgId, (&po.PpmPrsProcessProcessStatus{}).TableName(), "", int64(applyThis))
		if err != nil {
			return err
		}
		for _, v := range ids.Ids {
			contextMap["ProcessProcessStatusId"+strconv.Itoa(idForProcessProcessStatus)] = v.Id
			idForProcessProcessStatus++
		}
	}
	return nil
}

//注册ppm_prs_process_step主键id
func assemblyIdForProcessStep(orgId int64, contextMap map[string]interface{}) errs.SystemErrorInfo {
	//注册ppm_prs_process_step主键id
	idForProcessStep := 1
	for i := PrsProcessStepPrimaryKey; i >= 0; i -= MaxApply {
		applyThis := 0
		if i > MaxApply {
			applyThis = MaxApply
		} else {
			applyThis = i
		}
		ids, err := idfacade.ApplyMultipleIdRelaxed(orgId, (&po.PpmPrsProcessStep{}).TableName(), "", int64(applyThis))
		if err != nil {
			return err
		}
		for _, v := range ids.Ids {
			contextMap["ProcessStepId"+strconv.Itoa(idForProcessStep)] = v.Id
			idForProcessStep++
		}
	}
	return nil
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
		ids, err := idfacade.ApplyMultipleIdRelaxed(orgId, consts.TableIssueSource, "", int64(applyThis))
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
		ids, err := idfacade.ApplyMultipleIdRelaxed(orgId, consts.TableIssueObjectType, "", int64(applyThis))
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

func ProcessStatusInit(orgId int64, contextMap map[string]interface{}, tx sqlbuilder.Tx) errs.SystemErrorInfo {

	contextMap["OrgId"] = orgId
	//获取项目流程信息
	processInfo := &[]po.PpmPrsProcess{}
	err := tx.Select(consts.TcLangCode, consts.TcId).From((&po.PpmPrsProcess{}).TableName()).Where(db.Cond{
		consts.TcOrgId: orgId,
	}).All(processInfo)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if len(*processInfo) == 0 {
		return errs.BuildSystemErrorInfo(errs.ProjectNotInit)
	}
	for _, v := range *processInfo {
		contextMap[strings.ReplaceAll(v.LangCode, ".", "")] = v.Id
	}

	//注册ppm_prs_process_status主键id
	//idForProcessStatus := 1
	//for i := PrsProcessStatusPrimaryKey; i >= 0; i -= MaxApply {
	//	applyThis := 0
	//	if i > MaxApply {
	//		applyThis = MaxApply
	//	} else {
	//		applyThis = i
	//	}
	//	ids, err := idservice.ApplyMultipleIdRelaxed(orgId, (&po.PpmPrsProcessStatus{}).TableName(), "", int64(applyThis))
	//	if err != nil {
	//		return err
	//	}
	//	for _, v := range ids.Ids {
	//		context["ProcessStatusId"+strconv.Itoa(idForProcessStatus)] = v.Id
	//		idForProcessStatus++
	//	}
	//}

	//注册ppm_prs_process_status主键id
	err1 := assemblyIdForProcessStatus(orgId, contextMap)

	if err != nil {
		return err1
	}
	//注册ppm_prs_process_process_status主键id
	//idForProcessProcessStatus := 1
	//for i := PrsProcessProcessStatusPrimaryKey; i >= 0; i -= MaxApply {
	//	applyThis := 0
	//	if i > MaxApply {
	//		applyThis = MaxApply
	//	} else {
	//		applyThis = i
	//	}
	//	ids, err := idservice.ApplyMultipleIdRelaxed(orgId, (&po.PpmPrsProcessProcessStatus{}).TableName(), "", int64(applyThis))
	//	if err != nil {
	//		return err
	//	}
	//	for _, v := range ids.Ids {
	//		context["ProcessProcessStatusId"+strconv.Itoa(idForProcessProcessStatus)] = v.Id
	//		idForProcessProcessStatus++
	//	}
	//}

	//注册ppm_prs_process_process_status主键id
	err1 = assemblyIdForProcessProcessStatus(orgId, contextMap)

	if err != nil {
		return err1
	}

	//注册ppm_prs_process_step主键id
	//idForProcessStep := 1
	//for i := PrsProcessStepPrimaryKey; i >= 0; i -= MaxApply {
	//	applyThis := 0
	//	if i > MaxApply {
	//		applyThis = MaxApply
	//	} else {
	//		applyThis = i
	//	}
	//	ids, err := idservice.ApplyMultipleIdRelaxed(orgId, (&po.PpmPrsProcessStep{}).TableName(), "", int64(applyThis))
	//	if err != nil {
	//		return err
	//	}
	//	for _, v := range ids.Ids {
	//		context["ProcessStepId"+strconv.Itoa(idForProcessStep)] = v.Id
	//		idForProcessStep++
	//	}
	//}

	//注册ppm_prs_process_step主键id
	err1 = assemblyIdForProcessStep(orgId, contextMap)

	if err != nil {
		return err1
	}

	//注册ppm_prs_issue_source主键id
	//idForIssueSource := 1
	//for i := IssueSourcePrimaryKey; i >= 0; i -= MaxApply {
	//	applyThis := 0
	//	if i > MaxApply {
	//		applyThis = MaxApply
	//	} else {
	//		applyThis = i
	//	}
	//	ids, err := idservice.ApplyMultipleIdRelaxed(orgId, consts.TableIssueSource, "", int64(applyThis))
	//	if err != nil {
	//		return err
	//	}
	//	for _, v := range ids.Ids {
	//		context["IssueSource"+strconv.Itoa(idForIssueSource)] = v.Id
	//		idForIssueSource++
	//	}
	//}

	//err1 = assemblyIdForIssueSource(orgId, contextMap)
	//
	//if err != nil {
	//	return err1
	//}

	//注册ppm_prs_issue_object_type主键id
	//idForIssueType := 1
	//for i := IssueObjectTypePrimaryKey; i >= 0; i -= MaxApply {
	//	applyThis := 0
	//	if i > MaxApply {
	//		applyThis = MaxApply
	//	} else {
	//		applyThis = i
	//	}
	//	ids, err := idservice.ApplyMultipleIdRelaxed(orgId, consts.TableIssueObjectType, "", int64(applyThis))
	//	if err != nil {
	//		return err
	//	}
	//	for _, v := range ids.Ids {
	//		context["IssueObjectType"+strconv.Itoa(idForIssueType)] = v.Id
	//		idForIssueType++
	//	}
	//}

	//err1 = assemblyIdForIssueType(orgId, contextMap)
	//
	//if err != nil {
	//	return err1
	//}

	err2 := util.ReadAndWrite(ProcessStatusSql, contextMap, tx)
	if err2 != nil {
		return err2
	}
	err2 = util.ReadAndWrite(ProcessProcessStatusSql, contextMap, tx)
	if err2 != nil {
		return err2
	}
	err2 = util.ReadAndWrite(ProcessStepSql, contextMap, tx)
	if err2 != nil {
		return err2
	}
	//err2 = util.ReadAndWrite(IssueSourceSql, contextMap, tx)
	//if err2 != nil {
	//	return err2
	//}
	//err2 = util.ReadAndWrite(IssueObjectTypeSql, contextMap, tx)
	//if err2 != nil {
	//	return err2
	//}

	return nil
}
