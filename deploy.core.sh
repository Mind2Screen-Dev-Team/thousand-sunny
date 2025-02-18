#!/bin/bash

# Add deploy script
source deploy.script.sh

PARAM="$1"
REBUILD_PARAM="$2"

SOURCE_FILE="config.example.yaml"
DESTINATION_FILE="config.yaml"

ENV_FILE=".env"
APP_ENV=$(get_env_value "$ENV_FILE" "APP_ENV" "dev")

APP_NAME="core"
APP_SERVICE_NAME="$APP_NAME-$APP_ENV-app"
APP_STACK_NAME="$APP_SERVICE_NAME-stack"
APP_IMAGE_NAME="$APP_SERVICE_NAME:$PARAM"
APP_NETWORK_NAME="${APP_NAME}_${APP_ENV}_app_net"

DOCKER_FILE="Dockerfile.${APP_NAME}"
COMPOSE_FILE="compose.${APP_NAME}.yml"

# Check if Docker is running
check_docker_running

# Source profiles
source_profile "/etc/profile"
source_profile "$HOME/.profile"
source_profile "$HOME/.bashrc"

# Check and create if the Docker network exists
create_docker_network "$APP_NETWORK_NAME"

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
        EXISTING_IMAGE_VERSION=$(docker images --format "{{.Tag}}" $APP_SERVICE_NAME | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -n 1)

        # Skip version comparison if using 'latest'
        if [ "$PARAM" != "latest" ]; then
            if [ -n "$EXISTING_IMAGE_VERSION" ]; then
                echo "Existing $APP_SERVICE_NAME image version: $EXISTING_IMAGE_VERSION"
                if version_lt "$PARAM" "$EXISTING_IMAGE_VERSION"; then
                    echo "Warning: Provided version ($PARAM) is lower than existing version ($EXISTING_IMAGE_VERSION). Aborting build."
                    exit 1
                fi
            else
                echo "No existing $APP_SERVICE_NAME image found."
            fi
        else
            echo "Using 'latest' tag, skipping version comparison."
        fi

        # Build docker images
        if [ -z "$EXISTING_IMAGE_VERSION" ]; then
            # If image doesn't exist, build it
            echo "No existing $APP_SERVICE_NAME image found. Building the image."
            docker build -t "$APP_IMAGE_NAME" -f "$DOCKER_FILE" . || { echo "Error: Failed to build Docker image."; exit 1; }
        elif [[ "$REBUILD_PARAM" == "rebuild" ]]; then
            # If 'rebuild' param is passed, remove and rebuild the image
            rebuild_image "$APP_IMAGE_NAME" "$DOCKER_FILE"
        else
            # Skip rebuilding if the image exists and no 'rebuild' is passed
            echo "Image $APP_IMAGE_NAME already exists. Skipping rebuild."
        fi

        # Stop Docker Compose services
        docker compose -p "$APP_STACK_NAME" -f "$COMPOSE_FILE" down

        # Call the function with the parameter
        update_service_version_in_env "$PARAM" "$ENV_FILE"

        # Start Docker Compose services with build
        docker compose --env-file "$ENV_FILE" -p "$APP_STACK_NAME" -f "$COMPOSE_FILE" up -d || { echo "Error: Failed to start Docker Compose $APP_SERVICE_NAME services."; exit 1; }
    else
        echo "Skipping Docker Compose commands as Docker is not running."
    fi
else
    echo "Setup parameter detected. Skipping migrations and Docker Compose commands."
fi