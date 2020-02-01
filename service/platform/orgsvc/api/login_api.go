package api

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/service"
	"strings"
)

func (PostGreeter) UserLogin(req orgvo.UserLoginReqVo) orgvo.UserSMSLoginRespVo {
	if req.UserLoginReq.Name != nil {
		creatorName := strings.TrimSpace(*req.UserLoginReq.Name)
		//creatorNameLen := str.CountStrByGBK(creatorName)
		//if creatorNameLen == 0 || creatorNameLen > 20{
		//	log.Error("姓名长度错误")
		//	return orgvo.UserSMSLoginRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.UserNameLenError)), Data: nil}
		//}
		isNameRight := format.VerifyUserNameFormat(creatorName)
		if !isNameRight {
			return orgvo.UserSMSLoginRespVo{Err: vo.NewErr(errs.UserNameLenError), Data: nil}
		}
		req.UserLoginReq.Name = &creatorName
	}

	res, err := service.UserLogin(req.UserLoginReq)
	return orgvo.UserSMSLoginRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) UserQuit(req orgvo.UserQuitReqVo) vo.VoidErr {
	err := service.UserQuit(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}
