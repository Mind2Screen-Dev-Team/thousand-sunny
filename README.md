
# 🚢 _Thousand Sunny_ - Project Skeleton

Embark on your next adventure with the Thousand Sunny! Inspired by the legendary ship from One Piece, this Go project skeleton is designed for speed, flexibility, and scalability. Whether you’re building small tools or large applications, the Thousand Sunny will guide your journey to success.

![Going Merry](./storage/assets/Thousand-Sunny.webp "Thousand Sunny")

## 🗂 Project Structure

```bash
├── app
│   ├── dependency  # External dependencies or libraries used across the application.
│   ├── module      # Application-specific modules encapsulating core features.
│   └── registry    # Manages the initialization and registration of application components.
├── bin             # Compiled binaries or executables for the application.
├── cmd
│   └── core        # Main entry point for the application.
├── config          # Configuration files for managing application settings.
├── constant        # Application-wide constants to avoid hardcoding values.
├── database
│   ├── migrations  # Database schema migrations for version control.
│   ├── queries     # Custom Generator For SQLC queries for specific operations.
│   └── seeders     # Seed data for initializing or populating the database.
├── gen
│   └── repo        # Auto-generated repository code for data access.
├── internal        # Internal packages for application-specific functionality.
│   ├── asynq       # Handles asynchronous task queues.
│   │   └── handler # Specific handlers for processing tasks.
│   ├── http        # HTTP server and related components.
│   │   ├── handler
│   │   │   ├── health  # Ex. Handlers for application health check endpoints.
│   │   │   ├── user    # Ex. Handlers for user-related endpoints.
│   │   │   └── ....    # Other Handlers.
│   │   ├── middleware  # HTTP middleware for request processing.
│   │   │   ├── global  # Middleware applied to all requests globally.
│   │   │   └── private # Middleware for restricted/private routes.
│   │   └── router      # HTTP routing logic.
│   ├── repo         # Data access layer.
│   │   ├── api      # Repository interfaces for APIs.
│   │   ├── attr     # Repositories for handling attributes.
│   │   └── impl     # Implementation of repository interfaces.
│   └── service      # Business logic layer.
│       ├── api      # API-specific services.
│       ├── attr     # Services for managing attributes.
│       └── impl     # Implementation of service interfaces.
├── pkg             # Utility and reusable packages.
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
    └── logs        # Application log files.
        └── core
            ├── debug  # Debug-level logs.
            ├── io     # Input/output operation logs.
            └── trx    # Transaction logs for auditing or debugging.
```

## 📋 Features

Here's a quick look at what's done and what's still in progress:

### Done ✅
- 🗃️ **Base Structural Directory**: Well-organized code structure to get you started quickly.
- 🔧 **Setup Uber Config**: Configuration uber tool setup.
- 🔧 **Setup Uber Fx**: Uber Dependency injection tool setup.
- 📦 **SQLC Repositories Generator**: Repository generator tools.
- 🌐 **HTTP Handler and Router Loader**: Load and manage routes effortlessly.
- 📜 **DTO Validation**: Validate incoming data with ease.
- 📦 **DB Migrations and Seeders**: Database migration and seeding tools.
- 📄 **Logging**: Integrated logging for better observability.
- 📑 **Makefile Runner**: Simple command runners for building and testing.
- 🐳 **Docker Integration**: Containerize the application.

## 📦 Installation and Setup

To get started with Going-Merry-Go, follow these steps:

```bash
# Clone the repository
git clone git@github.com:Mind2Screen-Dev-Team/thousand-sunny.git

# Navigate to the project directory
cd thousand-sunny

# Install dependencies and set up the project
make setup

# Run the application
make go-run app=core
```

## ⚙️ Makefile Commands

The Makefile provides a set of commands to help you manage and interact with your Go project efficiently. Below is a list of the available commands:

### Setup Commands

- **`make setup`**: Sets up the project by installing necessary tools like `protoc-gen-go`, `protoc-gen-go-grpc`, `goose`, and `pkl-gen-go`.

### Go Commands

- **`make go-tidy`**: Cleans up the `go.mod` file by removing unnecessary dependencies.
- **`make go-run app=<application>`**: Runs the specified application.
- **`make go-build app=<application>`**: Builds the specified application.

### Migration Commands

- **`make migrate-up`**: Migrates the database to the most recent version.
- **`make migrate-up-by-one`**: Migrates the database up by one version.
- **`make migrate-down`**: Rolls back the database version by one.
- **`make migrate-status`**: Displays the migration status of the database.
- **`make migrate-create n=<migration_name> t=<sql|go>`**: Creates a new migration file.

### Seeder Commands

- **`make seeder-up`**: Runs the seeders to populate the database.
- **`make seeder-down`**: Rolls back the seeders by one version.
- **`make seeder-create n=<seeder_name> t=<sql|go>`**: Creates a new seeder file.

### Utility Commands

- **`make print-path`**: Displays the current `PATH` environment variable.
- **`make migrate-help`**: Provides help on migration commands.
- **`make go-help`**: Provides help on Go commands.

### Examples

```bash
# Setup your project workspace
make setup

# Run a Go application (example: restapi)
make go-run app=core

# Migrate the database to the latest version
make migrate-up

```

These commands make it easy to manage your Go application, including its dependencies, database migrations, and proto file generation.

## 📖 Documentation

For detailed documentation and advanced usage, please refer to the [Wiki](https://github.com/Mind2Screen-Dev-Team/thousand-sunny/wiki) page.

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

We welcome contributions! Feel free to submit issues, fork the repository, and send pull requests.

## 🌟 Show Your Support

Give a ⭐️ if you like this project!

## 📧 Contact

For more information or support, you can reach out to us.
