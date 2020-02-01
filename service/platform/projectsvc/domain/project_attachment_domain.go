package domain

import (
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

func AssemblyAttachmentList(resourceList []*vo.Resource, issueBoList []bo.IssueBo, relationMap map[int64][]int64) ([]bo.AttachmentBo, errs.SystemErrorInfo) {
	issueBoMap := map[int64]bo.IssueBo{}
	for _, value := range issueBoList {
		issueBoMap[value.Id] = value
	}
	resourceBoMap := map[int64]vo.Resource{}
	for _, value := range resourceList {
		resourceBoMap[value.ID] = *value
	}
	bos := make([]bo.AttachmentBo, 0)
	for resourceId, issueIdList := range relationMap {
		attachmentBo := bo.AttachmentBo{}
		if value, ok := resourceBoMap[resourceId]; ok {
			attachmentBo.Resource = value
		} else {
			log.Errorf("ResourceId %d not in resourceBoMap", resourceId)
			return nil, errs.ResourceNotExist
		}
		issueBoList := make([]bo.IssueBo, 0)
		for _, issueId := range issueIdList {
			if value, ok := issueBoMap[issueId]; ok {
				issueBo := value
				issueBoList = append(issueBoList, issueBo)
			} else {
				log.Errorf("IssueId %d not in issueBoMap", issueId)
				//return nil, errs.IssueNotExist
				continue
			}
		}
		attachmentBo.IssueList = issueBoList
		bos = append(bos, attachmentBo)
	}
	return bos, nil
}
