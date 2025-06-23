package dependency

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
)

func ProvidePostgres(c config.Cfg, s config.Server, d *xlog.DebugLogger, lc fx.Lifecycle) (*pgxpool.Pool, error) {
	var (
		cfg = c.DB["postgres"]
		ctx = context.Background()
		dsn = fmt.Sprintf(
			"application_name=%s host=%s port=%d dbname=%s user=%s password=%s TimeZone=%s sslmode=%s connect_timeout=%d",

			// # APP ID
			fmt.Sprintf("pgx-%s-%s", c.App.Project, s.Name),

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
	poolCfg.ConnConfig.Tracer = &xlog.PgxLogger{Log: xlog.NewLogger(d.Logger)}

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

// This Paramater:
//   - pool *pgxpool.Pool just for load the provider.
func ProvidePostgresSQLDB(c config.Cfg, s config.Server, pool *pgxpool.Pool, lc fx.Lifecycle) (*sql.DB, error) {
	cfg := c.DB["postgres"]
	dsn := _BuildSafePostgresDSN(c, s, "sqldb")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Optionally configure connection pool settings
	db.SetMaxOpenConns(cfg.Options.MaxOpenConnection)
	db.SetMaxIdleConns(cfg.Options.MaxIdleConnection)
	db.SetConnMaxLifetime(time.Duration(cfg.Options.MaxConnectionLifetime) * time.Second)

	// Close on shutdown
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

func _BuildSafePostgresDSN(c config.Cfg, s config.Server, app_prefix string) string {
	cfg := c.DB["postgres"]

	// Create the URL struct using the URL-encoded password
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Credential.Username, cfg.Credential.Password),
		Host:   fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		Path:   cfg.DBName,
	}

	// Add the query parameters
	q := u.Query()
	q.Set("application_name", fmt.Sprintf("%s-%s-%s", app_prefix, c.App.Project, s.Name))
	q.Set("TimeZone", cfg.Options.Timezone)
	q.Set("sslmode", cfg.Options.Sslmode)
	q.Set("connect_timeout", fmt.Sprintf("%d", cfg.Options.ConnectionTimeout))
	u.RawQuery = q.Encode()

	// Return the fully constructed DSN
	return u.String()
}
