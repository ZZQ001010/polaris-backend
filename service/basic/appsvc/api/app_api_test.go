package api

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/times"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/appvo"
	"github.com/galaxy-book/polaris-backend/facade/appfacade"
	services "github.com/galaxy-book/polaris-backend/service/basic/appsvc/service"
	"github.com/galaxy-book/polaris-backend/service/basic/appsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

func TestGetAppInfoByActive(t *testing.T) {
	convey.Convey("Test 行业列表", t, test.StartUp(func(ctx context.Context) {

		mysqlJson, _ := json.ToJson(config.GetMysqlConfig())
		redisJson, _ := json.ToJson(config.GetRedisConfig())

		fmt.Println("数据库配置:", mysqlJson)
		fmt.Println("redis配置", redisJson)

		userId := int64(1)
		orgId := int64(1)

		now := times.GetNowMillisecond()
		snow := "t-" + strconv.FormatInt(now, 10)
		reqVo := appvo.CreateAppInfoReqVo{
			OrgId:  orgId,
			UserId: userId,
			CreateAppInfo: vo.CreateAppInfoReq{
				Name:        "name-" + snow,
				Code:        snow,
				Secret1:     "s1-" + snow,
				Secret2:     "s2-" + snow,
				Owner:       "owner-" + snow,
				CheckStatus: consts.AppCheckStatusSuccess,
				Status:      consts.AppStatusEnable,
				Creator:     userId,
				CreateTime:  types.NowTime(),
				Updator:     userId,
				UpdateTime:  types.NowTime(),
				Version:     1,
				IsDelete:    consts.AppIsNoDelete,
			},
		}
		fmt.Println(json.ToJson(reqVo))
		resp, err := services.CreateAppInfo(reqVo)
		convey.ShouldBeNil(err)
		fmt.Println(resp.ID)

		req := appvo.AppInfoReqVo{
			AppCode: snow,
		}

		resp2 := appfacade.GetAppInfoByActiveV1(req)

		convey.ShouldEqual(resp2.AppInfo.Code, reqVo.CreateAppInfo.Code)

	}))
}
