package service

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestIssueAndProjectCountStat(t *testing.T) {

	convey.Convey("TestIssueAndProjectCountStat", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		resp, err := IssueAndProjectCountStat(projectvo.IssueAndProjectCountStatReqVo{
			OrgId: orgId,
			UserId: userId,
		})

		log.Error(err)
		log.Error(json.ToJsonIgnoreError(resp))
		assert.Equal(t, err, nil)
		assert.NotEqual(t, resp, nil)
		assert.Equal(t, resp.ProjectNotCompletedCount, int64(1))
		assert.Equal(t, resp.IssueNotCompletedCount, int64(1))
	}))

}

func TestIssueDailyPersonalWorkCompletionStat(t *testing.T) {
	convey.Convey("TestIssueDailyPersonalWorkCompletionStat", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		resp, err := IssueDailyPersonalWorkCompletionStat(projectvo.IssueDailyPersonalWorkCompletionStatReqVo{
			OrgId: orgId,
			UserId: userId,
		})

		log.Error(err)
		log.Error(json.ToJsonIgnoreError(resp))
		assert.Equal(t, err, nil)
		assert.NotEqual(t, resp, nil)
	}))

}