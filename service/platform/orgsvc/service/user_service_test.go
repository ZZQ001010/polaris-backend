package service

import (
	"context"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUserConfigInfo(t *testing.T) {
	convey.Convey("Test sql", t, test.StartUp(func(ctx context.Context) {

		config, err := UserConfigInfo(1001, 1001)

		convey.ShouldBeNil(err)
		convey.ShouldNotBeNil(config)
	}))
}

func TestUpdateUserInfo(t *testing.T) {

	convey.Convey("Test UpdateUserInfo", t, test.StartUpWithUserInfo(func(userId, orgId int64) {

		updateFields := []string{"name", "avatar", "sex", "birthday"}

		avatar := "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1576325272960&di=64d46c5b9578915bbc91ddf4a282d917&imgtype=0&src=http%3A%2F%2Fimg1.xiazaizhijia.com%2Fwalls%2F20150910%2Fmiddle_eea71d559bf1341.jpg"
		name := "更改后的名字"
		sex := 1
		nowTime := types.NowTime()

		input := vo.UpdateUserInfoReq{
			Name:         &name,
			Avatar:       &avatar,
			UpdateFields: updateFields,
			Sex:          &sex,
			Birthday:     &nowTime,
		}

		resp, err := UpdateUserInfo(0, userId, input)

		convey.ShouldBeNil(err)
		convey.ShouldNotBeNil(resp)
	}))

}

func TestVerifyOrgUsers(t *testing.T) {
	convey.Convey("Test UpdateUserInfo", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		t.Log(VerifyOrgUsers(1001, []int64{1,2}))
	}))

}