package notice

import (
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/polaris-team/dingtalk-sdk-golang/json"
	"testing"
)

func TestPushIssue(t *testing.T) {

	config.LoadEnvConfig("F:\\workspace-golang-polaris\\polaris-backend\\service\\platform\\projectsvc\\config", "application", "test")

	data := "{\"PushType\":2,\"OrgId\":1001,\"OperatorId\":1010,\"IssueId\":1190,\"ParentIssueId\":0,\"ProjectId\":1001,\"IssueTitle\":\"测试哦哦哦哦哦\",\"IssueStatusId\":7,\"BeforeOwner\":1010,\"AfterOwner\":0,\"BeforeChangeFollowers\":[],\"AfterChangeFollowers\":[],\"BeforeChangeParticipants\":[1014],\"AfterChangeParticipants\":[],\"OperateObjProperty\":\"status\",\"NewValue\":\"{\\\"status\\\":15}\",\"OldValue\":\"{\\\"status\\\":7}\"}"

	bo := &bo.IssueNoticeBo{}

	json.FromJson(data, bo)

	PushIssue(*bo)

}
