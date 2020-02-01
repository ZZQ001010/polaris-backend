package service

import (
	"github.com/galaxy-book/common/core/errors"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/processvo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"sort"
	"strconv"
	"upper.io/db.v3"
)

func ProjectObjectTypeList(orgId int64, page uint, size uint, params *vo.ProjectObjectTypesReq) (*vo.ProjectObjectTypeList, errs.SystemErrorInfo) {
	//cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//orgId := cacheUserInfo.OrgId

	cond := db.Cond{}
	if params != nil {
		if params.ObjectType != 0 {
			cond[consts.TcObjectType] = params.ObjectType
		}
		if len(params.Ids) > 0 {
			cond[consts.TcId] = db.In(params.Ids)
		}
	}
	cond[consts.TcIsDelete] = consts.AppIsNoDelete
	cond[consts.TcOrgId] = db.In([]int64{orgId, 0})

	bos, total, err := domain.GetProjectObjectTypeBoList(page, size, cond)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}

	resultList := &[]*vo.ProjectObjectType{}
	copyErr := copyer.Copy(bos, resultList)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &vo.ProjectObjectTypeList{
		Total: total,
		List:  *resultList,
	}, nil
}

func CreateProjectObjectType(orgId, currentUserId int64, input vo.CreateProjectObjectTypeReq) (*vo.Void, errs.SystemErrorInfo) {

	//用户角色权限校验
	err := domain.AuthProject(orgId, currentUserId, input.ProjectID, consts.RoleOperationPathOrgProjectObjectType, consts.RoleOperationCreate)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}
	//length := len(input.Name)
	//if length == 0 || length > 30 {
	//	return nil, errs.InvalidProjectObjectTypeName
	//}
	isNameRight := format.VerifyProjectObjectTypeNameFormat(input.Name)
	if !isNameRight {
		return nil, errs.InvalidProjectObjectTypeName
	}
	//校验任务类型是否是普通任务
	projectTypeBo, err2 := AuthProjectType(orgId, input.ProjectID)
	if err2 != nil {
		return nil, err2
	}

	//插入projectObjectType
	id, err := insertProjectObjectType(input, orgId, currentUserId, projectTypeBo)

	if err != nil {
		return nil, err
	}
	DeleteProjectExcel(orgId, input.ProjectID)
	return &vo.Void{
		ID: id,
	}, nil
}

func AuthProjectType(orgId, projectId int64) (*bo.ProjectTypeBo, errs.SystemErrorInfo) {

	projectTypeBo, err := domain.GetProjectTypeByLangCode(orgId, consts.ProjectTypeLangCodeNormalTask)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	//获取项目 里面包含projectTypeId 和普通任务的projectType的Id 做对比
	projectAuthBo, err := domain.LoadProjectAuthBo(orgId, projectId)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if projectAuthBo.ProjectType != projectTypeBo.Id {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectTypeNormalError)
	}

	return projectTypeBo, nil
}

