package domain

import (
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/common/sdk/dingtalk"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"testing"
)

func TestUserInit(t *testing.T) {

	config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\polaris-server\\configs", "application")

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
	cache.Set(consts.CacheDingTalkSuiteTicket, "abc")
	_, err := UserInitByOrg("manager5225", "ding8ac2bab2b708b3cc35c2f4657eb6378f", 0, tx)

	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	} else {
		err2 := tx.Commit()
		if err2 != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	}
}

func TestGetUserList(t *testing.T) {
	config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\polaris-server\\configs", "application")

	c, _ := dingtalk.GetDingTalkClient("ding8ac2bab2b708b3cc35c2f4657eb6378f", "_")
	fmt.Println(c.GetDepMemberIds("1"))
}
