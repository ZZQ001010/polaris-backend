package api

import (
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/service/platform/resourcesvc/domain"
	"github.com/galaxy-book/polaris-backend/service/platform/resourcesvc/service"
	"upper.io/db.v3/lib/sqlbuilder"
)

func (PostGreeter) CreateResource(req resourcevo.CreateResourceReqVo) resourcevo.CreateResourceRespVo {
	respVo := resourcevo.CreateResourceRespVo{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		resourceId, err := service.CreateResource(req.CreateResourceBo, tx)
		respVo.ResourceId = resourceId
		respVo.Err = vo.NewErr(err)
		return err
	})
	return respVo
}

func (PostGreeter) UpdateResourceInfo(reqVo resourcevo.UpdateResourceInfoReqVo) resourcevo.UpdateResourceInfoResVo {
	res, err := service.UpdateResourceInfo(reqVo.Input)
	return resourcevo.UpdateResourceInfoResVo{UpdateResourceData: res, Err: vo.NewErr(err)}
}

func (PostGreeter) UpdateResourceFolder(reqVo resourcevo.UpdateResourceFolderReqVo) resourcevo.UpdateResourceInfoResVo {
	res, err := service.UpdateResourceFolder(reqVo.Input)
	return resourcevo.UpdateResourceInfoResVo{UpdateResourceData: res, Err: vo.NewErr(err)}
}

//func (PostGreeter) UpdateResourceParentId(reqVo resourcevo.UpdateResourceParentIdReqVo) vo.CommonRespVo {
//	res, err := service.UpdateResourceParentId(reqVo.Input)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}

func (PostGreeter) DeleteResource(reqVo resourcevo.DeleteResourceReqVo) resourcevo.UpdateResourceInfoResVo {
	res, err := service.DeleteResource(reqVo.Input)
	return resourcevo.UpdateResourceInfoResVo{Err: vo.NewErr(err), UpdateResourceData: res}
}

func (PostGreeter) GetResource(reqVo resourcevo.GetResourceReqVo) resourcevo.GetResourceVoListRespVo {
	res, err := service.GetResource(reqVo.Input)
	return resourcevo.GetResourceVoListRespVo{ResourceList: res, Err: vo.NewErr(err)}
}

//func (PostGreeter) InsertResource(req resourcevo.InsertResourceReqVo) resourcevo.InsertResourceRespVo {
//	input := req.Input
//	respVo := resourcevo.InsertResourceRespVo{}
//	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//		resourceId, err := service.InsertResource(tx, input.ResourcePath, input.OrgId, input.CurrentUserId, input.ResourceType, input.FileName)
//		respVo.ResourceId = resourceId
//		respVo.Err = vo.NewErr(err)
//		return err
//	})
//	return respVo
//}

func (PostGreeter) GetResourceById(req resourcevo.GetResourceByIdReqVo) resourcevo.GetResourceByIdRespVo {
	resourceBos, err := domain.GetResourceByIds(req.GetResourceByIdReqBody.ResourceIds)
	return resourcevo.GetResourceByIdRespVo{ResourceBos: resourceBos, Err: vo.NewErr(err)}
}

func (GetGreeter) GetIdByPath(req resourcevo.GetIdByPathReqVo) resourcevo.GetIdByPathRespVo {
	resourceId, err := service.GetIdByPath(req.OrgId, req.ResourcePath, req.ResourceType)
	return resourcevo.GetIdByPathRespVo{ResourceId: resourceId, Err: vo.NewErr(err)}
}

func (PostGreeter) GetResourceBoList(req resourcevo.GetResourceBoListReqVo) resourcevo.GetResourceBoListRespVo {
	list, total, err := domain.GetResourceBoList(req.Page, req.Size, req.Input)
	return resourcevo.GetResourceBoListRespVo{GetResourceBoListRespData: resourcevo.GetResourceBoListRespData{ResourceBos: list, Total: total}, Err: vo.NewErr(err)}
}
