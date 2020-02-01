package schedule

import (
	"context"
	"github.com/galaxy-book/polaris-backend/schedule/test"
	"github.com/Jeffail/tunny"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestStatisticProjectIssueBurnDownChart(t *testing.T) {
	//定时任务通用协程池
	pool := tunny.NewFunc(5, func(payload interface{}) interface{} {
		fn := payload.(func() error)
		return fn()
	})
	defer pool.Close()

	convey.Convey("测试燃尽图统计", t, test.StartUp(func(ctx context.Context) {
		convey.Convey("测试燃尽图统计", func() {
			StatisticProjectIssueBurnDownChart(*pool)
		})
	}))

	time.Sleep(time.Duration(2) * time.Second)
}
