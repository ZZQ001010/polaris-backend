package dingtalk

import (
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/tests"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

var env = ""

const BaseConfigPath = "./../../../config"

const testPaht = "./"

func init() {
	env = os.Getenv(consts.RunEnvKey)
	if "" == env {
		env = "unittest"
	}

	//测试配置文件
	//err := config.LoadEnvConfig(testPaht, "application.common", env)
	//
	//if  err != nil {
	//	fmt.Printf("err:%s\n", err)
	//}

	//配置文件
	err := config.LoadEnvConfig(BaseConfigPath, "application.common", env)

	if err != nil {
		fmt.Printf("err:%s\n", err)
	}

	//为了加载钉钉的配置文件
	err = config.LoadEnvConfig(BaseConfigPath, "application.common", env)

	if err != nil {
		fmt.Printf("err:%s\n", err)
	}

	fmt.Println("load env finish")
}

func TestGetDingTalkUserRoleBos(t *testing.T) {

	convey.Convey("更新任务关联测试", t, tests.StartUp(func() {

		bos, err := GetDingTalkUserRoleBos("ding8ac2bab2b708b3cc35c2f4657eb6378f")
		t.Log(err)
		t.Log(json.ToJson(bos))
	}))
}
