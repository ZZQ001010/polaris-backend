package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	selfConsts "github.com/galaxy-book/polaris-backend/service/platform/rolesvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/domain"
	"strings"
)

func PermissionOperationList(orgId, roleId, userId int64, projectId *int64) ([]*vo.PermissionOperationListResp, errs.SystemErrorInfo) {
	roleInfo, err := domain.GetRole(orgId, 0, roleId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var permissionBo []bo.PermissionBo
	var permissionErr errs.SystemErrorInfo
	if projectId != nil && *projectId != 0 {
		//获取项目权限项
		permissionBo, permissionErr = domain.GetProjectPermission()
		if permissionErr != nil {
			return nil, permissionErr
		}
	} else {
		//获取组织权限项
		permissionBo, permissionErr = domain.GetPermissionByType(consts.PermissionTypeOrg)
		if permissionErr != nil {
			return nil, permissionErr
		}
	}

	//获取角色所有的操作权限
	rolePermissionOperation, err := GetRolePermissionOperationList(orgId, roleId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	roleOperationMap := maps.NewMap("PermissionId", *rolePermissionOperation)

	resBo := []bo.PermissionOperationListBo{}
	allCanUse := false
	for _, v := range permissionBo {
		if v.ParentId == 0 {
			if _, ok := roleOperationMap[v.Id]; ok {
				//默认设置父级权限，则拥有所有子权限
				allCanUse = true
			}
			continue
		}
		//如果是查询项目的，且角色为负责人
		if projectId != nil && *projectId != 0 && roleInfo.LangCode == consts.RoleGroupSpecialOwner {
			allCanUse = true
		}

		allPermissionHave := []int64{}
		operation, err := domain.GetPermissionOperationListByPermissionId(v.Id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if allCanUse {
			for _, val := range operation {
				allPermissionHave = append(allPermissionHave, val.Id)
			}
		} else if _, ok := roleOperationMap[v.Id]; ok {
			roleOperation := roleOperationMap[v.Id].(bo.RolePermissionOperationBo)
			for _, val := range operation {
				if judgeOperation(val.OperationCodes, roleOperation.OperationCodes) {
					allPermissionHave = append(allPermissionHave, val.Id)
				}
			}
		}

		resBo = append(resBo, bo.PermissionOperationListBo{
			PermissionInfo: v,
			OperationList:  operation,
			PermissionHave: allPermissionHave,
		})
	}

	resVo := &[]*vo.PermissionOperationListResp{}
	copyErr := copyer.Copy(resBo, resVo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return *resVo, nil
}

func judgeOperation(operation string, allOperation string) bool {
	allArr := []string{}
	if allOperation == "*" {
		return true
	} else if strings.Index(allOperation, "|") == -1 {
		allArr = []string{allOperation}
	} else {
		mid := strings.Split(allOperation, "|")
		for _, v := range mid {
			if len(v) > 2 {
				allArr = append(allArr, v[1:len(v)-1])
			}
		}
	}

	if strings.Index(operation, ",") == -1 {
		if ok, _ := slice.Contain(allArr, operation); ok {
			return true
		}
	} else {
		mid := strings.Split(operation, ",")
		for _, v := range mid {
			if ok, _ := slice.Contain(allArr, v); ok {
				return true
			}
		}
	}

	return false
}

func UpdateRolePermissionOperation(orgId int64, userId int64, input vo.UpdateRolePermissionOperationReq) (*vo.Void, errs.SystemErrorInfo) {
	// 1.系统权限判断
	authErr := Authenticate(orgId, userId, nil, nil, consts.RoleOperationPathOrgOrgConfig, consts.RoleOperationModify)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}

	if len(input.UpdatePermissions) == 0 {
		log.Info("权限无更新")
		return &vo.Void{ID: 0}, nil
	}

	roleInfo, err := domain.GetRole(orgId, 0, input.RoleID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//超管和组织超管权限不可编辑且默认最大，所以不考虑其权限修改的问题
	if selfConsts.IsDefaultRole(roleInfo.LangCode) {
		return nil, errs.OrgUserRoleModifyError
	}

	permissionOperation := map[int64][]string{}
	for _, permission := range input.UpdatePermissions {
		//拼装角色在该权限的操作项
		operation, err := domain.GetPermissionOperationListByPermissionId(permission.PermissionID)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		var operationCode []string
		//判断操作项
		for _, operationBo := range operation {
			if ok, _ := slice.Contain(permission.OperationIds, operationBo.Id); !ok {
				continue
			}
			operationCode = append(operationCode, strings.Split(operationBo.OperationCodes, ",")...)
		}
		permissionOperation[permission.PermissionID] = operationCode
	}

	updateErr := domain.UpdateRolePermissionOperation(orgId, userId, input.RoleID, permissionOperation)
	if updateErr != nil {
		log.Error(updateErr)
		return nil, updateErr
	}

	//清除缓存
	clearErr := ClearRolePermissionOperationList(orgId, input.RoleID)
	if clearErr != nil {
		log.Error(clearErr)
		return nil, clearErr
	}
	return &vo.Void{ID: 0}, nil
}

//获取个人权限信息
func GetPersonalPermissionInfo(orgId, userId int64, projectId, issueId *int64, sourceChannel string) (map[string]interface{}, errs.SystemErrorInfo) {
	var permissionBo []bo.PermissionBo
	var permissionErr errs.SystemErrorInfo
	var isProjectOwner = false
	var isIssueOwner = false
	var newProjectId int64
	var projectAuthInfo *bo.ProjectAuthBo
	var issueAuthInfo *bo.IssueAuthBo
	if projectId != nil && *projectId != 0 {
		newProjectId = *projectId
		//获取项目权限项
		permissionBo, permissionErr = domain.GetProjectPermission()
		if permissionErr != nil {
			return nil, permissionErr
		}

		//获取项目信息
		projectInfo := projectfacade.GetCacheProjectInfo(projectvo.GetCacheProjectInfoReqVo{
			ProjectId:*projectId,
			OrgId:         orgId,
		})
		if projectInfo.Failure() {
			log.Error(projectInfo.Error())
			return nil, projectInfo.Error()
		}
		projectAuthInfo = projectInfo.ProjectCacheBo
		if projectInfo.ProjectCacheBo.Owner == userId {
			isProjectOwner = true
		}
	} else {
		//获取组织权限项
		permissionBo, permissionErr = domain.GetPermissionByType(consts.PermissionTypeOrg)
		if permissionErr != nil {
			return nil, permissionErr
		}
	}

	if issueId != nil && *issueId != 0 {
		issueInfo := projectfacade.IssueInfo(projectvo.IssueInfoReqVo{
			IssueID:*issueId,
			UserId:userId,
			OrgId:orgId,
			SourceChannel:sourceChannel,
		})
		if issueInfo.Failure() {
			log.Error(issueInfo.Error())
			return nil, issueInfo.Error()
		}
		issueAuthInfo := &bo.IssueAuthBo{
			Owner:issueInfo.IssueInfo.Issue.Owner,
			ProjectId:issueInfo.IssueInfo.Issue.ProjectID,
			Id:issueInfo.IssueInfo.Issue.ID,
			Status:issueInfo.IssueInfo.Issue.Status,
			Creator:issueInfo.IssueInfo.Issue.Creator,
		}
		//issueAuthInfo.Owner = issueInfo.IssueInfo.Issue.Owner
		//issueAuthInfo.ProjectId = issueInfo.IssueInfo.Issue.ProjectID
		//issueAuthInfo.Id = issueInfo.IssueInfo.Issue.ID
		//issueAuthInfo.Status = issueInfo.IssueInfo.Issue.Status
		//issueAuthInfo.Creator = issueInfo.IssueInfo.Issue.Creator
		for _, info := range issueInfo.IssueInfo.FollowerInfos {
			issueAuthInfo.Followers = append(issueAuthInfo.Followers, info.UserID)
		}
		for _, info := range issueInfo.IssueInfo.ParticipantInfos {
			issueAuthInfo.Participants = append(issueAuthInfo.Participants, info.UserID)
		}

		if issueInfo.IssueInfo.Issue.Owner == userId {
			isIssueOwner = true
		}
	}

	//获取用户所有角色
	roleIds, roleErr := GetUserRoleIds(orgId, userId, newProjectId, projectAuthInfo, issueAuthInfo)
	log.Infof("用户%d所有角色:%s", userId, json.ToJsonIgnoreError(roleIds))
	if roleErr != nil {
		log.Error(roleErr)
		return nil, errs.BuildSystemErrorInfo(errs.GetUserRoleError, roleErr)
	}

	//获取角色所有的操作权限
	roleOperationMap := map[int64][]string{}
	for _, roleId := range roleIds {
		rolePermissionOperation, err := GetRolePermissionOperationList(orgId, roleId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		for _, operationBo := range *rolePermissionOperation {
			allCode := OperationCodeToSlice(operationBo.OperationCodes)
			if _, ok := roleOperationMap[operationBo.PermissionId]; ok {
				if len(roleOperationMap[operationBo.PermissionId]) == 1 && roleOperationMap[operationBo.PermissionId][0] == "*" {
					continue
				}
				if len(allCode) == 1 && allCode[0] == "*" {
					roleOperationMap[operationBo.PermissionId] = allCode
				} else {
					roleOperationMap[operationBo.PermissionId] = append(roleOperationMap[operationBo.PermissionId], allCode...)
				}
			} else {
				roleOperationMap[operationBo.PermissionId] = allCode
			}
		}
	}
	//去重
	for i, i2 := range roleOperationMap {
		roleOperationMap[i] = slice.SliceUniqueString(i2)
	}

	resBo := map[string][]string{}
	allCanUse := false
	for _, v := range permissionBo {
		code := v.Code
		if v.LangCode == consts.PermissionProIssue4 {
			code = "Issue"
		}
		if v.ParentId == 0 {
			//默认设置父级权限，则拥有所有子权限
			allCanUse = true
			continue
		}
		//如果是查询项目的，且角色为负责人
		if projectId != nil && *projectId != 0 && isProjectOwner {
			allCanUse = true
		}

		operation, err := domain.GetPermissionOperationListByPermissionId(v.Id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		//把每个权限的操作项改造成数组
		operationArr := []string{}
		for _, val := range operation {
			if strings.Index(val.OperationCodes, ",") == -1 {
				operationArr = append(operationArr, val.OperationCodes)
			} else {
				mid := strings.Split(val.OperationCodes, ",")
				operationArr = append(operationArr, mid...)
			}
		}

		if allCanUse {
			resBo[code] = append(resBo[code], operationArr...)
		} else if _, ok := roleOperationMap[v.Id]; ok {
			if len(roleOperationMap[v.Id]) == 1 && roleOperationMap[v.Id][0] == "*" {
				resBo[code] = append(resBo[code], operationArr...)
			} else {
				for _, val := range operationArr {
					if ok, _ := slice.Contain(roleOperationMap[v.Id], val); ok {
						resBo[code] = append(resBo[code], val)
					}
				}
			}
		}

		//如果是任务负责人，默认赋予任务和附件所有权限
		if ok, _ := slice.Contain([]string{consts.PermissionProIssue4, consts.PermissionProAttachment}, v.LangCode); ok {
			if !isIssueOwner {
				continue
			}
			resBo[code] = append(resBo[code], operationArr...)
		}
	}

	res := &map[string]interface{}{}
	copyErr := copyer.Copy(resBo, res)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return *res, nil

}

func OperationCodeToSlice(allOperation string) []string {
	allArr := []string{}
	if allOperation == "*" || strings.Index(allOperation, "|") == -1 {
		return []string{allOperation}
	} else {
		mid := strings.Split(allOperation, "|")
		for _, v := range mid {
			if len(v) > 2 {
				allArr = append(allArr, v[1:len(v)-1])
			}
		}
	}

	return allArr
}