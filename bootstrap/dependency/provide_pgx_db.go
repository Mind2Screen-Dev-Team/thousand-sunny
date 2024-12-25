package dependency

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
)

func ProvidePGxDB(c config.Cfg, lc fx.Lifecycle) (*pgx.Conn, error) {

	var (
		ctx = context.Background()
		dsn = fmt.Sprintf(
			"application_name=%s host=%s port=%d dbname=%s user=%s password=%s TimeZone=%s sslmode=%s connect_timeout=%d",

			// # APP ID
			c.App.Name,

			// # Connection
			c.DB.Address,
			c.DB.Port,
			c.DB.DBName,

			// # Credentials
			c.DB.Credential.Username,
			c.DB.Credential.Password,

			// # Options
			c.DB.Options.Timezone,
			c.DB.Options.Sslmode,
			c.DB.Options.ConnectionTimeout,
		)
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				return conn.Close(ctx)
			},
		},
	)

	return conn, nil
}

func PingPGxDB(conn *pgx.Conn) error {
	if err := conn.Ping(context.Background()); err != nil {
		return err
	}

	return nil
}
