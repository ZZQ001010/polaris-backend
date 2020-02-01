package domain

import (
	"errors"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/date"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/facade/resourcefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"gopkg.in/fatih/set.v0"
	"strings"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//func CreateProject(entity bo.ProjectBo, orgId, currentUserId int64, resourcePath string, resourceType int, memberIds []int64, followerIds []int64) (bo.ProjectBo, []int64, errs.SystemErrorInfo) {
func CreateProject(entity bo.ProjectBo, orgId, currentUserId int64, input vo.CreateProjectReq) (bo.ProjectBo, []int64, errs.SystemErrorInfo) {
	memberEntities, addedMemberIds, err := HandleProjectMember(orgId, currentUserId, entity.Owner, entity.Id, input.MemberIds, input.FollowerIds)
	if err != nil {
		return entity, nil, err
	}

	//查询资源是否已存在
	var resourceId int64
	//if resourcePath != "" {
	//	respVo := resourcefacade.GetIdByPath(
	//		resourcevo.GetIdByPathReqVo{
	//			OrgId:        orgId,
	//			ResourceType: resourceType,
	//			ResourcePath: resourcePath,
	//		})
	//	if !respVo.Failure() {
	//		log.Error(respVo.Message)
	//		resourceId = respVo.ResourceId
	//		entity.ResourceId = resourceId
	//	}
	//}

	dealResourcePath(&entity, orgId, input.ResourcePath, input.ResourceType, &resourceId)

	err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//插入项目成员
		insertErr := insertMemberEntities(tx, memberEntities)

		if insertErr != nil {
			return insertErr
		}

		//插入资源

		err = insertSource(&entity, input.ResourcePath, resourceId, tx, orgId, currentUserId, input.ResourceType)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}

		//插入项目
		//insert := &po.PpmProProject{}
		//err1 := copyer.Copy(entity, insert)
		//if err1 != nil {
		//	return err1
		//}
		//_, insertProjectErr := tx.Collection(consts.TableProject).Insert(insert)
		//if insertProjectErr != nil {
		//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertProjectErr)
		//}
		projectError := insertProject(tx, entity)

		if projectError != nil {
			return projectError
		}

		//插入项目公告
		err = insertProjectDetail(&entity, orgId, currentUserId, input.IsSyncOutCalendar)

		if err != nil {
			return err
		}

		//创建项目对象流程关联
		err = CreateProjectObjectTypeProcess(tx, entity.ProjectTypeId, orgId, entity.Id, currentUserId)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
		}

		//创建项目状态关联
		err = InitProjectStatus(entity, currentUserId, tx)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
		}

		return nil
	})
	if err1 != nil {
		log.Errorf("tx.Commit(): %q\n", err1)
		return entity, addedMemberIds, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}

	return entity, addedMemberIds, nil
}

func dealResourcePath(entity *bo.ProjectBo, orgId int64, resourcePath string, resourceType int, resourceId *int64) {

	if resourcePath != "" {
		respVo := resourcefacade.GetIdByPath(
			resourcevo.GetIdByPathReqVo{
				OrgId:        orgId,
				ResourceType: resourceType,
				ResourcePath: resourcePath,
			})
		if !respVo.Failure() {
			log.Error(respVo.Message)
			*resourceId = respVo.ResourceId
			(*entity).ResourceId = *resourceId
		}
	}
}

func insertProjectDetail(entity *bo.ProjectBo, orgId, currentUserId int64, isSyncOutCalendar *int) errs.SystemErrorInfo {

	detailId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectDetail)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	IsSyncOutCalendar := consts.IsSyncOutCalendar
	IsNotSyncOutCalendar := consts.IsNotSyncOutCalendar
	if isSyncOutCalendar != nil && *isSyncOutCalendar == consts.IsSyncOutCalendar {
		isSyncOutCalendar = &IsSyncOutCalendar
	} else {
		isSyncOutCalendar = &IsNotSyncOutCalendar
	}
	insertProjectDetailErr := dao.InsertProjectDetail(po.PpmProProjectDetail{
		Id:                detailId,
		OrgId:             orgId,
		ProjectId:         entity.Id,
		Notice:            consts.BlankString,
		IsSyncOutCalendar: *isSyncOutCalendar,
		Creator:           currentUserId,
		CreateTime:        time.Now(),
		Updator:           currentUserId,
		UpdateTime:        time.Now(),
		IsDelete:          consts.AppIsNoDelete,
	})
	if insertProjectDetailErr != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertProjectDetailErr)
	}

	return nil
}

