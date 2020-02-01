package domain

import (
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/model"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/bo/mqbo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/trendsvo"
	"github.com/galaxy-book/polaris-backend/facade/msgfacade"
	"github.com/galaxy-book/polaris-backend/facade/trendsfacade"
	"strconv"
	"time"
)

func PushIssueTrends(issueTrendsBo bo.IssueTrendsBo) {
	issueTrendsBo.OperateTime = time.Now()
	mqKeys := strconv.FormatInt(issueTrendsBo.OrgId, 10)
	//动态改成同步的
	resp := trendsfacade.AddIssueTrends(trendsvo.AddIssueTrendsReqVo{IssueTrendsBo: issueTrendsBo, Key: mqKeys})
	if resp.Failure(){
		log.Error(resp.Message)
	}
}

func PushIssueThirdPlatformNotice(issueTrendsBo bo.IssueTrendsBo){
	//issueTrendsBo.OperateTime = time.Now()
	//mqKeys := strconv.FormatInt(issueTrendsBo.OrgId, 10)
	//
	//orgId := issueTrendsBo.OrgId
	//pushType := int(issueTrendsBo.PushType)
	//issueTrendsBo.OperateTime = time.Now()
	//message, err := json.ToJson(issueTrendsBo)
	//if err != nil {
	//	log.Error(err)
	//}
	//
	//mqMessage := &model.MqMessage{
	//	Topic:          config.GetMqIssueTrendsTopicConfig().Topic,
	//	Keys:          mqKeys,
	//	Body:           message,
	//	DelayTimeLevel: 3,
	//}
	//
	//msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, pushType, orgId)
	//if msgErr != nil {
	//	log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	//}
}

func PushCreateIssue(createIssueBo mqbo.PushCreateIssueBo) {
	message, err := json.ToJson(createIssueBo)
	if err != nil {
		log.Error(err)
	}

	reqVo := createIssueBo.CreateIssueReqVo
	//这里key使用项目id，保证同一项目下导入的任务顺序的有效性
	mqMessage := &model.MqMessage{
		Topic:          config.GetMqImportIssueTopicConfig().Topic,
		Keys:           strconv.FormatInt(reqVo.CreateIssue.ProjectID, 10),
		Body:           message,
		DelayTimeLevel: 3,
	}
	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, 0, reqVo.OrgId)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}

func PushProjectTrends(projectMemberChangeBo bo.ProjectTrendsBo) {
	projectMemberChangeBo.OperateTime = time.Now()
	trendsfacade.AddProjectTrends(trendsvo.AddProjectTrendsReqVo{ProjectTrendsBo: projectMemberChangeBo})
}

func PushProjectThirdPlatformNotice(projectMemberChangeBo bo.ProjectTrendsBo){
	projectMemberChangeBo.OperateTime = time.Now()

	orgId := projectMemberChangeBo.OrgId
	pushType := int(projectMemberChangeBo.PushType)

	message, err := json.ToJson(projectMemberChangeBo)
	if err != nil {
		log.Error(err)
	}
	mqKeys := strconv.FormatInt(projectMemberChangeBo.OrgId, 10)

	mqMessage := &model.MqMessage{
		Topic:          config.GetMqProjectTrendsTopicConfig().Topic,
		Keys:           mqKeys,
		Body:           message,
		DelayTimeLevel: 3,
	}

	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, pushType, orgId)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}
