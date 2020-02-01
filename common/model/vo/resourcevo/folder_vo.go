package resourcevo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type CreateFolderReqVo struct {
	Input bo.CreateFolderBo `json:"createfolder"`
}

type UpdateFolderReqVo struct {
	Input bo.UpdateFolderBo `json:"updateFolderBo"`
}

type UpdateFolderRespVo struct {
	*UpdateFolderData
	vo.Err
}

type GetFolderReqVo struct {
	Input bo.GetFolderBo `json:"getFolderBo"`
}

type GetFolderVoListRespVo struct {
	*vo.FolderList `json:"data"`

	vo.Err
}

type UpdateFolderData struct {
	FolderId     int64
	FolderName   *string
	UpdateFields []string
	OldValue     *string
	NewValue     *string
}

type DeleteFolderReqVo struct {
	Input bo.DeleteFolderBo `json:"deleteFolder"`
}

type DeleteFolderRespVo struct {
	*DeleteFolderData
	vo.Err
}

type DeleteFolderData struct {
	FolderIds   []int64
	FolderNames []string
}
