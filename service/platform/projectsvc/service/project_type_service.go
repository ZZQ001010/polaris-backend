package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
)

func ProjectTypes(orgId int64) ([]*vo.ProjectType, errs.SystemErrorInfo) {
	projectTypeBo, err := domain.GetProjectTypeList(orgId)
	if err != nil {
		return nil, err
	}

	projectTypeVo := &[]*vo.ProjectType{}
	copyErr := copyer.Copy(projectTypeBo, projectTypeVo)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, copyErr)
	}

	return *projectTypeVo, nil
}
