package service

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/service/platform/rolesvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRoleInit(t *testing.T) {
	config.LoadConfig("F:\\polaris-backend\\polaris-server\\configs", "application")

	convey.Convey("Test GetRoleOperationList", t, test.StartUp(func(ctx context.Context) {
		convey.Convey("权限init", func() {
			conn, _ := mysql.GetConnect()
			tx, _ := conn.NewTx(nil)
			_, err := RoleInit(1, tx)
			if err != nil {
				tx.Rollback()
				fmt.Println(err)
			}
			tx.Commit()
		})
	}))
}
