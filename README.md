# 🏴‍☠️ *Thousand Sunny* – Project Skeleton

⚡ **Navigate the Grand Line of Go development** with **Thousand Sunny** – where dreams become deployable reality!

Just like Luffy's legendary vessel that conquered impossible seas, this battle-tested Go skeleton breaks through the storms of complex architecture. Built with the **spirit of adventure** and engineered for **legendary performance**, Thousand Sunny transforms your wildest coding ambitions into production-ready masterpieces.

🌊 From humble microservices to enterprise titans – **every great journey starts with the right ship**.

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
├── gen             # Generation packages
│   ├── sqlc        # Auto-generated SQLC repository code for database access.
│   └── gorm        # Auto-generated GORM repository code for database access.
│       ├── model   # Auto-generated GORM models.
│       └── query   # Auto-generated GORM query code.
├── infra           # Infrastructure packages
│   ├── sdk         # SDK packages.
│   │   ├── <sdk-name>                        # Some SDK packages.
│   │   │   ├── sdk.<sdk-name>.<action>.go    # Some SDK packages.
│   │   │   └── sdk.<sdk-name>.fx.modules.go  # Some SDK Uber Fx modules.
│   │   └── sdk.fx.modules.go                 # Main SDK Uber Fx modules.
│   └── http            # HTTP packages.
│       └── middleware  # HTTP middleware.
│           ├── mdl.<private|global>.<name>.go  # HTTP middleware.
│           └── mdl.fx.modules.go               # HTTP middleware Uber Fx modules.
├── internal                # Internal packages (application-specific).
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

* 🗃️ **Modular Go Structure** – Onion Architecture with Clean Architecture principles and DDD influence (Basically Clean Architecture).
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

### Run By Docker Compose

```bash
# Make it Executeable deploy script
chmod +x deploy.script.sh
chmod +x deploy.sh

# X.Y.Z = Sematic Version, ex: 1.2.3, 0.2.3, 0.0.3
make deploy-core-up v=X.Y.Z
```

### Run by Pure Docker
```bash
# 1. Build the image
docker build -f Dockerfile -t api-core-thousand-sunny:latest .

# 2. Run the container
docker run -d \
  --name api-core-thousand-sunny-app \
  -p 8080:8080 \
  -v $(pwd)/storage:/app/storage \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  --restart unless-stopped \
  api-core-thousand-sunny:latest

# 3. Check if it's running
docker ps

# 4. View logs
docker logs api-core-thousand-sunny-app
```

### Run by K8S With Minikube and Docker

```bash
# start / open your docker
# install minikube
brew install minikube # macos

# enable ingress
minikube addons enable ingress

# enable metric
minikube addons enable metrics-server

# build image
minikube image build -t api-core-thousand-sunny:latest -f Dockerfile .

# open new tab terminal, run minikube mount directory, ex:
#   - /Users/<username>/<some_path>/thousand-sunny:/mnt/thousand-sunny
minikube mount <your_project_path>:/mnt/thousand-sunny

# ------ START: ONLY FOR MAC ------
# add local domain
sudo nano /etc/hosts
# Add
# 127.0.0.1       thousand-sunny.local
sudo dscacheutil -flushcache
sudo killall -HUP mDNSResponder
# ------ END: ONLY FOR MAC ------

# unload deploy
kubectl delete -f k8s-manifest.yml -n internal

# deploy
kubectl apply -f k8s-manifest.yml -n internal 

# open new tab terminal, run minikube dashboard
minikube dashboard

# open new terminal, and run tunnel into service
minikube tunnel

# open http://thousand-sunny.local/health
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
# Make it Executeable git-export script
chmod +x git-export.script.sh

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

## 📖 Contact Us
Visit our website, [mindtoscreen.com](https://mindtoscreen.com/).
