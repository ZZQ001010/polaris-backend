package idfacade

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/idvo"
)

var log = logger.GetDefaultLogger()

func ApplyPrimaryIdRelaxed(code string) (int64, errs.SystemErrorInfo) {
	respVo := ApplyPrimaryId(idvo.ApplyPrimaryIdReqVo{Code: code,})
	if respVo.Failure() {
		return 0, respVo.Error()
	}
	return respVo.Id, nil
}

func ApplyMultipleIdRelaxed(orgId int64, code, preCode string, count int64) (*bo.IdCodes, errs.SystemErrorInfo) {
	respVo := ApplyMultipleId(idvo.ApplyMultipleIdReqVo{
		Code:    code,
		OrgId:   orgId,
		PreCode: preCode,
		Count:   int(count),
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.IdCodes, nil
}

func ApplyIdRelaxed(orgId int64, code string, preCode string) (*bo.IdCodes, errs.SystemErrorInfo) {
	respVo := ApplyMultipleId(idvo.ApplyMultipleIdReqVo{
		Code:    code,
		OrgId:   orgId,
		PreCode: preCode,
		Count:   int(1),
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.IdCodes, nil
}

func ApplyMultiplePrimaryIdRelaxed(code string, count int) (*bo.IdCodes, errs.SystemErrorInfo) {
	respVo := ApplyMultiplePrimaryId(idvo.ApplyMultiplePrimaryIdReqVo{
		Code:    code,
		Count:   count,
	})
	if respVo.Failure() {
		return nil, respVo.Error()
	}
	return respVo.IdCodes, nil
}