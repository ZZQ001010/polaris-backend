package service

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/test"
	"github.com/magiconair/properties/assert"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUpdateOrgMemberStatus(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		_, err := UpdateOrgMemberStatus(orgvo.UpdateOrgMemberStatusReq{
			UserId: 1004,
			OrgId:  orgId,
			Input: vo.UpdateOrgMemberStatusReq{
				MemberIds: []int64{1019},
				Status:    2,
			},
		})
		assert.Equal(t, err, nil)
	}))
}

func TestUpdateOrgMemberCheckStatus(t *testing.T) {
	convey.Convey("TestUpdateOrgMemberStatus", t, test.StartUpWithUserInfo(func(userId, orgId int64) {
		_, err := UpdateOrgMemberCheckStatus(orgvo.UpdateOrgMemberCheckStatusReq{
			UserId: 1004,
			OrgId:  orgId,
			Input: vo.UpdateOrgMemberCheckStatusReq{
				MemberIds:   []int64{1019},
				CheckStatus: 2,
			},
		})
		assert.Equal(t, err, nil)
	}))
}

func TestOrgUserList(t *testing.T) {
	type TestStruct struct {
		Id   int64
		Name string
	}
	list := &[]TestStruct{
		{
			Id:   1,
			Name: "hello",
		},
		{
			Id:   1,
			Name: "world",
		},
		{
			Id:   3,
			Name: "haha",
		},
	}

	cacheMap := maps.NewMap("Id", list)
	t.Log(json.ToJsonIgnoreError(cacheMap))
	t.Log(cacheMap)
	t.Log(cacheMap[int64(1)])
}

func TestGetOrgUserInfoListBySourceChannel(t *testing.T) {
	convey.Convey("TestGetOrgUserInfoListBySourceChannel", t, test.StartUpWithUserInfo(func(userId, orgId int64) {

		resp, err := GetOrgUserInfoListBySourceChannel(orgvo.GetOrgUserInfoListBySourceChannelReq{
			SourceChannel: consts.AppSourceChannelFeiShu,
			Page:          -1,
			Size:          -1,
		})
		t.Log(err)
		assert.Equal(t, err, nil)
		t.Log(json.ToJsonIgnoreError(resp))

		for _, userInfo := range resp.List {
			t.Log(json.ToJsonIgnoreError(userInfo))
		}

		t.Log("===================================")

		resp, err = GetOrgUserInfoListBySourceChannel(orgvo.GetOrgUserInfoListBySourceChannelReq{
			SourceChannel: consts.AppSourceChannelFeiShu,
			Page:          2,
			Size:          10,
		})
		t.Log(err)
		assert.Equal(t, err, nil)
		t.Log(json.ToJsonIgnoreError(resp))

		for _, userInfo := range resp.List {
			t.Log(json.ToJsonIgnoreError(userInfo))
		}

		t.Log("===================================")

		resp, err = GetOrgUserInfoListBySourceChannel(orgvo.GetOrgUserInfoListBySourceChannelReq{
			SourceChannel: consts.AppSourceChannelFeiShu,
			Page:          -1,
			Size:          -1,
		})
		t.Log(err)
		assert.Equal(t, err, nil)
		t.Log(json.ToJsonIgnoreError(resp))

		for _, userInfo := range resp.List {
			t.Log(json.ToJsonIgnoreError(userInfo))
		}
	}))
}
