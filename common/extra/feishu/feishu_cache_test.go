package feishu

import (
	"context"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/test"
	"github.com/magiconair/properties/assert"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetAppAccessToken(t *testing.T) {

	convey.Convey("TestGetAppAccessToken", t, test.StartUp(func(ctx context.Context) {

		cacheBo, err := GetAppAccessToken()
		assert.Equal(t, err, nil)
		t.Log(json.ToJsonIgnoreError(cacheBo))
	}))

}

func TestGetTenantAccessToken(t *testing.T) {
	convey.Convey("TestGetTenantAccessToken", t, test.StartUp(func(ctx context.Context) {

		cacheBo, err := GetTenantAccessToken("2ed263bf32cf1651")
		assert.Equal(t, err, nil)
		t.Log(json.ToJsonIgnoreError(cacheBo))
	}))
}

func TestGetScopeOpenIds(t *testing.T) {
	convey.Convey("TestGetScopeOpenIds", t, test.StartUp(func(ctx context.Context) {

		cacheBo, err := GetScopeOpenIds("2ed263bf32cf1651")
		assert.Equal(t, err, nil)
		t.Log(json.ToJsonIgnoreError(cacheBo))
	}))
}