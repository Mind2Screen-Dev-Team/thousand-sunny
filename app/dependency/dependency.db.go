package dependency

import (
	"context"
	"fmt"
	"time"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func ProvidePostgres(c config.Cfg, lc fx.Lifecycle) (*pgxpool.Pool, error) {
	var (
		cfg = c.DB["postgres"]
		ctx = context.Background()
		dsn = fmt.Sprintf(
			"application_name=%s host=%s port=%d dbname=%s user=%s password=%s TimeZone=%s sslmode=%s connect_timeout=%d",

			// # APP ID
			c.App.Name,

			// # Connection
			cfg.Address,
			cfg.Port,
			cfg.DBName,

			// # Credentials
			cfg.Credential.Username,
			cfg.Credential.Password,

			// # Options
			cfg.Options.Timezone,
			cfg.Options.Sslmode,
			cfg.Options.ConnectionTimeout,
		)
	)

	// Create a pgxpool.Config from the connection string
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// Customize pool settings
	poolCfg.MaxConns = int32(cfg.Options.MaxOpenConnection)
	poolCfg.MinConns = int32(cfg.Options.MaxIdleConnection)
	poolCfg.MaxConnLifetime = time.Duration(cfg.Options.MaxConnectionLifetime) * time.Second

	db, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	hook := fx.Hook{OnStop: func(ctx context.Context) error {
		db.Close()
		return nil
	}}
	lc.Append(hook)

	return db, nil
}

func InvokePostgres(conn *pgxpool.Pool) error {
	if err := conn.Ping(context.Background()); err != nil {
		return err
	}

	return nil
}
