# 🏴‍☠️ *Thousand Sunny* – Project Skeleton

Set sail on your next adventure with **Thousand Sunny**!  
Inspired by the legendary ship from *One Piece*, this Go project skeleton is built for **speed, flexibility, and scalability**. Whether you're building small tools or large-scale applications, Thousand Sunny provides the foundation to power your journey.

![Thousand Sunny](./storage/assets/thousand-sunny.png "Thousand Sunny")

---

## 🗂 Project Structure

```bash
├── app
│   ├── dependency  # Shared external dependencies and libraries.
│   └── injector    # Handles initialization and registration of dependencies.
├── bin             # Compiled binaries or executables.
├── cmd
│   └── core        # Main entry point for the 'core' application.
├── config          # Configuration files for application settings.
├── constant        # Centralized constants to avoid hardcoding.
├── database
│   ├── migrations  # Schema migrations (version-controlled).
│   ├── queries     # SQLC query generators for custom DB operations.
│   └── seeders     # Seed data for initial or demo setups.
├── docs            # Auto-generated OpenAPI specs (openapi.yaml/json) and static docs.
├── gen
│   └── repo        # Auto-generated repository code for database access.
├── internal        # Internal packages (application-specific).
│   ├── <domain>            # Domain modules.
│   │   ├── <sub-domain>    # Sub-domains within a domain.
│   │   │   ├── <domain>.<sub-domain>.<task-name>.handler.go  # Endpoint handlers.
│   │   │   ├── <domain>.<sub-domain>.repo.go                 # Data access layer.
│   │   │   ├── <domain>.<sub-domain>.service.go              # Business logic.
│   │   │   └── <domain>.<sub-domain>.fx.module.go            # Uber Fx modules.
│   │   ├── <domain>.<task-name>.handler.go                   # Endpoint handlers.
│   │   ├── <domain>.repo.go                                  # Data access layer.
│   │   ├── <domain>.service.go                               # Business logic.
│   │   └── <domain>.fx.modules.go                            # Uber Fx modules.
│   └── fx.modules  # Global Uber Fx module definitions.
├── pkg             # Reusable libraries and utility packages.
│   ├── xfiber       # Fiber server helpers and middleware.
│   ├── xfilter      # Data filtering helpers.
│   ├── xhuma        # Extensions for Huma (OpenAPI framework integration).
│   ├── xlog         # Logging utilities.
│   ├── xmail        # Email helpers.
│   ├── xpanic       # Panic recovery utilities.
│   ├── xresp        # Standardized HTTP response utilities.
│   ├── xsecurity    # Encryption/decryption utilities.
│   ├── xtracer      # OpenTelemetry tracing helpers.
│   ├── xutil        # Generic helper functions.
│   └── xvalidate    # Validation helpers (with error mapping).
└── storage
    ├── assets      # Static assets (images, documents, etc.).
    ├── backup      # Backup data.
    ├── cron        # Cron job configurations.
    ├── template    # Templates (emails, configs, etc.).
    └── logs
        └── <server.name> # Log folders (based on `config.yaml` server name).
            ├── debug     # Debug-level logs.
            ├── io        # Input/output logs.
            └── trx       # Transaction/audit logs.
````

---

## 📋 Features

* 🗃️ **Base Project Structure** – Modular Go skeleton.
* 🌐 **Huma Framework** – Auto-generates **OpenAPI** specifications.
* 📜 **Live Docs at `/docs`** – Interactive Swagger-style UI.
* 📂 **Exportable OpenAPI Specs** – Downloadable from:
  * `http://<host>:<port>/openapi.yaml`
  * `http://<host>:<port>/openapi.json`
* 🔧 **Uber Fx & Config Integration** – Simplified dependency and config management.
* 📦 **SQLC Repositories** – Auto-generated DB repositories.
* 📜 **DTO Validation** – Structured payload validation.
* 🐳 **Docker-Ready** – Containerized for dev/prod.
* 🌐 **OpenTelemetry** – Observability with traces, metrics, and logs.

---

## 🚀 Getting Started

Clone the repository and set up your environment:

```bash
# Clone the repository
git clone git@github.com:Mind2Screen-Dev-Team/thousand-sunny.git

cd thousand-sunny

# Install dependencies and prepare tools
make setup

# Copy and configure environment variables
cp stack.example.env stack.core.env

# Copy application configuration
cp config.example.yaml config.yaml
```

### Running Locally

```bash
# Start the application locally
make go-run a=core
```

### Running with Docker

```bash
# Copy configuration
cp config.example.yaml config.yaml

# Make deploy scripts executable
chmod +x ./deploy.*.sh

# Deploy (version must follow semantic versioning: X.Y.Z)
make deploy-core-up v=<version>

# Stop services
make deploy-core-down
```

---

Then Access It:
* **With Interactive Documentation:**
  - `http://<host>:<port>/docs`

* **OpenAPI Specification (YAML):**
  - `http://<host>:<port>/openapi.yaml`

* **OpenAPI Specification (JSON):**
  - `http://<host>:<port>/openapi.json`

The `/docs` UI is powered by Huma and reads from the same OpenAPI schema.

---

## ⚙️ Makefile Commands

Some useful commands:

```bash
# Install tools (sqlc, goose, etc.)
make setup

# Generate SQLC repositories
make sqlc-gen

# Run a Go service
make go-run a=core

# Run with specific config
make go-run a=core c=config.yaml

# Build the application
make go-build a=core

# Clean up go.mod
make go-tidy

# Print PATH variable
make print-path
```

## ⚙️ Git-Export Makefile Command:

The command below helps you generate a JSON file, which can then be analyzed by AI to produce a detailed summary.

```bash
# Make git-export scripts executable
chmod +x ./git-export.script.sh

# Show help
make go-help

# Export all commits (no filters)
make git-export

# Export last 5 commits
make git-export-last N=5

# Export commits within a date range
make git-export-range SINCE=2025-07-01 UNTIL=2025-07-28

# Clean up exported JSONs
make git-export-clean
```

---

## 📖 Documentation

For advanced guides, see the [Wiki](https://github.com/Mind2Screen-Dev-Team/thousand-sunny).
To integrate these OpenAPI specs with external tools (e.g., codegen for clients), use the `/openapi.yaml` or `/openapi.json` endpoints directly.