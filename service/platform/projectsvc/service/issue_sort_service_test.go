package service

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/test"
	"github.com/magiconair/properties/assert"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUpdateIssueSort(t *testing.T) {

	convey.Convey("TestUpdateIssueSort", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		afterId := test.IssueId100
		_, err := UpdateIssueSort(projectvo.UpdateIssueSortReqVo{
			UserId: userId,
			OrgId: orgId,
			Input: vo.UpdateIssueSortReq{
				ID: test.IssueId1083,
				AfterID: &afterId,
			},
		})
		assert.Equal(t, err, nil)
		log.Error(err)
	}))

}