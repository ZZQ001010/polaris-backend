package main

import (
	"flag"
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/polaris-backend/common/core/buildinfo"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/extra/gin/mid"
	"github.com/galaxy-book/polaris-backend/common/extra/gin/mvc"
	"github.com/galaxy-book/polaris-backend/service/basic/appsvc/api"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"os"
	"runtime"
	"strconv"
)

var log = logger.GetDefaultLogger()
var env = ""
var build = false

const BaseConfigPath = "./../../../config"
const SelfConfigPath = "./config"

func init() {
	env = os.Getenv(consts.RunEnvKey)
	if "" == env {
		env = consts.RunEnvLocal
	}
	//配置
	flag.BoolVar(&build, "build", false, "build facade")
	flag.StringVar(&env, "env", env, "env")
	flag.Parse()

	//配置文件
	if runtime.GOOS != consts.LinuxGOOS {
		//配置文件
		config.LoadEnvConfig(BaseConfigPath, "application.common", env)
		config.LoadEnvConfig(SelfConfigPath, "application", env)
	} else {
		//配置文件
		config.LoadEnvConfig(SelfConfigPath, "application.common", env)
		config.LoadEnvConfig(SelfConfigPath, "application", env)
	}
}

func main() {
	// 打印程序信息
	log.Info(buildinfo.StringifySingleLine())
	fmt.Println(buildinfo.StringifyMultiLine())

	port := config.GetServerConfig().Port
	host := config.GetServerConfig().Host

	applicationName := config.GetApplication().Name

	r := gin.New()

	r.Use(mid.StartTrace())
	r.Use(mid.GinContextToContextMiddleware())
	r.Use(mid.CorsMiddleware())
	r.Use(mid.AuthMiddleware())

	version := "v1"
	postGreeter := api.PostGreeter{Greeter: mvc.NewPostGreeter(applicationName, host, port, version)}
	getGreeter := api.GetGreeter{Greeter: mvc.NewGetGreeter(applicationName, host, port, version)}

	//build
	if build {
		facadeBuilder := mvc.FacadeBuilder{
			StorageDir: "./../../../facade/appfacade",
			Package:    "appfacade",
			VoPackage:  "appvo",
			Greeters:   []interface{}{&postGreeter, &getGreeter},
		}
		facadeBuilder.Build()
		return
	}

	// 多库库模式才会执行
	//if (consts.AppRunmodeSaas == config.GetApplication().RunMode) || (consts.AppRunmodePrivate == config.GetApplication().RunMode) {
	//	mysqlConfig := config.GetMysqlConfig()
	//	initErr := db.DbMigrations(env, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Usr, mysqlConfig.Pwd, mysqlConfig.Database)
	//	if initErr != nil {
	//		panic(" init db fail....")
	//	}
	//}

	ginHandler := mvc.NewGinHandler(r)
	ginHandler.RegisterGreeter(&postGreeter)
	ginHandler.RegisterGreeter(&getGreeter)

	if env != consts.RunEnvNull {
		log.Info("开启pprof监控")
		ginpprof.Wrap(r)
	}

	log.Infof("POL_ENV:%s, connect to http://%s:%d/ for %s service", env, host, port, applicationName)
	r.Run(":" + strconv.Itoa(port))
}
