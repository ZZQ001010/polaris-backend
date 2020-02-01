package bo

import "github.com/galaxy-book/polaris-backend/common/model/vo"

type AttachmentBo struct {
	vo.Resource
	IssueList []IssueBo
}
