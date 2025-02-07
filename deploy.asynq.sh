#!/bin/bash

# Function to check if Docker is running
check_docker_running() {
    if ! docker info >/dev/null 2>&1; then
        echo "Warning: Docker is not running. Skipping Docker-related commands."
        DOCKER_RUNNING=false
    else
        DOCKER_RUNNING=true
    fi
}

# Source environment profiles if they exist
source_profile() {
    local PROFILE_FILE=$1
    if [ -f "$PROFILE_FILE" ]; then
        echo "Sourcing $PROFILE_FILE..."
        source "$PROFILE_FILE"
    else
        echo "$PROFILE_FILE not found, skipping..."
    fi
}

# Start of the script
echo "Starting source script..."

# Source profiles
source_profile "/etc/profile"
source_profile "$HOME/.profile"
source_profile "$HOME/.bashrc"

# Define source and destination files
SOURCE_FILE="config.example.yaml"
DESTINATION_FILE="config.asynq.yaml"

# Check and copy the config file if needed
if [ -f "$DESTINATION_FILE" ]; then
    echo "Destination file $DESTINATION_FILE already exists. Skipping copy."
else
    if [ -f "$SOURCE_FILE" ]; then
        cp "$SOURCE_FILE" "$DESTINATION_FILE"
        echo "Copied $SOURCE_FILE to $DESTINATION_FILE successfully."
    else
        echo "Source file $SOURCE_FILE does not exist."
    fi
fi

# Check if Docker is running
check_docker_running

# Check and create if the Docker network exists
create_docker_network() {
  local network_name="$1"

  # Check if Docker is running before proceeding
  if [ "$DOCKER_RUNNING" = true ]; then
    # Validate that a network name has been provided
    if [ -z "$network_name" ]; then
      echo "Error: Network name parameter missing."
      return 1
    fi

    # Check if the Docker network exists
    if ! docker network inspect "$network_name" >/dev/null 2>&1; then
      echo "Network $network_name not found, creating it..."
      docker network create --driver bridge "$network_name"
    else
      echo "Network $network_name already exists, skipping creation."
    fi
  else
    echo "Skipping Docker network creation as Docker is not running."
  fi
}

# Function to get APP_ENV value from the .env file
get_env_value() {
  local env_file="$1"
  local key="$2"
  local default_value="$3"

  if [ ! -f "$env_file" ]; then
    echo "Warning: $env_file file not found. Defaulting to '$default_value'."
    echo "$default_value"
    return 0
  fi

  local value
  value=$(grep "^$key=" "$env_file" | cut -d '=' -f2- | tr -d '"' | tr -d "'")

  if [ -n "$value" ]; then
    echo "$value"
  else
    echo "Warning: $key not set in $env_file. Defaulting to '$default_value'."
    echo "$default_value"
  fi
}

PARAM="$1"
REBUILD_PARAM="$2"
ENV_FILE=".env"
DOCKER_FILE="Dockerfile.asynq"
COMPOSE_FILE="compose.asynq.yml"

APP_NAME="asynq"
APP_ENV=$(get_env_value "$ENV_FILE" "APP_ENV" "dev")
SERVICE_NAME="$APP_NAME-$APP_ENV-app"
APP_STACK_NAME="$SERVICE_NAME-stack"
IMAGE_NAME="$SERVICE_NAME:$PARAM"

NETWORK_NAME="${APP_NAME}_${APP_ENV}_app_net"

# Check and create if the Docker network exists
create_docker_network "$NETWORK_NAME"