func insertProject(tx sqlbuilder.Tx, entity bo.ProjectBo) errs.SystemErrorInfo {

	insert := &po.PpmProProject{}
	err1 := copyer.Copy(entity, insert)
	if err1 != nil {
		return errs.BuildSystemErrorInfo(errs.SystemError, err1)
	}
	_, insertProjectErr := tx.Collection(consts.TableProject).Insert(insert)
	if insertProjectErr != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertProjectErr)
	}
	return nil
}

func insertMemberEntities(tx sqlbuilder.Tx, memberEntities []interface{}) error {

	if len(memberEntities) != 0 {
		insertErr := mysql.TransBatchInsert(tx, &po.PpmProProjectRelation{}, memberEntities)
		if insertErr != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertErr)
		}
	}
	return nil
}

func insertSource(entity *bo.ProjectBo, resourcePath string, resourceId int64, tx sqlbuilder.Tx, orgId, currentUserId int64,
	resourceType int) errs.SystemErrorInfo {
	if resourcePath != "" && resourceId == 0 {
		fileName := util.ParseFileName(resourcePath)
		suffix := util.ParseFileSuffix(fileName)
		respVo := resourcefacade.CreateResource(resourcevo.CreateResourceReqVo{
			CreateResourceBo: bo.CreateResourceBo{
				Path:       resourcePath,
				Name:       fileName,
				Suffix:     suffix,
				OrgId:      orgId,
				OperatorId: currentUserId,
				Type:       resourceType,
			},
		})
		if respVo.Failure() {
			return respVo.Error()
		}
		entity.ResourceId = respVo.ResourceId
	}
	return nil
}

