package domain

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConstlog(t *testing.T) {

	convey.Convey("Test log", t, func() {
		convey.Convey("log", func() {
			convey.So(log, convey.ShouldNotBeNil)
		})
	})
}
