package domain

import (
	"github.com/galaxy-book/common/core/util/tests"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAppendIterationStat(t *testing.T) {

	convey.Convey("Test AppendIterationStat", t, tests.StartUp(func() {
		convey.Convey("Test AppendIterationStat", func() {
			iterationBo := bo.IterationBo{
				Id: 1010,
				OrgId: 17,
				ProjectId: 1139,
			}
			convey.So(AppendIterationStat(iterationBo, consts.BlankDate), convey.ShouldEqual, nil)
		})
	}))


}