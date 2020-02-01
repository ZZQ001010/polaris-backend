package api

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/extra/gin/util"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSwitchUserOrganization(t *testing.T) {
	convey.Convey("Test 切换用户组织", t, test.StartUp(func(ctx context.Context) {
		token, _ := util.GetCtxToken(ctx)

		orgfacade.SendSMSLoginCode(orgvo.SendSMSLoginCodeReqVo{
			PhoneNumber: "17621142248",
		})

		mysqlJson, _ := json.ToJson(config.GetMysqlConfig())
		redisJson, _ := json.ToJson(config.GetRedisConfig())

		fmt.Println("数据库配置:", mysqlJson)
		fmt.Println("redis配置", redisJson)

		user := orgfacade.GetCurrentUser(ctx)

		reqVo := orgvo.SwitchUserOrganizationReqVo{
			OrgId:  user.CacheInfo.OrgId,
			UserId: user.CacheInfo.UserId,
			Token:  token,
		}

		resp := orgfacade.SwitchUserOrganization(reqVo)

		convey.ShouldBeFalse(resp.Failure())
	}))
}

func TestUserOrganizationList(t *testing.T) {
	convey.Convey("Test 测试用户组织", t, test.StartUp(func(ctx context.Context) {

		mysqlJson, _ := json.ToJson(config.GetMysqlConfig())
		redisJson, _ := json.ToJson(config.GetRedisConfig())

		fmt.Println("数据库配置:", mysqlJson)
		fmt.Println("redis配置", redisJson)

		user := orgfacade.GetCurrentUser(ctx)

		reqVo := orgvo.UserOrganizationListReqVo{
			UserId: user.CacheInfo.UserId,
		}

		resp := orgfacade.UserOrganizationList(reqVo)

		convey.ShouldBeFalse(resp.Failure())
	}))
}

//func TestOrganizationSettingsOwnerFail(t *testing.T) {
//	convey.Convey("Test 测试组织设置", t, test.StartUp(func(ctx context.Context) {
//
//		mysqlJson, _ := json.ToJson(config.GetMysqlConfig())
//		redisJson, _ := json.ToJson(config.GetRedisConfig())
//
//		fmt.Println("数据库配置:", mysqlJson)
//		fmt.Println("redis配置", redisJson)
//
//		user := orgfacade.GetCurrentUser(ctx)
//
//		website := "www.baidu.com"
//
//		reqVo := orgvo.OrganizationSettingsReqVo{
//			UserId: user.CacheInfo.UserId,
//			Input: vo.OrganizationSettingsReq{
//				OrgID:   user.CacheInfo.OrgId,
//				OrgName: "alan的测试组织名",
//				WebSite: &website},
//		}
//
//		resp := orgfacade.OrganizationSettings(reqVo)
//
//		convey.ShouldBeTrue(resp.Failure())
//	}))
//}
//
//func TestOrganizationSettingsOwner(t *testing.T) {
//	convey.Convey("Test 测试组织设置", t, test.StartUp(func(ctx context.Context) {
//
//		mysqlJson, _ := json.ToJson(config.GetMysqlConfig())
//		redisJson, _ := json.ToJson(config.GetRedisConfig())
//
//		fmt.Println("数据库配置:", mysqlJson)
//		fmt.Println("redis配置", redisJson)
//
//		user := orgfacade.GetCurrentUser(ctx)
//
//		website := "www.baidu.com"
//
//		reqVo := orgvo.OrganizationSettingsReqVo{
//			UserId: 1029,
//			Input: vo.OrganizationSettingsReq{
//				OrgID:   user.CacheInfo.OrgId,
//				OrgName: "alan的测试组织名",
//				WebSite: &website},
//		}
//
//		resp := orgfacade.OrganizationSettings(reqVo)
//
//		convey.ShouldBeFalse(resp.Failure())
//	}))
//}
