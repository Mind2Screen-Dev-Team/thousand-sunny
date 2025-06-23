package dependency

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func ProvideRedis(c config.Cfg, s config.Server, lc fx.Lifecycle) *redis.Client {

	durationFn := func(v int) time.Duration {
		if v > 0 {
			return time.Duration(v) * time.Second
		}
		return time.Duration(v)
	}

	var (
		cfg = c.Cache["redis"]

		dialTimeout  = cfg.Options.DialTimeout
		readTimeout  = cfg.Options.ReadTimeout
		writeTimeout = cfg.Options.WriteTimeout

		dbIdx, _ = strconv.Atoi(cfg.DBName)
		rdb      = redis.NewClient(
			&redis.Options{
				ClientName:   fmt.Sprintf("%s/%s", c.App.Project, s.Name),
				Addr:         fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
				DB:           dbIdx,
				Username:     cfg.Credential.Username,
				Password:     cfg.Credential.Password,
				DialTimeout:  durationFn(dialTimeout),
				ReadTimeout:  durationFn(readTimeout),
				WriteTimeout: durationFn(writeTimeout),
			},
		)
	)

	hook := fx.Hook{OnStop: func(ctx context.Context) error {
		return rdb.Close()
	}}
	lc.Append(hook)

	return rdb
}

func InvokeRedis(rdb *redis.Client) error {
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return err
	}
	return nil
}

func ProvideRedisLock(c *redis.Client) *redislock.Client {
	return redislock.New(c)
}
