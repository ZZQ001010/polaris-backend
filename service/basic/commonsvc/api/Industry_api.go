package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/commonvo"
	"github.com/galaxy-book/polaris-backend/service/basic/commonsvc/service"
)

func (GetGreeter) IndustryList() commonvo.IndustryListRespVo {
	res, err := service.IndustryList()
	return commonvo.IndustryListRespVo{Err: vo.NewErr(err), IndustryList: res}
}
