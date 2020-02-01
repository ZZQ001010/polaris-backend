package service

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/websitevo"
	"github.com/galaxy-book/polaris-backend/service/platform/websitesvc/test"
	"github.com/magiconair/properties/assert"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRegisterWebSiteContact(t *testing.T) {

	convey.Convey("Test RegisterWebSiteContact", t,test.StartUp(func(ctx context.Context) {
		convey.Convey("test", func() {

			void, err := RegisterWebSiteContact(websitevo.RegisterWebSiteContactReqVo{
				Input: vo.RegisterWebSiteContactReq{
					Sex: 1,
					Name: "Nico",
					ContactInfo: "1111111112",
				},
			})
			t.Log(void, err)
			assert.Equal(t, err, nil)
		})
	}))

}