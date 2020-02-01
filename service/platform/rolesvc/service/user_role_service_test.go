package service

import (
	"context"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/polaris-backend/common/model/vo/rolevo"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetOrgRoleUser(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		t.Log(GetOrgRoleUser(1001, 0))
	}))
}

func TestMap(t *testing.T) {
	a := []map[string]string{}
	a = append(a, map[string]string{
		"a":    "a",
		"name": "jj",
	})
	a = append(a, map[string]string{
		"a":    "b",
		"name": "dd",
	})
	t.Log(maps.NewMap("a", a))
}

func TestUpdateUserOrgRole(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {

		req := rolevo.UpdateUserOrgRoleReqVo{
			OrgId:         10101,
			CurrentUserId: 10201,
			UserId:        10202,
			RoleId:        10280,
		}

		t.Log(UpdateUserOrgRole(req))
	}))
}

func TestGetProjectRoleList(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {

		data, err := GetProjectRoleList(10101, 10101)
		t.Log(json.ToJsonIgnoreError(data))
		t.Log(len(data))
		t.Log(err)
	}))
}