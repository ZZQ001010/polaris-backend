package mid

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/threadlocal"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	consts1 "github.com/galaxy-book/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"time"
)

//var log = logger.GetDefaultLogger()

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,PM-TOKEN,PM-ORG,PM-PRO,PM-ENV,PM-PLAT,PM-VER,PM-TRACE-ID")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length,  Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,PM-TOKEN,PM-ORG,PM-PRO,PM-ENV,PM-PLAT,PM-VER,PM-TRACE-ID")

		c.Header("Access-Control-Allow-Credentials", "true")

		//表示隔24个小时才发起预检请求。也就是说，发送两次请求
		c.Header("Access-Control-Max-Age", "86400")

		//fmt.Println(c.Request.Context().Value(consts.TraceIdKey))
		//fmt.Println(c.Request.Context().Value(consts.HttpContextKey).(domains.HttpContext).StartTime)

		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(consts.AppHeaderTokenName)
		//兼容
		if token == "" {
			token = c.GetHeader("Token")
		}
		if token != "" {
			c.Set(consts.AppHeaderTokenName, token)
		}

		// 处理请求
		c.Next()
	}
}

func SentryMiddleware(applicationName string, env string, dsn string) gin.HandlerFunc {
	connectSentryErr := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		ServerName:  applicationName,
		Environment: env,
	})
	return func(c *gin.Context) {
		//捕获异常
		defer func() {
			if p := recover(); p != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				stackInfo := fmt.Sprintf("%s", buf[:n])

				errMsg := ""
				if err, ok := p.(error); ok {
					errMsg = err.Error()
				} else {
					errMsg = fmt.Sprint(p)
				}
				event := sentry.NewEvent()
				event.EventID = sentry.EventID(errMsg)
				event.Message = errMsg + "\n" + stackInfo
				event.Environment = env
				event.ServerName = applicationName
				event.Tags[consts1.TraceIdKey] = threadlocal.GetTraceId()
				if connectSentryErr == nil {
					sentry.CaptureEvent(event)
					sentry.Flush(time.Second * 5)
				} else {
					fmt.Println(connectSentryErr)
				}
				fmt.Println(event.Message)

				respObj := vo.VoidErr{Err: vo.NewErr(errs.BuildSystemErrorInfoWithMessage(errs.ServerError, errMsg))}
				c.String(200, json.ToJsonIgnoreError(respObj))
			}
		}()

		// 处理请求
		c.Next()
	}
}

//func OpenTracingMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		spanCtx, _ := tracing.SpanFromContext(c)
//		if spanCtx != nil {
//			span := opentracing.StartSpan(c.Request.RequestURI, opentracing.ChildOf(spanCtx))
//			ctx = opentracing.ContextWithSpan(ctx, span)
//			defer span.Finish()
//		}
//		// 处理请求
//		c.Next()
//	}
//}
