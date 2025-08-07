.PHONY: help setup verify-go sqlc-gen go-help go-tidy go-run go-build gorm-gen sqlc-gen \
        print-path deploy-help deploy-core-up deploy-core-down \
        deploy-core-down-clean deploy-core-logs disconnect-core-dependency \
        git-export-all git-export-last git-export-range git-export-clean \
        migrate-help migrate-up migrate-down migrate-status migrate-reset \
        migrate-up-to migrate-down-to \
        seed-up seed-down seed-status seed-reset \
        seed-up-to seed-down-to

# ------------------------------
# Binary Availability Checks
# ------------------------------

# Check if required binaries exist
GO_EXISTS := $(shell command -v go 2>/dev/null)
SQLC_EXISTS := $(shell command -v sqlc 2>/dev/null)

# Function to check binary availability
check-binary = $(if $(shell command -v $(1) 2>/dev/null),$(info ✅ $(1) found),$(error "$(1) not found in PATH. Please install it first"))

# Set Go environment if Go is installed
ifeq ($(GO_EXISTS),)
  GOPATH :=
  GOPATH_BIN :=
else
  GOPATH := $(shell go env GOPATH)
  GOPATH_BIN := $(GOPATH)/bin
  PATH := $(GOPATH_BIN):$(PATH)
endif

# ------------------------------
# Help Section
# ------------------------------
help:
	@echo "Available Commands:"
	@echo ""
	@echo "Setup:"
	@echo "  make setup                      Install required tools (Goose, jq)"
	@echo "  make print-path                 Show GOPATH and toolchain versions"
	@echo ""
	@echo "Go Development:"
	@echo "  make go-tidy                    Tidy and clean Go dependencies"
	@echo "  make go-run a=<app>             Run Go app from ./cmd/<app> (use c=<cfg> for config)"
	@echo "  make go-build a=<app>           Build Go app to ./bin/<app>"
	@echo "  make sqlc-gen                   Generate SQLC repositories"
	@echo "  make gorm-gen                   Generate GORM models/queries via ./cmd/gorm"
	@echo ""
	@echo "Database Migrations (Goose):"
	@echo "  make migrate-up                   Apply all migrations"
	@echo "  make migrate-down                 Roll back last migration"
	@echo "  make migrate-status               Show migration status"
	@echo "  make migrate-reset                Roll back ALL migrations"
	@echo "  make migrate-up-to v=<version>    Migrate up to a specific version"
	@echo "  make migrate-down-to v=<version>  Roll back down to a specific version"
	@echo "  make migrate-create n=<name>      Create new migration (default: .sql)"
	@echo "  make migrate-fix                  Fix migration filenames"
	@echo ""
	@echo "Database Seeders (Goose):"
	@echo "  make seed-up                      Apply all seeders"
	@echo "  make seed-down                    Roll back last seeder"
	@echo "  make seed-status                  Show seeder status"
	@echo "  make seed-reset                   Roll back ALL seeders"
	@echo "  make seed-up-to v=<version>       Apply seeders up to version"
	@echo "  make seed-down-to v=<version>     Roll back seeders down to version"
	@echo "  make seed-create n=<name>         Create new seeder (default: .sql)"
	@echo "  make seed-fix                     Fix seeder filenames"
	@echo ""
	@echo "Docker Deployment:"
	@echo "  make deploy-core-up v=<ver>     Deploy core stack (version X.Y.Z)"
	@echo "  make deploy-core-down           Stop core stack"
	@echo "  make deploy-core-down-clean     Stop & clean core stack (DANGEROUS)"
	@echo "  make deploy-core-logs [f=1]     Tail logs (f=1 to follow)"
	@echo "  make disconnect-core-dependency Disconnect dependencies (core)"
	@echo ""
	@echo "Git Export:"
	@echo "  make git-export-all [s=YYYY-MM-DD u=YYYY-MM-DD l=N]"
	@echo "  make git-export-last N=5"
	@echo "  make git-export-range s=YYYY-MM-DD u=YYYY-MM-DD"
	@echo "  make git-export-clean"
	@echo ""
	@echo "Examples:"
	@echo "  make setup"
	@echo "  make go-run a=core"
	@echo "  make gorm-gen"
	@echo "  make migrate-up"
	@echo "  make migrate-up-to v=20250728120000"
	@echo "  make seed-down-to v=20250728100000"
	@echo "  make deploy-core-up v=1.2.3"

# ------------------------------
# Setup Section
# ------------------------------

