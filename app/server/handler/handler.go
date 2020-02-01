package handler

import (
	"github.com/galaxy-book/common/core/consts"
	"github.com/galaxy-book/common/core/model"
	"github.com/galaxy-book/common/core/threadlocal"
	"github.com/galaxy-book/polaris-backend/app/api"

	"github.com/galaxy-book/polaris-backend/app/generated"
	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
	"github.com/jtolds/gls"
)

// Defining the Graphql handler
func GraphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.GraphQL(generated.NewExecutableSchema(generated.Config{
		Resolvers:  &api.Resolver{},
		Directives: api.DefaultDirective,
	}))

	return func(c *gin.Context) {
		//token := c.GetHeader(consts.APP_HEADER_TOKEN_NAME)

		httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
		threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
			h.ServeHTTP(c.Writer, c.Request)
		})
	}
}

// Defining the Playground handler
func PlaygroundHandler() gin.HandlerFunc {
	h := handler.Playground("GraphQL", "/query")
	return func(c *gin.Context) {
		httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
		threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
			h.ServeHTTP(c.Writer, c.Request)
		})
	}
}

func DingTalkCallBackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpContext := c.Request.Context().Value(consts.HttpContextKey).(model.HttpContext)
		threadlocal.Mgr.SetValues(gls.Values{consts.HttpContextKey: httpContext, consts.TraceIdKey: httpContext.TraceId}, func() {
			//dingtalk.DingTalkCallbackHandler(c.Writer, c.Request)
		})
	}
}

func HeartbeatHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "ok")
	}
}
