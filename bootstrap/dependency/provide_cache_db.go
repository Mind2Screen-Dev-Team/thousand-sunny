package dependency

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func ProvideRedisDB(c config.Cfg, lc fx.Lifecycle) *redis.Client {
	var (
		dbIdx, _ = strconv.Atoi(c.Cache.DBName)
		rdb      = redis.NewClient(
			&redis.Options{
				ClientName: c.App.Name,
				Addr:       fmt.Sprintf("%s:%d", c.Cache.Address, c.Cache.Port),
				DB:         dbIdx,
				Username:   c.Cache.Credential.Username,
				Password:   c.Cache.Credential.Password,
			},
		)
	)

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				return rdb.Close()
			},
		},
	)

	return rdb
}

func PingRedisDB(rdb *redis.Client) error {
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return err
	}
	return nil
}
