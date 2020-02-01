package api

import (
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
	"upper.io/db.v3/lib/sqlbuilder"
)

func (PostGreeter) ProjectInit(req projectvo.ProjectInitReqVo) projectvo.ProjectInitRespVo {
	respVo := projectvo.ProjectInitRespVo{}

	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		contextMap, err := service.ProjectInit(req.OrgId, tx)

		respVo.ContextMap = contextMap
		respVo.Err = vo.NewErr(err)

		return err
	})
	return respVo
}

func (PostGreeter) DataInitForLarkApplet(req vo.BasicInfoReqVo) vo.CommonRespVo {
	err := service.DataInitForLarkApplet(req.OrgId, req.UserId)
	return vo.CommonRespVo{Err: vo.NewErr(err)}
}