//判断项目名是否重复
func JudgeRepeatProjectName(name *string, orgId int64, projectId *int64) (string, errs.SystemErrorInfo) {
	if name == nil {
		*name = consts.BlankString
	}
	cond := make(db.Cond)
	cond = db.Cond{
		consts.TcIsDelete: db.Eq(consts.AppIsNoDelete),
		consts.TcName:     db.Eq(name),
		consts.TcOrgId:    orgId,
	}
	//如果传项目id
	if projectId != nil {
		cond[consts.TcId] = db.NotEq(projectId)
	}
	exist, err := mysql.IsExistByCond(consts.TableProject, cond)
	if err != nil {
		return *name, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if exist {
		return *name, errs.BuildSystemErrorInfo(errs.RepeatProjectName)
	}

	return *name, nil
}

//判断前缀编号是否重复
func JudgeRepeatProjectPreCode(preCode *string, orgId int64, projectId *int64) (string, errs.SystemErrorInfo) {
	if preCode == nil {
		*preCode = consts.BlankString
	}
	cond := make(db.Cond)
	cond = db.Cond{
		consts.TcIsDelete: db.Eq(consts.AppIsNoDelete),
		consts.TcPreCode:  preCode,
		consts.TcOrgId:    orgId,
	}
	//如果传项目id
	if projectId != nil {
		cond[consts.TcId] = db.NotEq(projectId)
	}
	exist, err := mysql.IsExistByCond(consts.TableProject, cond)
	if err != nil {
		return *preCode, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if exist {
		return *preCode, errs.BuildSystemErrorInfo(errs.RepeatProjectPrecode)
	}

	return *preCode, nil
}

//获取项目类型和初始状态
func GetTypeAndStatus(orgId int64, projectTypeId int64, status int64) (int64, int64, errs.SystemErrorInfo) {
	var resultType, resultStatus int64
	projectType := &po.PpmPrsProjectType{}
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    db.In([]int64{orgId, 0}),
	}
	//如果不传项目类型默认是普通项目
	if projectTypeId != 0 {
		cond[consts.TcId] = projectTypeId
	} else {
		cond[consts.TcLangCode] = po.ProjectTypeNormalTask.LangCode
	}
	err := mysql.SelectOneByCond(consts.TableProjectType, cond, projectType)
	if err == nil {
		if projectTypeId == 0 {
			resultType = projectType.Id
		} else {
			resultType = projectTypeId
		}
	} else {
		return resultType, resultStatus, errs.BuildSystemErrorInfo(errs.MysqlOperateError, errors.New("项目类型未初始化"))
	}

	//获取项目初始状态
	if status == 0 {
		defaultStatusId, err := processfacade.GetDefaultProcessStatusIdRelaxed(orgId, projectType.DefaultProcessId, consts.ProcessStatusCategoryProject)
		if err == nil {
			resultStatus = defaultStatusId
		} else {
			return resultType, resultStatus, errs.BuildSystemErrorInfo(errs.MysqlOperateError, errors.New("项目流程状态未初始化"))
		}
	}

	return resultType, resultStatus, nil
}

//创建项目对象流程关联
func CreateProjectObjectTypeProcess(tx sqlbuilder.Tx, projectTypeId int64, orgId int64, projectId int64, currentUserId int64) errs.SystemErrorInfo {
	projectTypeProjectObjectTypeList, err := GetProjectTypeProjectObjectTypeList(orgId)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	//查通用的项目对象类型
	projectObjectTypeList, err := GetProjectObjectTypeList(0, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	projectObjectTypeMap := maps.NewMap("Id", projectObjectTypeList)

	projectObjectTypePos := make([]po.PpmPrsProjectObjectType, 0)
	projectObjectTypeProcessInsert := []interface{}{}
	for _, v := range *projectTypeProjectObjectTypeList {
		if projectTypeId != 0 && projectTypeId == v.ProjectTypeId {
			projectObjectTypeInterface, ok := projectObjectTypeMap[v.ProjectObjectTypeId]
			if !ok {
				log.Errorf("项目对象类型集合%s中不存在%d", json.ToJsonIgnoreError(projectObjectTypeMap), v.ProjectObjectTypeId)
				continue
			}
			projectObjectType := projectObjectTypeInterface.(bo.ProjectObjectTypeBo)
			projectObjectTypePo := &po.PpmPrsProjectObjectType{}
			_ = copyer.Copy(projectObjectType, projectObjectTypePo)

			id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectObjectType)
			if err != nil {
				log.Error(err)
				return errs.ApplyIdError
			}
			projectObjectTypePo.Id = id
			projectObjectTypePo.OrgId = orgId
			projectObjectTypePos = append(projectObjectTypePos, *projectObjectTypePo)

			processId, err := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())
			if err != nil {
				return errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
			}
			projectObjectTypeProcessInsert = append(projectObjectTypeProcessInsert, po.PpmPrsProjectObjectTypeProcess{
				Id:                  processId,
				OrgId:               orgId,
				ProjectId:           projectId,
				ProjectObjectTypeId: projectObjectTypePo.Id,
				ProcessId:           v.DefaultProcessId,
				Creator:             currentUserId,
				CreateTime:          time.Now(),
				Updator:             currentUserId,
				UpdateTime:          time.Now(),
				Version:             1,
			})
		}
	}
	if len(projectObjectTypeProcessInsert) != 0 {
		err := mysql.TransBatchInsert(tx, &po.PpmPrsProjectObjectTypeProcess{}, projectObjectTypeProcessInsert)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}
	if len(projectObjectTypePos) != 0 {
		err := mysql.TransBatchInsert(tx, &po.PpmPrsProjectObjectType{}, slice.ToSlice(projectObjectTypePos))
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}
	return nil
}