setup: verify-go
	@echo "\nInstalling development tools..."
	@echo "- Installing Goose (migration tool)"
	@go install github.com/pressly/goose/v3/cmd/goose@latest || { echo 'Error: goose installation failed.'; exit 1; }
	@echo "- Installing Sqlc (codegen tool)"
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest || { echo 'Error: sqlc installation failed.'; exit 1; }
	@echo "- Checking jq (JSON processor)"
	@if ! command -v jq >/dev/null 2>&1; then \
		echo "  jq not found, installing..."; \
		if [ "$$(uname)" = "Darwin" ]; then \
			brew install jq || { echo 'Error: Failed to install jq with brew.'; exit 1; }; \
		else \
			sudo apt-get update && sudo apt-get install -y jq || { echo 'Error: Failed to install jq with apt.'; exit 1; }; \
		fi \
	else \
		echo "  jq already installed."; \
	fi
	@echo "- Installed tools in $(GOPATH_BIN):"
	@ls $(GOPATH_BIN) || { echo 'Error: Failed to list tools in GOPATH/bin.'; exit 1; }
	@echo "\n✅ Setup completed successfully"

sqlc-gen: verify-sqlc
	@echo "Generating SQL code..."
	@sqlc generate
	@echo "✅ SQL code generation completed"

verify-sqlc:
ifndef SQLC_EXISTS
	$(error "sqlc not found. Please run 'make setup' first")
endif
	@echo "✅ sqlc is available"

verify-go:
	$(call check-binary,go)

# ------------------------------
# Go Development Section
# ------------------------------

print-path: verify-go
	@echo "Development Environment:"
	@echo "GOPATH:     $(GOPATH)"
	@echo "GOPATH_BIN: $(GOPATH_BIN)"
	@echo "PATH:       $(PATH)"
	@echo "\nTool Versions:"
	@go version || true
	@sqlc version || true

go-help: verify-go
	@echo "Go Commands:"
	@echo "  make go-tidy             - Tidy and clean dependencies"
	@echo "  make go-run a=<app>      - Run app in ./cmd/<app>"
	@echo "  make go-build a=<app>    - Build app to ./bin/<app>"
	@echo "  make sqlc-gen            - Generate SQLC repositories"
	@echo "  make gorm-gen            - Generate GORM models/queries"
	@echo ""

go-tidy: verify-go
	@go mod tidy
	@echo "✅ Go dependencies cleaned"

go-run: verify-go
	@if [ -z "$(a)" ]; then \
		echo "Error: Application name required (a=app_name)"; \
		exit 1; \
	fi
	@if [ -z "$(c)" ]; then \
		go run ./cmd/$(a); \
	else \
		go run ./cmd/$(a) -cfg=$(c); \
	fi

go-build: verify-go
	@if [ -z "$(a)" ]; then \
		echo "Error: Application name required (a=app_name)"; \
		exit 1; \
	fi
	@go build -o ./bin/$(a) ./cmd/$(a)
	@echo "✅ Built application: ./bin/$(a)"

gorm-gen: verify-go
	@echo "Generating GORM models/queries..."
	@go run ./cmd/gorm
	@echo "✅ GORM generation completed"

gorm-gen-clean: verify-go
	@echo "Deleting GORM models/queries..."
	@rm -rf ./gen/gorm
	@echo "✅ Delete GORM generation completed"

# ------------------------------
# Docker Deployment Section (unchanged)
# ------------------------------

deploy-help:
	@echo "Docker Deployment Commands"
	@echo ""
	@echo "Usage: make deploy-[core]-[up|down] v=<version>"
	@echo "  f=1 -> follow-logs (default f=0)"
	@echo ""
	@echo "Examples:"
	@echo "  make deploy-core-up v=0.0.1"
	@echo "  make deploy-core-logs f=1"

deploy-core-up:
	@[ -n "$(v)" ] || { echo "Error: v is not set."; exit 1; }
	@echo "Validating version: $(v)"
	@echo "$(v)" | grep -Eq '^v?[0-9]+\.[0-9]+\.[0-9]+$$' || { echo "Error: Version must follow 'vX.Y.Z' or 'X.Y.Z'. See https://semver.org"; exit 1; }
	@./deploy.sh stack.core.env $(v)

deploy-core-down:
	@./deploy.sh stack.core.env down

deploy-core-down-clean:
	@echo "WARNING: This will DOWN and CLEAN the core stack! This is irreversible."
	@read -p "Are you sure? (y/N): " ans; \
	if [ "$$ans" = "y" ] || [ "$$ans" = "Y" ]; then \
		./deploy.sh stack.core.env down-clean; \
	else \
		echo "Aborted."; \
	fi

deploy-core-logs:
	@f=$(f) ./deploy.sh stack.core.env logs

disconnect-core-dependency:
	@./deploy.sh stack.core.env dc-ctr

# ------------------------------
# Git Export Section (unchanged)
# ------------------------------

# Variables (can be overridden)
s ?=     # since
u ?=     # until
l ?=     # limit
n ?=     # last N commits

SCRIPT := git-export.script.sh
TIMESTAMP := $(shell date +%s)
OUTPUT := git-export-commits-$(TIMESTAMP).json

git-export-all:
	@echo "Exporting commits to $(OUTPUT)..."
	@bash $(SCRIPT) $(if $(s),--since=$(s)) $(if $(u),--until=$(u)) $(if $(l),--limit=$(l)) > $(OUTPUT)
	@echo "Done. Output saved to $(OUTPUT)."

