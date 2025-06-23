package dependency

import (
	"database/sql"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/database"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

func InvokeMigrations(c config.Cfg, s config.Server, db *sql.DB, lc fx.Lifecycle) error {
	var (
		dialect = "postgres"
		dir     = "migrations"
	)

	if err := db.Ping(); err != nil {
		return err
	}

	goose.SetBaseFS(database.EmbedMigrations)
	goose.SetTableName("migration_db_version")

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.Up(db, dir); err != nil {
		return err
	}

	if err := goose.Status(db, dir); err != nil {
		return err
	}

	return nil
}

func InvokeSeeders(c config.Cfg, s config.Server, db *sql.DB, lc fx.Lifecycle) error {
	var (
		dialect = "postgres"
		dir     = "seeders"
	)

	if err := db.Ping(); err != nil {
		return err
	}

	goose.SetBaseFS(database.EmbedSeeders)
	goose.SetTableName("seeder_db_version")

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.Up(db, dir); err != nil {
		return err
	}

	if err := goose.Status(db, dir); err != nil {
		return err
	}

	return nil
}
