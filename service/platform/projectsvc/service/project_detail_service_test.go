package service

import (
	"fmt"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	bo2 "github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"testing"
	"time"
)

func TestUpdateProjectDetail(t *testing.T) {

	bo := &bo2.ProjectDetailBo{}
	bo.UpdateTime = time.Now()
	po := &po.PpmProProjectDetail{}
	mayBlank, _ := time.Parse(consts.AppTimeFormat, "")
	fmt.Println(mayBlank.Format(consts.AppTimeFormat))
	copyer.Copy(bo, po)

	t.Log(po)

}
