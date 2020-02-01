package consume

import (
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/errors"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/model"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/mq"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo/mqbo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/msgfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
)

var log = logger.GetDefaultLogger()

func ImportIssueConsume() {

	log.Infof("mq消息-任务动态消费者启动成功")

	importIssueTopicConfig := config.GetMqImportIssueTopicConfig()

	client := *mq.GetMQClient()
	_ = client.ConsumeMessage(importIssueTopicConfig.Topic, importIssueTopicConfig.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("捕获到的错误：%s\n", r)
			}
		}()

		log.Infof("mq消息-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		createIssueBo := &mqbo.PushCreateIssueBo{}
		err := json.FromJson(message.Body, createIssueBo)
		if err != nil{
			log.Error(err)
			return errs.JSONConvertError
		}

		issueVo, respErr := service.CreateIssueWithId(createIssueBo.CreateIssueReqVo, createIssueBo.IssueId)
		if respErr != nil{
			log.Error(respErr)
			return respErr
		}

		log.Infof("创建任务成功 %s", json.ToJsonIgnoreError(issueVo))

		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage

		log.Infof("mq消息消费失败-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		createIssueReq := &projectvo.CreateIssueReqVo{}
		err := json.FromJson(message.Body, createIssueReq)
		if err != nil{
			log.Error(err)
		}

		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, 0, createIssueReq.OrgId)
		if msgErr != nil {
			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})
}
