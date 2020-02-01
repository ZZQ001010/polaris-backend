package idvo

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"testing"
)

func Test_Convert(t *testing.T){
	obj := ApplyPrimaryIdRespVo{
		Id: 123,
		Err: vo.NewErr(errs.SystemError),
	}
	t.Log(json.ToJsonIgnoreError(obj))

	obj1 := vo.VoidErr{
		Err: vo.NewErr(errs.SystemError),
	}
	t.Log(json.ToJsonIgnoreError(obj1))
}