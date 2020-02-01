package domain

import (
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetBaseUserOutInfoByUserIds(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		t.Log(GetBaseUserOutInfoByUserIds("fs", 1082, []int64{1167,1168,1166,1862,1871}))
	}))
}

func TestGetBaseUserInfoBatch(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		t.Log(GetBaseUserInfoBatch("", 1003, []int64{1006,1007,1008,1012,1013}))
		t.Log(GetBaseUserInfoBatch("", 1003, []int64{}))
	}))
}