package api

import (
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/processvo"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/service"
	"upper.io/db.v3/lib/sqlbuilder"
)

func (PostGreeter) InitProcess(req processvo.InitProcessReqVo) vo.VoidErr {
	err := service.InitProcess(req.OrgId)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) GetProcessByLangCode(req processvo.GetProcessByLangCodeReqVo) processvo.GetProcessByLangCodeRespVo {
	res, err := service.GetProcessByLangCode(req.OrgId, req.LangCode)
	return processvo.GetProcessByLangCodeRespVo{ProcessBo: res, Err: vo.NewErr(err)}
}

func (PostGreeter) GetProcessBo(req processvo.GetProcessBoReqVo) processvo.GetProcessBoRespVo {
	res, err := service.GetProcessBo(req.Cond)
	return processvo.GetProcessBoRespVo{ProcessBo: res, Err: vo.NewErr(err)}
}

func (PostGreeter) AssignValueToField(reqVo processvo.AssignValueToFieldReqVo) vo.VoidErr {
	respVo := vo.VoidErr{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := service.AssignValueToField(reqVo.ProcessRes, tx, reqVo.OrgId)
		respVo.Err = vo.NewErr(err)
		return err
	})

	return respVo
}

func (GetGreeter) GetNextProcessStepStatusList(reqVo processvo.GetNextProcessStepStatusListReqVo) processvo.GetNextProcessStepStatusListRespVo {
	res, err := service.GetNextProcessStepStatusList(reqVo.OrgId, reqVo.ProcessId, reqVo.StartStatusId)
	return processvo.GetNextProcessStepStatusListRespVo{Err: vo.NewErr(err), CacheProcessStatus: res}
}

func (GetGreeter) GetProcessById(req processvo.GetProcessByIdReqVo) processvo.GetProcessByIdRespVo {
	res, err := service.GetProcessById(req.OrgId, req.Id)
	return processvo.GetProcessByIdRespVo{ProcessBo: res, Err: vo.NewErr(err)}
}
