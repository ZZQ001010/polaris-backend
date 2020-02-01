package domain

import (
	"github.com/galaxy-book/common/core/util/tests"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUpdateIssueRelationSingle(t *testing.T) {

	convey.Convey("更新任务关联测试", t, tests.StartUp(func() {
		convey.Convey("更新任务关联测试", func() {
			issueBo, err1 := GetIssueBo(17, 1011)
			if err1 != nil {
				return
			}

			_, err3 := UpdateIssueRelationSingle(1, *issueBo, consts.IssueRelationTypeParticipant, 1083)
			t.Log(err3)

		})
	}))

}
