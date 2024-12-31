package service_api

import (
	"context"

	service_attr "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/service/attr"
)

type UserServiceAPI interface {
	Create(ctx context.Context, user service_attr.UserCreate) error
	Delete(ctx context.Context, id string) error
	Read(ctx context.Context, id string) (*service_attr.User, error)
	ReadAll(ctx context.Context, limit int, offset int) ([]service_attr.User, error)
	Update(ctx context.Context, user service_attr.User) error
}
