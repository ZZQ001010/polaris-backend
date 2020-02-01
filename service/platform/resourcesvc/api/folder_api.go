package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/service/platform/resourcesvc/service"
)

func (PostGreeter) CreateFolder(reqVo resourcevo.CreateFolderReqVo) vo.CommonRespVo {
	//if len(reqVo.Input.Name) > 15 || reqVo.Input.Name == "" {
	//	return vo.CommonRespVo{Err: vo.NewErr(errs.InvalidResourceNameError), Void: nil}
	//}
	res, err := service.CreateFolder(reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateFolder(reqVo resourcevo.UpdateFolderReqVo) resourcevo.UpdateFolderRespVo {
	//if reqVo.Input.Name != nil && (len(*reqVo.Input.Name) > 15 || *reqVo.Input.Name == "") {
	//	return resourcevo.UpdateFolderRespVo{Err: vo.NewErr(errs.InvalidResourceNameError), FolderIds: nil}
	//}
	res, err := service.UpdateFolder(reqVo.Input)
	return resourcevo.UpdateFolderRespVo{Err: vo.NewErr(err), UpdateFolderData: res}
}

func (PostGreeter) DeleteFolder(reqVo resourcevo.DeleteFolderReqVo) resourcevo.DeleteFolderRespVo {
	res, err := service.DeleteFolder(reqVo.Input)
	return resourcevo.DeleteFolderRespVo{Err: vo.NewErr(err), DeleteFolderData: res}
}

func (PostGreeter) GetFolder(reqVo resourcevo.GetFolderReqVo) resourcevo.GetFolderVoListRespVo {
	res, err := service.GetFolder(reqVo.Input)
	return resourcevo.GetFolderVoListRespVo{FolderList: res, Err: vo.NewErr(err)}
}
