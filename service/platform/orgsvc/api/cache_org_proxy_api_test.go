package api

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/service"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetGreeter_IssueAssignRank(t *testing.T) {

	convey.Convey("Test 获取组织信息单元测试", t, test.StartUp(func(ctx context.Context) {

		_, err := service.GetCurrentUser(ctx)
		if err != nil {
			log.Error(err)
			return
		}
		//暂时用北极星的组织id
		info := getGreeter.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{
			SourceChannel: consts.AppSourceChannelDingTalk,
			OrgId:         1,
		})

		convey.ShouldNotBeNil(info.Err)
	}))
}
