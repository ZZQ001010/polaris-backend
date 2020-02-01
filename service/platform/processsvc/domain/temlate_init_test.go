package domain

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/service/platform/processsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestVariableInit(t *testing.T) {

	convey.Convey("Test InitUserRoles", t, test.StartUp(func(ctx context.Context) {

		conn, _ := mysql.GetConnect()
		tx, _ := conn.NewTx(nil)

		contentMap := make(map[string]interface{})

		err := ProcessStatusInit(17, contentMap, tx)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
		}
		tx.Commit()
	}))
}

func TestWriteWithIoutil(t *testing.T) {
	util.WriteWithIoutil("F:\\polaris-backend\\polaris-sysinit\\inits\\ppm_prs_process_status_new_test.template", "abcd")
}

func TestFileRead(t *testing.T) {

	str, err := util.FileRead("../resources/template/ppm_prs_process_status.template")
	fmt.Println(str, err)
}

func TestTemplate(t *testing.T) {
	convey.Convey("Test InitUserRoles", t, test.StartUp(func(ctx context.Context) {
		conn, _ := mysql.GetConnect()
		tx, _ := conn.NewTx(nil)

		contentMap := make(map[string]interface{})
		err := util.ReadAndWrite("../resources/template/test.template", contentMap, tx)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
		}
		tx.Commit()
	}))
}
