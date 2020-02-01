package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/facade/resourcefacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"upper.io/db.v3"
)

func DeleteProjectAttachment(orgId, operatorId int64, input vo.DeleteProjectAttachmentReq) (*vo.DeleteProjectAttachmentResp, errs.SystemErrorInfo) {
	err := domain.AuthProject(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProAttachment, consts.RoleOperationDelete)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = domain.DeleteProjectAttachment(orgId, operatorId, input.ProjectID, input.ResourceIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.DeleteProjectAttachmentResp{
		ResourceIds: input.ResourceIds,
	}, nil
}

func GetProjectAttachment(orgId, operatorId int64, page, size int, input vo.ProjectAttachmentReq) (*vo.AttachmentList, errs.SystemErrorInfo) {
	err := domain.AuthProjectWithOutPermission(orgId, operatorId, input.ProjectID, consts.RoleOperationPathOrgProAttachment, consts.RoleOperationView)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resourceInput := bo.GetResourceBo{
		UserId:     operatorId,
		OrgId:      orgId,
		ProjectId:  input.ProjectID,
		Page:       page,
		Size:       size,
		SourceType: consts.OssPolicyTypeIssueResource,
	}
	if input.FileType != nil {
		resourceInput.FileType = input.FileType
	}
	if input.KeyWord != nil {
		resourceInput.KeyWord = input.KeyWord
	}
	resp := resourcefacade.GetResource(resourcevo.GetResourceReqVo{
		Input: resourceInput,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}
	resourceList := resp.ResourceList.List
	total := resp.ResourceList.Total
	resourceIds := make([]int64, len(resourceList))
	for i, value := range resourceList {
		resourceIds[i] = value.ID
	}
	cond := db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcProjectId:    input.ProjectID,
		consts.TcOrgId:        orgId,
		consts.TcRelationType: consts.IssueRelationTypeResource,
		consts.TcRelationId:   db.In(resourceIds),
	}
	issueRelationPos, err := domain.GetTotalResourceByRelationCond(cond)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//获取resourceId和issueId的映射关系
	relationMap := map[int64][]int64{}
	issueIdMap := map[int64]bool{}
	for _, value := range *issueRelationPos {
		issueIdMap[value.IssueId] = true
		issueId := value.IssueId
		resourceId := value.RelationId
		relationMap[resourceId] = append(relationMap[resourceId], issueId)
		//fmt.Printf("relationId:%d,issueId:%d,resourceId:%d\n", value.Id, value.IssueId, value.RelationId)
	}
	issueIds := make([]int64, len(issueIdMap))
	for issueId, _ := range issueIdMap {
		issueIds = append(issueIds, issueId)
	}
	issueBoList, err := domain.GetIssueInfoList(issueIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	attachmentBos, err := domain.AssemblyAttachmentList(resourceList, issueBoList, relationMap)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	list := make([]*vo.Attachment, 0)
	err1 := copyer.Copy(attachmentBos, &list)
	if err1 != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	return &vo.AttachmentList{Total: total, List: list}, nil
}
