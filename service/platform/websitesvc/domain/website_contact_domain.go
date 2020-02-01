package domain

import (
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/websitesvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/websitesvc/po"
	"upper.io/db.v3"
)

func RegisterWebSiteContact(bo bo.ContactBo) (int64, errs.SystemErrorInfo) {
	po := &po.PpmWstContact{}
	copyErr := copyer.Copy(bo, po)
	if copyErr != nil {
		log.Error(copyErr)
		return 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	contactId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableContact)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	po.Id = contactId

	fmt.Println(json.ToJsonIgnoreError(config.GetMysqlConfig()))
	err1 := basedao.InsertContact(*po)
	if err1 != nil {
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}
	return contactId, nil
}

func CheckContactRepetition(contactInfo string, status int) (bool, errs.SystemErrorInfo) {
	total, err := basedao.SelectCountContact(db.Cond{
		consts.TcMobile:   contactInfo,
		consts.TcStatus:   status,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err != nil {
		log.Error(err)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	if total > 0 {
		return true, nil
	}
	return false, nil
}
