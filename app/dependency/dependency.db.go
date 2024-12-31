package dependency

import (
	"context"
	"fmt"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/jackc/pgx/v5"
	"go.uber.org/fx"
)

func ProvidePostgres(c config.Cfg, lc fx.Lifecycle) (*pgx.Conn, error) {
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

	db, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	hook := fx.Hook{OnStop: func(ctx context.Context) error {
		return db.Close(ctx)
	}}
	lc.Append(hook)

	return db, nil
}

func InvokePostgres(conn *pgx.Conn) error {
	if err := conn.Ping(context.Background()); err != nil {
		return err
	}

	return nil
}
