package service

import (
	"context"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/test"
	"github.com/smartystreets/goconvey/convey"
	"gotest.tools/assert"
	"testing"
)

func TestPermissionOperationList(t *testing.T) {
	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		var projectId int64 = 10101
		data, err := PermissionOperationList(10101, 10231, 10201, &projectId)
		t.Log(json.ToJsonIgnoreError(data))
		assert.Equal(t, err, nil)
	}))
}

func TestUpdateRolePermissionOperation(t *testing.T) {
	convey.Convey("Test UpdateRolePermissionOperation", t, test.StartUp(func(ctx context.Context) {
		updArr := []*vo.EveryPermission{}
		upd := vo.EveryPermission{PermissionID:15, OperationIds:[]int64{51}}
		data, err := UpdateRolePermissionOperation(10101, 10201, vo.UpdateRolePermissionOperationReq{
			RoleID:10231,
			UpdatePermissions: append(updArr, &upd),
		})
		t.Log(json.ToJsonIgnoreError(data))
		assert.Equal(t, err, nil)
	}))
}

func TestGetPersonalPermissionInfo(t *testing.T) {
	convey.Convey("Test UpdateRolePermissionOperation", t, test.StartUp(func(ctx context.Context) {
		var projectId int64 = 1119
		var issueId int64 = 11561
		//t.Log(GetPersonalPermissionInfo(1001, 1002, &projectId, nil, ""))
		t.Log(GetPersonalPermissionInfo(1070, 1130, &projectId, &issueId, ""))
	}))
}