func GetProjectList(currentUserId int64, joinParams db.Cond, order []*string, size int, page int) ([]*bo.ProjectBo, int64, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	entities := &[]*po.PpmProProject{}
	mid := conn.Collection(consts.TableProject).Find(joinParams)
	if order != nil {
		for _, v := range order {
			mid = mid.OrderBy(*v)
		}
	}
	mid = mid.OrderBy("id asc")

	count, err := mid.TotalEntries()
	if err != nil {
		return nil, int64(count), errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if size > 0 && page > 0 {
		err = mid.Paginate(uint(size)).Page(uint(page)).All(entities)
	} else {
		err = mid.All(entities)
	}
	if err != nil {
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	projectBos := &[]*bo.ProjectBo{}
	_ = copyer.Copy(entities, projectBos)

	return *projectBos, int64(count), nil
}

func GetProjectInfo(id int64, orgId int64) (bo.ProjectBo, errs.SystemErrorInfo) {
	project := &po.PpmProProject{}
	err := mysql.SelectOneByCond(project.TableName(), db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
	}, project)
	projectBo := &bo.ProjectBo{}
	if err != nil {
		return *projectBo, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	_ = copyer.Copy(project, projectBo)
	return *projectBo, nil
}

func QuitProject(currentUserId, orgId, owner, projectId, memberId int64) errs.SystemErrorInfo {
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := DeleteRelationByDeleteMember(tx, []interface{}{currentUserId}, owner, projectId, orgId, currentUserId)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}

		deleteErr := mysql.TransUpdateSmart(tx, consts.TableProjectRelation, memberId, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if deleteErr != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, deleteErr)
		}

		return nil
	})
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

func UpdateProject(orgId int64, currentUserId int64, updPoint *mysql.Upd, oldOwner int64, input bo.UpdateProjectBo) ([]int64, []int64, errs.SystemErrorInfo) {
	var resourceId int64
	if input.ResourcePath != nil && input.ResourceType != nil {
		respVo := resourcefacade.GetIdByPath(
			resourcevo.GetIdByPathReqVo{
				OrgId:        orgId,
				ResourceType: *input.ResourceType,
				ResourcePath: *input.ResourcePath,
			})
		if !respVo.Failure() {
			resourceId = respVo.ResourceId
		}
	}
	upd := *updPoint
	oldMembers := set.New(set.ThreadSafe)
	thisMembers := set.New(set.ThreadSafe)

	err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//插入资源
		err := updateSource(input, resourceId, tx, orgId, currentUserId, &upd)
		if err != nil {
			log.Error(err)
			return err
		}

		//查看团队原有成员
		oldMembers, thisMembers, err = GetChangeMembersAndDeleteOld(tx, input, orgId, oldOwner, &upd)

		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		err = toViewMembers(input, tx, orgId, currentUserId, oldOwner, oldMembers, thisMembers)

		if err != nil {
			log.Error(err)
			return err
		}

		//更新项目
		projectError := updateProject(upd, tx, currentUserId, input)

		if projectError != nil {
			log.Error(projectError)
			return projectError
		}

		return nil
	})

	if err1 != nil {
		log.Error(err1)
		return nil, nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
	}

	beforeMemberIds := make([]int64, oldMembers.Size())
	afterMemberIds := make([]int64, thisMembers.Size())

	for i, member := range oldMembers.List() {
		beforeMemberIds[i] = member.(int64)
	}
	for i, member := range thisMembers.List() {
		afterMemberIds[i] = member.(int64)
	}

	return beforeMemberIds, afterMemberIds, nil
}

