package processfacade

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/processvo"
	"upper.io/db.v3"
)

func AssignValueToFieldRelaxed(processRes *map[string]int64, orgId int64) errs.SystemErrorInfo {
	respVo := AssignValueToField(processvo.AssignValueToFieldReqVo{OrgId: orgId, ProcessRes: processRes})
	if respVo.Failure() {
		return respVo.Error()
	}
	return nil
}

func GetDefaultProcessStatusIdRelaxed(orgId int64, processId int64, category int) (int64, errs.SystemErrorInfo) {
	respVo := GetDefaultProcessStatusId(processvo.GetDefaultProcessIdReqVo{OrgId: orgId, ProcessId: processId, Category: category})
	if respVo.Failure() {
		return 0, respVo.Error()
	}
	return respVo.ProcessId, nil
}

func GetNextProcessStepStatusListRelaxed(orgId, processId, startStatusId int64) (*[]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	respVo := GetNextProcessStepStatusList(processvo.GetNextProcessStepStatusListReqVo{OrgId: orgId, ProcessId: processId, StartStatusId: startStatusId})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.CacheProcessStatus, nil
}

func GetProcessBoRelaxed(cond db.Cond) (bo.ProcessBo, errs.SystemErrorInfo) {
	respVo := GetProcessBo(processvo.GetProcessBoReqVo{Cond: cond})
	if respVo.Failure() {
		return bo.ProcessBo{}, respVo.Error()
	}
	return respVo.ProcessBo, nil
}

func GetProcessByLangCodeRelaxed(orgId int64, langCode string) (*bo.ProcessBo, errs.SystemErrorInfo) {
	respVo := GetProcessByLangCode(processvo.GetProcessByLangCodeReqVo{OrgId: orgId, LangCode: langCode})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.ProcessBo, nil
}

func GetProcessInitStatusIdRelaxed(orgId, projectId, projectObjectTypeId int64, category int) (int64, errs.SystemErrorInfo) {
	respVo := GetProcessInitStatusId(processvo.GetProcessInitStatusIdReqVo{OrgId: orgId, ProjectId: projectId, ProjectObjectTypeId: projectObjectTypeId, Category: category})
	if respVo.Failure() {
		return 0, respVo.Error()
	}
	return respVo.ProcessInitStatusId, nil
}

func GetProcessStatusByCategoryRelaxed(orgId int64, statusId int64, category int) (*bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	respVo := GetProcessStatusByCategory(processvo.GetProcessStatusByCategoryReqVo{OrgId: orgId, StatusId: statusId, Category: category})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.CacheProcessStatusBo, nil
}

func GetProcessStatusIdsRelaxed(orgId int64, category int, typ int) (*[]int64, errs.SystemErrorInfo) {
	respVo := GetProcessStatusIds(processvo.GetProcessStatusIdsReqVo{OrgId: orgId, Typ: typ, Category: category})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.ProcessStatusIds, nil
}

func GetProcessStatusListRelaxed(orgId, processId int64) (*[]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	respVo := GetProcessStatusList(processvo.GetProcessStatusListReqVo{OrgId: orgId, ProcessId: processId})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.ProcessStatusBoList, nil
}

func GetProcessStatusListByCategoryRelaxed(orgId int64, category int) ([]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	respVo := GetProcessStatusListByCategory(processvo.GetProcessStatusListByCategoryReqVo{OrgId: orgId, Category: category})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.CacheProcessStatusBoList, nil
}

func InitProcessRelaxed(orgId int64) errs.SystemErrorInfo {
	respVo := InitProcess(processvo.InitProcessReqVo{OrgId: orgId})
	if respVo.Failure() {
		return respVo.Error()
	}
	return nil
}

func ProcessStatusInitRelaxed(orgId int64, contextMap map[string]interface{}) errs.SystemErrorInfo {
	respVo := ProcessStatusInit(processvo.ProcessStatusInitReqVo{OrgId: orgId, ContextMap: contextMap})
	if respVo.Failure() {
		return respVo.Error()
	}
	return nil
}

func GetProcessStatusRelaxed(orgId int64, id int64) (*bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	respVo := GetProcessStatus(processvo.GetProcessStatusReqVo{OrgId: orgId, Id: id})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.CacheProcessStatusBo, nil
}

func GetProcessByIdRelaxed(orgId, id int64) (*bo.ProcessBo, errs.SystemErrorInfo) {
	respVo := GetProcessById(processvo.GetProcessByIdReqVo{OrgId: orgId, Id: id})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.ProcessBo, nil
}
