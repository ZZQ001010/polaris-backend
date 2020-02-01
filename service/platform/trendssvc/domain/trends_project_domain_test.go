package domain

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/trendssvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAddProjectTrends(t *testing.T) {
	convey.Convey("Test ArchiveProject", t, test.StartUp(func(ctx context.Context) {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = "二号测试应用aa"
		AddProjectTrends(bo.ProjectTrendsBo{
			PushType:              consts.PushTypeUpdateProjectMembers,
			OrgId:                 10101,
			ProjectId:             10113,
			OperatorId:            10201,
			BeforeChangeMembers:   []int64{},
			AfterChangeMembers:    []int64{},
			BeforeChangeFollowers: []int64{},
			AfterChangeFollowers:  []int64{10202},
			Ext:                   ext,
			SourceChannel:         "",
		})
	}))
}