func UpdateFollower(input bo.UpdateProjectBo, currentUserId, orgId int64) ([]int64, []int64, errs.SystemErrorInfo) {
	if !util.FieldInUpdate(input.UpdateFields, "followerIds") {
		return nil, nil, nil
	}
	followerIds := slice.SliceUniqueInt64(input.FollowerIds)
	memberEntities := &[]po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    input.ID,
		consts.TcRelationType: consts.IssueRelationTypeFollower,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, memberEntities)
	if err != nil {
		return nil, nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
	}
	oldFollower := set.New(set.ThreadSafe)
	oldSlice := []int64{}
	newSlice := []int64{}
	newFollower := set.New(set.ThreadSafe)
	for _, v := range *memberEntities {
		oldFollower.Add(v.RelationId)
		oldSlice = append(oldSlice, v.RelationId)
	}
	for _, v := range followerIds {
		newFollower.Add(v)
		newSlice = append(newSlice, v)
	}

	delFollower := set.Difference(oldFollower, newFollower)
	addFollower := set.Difference(newFollower, oldFollower)
	_, delErr := mysql.UpdateSmartWithCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    input.ID,
		consts.TcRelationType: consts.IssueRelationTypeFollower,
		consts.TcRelationId:   db.In(delFollower.List()),
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	if delErr != nil {
		return nil, nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, delErr)
	}

	if addFollower.Size() == 0 {
		return oldSlice, newSlice, nil
	}

	addFollowerIds := make([]int64, 0)
	for _, followerIdInterface := range addFollower.List() {
		if followerId, ok := followerIdInterface.(int64); ok {
			addFollowerIds = append(addFollowerIds, followerId)
		}
	}

	//校验组织
	orgUserVerifyFlag := orgfacade.VerifyOrgUsersRelaxed(orgId, addFollowerIds)
	if !orgUserVerifyFlag {
		log.Error("存在用户组织校验失败")
		return nil, nil, errs.VerifyOrgError
	}

	updateProjectRelationErr := UpdateProjectRelation(currentUserId, orgId, input.ID, consts.IssueRelationTypeFollower, addFollowerIds)
	if updateProjectRelationErr != nil {
		log.Error(updateProjectRelationErr)
		return nil, nil, updateProjectRelationErr
	}

	return oldSlice, newSlice, nil
}

func updateProject(upd mysql.Upd, tx sqlbuilder.Tx, currentUserId int64, input bo.UpdateProjectBo) errs.SystemErrorInfo {
	if len(upd) > 0 {
		//更新项目
		upd[consts.TcUpdator] = currentUserId
		upd[consts.TcUpdateTime] = time.Now()
		_, updateProjectErr := mysql.TransUpdateSmartWithCond(tx, consts.TableProject, db.Cond{
			consts.TcId: input.ID,
		}, upd)
		if updateProjectErr != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, updateProjectErr)
		}
	}
	return nil
}

func toViewMembers(input bo.UpdateProjectBo, tx sqlbuilder.Tx, orgId, currentUserId int64, oldOwner int64, oldMembers, thisMembers set.Interface) errs.SystemErrorInfo {

	if util.FieldInUpdate(input.UpdateFields, "memberIds") || util.FieldInUpdate(input.UpdateFields, "owner") {
		err := ChangeMembers(tx, input, orgId, currentUserId, oldOwner, oldMembers, thisMembers)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}
	return nil
}

func updateSource(input bo.UpdateProjectBo, resourceId int64, tx sqlbuilder.Tx, orgId, currentUserId int64, upd *mysql.Upd) errs.SystemErrorInfo {

	if util.FieldInUpdate(input.UpdateFields, "resourcePath") && util.FieldInUpdate(input.UpdateFields, "resourceType") {
		if input.ResourcePath != nil && input.ResourceType != nil {
			if resourceId != 0 {
				(*upd)[consts.TcResourceId] = resourceId
			} else {
				fileName := util.ParseFileName(*input.ResourcePath)
				suffix := util.ParseFileSuffix(fileName)
				respVo := resourcefacade.CreateResource(resourcevo.CreateResourceReqVo{
					CreateResourceBo: bo.CreateResourceBo{
						Path:       *input.ResourcePath,
						Name:       fileName,
						Suffix:     suffix,
						OrgId:      orgId,
						OperatorId: currentUserId,
						Type:       *input.ResourceType,
					},
				})
				if respVo.Failure() {
					return respVo.Error()
				}
				(*upd)[consts.TcResourceId] = respVo.ResourceId
			}
		} else {
			(*upd)[consts.TcResourceId] = 0
		}
	}

	return nil
}

