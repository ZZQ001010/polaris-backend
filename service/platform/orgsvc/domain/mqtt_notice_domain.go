package domain

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/extra/mqtt"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
)

func PushAddOrgMemberNotice(orgId, depId int64, memberIds []int64) {
	globalRefreshList := make([]bo.MQTTGlobalRefresh, 0)

	baseUserInfos, err := GetBaseUserInfoBatch("", orgId, memberIds)
	if err != nil{
		log.Error(err)
		return
	}

	for _, baseUserInfo := range baseUserInfos{
		globalRefreshList = append(globalRefreshList, bo.MQTTGlobalRefresh{
			ObjectId: baseUserInfo.UserId,
			ObjectValue: bo.BaseUserInfoExtBo{
				BaseUserInfoBo: baseUserInfo,
				DepartmentId: depId,
			},
		})
	}

	//推送refresh
	err = mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId: orgId,
		Action: consts.MQTTDataRefreshActionAdd,
		Type: consts.MQTTDataRefreshTypeMember,
		GlobalRefresh: globalRefreshList,
	})
	if err != nil{
		log.Error(err)
	}
}

func PushRemoveOrgMemberNotice(orgId int64, memberIds []int64) {
	globalRefreshList := make([]bo.MQTTGlobalRefresh, 0)

	baseUserInfos, err := GetBaseUserInfoBatch("", orgId, memberIds)
	if err != nil{
		log.Error(err)
		return
	}

	for _, baseUserInfo := range baseUserInfos{
		globalRefreshList = append(globalRefreshList, bo.MQTTGlobalRefresh{
			ObjectId: baseUserInfo.UserId,
			ObjectValue: baseUserInfo,
		})
	}

	//推送refresh
	err = mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId: orgId,
		Action: consts.MQTTDataRefreshActionDel,
		Type: consts.MQTTDataRefreshTypeMember,
		GlobalRefresh: globalRefreshList,
	})
	if err != nil{
		log.Error(err)
	}
}