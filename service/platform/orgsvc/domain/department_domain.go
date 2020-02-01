package domain

import (
	"fmt"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetDepartmentBoList(page uint, size uint, cond db.Cond) (*[]bo.DepartmentBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectDepartmentByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.DepartmentBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetDepartmentBoWithOrg(id int64, orgId int64) (*bo.DepartmentBo, errs.SystemErrorInfo) {
	po, err := dao.SelectDepartmentByIdAndOrg(id, orgId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.DepartmentBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func JudgeDepartmentIsExist(orgId int64, name string) (bool, errs.SystemErrorInfo) {
	exist, err := mysql.IsExistByCond(consts.TableDepartment, db.Cond{
		consts.TcName:     name,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return exist, nil
}

func GetDepartmentMembers(orgId int64, departmentId *int64) ([]bo.DepartmentMemberInfoBo, errs.SystemErrorInfo) {
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if departmentId != nil {
		cond[consts.TcDepartmentId] = *departmentId
	}
	userDepartmentList, err := dao.SelectUserDepartment(cond)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	userIds := make([]int64, 0)
	for _, userDepartmentInfo := range *userDepartmentList {
		userId := userDepartmentInfo.UserId
		userIdExist, _ := slice.Contain(userIds, userId)
		if !userIdExist {
			userIds = append(userIds, userId)
		}
	}
	userInfos, err1 := GetBaseUserInfoBatch(consts.AppSourceChannelDingTalk, orgId, userIds)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}
	userMap := maps.NewMap("UserId", userInfos)

	userIdInfoBoList := make([]bo.DepartmentMemberInfoBo, 0)
	for _, userDepartment := range *userDepartmentList {
		if userCacheInfo, ok := userMap[userDepartment.UserId]; ok {
			baseUserInfo := userCacheInfo.(bo.BaseUserInfoBo)
			userIdInfoBoList = append(userIdInfoBoList, bo.DepartmentMemberInfoBo{
				UserID:       baseUserInfo.UserId,
				Name:         baseUserInfo.Name,
				NamePy:		  baseUserInfo.NamePy,
				Avatar:       baseUserInfo.Avatar,
				EmplID:       baseUserInfo.OutUserId,
				DepartmentID: userDepartment.DepartmentId,
			})
		} else {
			log.Errorf("GetDepartmentMembers: 查询不到部门%d下的用户%d信息", userDepartment.DepartmentId, userDepartment.UserId)
		}
	}
	return userIdInfoBoList, nil
}

func GetTopDepartmentInfoList(orgId int64) ([]bo.DepartmentBo, errs.SystemErrorInfo) {
	departmentInfo := &[]po.PpmOrgDepartment{}
	err := mysql.SelectAllByCond(consts.TableDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcStatus:   consts.AppIsInitStatus,
		consts.TcParentId: 0,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, departmentInfo)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	departmentInfoBo := &[]bo.DepartmentBo{}
	err1 := copyer.Copy(departmentInfo, departmentInfoBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *departmentInfoBo, nil
}

func GetTopDepartmentInfo(orgId int64) (*bo.DepartmentBo, errs.SystemErrorInfo) {
	departmentInfoList, err := GetTopDepartmentInfoList(orgId)
	if err != nil {
		log.Errorf("获取部门信息错误 %v", err)
		return nil, err
	}
	if len(departmentInfoList) == 0 {
		log.Errorf("组织%d下不存在顶级部门", orgId)
		return nil, errs.BuildSystemErrorInfo(errs.TopDepartmentNotExist)
	}
	departmentInfo := departmentInfoList[0]
	return &departmentInfo, nil
}

func BoundOrgMemberToTopDepartment(orgId int64, userIds []int64, operatorId int64) (int64, errs.SystemErrorInfo) {
	departmentInfo, err := GetTopDepartmentInfoList(orgId)
	if err != nil {
		log.Error("获取部门信息错误 " + strs.ObjectToString(err))
		return 0, err
	}
	var departmentId int64
	for _, v := range departmentInfo {
		departmentId = v.Id
		break
	}
	return departmentId, BoundDepartmentUser(orgId, userIds, departmentId, operatorId, false)
}

//解绑部门用户，解绑当前用户所在的所有部门
func UnBoundDepartmentUser(orgId int64, userIds []int64, operatorId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	//查询已有的绑定关系
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   db.In(userIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	_, err := dao.UpdateUserDepartmentByCond(cond, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operatorId,
	}, tx)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

//绑定部门用户，带分布式锁
func BoundDepartmentUser(orgId int64, userIds []int64, departmentId, operatorId int64, isLeaderFlag bool) errs.SystemErrorInfo {
	isLeader := 2
	if isLeaderFlag {
		isLeader = 1
	}

	//先上锁
	lockKey := consts.UserAndDepartmentRelationLockKey + fmt.Sprintf("%d:%d", orgId, departmentId)
	uuid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uuid)
	if lockErr != nil {
		log.Error(lockErr)
		return errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
	}
	if !suc {
		log.Errorf("绑定用户时没有获取到锁 orgId %d departmentId %d", orgId, departmentId)
		return errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
	}
	defer func() {
		if _, err := cache.ReleaseDistributedLock(lockKey, uuid); err != nil {
			log.Error(err)
		}
	}()
	//查询已有的绑定关系
	cond := db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcDepartmentId: departmentId,
		consts.TcUserId:       db.In(userIds),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}
	userDepartmentList, dbErr := dao.SelectUserDepartment(cond)
	if dbErr != nil {
		log.Error(dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	notRelationUserIds := make([]int64, 0)
	alreadyRelationUserIdMap := map[int64]bool{}
	if userDepartmentList != nil && len(*userDepartmentList) > 0 {
		for _, userDepartment := range *userDepartmentList {
			alreadyRelationUserIdMap[userDepartment.UserId] = true
		}
	}
	//获取没有关联关系的用户
	for _, userId := range userIds {
		if _, ok := alreadyRelationUserIdMap[userId]; !ok {
			notRelationUserIds = append(notRelationUserIds, userId)
			alreadyRelationUserIdMap[userId] = true
		}
	}
	if len(notRelationUserIds) == 0 {
		return nil
	}

	userIdsLen := len(notRelationUserIds)

	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, userIdsLen)
	if err != nil {
		log.Error(err)
		return err
	}

	userDepartments := make([]po.PpmOrgUserDepartment, userIdsLen)
	for i, userId := range notRelationUserIds {
		userDepartments[i] = po.PpmOrgUserDepartment{
			Id:           ids.Ids[i].Id,
			OrgId:        orgId,
			UserId:       userId,
			DepartmentId: departmentId,
			IsLeader:     isLeader,
			Creator:      operatorId,
			Updator:      operatorId,
		}
	}

	err1 := dao.InsertUserDepartmentBatch(userDepartments)
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}