func UpdateProjectCondAssembly(input bo.UpdateProjectBo, orgId int64, old, new *map[string]interface{}, originProjectInfo bo.ProjectBo, changeList *[]bo.TrendChangeListBo) (mysql.Upd, errs.SystemErrorInfo) {
	planStartTime := time.Time(originProjectInfo.PlanStartTime)
	planEndTime := time.Time(originProjectInfo.PlanEndTime)
	upd := mysql.Upd{}

	repeatErr := needUpdateVertifyRepeat(input, &upd, orgId, old, new, originProjectInfo, changeList)
	if repeatErr != nil {
		return nil, repeatErr
	}
	priorityIdErr := needUpdatePriorityId(input, &upd, orgId, old, new, originProjectInfo, changeList)
	if priorityIdErr != nil {
		return nil, priorityIdErr
	}
	needUpdateVertifyValidField(input, &upd, old, new, originProjectInfo, changeList)
	simpleErr := needUpdateSimpleField(input, &upd, old, new, originProjectInfo, changeList)
	if simpleErr != nil {
		return nil, simpleErr
	}
	planTimeErr := needUpdatePlanTime(input, &planStartTime, &planEndTime, &upd, old, new, originProjectInfo, changeList)
	if planTimeErr != nil {
		return nil, planTimeErr
	}

	return upd, nil
}

func needUpdateVertifyValidField(input bo.UpdateProjectBo, upd *mysql.Upd, old, new *map[string]interface{}, originProjectInfo bo.ProjectBo, changeList *[]bo.TrendChangeListBo) {
	publicStatus := map[int]string{
		consts.PrivateProject: "私有",
		consts.PublicProject:  "公开",
	}
	if util.FieldInUpdate(input.UpdateFields, "publicStatus") {
		if input.PublicStatus != nil {
			if ok, _ := slice.Contain([]int{consts.PrivateProject, consts.PublicProject}, *input.PublicStatus); ok {
				(*upd)[consts.TcPublicStatus] = input.PublicStatus
				(*old)["publicStatus"] = originProjectInfo.PublicStatus
				(*new)["publicStatus"] = input.PublicStatus
				*changeList = append(*changeList, bo.TrendChangeListBo{
					Field:     "publicStatus",
					FieldName: consts.PublicStatus,
					OldValue:  publicStatus[originProjectInfo.PublicStatus],
					NewValue:  publicStatus[*input.PublicStatus],
				})
			}
		}
	}

	if util.FieldInUpdate(input.UpdateFields, "isFiling") {
		//todo 暂时归档项目不放在更新项目的动态
		if input.IsFiling != nil {
			if ok, _ := slice.Contain([]int{consts.ProjectIsFiling, consts.ProjectIsNotFiling}, *input.IsFiling); ok {
				(*upd)[consts.TcIsFiling] = input.IsFiling
				(*old)["isFiling"] = originProjectInfo.IsFiling
				(*new)["isFiling"] = input.IsFiling
			}
		}
	}
}

func needUpdatePriorityId(input bo.UpdateProjectBo, upd *mysql.Upd, orgId int64, old, new *map[string]interface{}, originProjectInfo bo.ProjectBo, changeList *[]bo.TrendChangeListBo) errs.SystemErrorInfo {
	if util.FieldInUpdate(input.UpdateFields, "priorityId") {
		if input.PriorityID != nil {
			suc, err := VerifyPriority(orgId, consts.PriorityTypeProject, *input.PriorityID)
			if err != nil {
				return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
			}
			if !suc {
				return errs.BuildSystemErrorInfo(errs.IllegalityPriority, err)
			}
			(*upd)[consts.TcPriorityId] = input.PriorityID
			(*old)["priorityId"] = originProjectInfo.PriorityId
			(*new)["priorityId"] = input.PriorityID
			old, err := GetPriorityBo(originProjectInfo.PriorityId)
			if err != nil {
				return err
			}
			new, err := GetPriorityBo(*input.PriorityID)
			if err != nil {
				return err
			}
			*changeList = append(*changeList, bo.TrendChangeListBo{
				Field:     "priorityId",
				FieldName: consts.Priority,
				OldValue:  old.Name,
				NewValue:  new.Name,
			})
		}
	}

	return nil
}

