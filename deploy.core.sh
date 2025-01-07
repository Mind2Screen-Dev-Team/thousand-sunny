#!/bin/bash

# Check and source /etc/profile if it exists
if [ -f /etc/profile ]; then
    echo "Sourcing /etc/profile..."
    source /etc/profile
else
    echo "/etc/profile not found, skipping..."
fi

# Check and source ~/.profile if it exists
if [ -f ~/.profile ]; then
    echo "Sourcing ~/.profile..."
    source ~/.profile
else
    echo "~/.profile not found, skipping..."
fi

# Check and source ~/.bashrc if it exists
if [ -f ~/.bashrc ]; then
    echo "Sourcing ~/.bashrc..."
    source ~/.bashrc
else
    echo "~/.bashrc not found, skipping..."
fi

# Define source and destination files
SOURCE_FILE="example.config.yaml"
DESTINATION_FILE="config.core.yaml"

# Check if the destination file exists
if [ -f "$DESTINATION_FILE" ]; then
    echo "Destination file $DESTINATION_FILE already exists. Skipping copy."
else
    # Check if the source file exists
    if [ -f "$SOURCE_FILE" ]; then
        # Copy the file
        cp "$SOURCE_FILE" "$DESTINATION_FILE"
        echo "Copied $SOURCE_FILE to $DESTINATION_FILE successfully."
    else
        echo "Source file $SOURCE_FILE does not exist."
        exit 1
    fi
fi

# Check if the Docker network exists
NETWORK_NAME="core_app_net"
if ! docker network inspect $NETWORK_NAME >/dev/null 2>&1; then
  echo "Network $NETWORK_NAME not found, creating it..."
  docker network create --driver bridge $NETWORK_NAME
else
  echo "Network $NETWORK_NAME already exists, skipping creation."
fi

# Check if 'setup' parameter is passed
if [[ "$1" != "setup" ]]; then
  # Migrate DB Up
  make migrate-up || { echo 'Error: migrate up a db migrations.'; exit 1; }

  # Migrate DB Status
  make migrate-status || { echo 'Error: migrate status a db migrations.'; exit 1; }

  # Docker Build Image services
  docker build -t core-app:latest -f Dockerfile.core . || { echo 'Error: build docker build image service.'; exit 1; }

  # Stop Docker Compose services
  docker compose -f compose.core.yml down || { echo 'Error: take-down a service via docker compose.'; exit 1; }

  # Start Docker Compose services with build
  docker compose -f compose.core.yml up -d || { echo 'Error: deploy a service via docker compose.'; exit 1; }
else
  echo "Setup parameter detected. Skipping migrations and Docker Compose commands."
fi