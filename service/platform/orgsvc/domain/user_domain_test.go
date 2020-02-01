package domain

import (
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"testing"
)

func TestGetUserInfoListByOrg(t *testing.T) {
	config.LoadEnvConfig("F:\\polaris-backend-clone\\config", "application.common", "local")
	//config.LoadEnvConfig("F:\\polaris-backend-clone\\service\\\\config", "application", "local")

	fmt.Println("1111", config.GetConfig().Mysql.Host)
	res, _ := GetUserInfoListByOrg(17)
	t.Log(json.ToJsonIgnoreError(res))
}
