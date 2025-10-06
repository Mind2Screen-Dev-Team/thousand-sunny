package dependency

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"

	"go.uber.org/fx"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/gen/gorm/query"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
)

func ProvideGormQuery(db *gorm.DB) *query.Query {
	return query.Use(db)
}

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

	// Exclude migration-command and gorm-generator-command
	if !slices.Contains([]string{"migration-command", "gorm-generator-command"}, s.Name) {
		poolCfg.ConnConfig.Tracer = &xlog.PgxLogger{Log: xlog.NewLogger(d.Logger)}
	}

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

func ProvidePostgresSQLDB(c config.Cfg, s config.Server, pool *pgxpool.Pool, lc fx.Lifecycle) (*sql.DB, error) {
	return stdlib.OpenDBFromPool(pool), nil
}

// ProvideGormPostgres provides a *gorm.DB for ORM usage
type ProvideGormPostgresParamFx struct {
	fx.In

	Cfg       config.Cfg
	Server    config.Server
	SqlDB     *sql.DB
	Lifecycle fx.Lifecycle

	DebugLog *xlog.DebugLogger `optional:"true"`
}

func ProvideGormPostgres(p ProvideGormPostgresParamFx) (*gorm.DB, error) {
	p.Server.Name = fmt.Sprintf("%s-gorm-db", p.Server.Name)

	var (
		cfg        = p.Cfg.DB["postgres"]
		gormLogger = logger.Default
	)

	// Configure GORM logger (optional: Silent for less noise)
	if p.DebugLog != nil {
		gormLogger = xlog.NewGormLogger(
			xlog.NewLogger(p.DebugLog.Logger),
			logger.Warn,
			200*time.Millisecond,
		)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: p.SqlDB}), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, err
	}

	// Get the underlying *sql.DB for pooling and lifecycle management
	sdb, err := db.DB()
	if err != nil {
		return nil, err
	}

	sdb.SetMaxOpenConns(cfg.Options.MaxOpenConnection)
	sdb.SetMaxIdleConns(cfg.Options.MaxIdleConnection)
	sdb.SetConnMaxLifetime(time.Duration(cfg.Options.MaxConnectionLifetime) * time.Second)

	// Verify connection
	if err := sdb.Ping(); err != nil {
		return nil, err
	}

	p.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return sdb.Close()
		},
	})

	return db, nil
}
