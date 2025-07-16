package impl

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	repo_api "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repository/user/api"
	repo_attr "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repository/user/attr"

	svc_api "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/user/api"
	svc_attr "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/user/attr"
)

type (
	ExampleUserServiceParamFx struct {
		fx.In

		RDB             *redis.Client
		ExampleUserRepo repo_api.ExampleUserRepoAPI `optional:"false"`
	}

	ExampleUserImplServiceFx struct {
		p ExampleUserServiceParamFx
	}
)

func NewExampleUserService(p ExampleUserServiceParamFx) (svc_api.ExampleUserServiceAPI, error) {
	if p.ExampleUserRepo == nil {
		return nil, errors.New("failed to load user curd repo")
	}
	return &ExampleUserImplServiceFx{p}, nil
}

func (s *ExampleUserImplServiceFx) Create(ctx context.Context, user svc_attr.ExampleUserCreate) (*svc_attr.ExampleUser, error) {
	d := repo_attr.ExampleUser{
		ID:        uuid.Must(uuid.NewV7()),
		Name:      user.Name,
		Age:       user.Age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r, err := s.p.ExampleUserRepo.Create(ctx, d)
	if err != nil {
		return nil, err
	}

	return &svc_attr.ExampleUser{
		ID:        r.ID,
		Name:      r.Name,
		Age:       r.Age,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func (s *ExampleUserImplServiceFx) Read(ctx context.Context, id string) (*svc_attr.ExampleUser, error) {
	r, err := s.p.ExampleUserRepo.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	return &svc_attr.ExampleUser{
		ID:        r.ID,
		Age:       r.Age,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func (s *ExampleUserImplServiceFx) ReadAll(ctx context.Context, limit int, offset int) ([]svc_attr.ExampleUser, error) {
	r, err := s.p.ExampleUserRepo.ReadAll(ctx, limit, offset)
	if err != nil {
		return make([]svc_attr.ExampleUser, 0), err
	}

	result := make([]svc_attr.ExampleUser, len(r))
	for i := range r {
		result[i] = svc_attr.ExampleUser{
			ID:        r[i].ID,
			Name:      r[i].Name,
			Age:       r[i].Age,
			CreatedAt: r[i].CreatedAt,
			UpdatedAt: r[i].UpdatedAt,
		}
	}

	return result, nil
}

func (s *ExampleUserImplServiceFx) Update(ctx context.Context, user svc_attr.ExampleUser) (*svc_attr.ExampleUser, error) {
	d := repo_attr.ExampleUser{
		ID:        user.ID,
		Name:      user.Name,
		Age:       user.Age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r, err := s.p.ExampleUserRepo.Update(ctx, d)
	if err != nil {
		return nil, err
	}

	return &svc_attr.ExampleUser{
		ID:        r.ID,
		Name:      r.Name,
		Age:       r.Age,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func (s *ExampleUserImplServiceFx) Delete(ctx context.Context, id string) (*svc_attr.ExampleUser, error) {
	r, err := s.p.ExampleUserRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return &svc_attr.ExampleUser{
		ID:        r.ID,
		Name:      r.Name,
		Age:       r.Age,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}
