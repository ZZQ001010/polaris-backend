package projectfacade

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
)

func GetProjectBoListByProjectTypeLangCodeRelaxed(orgId int64, projectTypeLangCode *string) ([]bo.ProjectBo, errs.SystemErrorInfo) {
	respVo := GetProjectBoListByProjectTypeLangCode(projectvo.GetProjectBoListByProjectTypeLangCodeReqVo{
		OrgId: orgId,
		ProjectTypeLangCode: projectTypeLangCode,
	})

	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.ProjectBoList, nil
}

func GetPriorityByIdRelaxed(orgId int64, id int64) (*bo.PriorityBo, errs.SystemErrorInfo) {
	respVo := GetPriorityById(projectvo.GetPriorityByIdReqVo{OrgId: orgId, Id: id})
	if respVo.Failure(){
		return nil, respVo.Error()
	}
	return respVo.PriorityBo, nil
}

func InitPriorityRelaxed(orgId int64) errs.SystemErrorInfo {
	respVo := InitPriority(projectvo.InitPriorityReqVo{OrgId: orgId})
	if respVo.Failure(){
		return respVo.Error()
	}
	return nil
}

func VerifyPriorityRelaxed(orgId int64, typ int, beVerifyId int64) (bool, errs.SystemErrorInfo) {
	respVo := VerifyPriority(projectvo.VerifyPriorityReqVo{OrgId: orgId, Typ: typ, BeVerifyId: beVerifyId})
	if respVo.Failure(){
		return false, respVo.Error()
	}
	return respVo.Successful, nil
}
