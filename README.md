
# 🏴‍☠️ _Thousand Sunny_ 🏴‍☠️ - Project Skeleton

Embark on your next adventure with the Thousand Sunny! Inspired by the legendary ship from One Piece, this Go project skeleton is designed for speed, flexibility, and scalability. Whether you’re building small tools or large applications, the Thousand Sunny will guide your journey to success.

![Thousand-Sunny-Image](./storage/assets/thousand-sunny.png "Thousand Sunny")

## 🗂 Project Structure

```bash
├── app
│   ├── dependency  # External dependencies or libraries used across the application.
│   │   └── ...     # Create Dependency Here.
│   ├── module      # Application-specific modules encapsulating core features.
│   │   └── ...     # Create Module Here.
│   └── registry    # Manages the initialization and registration of application components.
│       └── ...     # Create Registry Here.
├── bin             # Compiled binaries or executables for the application.
├── cmd
│   ├── asynq       # Main entry point for the 'asynq' application.
│   └── core        # Main entry point for the 'core' application.
├── config          # Configuration files for managing application settings.
├── constant        # Application-wide constants to avoid hardcoding values.
├── database
│   ├── migrations  # Database schema migrations for version control.
│   │   └── ...     # List of migration files.
│   ├── queries     # Custom Generator For SQLC queries for specific operations.
│   │   └── ...     # List of sqlc generator queries.
│   └── seeders     # Seed data for initializing or populating the database.
│       └── ...     # List of seeder files.
├── gen
│   └── repo        # Auto-generated repository code for data access.
├── internal        # Internal packages for application-specific functionality.
│   ├── asynq       # Handles asynchronous task queues.
│   │   ├── router      # Base Asynq routing configuration (No need to add Something here).
│   │   ├── worker      # Specific handlers for processing workers.
│   │   │   └── ...     # Other Worker Routing Handlers.
│   │   └── scheduler   # Specific handlers for processing schedulers.
│   │       └── ...     # Other Scheduler Routing Handlers.
│   ├── http        # HTTP server and related components.
│   │   ├── handler
│   │   │   ├── health  # Example. Handlers for health-related endpoints.
│   │   │   ├── user    # Example. Handlers for user-related endpoints.
│   │   │   └── ...     # Other Routing Handlers.
│   │   ├── middleware  # HTTP middleware for request processing.
│   │   │   ├── global  # Middleware applied to all requests globally.
│   │   │   │   └── ... # Other Global Middleware.
│   │   │   └── private # Middleware for restricted/private routes.
│   │   │       └── ... # Other Private Middleware.
│   │   └── router      # Base HTTP routing configuration (No need to add Something here).
│   ├── provider     # External Provider Data access layer.
│   │   ├── api      # External Provider interfaces for APIs.
│   │   │   └── ...  # Other API External Provider.
│   │   ├── attr     # External Providers for handling attributes.
│   │   │   └── ...  # Other Attribute External Provider.
│   │   └── impl     # Implementation of repository interfaces.
│   │       └── ...  # Other Implementaion External Provider.
│   ├── repo         # Data access layer.
│   │   ├── api      # Repository interfaces for APIs.
│   │   │   └── ...  # Other API Repository.
│   │   ├── attr     # Repositories for handling attributes.
│   │   │   └── ...  # Other Attribute Repository.
│   │   └── impl     # Implementation of repository interfaces.
│   │       └── ...  # Other Implementaion Repository.
│   └── service      # Business logic layer.
│       ├── api      # API-specific services.
│       │   └── ...  # Other API Service.
│       ├── attr     # Services for managing attributes.
│       │   └── ...  # Other Attribute Service.
│       └── impl     # Implementation of service interfaces.
│           └── ...  # Other Implementaion Service.
├── pkg             # Utility and reusable packages.
│   ├── xasynq       # Asynq sever helpers and utilities.
│   ├── xmail        # Email helpers and utilities.
│   ├── xauth        # Authentication helpers and utilities.
│   ├── xecho        # Extensions for the Echo web framework.
│   ├── xfilter      # Utilities for filtering data in requests.
│   ├── xhttp        # General HTTP helpers and utilities.
│   ├── xlog         # Logging utilities.
│   ├── xpanic       # Panic recovery utilities for error handling.
│   ├── xresp        # Response utilities for standardizing HTTP responses.
│   └── xrsa         # RSA encryption and decryption utilities.
└── storage         # Storage for static files and logs.
    ├── assets      # Static assets like images or documents.
    │   └── ...     # Add other assets here.
    ├── template    # Template files.
    │   └── ...     # Add other template here.
    └── logs        # Application log files.
        ├── asynq
        │   ├── debug   # Debug-level logs.
        │   ├── io      # Input/output (incoming logs) operation logs.
        │   └── trx     # Transaction logs for auditing or debugging.
        └── core
            ├── debug   # Debug-level logs.
            ├── io      # Input/output (incoming logs) operation logs.
            └── trx     # Transaction logs for auditing or debugging.
```

## 📋 Features

Here's a quick look at what's done and what's still in progress:

