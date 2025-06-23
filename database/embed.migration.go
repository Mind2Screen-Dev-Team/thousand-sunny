package database

import "embed"

//go:embed migrations/*.sql
var EmbedMigrations embed.FS

//go:embed seeders/*.sql
var EmbedSeeders embed.FS
