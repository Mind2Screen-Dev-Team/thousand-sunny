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
│   ├── core        # Main entry point for the 'core' application.
│   ├── gorm        # GORM model/query generator.
│   └── migrate     # Goose migration/seeder CLI runner.
├── config          # Configuration files for application settings.
├── constant        # Centralized constants to avoid hardcoding.
├── database
│   ├── migrations  # Schema migrations (version-controlled).
│   ├── queries     # SQLC query generators for custom DB operations.
│   └── seeders     # Seed data for initial or demo setups.
├── gen
│   ├── sqlc        # Auto-generated SQLC repository code for database access.
│   └── gorm        # Auto-generated GORM repository code for database access.
│       ├── model   # Auto-generated GORM models.
│       └── query   # Auto-generated GORM query code.
├── internal        # Internal packages (application-specific).
│   ├── <domain>            # Domain modules.
│   │   ├── <sub-domain>    # Sub-domains within a domain.
│   │   │   ├── <domain>.<sub-domain>.<action>.handler.go     # Endpoint handlers.
│   │   │   ├── <domain>.<sub-domain>.repo.go                 # Data access layer.
│   │   │   ├── <domain>.<sub-domain>.service.go              # Business logic.
│   │   │   └── <domain>.<sub-domain>.fx.modules.go           # Uber Fx modules.
│   │   ├── <domain>.<action>.handler.go                   # Endpoint handlers.
│   │   ├── <domain>.repo.go                               # Data access layer.
│   │   ├── <domain>.service.go                            # Business logic.
│   │   └── <domain>.fx.modules.go                         # Uber Fx modules.
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

* 🗃️ **Modular Go Structure** – Clean architecture.
* 🌐 **Huma Framework** – Auto-generates **OpenAPI** specifications.
* 📜 **Live Docs at `/docs`** – Interactive Swagger UI.
* 📂 **OpenAPI Specs** – Download at:

  * `http://<host>:<port>/openapi.yaml`
  * `http://<host>:<port>/openapi.json`
* 🔧 **Uber Fx & Config** – Dependency management.
* 📦 **SQLC Gen** – SQLC DB code generation.
* 📦 **GORM Gen** – GORM DB code generation.
* 🐘 **Goose Migrations & Seeders** – Easy DB version control.
* 🐳 **Docker-Ready** – Containerized for dev/prod.
* 🌐 **OpenTelemetry** – Observability with traces, metrics, and logs.

---

## 🚀 Getting Started

```bash
git clone git@github.com:Mind2Screen-Dev-Team/thousand-sunny.git
cd thousand-sunny

# Install tools and dependencies (Goose, jq, etc.)
make setup

# Copy env & config
cp stack.example.env stack.core.env
cp config.example.yaml config.yaml
```

### Run Locally

```bash
make go-run a=core
```

### Run By Docker

```bash
# X.Y.Z = Sematic Version, ex: 1.2.3, 0.2.3, 0.0.3
make deploy-core-up v=X.Y.Z
```

### Generate Code Repository Queries by SQLC

```bash
make sqlc-gen
```

### Generate ORM Models & Queries by GORM

```bash
make gorm-gen
```

---

### Database Migrations & Seeders (Quick Reference)

| Command                        | Description                                               |
| ------------------------------ | --------------------------------------------------------- |
| `make migrate-up`              | Apply all new migrations                                  |
| `make migrate-down`            | Rollback the last migration                               |
| `make migrate-up-to v=1234`    | Migrate **up** to version `1234` (use `v=` param)         |
| `make migrate-down-to v=1234`  | Rollback **down** to version `1234` (use `v=` param)      |
| `make migrate-status`          | Show migration status                                     |
| `make migrate-version`         | Print current DB version                                  |
| `make migrate-create n=xyz`    | Create a new migration file with name `datetime_xyz.sql`  |
| `make migrate-fix`             | Fix migration filenames (zero-padding)                    |
| `make seed-up`                 | Apply all new seeders                                     |
| `make seed-down`               | Rollback the last seeder                                  |
| `make seed-up-to v=1234`       | Seed **up** to version `1234` (use `v=` param)            |
| `make seed-down-to v=1234`     | Rollback **down** to version `1234` (use `v=` param)      |
| `make seed-status`             | Show seeder status                                        |
| `make seed-version`            | Print current seeder version                              |
| `make seed-create n=xyz`       | Create a new seeder file with name `datetime_xyz.sql`     |
| `make seed-fix`                | Fix seeder filenames (zero-padding)                       |

---

### Git Export
#### This will help you to summarize `Commit / Merge / Pull Request` feed a commit changes into AI like ChatGPT.

```bash
# Export all commits (filters: s=since, u=until, l=limit)
make git-export-all s=2025-07-01 u=2025-07-28 l=10

# Export last N commits
make git-export-last n=5

# Export commits by date range
make git-export-range s=2025-07-01 u=2025-07-28

# Clean exported JSON files
make git-export-clean
```

---

### Access Docs

* Docs UI: `http://<host>:<port>/docs`
* OpenAPI YAML: `http://<host>:<port>/openapi.yaml`
* OpenAPI JSON: `http://<host>:<port>/openapi.json`

## 📖 Documentation

For advanced guides, see the Wiki. To integrate these OpenAPI specs with external tools (e.g., codegen for clients), use the `/openapi.yaml` or `/openapi.json` endpoints directly.
