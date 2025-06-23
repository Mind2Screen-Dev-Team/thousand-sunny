#!/bin/bash

# Exit on error, unset variable, or failed pipe
set -euo pipefail

# Show usage instructions if help flag is passed
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
  echo "Usage: $0 [env_file] [new_version|down]"
  exit 0
fi

# Load shared helper functions
if ! source deploy.helper.sh; then
  echo "Error: Failed to source deploy.helper.sh"
  exit 1
fi

# Validate required core env file
if [ ! -f "stack.core.env" ]; then
  echo "Error: File env does not exist: stack.core.env" >&2
  exit 1
fi

# Validate required asynq env file
if [ ! -f "stack.asynq.env" ]; then
  echo "Error: File env does not exist: stack.asynq.env" >&2
  exit 1
fi

# Ensure config file exists
if [ ! -f "config.yaml" ]; then
  echo "Error: File config does not exist: config.yaml" >&2
  exit 1
fi

# Ensure Docker daemon is running
if ! check_docker_running; then
  exit 1
fi

# Load environment variables from the given file (default to stack.core.env)
ENV_FILE="${1:-stack.core.env}"
if [ -f "$ENV_FILE" ]; then
    echo "Loading environment from: $ENV_FILE"
    set -a
    source "$ENV_FILE"
    set +a

    # Validate required variables after sourcing
    echo "Validating required environment variables from $ENV_FILE..."
    : "${APP_DOCKER_NAME:?APP_DOCKER_NAME not set in $ENV_FILE}"
    : "${APP_DOCKER_PORT:?APP_DOCKER_PORT not set in $ENV_FILE}"
    : "${APP_PROJECT_NAME:?APP_PROJECT_NAME not set in $ENV_FILE}"
    : "${APP_NETWORK_NAME:?APP_NETWORK_NAME not set in $ENV_FILE}"
    : "${APP_SERVICE_NAME:?APP_SERVICE_NAME not set in $ENV_FILE}"
else
    echo "Environment file not found: $ENV_FILE"
    exit 1
fi

# Parse the second argument: either version number or 'down'
ACTION="${2:-}"

# If "dc-ctr" is passed, disconnect containers from network and exit
if [[ "$ACTION" == "dc-ctr" ]]; then
    echo "Disconnect container from network mode: disconnecting some containers from current Docker network."
    disconnect_containers_from_network "$APP_NETWORK_NAME" "$APP_DISCONNECT_CONTAINER_DEPENDENCY"
    echo "✅ Successfully disconnected containers: '$APP_DISCONNECT_CONTAINER_DEPENDENCY' from network: '$APP_NETWORK_NAME'."
    exit 0
fi

# Determine Docker Compose and Dockerfile paths
DOCKER_FILE="Dockerfile"
COMPOSE_FILE="docker-compose.$APP_DOCKER_NAME.yml"

# Handle logs command (with optional --follow)
if [[ "$ACTION" == "logs" ]]; then
    FOLLOW_FLAG=""
    if [[ "${f:-0}" == "1" ]]; then
        FOLLOW_FLAG="--follow"
        echo "Following logs for project: $APP_PROJECT_NAME"
    else
        echo "Showing logs for project: $APP_PROJECT_NAME"
    fi

    if is_old_docker; then
        compose_exec \
        --env-file "$ENV_FILE" \
        -f "$COMPOSE_FILE" logs \
        $FOLLOW_FLAG --timestamps app
    else
        compose_exec \
        --env-file "$ENV_FILE" \
        -p "$APP_PROJECT_NAME" \
        -f "$COMPOSE_FILE" logs \
        $FOLLOW_FLAG --timestamps app
    fi
    exit 0
fi

# If "down" is passed, stop containers and exit
if [[ "$ACTION" == "down" ]]; then
    echo "Take-down mode: stopping Docker Compose services only."
    if is_old_docker; then
        compose_exec \
        --env-file "$ENV_FILE" \
        -f "$COMPOSE_FILE" down
    else
        compose_exec \
        --env-file "$ENV_FILE" \
        -p "$APP_PROJECT_NAME" \
        -f "$COMPOSE_FILE" down
    fi

    # Backup current config with timestamp
    backup_config "$APP_SERVICE_VERSION"
    echo "Project $APP_PROJECT_NAME services has been stopped."
    exit 0
fi

# If "down-clean" is passed, stop containers and exit
if [[ "$ACTION" == "down-clean" ]]; then
    echo "Take-down-clean mode: stopping Docker Compose services only."
    if is_old_docker; then
        compose_exec \
        --env-file "$ENV_FILE" \
        -f "$COMPOSE_FILE" down \
        -v --rmi all --remove-orphans
    else
        compose_exec \
        --env-file "$ENV_FILE" \
        -p "$APP_PROJECT_NAME" \
        -f "$COMPOSE_FILE" down \
        -v --rmi all --remove-orphans
    fi

    # Backup current config with timestamp
    backup_config "$APP_SERVICE_VERSION"
    echo "Project $APP_PROJECT_NAME services has been stopped and cleaned."
    exit 0
fi

# Otherwise, treat second argument as a version string
APP_NEW_VERSION="$ACTION"
APP_IMAGE_NAME="$APP_SERVICE_NAME:v$APP_NEW_VERSION"
if [ -z "$APP_NEW_VERSION" ]; then
    echo "Error: New version not provided."
    echo "Usage: $0 [env_file] [new_version]"
    exit 1
fi

# Validate version string format (e.g., v1.2.3)
validate_version "v$APP_NEW_VERSION"

# Create Docker network if missing
create_docker_network "$APP_NETWORK_NAME"

# Attach dependent containers (if any) to the network
connect_containers_to_network "$APP_NETWORK_NAME" "$APP_CONTAINER_DEPENDENCY"

# Update the version inside the .env file
update_env_var APP_SERVICE_VERSION "v$APP_NEW_VERSION" "$ENV_FILE"
echo "Updated 'APP_SERVICE_VERSION' in $ENV_FILE to v$APP_NEW_VERSION"

# Reload environment file after update to apply new version
set -a
source "$ENV_FILE"
set +a

# Check if the Docker image already exists
if docker image inspect "$APP_IMAGE_NAME" > /dev/null 2>&1; then
    echo "Docker image $APP_IMAGE_NAME already exists. Skipping build."
    COMPOSE_BUILD_FLAG=""
else
    echo "Docker image $APP_IMAGE_NAME not found. Building image..."
    COMPOSE_BUILD_FLAG="--build"
fi

# Launch services using the correct Docker Compose version
if is_old_docker; then
    compose_exec \
        --env-file="$ENV_FILE" \
        -f "$COMPOSE_FILE" up \
        -d $COMPOSE_BUILD_FLAG || {
            echo "Error: Failed to start (legacy) docker compose services for $APP_SERVICE_NAME"
            exit 1
        }

    # Backup new config with updated version
    backup_config "v$APP_NEW_VERSION"
    exit 0
fi

# Launch services (Docker Compose v2+)
compose_exec \
    --env-file="$ENV_FILE" \
    -p "$APP_PROJECT_NAME" \
    -f "$COMPOSE_FILE" up \
    -d $COMPOSE_BUILD_FLAG || {
        echo "Error: Failed to start (modern) docker compose services for $APP_SERVICE_NAME"
        exit 1
    }

# Backup new config after successful deployment
backup_config "v$APP_NEW_VERSION"
exit 0
