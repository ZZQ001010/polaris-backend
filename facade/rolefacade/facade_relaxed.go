package rolefacade

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
)

func GetUserAdminFlagRelaxed(orgId int64, userId int64) (*bo.UserAdminFlagBo, errs.SystemErrorInfo) {
	respVo := GetUserAdminFlag(rolevo.GetUserAdminFlagReqVo{
		OrgId: orgId,
		UserId: userId,
	})
	if respVo.Failure(){
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.Data, nil
}