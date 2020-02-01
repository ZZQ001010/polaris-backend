package service

import (
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestCreateOrg(t *testing.T) {

	orgName := "0123456789123456"
	orgNameLen := strs.Len(orgName)

	var err error = nil
	if orgNameLen == 0 || orgNameLen > 256 {
		log.Error("组织名称长度错误")
		err = errs.BuildSystemErrorInfo(errs.OrgNameLenError)
	}

	assert.Equal(t, err, nil)

}