package domain

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InitProjectStatus(projectBo bo.ProjectBo, operatorId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	orgId := projectBo.OrgId
	projectId := projectBo.Id

	projectType, err1 := GetProjectTypeById(orgId, projectBo.ProjectTypeId)
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err1)
	}

	processLangCode := consts.ProcessLangCodeDefaultProject
	if projectType.LangCode == consts.ProjectTypeLangCodeAgile {
		processLangCode = consts.ProcessLangCodeDefaultAgileTask
	}

	process, err1 := processfacade.GetProcessByLangCodeRelaxed(orgId, processLangCode)
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err1)
	}

	statusBos, err := processfacade.GetProcessStatusListRelaxed(orgId, process.Id)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	pos := make([]po.PpmProProjectRelation, 0, 10)

	for _, statusBo := range *statusBos {
		id, err3 := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectRelation)
		if err3 != nil {
			log.Error(err3)
			return errs.BuildSystemErrorInfo(errs.ApplyIdError)
		}
		projectRelationPo := po.PpmProProjectRelation{
			Id:           id,
			OrgId:        orgId,
			ProjectId:    projectId,
			RelationId:   statusBo.StatusId,
			RelationType: consts.IssueRelationTypeStatus,
			Status:       consts.AppStatusEnable,
			Creator:      operatorId,
			Updator:      operatorId,
			IsDelete:     consts.AppIsNoDelete,
		}
		pos = append(pos, projectRelationPo)
	}

	err2 := dao.InsertProjectRelationBatch(pos, tx...)
	if err2 != nil {
		log.Error(err2)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	return nil
}