func insertProjectObjectType(input vo.CreateProjectObjectTypeReq, orgId, currentUserId int64, projectType *bo.ProjectTypeBo) (int64, errors.SystemErrorInfo) {

	boMap, existErr := domain.CheckSameProjectObjectTypeName(orgId, input.ProjectID, input.Name, nil)

	if existErr != nil {
		return 0, existErr
	}

	entity := &bo.ProjectObjectTypeBo{}
	copyErr := copyer.Copy(input, entity)
	if copyErr != nil {
		return 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectObjectType)
	if err != nil {
		return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	entity.Id = id
	entity.Name = input.Name
	entity.OrgId = orgId
	entity.ObjectType = input.ObjectType
	entity.Creator = currentUserId
	entity.Updator = currentUserId

	//获取这个通用项目之前的projectObjectType的sort
	if input.BeforeID != 0 {

		if v, ok := boMap[input.BeforeID]; ok {
			if data, ok := v.(bo.ProjectObjectTypeBo); ok == true {
				entity.Sort = data.Sort + 1
			}
		} else {
			return 0, errs.BuildSystemErrorInfo(errs.ProjectTypeNotExist)
		}

	}

	//组装关联表数据
	projectObjectTypeProcessBo, err := assemblyPpmProjectObjectTypeProcess(projectType, orgId, entity, input.ProjectID)

	if err != nil {
		return 0, err
	}
	//赋值用于 事务嵌套
	entity.PpmPrsProjectObjectTypeProcessBo = projectObjectTypeProcessBo

	err1 := domain.CreateProjectObjectType(entity)
	if err1 != nil {

		log.Error(err1)
		return 0, err1
	}

	domain.ClearProjectObjectTypeList(orgId, input.ProjectID)

	return id, nil
}

func assemblyPpmProjectObjectTypeProcess(typeBo *bo.ProjectTypeBo, userOrgId int64, entity *bo.ProjectObjectTypeBo, projectId int64) (bo.PpmPrsProjectObjectTypeProcessBo, errors.SystemErrorInfo) {

	bo := bo.PpmPrsProjectObjectTypeProcessBo{}
	//系统初始化的数据 orgId 就是0 typeBo的orgId =0  跟着模板走的

	processResp := processfacade.GetProcessByLangCode(processvo.GetProcessByLangCodeReqVo{
		OrgId:    typeBo.OrgId,
		LangCode: consts.ProcessLangCodeDefaultTask})

	if processResp.Err.Failure() {
		logger.GetDefaultLogger().Error("Rollback:" + processResp.Err.Message)
		return bo, processResp.Err.Error()
	}
	//默认任务流程的processId
	processId := processResp.ProcessBo.Id

	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectObjectTypeProcess)
	if err != nil {
		return bo, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	//赋值时候用用户自己的orgId
	bo.Id = id
	bo.OrgId = userOrgId
	bo.ProjectId = projectId
	bo.ProcessId = processId
	bo.ProjectObjectTypeId = entity.Id
	bo.Creator = entity.Creator
	bo.Updator = entity.Updator

	return bo, nil
}

func UpdateProjectObjectType(orgId, currentUserId int64, input vo.UpdateProjectObjectTypeReq) (*vo.Void, errs.SystemErrorInfo) {

	//用户角色权限校验
	err := domain.AuthProject(orgId, currentUserId, input.ProjectID, consts.RoleOperationPathOrgProjectObjectType, consts.RoleOperationModify)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	//length := len(input.Name)
	//if length == 0 || length > 30 {
	//	return nil, errs.InvalidProjectObjectTypeName
	//}
	isNameRight := format.VerifyProjectObjectTypeNameFormat(input.Name)
	if !isNameRight {
		return nil, errs.InvalidProjectObjectTypeName
	}
	////校验任务类型是否是普通任务
	//_, err2 := AuthProjectType(orgId, input.ProjectID)
	//if err2 != nil {
	//	return nil, err2
	//}

	projectObjectTypeId := input.ID

	boMap, existErr := domain.CheckSameProjectObjectTypeName(orgId, input.ProjectID, input.Name, &projectObjectTypeId)

	if existErr != nil {
		return nil, existErr
	}

	//查询原本的projectObjectType 如果sort不一样把之后所有的projectObjectType的sort都+1

	orginalProjectObjectTypeBo := &bo.ProjectObjectTypeBo{}

	if v, ok := boMap[input.ID]; ok {
		if data, ok := v.(bo.ProjectObjectTypeBo); ok == true {
			orginalProjectObjectTypeBo = &data
		}
	} else {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectTypeNotExist)
	}

	//orginalProjectObjectTypeBo, err2 := domain.GetProjectObjectTypeBo(input.ID)
	//if err2 != nil {
	//	log.Error(err2)
	//	return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	//}

	//目前这里更改正常逻辑上只能改名字和排序
	entity := &bo.ProjectObjectTypeBo{}
	copyErr := copyer.Copy(input, entity)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	entity.Updator = currentUserId
	entity.OrgId = orgId
	entity.OrginalSort = orginalProjectObjectTypeBo.Sort

	//更新操作
	err1 := domain.UpdateProjectObjectType(entity, input.ProjectID)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	DeleteProjectExcel(orgId, input.ProjectID)
	//删除时候清除缓存
	cacheErr := domain.ClearProjectObjectTypeList(orgId, input.ProjectID)
	if cacheErr != nil {
		log.Error(cacheErr)
	}

	return &vo.Void{
		ID: input.ID,
	}, nil
}

