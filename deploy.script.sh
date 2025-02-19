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

# Function to copy a config file if it doesn't exist
copy_config_if_needed() {
    local source_file="$1"
    local destination_file="$2"
    local file_permission="${3:-644}"  # Default to 644 if no permission is provided

    # Check if both parameters are provided
    if [ -z "$source_file" ] || [ -z "$destination_file" ]; then
        echo "Usage: copy_config_if_needed <source_file> <destination_file>"
        return 1
    fi

    echo "Checking config files..."

    # If the destination file exists, skip copying
    if [ -f "$destination_file" ]; then
        echo "Destination file $destination_file already exists. Skipping copy."
    else
        # Check if the source file exists before copying
        if [ -f "$source_file" ]; then
            cp "$source_file" "$destination_file"  # Copy the source file to the destination
            chmod "$file_permission" "$destination_file"  # Set appropriate file permissions
            echo "Copied $source_file to $destination_file successfully."
        else
            echo "Source file $source_file does not exist."
            return 1
        fi
    fi
}

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

# Function to rebuild a Docker image with parameters
rebuild_image() {
    local image_name="$1"
    local docker_file="$2"

    # Exit if required parameters are missing
    [ -z "$image_name" ] || [ -z "$docker_file" ] && return 1

    echo "Force rebuilding image: $image_name"

    # Remove the Docker image
    docker rmi -f "$image_name" || { echo "Error: Failed to remove image $image_name."; exit 1; }

    # Rebuild the Docker image
    docker build -t "$image_name" -f "$docker_file" . || { echo "Error: Failed to build Docker image."; exit 1; }
}

# Function to strip the 'v' prefix and compare two semantic versions
version_lt() {
    local ver1="${1#v}"
    local ver2="${2#v}"
    [ "$(printf '%s\n' "$ver1" "$ver2" | sort -V | head -n 1)" = "$ver1" ] && [ "$ver1" != "$ver2" ]
}

# Ensure the input version follows the vX.Y.Z sematic versioning format, (X.Y.Z) is sematic versioning and prefix "v" is just for information mean "version"
validate_version() {
    if ! [[ "$1" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "Error: Version must follow the sematic versioning format 'vX.Y.Z' (e.g., v1.0.0). Please refer on this docs: https://semver.org/"
        exit 1
    fi
}

# Function for convert text into title case.
to_title_case() {
  echo "$1" | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2)); print}'
}

# Function to update SERVICE_CORE_VERSION in .env file based on OS
update_service_version_in_env() {
  local SERVICE_VERSION="$1"
  local ENV_FILE="$2"
  local APP_NAME="$3"
  APP_NAME="$(echo "$APP_NAME" | tr '[:lower:]' '[:upper:]')"  # Convert to uppercase

  # Detect the operating system type
  local OS=$(uname)

  # Replace SERVICE_CORE_VERSION value in .env file based on OS
  if [[ "$OS" == "Darwin" ]]; then
      # macOS (BSD sed requires -i with empty string)
      sed -i '' "s/^SERVICE_${APP_NAME}_VERSION=\"[^\"']*\"/SERVICE_${APP_NAME}_VERSION=\"$SERVICE_VERSION\"/" "$ENV_FILE" || { echo "Error: Failed to update $(to_title_case "$APP_NAME") Service Version in .env file."; exit 1; }
      echo "$(to_title_case "$APP_NAME") Service Version updated to $SERVICE_VERSION in .env file (macOS)."

  elif [[ "$OS" == "Linux" ]]; then
      # Linux (GNU sed allows -i without empty string)
      sed -i "s/^SERVICE_${APP_NAME}_VERSION=\"[^\"']*\"/SERVICE_${APP_NAME}_VERSION=\"$SERVICE_VERSION\"/" "$ENV_FILE" || { echo "Error: Failed to update $(to_title_case "$APP_NAME") Service Version in .env file."; exit 1; }
      echo "$(to_title_case "$APP_NAME") Service Version updated to $SERVICE_VERSION in .env file (Linux)."

  elif [[ "$OS" == *"MINGW"* || "$OS" == *"MSYS"* ]]; then
      # Windows (Git Bash or MSYS environments)
      sed -i "s/^SERVICE_${APP_NAME}_VERSION=\"[^\"']*\"/SERVICE_${APP_NAME}_VERSION=\"$SERVICE_VERSION\"/" "$ENV_FILE" || { echo "Error: Failed to update $(to_title_case "$APP_NAME") Service Version in .env file."; exit 1; }
      echo "$(to_title_case "$APP_NAME") Service Version updated to $SERVICE_VERSION in .env file (Windows)."

  else
      echo "Unsupported OS: $OS"
      exit 1
  fi
}
