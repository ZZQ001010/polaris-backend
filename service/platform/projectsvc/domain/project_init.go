package domain

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//初始化项目类型
func ProjectTypeInit(orgId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	logger.GetDefaultLogger().Infof("start:组织项目类型初始化")
	//获取默认项目流程id
	processInfo, err := processfacade.GetProcessBoRelaxed(db.Cond{
		consts.TcIsDelete:  db.Eq(consts.AppIsNoDelete),
		consts.TcIsDefault: db.Eq(1),
		consts.TcType:      db.Eq(1),
	})
	if err != nil {
		//默认项目流程id不存在
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	count, projectErr := mysql.SelectCountByCond((&po.PpmPrsProjectType{}).TableName(), db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	})
	if projectErr != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, projectErr)
	}
	if count != 0 {
		logger.GetDefaultLogger().Infof("组织项目类型已初始化")
		return nil
	}

	projectType := []interface{}{}
	for _, v := range po.ProjectType {
		id, err := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectType{}).TableName())
		if err != nil {
			return err
		}
		temp := v
		temp.OrgId = orgId
		temp.Id = id
		temp.DefaultProcessId = processInfo.Id
		projectType = append(projectType, temp)
	}

	insertErr := mysql.TransBatchInsert(tx, &po.PpmPrsProjectType{}, projectType)
	if insertErr != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertErr)
	}
	logger.GetDefaultLogger().Infof("success:组织项目类型初始化")

	return nil
}

//初始化项目对象类型
func ProjectObjectTypeInit(orgId int64, contextMap *map[string]interface{}, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	context := *contextMap

	logger.GetDefaultLogger().Infof("start:项目对象类型初始化")
	count, err := mysql.SelectCountByCond((&po.PpmPrsProjectObjectType{}).TableName(), db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if count != 0 {
		logger.GetDefaultLogger().Infof("项目对象类型已初始化")
		return nil
	}

	projectObjectTypeInsert := []interface{}{}
	for k, v := range po.ProjectObjectType {
		id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectObjectType)
		if err != nil {
			log.Error(err)
			return err
		}
		temp := v
		temp.OrgId = orgId
		temp.Id = id
		projectObjectTypeInsert = append(projectObjectTypeInsert, temp)

		switch k {
		case "feature":
			context["ProjectObjectFeatureId"] = id
		case "demand":
			context["ProjectObjectDemandId"] = id
		case "task":
			context["ProjectObjectTaskId"] = id
		case "bug":
			context["ProjectObjectBugId"] = id
		}
	}
	err = mysql.TransBatchInsert(tx, &po.PpmPrsProjectObjectType{}, projectObjectTypeInsert)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	logger.GetDefaultLogger().Infof("success:项目对象类型初始化")

	return nil
}

//初始化默认流程
func ProcessInit(orgId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	logger.GetDefaultLogger().Infof("start:默认流程初始化")
	err := processfacade.InitProcessRelaxed(orgId)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	logger.GetDefaultLogger().Infof("success:默认流程初始化")

	return nil
}

//初始化项目类型与项目对象类型关联
func ProjectTypeProjectObjectTypeInit(orgId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	logger.GetDefaultLogger().Infof("start:初始化项目类型与项目对象类型关联")
	count, err := mysql.SelectCountByCond((&po.PpmPrsProjectTypeProjectObjectType{}).TableName(), db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if count != 0 {
		logger.GetDefaultLogger().Infof("项目类型与项目对象类型关联已初始化")
		return nil
	}

	//获取项目类型
	projectTypeNormalTask, projectTypeAgile, err := getProjectType(orgId, tx)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}

	//获取项目对象类型
	var projectObjectTypeRes = map[string]int64{
		"bug":       0,
		"demand":    0,
		"feature":   0,
		"iteration": 0,
		"task":      0,
		"test_task": 0,
	}
	err = getProjectObjectType(&projectObjectTypeRes, tx, orgId)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}

	//获取默认流程
	var processRes = map[string]int64{
		"test_task":  0,
		"bug":        0,
		"task":       0,
		"demand":     0,
		"feature":    0,
		"agile_task": 0,
		"iteration":  0,
	}
	err = getProcess(&processRes, tx, orgId)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}

	insert := &[]interface{}{}

	//普通任务项目+任务+任务流程
	err = insertProjectTypeProjectObjectType(insert, projectTypeNormalTask, projectObjectTypeRes["task"], processRes["task"], orgId)
	if err != nil {
		log.Error(err)
		return nil
	}
	//敏捷研发项目+迭代+迭代流程
	err = insertProjectTypeProjectObjectType(insert, projectTypeAgile, projectObjectTypeRes["iteration"], processRes["iteration"], orgId)
	if err != nil {
		log.Error(err)
		return nil
	}
	//敏捷研发项目+特性+特性流程
	err = insertProjectTypeProjectObjectType(insert, projectTypeAgile, projectObjectTypeRes["feature"], processRes["feature"], orgId)
	if err != nil {
		log.Error(err)
		return nil
	}
	//敏捷研发项目+需求+需求流程
	err = insertProjectTypeProjectObjectType(insert, projectTypeAgile, projectObjectTypeRes["demand"], processRes["demand"], orgId)
	if err != nil {
		log.Error(err)
		return nil
	}
	//敏捷研发项目+任务+敏捷项目任务流程
	err = insertProjectTypeProjectObjectType(insert, projectTypeAgile, projectObjectTypeRes["task"], processRes["agile_task"], orgId)
	if err != nil {
		log.Error(err)
		return nil
	}
	//敏捷研发项目+缺陷+缺陷流程
	err = insertProjectTypeProjectObjectType(insert, projectTypeAgile, projectObjectTypeRes["bug"], processRes["bug"], orgId)
	if err != nil {
		log.Error(err)
		return nil
	}
	//敏捷研发项目+测试任务+测试任务流程
	err = insertProjectTypeProjectObjectType(insert, projectTypeAgile, projectObjectTypeRes["test_task"], processRes["test_task"], orgId)
	if err != nil {
		log.Error(err)
		return nil
	}
	if len(*insert) == 0 {
		logger.GetDefaultLogger().Infof("初始化项目类型与项目对象类型关联：没有需要插入的数据")
		return nil
	}

	err = mysql.TransBatchInsert(tx, &po.PpmPrsProjectTypeProjectObjectType{}, *insert)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	logger.GetDefaultLogger().Infof("success:初始化项目类型与项目对象类型关联")

	return nil
}

