package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/processvo"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/service"
)

func (PostGreeter) ProcessStatusList(input vo.BasicReqVo) processvo.ProcessStatusListRespVo {
	res, err := service.ProcessStatusList(input.Page, input.Size)
	return processvo.ProcessStatusListRespVo{Err: vo.NewErr(err), ProcessStatusList: res}
}

func (PostGreeter) CreateProcessStatus(req processvo.CreateProcessStatusReqVo) vo.CommonRespVo {
	res, err := service.CreateProcessStatus(req)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateProcessStatus(req processvo.UpdateProcessStatusReqVo) vo.CommonRespVo {
	res, err := service.UpdateProcessStatus(req)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) DeleteProcessStatus(req processvo.DeleteProcessStatusReq) vo.CommonRespVo {
	res, err := service.DeleteProcessStatus(req)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (GetGreeter) GetProcessStatus(req processvo.GetProcessStatusReqVo) processvo.GetProcessStatusRespVo {
	cacheProcessStatusBo, err := service.GetProcessStatus(req.OrgId, req.Id)
	return processvo.GetProcessStatusRespVo{CacheProcessStatusBo: cacheProcessStatusBo, Err: vo.NewErr(err)}
}

func (PostGreeter) ProcessStatusInit(req processvo.ProcessStatusInitReqVo) vo.VoidErr {
	err := service.ProcessStatusInit(req.OrgId, req.ContextMap)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (GetGreeter) GetProcessStatusByCategory(req processvo.GetProcessStatusByCategoryReqVo) processvo.GetProcessStatusByCategoryRespVo {
	res, err := service.GetProcessStatusByCategory(req.OrgId, req.StatusId, req.Category)
	return processvo.GetProcessStatusByCategoryRespVo{CacheProcessStatusBo: res, Err: vo.NewErr(err)}
}

func (GetGreeter) GetProcessStatusListByCategory(req processvo.GetProcessStatusListByCategoryReqVo) processvo.GetProcessStatusListByCategoryRespVo {
	res, err := service.GetProcessStatusListByCategory(req.OrgId, req.Category)
	return processvo.GetProcessStatusListByCategoryRespVo{CacheProcessStatusBoList: res, Err: vo.NewErr(err)}
}

func (GetGreeter) GetProcessStatusIds(req processvo.GetProcessStatusIdsReqVo) processvo.GetProcessStatusIdsRespVo {
	res, err := service.GetProcessStatusIds(req.OrgId, req.Category, req.Typ)
	return processvo.GetProcessStatusIdsRespVo{ProcessStatusIds: res, Err: vo.NewErr(err)}
}

func (GetGreeter) GetProcessStatusList(req processvo.GetProcessStatusListReqVo) processvo.GetProcessStatusListRespVo {
	res, err := service.GetProcessStatusList(req.OrgId, req.ProcessId)
	return processvo.GetProcessStatusListRespVo{ProcessStatusBoList: res, Err: vo.NewErr(err)}
}

func (GetGreeter) GetProcessInitStatusId(req processvo.GetProcessInitStatusIdReqVo) processvo.GetProcessInitStatusIdRespVo {
	res, err := service.GetProcessInitStatusId(req.OrgId, req.ProjectId, req.ProjectObjectTypeId, req.Category)
	return processvo.GetProcessInitStatusIdRespVo{ProcessInitStatusId: res, Err: vo.NewErr(err)}
}

func (GetGreeter) GetDefaultProcessStatusId(reqVo processvo.GetDefaultProcessIdReqVo) processvo.GetDefaultProcessIdRespVo {
	res, err := service.GetDefaultProcessStatusId(reqVo.OrgId, reqVo.ProcessId, reqVo.Category)
	return processvo.GetDefaultProcessIdRespVo{Err: vo.NewErr(err), ProcessId: res}
}