### Done ✅
- 🗃️ **Base Structural Directory**: Well-organized code structure to get you started quickly.
- 🔧 **Setup Uber Fx**: Uber Dependency injection tool setup.
- 🔧 **Setup Uber Config**: Uber Configuration tool setup.
- 📦 **SQLC Repositories Generator**: Repository generator tools.
- 🌐 **Asynq Redis Queue Worker and Scheduler Handler and Router Loader**: Load and manage routes effortlessly.
- 🌐 **HTTP Handler and Router Loader**: Load and manage routes effortlessly.
- 📜 **DTO Validation**: Validate incoming data with ease.
- 📦 **DB Migrations and Seeders**: Database migration and seeding tools.
- 📄 **Logging**: Integrated logging for better observability.
- 📑 **Makefile Runner**: Simple command runners for building and testing.
- 🐳 **Docker Integration**: Containerize the application.
- 🌐 **Open-Telemetry**: Add Tracer and Metric Configuration.

## 📦 Installation and Setup

To get started, follow these steps:

```bash
# Clone the repository
git clone git@github.com:Mind2Screen-Dev-Team/thousand-sunny.git

# Navigate to the project directory
cd thousand-sunny

# Install dependencies and set up the project
make setup

# Copy example config and fill the value of configuration
# The `.env` file for deployment time and sql db migrations:
#   + NOTE: There is a some key no need to fill a value, that is:
#     - SERVICE_CORE_VERSION
#     - SERVICE_ASYNQ_VERSION
cp .example.env .env

# The `config.yaml` file for application configuration.
cp config.example.yaml config.yaml

# Run LOCAL for simplify step

# Run the application
make go-run a=core

# Run On PROD / LOCAL and deploy into Docker

# The `config.yaml` file for application configuration.
cp config.example.yaml config.yaml

# Make it script deployment executeable
chmod +x ./deploy.sh
chmod +x ./deploy.*.sh

# Note For (AUTOMATICALLY) Set Version:
# Please refer on this docs: https://semver.org/
#   - [major]: increment by 1 major version (major reset existing minor and patch version)). [ex. before -> v1.1.1 -> after -> v2.0.0]
#   - [minor]: increment by 1 minor version (minor reset existing patch version, and do not reset existing major version). [ex. before -> v1.1.91 -> after -> v1.2.0]
#   - [patch]: increment by 1 patch version (patch do not reset existing major and minor version). [ex. before -> v0.1.24 -> after -> v0.1.25]

# For deploy core app
make deploy-core v=major
make deploy-core v=minor
make deploy-core v=patch

# For deploy asynq app
make deploy-asynq v=major
make deploy-asynq v=minor
make deploy-asynq v=patch

# Note For (MANUAL) Set Version:
# Version must follow the sematic versioning format 'X.Y.Z' (e.g., 1.0.0).
# Please refer on this docs: https://semver.org/

# For force re-build docker image to deploy core app, ex: 0.0.1
make deploy-core-rebuild <version>

# For force re-build docker image to deploy asynq app, ex: 0.0.1
make deploy-asynq-rebuild <version>
```

## ⚙️ Makefile Commands

The Makefile provides a set of commands to help you manage and interact with your Go project efficiently. Below is a list of the available commands:

### Note:
- Please escape strings in environment variables in file .env when handling errors, especially when using migrations.

### Setup Commands

- **`make setup`**: Sets up the project by installing necessary tools like `goose` and `sqlc`.

### Go Commands

- **`make go-tidy`**: Cleans up the `go.mod` file by removing unnecessary dependencies.
- **`make go-run a=<application>`**: Runs the specified application.
- **`make go-run a=<application> c=<configuration file>`**: Runs the specified application with configuration.
- **`make go-build a=<application>`**: Builds the specified application.

### Migration Commands

- **`make migrate-up`**: Migrates the database to the most recent version.
- **`make migrate-up-by-one`**: Migrates the database up by one version.
- **`make migrate-down`**: Rolls back the database version by one.
- **`make migrate-status`**: Displays the migration status of the database.
- **`make migrate-create n=<migration_name>`**: Creates a new migration file.

### Seeder Commands

- **`make seeder-up`**: Runs the seeders to populate the database.
- **`make seeder-down`**: Rolls back the seeders by one version.
- **`make seeder-create n=<seeder_name>`**: Creates a new seeder file.

### Utility Commands

- **`make print-path`**: Displays the current `PATH` environment variable.
- **`make migrate-help`**: Provides help on migration commands.
- **`make go-help`**: Provides help on Go commands.

### Examples

```bash
# Setup your project workspace
make setup

# Run a Go application (example: core, asynq)
make go-run a=core

# Migrate the database to the latest version
make migrate-up

```

These commands make it easy to manage your Go application, including its dependencies, database migrations, and proto file generation.

## 📖 Documentation

For detailed documentation and advanced usage, please refer to the [Wiki](https://github.com/Mind2Screen-Dev-Team/thousand-sunny) page.

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

We welcome contributions! Feel free to submit issues, fork the repository, and send pull requests.

## 🌟 Show Your Support

Give a ⭐️ if you like this project!

## 📧 Contact

For more information or support, you can reach out to us.
