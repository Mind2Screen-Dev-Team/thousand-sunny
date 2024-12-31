package service_impl

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/rs/xid"
	"go.uber.org/fx"

	repo_api "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repo/api"
	repo_attr "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repo/attr"
	service_api "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/api"
	service_attr "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/attr"
)

type (
	UserServiceParamFx struct {
		fx.In

		RDB      *redis.Client
		UserRepo repo_api.UserRepoAPI `optional:"false"`
	}

	UserImplServiceFx struct {
		p UserServiceParamFx
	}
)

func NewUserCURDService(p UserServiceParamFx) (service_api.UserServiceAPI, error) {
	if p.UserRepo == nil {
		return nil, errors.New("failed to load user curd repo")
	}

	return &UserImplServiceFx{p}, nil
}

func (r *UserImplServiceFx) Create(ctx context.Context, user service_attr.UserCreate) error {
	return r.p.UserRepo.Create(ctx, repo_attr.User{
		ID:   xid.New().String(),
		Name: user.Name,
		Age:  user.Age,
	})
}

func (r *UserImplServiceFx) ReadAll(ctx context.Context, limit int, offset int) ([]service_attr.User, error) {
	res, err := r.p.UserRepo.ReadAll(ctx, limit, offset)
	if err != nil {
		return make([]service_attr.User, 0), err
	}

	result := make([]service_attr.User, len(res))
	for i := range res {
		result[i] = service_attr.User{
			ID:   res[i].ID,
			Name: res[i].Name,
			Age:  res[i].Age,
		}
	}

	return result, nil
}

func (r *UserImplServiceFx) Read(ctx context.Context, id string) (*service_attr.User, error) {
	res, err := r.p.UserRepo.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	return &service_attr.User{
		ID:   res.ID,
		Age:  res.Age,
		Name: res.Name,
	}, nil
}

func (r *UserImplServiceFx) Update(ctx context.Context, user service_attr.User) error {
	return r.p.UserRepo.Update(ctx, repo_attr.User{
		ID:   user.ID,
		Name: user.Name,
		Age:  user.Age,
	})
}

func (r *UserImplServiceFx) Delete(ctx context.Context, id string) error {
	return r.p.UserRepo.Delete(ctx, id)
}
