package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type ExampleUserServiceAPI interface {
	Create(ctx context.Context, user ExampleUser) (*ExampleUser, error)
	Read(ctx context.Context, id string) (*ExampleUser, error)
	ReadAll(ctx context.Context, limit int, offset int) ([]ExampleUser, error)
	Update(ctx context.Context, user ExampleUser) (*ExampleUser, error)
	Delete(ctx context.Context, id string) (*ExampleUser, error)
}

type (
	ExampleUserServiceParamFx struct {
		fx.In

		RDB             *redis.Client
		ExampleUserRepo ExampleUserRepoAPI `optional:"false"`
	}

	ExampleUserImplServiceFx struct {
		p ExampleUserServiceParamFx
	}
)

func NewService(p ExampleUserServiceParamFx) (ExampleUserServiceAPI, error) {
	if p.ExampleUserRepo == nil {
		return nil, errors.New("failed to load user curd repo")
	}
	return &ExampleUserImplServiceFx{p}, nil
}

func (s *ExampleUserImplServiceFx) Create(ctx context.Context, user ExampleUser) (*ExampleUser, error) {
	d := ExampleUser{
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

	return &ExampleUser{
		ID:        r.ID,
		Name:      r.Name,
		Age:       r.Age,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func (s *ExampleUserImplServiceFx) Read(ctx context.Context, id string) (*ExampleUser, error) {
	r, err := s.p.ExampleUserRepo.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ExampleUser{
		ID:        r.ID,
		Age:       r.Age,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func (s *ExampleUserImplServiceFx) ReadAll(ctx context.Context, limit int, offset int) ([]ExampleUser, error) {
	r, err := s.p.ExampleUserRepo.ReadAll(ctx, limit, offset)
	if err != nil {
		return make([]ExampleUser, 0), err
	}

	result := make([]ExampleUser, len(r))
	for i := range r {
		result[i] = ExampleUser{
			ID:        r[i].ID,
			Name:      r[i].Name,
			Age:       r[i].Age,
			CreatedAt: r[i].CreatedAt,
			UpdatedAt: r[i].UpdatedAt,
		}
	}

	return result, nil
}

func (s *ExampleUserImplServiceFx) Update(ctx context.Context, user ExampleUser) (*ExampleUser, error) {
	d := ExampleUser{
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

	return &ExampleUser{
		ID:        r.ID,
		Name:      r.Name,
		Age:       r.Age,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func (s *ExampleUserImplServiceFx) Delete(ctx context.Context, id string) (*ExampleUser, error) {
	r, err := s.p.ExampleUserRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ExampleUser{
		ID:        r.ID,
		Name:      r.Name,
		Age:       r.Age,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}
