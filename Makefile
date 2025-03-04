# Check and create .env file if not present using a shell command
.env:
	@echo ""
	@echo "Read .env file not found. Copying from .example.env..."
	@echo ""
	@cp .example.env .env

include .env

# Set GOPATH to your desired Go path, or it defaults to the system GOPATH if not set
GOPATH ?= $(shell go env GOPATH)
GOPATH_BIN := $(GOPATH)/bin

# Prepend GOPATH/bin to PATH if not already present
ifneq ($(filter $(GOPATH_BIN), $(PATH)),)
    export PATH := $(PATH)
else
    export PATH := $(GOPATH_BIN):$(PATH)
endif

help:
	@echo "Application Available Commands"
	@echo ""
	@echo "Usage: make [commands]"
	@echo ""
	@echo "Commands:"
	@echo "  setup                      Make setup your project workspace"
	@echo "  migrate-help               Make migarte help command"
	@echo "  go-help                    Make go help command"
	@echo "  print-path                 Make print current path variable"
	@echo ""
	@echo "Examples:"
	@echo "  make setup"
	@echo "  make sqlc-gen"
	@echo "  make go-help"
	@echo "  make migrate-help"
	@echo "  make deploy-help"

deploy-help:
	@echo "Application Deployment Available Commands"
	@echo ""
	@echo "Usage: make deploy-[commands]"
	@echo ""
	@echo "Commands:"
	@echo "  deploy-[asynq|core] v=[major|minor|patch]  Make deploy app with automatic increment sematic version"
	@echo "  deploy-[asynq|core]-rebuild v=0.0.1        Make deploy app with force rebuild docker image with specific tag version"
	@echo "  deploy-[asynq|core]-down                   Make deployed app shutdown from docker"
	@echo ""
	@echo "Examples:"
	@echo "  make deploy-[asynq|core] v=major"
	@echo "  make deploy-[asynq|core] v=minor"
	@echo "  make deploy-[asynq|core] v=patch"
	@echo "  make deploy-[asynq|core]-rebuild v=0.0.1"
	@echo "  make deploy-[asynq|core]-rebuild v=0.1.1"
	@echo "  make deploy-[asynq|core]-down"

# Command Deployment
deploy-asynq:
	@[ -n "$(v)" ] || { echo "Error: v is not set."; exit 1; }
	@echo "Validating version: $(v)"
	@echo "$(v)" | grep -Eq '^(major|minor|patch|latest)$$' || { echo "Error: v must be one of: major, minor, patch, latest."; exit 1; }
	@./deploy.sh asynq $(v)

deploy-asynq-rebuild:
	@[ -n "$(v)" ] || { echo "Error: v is not set."; exit 1; }
	@echo "Validating version: $(v)"
	@echo "$(v)" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$$' || { echo "Error: Version must follow the semantic versioning format 'vX.Y.Z' (e.g., v1.0.0). Please refer to: https://semver.org/"; exit 1; }
	@./deploy.sh asynq v$(v) rebuild

deploy-asynq-down:
	@docker compose -p asynq-$(APP_ENV)-app-stack down

deploy-core:
	@[ -n "$(v)" ] || { echo "Error: v is not set."; exit 1; }
	@echo "Validating version: $(v)"
	@echo "$(v)" | grep -Eq '^(major|minor|patch|latest)$$' || { echo "Error: v must be one of: major, minor, patch, latest."; exit 1; }
	@./deploy.sh core $(v)

deploy-core-rebuild:
	@[ -n "$(v)" ] || { echo "Error: v is not set."; exit 1; }
	@echo "Validating version: $(v)"
	@echo "$(v)" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$$' || { echo "Error: Version must follow the semantic versioning format 'vX.Y.Z' (e.g., v1.0.0). Please refer to: https://semver.org/"; exit 1; }
	@./deploy.sh core v$(v) rebuild

deploy-core-down:
	@docker compose -p core-$(APP_ENV)-app-stack down

# Command Setup
setup:
	@echo "- Install Goose For Migration & Seeder Tool"
	@go install github.com/pressly/goose/v3/cmd/goose@latest || { echo 'Error: goose installation failed.'; exit 1; }
	@echo "- Install Sqlc For Code-Gen Repository Pattern Tool"
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest || { echo 'Error: sqlc installation failed.'; exit 1; }
	@echo "- Checking installed tools in GOPATH/bin..."
	@echo ""
	@ls $(shell go env GOPATH)/bin || { echo 'Error: check installation tools in GOPATH/bin.'; exit 1; }
	@echo ""
	@echo "Setup Installation Tools Is Complete!"

sqlc-gen:
	@sqlc generate

# Command to run golang commands

go-help:
	@echo "Application Golang Commands"
	@echo ""
	@echo "Usage: make [commands] [OPTIONS]"
	@echo ""
	@echo "Commands:"
	@echo "  go-tidy                    Go module tidy"
	@echo "  go-run APPLICATION         Go run application"
	@echo "  go-build APPLICATION       Go build application"
	@echo ""
	@echo "Examples:"
	@echo "  make go-tidy"
	@echo "  make go-run a=core"
	@echo "  make go-run a=core c=<your_config_yaml_file>"
	@echo "  make go-build a=core"

go-tidy:
	@go mod tidy

go-run:
	@if [ -z "$(c)" ]; then \
		go run ./cmd/$(a); \
	else \
		go run ./cmd/$(a) -cfg=$(c); \
	fi

