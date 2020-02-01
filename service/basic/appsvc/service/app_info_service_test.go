package services

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/times"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/appvo"
	"github.com/galaxy-book/polaris-backend/service/basic/appsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

func TestCreateGetUpdateAppInfo(t *testing.T) {
	convey.Convey("Test 创建应用", t, test.StartUp(func(ctx context.Context) {

		mysqlJson, _ := json.ToJson(config.GetMysqlConfig())
		redisJson, _ := json.ToJson(config.GetRedisConfig())

		fmt.Println("数据库配置:", mysqlJson)
		fmt.Println("redis配置", redisJson)

		userId := int64(1)
		orgId := int64(1)

		now := times.GetNowMillisecond()
		snow := strconv.FormatInt(now, 10)
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
		resp, err := CreateAppInfo(reqVo)
		convey.ShouldBeNil(err)
		fmt.Println(resp.ID)

		// GET
		appInfo, err := GetAppInfoByActive(reqVo.CreateAppInfo.Code)
		convey.ShouldBeNil(err)
		fmt.Println(json.ToJson(appInfo))
		convey.ShouldEqual(reqVo.CreateAppInfo.Name, appInfo.Name)

		// Update
		updAppInfo := &vo.UpdateAppInfoReq{}
		copyErr := copyer.Copy(appInfo, updAppInfo)
		if copyErr != nil {
			convey.ShouldBeNil(copyErr)
		}

		updAppInfo.Name = snow

		updVo := appvo.UpdateAppInfoReqVo{
			Input:  *updAppInfo,
			UserId: userId,
			OrgId:  orgId,
		}
		result, err := UpdateAppInfo(updVo)

		convey.ShouldBeNil(err)
		convey.ShouldEqual(updAppInfo.ID, result.ID)

		appInfo2, err := GetAppInfoByActive(updVo.Input.Code)
		convey.ShouldEqual(appInfo2.Name, updAppInfo.Name)

		updVo.Input.Status = consts.AppStatusDisabled
		result, err = UpdateAppInfo(updVo)

		appInfo4, err := GetAppInfoByActive(updVo.Input.Code)
		fmt.Println(appInfo4, err)
		convey.ShouldEqual(err.Code(), errs.TargetNotExist.Code())

	}))
}
