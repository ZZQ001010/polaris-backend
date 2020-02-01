package service

import (
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/rolefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
)

func RemoveProjectMember(orgId, userId int64, input vo.RemoveProjectMemberReq) (*vo.Void, errs.SystemErrorInfo) {
	//校验当前用户是否具有修改删除成员的权限
	authErr := AuthProjectPermission(orgId, userId, input.ProjectID, consts.RoleOperationPathOrgProMember, consts.RoleOperationUnbind, false)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}
	err := domain.RemoveProjectMember(orgId, userId, input)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &vo.Void{ID: 0}, nil
}

func AddProjectMember(orgId, userId int64, input vo.RemoveProjectMemberReq) (*vo.Void, errs.SystemErrorInfo) {
	//校验当前用户是否具有修改删除成员的权限
	authErr := AuthProjectPermission(orgId, userId, input.ProjectID, consts.RoleOperationPathOrgProMember, consts.RoleOperationBind, false)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}

	err := domain.AddProjectMember(orgId, userId, input)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &vo.Void{ID: 0}, nil
}

func ProjectUserList(orgId int64, page int, size int, input vo.ProjectUserListReq) (*vo.ProjectUserListResp, errs.SystemErrorInfo) {
	count, bos, err := domain.GetProjectAllMember(orgId, input.ProjectID, page, size)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	res := &vo.ProjectUserListResp{Total: count}

	if len(bos) == 0 {
		return res, nil
	}

	var allRelatedUser []int64
	for _, bo := range bos {
		allRelatedUser = append(allRelatedUser, bo.RelationId, bo.Creator)
	}

	//获取所有相关人员信息
	userInfo := orgfacade.BatchGetUserDetailInfo(orgvo.BatchGetUserInfoReq{UserIds: allRelatedUser})
	if userInfo.Failure() {
		log.Error(userInfo.Error())
		return nil, userInfo.Error()
	}
	userInfoVos := &[]vo.PersonalInfo{}
	copyErr := copyer.Copy(userInfo.Data, userInfoVos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	userInfoMap := maps.NewMap("ID", *userInfoVos)

	//特殊角色（项目成员和负责人）
	roleGroup := rolefacade.GetProjectRoleList(rolevo.GetProjectRoleListReqVo{
		OrgId:     orgId,
		ProjectId: input.ProjectID,
	})
	if roleGroup.Failure() {
		log.Error(roleGroup.Error())
		return nil, roleGroup.Error()
	}
	memberRoleInfo := vo.Role{}
	ownerRoleInfo := vo.Role{}
	for _, v := range roleGroup.Data {
		if v.LangCode == consts.RoleGroupProMember {
			memberRoleInfo = *v
		} else if v.LangCode == consts.RoleGroupSpecialOwner {
			ownerRoleInfo = *v
		}
	}

	//查询成员角色
	roleUserResp := rolefacade.GetOrgRoleUser(rolevo.GetOrgRoleUserReqVo{
		OrgId:     orgId,
		ProjectId: input.ProjectID,
	})
	if roleUserResp.Failure() {
		log.Error(roleUserResp.Error())
		return nil, roleUserResp.Error()
	}
	userRoleMap := map[int64]rolevo.RoleUser{}
	for _, datum := range roleUserResp.Data {
		if _, ok := userRoleMap[datum.UserId]; !ok {
			userRoleMap[datum.UserId] = datum
		}
	}
	for _, v := range bos {
		tempInfo := &vo.ProjectUser{}
		tempInfo.Creator = v.Creator
		tempInfo.CreateTime = types.Time(v.CreateTime)
		if _, ok := userInfoMap[v.Creator]; ok {
			user := userInfoMap[v.Creator].(vo.PersonalInfo)
			tempInfo.CreatorInfo = &user
		}
		if _, ok := userInfoMap[v.RelationId]; ok {
			user := userInfoMap[v.RelationId].(vo.PersonalInfo)
			tempInfo.UserInfo = &user
		}
		if v.RelationType == consts.IssueRelationTypeOwner {
			//负责人
			tempInfo.UserRole = &vo.UserRoleInfo{
				ID:       ownerRoleInfo.ID,
				Name:     ownerRoleInfo.Name,
				LangCode: ownerRoleInfo.LangCode,
			}
		} else if _, ok := userRoleMap[v.RelationId]; ok {
			//特殊角色
			role := userRoleMap[v.RelationId]
			tempInfo.UserRole = &vo.UserRoleInfo{
				ID:       role.RoleId,
				Name:     role.RoleName,
				LangCode: role.RoleLangCode,
			}
		} else {
			//项目成员
			tempInfo.UserRole = &vo.UserRoleInfo{
				ID:       memberRoleInfo.ID,
				Name:     memberRoleInfo.Name,
				LangCode: memberRoleInfo.LangCode,
			}
		}
		res.List = append(res.List, tempInfo)
	}

	return res, nil
}
