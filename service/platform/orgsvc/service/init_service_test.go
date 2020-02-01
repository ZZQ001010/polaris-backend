package service

import (
	"context"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"upper.io/db.v3/lib/sqlbuilder"
)

func TestLarkInit(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo := bo.CacheUserInfoBo{OutUserId: "213213", SourceChannel: "lark_test", UserId: int64(1001), CorpId: "1", OrgId: 222}

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		cache.Set(consts.CacheUserToken+"abc", cacheUserInfoJson)

		mysql.TransX(func(tx sqlbuilder.Tx) error {
			err := LarkInit(cacheUserInfo.OrgId, cacheUserInfo.UserId, "1", "2")
			return err
		})
	}))
}
