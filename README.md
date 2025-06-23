
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
│   ├── seeders     # Seed data for initializing or populating the database.
│   │   └── ...     # List of seeder files.
│   └── ...         # database codes.
├── gen
│   └── repo        # Auto-generated repository code for data access.
├── internal        # Internal packages for application-specific functionality.
│   ├── schema      # Collections of schema.
│   │   └── ...     # list of schema.
│   ├── helper      # Collections of short function for helpers.
│   │   └── ...     # list of helpers
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
│   ├── schema       # Defind Schema Database Object Table and column.
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
│   ├── xauth        # Authentication helpers and utilities.
│   ├── xecho        # Extensions for the Echo web framework.
│   ├── xfilter      # Utilities for filtering data in requests.
│   ├── xhttp        # General HTTP helpers and utilities.
│   ├── xlog         # Logging utilities.
│   ├── xmail        # Email helpers and utilities.
│   ├── xpanic       # Panic recovery utilities for error handling.
│   ├── xresp        # Response utilities for standardizing HTTP responses.
│   ├── xsecurity    # Security for encryption and decryption utilities.
│   ├── xtracer      # Open-Telemtry Pkg Helper.
│   └── xvalidate    # Validation Pkg for helper mapping / defind error.
└── storage         # Storage for static files and logs.
    ├── assets      # Static assets like images or documents.
    │   └── ...     # Add other assets here.
    ├── cron        # Cron configuration.
    │   └── ...     # Add other cron configuration here.
    ├── template    # Template files.
    │   └── ...     # Add other template here.
    └── logs        # Application log files.
        └── <server.name> # Based `config.yaml` on server section.
            ├── debug     # Debug-level logs.
            ├── io        # Input/output (incoming logs) operation logs.
            └── trx       # Transaction logs for auditing or debugging.
```

## 📋 Features

Here's a quick look at what's done and what's still in progress:

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
- 🌐 **Open-Telemetry**: Add Tracer, Metric and Logs Configuration.

## 📦 Installation and Setup

To get started, follow these steps:

```bash
# Clone the repository
git clone git@github.com:Mind2Screen-Dev-Team/thousand-sunny.git

# Navigate to the project directory
cd thousand-sunny

# Install dependencies and set up the project
make setup

# Copy example config and fill the value of configuration for deployment.
cp stack.example.env stack.asynq.env
cp stack.example.env stack.core.env

# The `config.yaml` file for application configuration.
cp config.example.yaml config.yaml

# Run LOCAL for simplify step

# Run the application
make go-run a=core

# Run on Docker

# The `config.yaml` file for application configuration.
cp config.example.yaml config.yaml

# Make it script deployment executeable
chmod +x ./deploy.*.sh

# Version must follow the sematic versioning format 'X.Y.Z' (e.g., 1.0.0).
# Please refer on this docs: https://semver.org/

# For deploy up
make deploy-asynq-up v=<version>
make deploy-core-up v=<version>

# For deploy down
make deploy-asynq-down
make deploy-core-down
```

## ⚙️ Makefile Commands

The Makefile provides a set of commands to help you manage and interact with your Go project efficiently. Below is a list of the available commands:

### Setup Commands

- **`make setup`**: Sets up the project by installing necessary tools like `goose` and `sqlc`.

### Sqlc Commands

- **`make sqlc-gen`**: Sets up the project by generating `sqlc` repositories.

### Go Commands

- **`make go-tidy`**: Cleans up the `go.mod` file by removing unnecessary dependencies.
- **`make go-run a=<application>`**: Runs the specified application.
- **`make go-run a=<application> c=<configuration file>`**: Runs the specified application with configuration.
- **`make go-build a=<application>`**: Builds the specified application.

### Utility Commands

- **`make print-path`**: Displays the current `PATH` environment variable.
- **`make go-help`**: Provides help on Go commands.

### Examples

```bash
# Setup your project workspace
make setup

# Run a Go application (example: core, asynq)
make go-run a=asynq
make go-run a=core

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
