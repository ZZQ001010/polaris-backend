package main

import (
	"fmt"
	"github.com/galaxy-book/polaris-backend/common/core/buildinfo"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	consts2 "github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/extra/gin/mid"
	"github.com/galaxy-book/polaris-backend/common/extra/gin/mvc"
	"github.com/galaxy-book/polaris-backend/schedule/api"
	"github.com/galaxy-book/polaris-backend/schedule/schedule"
	"github.com/DeanThompson/ginpprof"
	"github.com/Jeffail/tunny"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/network"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"os"
	"runtime"
	"strconv"
)

var log = logger.GetDefaultLogger()
var env = ""

const BaseConfigPath = "./../config"
const SelfConfigPath = "./config"

func init() {
	env = os.Getenv(consts.RunEnvKey)
	if "" == env {
		env = consts.RunEnvLocal
	}
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

const (
	//每天凌晨执行
	ZeroClock1Minute = "5 0 0 * * ?"
)

func main() {
	// 打印程序信息
	log.Info(buildinfo.StringifySingleLine())
	fmt.Println(buildinfo.StringifyMultiLine())

	port := config.GetServerConfig().Port
	host := config.GetServerConfig().Host

	msg := json.ToJsonIgnoreError(config.GetConfig())

	log.Info("config配置:" + msg)

	applicationName := config.GetApplication().Name

	//定时任务通用协程池
	pool := tunny.NewFunc(300, func(payload interface{}) interface{} {
		fn := payload.(func() error)
		return fn()
	})
	defer pool.Close()

	c := cron.New()

	//每天凌晨，迭代燃尽图
	_ = c.AddFunc(ZeroClock1Minute, func() {
		log.Info("开始迭代燃尽图统计")
		schedule.StatisticIterationBurnDownChart(*pool)
	})
	//每天凌晨，项目燃尽图
	_ = c.AddFunc(ZeroClock1Minute, func() {
		log.Info("开始项目燃尽图统计")
		schedule.StatisticProjectIssueBurnDownChart(*pool)
	})


	c.Start()
	log.Info("定时任务启动成功！")

	r := gin.New()
	r.Use(mid.StartTrace())
	r.Use(mid.GinContextToContextMiddleware())
	r.Use(mid.CorsMiddleware())
	r.Use(mid.AuthMiddleware())
	version := ""
	getGreeter := api.GetGreeter{Greeter: mvc.NewGetGreeter(applicationName, host, port, version)}

	ginHandler := mvc.NewGinHandler(r)
	ginHandler.RegisterGreeter(&getGreeter)

	if env != consts2.RunEnvNull {
		log.Info("开启pprof监控")
		ginpprof.Wrap(r)
	}

	log.Infof("POL_ENV:%s, connect to http://%s:%d/ for %s service", env, network.GetIntranetIp(), port, applicationName)
	r.Run(":" + strconv.Itoa(port))
}
