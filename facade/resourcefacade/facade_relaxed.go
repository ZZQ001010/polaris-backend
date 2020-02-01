package resourcefacade

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
)

func GetResourceBoListRelaxed(page uint, size uint, cond resourcevo.GetResourceBoListCond) (*[]bo.ResourceBo, int64, errs.SystemErrorInfo) {
	respVo := GetResourceBoList(resourcevo.GetResourceBoListReqVo{
		Page: page,
		Size: size,
		Input: cond,
	})

	if respVo.Failure() {
		return nil, 0, respVo.Error()
	}

	return respVo.ResourceBos, respVo.Total, nil
}
