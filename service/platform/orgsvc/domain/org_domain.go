package domain

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"strings"
	"upper.io/db.v3"
)

func GetOrgBoList() ([]bo.OrganizationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmOrgOrganization{}
	err := mysql.SelectAllByCond(consts.TableOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
	}, pos)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bos := &[]bo.OrganizationBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return *bos, nil
}

func GetOrgBoListByIds(orgIds []int64) (*[]bo.OrganizationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmOrgOrganization{}
	err := mysql.SelectAllByCond(consts.TableOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcId:       db.In(orgIds),
	}, pos)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bos := &[]bo.OrganizationBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return bos, nil
}

func GetOrgBoById(orgId int64) (*bo.OrganizationBo, errs.SystemErrorInfo) {
	po := &po.PpmOrgOrganization{}
	err := mysql.SelectOneByCond(consts.TableOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcId:       orgId,
	}, po)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bo := &bo.OrganizationBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return bo, nil
}

func GetOrgBoByCode(code string) (*bo.OrganizationBo, errs.SystemErrorInfo) {
	po := &po.PpmOrgOrganization{}
	err := mysql.SelectOneByCond(consts.TableOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcCode:     code,
	}, po)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bo := &bo.OrganizationBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func ScheduleOrganizationPageList(size int, page int) (*[]*bo.ScheduleOrganizationListBo, int64, errs.SystemErrorInfo) {

	pos, count, err := dao.OrcConfigPageList(size, page)

	if err != nil {
		return nil, int64(count), errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	bos := &[]*bo.ScheduleOrganizationListBo{}

	err = copyer.Copy(pos, bos)

	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return bos, count, nil

}

//校验当前用户是否有效
func VerifyOrg(orgId int64, userId int64) bool {
	baseUserInfo, err := GetBaseUserInfo("", orgId, userId)
	if err != nil {
		log.Error(err)
		return false
	}
	//未被移除且通过审核，认为当前用户有效
	return baseUserInfo.OrgUserIsDelete == consts.AppIsNoDelete && baseUserInfo.OrgUserCheckStatus == consts.AppCheckStatusSuccess
}

func VerifyOrgUsers(orgId int64, userIds []int64) bool {
	userIds = slice.SliceUniqueInt64(userIds)
	baseUserInfos, err := GetBaseUserInfoBatch("", orgId, userIds)
	if err != nil {
		log.Error(err)
		return false
	}
	if len(baseUserInfos) != len(userIds) {
		log.Error("部分用户无效")
		return false
	}
	for _, userInfo := range baseUserInfos {
		//如果存在待审核或者已移除的，则不允许更新
		if userInfo.OrgUserIsDelete != consts.AppIsNoDelete || userInfo.OrgUserCheckStatus != consts.AppCheckStatusSuccess {
			return false
		}
	}
	return true
}

func GetOrgByOutOrgId(sourceChannel, outOrgId string) (*bo.OrganizationBo, errs.SystemErrorInfo) {
	outInfo := &po.PpmOrgOrganizationOutInfo{}
	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcSourceChannel: sourceChannel,
		consts.TcOutOrgId:      outOrgId,
	}, outInfo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	orgId := outInfo.OrgId
	orgPo := &po.PpmOrgOrganization{}
	err = mysql.SelectById(consts.TableOrganization, orgId, orgPo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	orgBo := &bo.OrganizationBo{}
	_ = copyer.Copy(orgPo, orgBo)
	return orgBo, nil
}

func CreateOrg(createOrgBo bo.CreateOrgBo, creatorId int64, sourceChannel, sourcePlatform string, outOrgId string) (int64, errs.SystemErrorInfo) {
	orgId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrganization)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return 0, err
	}

	orgName := strings.TrimSpace(createOrgBo.OrgName)
	//orgNameLen := strs.Len(orgName)
	//if orgNameLen == 0 || orgNameLen > 256 {
	//	log.Error("组织名称长度错误")
	//	return 0, errs.BuildSystemErrorInfo(errs.OrgNameLenError)
	//}
	isOrgNameRight := format.VerifyOrgNameFormat(orgName)
	if !isOrgNameRight {
		return 0, errs.OrgNameLenError
	}

	//组织
	org := &po.PpmOrgOrganization{}
	org.Id = orgId
	org.Status = consts.AppStatusEnable
	org.IsDelete = consts.AppIsNoDelete
	org.Creator = creatorId
	org.Owner = creatorId
	org.Updator = creatorId
	org.Name = orgName
	org.SourceChannel = sourceChannel
	org.SourcePlatform = sourcePlatform

	err1 := mysql.Insert(org)
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	orgOutId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrganizationOutInfo)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return 0, err
	}
	//外部组织信息
	orgOutInfo := &po.PpmOrgOrganizationOutInfo{
		Id:             orgOutId,
		OrgId:          orgId,
		OutOrgId:       outOrgId,
		SourceChannel:  sourceChannel,
		SourcePlatform: sourcePlatform,
		Name:           createOrgBo.OrgName,
		Creator:        creatorId,
		Updator:        creatorId,
	}

	err1 = mysql.Insert(orgOutInfo)
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	//组织配置
	orgConfigId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrgConfig)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return 0, err
	}
	//组织配置信息
	orgConfig := &po.PpmOrcConfig{
		Id:    orgConfigId,
		OrgId: orgId,
	}
	err1 = mysql.Insert(orgConfig)
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return orgId, nil
}

func GetOrgInfoByOutOrgId(outOrgId string, sourceChannel string) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	outOrgInfo := &po.PpmOrgOrganizationOutInfo{}
	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcOutOrgId:      outOrgId,
		consts.TcSourceChannel: sourceChannel,
		consts.TcIsDelete:      consts.AppIsNoDelete,
	}, outOrgInfo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.OrgOutInfoNotExist)
	}
	orgInfo, err1 := GetOrgBoById(outOrgInfo.OrgId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.OrgNotExist)
	}
	return &bo.BaseOrgInfoBo{
		OrgId:         orgInfo.Id,
		OrgName:       orgInfo.Name,
		OutOrgId:      outOrgId,
		SourceChannel: sourceChannel,
	}, nil
}

func UpdateOrg(updateBo bo.UpdateOrganizationBo) errs.SystemErrorInfo {

	organizationBo := updateBo.Bo
	upds := updateBo.OrganizationUpdateCond

	_, err := mysql.UpdateSmartWithCond(consts.TableOrganization, db.Cond{
		consts.TcId: organizationBo.Id,
	}, upds)

	if err != nil {
		log.Errorf("mysql.TransUpdateSmart: %q\n", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	err = ClearCacheBaseOrgInfo(organizationBo.SourceChannel, organizationBo.Id)

	if err != nil {
		log.Errorf("redis err: %q\n", err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}
