package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) SetPassword(req orgvo.SetPasswordReqVo) vo.VoidErr {
	err := service.SetPassword(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) ResetPassword(req orgvo.ResetPasswordReqVo) vo.VoidErr {
	err := service.ResetPassword(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) RetrievePassword(req orgvo.RetrievePasswordReqVo) vo.VoidErr {
	err := service.RetrievePassword(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) UnbindLoginName(req orgvo.UnbindLoginNameReqVo) vo.VoidErr {
	err := service.UnbindLoginName(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}


func (PostGreeter) BindLoginName(req orgvo.BindLoginNameReqVo) vo.VoidErr {
	err := service.BindLoginName(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) CheckLoginName(req orgvo.CheckLoginNameReqVo) vo.VoidErr {
	err := service.CheckLoginName(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}