package service

import (
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/library/mqtt/emt"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/trendssvc/domain"
)

func UnreadNoticeCount(orgId, userId int64) (uint64, errs.SystemErrorInfo) {
	return domain.UnreadNoticeCount(orgId, userId)
}

func NoticeList(orgId, userId int64, page, size int, input *vo.NoticeListReq) (*vo.NoticeList, errs.SystemErrorInfo) {
	count, list, err := domain.GetNoticeList(orgId, userId, page, size, input)
	if err != nil {
		return nil, err
	}

	if len(*list) == 0 {
		return &vo.NoticeList{
			Total: int64(count),
			List:  nil,
		}, nil
	}

	creators := []int64{}
	issueIds := []int64{}
	projectIds := []int64{}
	for _, v := range *list {
		creators = append(creators, v.Creator)
		issueIds = append(issueIds, v.IssueId)
		projectIds = append(projectIds, v.ProjectId)
	}

	//用户信息
	creatorResp, err := orgfacade.GetBaseUserInfoBatchRelaxed("", orgId, slice.SliceUniqueInt64(creators))
	if err != nil {
		return nil, err
	}
	creatorInfo := map[int64]vo.UserIDInfo{}
	for _, v := range creatorResp {
		creatorInfo[v.UserId] = vo.UserIDInfo{
			UserID: v.UserId,
			Avatar: v.Avatar,
			Name:   v.Name,
			EmplID: v.OutUserId,
		}
	}
	//任务信息
	issueResp := projectfacade.GetSimpleIssueInfoBatch(projectvo.GetSimpleIssueInfoBatchReqVo{OrgId: orgId, Ids: slice.SliceUniqueInt64(issueIds)})
	if issueResp.Failure() {
		log.Error(issueResp.Error())
		return nil, issueResp.Error()
	}
	issueInfo := map[int64]string{}
	for _, v := range *(issueResp.Data) {
		issueInfo[v.ID] = v.Title
	}
	//项目信息
	projectResp := projectfacade.GetSimpleProjectInfo(projectvo.GetSimpleProjectInfoReqVo{OrgId: orgId, Ids: slice.SliceUniqueInt64(projectIds)})
	if projectResp.Failure() {
		log.Error(projectResp.Error())
		return nil, projectResp.Error()
	}
	projectInfo := map[int64]string{}
	for _, v := range *projectResp.Data {
		projectInfo[v.ID] = v.Name
	}

	noticeVo := &[]*vo.Notice{}
	copyErr := copyer.Copy(list, noticeVo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	for k, v := range *noticeVo {
		if user, ok := creatorInfo[v.Creator]; ok {
			(*noticeVo)[k].CreatorInfo = &user
		}
		if issueName, ok := issueInfo[v.IssueID]; ok {
			(*noticeVo)[k].IssueName = issueName
		}
		if projectName, ok := projectInfo[v.ProjectID]; ok {
			(*noticeVo)[k].ProjectName = projectName
		}
	}

	return &vo.NoticeList{
		Total: int64(count),
		List:  *noticeVo,
	}, nil
}


func GetMQTTChannelKey(orgId int64, userId int64, input vo.GetMQTTChannelKeyReq) (*vo.GetMQTTChannelKeyResp, errs.SystemErrorInfo) {
	mqttConfig := config.GetMQTTConfig()
	if mqttConfig == nil{
		log.Error("mqtt config is nil")
		return nil, errs.MQTTMissingConfigError
	}

	channelType := input.ChannelType
	channel := ""

	switch channelType {
	case consts.MQTTChannelTypeProject:
		if input.ProjectID == nil{
			return nil, errs.ProjectNotExist
		}
		projectId := *input.ProjectID

		//判断项目权限
		authResp := projectfacade.AuthProjectPermission(projectvo.AuthProjectPermissionReqVo{
			Input: projectvo.AuthProjectPermissionReqData{
				OrgId:      orgId,
				UserId:     userId,
				ProjectId:  projectId,
				Path:       consts.RoleOperationPathOrgPro,
				Operation:  consts.RoleOperationView,
				AuthFiling: false,
			},
		})
		if authResp.Failure() {
			log.Error(authResp.Message)
			return nil, authResp.Error()
		}
		channel = util.GetMQTTProjectChannel(orgId, projectId)
	case consts.MQTTChannelTypeOrg:
		channel = util.GetMQTTOrgChannel(orgId)
	case consts.MQTTChannelTypeUser:
		channel = util.GetMQTTUserChannel(orgId, userId)
	default:
		return nil, errs.IllegalityMQTTChannelType
	}

	key, err := emt.GenerateKey(channel, consts.MQTTDefaultPermissions, consts.MQTTDefaultTTL)
	if err != nil{
		log.Error(err)
		return nil, errs.MQTTKeyGenError
	}

	port := (*int)(nil)
	if mqttConfig.Port > 0{
		port = &mqttConfig.Port
	}
	return &vo.GetMQTTChannelKeyResp{
		Key: key,
		Host: mqttConfig.Host,
		Port: port,
		Channel: channel,
		Address: mqttConfig.Address,
	}, nil
}

