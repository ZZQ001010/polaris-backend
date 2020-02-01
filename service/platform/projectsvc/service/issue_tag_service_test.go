package service

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"gotest.tools/assert"
	"testing"
)

func TestCreateIssueRelationTags(t *testing.T) {
	convey.Convey("TestIssueAndProjectCountStat", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		err := CreateIssueRelationTags(projectvo.CreateIssueRelationTagsReqVo{
			OrgId: orgId,
			UserId: userId,
			Input: vo.UpdateIssueTagsReq{
				ID: test.IssueId1083,
				AddTags: []*vo.IssueTagReqInfo{
					{
						ID: 1001,
						Name: "测试",
					},
				},
			},
		})
		assert.Equal(t, err, nil)
	}))


}

