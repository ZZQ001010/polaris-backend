package main

import (
	"flag"
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/network"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/app/server/handler"
	"github.com/galaxy-book/polaris-backend/common/core/buildinfo"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/util/http"
	"github.com/galaxy-book/polaris-backend/common/extra/gin/mid"
	"github.com/galaxy-book/polaris-backend/common/extra/init/db"
	"github.com/DeanThompson/ginpprof"
	"github.com/ainilili/go2sky"
	skyGin "github.com/ainilili/go2sky/plugins/gin"
	"github.com/ainilili/go2sky/reporter"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"os"
	"runtime"
	"strconv"
)

var (
	env = ""
	log = logger.GetDefaultLogger()
)

const BaseConfigPath = "./../config"
const SelfConfigPath = "./config"

func init() {
	env = os.Getenv(consts.RunEnvKey)
	if "" == env {
		env = consts.RunEnvLocal
	}
	//配置
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

	serverConfig := config.GetServerConfig()
	port := strconv.Itoa(serverConfig.Port)

	applicationName := config.GetApplication().Name

	// 单库模式才会执行
	if (consts.AppRunmodePrivateSingleDb == config.GetApplication().RunMode) || (consts.AppRunmodeSingle == config.GetApplication().RunMode) {
		mysqlConfig := config.GetMysqlConfig()
		initErr := db.DbMigrations(env, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Usr, mysqlConfig.Pwd, mysqlConfig.Database)
		if initErr != nil {
			panic(" init db fail....")
		}
	}

	r := gin.New()

	//sky walking
	if config.GetSkyWalkingConfig() != nil {
		skyReporter, err := reporter.NewGRPCReporter(config.GetSkyWalkingConfig().ReportAddress)
		if err != nil {
			log.Error("new reporter error " + strs.ObjectToString(err))
		}
		defer skyReporter.Close()
		tracer, err := go2sky.NewTracer(applicationName+"-"+env, go2sky.WithReporter(skyReporter))
		if err != nil {
			log.Error("create tracer error " + strs.ObjectToString(err))
		}
		tracer.WaitUntilRegister()
		r.Use(skyGin.Middleware(r, tracer))
		http.InitSkyWalking(tracer)
		log.Info("skywaking uping!")
	}

	sentryConfig := config.GetSentryConfig()
	sentryDsn := ""
	if sentryConfig != nil {
		sentryDsn = sentryConfig.Dsn
	}
	r.Use(mid.SentryMiddleware(applicationName, env, sentryDsn))
	r.Use(mid.StartTrace())
	r.Use(mid.GinContextToContextMiddleware())
	r.Use(mid.CorsMiddleware())
	r.Use(mid.AuthMiddleware())

	captcha.SetCustomStore(&handler.RedisCache{})

	//1.获取验证码
	r.POST("/api/task/captcha", handler.CaptchaGetHandler())

	//2.获取验证码图片
	r.GET("/api/task/captcha/:captchaId", handler.CaptchaShowHandler())

	r.POST("/task", handler.GraphqlHandler())
	r.GET("/task/health", handler.HeartbeatHandler())
	r.GET("/", handler.PlaygroundHandler())
	r.POST("/callback", handler.DingTalkCallBackHandler())

	r.POST("/api/task", handler.GraphqlHandler())
	r.GET("/api/task/health", handler.HeartbeatHandler())
	r.GET("/admin/pg", handler.PlaygroundHandler())
	r.POST("/api/callback", handler.DingTalkCallBackHandler())

	r.POST("/api/task/upload", handler.FileUploadHandler())
	r.GET("/api/task/read/*path", handler.FileReadHandler())
	r.POST("/api/task/importData", handler.ImportDataHandler())

	//测试apk分组接口
	apk := r.Group("/api/mct/apk")
	{
		apk.POST("/uploadApkInfo", handler.UploadApkInfoHandler())
		apk.GET("/getAllApk", handler.GetAllApkHandler())
		apk.POST("/deleteApk", handler.DeleteApkHandler())
	}
	//测试设备分组接口
	devices := r.Group("/api/mct/devices")
	{
		devices.GET("/getAllDevices", handler.GetAllDevicesHandler())
		devices.GET("/getDeviceFiltrate", handler.GetDeviceFiltrateHandler())
		devices.GET("/getAllDevicesStatus", handler.GetAllDevicesStatusHandler())
		devices.POST("/startCompat", handler.StartCompatHandler())
	}
	//获取报告
	report := r.Group("/api/mct/report")
	{
		report.POST("/getReport", handler.GetReportHandler())
		report.POST("/deleteReport", handler.DeleteReportHandler())
		report.POST("/getReportApkInfo", handler.GetReportApkInfoHandler())
		report.POST("/getReportDetailOverView", handler.GetReportDetailOverViewHandler())
		report.POST("/getReportDetailSingle", handler.GetReportDetailSingleHandler())
		report.POST("/getReportDetailError", handler.GetReportDetailErrorHandler())
		report.POST("/getReportDetailPerformance", handler.GetReportDetailPerformanceHandler())
	}

	//fmt.Println("connect to http://" + network.GetIntranetIp() + ":" + port + "/ for GraphQL playground")

	if env != consts.RunEnvNull {
		log.Info("开启pprof监控")
		ginpprof.Wrap(r)
	}

	log.Infof("POL_ENV:%s, connect to http://%s:%s/ for GraphQL playground", env, network.GetIntranetIp(), port)
	r.Run(":" + port)
}