func needUpdateVertifyRepeat(input bo.UpdateProjectBo, upd *mysql.Upd, orgId int64, old, new *map[string]interface{}, originProjectInfo bo.ProjectBo, changeList *[]bo.TrendChangeListBo) errs.SystemErrorInfo {
	//判断项目名是否重复
	if util.FieldInUpdate(input.UpdateFields, "name") {
		if input.Name == nil {
			return nil
		} else if strings.Trim(*input.Name, " ") == consts.BlankString {
			return errs.ProjectNameEmpty
		}
		isNameRight := format.VerifyProjectNameFormat(*input.Name)
		if !isNameRight {
			log.Error(errs.InvalidProjectNameError)
			return errs.InvalidProjectNameError
		}
		name, err := JudgeRepeatProjectName(input.Name, orgId, &input.ID)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
		} else {
			(*upd)[consts.TcName] = name
			(*old)["name"] = originProjectInfo.Name
			(*new)["name"] = input.Name
			*changeList = append(*changeList, bo.TrendChangeListBo{
				Field:     "title",
				FieldName: consts.Title,
				OldValue:  originProjectInfo.Name,
				NewValue:  *input.Name,
			})
		}
	}

	return nil
}

func needUpdateSimpleField(input bo.UpdateProjectBo, upd *mysql.Upd, old, new *map[string]interface{}, originProjectInfo bo.ProjectBo, changeList *[]bo.TrendChangeListBo) errs.SystemErrorInfo {
	if util.FieldInUpdate(input.UpdateFields, "remark") {
		if input.Remark != nil {
			//if strs.Len(*input.Remark) > 500 {
			//	return errs.TooLongProjectRemark
			//}
			isRemarkRight := format.VerifyProjectRemarkFormat(*input.Remark)
			if !isRemarkRight {
				log.Error(errs.InvalidProjectRemarkError)
				return errs.InvalidProjectRemarkError
			}
			(*upd)[consts.TcRemark] = input.Remark
		} else {
			(*upd)[consts.TcRemark] = consts.BlankString
		}
		(*old)["remark"] = originProjectInfo.Remark
		(*new)["remark"] = input.Remark
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "remark",
			FieldName: consts.Remark,
			OldValue:  originProjectInfo.Remark,
			NewValue:  *input.Remark,
		})
	}
	return nil
}

func needUpdatePlanTime(input bo.UpdateProjectBo, planStartTime, planEndTime *time.Time, upd *mysql.Upd, old, new *map[string]interface{}, originProjectInfo bo.ProjectBo, changeList *[]bo.TrendChangeListBo) errs.SystemErrorInfo {
	if util.FieldInUpdate(input.UpdateFields, "planStartTime") {
		if input.PlanStartTime != nil && input.PlanStartTime.IsNotNull() {
			(*upd)[consts.TcPlanStartTime] = date.FormatTime(*input.PlanStartTime)
			*planStartTime = time.Time(*input.PlanStartTime)
		} else {
			(*upd)[consts.TcPlanStartTime] = consts.BlankTime
			*planStartTime = consts.BlankTimeObject
		}
		(*old)["planStartTime"] = originProjectInfo.PlanStartTime
		(*new)["planStartTime"] = (*upd)[consts.TcPlanStartTime]
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "planStartTime",
			FieldName: consts.PlanStartTime,
			OldValue:  originProjectInfo.PlanStartTime.String(),
			NewValue:  (*upd)[consts.TcPlanStartTime].(string),
		})
	}

	if util.FieldInUpdate(input.UpdateFields, "planEndTime") {
		if input.PlanEndTime != nil && input.PlanEndTime.IsNotNull() {
			(*upd)[consts.TcPlanEndTime] = date.FormatTime(*input.PlanEndTime)
			*planEndTime = time.Time(*input.PlanEndTime)
		} else {
			(*upd)[consts.TcPlanEndTime] = consts.BlankTime
			*planEndTime = consts.BlankTimeObject
		}
		(*old)["planEndTime"] = originProjectInfo.PlanEndTime
		(*new)["planEndTime"] = (*upd)[consts.TcPlanEndTime]
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "planEndTime",
			FieldName: consts.PlanEndTime,
			OldValue:  originProjectInfo.PlanEndTime.String(),
			NewValue:  (*upd)[consts.TcPlanEndTime].(string),
		})
	}

	if (*planEndTime).After(consts.BlankTimeObject) && planStartTime.After(*planEndTime) {
		return errs.BuildSystemErrorInfo(errs.CreateProjectTimeError)
	}

	return nil
}

