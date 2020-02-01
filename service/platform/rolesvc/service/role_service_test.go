package service

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCreateRole(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		t.Log(CreateRole(10113, 10201, vo.CreateRoleReq{
			RoleGroupType: 1,
			Name:          "测试角色",
		}))
	}))
}
