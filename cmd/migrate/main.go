package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/injector"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/database"

	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Open-Detail: make migrate-help")
		os.Exit(1)
	}

	var (
		dirType = "migration"
		cmdArgs = args
	)

	if args[0] == "migration" || args[0] == "seeder" {
		dirType = args[0]
		cmdArgs = args[1:]
	}

	Run(
		fx.NopLogger,
		injector.GlobalConfig,
		injector.Database,
		fx.Provide(dependency.ProvideOtelConfig),
		fx.Provide(dependency.ProvideOtelGrpcClient),
		fx.Provide(dependency.ProvideOtelResource),
		fx.Provide(dependency.ProvideOtelLog),
		fx.Provide(dependency.ProvideDebugLogger),
		fx.Provide(func() config.Server { return config.Server{Name: "migration-command"} }),
		fx.Invoke(func(c config.Cfg, s config.Server, db *sql.DB) {
			if err := RunGooseCommand(c, s, db, dirType, cmdArgs); err != nil {
				fmt.Println("[Goose] ERROR:", err)
			}
		}),
	)
}

func Run(opts ...fx.Option) {
	app := fx.New(opts...)

	if err := app.Start(context.Background()); err != nil {
		fmt.Println("[Fx] START FAILED\t" + err.Error())
		return
	}
	if err := app.Stop(context.Background()); err != nil {
		fmt.Println("[Fx] STOP FAILED\t" + err.Error())
		return
	}
}

// RunGooseCommand runs a Goose command (up, down, status, etc.) for either migrations or seeders.
func RunGooseCommand(c config.Cfg, s config.Server, db *sql.DB, dirType string, args []string) error {
	var (
		dialect = "postgres"
	)

	var (
		dir   string
		table string
	)

	// Choose directory and table based on type
	switch dirType {
	case "migration":
		dir = "migrations"
		table = "migration_db_version"
		goose.SetBaseFS(database.EmbedMigrations)
	case "seeder":
		dir = "seeders"
		table = "seeder_db_version"
		goose.SetBaseFS(database.EmbedSeeders)
	default:
		return fmt.Errorf("unknown type: %s (must be migrations or seeders)", dirType)
	}

	if err := db.Ping(); err != nil {
		return err
	}

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}
	goose.SetTableName(table)

	if len(args) == 0 {
		args = []string{"status"} // default command
	}

	cmd := args[0]
	cmdArgs := args[1:]

	// Dispatch commands
	switch cmd {
	case "up":
		return goose.Up(db, dir)
	case "up-by-one":
		return goose.UpByOne(db, dir)
	case "up-to":
		if len(cmdArgs) < 1 {
			return fmt.Errorf("up-to requires VERSION")
		}
		return goose.UpTo(db, dir, mustParseVersion(cmdArgs[0]))
	case "down":
		return goose.Down(db, dir)
	case "down-to":
		if len(cmdArgs) < 1 {
			return fmt.Errorf("down-to requires VERSION")
		}
		return goose.DownTo(db, dir, mustParseVersion(cmdArgs[0]))
	case "redo":
		return goose.Redo(db, dir)
	case "reset":
		return goose.Reset(db, dir)
	case "status":
		return goose.Status(db, dir)
	case "version":
		v, err := goose.GetDBVersion(db)
		if err != nil {
			return err
		}
		fmt.Printf("Current DB version: %d\n", v)
		return nil

	case "create":
		if len(cmdArgs) < 1 {
			return fmt.Errorf("create requires at least NAME")
		}
		name := cmdArgs[0]
		if err := goose.Create(nil, dir, name, "sql"); err != nil {
			return err
		}
		fmt.Printf("Created new %s: %s\n", "sql", name)
		return nil

	case "fix":
		return goose.Fix(dir)

	default:
		return fmt.Errorf("unsupported goose command: %s", cmd)
	}
}

// mustParseVersion converts a string to int64 or exits on error.
func mustParseVersion(s string) int64 {
	var v int64
	_, err := fmt.Sscan(s, &v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid version: %s\n", s)
		os.Exit(1)
	}
	return v
}