func DeleteProjectObjectType(orgId, currentUserId int64, input vo.DeleteProjectObjectTypeReq) (*vo.Void, errs.SystemErrorInfo) {
	projectObjectTypeId := input.ID

	//用户角色权限校验
	err := domain.AuthProject(orgId, currentUserId, input.ProjectID, consts.RoleOperationPathOrgProjectObjectType, consts.RoleOperationDelete)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	//判断是否有任务挂在这个projectObjectType下面有就不可以删除
	issueExist, err := domain.GetCountIssueByProjectObjectTypeId(projectObjectTypeId, input.ProjectID, orgId)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectObjectTypeDeleteFailExistIssue, err)
	}

	if issueExist {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectObjectTypeDeleteFailExistIssue)
	}

	bo, err1 := domain.GetProjectObjectTypeBo(projectObjectTypeId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	//判断是否是最后一个是的话不给删除
	err = domain.JudgeLastProjectObjectType(orgId, input.ProjectID, projectObjectTypeId)

	if err != nil {
		return nil, err
	}

	uuid := uuid.NewUuid()

	strId := strconv.FormatInt(projectObjectTypeId, 10)
	suc, err5 := cache.TryGetDistributedLock(consts.ProjectObjectTypeLockKey+strId, uuid)

	if err5 != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
	}
	if !suc {
		log.Errorf("删除项目泳道异常 %q 原因:同时删除", strId)
		return nil, errs.BuildSystemErrorInfo(errs.OrgNotInitError)
	}
	defer func() {
		if _, lockErr := cache.ReleaseDistributedLock(consts.FeiShuCorpInitKey+strId, uuid); lockErr != nil {
			log.Error(lockErr)
		}
	}()

	err = domain.JudgeLastProjectObjectType(orgId, input.ProjectID, projectObjectTypeId)

	if err != nil {
		return nil, err
	}

	err2 := domain.DeleteProjectObjectType(bo, input.ProjectID, projectObjectTypeId, currentUserId)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err2)
	}
	DeleteProjectExcel(orgId, input.ProjectID)
	//删除时候清除缓存
	cacheErr := domain.ClearProjectObjectTypeList(orgId, input.ProjectID)
	if cacheErr != nil {
		log.Error(cacheErr)
	}

	return &vo.Void{
		ID: projectObjectTypeId,
	}, nil
}

func ProjectObjectTypesWithProject(orgId, projectId int64) (*vo.ProjectObjectTypeWithProjectList, errs.SystemErrorInfo) {

	bos, err := domain.ProjectObjectTypesWithProjectByOrder(orgId, projectId, "sort asc")

	sorters := bo.ProjectObjectTypeBoSorter{}

	for _, bo := range *bos {
		sorters = append(sorters, bo)
	}

	//排序
	sort.Sort(sorters)

	if err != nil {
		return nil, err
	}

	resultList := &[]*vo.ProjectObjectType{}
	copyErr := copyer.Copy(bos, resultList)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &vo.ProjectObjectTypeWithProjectList{
		List: *resultList,
	}, nil
}

func ProjectSupportObjectTypes(orgId int64, input vo.ProjectSupportObjectTypeListReq) (*vo.ProjectSupportObjectTypeListResp, errs.SystemErrorInfo) {

	projectBo, err := domain.GetProjectInfo(input.ProjectID, orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	projectSupportObjectTypeListBo, err1 := domain.GetProjectSupportObjectTypeList(projectBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectTypeDomainError, err1)
	}

	resp := &vo.ProjectSupportObjectTypeListResp{}

	err2 := copyer.Copy(projectSupportObjectTypeListBo, resp)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return resp, nil
}