func getProcess(processRes *map[string]int64, tx sqlbuilder.Tx, orgId int64) errs.SystemErrorInfo {
	err := processfacade.AssignValueToFieldRelaxed(processRes, orgId)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	return nil
}

func getProjectObjectType(projectObjectTypeRes *map[string]int64, tx sqlbuilder.Tx, orgId int64) errs.SystemErrorInfo {
	//获取项目对象类型
	projectObjectType := &[]po.PpmPrsProjectObjectType{}
	err := tx.Select(consts.TcId, consts.TcLangCode).From((&po.PpmPrsProjectObjectType{}).TableName()).Where(db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}).All(projectObjectType)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if len(*projectObjectType) == 0 {
		logger.GetDefaultLogger().Infof("初始化项目类型与项目对象类型关联：无对应项目对象类型")
		return errs.BuildSystemErrorInfo(errs.ProjectObjectTypeNotExist)
	}
	for _, v := range *projectObjectType {
		switch v.LangCode {
		case po.ProjectObjectTypeBug.LangCode:
			(*projectObjectTypeRes)["bug"] = v.Id
		case po.ProjectObjectTypeDemand.LangCode:
			(*projectObjectTypeRes)["demand"] = v.Id
		case po.ProjectObjectTypeFeature.LangCode:
			(*projectObjectTypeRes)["feature"] = v.Id
		case po.ProjectObjectTypeIteration.LangCode:
			(*projectObjectTypeRes)["iteration"] = v.Id
		case po.ProjectObjectTypeTask.LangCode:
			(*projectObjectTypeRes)["task"] = v.Id
		case po.ProjectObjectTypeTestTask.LangCode:
			(*projectObjectTypeRes)["test_task"] = v.Id
		}
	}

	return nil
}

func getProjectType(orgId int64, tx sqlbuilder.Tx) (int64, int64, errs.SystemErrorInfo) {
	var projectTypeNormalTask,
		projectTypeAgile int64
	//获取项目类型
	projectType := &[]po.PpmPrsProjectType{}
	err := tx.Select(consts.TcId, consts.TcLangCode).From((&po.PpmPrsProjectType{}).TableName()).Where(db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    db.In([]int64{orgId, 0}),
	}).All(projectType)
	if err != nil {
		log.Error(err)
		return projectTypeNormalTask, projectTypeAgile, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if len(*projectType) == 0 {
		logger.GetDefaultLogger().Infof("初始化项目类型与项目对象类型关联：无对应项目类型")
		return projectTypeNormalTask, projectTypeAgile, errs.BuildSystemErrorInfo(errs.ProjectTypeNotExist)
	}

	for _, v := range *projectType {
		if v.LangCode == po.ProjectTypeNormalTask.LangCode {
			projectTypeNormalTask = v.Id
		} else if v.LangCode == po.ProjectTypeAgile.LangCode {
			projectTypeAgile = v.Id
		}
	}

	return projectTypeNormalTask, projectTypeAgile, nil
}

func insertProjectTypeProjectObjectType(insert *[]interface{}, projectTypeId int64, projectObjectTypeId int64, processId int64, orgId int64) errs.SystemErrorInfo {
	if projectTypeId == 0 || projectObjectTypeId == 0 || processId == 0 {
		return nil
	}
	id, err := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectTypeProjectObjectType{}).TableName())
	if err != nil {
		log.Error(err)
		return err
	}
	*insert = append(*insert, po.PpmPrsProjectTypeProjectObjectType{
		Id:                  id,
		OrgId:               orgId,
		ProjectTypeId:       projectTypeId,
		ProjectObjectTypeId: projectObjectTypeId,
		DefaultProcessId:    processId,
		IsReadonly:          1,
	})

	return nil
}

//优先级初始化
func PriorityInit(orgId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	logger.GetDefaultLogger().Infof("start:优先级初始化")
	err := InitPriority(orgId, tx)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	logger.GetDefaultLogger().Infof("success:优先级初始化")
	return nil
}
