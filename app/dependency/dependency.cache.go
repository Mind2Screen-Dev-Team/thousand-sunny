package dependency

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func ProvideRedis(c config.Cfg, lc fx.Lifecycle) *redis.Client {
	var (
		cfg      = c.Cache["redis"]
		dbIdx, _ = strconv.Atoi(cfg.DBName)
		rdb      = redis.NewClient(
			&redis.Options{
				ClientName: c.App.Name,
				Addr:       fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
				DB:         dbIdx,
				Username:   cfg.Credential.Username,
				Password:   cfg.Credential.Password,
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
