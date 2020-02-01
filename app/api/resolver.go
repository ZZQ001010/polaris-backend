package api

import (
	"context"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/polaris-backend/app/generated"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/99designs/gqlgen/graphql"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

var log = logger.GetDefaultLogger()

type Resolver struct{}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

var DefaultDirective = generated.DirectiveRoot{
	HasRole: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		graphqlContext := graphql.GetResolverContext(ctx)
		if graphqlContext == nil {
			return nil, errs.BuildSystemErrorInfo(errs.SystemError)
		}
		//fmt.Println(graphqlContext.Field.Name)
		//fmt.Println(graphqlContext.Args["input"])
		//等待统一接口入参，做权限判断
		return next(ctx)
	},
}
