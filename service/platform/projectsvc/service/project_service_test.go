package service

import (
	"context"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetProjectRelation(t *testing.T) {
	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
		t.Log(GetProjectRelation(1001, []int64{1, 2, 3}))
	}))
}
func TestArchiveProject(t *testing.T) {
	convey.Convey("Test ArchiveProject", t, test.StartUp(func(ctx context.Context) {
		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1001), CorpId: "1", OrgId: 1001}

		}

		t.Log(ArchiveProject(cacheUserInfo.OrgId, cacheUserInfo.UserId, 1))
		t.Log(CancelArchivedProject(cacheUserInfo.OrgId, cacheUserInfo.UserId, 1))
	}))
}

func TestOrgProjectMembers(t *testing.T) {

	convey.Convey("Test OrgProjectMembers", t, test.StartUp(func(ctx context.Context) {

		resp, err := OrgProjectMembers(projectvo.OrgProjectMemberReqVo{
			ProjectId: 1001,
			OrgId:     1001,
			UserId:    1001,
		})

		convey.ShouldBeNil(err)
		convey.ShouldNotBeNil(resp)
	}))
}

func TestProjects(t *testing.T) {
	convey.Convey("Test ArchiveProject", t, test.StartUp(func(ctx context.Context) {
		t.Log(Projects(projectvo.ProjectsRepVo{OrgId:10101, UserId:10201,Page:1,Size:10}))
	}))
}

func TestUpdateProject(t *testing.T) {
	convey.Convey("Test ArchiveProject", t, test.StartUp(func(ctx context.Context) {
		res, err := UpdateProject(projectvo.UpdateProjectReqVo{
			OrgId:10101,
			UserId:10201,
			Input: vo.UpdateProjectReq{
				ID:10113,
				FollowerIds:[]int64{10202},
				UpdateFields:[]string{"followerIds"},
			},
		})
		t.Log(res)
		t.Log(err)
	}))
}

func TestStarProject(t *testing.T) {
	convey.Convey("Test StarProject", t, test.StartUp(func(ctx context.Context) {
		t.Log(StarProject(projectvo.ProjectIdReqVo{ProjectId:1006, SourceChannel:"",UserId:1007,OrgId:1003}))
	}))
}