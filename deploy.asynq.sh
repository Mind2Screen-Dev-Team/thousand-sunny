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

# Docker-related operations (only execute if Docker is running)
if [ "$DOCKER_RUNNING" = true ]; then
    # Check if the Docker network exists
    NETWORK_NAME="asynq_app_net"
    if ! docker network inspect $NETWORK_NAME >/dev/null 2>&1; then
        echo "Network $NETWORK_NAME not found, creating it..."
        docker network create --driver bridge $NETWORK_NAME
    else
        echo "Network $NETWORK_NAME already exists, skipping creation."
    fi
else
    echo "Skipping Docker network creation as Docker is not running."
fi

# Check if 'setup' parameter is passed
if [[ "$1" != "setup" ]]; then
    # Migrate DB Up
    make migrate-up || { echo 'Error: Failed to migrate DB up.'; exit 1; }

    # Migrate DB Status
    make migrate-status || { echo 'Error: Failed to get migration status.'; exit 1; }

    # Docker-related commands (only if Docker is running)
    if [ "$DOCKER_RUNNING" = true ]; then
        # Docker Build Image services
        docker build -t asynq-app:latest -f Dockerfile.asynq . || { echo 'Error: Failed to build Docker asynq image.'; exit 1; }

        # Stop Docker Compose services
        docker compose -p asynq-app-stack -f compose.asynq.yml down

        # Start Docker Compose services with build
        docker compose -p asynq-app-stack -f compose.asynq.yml up -d || { echo 'Error: Failed to start Docker Compose asynq services.'; exit 1; }
    else
        echo "Skipping Docker Compose commands as Docker is not running."
    fi
else
    echo "Setup parameter detected. Skipping migrations and Docker Compose commands."
fi
