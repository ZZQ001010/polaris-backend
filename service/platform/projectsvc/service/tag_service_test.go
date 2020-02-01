package service

import (
	"context"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTagList(t *testing.T) {
	convey.Convey("tag", t, test.StartUp(func(ctx context.Context) {
		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)
		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1002), CorpId: "1", OrgId: 1001}

		}

		log.Info("缓存用户信息" + cacheUserInfoJson)
		t.Log(CreateTag(cacheUserInfo.OrgId, cacheUserInfo.UserId, vo.CreateTagReq{
			Name:    "中华",
			BgStyle: "#33333",
			ProjectID:1001,
		}))

		t.Log(CreateTag(cacheUserInfo.OrgId, cacheUserInfo.UserId, vo.CreateTagReq{
			Name:    "中华和",
			BgStyle: "#E0E6E8",
			ProjectID:1001,
		}))

		res, err := TagList(cacheUserInfo.OrgId, 0, 0, vo.TagListReq{ProjectID:1})
		t.Log(res.Total, err)
	}))
}

func TestUpdateTag(t *testing.T) {
	convey.Convey("tag", t, test.StartUp(func(ctx context.Context) {

		//t.Log(HotTagList(1001, 1001))
		//name := "还会f"
		//bgStyle := "#FFE3C4"
		//t.Log(UpdateTag(1001, 1002, vo.UpdateTagReq{ID:10048, Name:&name, BgStyle:&bgStyle}))
		//t.Log(DeleteTag(1001, 1002, vo.DeleteTagReq{ProjectID:1001, Ids:[]int64{10048}}))
		//t.Log(HotTagList(1001, 1001))
		t.Log(CreateTag(1001, 1002, vo.CreateTagReq{
			Name:    "中华和",
			BgStyle: "#E0E6E8",
			ProjectID:1001,
		}))
	}))
}