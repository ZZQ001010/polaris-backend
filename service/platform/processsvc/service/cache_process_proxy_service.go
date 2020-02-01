package service

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/domain"
)

func GetProcessByLangCode(orgId int64, langCode string) (*bo.ProcessBo, errs.SystemErrorInfo) {
	return domain.GetProcessByLangCode(orgId, langCode)
}

func GetNextProcessStepStatusList(orgId, processId, startStatusId int64) (*[]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	return domain.GetNextProcessStepStatusList(orgId, processId, startStatusId)
}
