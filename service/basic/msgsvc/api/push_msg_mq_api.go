package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/msgvo"
	"github.com/galaxy-book/polaris-backend/service/basic/msgsvc/service"
)

func (PostGreeter) PushMsgToMq(msg msgvo.PushMsgToMqReqVo) vo.VoidErr{
	err := service.PushMsgToMq(msg)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) InsertMqConsumeFailMsg(msg msgvo.InsertMqConsumeFailMsgReqVo) vo.VoidErr{
	err := service.InsertMqConsumeFailMsg(msg)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (GetGreeter) GetFailMsgList(req msgvo.GetFailMsgListReqVo) msgvo.GetFailMsgListRespVo{
	res, err := service.GetFailMsgList(req)
	return msgvo.GetFailMsgListRespVo{Data:msgvo.GetFailMsgListRespData{MsgList:res,}, Err: vo.NewErr(err)}
}

func (PostGreeter) UpdateMsgStatus(req msgvo.UpdateMsgStatusReqVo) vo.VoidErr{
	err := service.UpdateMsgStatus(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}