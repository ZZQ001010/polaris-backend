package service

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
)

func TestImportIssues(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		config.LoadEnvConfig("F:\\polaris-backend-clone\\config", "application.common", "local")

		testMysqlJson, _ := json.ToJson(config.GetMysqlConfig())
		fmt.Println("unittest Mysql配置json:", testMysqlJson)

		testRedisJson, _ := json.ToJson(config.GetRedisConfig())
		fmt.Println("unittest redis配置json:", testRedisJson)

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		t.Log(ImportIssues(cacheUserInfo.OrgId, cacheUserInfo.UserId, vo.ImportIssuesReq{
			ProjectID: 1007,
			URL:       "",
		}))
	}))

}

func TestExportIssueTemplate(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		//config.LoadEnvConfig("F:\\polaris-backend-clone\\config", "application.common", "local")
		testMysqlJson, _ := json.ToJson(config.GetMysqlConfig())
		fmt.Println("unittest Mysql配置json:", testMysqlJson)

		testRedisJson, _ := json.ToJson(config.GetRedisConfig())
		fmt.Println("unittest redis配置json:", testRedisJson)

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1006), CorpId: "1", OrgId: 1001}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
		t.Log(cacheUserInfo.OrgId)

		t.Log(ExportIssueTemplate(cacheUserInfo.OrgId, 1))
	}))
}

func Test_BatchImport_Sync(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		for i := 0; i < 1; i++ {
			go func() {
				defer wg.Add(-1)
				issueDetailId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueDetail)
				t.Log(issueDetailId, err)
			}()
		}
		wg.Wait()

		t.Log("ok")
	}))

}