git-export-last:
	@if [ -z "$(n)" ]; then echo "Usage: make git-export-last n=5"; exit 1; fi
	@echo "Exporting last $(n) commits to $(OUTPUT)..."
	@bash $(SCRIPT) --limit=$(n) > $(OUTPUT)
	@echo "Done. Output saved to $(OUTPUT)."

git-export-range:
	@if [ -z "$(s)" ] || [ -z "$(u)" ]; then echo "Usage: make git-export-range s=YYYY-MM-DD u=YYYY-MM-DD"; exit 1; fi
	@echo "Exporting commits from $(s) to $(u) to $(OUTPUT)..."
	@bash $(SCRIPT) --since=$(s) --until=$(u) > $(OUTPUT)
	@echo "Done. Output saved to $(OUTPUT)."

git-export-clean:
	@rm -f git-export-commits-*.json
	@echo "Removed all exported JSON files."

# ------------------------------
# Database Migration (Goose CLI)
# ------------------------------

define MIGRATE_CMD
	@if [ -z "$(c)" ]; then \
		go run ./cmd/migrate migration $(1); \
	else \
		go run ./cmd/migrate migration -cfg=$(c) $(1); \
	fi
endef

define SEED_CMD
	@if [ -z "$(c)" ]; then \
		go run ./cmd/migrate seeder $(1); \
	else \
		go run ./cmd/migrate seeder -cfg=$(c) $(1); \
	fi
endef

migrate-help:
	@echo "Database Migration Commands (via Fx + Goose)"
	@echo ""
	@echo "Usage:"
	@echo "  make migrate-up                   Apply all migrations"
	@echo "  make migrate-down                 Roll back last migration"
	@echo "  make migrate-status               Show migration status"
	@echo "  make migrate-reset                Roll back ALL migrations"
	@echo "  make migrate-up-to v=<version>    Migrate up to a specific version"
	@echo "  make migrate-down-to v=<version>  Roll back down to a specific version"
	@echo "  make migrate-create n=<name>      Create new migration (default: .sql)"
	@echo "  make migrate-fix                  Fix migration filenames"
	@echo ""
	@echo "Seeder Commands:"
	@echo "  make seed-up                      Apply all seeders"
	@echo "  make seed-down                    Roll back last seeder"
	@echo "  make seed-status                  Show seeder status"
	@echo "  make seed-reset                   Roll back ALL seeders"
	@echo "  make seed-up-to v=<version>       Apply seeders up to version"
	@echo "  make seed-down-to v=<version>     Roll back seeders down to version"
	@echo "  make seed-create n=<name>         Create new seeder (default: .sql)"
	@echo "  make seed-fix                     Fix seeder filenames"
	@echo ""
	@echo "Examples:"
	@echo "  make migrate-up"
	@echo "  make migrate-up-to v=20250728120000"
	@echo "  make seed-down-to v=20250728110000"
	@echo "  make migrate-create n=create_users_table"
	@echo "  make seed-create n=init_data"

# Migration commands
migrate-up:
	$(call MIGRATE_CMD,up)

migrate-down:
	$(call MIGRATE_CMD,down)

migrate-status:
	$(call MIGRATE_CMD,status)

migrate-reset:
	$(call MIGRATE_CMD,reset)

migrate-up-to:
	@[ -n "$(v)" ] || { echo "Error: v is required (e.g., make migrate-up-to v=20250728120000)"; exit 1; }
	$(call MIGRATE_CMD,up-to $(v))

migrate-down-to:
	@[ -n "$(v)" ] || { echo "Error: v is required (e.g., make migrate-down-to v=20250728120000)"; exit 1; }
	$(call MIGRATE_CMD,down-to $(v))

migrate-create:
	@[ -n "$(n)" ] || { echo "Error: name is required (e.g., make migrate-create n=create_users_table)"; exit 1; }
	$(call MIGRATE_CMD,create $(n))

migrate-fix:
	$(call MIGRATE_CMD,fix)

# Seeder commands
seed-up:
	$(call SEED_CMD,up)

seed-down:
	$(call SEED_CMD,down)

seed-status:
	$(call SEED_CMD,status)

seed-reset:
	$(call SEED_CMD,reset)

seed-up-to:
	@[ -n "$(v)" ] || { echo "Error: v is required (e.g., make seed-up-to v=20250728120000)"; exit 1; }
	$(call SEED_CMD,up-to $(v))

seed-down-to:
	@[ -n "$(v)" ] || { echo "Error: v is required (e.g., make seed-down-to v=20250728120000)"; exit 1; }
	$(call SEED_CMD,down-to $(v))

seed-create:
	@[ -n "$(n)" ] || { echo "Error: name is required (e.g., make seed-create n=init_data)"; exit 1; }
	$(call SEED_CMD,create $(n))

seed-fix:
	$(call SEED_CMD,fix)