package coding_carefree

import (
	"context"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Register(ctx context.Context, input RegisterInVo) (*Void, error) {
	return RegisterService(ctx, input)
}
func (r *mutationResolver) Login(ctx context.Context, input LoginInVo) (*LoginOutVo, error) {
	return LoginService(ctx, input)
}
func (r *mutationResolver) SendMail(ctx context.Context, input SendMailInVo) (*Void, error) {
	return SendMailService(ctx, input)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]*User, error) {
	panic("not implemented")
}
