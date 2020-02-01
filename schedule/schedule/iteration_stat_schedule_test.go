package schedule

import (
	"context"
	"github.com/galaxy-book/common/core/util/tests"
	"github.com/galaxy-book/polaris-backend/schedule/test"
	"github.com/Jeffail/tunny"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStatisticIterationBurnDownChart(t *testing.T) {

	convey.Convey("Test config", t, test.StartUp(func(ctx context.Context) {

		//定时任务通用协程池
		pool := tunny.NewFunc(5, func(payload interface{}) interface{} {
			fn := payload.(func() error)
			return fn()
		})
		defer pool.Close()

		convey.Convey("测试燃尽图统计", t, tests.StartUp(func() {
			convey.Convey("测试燃尽图统计", func() {
				StatisticIterationBurnDownChart(*pool)
			})
		}))
	}))
}
