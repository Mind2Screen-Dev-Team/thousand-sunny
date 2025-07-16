package api

import (
	"context"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repository/user/attr"
)

type ExampleUserRepoAPI interface {
	Create(ctx context.Context, user attr.ExampleUser) (*attr.ExampleUser, error)
	Read(ctx context.Context, id string) (*attr.ExampleUser, error)
	ReadAll(ctx context.Context, limit int, offset int) ([]attr.ExampleUser, error)
	Update(ctx context.Context, user attr.ExampleUser) (*attr.ExampleUser, error)
	Delete(ctx context.Context, id string) (*attr.ExampleUser, error)
}
