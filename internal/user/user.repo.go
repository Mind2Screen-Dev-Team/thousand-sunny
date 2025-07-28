package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type ExampleUserRepoAPI interface {
	Create(ctx context.Context, user ExampleUser) (*ExampleUser, error)
	Read(ctx context.Context, id string) (*ExampleUser, error)
	ReadAll(ctx context.Context, limit int, offset int) ([]ExampleUser, error)
	Update(ctx context.Context, user ExampleUser) (*ExampleUser, error)
	Delete(ctx context.Context, id string) (*ExampleUser, error)
}

type (
	ExampleUserRepoParamFx struct {
		fx.In

		RDB *redis.Client
	}

	ExampleUserImplRepoFx struct {
		p ExampleUserRepoParamFx
	}
)

func NewRepo(p ExampleUserRepoParamFx) (ExampleUserRepoAPI, error) {
	if p.RDB == nil {
		return nil, errors.New("field 'RDB' with type '*redis.Client' is not provided")
	}

	return &ExampleUserImplRepoFx{p}, nil
}

func (r *ExampleUserImplRepoFx) Create(ctx context.Context, user ExampleUser) (*ExampleUser, error) {
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

func (r *ExampleUserImplRepoFx) Read(ctx context.Context, id string) (*ExampleUser, error) {
	key := fmt.Sprintf("user:%s", id)
	data, err := r.p.RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("user with ID %s not found", id)
	} else if err != nil {
		return nil, err
	}

	var user ExampleUser
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ExampleUserImplRepoFx) ReadAll(ctx context.Context, limit int, offset int) ([]ExampleUser, error) {
	keys, err := r.p.RDB.Keys(ctx, "user:*").Result()
	if err != nil {
		return nil, err
	}

	// Sort keys to ensure consistent pagination
	sort.Strings(keys)

	start := offset
	if start > len(keys) {
		return []ExampleUser{}, nil
	}

	end := offset + limit
	if limit <= 0 || end > len(keys) {
		end = len(keys)
	}

	var users []ExampleUser
	for _, key := range keys[start:end] {
		data, err := r.p.RDB.Get(ctx, key).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return nil, err
		}

		var user ExampleUser
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *ExampleUserImplRepoFx) Update(ctx context.Context, user ExampleUser) (*ExampleUser, error) {
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

func (r *ExampleUserImplRepoFx) Delete(ctx context.Context, id string) (*ExampleUser, error) {
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
