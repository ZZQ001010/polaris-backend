package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/msgvo"
	"github.com/galaxy-book/polaris-backend/service/basic/msgsvc/service"
)

func (PostGreeter) SendLoginSMS(req msgvo.SendLoginSMSReqVo) vo.VoidErr{
	err := service.SendLoginSMS(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) SendMail(req msgvo.SendMailReqVo) vo.VoidErr{
	err := service.SendMail(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) SendSMS(req msgvo.SendSMSReqVo) vo.VoidErr{
	err := service.SendSMS(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}