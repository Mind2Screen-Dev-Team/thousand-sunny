.PHONY: help setup sqlc-gen go-help go-tidy go-run go-build print-path \
        deploy-help deploy-core-up deploy-core-down \
        deploy-async-up deploy-async-down \
        deploy-core-down-clean deploy-async-down-clean \
        deploy-async-logs deploy-core-logs \
        disconnect-core-dependency disconnect-async-dependency

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
	@echo "Application Available Commands"
	@echo ""
	@echo "Usage: make [command]"
	@echo ""
	@echo "Commands:"
	@echo "  setup                      Install required tools"
	@echo "  sqlc-gen                   Run sqlc code generation"
	@echo "  go-help                    Show Go-related commands"
	@echo "  deploy-help                Show Docker deployment commands"
	@echo "  print-path                 Print current GOPATH and PATH"
	@echo ""
	@echo "Examples:"
	@echo "  make setup"
	@echo "  make sqlc-gen"
	@echo "  make go-run a=core"
	@echo "  make deploy-core-up v=0.0.1"

# ------------------------------
# Project Setup with Enhanced Checks
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

# Verify Go is installed
verify-go:
	$(call check-binary,go)

# ------------------------------
# Development Tools Section
# ------------------------------

print-path: verify-go
	@echo "Development Environment:"
	@echo "GOPATH:    $(GOPATH)"
	@echo "GOPATH_BIN: $(GOPATH_BIN)"
	@echo "PATH:      $(PATH)"
	@echo "\nTool Versions:"
	@go version || true
	@sqlc version || true

# SQL Code Generation with verification
sqlc-gen: verify-sqlc
	@echo "Generating SQL code..."
	@sqlc generate
	@echo "✅ SQL code generation completed"

# Verify sqlc is installed
verify-sqlc:
ifndef SQLC_EXISTS
	$(error "sqlc not found. Please run 'make setup' first")
endif
	@echo "✅ sqlc is available"

# ------------------------------
# Go Commands Section
# ------------------------------

go-help: verify-go
	@echo "Go Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make go-tidy"
	@echo "  make go-run a=<app> [c=<config>]"
	@echo "  make go-build a=<app>"
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

# ------------------------------
# Docker Deployment Section 
# (Remains unchanged from original)
# ------------------------------

deploy-help:
	@echo "Docker Deployment Commands"
	@echo ""
	@echo "Usage: make deploy-[async|core]-[up|down] v=<version>"
	@echo ""
	@echo "Paramater:"
	@echo "       f=1 -> follow-logs"
	@echo "       f=0 -> not-follow-logs (default)"
	@echo ""
	@echo "Examples:"
	@echo "  make deploy-core-up v=0.0.1"
	@echo "  make deploy-async-down"
	@echo "  make deploy-core-logs"
	@echo "  make deploy-async-logs f=1"

deploy-async-up:
	@[ -n "$(v)" ] || { echo "Error: v is not set."; exit 1; }
	@echo "Validating version: $(v)"
	@echo "$(v)" | grep -Eq '^v?[0-9]+\.[0-9]+\.[0-9]+$$' || { echo "Error: Version must follow 'vX.Y.Z' or 'X.Y.Z'. See https://semver.org"; exit 1; }
	@./deploy.sh stack.async.env $(v)

deploy-core-up:
	@[ -n "$(v)" ] || { echo "Error: v is not set."; exit 1; }
	@echo "Validating version: $(v)"
	@echo "$(v)" | grep -Eq '^v?[0-9]+\.[0-9]+\.[0-9]+$$' || { echo "Error: Version must follow 'vX.Y.Z' or 'X.Y.Z'. See https://semver.org"; exit 1; }
	@./deploy.sh stack.core.env $(v)

deploy-async-down:
	@./deploy.sh stack.async.env down

deploy-core-down:
	@./deploy.sh stack.core.env down

deploy-async-down-clean:
	@echo "WARNING: This action will DOWN and CLEAN the async stack! This is dangerous and irreversible."
	@read -p "Are you sure you want to down-clean async stack? (y/N): " ans; \
	if [ "$$ans" = "y" ] || [ "$$ans" = "Y" ]; then \
		./deploy.sh stack.async.env down-clean; \
	else \
		echo "Aborted."; \
	fi

deploy-core-down-clean:
	@echo "WARNING: This action will DOWN and CLEAN the core stack! This is dangerous and irreversible."
	@read -p "Are you sure you want to down-clean Core stack? (y/N): " ans; \
	if [ "$$ans" = "y" ] || [ "$$ans" = "Y" ]; then \
		./deploy.sh stack.core.env down-clean; \
	else \
		echo "Aborted."; \
	fi

deploy-async-logs:
	@f=$(f) ./deploy.sh stack.async.env logs

deploy-core-logs:
	@f=$(f) ./deploy.sh stack.core.env logs

disconnect-async-dependency:
	@./deploy.sh stack.async.env dc-ctr

disconnect-core-dependency:
	@./deploy.sh stack.core.env dc-ctr

# Variables (can be overridden)
SINCE ?=
UNTIL ?=
LIMIT ?=
N ?=

SCRIPT := git-export.script.sh
TIMESTAMP := $(shell date +%s)
OUTPUT := git-export-commits-$(TIMESTAMP).json

# Export all commits (with optional filters)
git-export-all:
	@echo "Exporting commits to $(OUTPUT)..."
	@bash $(SCRIPT) $(if $(SINCE),--since=$(SINCE)) $(if $(UNTIL),--until=$(UNTIL)) $(if $(LIMIT),--limit=$(LIMIT)) > $(OUTPUT)
	@echo "Done. Output saved to $(OUTPUT)."

# Export last N commits (e.g., make git-export-last N=5)
git-export-last:
	@if [ -z "$(N)" ]; then echo "Usage: make git-export-last N=5"; exit 1; fi
	@echo "Exporting last $(N) commits to $(OUTPUT)..."
	@bash $(SCRIPT) --limit=$(N) > $(OUTPUT)
	@echo "Done. Output saved to $(OUTPUT)."

# Export commits within a date range
# Usage: make git-export-range SINCE=YYYY-MM-DD UNTIL=YYYY-MM-DD
git-export-range:
	@if [ -z "$(SINCE)" ] || [ -z "$(UNTIL)" ]; then echo "Usage: make git-export-range SINCE=YYYY-MM-DD UNTIL=YYYY-MM-DD"; exit 1; fi
	@echo "Exporting commits from $(SINCE) to $(UNTIL) to $(OUTPUT)..."
	@bash $(SCRIPT) --since=$(SINCE) --until=$(UNTIL) > $(OUTPUT)
	@echo "Done. Output saved to $(OUTPUT)."

# Clean all exported JSON files
git-export-clean:
	@rm -f git-export-commits-*.json
	@echo "Removed all exported JSON files."