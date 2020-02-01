package domain

import (
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/po"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"upper.io/db.v3"
)

func TestTeamInit(t *testing.T) {

	config.LoadEnvConfig("F:\\polaris-backend-clone\\config", "application.common", "local")

	conn, _ := mysql.GetConnect()
	tx, _ := conn.NewTx(nil)

	TeamInit(2, tx)
	tx.Commit()

}

func TestGetCorpAuthInfo(t *testing.T) {
	convey.Convey("测试加载env2", t, func() {
		config.LoadEnvConfig("config", "application", "dev")
		info := &[]po.PpmOrgDepartmentOutInfo{}
		mysql.SelectAllByCond(consts.TableDepartmentOutInfo, db.Cond{}, info)
		fmt.Println(info)
	})
}
