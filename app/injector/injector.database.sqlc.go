package injector

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/gen/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var (
	RepoGenerationSqlc = fx.Options(
		fx.Module("dependency:database:sqlc:repo",
			fx.Provide(func(db *pgxpool.Pool) sqlc.DBTX {
				return db
			}),
			fx.Provide(sqlc.New),
		),
	)
)
