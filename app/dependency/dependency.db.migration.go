package dependency

import (
	"context"
	"database/sql"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/database"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/pressly/goose/v3"
	"github.com/rs/xid"
	"go.uber.org/fx"
)

func InvokeMigrations(c config.Cfg, s config.Server, db *sql.DB, l *xlog.DebugLogger, lc fx.Lifecycle) error {
	var (
		dialect = "postgres"
		dir     = "migrations"
		ctx     = context.WithValue(context.Background(), xlog.XLOG_REQ_TRACE_ID_CTX_KEY, xid.New().String())
	)

	if err := db.Ping(); err != nil {
		return err
	}

	if c.Log.ConsoleFormat == "json" {
		goose.SetLogger(xlog.NewGooseLogger(xlog.NewLogger(l.Logger), context.Background()))
	}
	goose.SetBaseFS(database.EmbedMigrations)
	goose.SetTableName("migration_db_version")

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.UpContext(ctx, db, dir); err != nil {
		return err
	}

	if err := goose.StatusContext(ctx, db, dir); err != nil {
		return err
	}

	return nil
}

func InvokeSeeders(c config.Cfg, s config.Server, db *sql.DB, l *xlog.DebugLogger, lc fx.Lifecycle) error {
	var (
		dialect = "postgres"
		dir     = "seeders"
		ctx     = context.WithValue(context.Background(), xlog.XLOG_REQ_TRACE_ID_CTX_KEY, xid.New().String())
	)

	if err := db.Ping(); err != nil {
		return err
	}

	if c.Log.ConsoleFormat == "json" {
		goose.SetLogger(xlog.NewGooseLogger(xlog.NewLogger(l.Logger), context.Background()))
	}
	goose.SetBaseFS(database.EmbedSeeders)
	goose.SetTableName("seeder_db_version")

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.UpContext(ctx, db, dir); err != nil {
		return err
	}

	if err := goose.StatusContext(ctx, db, dir); err != nil {
		return err
	}

	return nil
}
