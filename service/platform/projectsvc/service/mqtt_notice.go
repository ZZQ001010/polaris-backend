package service

import (
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/extra/mqtt"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
)

func PushAddProjectNotice(orgId, projectId int64) {
	projectInfo, err := ProjectInfo(orgId, vo.ProjectInfoReq{
		ProjectID: projectId,
	}, "")
	if err != nil{
		log.Error(err)
		return
	}

	//推送refresh
	err = mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId: orgId,
		Action: consts.MQTTDataRefreshActionAdd,
		Type: consts.MQTTDataRefreshTypePro,
		GlobalRefresh: []bo.MQTTGlobalRefresh{
			{
				ObjectId: projectId,
				ObjectValue: projectInfo,
			},
		},
	})
	if err != nil{
		log.Error(err)
	}
}

func PushModifyProjectNotice(orgId, projectId int64){
	projectInfo, err := ProjectInfo(orgId, vo.ProjectInfoReq{
		ProjectID: projectId,
	}, "")
	if err != nil{
		log.Error(err)
		return
	}

	err = mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId: orgId,
		Action: consts.MQTTDataRefreshActionModify,
		Type: consts.MQTTDataRefreshTypePro,
		GlobalRefresh: []bo.MQTTGlobalRefresh{
			{
				ObjectId: projectId,
				ObjectValue: projectInfo,
			},
		},
	})
	if err != nil{
		log.Error(err)
	}

	//推送refresh
	err = mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId: orgId,
		ProjectId: projectId,
		Action: consts.MQTTDataRefreshActionModify,
		Type: consts.MQTTDataRefreshTypePro,
		GlobalRefresh: []bo.MQTTGlobalRefresh{
			{
				ObjectId: projectId,
				ObjectValue: projectInfo,
			},
		},
	})
	if err != nil{
		log.Error(err)
	}
}

func PushAddIssueNotice(orgId, projectId, issueId int64, currentUserId int64){
	issueInfo, err := IssueInfo(orgId, currentUserId, issueId, "")
	if err != nil{
		log.Error(err)
		return
	}

	//推送refresh
	err = mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId: orgId,
		ProjectId: projectId,
		Action: consts.MQTTDataRefreshActionAdd,
		Type: consts.MQTTDataRefreshTypeIssue,
		GlobalRefresh: []bo.MQTTGlobalRefresh{
			{
				ObjectId: currentUserId,
				ObjectValue: issueInfo,
			},
		},
	})
	if err != nil{
		log.Error(err)
	}
}

func PushModifyIssueNotice(orgId, projectId, issueId int64, currentUserId int64){
	issueInfo, err := IssueInfo(orgId, currentUserId, issueId, "")
	if err != nil{
		log.Error(err)
		return
	}

	//推送refresh
	err = mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId: orgId,
		ProjectId: projectId,
		Action: consts.MQTTDataRefreshActionModify,
		Type: consts.MQTTDataRefreshTypeIssue,
		GlobalRefresh: []bo.MQTTGlobalRefresh{
			{
				ObjectId: currentUserId,
				ObjectValue: issueInfo,
			},
		},
	})
	if err != nil{
		log.Error(err)
	}
}

func PushDelIssueNotice(orgId, projectId int64, issueIds []int64){
	globalRefreshList := make([]bo.MQTTGlobalRefresh, 0)
	for _, issueId := range issueIds{
		globalRefreshList = append(globalRefreshList, bo.MQTTGlobalRefresh{
			ObjectId: issueId,
			ObjectValue: nil,
		})
	}

	//推送refresh
	err := mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId: orgId,
		ProjectId: projectId,
		Action: consts.MQTTDataRefreshActionDel,
		Type: consts.MQTTDataRefreshTypeIssue,
		GlobalRefresh: globalRefreshList,
	})
	if err != nil{
		log.Error(err)
	}
}

func PushAddTagNotice(orgId, projectId int64, tags []bo.TagBo){
	if len(tags) == 0{
		return
	}

	refreshList := make([]bo.MQTTGlobalRefresh, len(tags))
	for i, tagBo := range tags{
		refreshList[i] = bo.MQTTGlobalRefresh{
			ObjectId: tagBo.Id,
			ObjectValue: tagBo,
		}
	}

	//推送refresh
	err := mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId: orgId,
		ProjectId: projectId,
		Action: consts.MQTTDataRefreshActionAdd,
		Type: consts.MQTTDataRefreshTypeTag,
		GlobalRefresh: refreshList,
	})
	if err != nil{
		log.Error(err)
	}
}

func PushModifyTagNotice(orgId, projectId, tagId int64){
	tagInfo, err := domain.GetTagInfo(tagId)
	if err != nil {
		log.Error(err)
		return
	}
	//推送refresh
	err = mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId: orgId,
		ProjectId: projectId,
		Action: consts.MQTTDataRefreshActionModify,
		Type: consts.MQTTDataRefreshTypeTag,
		GlobalRefresh: []bo.MQTTGlobalRefresh{
			{
				ObjectId: tagId,
				ObjectValue: tagInfo,
			},
		},
	})
	if err != nil{
		log.Error(err)
	}
}

func PushRemoveTagNotice(orgId, projectId int64, tagIds []int64){
	globalRefreshList := make([]bo.MQTTGlobalRefresh, 0)
	for _, tagId := range tagIds{
		globalRefreshList = append(globalRefreshList, bo.MQTTGlobalRefresh{
			ObjectId: tagId,
			ObjectValue: nil,
		})
	}

	//推送refresh
	err := mqtt.PushMQTTDataRefreshMsg(bo.MQTTDataRefreshNotice{
		OrgId: orgId,
		ProjectId: projectId,
		Action: consts.MQTTDataRefreshActionDel,
		Type: consts.MQTTDataRefreshTypeTag,
		GlobalRefresh: globalRefreshList,
	})
	if err != nil{
		log.Error(err)
	}
}



