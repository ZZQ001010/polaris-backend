package domain

import (
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"testing"
)

func TestOrgInit(t *testing.T) {
	config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\polaris-server\\configs", "application")

	cache.Set(consts.CacheDingTalkSuiteTicket, "abc")

	conn, _ := mysql.GetConnect()

	tx, _ := conn.NewTx(nil)
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
		if tx != nil{
			if err := tx.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()

	_, err := OrgInit("ding8ac2bab2b708b3cc35c2f4657eb6378f", "efg", tx)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	} else {
		err := tx.Commit()
		if err != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	}

}
