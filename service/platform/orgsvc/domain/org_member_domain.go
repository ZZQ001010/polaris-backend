package domain

import (
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//修改组织成员状态，企业用户状态, 1可用,2禁用
func ModifyOrgMemberStatus(orgId int64, memberIds []int64, status int, operatorId int64) errs.SystemErrorInfo {
	//组织负责人不允许被修改状态
	orgInfo, err := GetBaseOrgInfo("", orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	orgOwnerId := orgInfo.OrgOwnerId
	memberIds = filterMemberIds(memberIds, operatorId)
	if orgOwnerId != operatorId {
		memberIds = filterMemberIds(memberIds, orgOwnerId)
	}
	if len(memberIds) == 0 {
		return errs.UpdateMemberIdsIsEmptyError
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		modifyCount, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcUserId:      db.In(memberIds),
			consts.TcIsDelete:    consts.AppIsNoDelete,
			consts.TcStatus:      db.NotEq(status),
			consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		}, mysql.Upd{
			consts.TcStatus:           status,
			consts.TcStatusChangerId:  operatorId,
			consts.TcUpdator:          operatorId,
			consts.TcStatusChangeTime: time.Now(),
		})
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
		//如果更新数量与预期不符，认为动作失败
		if modifyCount != int64(len(memberIds)) {
			return errs.BuildSystemErrorInfo(errs.UpdateMemberStatusFail)
		}
		//禁用的用户是否可以在选人界面显示，待定
		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}
	//最后将用户信息缓存清掉
	clearErr := ClearBaseUserInfoBatch(orgId, memberIds)
	if clearErr != nil {
		log.Error(clearErr)
	}
	return nil
}

//修改组织成员审核状态, 审核状态,1待审核,2审核通过,3审核不过
func ModifyOrgMemberCheckStatus(orgId int64, memberIds []int64, checkStatus int, operatorId int64) errs.SystemErrorInfo {
	//组织负责人不允许被修改状态
	orgInfo, err := GetBaseOrgInfo("", orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	orgOwnerId := orgInfo.OrgOwnerId
	memberIds = filterMemberIds(memberIds, operatorId)
	if orgOwnerId != operatorId {
		memberIds = filterMemberIds(memberIds, orgOwnerId)
	}
	if len(memberIds) == 0 {
		return errs.UpdateMemberIdsIsEmptyError
	}

	isCheckPass := checkStatus == consts.AppCheckStatusSuccess
	departmentId := int64(0)

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		upd := mysql.Upd{
			consts.TcCheckStatus: checkStatus,
			consts.TcAuditorId:   operatorId,
			consts.TcUpdator:     operatorId,
			consts.TcAuditTime:   time.Now(),
		}

		if isCheckPass {
			upd[consts.TcStatus] = consts.AppStatusEnable
		}
		modifyCount, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcUserId:      db.In(memberIds),
			consts.TcCheckStatus: consts.AppCheckStatusWait,
			consts.TcIsDelete:    consts.AppIsNoDelete,
		}, upd)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
		//如果更新数量与预期不符，认为动作失败
		if modifyCount != int64(len(memberIds)) {
			return errs.BuildSystemErrorInfo(errs.UpdateMemberStatusFail)
		}
		//审核通过，加入部门
		if isCheckPass {
			depId, err := BoundOrgMemberToTopDepartment(orgId, memberIds, operatorId)
			if err != nil {
				log.Error(err)
				return err
			}
			departmentId = depId
		}
		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}
	//推送消息
	if isCheckPass{
		asyn.Execute(func(){
			PushAddOrgMemberNotice(orgId, departmentId, memberIds)
		})
	}

	//最后将用户信息缓存清掉
	clearErr := ClearBaseUserInfoBatch(orgId, memberIds)
	if clearErr != nil {
		log.Error(clearErr)
	}
	return nil
}

