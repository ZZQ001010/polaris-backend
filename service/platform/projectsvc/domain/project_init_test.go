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

func ConfigInit() {
	config.LoadConfig("F:\\polaris-backend\\polaris-server\\configs", "application")

	cache.Set(consts.CacheDingTalkSuiteTicket, "abc")
}

func TestProjectTypeInit(t *testing.T) {
	ConfigInit()

	conn, _ := mysql.GetConnect()

	tx, _ := conn.NewTx(nil)
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
		if err := tx.Close(); err != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	}()

	err := ProjectTypeInit(1000, tx)
	if err != nil {
		logger.GetDefaultLogger().Info("初始化project_type失败:" + strs.ObjectToString(err))
		err1 := tx.Rollback()
		if err1 != nil {
			logger.GetDefaultLogger().Info("Rollback error:" + strs.ObjectToString(err1))
		}
	} else {
		err2 := tx.Commit()
		if err2 != nil {
			logger.GetDefaultLogger().Info("Commit error :" + strs.ObjectToString(err2))
		}
	}
}

func TestProjectProjectTypeInit(t *testing.T) {
	ConfigInit()

	conn, _ := mysql.GetConnect()

	tx, _ := conn.NewTx(nil)
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
		if tx != nil {
			if err := tx.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()

	contextMap := make(map[string]interface{})

	err := ProjectObjectTypeInit(1000, &contextMap, tx)
	if err != nil {
		logger.GetDefaultLogger().Info("初始化project_object_type失败:" + strs.ObjectToString(err))
		err1 := tx.Rollback()
		if err1 != nil {
			logger.GetDefaultLogger().Info("Rollback error:" + strs.ObjectToString(err1))
		}
	} else {
		err2 := tx.Commit()
		if err2 != nil {
			logger.GetDefaultLogger().Info("Commit error:" + strs.ObjectToString(err2))
		}
	}
}

func TestProcessInit(t *testing.T) {
	ConfigInit()

	conn, _ := mysql.GetConnect()

	tx, _ := conn.NewTx(nil)
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
		if tx != nil {
			if err := tx.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()

	err := ProcessInit(1, tx)
	if err != nil {
		logger.GetDefaultLogger().Info("初始化process失败" + strs.ObjectToString(err))
		err1 := tx.Rollback()
		if err1 != nil {
			logger.GetDefaultLogger().Info("Rollback error" + strs.ObjectToString(err))
		}
	} else {
		err2 := tx.Commit()
		if err2 != nil {
			logger.GetDefaultLogger().Info("Commit error" + strs.ObjectToString(err))
		}
	}
}

func TestProjectTypeProjectObjectTypeInit(t *testing.T) {
	ConfigInit()

	conn, _ := mysql.GetConnect()

	tx, _ := conn.NewTx(nil)
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
		if tx != nil {
			if err := tx.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()

	err := ProjectTypeProjectObjectTypeInit(1000, tx)
	if err != nil {
		logger.GetDefaultLogger().Info("初始化project_type_project_object_type失败" + strs.ObjectToString(err))
		err1 := tx.Rollback()
		if err1 != nil {
			logger.GetDefaultLogger().Info("Rollback error" + strs.ObjectToString(err1))
		}
	} else {
		err2 := tx.Commit()
		if err2 != nil {
			logger.GetDefaultLogger().Info("Commit error" + strs.ObjectToString(err2))
		}
	}
}

func TestPriorityInit(t *testing.T) {
	ConfigInit()

	conn, _ := mysql.GetConnect()

	tx, _ := conn.NewTx(nil)
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
		if tx != nil {
			if err := tx.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()

	err := PriorityInit(1001, tx)
	if err != nil {
		logger.GetDefaultLogger().Info("初始化priority失败" + strs.ObjectToString(err))
		err1 := tx.Rollback()
		if err1 != nil {
			logger.GetDefaultLogger().Info("Rollback error" + strs.ObjectToString(err1))
		}
	} else {
		err2 := tx.Commit()
		if err2 != nil {
			logger.GetDefaultLogger().Info("Commit error" + strs.ObjectToString(err2))
		}
	}
}

//func TestProjectInit(t *testing.T) {
//	ConfigInit()
//
//	conn, _ := mysql.GetConnect()
//
//	tx, _ := conn.NewTx(nil)
//	defer func() {
//		if err := conn.Close(); err != nil {
//			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
//		}
//		if err := tx.Close(); err != nil {
//			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
//		}
//	}()
//
//	contextMap := make(map[string]interface{})
//
//	err := ProjectInit(6, &contextMap, tx)
//	if err != nil {
//		logger.GetDefaultLogger().Info("初始化project失败", err)
//		err1 := tx.Rollback()
//		if err1 != nil {
//			logger.GetDefaultLogger().Info("Rollback error", err1)
//		}
//	} else {
//		err2 := tx.Commit()
//		if err2 != nil {
//			logger.GetDefaultLogger().Info("Commit error", err2)
//		}
//	}
//}
