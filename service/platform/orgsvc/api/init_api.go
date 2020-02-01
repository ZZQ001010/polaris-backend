package api

import (
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/orgvo"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/service"
	"upper.io/db.v3/lib/sqlbuilder"
)

func (PostGreeter) OrgInit(reqVo orgvo.OrgInitReqVo) orgvo.OrgInitRespVo {
	respVo := orgvo.OrgInitRespVo{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		orgId, err := service.OrgInit(reqVo.CorpId, reqVo.PermanentCode, tx)
		respVo.OrgId = orgId
		respVo.Err = vo.NewErr(err)
		return err
	})

	return respVo
}

//飞书初始化调用
func (PostGreeter) InitOrg(reqVo orgvo.InitOrgReqVo) orgvo.OrgInitRespVo {
	respVo := orgvo.OrgInitRespVo{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		orgId, err := service.InitOrg(reqVo.InitOrg, tx)
		respVo.OrgId = orgId
		respVo.Err = vo.NewErr(err)
		return err
	})
	return respVo
}

//通用初始化调用
func (PostGreeter) GeneralInitOrg(reqVo orgvo.InitOrgReqVo) orgvo.OrgInitRespVo {
	respVo := orgvo.OrgInitRespVo{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		orgId, err := service.GeneralInitOrg(reqVo.InitOrg, tx)
		respVo.OrgId = orgId
		respVo.Err = vo.NewErr(err)
		return err
	})
	return respVo
}

func (PostGreeter) OrgOwnerInit(reqVo orgvo.OrgOwnerInitReqVo) vo.VoidErr {
	respVo := vo.VoidErr{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := service.OrgOwnerInit(reqVo.OrgId, reqVo.Owner, tx)
		respVo.Err = vo.NewErr(err)
		return err
	})

	return respVo
}

func (PostGreeter) OrgSysConfigInit(reqVo orgvo.OrgVo) vo.VoidErr {
	respVo := vo.VoidErr{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := service.OrgSysConfigInit(tx, reqVo.OrgId)
		respVo.Err = vo.NewErr(err)
		return err
	})

	return respVo
}

func (PostGreeter) TeamInit(reqVo orgvo.OrgVo) orgvo.TeamInitRespVo {
	respVo := orgvo.TeamInitRespVo{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		id, err := service.TeamInit(reqVo.OrgId, tx)
		respVo.TeamId = id
		respVo.Err = vo.NewErr(err)
		return err
	})

	return respVo
}

func (PostGreeter) TeamOwnerInit(reqVo orgvo.TeamOwnerInitReqVo) vo.VoidErr {
	respVo := vo.VoidErr{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := service.TeamOwnerInit(reqVo.TeamId, reqVo.Owner, tx)
		respVo.Err = vo.NewErr(err)
		return err
	})

	return respVo
}

func (PostGreeter) TeamUserInit(reqVo orgvo.TeamUserInitReqVo) vo.VoidErr {
	respVo := vo.VoidErr{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := service.TeamUserInit(reqVo.TeamId, reqVo.TeamId, reqVo.UserId, reqVo.IsRoot, tx)
		respVo.Err = vo.NewErr(err)
		return err
	})

	return respVo
}

func (PostGreeter) UserInitByOrg(reqVo orgvo.UserInitByOrgReqVo) orgvo.UserInitByOrgRespVo {
	respVo := orgvo.UserInitByOrgRespVo{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		id, err := service.UserInitByOrg(reqVo.UserId, reqVo.CorpId, reqVo.OrgId, tx)
		respVo.UserId = id
		respVo.Err = vo.NewErr(err)
		return err
	})

	return respVo
}

func (PostGreeter) DepartmentInit(reqVo orgvo.DepartmentInitReqVo) vo.VoidErr {
	respVo := vo.VoidErr{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := service.InitDepartment(reqVo.OrgId, reqVo.CorpId, reqVo.SourceChannel, tx)
		respVo.Err = vo.NewErr(err)
		return err
	})
	return respVo
}
