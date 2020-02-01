package domain

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
)

func AuthIssueWithIssueId(orgId, userId int64, issueId int64, path string, operation string) errs.SystemErrorInfo {
	issueBo, err1 := GetIssueBo(orgId, issueId)
	if err1 != nil {
		log.Error(err1)
		return err1
	}
	return AuthIssue(orgId, userId, *issueBo, path, operation)
}

//验证任务操作权，已归档的项目下的任务不允许编辑
func AuthIssue(orgId, userId int64, issueBo bo.IssueBo, path string, operation string) errs.SystemErrorInfo {
	issueAuthBo, err := GetIssueAuthBo(issueBo, userId)
	if err != nil {
		log.Error(err)
		return err
	}
	projectId := issueAuthBo.ProjectId

	projectAuthBo, err := LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return err
	}
	//校验私有项目
	if projectAuthBo.PublicStatus == consts.PrivateProject{
		authPrivateProjectErr := AuthPrivateProject(orgId, userId, projectAuthBo)
		if authPrivateProjectErr != nil{
			//判断当前用户是不是超管
			adminFlagBo, err := rolefacade.GetUserAdminFlagRelaxed(orgId, userId)
			if err != nil{
				log.Error(err)
				return authPrivateProjectErr
			}

			if ! adminFlagBo.IsAdmin{
				log.Error(authPrivateProjectErr)
				return authPrivateProjectErr
			}else{
				return nil
			}
		}
	}
	//校验项目是否归档
	if projectAuthBo.IsFilling == consts.AppIsFilling && operation != consts.RoleOperationView{
		return errs.BuildSystemErrorInfo(errs.ProjectIsArchivedWhenModifyIssue)
	}
	return rolefacade.AuthenticateRelaxed(orgId, userId, projectAuthBo, issueAuthBo, path, operation)
}

//同时会校验项目是否已归档
func AuthProject(orgId, userId, projectId int64, path string, operation string) errs.SystemErrorInfo {
	return AuthProjectWithCond(orgId, userId, projectId, path, operation, true,false)
}

//项目权限校验，不需要检测是否已归档
func AuthProjectWithOutArchivedCheck(orgId, userId, projectId int64, path string, operation string) errs.SystemErrorInfo {
	return AuthProjectWithCond(orgId, userId, projectId, path, operation, false,false)
}

//校验项目,跳过角色权限
func AuthProjectWithOutPermission(orgId, userId, projectId int64, path string, operation string) errs.SystemErrorInfo {
	return AuthProjectWithCond(orgId, userId, projectId, path, operation, false,true)
}

func AuthProjectWithCond(orgId, userId, projectId int64, path string, operation string, authFiling bool,skipAuthPermission bool) errs.SystemErrorInfo {
	projectAuthBo, err := LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	//校验私有项目
	if projectAuthBo.PublicStatus == consts.PrivateProject{
		authPrivateProjectErr := AuthPrivateProject(orgId, userId, projectAuthBo)
		if authPrivateProjectErr != nil{
			//判断当前用户是不是超管
			adminFlagBo, err := rolefacade.GetUserAdminFlagRelaxed(orgId, userId)
			if err != nil{
				log.Error(err)
				return authPrivateProjectErr
			}

			if ! adminFlagBo.IsAdmin{
				log.Error(authPrivateProjectErr)
				return authPrivateProjectErr
			}else{
				return nil
			}
		}
	}

	//校验项目是否归档
	if authFiling && projectAuthBo.IsFilling == consts.AppIsFilling && operation != consts.RoleOperationView{
		return errs.BuildSystemErrorInfo(errs.ProjectIsArchivedWhenModifyIssue)
	}
	if skipAuthPermission {
		return nil
	}else {
		return rolefacade.AuthenticateRelaxed(orgId, userId, projectAuthBo, nil, path, operation)
	}
}

//校验私有项目，只有成员才能编辑
func AuthPrivateProject(orgId, userId int64, projectAuthBo *bo.ProjectAuthBo) errs.SystemErrorInfo{
	memberIdMap := map[int64]bool{}
	memberIdMap[projectAuthBo.Owner] = true
	if projectAuthBo.Followers != nil{
		for _, follower := range projectAuthBo.Followers{
			memberIdMap[follower] = true
		}
	}
	if projectAuthBo.Participants != nil{
		for _, participant := range projectAuthBo.Participants{
			memberIdMap[participant] = true
		}
	}
	//如果项目成员不包括当前操作人
	if _, ok := memberIdMap[userId]; ! ok{
		return errs.BuildSystemErrorInfo(errs.NoPrivateProjectPermissions)
	}
	return nil
}

func AuthOrg(orgId, useId int64, path string, operation string) errs.SystemErrorInfo {
	return rolefacade.AuthenticateRelaxed(orgId, useId, nil, nil, path, operation)
}