//
//func GetProjectCondAssembly(params map[string]interface{}, currentUserId int64, orgId int64) (db.Cond, errs.SystemErrorInfo) {
//	var relationType interface{}
//	if _, ok := params["relation_type"]; ok {
//		if val, ok := params["relation_type"].(map[string]interface{}); ok {
//			if val["type"] != nil && val["value"] != nil {
//				relationType = val["value"]
//			}
//		} else {
//			relationType = params["relation_type"]
//		}
//		delete(params, "relation_type")
//	}
//	var relateType int64 = 0
//	if val, ok := relationType.(json.Number); ok {
//		relateType, _ = val.Int64()
//	} else if val, ok := relationType.(int64); ok {
//		relateType = val
//	}
//
//	condParam, err := cond.HandleParams(params)
//	if err != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.ConditionHandleError, err)
//	}
//	switch relateType {
//	case 0:
//		//所有
//	case 1:
//		//我发起的
//		condParam[consts.TcCreator] = currentUserId
//	case 2:
//		//我负责的
//		condParam[consts.TcOwner] = db.Eq(currentUserId)
//	case 3:
//		//我参与的
//		need, err := GetParticipantMembers(orgId, currentUserId)
//		if err != nil {
//			return nil, errs.BuildSystemErrorInfo(errs.ConditionHandleError, err)
//		}
//		condParam[consts.TcId] = db.In(need)
//	}
//
//	//默认查询没有被删除的
//	condParam[consts.TcIsDelete] = consts.AppIsNoDelete
//	joinParams := make(db.Cond)
//	for k, v := range condParam {
//		joinParams["p."+k.(string)] = v
//	}
//
//	return joinParams, nil
//}

func JudgeProjectIsExist(orgId, id int64) bool {
	_, err := LoadProjectAuthBo(orgId, id)
	if err != nil {
		return false
	}

	return true
}

func StatProject(orgId, id int64) (bo.ProjectStatBo, errs.SystemErrorInfo) {
	projectStat := bo.ProjectStatBo{}
	_, err := dao.SelectOneProject(db.Cond{
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err != nil {
		return projectStat, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
	}
	projectStat.MemberTotal = dao.GetProjectMemberCount(orgId, id)

	iterationTotal, err := mysql.SelectCountByCond(consts.TableIteration, db.Cond{
		consts.TcProjectId: id,
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	})
	if err != nil {
		return projectStat, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	projectStat.IterationTotal = int64(iterationTotal)

	taskTotal, err := mysql.SelectCountByCond(consts.TableIssue, db.Cond{
		consts.TcProjectId: id,
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	})
	if err != nil {
		return projectStat, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	projectStat.TaskTotal = int64(taskTotal)

	return projectStat, nil
}

func UpdateProjectStatus(projectBo bo.ProjectBo, nextStatusId int64) errs.SystemErrorInfo {
	orgId := projectBo.OrgId
	projectId := projectBo.Id

	if projectBo.Status == nextStatusId {
		log.Error("更新项目状态-要更新的状态和当前状态一样")
		return errs.BuildSystemErrorInfo(errs.ProjectStatusUpdateError)
	}

	//验证状态有效性
	_, err := GetProjectRelationBo(projectBo, consts.IssueRelationTypeStatus, nextStatusId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	_, err2 := dao.UpdateProjectByOrg(projectId, orgId, mysql.Upd{
		consts.TcStatus: nextStatusId,
	})
	if err2 != nil {
		log.Error(err2)
		return errs.BuildSystemErrorInfo(errs.IterationStatusUpdateError)
	}

	return nil
}
