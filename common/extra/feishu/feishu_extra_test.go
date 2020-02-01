package feishu

import (
	"context"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetDeptUserInfosByDeptIds(t *testing.T) {

	convey.Convey("TestGetDeptUserInfosByDeptIds", t, test.StartUp(func(ctx context.Context) {

		list, err := GetDeptUserInfosByDeptIds("2ed263bf32cf1651", []string{"0"})
		t.Log(err)
		for _, u := range list{
			t.Log(json.ToJsonIgnoreError(u))
		}
	}))

}

func TestGetScopeOpenIds2(t *testing.T) {

	convey.Convey("TestGetScopeOpenIds2", t, test.StartUp(func(ctx context.Context) {

		list, err := GetScopeOpenIds("2ed263bf32cf1651")
		t.Log(err)
		for _, u := range list{
			t.Log(json.ToJsonIgnoreError(u))
		}
	}))

}

func TestGetDeptList(t *testing.T) {
	convey.Convey("TestGetDeptList", t, test.StartUp(func(ctx context.Context) {

		list, err := GetDeptList("2ed263bf32cf1651")
		t.Log(err)
		for _, u := range list{
			t.Log(json.ToJsonIgnoreError(u))
		}
	}))
}

func TestGetScopeDeps(t *testing.T) {
	convey.Convey("TestGetScopeDeps", t, test.StartUp(func(ctx context.Context) {

		list, err := GetScopeDeps("2ed263bf32cf1651")
		t.Log(err)
		for _, u := range list{
			t.Log(json.ToJsonIgnoreError(u))
		}
	}))
}