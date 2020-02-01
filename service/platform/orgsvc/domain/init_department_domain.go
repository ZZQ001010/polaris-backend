package domain

import (
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	dingtalk2 "github.com/galaxy-book/polaris-backend/common/extra/dingtalk"
	"github.com/galaxy-book/polaris-backend/common/extra/feishu"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
	"strconv"
	"upper.io/db.v3/lib/sqlbuilder"
)

//初始化部门
func InitDepartment(orgId int64, outOrgId string, sourceChannel string, superAdminRoleId, normalAdminRoleId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	if sourceChannel == consts.AppSourceChannelDingTalk {
		return InitDingTalkDepartment(orgId, outOrgId, superAdminRoleId, normalAdminRoleId, tx)
	} else if sourceChannel == consts.AppSourceChannelFeiShu {
		return InitFsDepartment(orgId, outOrgId, superAdminRoleId, normalAdminRoleId, tx)
	}
	return errs.BuildSystemErrorInfo(errs.SourceChannelNotDefinedError)
}

func InitDingTalkDepartment(orgId int64, corpId string, superAdminRoleId, normalAdminRoleId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	deptList, err := dingtalk2.GetScopeDeps(corpId)
	if err != nil {
		return err
	}
	//把组织作为根部门插入
	client, err1 := dingtalk2.GetDingTalkClientRest(corpId)
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.DingTalkClientError, err1)
	}

	deptMap := map[int64]sdk.DepartmentInfo{}
	for _, dep := range deptList{
		deptMap[dep.Id] = dep
	}

	//判断根部门是否存在
	rootDep, ok := deptMap[1]
	if ! ok{
		depResp, rootErr := client.GetDeptDetail("1", nil)
		if rootErr != nil {
			log.Error(rootErr)
			return errs.BuildSystemErrorInfo(errs.DingTalkClientError, rootErr)
		}
		if depResp.ErrCode != 0{
			log.Error(depResp.ErrMsg)
			return errs.DingTalkClientError
		}
		rootDep = sdk.DepartmentInfo{
			Id: depResp.Id,
			Name: depResp.Name,
			ParentId: -1,
			CreateDeptGroup: depResp.CreateDeptGroup,
			AutoAddUser: depResp.AutoAddUser,
		}
		deptList = append(deptList, rootDep)
	}

	depSize := len(deptList)
	departmentInfo := make([]interface{}, len(deptList))
	outDepartmentInfo := make([]interface{}, len(deptList))
	fsDepIdMap := map[int64]int64{}

	depIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableDepartment, depSize)
	if err != nil {
		log.Error(err)
		return err
	}

	depOutIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableDepartmentOutInfo, depSize)
	if err != nil {
		log.Error(err)
		return err
	}

	for k, v := range deptList {
		fsDepIdMap[v.Id] = depIds.Ids[k].Id
	}

	rootId := fsDepIdMap[1]

	for k, v := range deptList {
		depId := depIds.Ids[k].Id
		depOutId := depOutIds.Ids[k].Id

		parentDepId := int64(0)

		if id, ok := fsDepIdMap[v.ParentId]; ok{
			parentDepId = id
		}else{
			parentDepId = rootId
		}
		if depId == rootId{
			parentDepId = 0
		}

		departmentInfo[k] = &po.PpmOrgDepartment{
			Id:            depId,
			OrgId:         orgId,
			Name:          v.Name,
			ParentId:      parentDepId,
			SourceChannel: consts.AppSourceChannelDingTalk,
		}

		outDepartmentInfo[k] = po.PpmOrgDepartmentOutInfo{
			Id:                       depOutId,
			OrgId:                    orgId,
			DepartmentId:             depId,
			SourceChannel:            consts.AppSourceChannelDingTalk,
			OutOrgDepartmentId:       strconv.FormatInt(v.Id, 10),
			Name:                     v.Name,
			OutOrgDepartmentParentId: strconv.FormatInt(v.ParentId, 10),
		}
	}

	//初始化用户
	userInitErr := InitDingTalkUserList(orgId, corpId, fsDepIdMap, superAdminRoleId, normalAdminRoleId, tx)
	if userInitErr != nil {
		log.Error(userInitErr)
		return userInitErr
	}

	departErr := mysql.TransBatchInsert(tx, &po.PpmOrgDepartment{}, departmentInfo)
	if departErr != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, departErr)
	}
	outDepartErr := mysql.TransBatchInsert(tx, &po.PpmOrgDepartmentOutInfo{}, outDepartmentInfo)
	if outDepartErr != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, outDepartErr)
	}
	return nil
}

func InitFsDepartment(orgId int64, tenantKey string, superAdminRoleId, normalAdminRoleId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	deptList, err := feishu.GetScopeDeps(tenantKey)
	if err != nil {
		log.Error(err)
		return err
	}
	deptMap := map[string]vo.DepartmentRestInfoVo{}
	for _, dep := range deptList{
		deptMap[dep.Id] = dep
	}

	_, ok := deptMap["0"]
	if !ok{
		deptList = append(deptList, vo.DepartmentRestInfoVo{
			Id: "0",
			Name: "飞书平台组织",
			ParentId: "Root_Department_Identification",
		})
	}

	depSize := len(deptList)

	departmentInfo := make([]interface{}, len(deptList))
	outDepartmentInfo := make([]interface{}, len(deptList))
	fsDepIdMap := map[string]int64{}

	depIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableDepartment, depSize)
	if err != nil {
		log.Error(err)
		return err
	}

	depOutIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableDepartmentOutInfo, depSize)
	if err != nil {
		log.Error(err)
		return err
	}

	for k, v := range deptList {
		fsDepIdMap[v.Id] = depIds.Ids[k].Id
	}

	rootId := fsDepIdMap["0"]

	for k, v := range deptList {
		depId := depIds.Ids[k].Id
		depOutId := depOutIds.Ids[k].Id

		parentDepId := int64(0)

		if id, ok := fsDepIdMap[v.ParentId]; ok{
			parentDepId = id
		}else{
			parentDepId = rootId
		}
		if depId == rootId{
			parentDepId = 0
		}

		departmentInfo[k] = &po.PpmOrgDepartment{
			Id:            depId,
			OrgId:         orgId,
			Name:          v.Name,
			ParentId:      parentDepId,
			SourceChannel: consts.AppSourceChannelFeiShu,
		}

		outDepartmentInfo[k] = po.PpmOrgDepartmentOutInfo{
			Id:                       depOutId,
			OrgId:                    orgId,
			DepartmentId:             depId,
			SourceChannel:            consts.AppSourceChannelFeiShu,
			OutOrgDepartmentId:       v.Id,
			Name:                     v.Name,
			OutOrgDepartmentParentId: v.ParentId,
		}
	}

	//初始化用户
	userInitErr := InitFsUserList(orgId, tenantKey, fsDepIdMap, superAdminRoleId, normalAdminRoleId, tx)
	if userInitErr != nil {
		log.Error(userInitErr)
		return userInitErr
	}

	departErr := mysql.TransBatchInsert(tx, &po.PpmOrgDepartment{}, departmentInfo)
	if departErr != nil {
		log.Error(departErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, departErr)
	}
	outDepartErr := mysql.TransBatchInsert(tx, &po.PpmOrgDepartmentOutInfo{}, outDepartmentInfo)
	if outDepartErr != nil {
		log.Error(outDepartErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, outDepartErr)
	}
	return nil
}
