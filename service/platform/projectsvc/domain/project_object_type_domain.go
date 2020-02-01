package domain

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
)

func GetProjectSupportObjectTypeList(projectBo bo.ProjectBo) (*bo.ProjectSupportObjectTypeListBo, errs.SystemErrorInfo) {
	orgId := projectBo.OrgId
	projectId := projectBo.Id

	projectTypeBo, err := GetProjectTypeByLangCode(orgId, consts.ProjectTypeLangCodeNormalTask)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	projectSupportList := make([]*bo.ProjectObjectTypeRestInfoBo, 0, 10)
	iterationSupportList := make([]*bo.ProjectObjectTypeRestInfoBo, 0, 10)

	//通用项目只返回任务, 写死
	if projectBo.ProjectTypeId == projectTypeBo.Id {
		projectSupportList = append(projectSupportList, &bo.ProjectObjectTypeRestInfoBo{
			ID:         0,
			LangCode:   consts.ProjectObjectTypeLangCodeTask,
			Name:       "任务",
			ObjectType: consts.ProjectObjectTypeTask,
		})
	} else {
		projectObjectList, err1 := GetProjectObjectTypeList(orgId, projectId)
		log.Infof("获取组织%d下的项目类型列表为%s", orgId, json.ToJsonIgnoreError(projectObjectList))
		if err1 != nil {
			log.Error(err1)
			return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err1)
		}

		hasIteration := false
		for _, projectObjectType := range *projectObjectList {
			if projectObjectType.LangCode == consts.ProjectObjectTypeLangCodeTestTask {
				continue
			}
			projectObjectTypeRestInfo := &bo.ProjectObjectTypeRestInfoBo{
				ID:         projectObjectType.Id,
				Name:       projectObjectType.Name,
				LangCode:   projectObjectType.LangCode,
				ObjectType: projectObjectType.ObjectType,
			}
			projectSupportList = append(projectSupportList, projectObjectTypeRestInfo)
			if projectObjectType.LangCode != consts.ProjectObjectTypeLangCodeFeature && projectObjectType.LangCode != consts.ProjectObjectTypeLangCodeIteration {
				iterationSupportList = append(iterationSupportList, projectObjectTypeRestInfo)
			}
			if projectObjectType.LangCode == consts.ProjectObjectTypeLangCodeIteration {
				hasIteration = true
			}
		}
		if !hasIteration {
			iterationSupportList = make([]*bo.ProjectObjectTypeRestInfoBo, 0)
		}
	}

	return &bo.ProjectSupportObjectTypeListBo{
		ProjectSupportList:   projectSupportList,
		IterationSupportList: iterationSupportList,
	}, nil
}