//修改组织成员
func RemoveOrgMember(orgId int64, memberIds []int64, operatorId int64) errs.SystemErrorInfo {
	//组织负责人不允许被修改状态
	orgInfo, err := GetBaseOrgInfo("", orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	orgOwnerId := orgInfo.OrgOwnerId
	memberIds = filterMemberIds(memberIds, operatorId)
	if orgOwnerId != operatorId {
		memberIds = filterMemberIds(memberIds, orgOwnerId)
	}
	if len(memberIds) == 0 {
		return errs.UpdateMemberIdsIsEmptyError
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		modifyCount, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcUserId:      db.In(memberIds),
			consts.TcCheckStatus: consts.AppCheckStatusSuccess,
			consts.TcIsDelete:    consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcAuditorId: operatorId,
			consts.TcUpdator:   operatorId,
			consts.TcIsDelete:  consts.AppIsDeleted,
		})
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
		//如果更新数量与预期不符，认为动作失败
		if modifyCount != int64(len(memberIds)) {
			return errs.BuildSystemErrorInfo(errs.UpdateMemberStatusFail)
		}
		//将用户从组织移除之后 - 将该用户从部门移除
		err = UnBoundDepartmentUser(orgId, memberIds, operatorId, tx)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}

	asyn.Execute(func(){
		PushRemoveOrgMemberNotice(orgId, memberIds)
	})

	//清掉用户的角色，弱业务（因为用户已经禁掉），不需要事务
	roleErr := rolefacade.RemoveRoleUserRelation(rolevo.RemoveRoleUserRelationReqVo{
		OrgId:      orgId,
		UserIds:    memberIds,
		OperatorId: operatorId,
	})
	if roleErr.Failure() {
		log.Error(roleErr.Message)
	}
	//最后将用户信息缓存清掉
	clearErr := ClearBaseUserInfoBatch(orgId, memberIds)
	if clearErr != nil {
		log.Error(clearErr)
	}
	return nil
}

//过滤掉当前操作人
func filterMemberIds(memberIds []int64, operatorId int64) []int64 {
	memberIds = slice.SliceUniqueInt64(memberIds)
	newMemberIds := make([]int64, 0)
	for _, memberId := range memberIds {
		if memberId != operatorId {
			newMemberIds = append(newMemberIds, memberId)
		}
	}
	return newMemberIds
}

//通过渠道获取组织用户信息
func GetOrgUserInfoListBySourceChannel(orgId int64, sourceChannel string, page, size int) ([]bo.OrgUserInfo, int64, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		log.Error(err)
		return nil, 0, errs.MysqlOperateError
	}

	mid := conn.Select(
		"userOrg.user_id as user_id",
		"userOutInfo.out_user_id as out_user_id",
		"userOrg.org_id as org_id",
		"userOrg.status as org_user_status",
		"userOrg.check_status as org_user_check_status",
	).
		From("ppm_org_user_organization userOrg").
		LeftJoin("ppm_org_user_out_info userOutInfo").
		On("userOrg.user_id = userOutInfo.user_id").
		Where(db.Cond{
			"userOrg.org_id":             orgId,
			"userOrg.is_delete":          consts.AppIsNoDelete,
			"userOutInfo.is_delete":      consts.AppIsNoDelete,
			"userOutInfo.source_channel": sourceChannel,
		})

	result := &[]bo.OrgUserInfo{}
	total := int64(0)
	if size > 0 && page > 0 {
		pageResult := mid.Paginate(uint(size)).Page(uint(page))
		rowSize, err := pageResult.TotalEntries()
		if err != nil {
			log.Error(err)
			return nil, 0, errs.MysqlOperateError
		}
		total = int64(rowSize)
		err = pageResult.All(result)
		if err != nil {
			log.Error(err)
			return nil, 0, errs.MysqlOperateError
		}
	} else {
		err := mid.All(result)
		if err != nil {
			log.Error(err)
			return nil, 0, errs.MysqlOperateError
		}
		total = int64(len(*result))
	}
	return *result, total, nil
}
