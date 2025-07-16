package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repository/user/api"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/repository/user/attr"
)

type (
	ExampleUserRepoParamFx struct {
		fx.In

		RDB *redis.Client
	}

	ExampleUserImplRepoFx struct {
		p ExampleUserRepoParamFx
	}
)

func NewExampleUserRepo(p ExampleUserRepoParamFx) (api.ExampleUserRepoAPI, error) {
	if p.RDB == nil {
		return nil, errors.New("field 'RDB' with type '*redis.Client' is not provided")
	}

	return &ExampleUserImplRepoFx{p}, nil
}

func (r *ExampleUserImplRepoFx) Create(ctx context.Context, user attr.ExampleUser) (*attr.ExampleUser, error) {
	key := fmt.Sprintf("user:%s", user.ID.String())
	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}

	return &user, r.p.RDB.Set(ctx, key, data, 0).Err()
}

func (r *ExampleUserImplRepoFx) Read(ctx context.Context, id string) (*attr.ExampleUser, error) {
	key := fmt.Sprintf("user:%s", id)
	data, err := r.p.RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("user with ID %s not found", id)
	} else if err != nil {
		return nil, err
	}

	var user attr.ExampleUser
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ExampleUserImplRepoFx) ReadAll(ctx context.Context, limit int, offset int) ([]attr.ExampleUser, error) {
	keys, err := r.p.RDB.Keys(ctx, "user:*").Result()
	if err != nil {
		return nil, err
	}

	// Sort keys to ensure consistent pagination
	sort.Strings(keys)

	start := offset
	if start > len(keys) {
		return []attr.ExampleUser{}, nil
	}

	end := offset + limit
	if limit <= 0 || end > len(keys) {
		end = len(keys)
	}

	var users []attr.ExampleUser
	for _, key := range keys[start:end] {
		data, err := r.p.RDB.Get(ctx, key).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return nil, err
		}

		var user attr.ExampleUser
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *ExampleUserImplRepoFx) Update(ctx context.Context, user attr.ExampleUser) (*attr.ExampleUser, error) {
	d, err := r.Read(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	user.CreatedAt = d.CreatedAt
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}

	key := fmt.Sprintf("user:%s", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", "repository failed to marshal", err)
	}

	if err := r.p.RDB.Set(ctx, key, data, 0).Err(); err != nil {
		return nil, fmt.Errorf("%s:%w", "repository failed to save to redis", err)
	}

	return &user, nil
}

func (r *ExampleUserImplRepoFx) Delete(ctx context.Context, id string) (*attr.ExampleUser, error) {
	d, err := r.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("user:%s", id)
	if err := r.p.RDB.Del(ctx, key).Err(); err != nil {
		return nil, err
	}

	return d, err
}
