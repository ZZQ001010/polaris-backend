package domain

import (
	"fmt"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"strings"
	"upper.io/db.v3"
)

//用来返回用户组织列表
func GetUserOrganizationIdList(userId int64) (*[]bo.PpmOrgUserOrganizationBo, errs.SystemErrorInfo) {

	UserOrganizationPo := &[]po.PpmOrgUserOrganization{}
	UserOrganizationBo := &[]bo.PpmOrgUserOrganizationBo{}

	err := mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
		consts.TcUserId:      userId,
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		//consts.TcStatus:   consts.AppStatusEnable,
	}, UserOrganizationPo)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	_ = copyer.Copy(UserOrganizationPo, UserOrganizationBo)

	return UserOrganizationBo, nil
}

//用来获取用户最新的组织关系
func GetUserOrganizationNewestRelation(orgId, userId int64) (*bo.PpmOrgUserOrganizationBo, errs.SystemErrorInfo) {
	UserOrganizationPo := &po.PpmOrgUserOrganization{}
	UserOrganizationBo := &bo.PpmOrgUserOrganizationBo{}

	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return nil, errs.MysqlOperateError
	}
	err = conn.Collection(consts.TableUserOrganization).Find(db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcUserId: userId,
	}).OrderBy("id desc").Limit(1).One(UserOrganizationPo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.BuildSystemErrorInfo(errs.UserOrgNotRelation)
		} else {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
	}
	_ = copyer.Copy(UserOrganizationPo, UserOrganizationBo)
	return UserOrganizationBo, nil
}

func handleGetOrganizationUserListCond(input *vo.OrgUserListReq, cond db.Cond) {
	if len(input.CheckStatus) > 0 {
		cond["o."+consts.TcCheckStatus] = db.In(input.CheckStatus)
	}
	if input.Status != nil && *input.Status != 0 {
		cond["o."+consts.TcStatus] = *input.Status
	}
	if input.UseStatus != nil && *input.UseStatus != 0 {
		cond["o."+consts.TcUseStatus] = *input.UseStatus
	}
}

func GetOrganizationUserList(orgId int64, page, size int, input *vo.OrgUserListReq, allUserHaveRoleIds []int64) (uint64, []bo.PpmOrgUserOrganizationBo, errs.SystemErrorInfo) {

	conn, err := mysql.GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	cond := db.Cond{
		"o." + consts.TcOrgId:    orgId,
		"o." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcIsDelete: consts.AppIsNoDelete,
		"o." + consts.TcUserId:   db.Raw("u." + consts.TcId),
	}
	if input != nil {
		handleGetOrganizationUserListCond(input, cond)
	}
	total := &po.PpmOrgUserOrganizationCount{}
	totalErr := conn.Select(db.Raw("count(*) as total")).From("ppm_org_user_organization o", "ppm_org_user u").Where(cond).One(total)
	if totalErr != nil {
		log.Error(totalErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, totalErr)
	}
	count := total.Total

	//count, err := mysql.SelectCountByCond(consts.TableUserOrganization, cond)
	//if err != nil {
	//	return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}

	orgUserPo := &[]po.PpmOrgUserOrganization{}
	//默认是审核时间升序，创建时间降序
	order := db.Raw("o.check_status asc, o.create_time desc")
	if len(allUserHaveRoleIds) > 0 {
		idStr := strings.Replace(strings.Trim(fmt.Sprint(allUserHaveRoleIds), "[]"), " ", ",", -1)
		order = db.Raw("FIELD(o.user_id," + idStr + ") desc, o.check_status asc, o.create_time desc")
	}

	//err = mysql.SelectAllByCondWithNumAndOrder(consts.TableUserOrganization, cond, nil, page, size, order, orgUserPo)
	//if err != nil {
	//	return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}

	selectErr := conn.Select(db.Raw("o.*")).From("ppm_org_user_organization o", "ppm_org_user u").Where(cond).Offset((page - 1) * size).Limit(size).
		OrderBy(order).All(orgUserPo)
	if selectErr != nil {
		log.Error(selectErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
	}

	orgUserBo := &[]bo.PpmOrgUserOrganizationBo{}
	copyErr := copyer.Copy(orgUserPo, orgUserBo)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return count, *orgUserBo, nil
}

func GetOrgIdListBySourceChannel(sourceChannel string, page int, size int) ([]int64, errs.SystemErrorInfo) {
	orgOutInfos := &[]po.PpmOrgOrganizationOutInfo{}
	_, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcSourceChannel: sourceChannel,
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcStatus:        consts.AppStatusEnable,
	}, nil, page, size, "", orgOutInfos)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	orgIds := make([]int64, 0)

	for _, outInfo := range *orgOutInfos {
		orgIds = append(orgIds, outInfo.OrgId)
	}

	orgIds = slice.SliceUniqueInt64(orgIds)
	return orgIds, nil
}
