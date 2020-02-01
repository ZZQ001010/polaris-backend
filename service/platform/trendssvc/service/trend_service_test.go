package service

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/util/slice"
	bo2 "github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/trendssvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTrendList(t *testing.T) {
	//type Exte struct {
	//	ObjName string
	//}
	//a := "{\"issueType\":\"T\",\"objName\":\"测试任务\",\"ccc\":\"aa\"}"
	//ex := map[string]interface{}{
	//	"issueType":"T",
	//	"objName":"测试",
	//}
	////ex := Ext{
	////	ObjName:"ceshi",
	////}
	//b := &Exte{}
	//c := &Exte{}
	//d := &map[string]interface{}{}
	//fmt.Println(json.FromJson(json.ToJsonIgnoreError(ex), b), b)
	//fmt.Println(json.FromJson(a, c), c)
	//fmt.Println(json.FromJson(a, d), d, (*d)["objName"])
	old := "{\"participant\":[1081,1064,1069,1077,1078,1079,1080], \"name\":\"ss\"}"
	new := "{\"participant\":[1080,1069,1079,1077,1078], \"name\":\"ssb\"}"
	old1 := &map[string]interface{}{}
	new1 := &map[string]interface{}{}
	json.FromJson(old, old1)
	json.FromJson(new, new1)
	fmt.Println(slice.JsonCompare(*old1, *new1))
}

func TestCreateTrends(t *testing.T) {
	ext := "{111}"
	bo := &bo2.ProjectObjectTypeBo{}
	_ = json.FromJson(ext, bo)
	fmt.Println(bo)
}

func TestTrendList2(t *testing.T) {
	convey.Convey("UnreadNoticeCount", t, test.StartUp(func(ctx context.Context) {
		operId := int64(1458)
		OrderType := int(2)
		t.Log(TrendList(1220, 1458, &vo.TrendReq{
			LastTrendID: nil,
			ObjType:     nil,
			ObjID:       nil,
			OperID:      &operId,
			StartTime:   nil,
			EndTime:     nil,
			Type:        nil,
			Page:        nil,
			Size:        nil,
			OrderType:   &OrderType,
		}))
	}))
	
}