package main

import (
	"context"
	"fmt"
	"slices"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/dependency"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/app/injector"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xutil"

	"go.uber.org/fx"

	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	app := fx.New(
		fx.NopLogger,
		injector.GlobalConfig,
		injector.Database,
		injector.GormDatabase,
		fx.Provide(dependency.ProvideOtelConfig),
		fx.Provide(dependency.ProvideOtelGrpcClient),
		fx.Provide(dependency.ProvideOtelResource),
		fx.Provide(dependency.ProvideOtelLog),
		fx.Provide(dependency.ProvideDebugLogger),
		fx.Provide(func() config.Server { return config.Server{Name: "generation-command"} }),
		fx.Invoke(RunGormGen),
	)
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	if err := app.Stop(context.Background()); err != nil {
		panic(err)
	}
}

func RunGormGen(db *gorm.DB) {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./gen/gorm/query",
		ModelPkgPath:      "./gen/gorm/model",
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})
	g.UseDB(db)

	g.Config.WithJSONTagNameStrategy(xutil.SnakeToCamel)

	// # Example For Mapping Data Type
	// g.Config.WithImportPkgPath("github.com/google/uuid")
	// g.WithDataTypeMap(map[string]func(gorm.ColumnType) string{
	// 	"UNIQUEIDENTIFIER": func(ct gorm.ColumnType) string { return "uuid.UUID" },
	// 	"BIT":              func(ct gorm.ColumnType) string { return "bool" },
	// 	"DECIMAL":          func(ct gorm.ColumnType) string { return "float64" },
	// 	"DATETIME":         func(ct gorm.ColumnType) string { return "time.Time" },
	// 	"NVARCHAR":         func(ct gorm.ColumnType) string { return "string" },
	// 	"VARCHAR":          func(ct gorm.ColumnType) string { return "string" },
	// })

	tables, err := db.Migrator().GetTables()
	if err != nil {
		panic(fmt.Errorf("get tables failed: %w", err))
	}
	tables = filterTables(tables, []string{"migration_db_version", "seeder_db_version"})

	fmt.Printf("Generating %d models: %v\n", len(tables), tables)

	var models []any
	for _, table := range tables {
		switch table {
		// # Example Preload Relationship
		// case "user_relations":
		// 	models = append(models, g.GenerateModel("user_relations",
		// 		gen.FieldRelate(field.BelongsTo, "RelatedUser", g.GenerateModel("users"),
		// 			&field.RelateConfig{RelatePointer: true, GORMTag: field.GormTag{"foreignKey": []string{"user_id"}, "references": []string{"id"}}, JSONTag: "relatedUser"}),
		// 		gen.FieldRelate(field.BelongsTo, "RelatedOrganization", g.GenerateModel("organizations"),
		// 			&field.RelateConfig{RelatePointer: true, GORMTag: field.GormTag{"foreignKey": []string{"organization_id"}, "references": []string{"id"}}, JSONTag: "relatedOrganization"}),
		// 		gen.FieldRelate(field.BelongsTo, "RelatedRole", g.GenerateModel("roles"),
		// 			&field.RelateConfig{RelatePointer: true, GORMTag: field.GormTag{"foreignKey": []string{"role_id"}, "references": []string{"id"}}, JSONTag: "relatedRole"}),
		// 	))
		default:
			models = append(models, g.GenerateModel(table))
		}
	}

	g.ApplyBasic(models...)
	g.Execute()
}

func filterTables(tables, exclude []string) []string {
	return slices.DeleteFunc(tables, func(t string) bool {
		return slices.Contains(exclude, t)
	})
}
