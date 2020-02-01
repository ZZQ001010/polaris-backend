package consume

import (
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/errors"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/model"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/mq"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/msgfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/domain"
)

var log = logger.GetDefaultLogger()

func OrgMemberChangeConsume() {

	log.Infof("mq消息-组织成员变动消费者启动成功")

	orgMemberChangeTopicConfig := config.GetMqOrgMemberChangeConfig()

	client := *mq.GetMQClient()
	_ = client.ConsumeMessage(orgMemberChangeTopicConfig.Topic, orgMemberChangeTopicConfig.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("捕获到的错误：%s\n", r)
			}
		}()

		log.Infof("mq消息-组织成员变动消费信息 topic %s, value %s", message.Topic, message.Body)

		orgMemberChange := &bo.OrgMemberChangeBo{}
		err := json.FromJson(message.Body, orgMemberChange)
		if err != nil{
			log.Error(err)
			return errs.JSONConvertError
		}

		orgId := orgMemberChange.OrgId

		var businessErr errs.SystemErrorInfo = nil

		changeType := orgMemberChange.ChangeType
		//业务处理
		switch changeType {
		//禁用
		case consts.OrgMemberChangeTypeDisable:
			businessErr = domain.ModifyOrgMemberStatus(orgId, []int64{orgMemberChange.UserId}, consts.AppStatusDisabled, 0)
		//启用
		case consts.OrgMemberChangeTypeEnable:
			businessErr = domain.ModifyOrgMemberStatus(orgId, []int64{orgMemberChange.UserId}, consts.AppStatusEnable, 0)
		//添加用户
		case consts.OrgMemberChangeTypeAdd, consts.OrgMemberChangeTypeAddDisable:
			if orgMemberChange.SourceChannel == consts.AppSourceChannelFeiShu{
				orgOutInfo, err := domain.GetBaseOrgInfo(consts.AppSourceChannelFeiShu, orgId)
				if err != nil{
					log.Error(err)
					return err
				}
				baseUserInfo, err := domain.FsAuth(orgOutInfo.OutOrgId, orgMemberChange.OpenId)
				if err != nil{
					log.Error(err)
					return err
				}
				if baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted{
					inDisabled := changeType == consts.OrgMemberChangeTypeAddDisable
					log.Infof("添加用户是否被禁用 %v", inDisabled)
					err = domain.AddOrgMember(baseUserInfo.OrgId, baseUserInfo.UserId, 0, false, inDisabled)
					if err != nil{
						log.Error(err)
						return err
					}
				}
			}
		case consts.OrgMemberChangeTypeRemove:
			businessErr = domain.RemoveOrgMember(orgId, []int64{orgMemberChange.UserId}, 0)
		}

		if businessErr != nil{
			log.Error(businessErr)
		}

		//在并发操作时，有几率更新失败，所以忽略异常
		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage

		log.Infof("mq消息消费失败-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		orgMemberChange := &bo.OrgMemberChangeBo{}
		err := json.FromJson(message.Body, orgMemberChange)
		if err != nil{
			log.Error(err)
			return
		}

		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, 0, orgMemberChange.OrgId)
		if msgErr != nil {
			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})
}
