package api

import (
	"context"

	repo_attr "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repository/user/attr"
)

type UserRepoAPI interface {
	Create(ctx context.Context, user repo_attr.User) error
	Delete(ctx context.Context, id string) error
	Read(ctx context.Context, id string) (*repo_attr.User, error)
	ReadAll(ctx context.Context, limit int, offset int) ([]repo_attr.User, error)
	Update(ctx context.Context, user repo_attr.User) error
}