# Ensure the input version follows the vX.Y.Z sematic versioning format, (X.Y.Z) is sematic versioning and prefix "v" is just for information mean "version"
validate_version() {
    if ! [[ "$1" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "Error: Version must follow the sematic versioning format 'vX.Y.Z' (e.g., v1.0.0). Please refer on this docs: https://semver.org/"
        exit 1
    fi
}

# Function to rebuild the image
rebuild_image() {
    echo "Force rebuild image."

    # Remove the image with tag param version
    docker rmi -f "$IMAGE_NAME" || { echo "Error: Failed to remove image $IMAGE_NAME."; exit 1; }

    # Rebuild image with tag param version
    docker build -t "$IMAGE_NAME" -f "$DOCKER_FILE" . || { echo "Error: Failed to build Docker image."; exit 1; }
}

# Function to strip the 'v' prefix and compare two semantic versions
version_lt() {
    local ver1="${1#v}"
    local ver2="${2#v}"
    [ "$(printf '%s\n' "$ver1" "$ver2" | sort -V | head -n 1)" = "$ver1" ] && [ "$ver1" != "$ver2" ]
}

# Function to update SERVICE_ASYNQ_VERSION in .env file based on OS
update_service_version_in_env() {
  local SERVICE_VERSION="$1"
  local ENV_FILE="$2"

  # Detect the operating system type
  local OS=$(uname)

  # Replace SERVICE_ASYNQ_VERSION value in .env file based on OS
  if [[ "$OS" == "Darwin" ]]; then
      # macOS (BSD sed requires -i with empty string)
      sed -i '' "s/^SERVICE_ASYNQ_VERSION=\"[^\"']*\"/SERVICE_ASYNQ_VERSION=\"$SERVICE_VERSION\"/" "$ENV_FILE" || { echo "Error: Failed to update Asynq Service Version in .env file."; exit 1; }
      echo "Asynq Service Version updated to $SERVICE_VERSION in .env file (macOS)."

  elif [[ "$OS" == "Linux" ]]; then
      # Linux (GNU sed allows -i without empty string)
      sed -i "s/^SERVICE_ASYNQ_VERSION=\"[^\"']*\"/SERVICE_ASYNQ_VERSION=\"$SERVICE_VERSION\"/" "$ENV_FILE" || { echo "Error: Failed to update Asynq Service Version in .env file."; exit 1; }
      echo "Asynq Service Version updated to $SERVICE_VERSION in .env file (Linux)."

  elif [[ "$OS" == *"MINGW"* || "$OS" == *"MSYS"* ]]; then
      # Windows (Git Bash or MSYS environments)
      sed -i "s/^SERVICE_ASYNQ_VERSION=\"[^\"']*\"/SERVICE_ASYNQ_VERSION=\"$SERVICE_VERSION\"/" "$ENV_FILE" || { echo "Error: Failed to update Asynq Service Version in .env file."; exit 1; }
      echo "Asynq Service Version updated to $SERVICE_VERSION in .env file (Windows)."

  else
      echo "Unsupported OS: $OS"
      exit 1
  fi
}

# Check if 'setup' parameter is passed
if [[ "$PARAM" != "setup" ]]; then
    # Skip version checker if using 'latest'
    if [ "$PARAM" != "latest" ]; then
        # Validate the input version
        validate_version "$PARAM"
    fi

    # Migrate DB Up
    make migrate-up || { echo 'Error: Failed to migrate DB up.'; exit 1; }

    # Migrate DB Status
    make migrate-status || { echo 'Error: Failed to get migration status.'; exit 1; }

    # Docker-related commands (only if Docker is running)
    if [ "$DOCKER_RUNNING" = true ]; then
        # Check existing image version
        EXISTING_IMAGE_VERSION=$(docker images --format "{{.Tag}}" $SERVICE_NAME | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -n 1)

        # Skip version comparison if using 'latest'
        if [ "$PARAM" != "latest" ]; then
            if [ -n "$EXISTING_IMAGE_VERSION" ]; then
                echo "Existing $SERVICE_NAME image version: $EXISTING_IMAGE_VERSION"
                if version_lt "$PARAM" "$EXISTING_IMAGE_VERSION"; then
                    echo "Warning: Provided version ($PARAM) is lower than existing version ($EXISTING_IMAGE_VERSION). Aborting build."
                    exit 1
                fi
            else
                echo "No existing $SERVICE_NAME image found."
            fi
        else
            echo "Using 'latest' tag, skipping version comparison."
        fi

        # Build docker images
        if [ -z "$EXISTING_IMAGE_VERSION" ]; then
            # If image doesn't exist, build it
            echo "No existing $SERVICE_NAME image found. Building the image."
            docker build -t "$IMAGE_NAME" -f "$DOCKER_FILE" . || { echo "Error: Failed to build Docker image."; exit 1; }
        elif [[ "$REBUILD_PARAM" == "rebuild" ]]; then
            # If 'rebuild' param is passed, remove and rebuild the image
            rebuild_image
        else
            # Skip rebuilding if the image exists and no 'rebuild' is passed
            echo "Image $IMAGE_NAME already exists. Skipping rebuild."
        fi

        # Stop Docker Compose services
        docker compose -p "$APP_STACK_NAME" -f "$COMPOSE_FILE" down

        # Call the function with the parameter
        update_service_version_in_env "$PARAM" "$ENV_FILE"

        # Start Docker Compose services with build
        docker compose --env-file "$ENV_FILE" -p "$APP_STACK_NAME" -f "$COMPOSE_FILE" up -d || { echo "Error: Failed to start Docker Compose $SERVICE_NAME services."; exit 1; }
    else
        echo "Skipping Docker Compose commands as Docker is not running."
    fi
else
    echo "Setup parameter detected. Skipping migrations and Docker Compose commands."
fi