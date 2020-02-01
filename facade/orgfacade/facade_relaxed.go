package orgfacade

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
)

func GetBaseUserInfoBatchRelaxed(sourceChannal string, orgId int64, userIds []int64) ([]bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	respVo := GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		SourceChannel: sourceChannal,
		OrgId: orgId,
		UserIds: userIds,
	})
	return respVo.BaseUserInfos, respVo.Error()
}