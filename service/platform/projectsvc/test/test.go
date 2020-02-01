package test

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/consts"
	"github.com/gin-gonic/gin"
)

var env = ""

const BaseConfigPath = "./../../../../config"
const SelfConfigPath = "./../config"

const User1001 = int64(1001)
const Org1001 = int64(1001)

const IssueId100 = int64(100)
const IssueId1083 = int64(1083)

func StartUp(f func(ctx context.Context)) func() {
	return func() {
		env = "local"
		//配置文件
		err := config.LoadEnvConfig(BaseConfigPath, "application.common", env)

		if err != nil {
			fmt.Printf("err:%s\n", err)
		}

		err = config.LoadEnvConfig(SelfConfigPath, "application", env)

		if err != nil {
			fmt.Printf("err:%s\n", err)
		}

		//添加token操作
		ginCtx := gin.Context{}
		ginCtx.Set(consts.AppHeaderTokenName, "abc")
		//获得一个顶级上下文
		ctx := context.Background()
		//返回父上下文  在父上下文中设置 key/value  这边是 GinContextKey:ginCtx
		ctx = context.WithValue(ctx, "GinContextKey", &ginCtx)

		//丢进来的方法立刻执行
		f(ctx)
	}
}

func StartUpWithUserInfo(f func(userId, orgId int64)) func() {
	return func() {
		env = "unittest"
		//配置文件
		err := config.LoadEnvConfig(BaseConfigPath, "application.common", env)

		if err != nil {
			fmt.Printf("err:%s\n", err)
		}

		err = config.LoadEnvConfig(SelfConfigPath, "application", env)

		if err != nil {
			fmt.Printf("err:%s\n", err)
		}
		f(User1001, Org1001)
	}
}
