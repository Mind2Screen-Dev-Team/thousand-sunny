package repo_impl

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	repo_api "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repo/api"
	repo_attr "github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repo/attr"
)

type (
	UserRepoParamFx struct {
		fx.In

		RDB *redis.Client
	}

	UserImplRepoFx struct {
		p UserRepoParamFx
	}
)

func NewUserCURDRepo(p UserRepoParamFx) repo_api.UserRepoAPI {
	return &UserImplRepoFx{p}
}

func (r *UserImplRepoFx) Create(ctx context.Context, user repo_attr.User) error {
	key := fmt.Sprintf("user:%s", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.p.RDB.Set(ctx, key, data, 0).Err()
}

func (r *UserImplRepoFx) ReadAll(ctx context.Context, limit int, offset int) ([]repo_attr.User, error) {
	keys, err := r.p.RDB.Keys(ctx, "user:*").Result()
	if err != nil {
		return nil, err
	}

	// Sort keys to ensure consistent pagination
	sort.Strings(keys)

	start := offset
	if start > len(keys) {
		return []repo_attr.User{}, nil
	}

	end := offset + limit
	if limit <= 0 || end > len(keys) {
		end = len(keys)
	}

	var users []repo_attr.User
	for _, key := range keys[start:end] {
		data, err := r.p.RDB.Get(ctx, key).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return nil, err
		}

		var user repo_attr.User
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserImplRepoFx) Read(ctx context.Context, id string) (*repo_attr.User, error) {
	key := fmt.Sprintf("user:%s", id)
	data, err := r.p.RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("user with ID %s not found", id)
	} else if err != nil {
		return nil, err
	}

	var user repo_attr.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserImplRepoFx) Update(ctx context.Context, user repo_attr.User) error {
	key := fmt.Sprintf("user:%s", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.p.RDB.Set(ctx, key, data, 0).Err()
}

func (r *UserImplRepoFx) Delete(ctx context.Context, id string) error {
	key := fmt.Sprintf("user:%s", id)
	return r.p.RDB.Del(ctx, key).Err()
}
