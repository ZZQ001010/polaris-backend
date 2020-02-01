package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
)

func ProjectIssueRelatedStatus(orgId int64, input vo.ProjectIssueRelatedStatusReq) ([]*vo.HomeIssueStatusInfo, errs.SystemErrorInfo) {

	projectId := input.ProjectID
	projectObjectTypeId := input.ProjectObjectTypeID

	_, err := domain.LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IllegalityProject, err)
	}

	homeStatusInfoBos, err := domain.GetProjectRelatedStatus(orgId, projectId, projectObjectTypeId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	results := &[]*vo.HomeIssueStatusInfo{}
	err1 := copyer.Copy(homeStatusInfoBos, results)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *results, nil
}
