package vo

import (
	"github.com/galaxy-book/common/core/errors"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
)

type Err struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e Err) Successful() bool {
	if e.Code == errs.OK.Code() {
		return true
	}
	return false
}

func (e Err) Failure() bool {
	return !e.Successful()
}

func (e Err) Error() errors.SystemErrorInfo {
	if e.Successful() {
		return nil
	}
	return errs.BuildSystemErrorInfo(errs.GetResultCodeInfoByCode(e.Code))
}

type VoidErr struct {
	Err
}

func NewErr(err errs.SystemErrorInfo) Err {
	if err == nil {
		err = errs.OK
	}
	return Err{
		Code:    err.Code(),
		Message: err.Message(),
	}
}

type CommonReqVo struct {
	UserId int64 `json:"userId"`
	OrgId int64 `json:"orgId"`
	SourceChannel string `json:"sourceChannel"`
}

type CommonRespVo struct {
	Err
	Void *Void `json:"data"`
}

type BasicReqVo struct {
	Page uint
	Size uint
}

type BoolRespVo struct {
	Err
	IsTrue bool `json:"data"`
}

type BasicInfoReqVo struct {
	UserId int64 `json:"userId"`
	OrgId  int64 `json:"orgId"`
}
