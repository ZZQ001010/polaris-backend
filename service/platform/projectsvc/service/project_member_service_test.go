package service

import (
	"context"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"gotest.tools/assert"
	"testing"
)

func TestProjectUserList(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		data, err := ProjectUserList(1001, 1, 10, vo.ProjectUserListReq{ProjectID:1001})
		t.Log(json.ToJsonIgnoreError(data))
		assert.Equal(t, err, nil)

		t.Log(AddProjectMember(1001, 1001, vo.RemoveProjectMemberReq{ProjectID:1001, MemberIds:[]int64{1002}}))
		t.Log(RemoveProjectMember(1001, 1001, vo.RemoveProjectMemberReq{ProjectID:1001, MemberIds:[]int64{1002}}))
	}))
}

func TestAddProjectMember(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		//t.Log(AddProjectMember(10101, 10201, vo.RemoveProjectMemberReq{ProjectID:10101, MemberIds:[]int64{10201, 10203}}))
		t.Log(AddProjectMember(1323, 1608, vo.RemoveProjectMemberReq{ProjectID:1521, MemberIds:[]int64{1464, 1609}}))
	}))
}