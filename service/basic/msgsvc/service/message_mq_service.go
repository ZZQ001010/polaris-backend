package service

import (
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/mq"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/msgvo"
	"github.com/galaxy-book/polaris-backend/service/basic/msgsvc/domain"
	"upper.io/db.v3"
)

func PushMsgToMq(msg msgvo.PushMsgToMqReqVo) errs.SystemErrorInfo {
	mqMsg := msg.Msg
	mqMsgType := msg.MsgType
	mqClient := *mq.GetMQClient()
	key := mqMsg.Keys

	log.Infof("kafka配置 %s", json.ToJsonIgnoreError(config.GetMQ()))
	_, err1 := mqClient.PushMessage(&mqMsg)
	if err1 != nil {
		log.Errorf("消息推送失败, key: %s，准备入表，推送失败原因%v", key, err1)
		//落库
		id, err1 := domain.InsertMqFailMsgToDB(mqMsg, mqMsgType, msg.OrgId)
		if err1 != nil {
			log.Errorf("消息入表失败, key: %s，原因：%v", key, err1)
			return err1
		}
		log.Infof("消息入表成功, key: %s，id: %d", key, id)
	}
	return nil
}

func InsertMqConsumeFailMsg(msg msgvo.InsertMqConsumeFailMsgReqVo) errs.SystemErrorInfo {
	mqMsg := msg.Msg
	mqMsgType := msg.MsgType
	key := mqMsg.Keys
	log.Errorf("消息消费失败，准备入表, key: %s", key)
	//落库
	id, err1 := domain.InsertMqFailMsgToDB(mqMsg, mqMsgType, msg.OrgId)
	if err1 != nil {
		log.Errorf("消息入表失败, key: %s，原因：%v", key, err1)
		return err1
	}
	log.Infof("消息入表成功, key: %s，id: %d", key, id)
	return nil
}

func GetFailMsgList(req msgvo.GetFailMsgListReqVo) (*[]bo.MessageBo, errs.SystemErrorInfo) {
	pageA, sizeA := util.PageOption(req.Page, req.Size)
	cond := db.Cond{}
	if req.OrgId != nil {
		cond[consts.TcOrgId] = *req.OrgId
	}
	if req.MsgType != nil {
		cond[consts.TcType] = *req.MsgType
	}

	list, _, err := domain.GetMessageBoList(pageA, sizeA, cond)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return list, nil
}

func UpdateMsgStatus(req msgvo.UpdateMsgStatusReqVo) errs.SystemErrorInfo {
	return domain.UpdateMsgStatus(req.MsgId, req.NewStatus)
}
