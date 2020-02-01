package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/domain"
	"upper.io/db.v3"
)

func Departments(page uint, size uint, params *vo.DepartmentListReq, orgId int64) (*vo.DepartmentList, errs.SystemErrorInfo) {

	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcIsHide:   consts.AppIsNotHiding, //默认只查询非隐藏部门
	}

	if params != nil {
		//查询父部门的子部门信息
		if params.ParentID != nil {
			cond[consts.TcParentId] = params.ParentID
		}
		//名称
		if params.Name != nil {
			cond[consts.TcName] = db.Like("%" + *params.Name + "%")
		}
		//查询顶级部门
		if params.IsTop != nil && *params.IsTop == 1 {
			cond[consts.TcParentId] = 0
		}
		//展示隐藏的部门
		if params.ShowHiding != nil && *params.ShowHiding == 1 {
			delete(cond, consts.TcIsHide)
		}
	}

	departmentBos, total, err := domain.GetDepartmentBoList(page, size, cond)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.DepartmentDomainError, err)
	}

	resultList := &[]*vo.Department{}
	copyErr := copyer.Copy(departmentBos, resultList)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &vo.DepartmentList{
		Total: total,
		List:  *resultList,
	}, nil
}

func DepartmentMembers(params vo.DepartmentMemberListReq, orgId int64) ([]*vo.DepartmentMemberInfo, errs.SystemErrorInfo) {
	//currentUserInfo, err := GetCurrentUser(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//orgId := currentUserInfo.OrgId
	departmentId := params.DepartmentID

	//departmentBo, err := domain.GetDepartmentBoWithOrg(departmentId, orgId)
	//if err != nil{
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.DepartmentNotExist)
	//}

	userIdInfoBoList, err := domain.GetDepartmentMembers(orgId, departmentId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.DepartmentDomainError, err)
	}

	userIdInfoVos := &[]*vo.DepartmentMemberInfo{}
	err1 := copyer.Copy(userIdInfoBoList, userIdInfoVos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *userIdInfoVos, nil
}