go-build:
	@go build -o ./bin/$(a) ./cmd/$(a)

# Help messages
seeder-help:
	@make migrate-help

migrate-help:
	@echo "Migrations and Seeders Docs"
	@echo ""
	@echo "Usage: make [migrate | seeder]-[command] [OPTIONS]"
	@echo ""
	@echo "Commands:"
	@echo "  [ migrate | seeder ]-up                    Migrate the DB to the most recent version available"
	@echo "  [ migrate | seeder ]-up-by-one             Migrate the DB up by 1"
	@echo "  [ migrate | seeder ]-up-to VERSION         Migrate the DB to a specific VERSION"
	@echo "  [ migrate | seeder ]-down                  Roll back the version by 1"
	@echo "  [ migrate | seeder ]-down-to VERSION       Roll back to a specific VERSION"
	@echo "  [ migrate | seeder ]-create NAME			Creates new migration file with the current timestamp"
	@echo "  [ migrate | seeder ]-redo                  Re-run the latest migration"
	@echo "  [ migrate | seeder ]-reset                 Roll back all migrations"
	@echo "  [ migrate | seeder ]-status                Dump the migration status for the current DB"
	@echo "  [ migrate | seeder ]-version               Print the current version of the database"
	@echo "  [ migrate | seeder ]-fix                   Apply sequential ordering to migrations"
	@echo "  [ migrate | seeder ]-validate              Check migration files without running them"
	@echo ""
	@echo "Options by env file:"
	@echo "  DB_GOOSE_DRIVER         Database driver (postgres, mysql, sqlite3, mssql, redshift, tidb, clickhouse, vertica, ydb, turso)"
	@echo "  DB_GOOSE_DBSTRING       Connection string for the database"
	@echo "  DB_GOOSE_MIGRATION_DIR  Directory for migration files (default: current directory)"
	@echo ""
	@echo "Examples:"
	@echo "  make [ migrate | seeder ]-up" 
	@echo "  make [ migrate | seeder ]-up-by-one"
	@echo "  make [ migrate | seeder ]-up-to v=20240922160357"
	@echo "  make [ migrate | seeder ]-down"
	@echo "  make [ migrate | seeder ]-down-to v=20240922160357"
	@echo "  make [ migrate | seeder ]-status"
	@echo "  make [ migrate | seeder ]-version"
	@echo "  make [ migrate | seeder ]-create n=<migration_name>"
	@echo "  make [ migrate | seeder ]-validate"

# Command to run goose with the specified options
GOOSE_MIGRATE_CMD = GOOSE_DRIVER=$(DB_GOOSE_DRIVER) GOOSE_DBSTRING="$(DB_GOOSE_DBSTRING)" GOOSE_MIGRATION_DIR=$(DB_GOOSE_MIGRATION_DIR) goose -table $(DB_GOOSE_MIGRATION_TABLE)

# Migration commands
migrate-up:
	@$(GOOSE_MIGRATE_CMD) up

migrate-up-by-one:
	@$(GOOSE_MIGRATE_CMD) up-by-one

migrate-up-to:
	@$(GOOSE_MIGRATE_CMD) up-to $(v)

migrate-down:
	@$(GOOSE_MIGRATE_CMD) down

migrate-down-to:
	@$(GOOSE_MIGRATE_CMD) down-to $(v)

migrate-create:
	@$(GOOSE_MIGRATE_CMD) create $(n) sql

migrate-redo:
	@$(GOOSE_MIGRATE_CMD) redo

migrate-reset:
	@$(GOOSE_MIGRATE_CMD) reset

migrate-status:
	@$(GOOSE_MIGRATE_CMD) status

migrate-version:
	@$(GOOSE_MIGRATE_CMD) version

migrate-fix:
	@$(GOOSE_MIGRATE_CMD) fix

migrate-validate:
	@$(GOOSE_MIGRATE_CMD) validate

# Command to run goose with the specified options
GOOSE_SEEDER_CMD = GOOSE_DRIVER=$(DB_GOOSE_DRIVER) GOOSE_DBSTRING="$(DB_GOOSE_DBSTRING)" GOOSE_MIGRATION_DIR=$(DB_GOOSE_MIGRATION_SEEDER_DIR) goose -table $(DB_GOOSE_MIGRATION_SEEDER_TABLE)

# Seeders commands
seeder-up:
	@$(GOOSE_SEEDER_CMD) up

seeder-up-by-one:
	@$(GOOSE_SEEDER_CMD) up-by-one

seeder-up-to:
	@$(GOOSE_SEEDER_CMD) up-to $(v)

seeder-down:
	@$(GOOSE_SEEDER_CMD) down

seeder-down-to:
	@$(GOOSE_SEEDER_CMD) down-to $(v)

seeder-create:
	@$(GOOSE_SEEDER_CMD) create $(n) sql

seeder-redo:
	@$(GOOSE_SEEDER_CMD) redo

seeder-reset:
	@$(GOOSE_SEEDER_CMD) reset

seeder-status:
	@$(GOOSE_SEEDER_CMD) status

seeder-version:
	@$(GOOSE_SEEDER_CMD) version

seeder-fix:
	@$(GOOSE_SEEDER_CMD) fix

seeder-validate:
	@$(GOOSE_SEEDER_CMD) validate