#!/bin/bash

# ---------------------------------------------------------------------------
# Function: check_docker_running
# Description: Verifies if Docker is currently running.
# Returns:
#   - 0 if Docker is running.
#   - 1 if Docker is not running.
# ---------------------------------------------------------------------------
check_docker_running() {
    if ! docker info >/dev/null 2>&1; then
        echo "Error: Docker is not running. Please start Docker and try again."
        return 1
    fi
    return 0
}

# ---------------------------------------------------------------------------
# Function: is_old_docker
# Description: Checks if the system is using legacy `docker-compose`
#              instead of the modern `docker compose` plugin.
# Returns:
#   - 0 if legacy `docker-compose` is available (old Docker).
#   - 1 if modern `docker compose` is available.
# ---------------------------------------------------------------------------
is_old_docker() {
    if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
        return 1  # Using modern docker CLI
    fi
    return 0  # Legacy docker-compose available
}

# ---------------------------------------------------------------------------
# Function: compose_exec
# Description: Executes docker-compose commands using the appropriate
#              binary (modern `docker compose` or legacy `docker-compose`)
# Arguments:
#   - All arguments are passed directly to the compose command.
# ---------------------------------------------------------------------------
compose_exec() {
  if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
    docker compose "$@"
  else
    docker-compose "$@"
  fi
}

# ---------------------------------------------------------------------------
# Function: create_docker_network
# Description: Ensures a Docker network exists. Creates it if missing.
# Arguments:
#   - $1: Name of the Docker network to check/create.
# Notes:
#   - Only runs if Docker is running and the network name is provided.
# ---------------------------------------------------------------------------
create_docker_network() {
  local network_name="$1"

  if ! check_docker_running; then
    echo "Skipping Docker network creation as Docker is not running."
    return 0
  fi

  if [ -z "$network_name" ]; then
    echo "Error: Network name parameter missing."
    return 1
  fi

  if ! docker network inspect "$network_name" >/dev/null 2>&1; then
    echo "Network $network_name not found, creating it..."
    docker network create --driver bridge "$network_name"
    return 0
  fi

  echo "Network $network_name already exists, skipping creation."
  return 0
}

# ---------------------------------------------------------------------------
# Function: update_env_var
# Description: Updates an existing key or inserts a new key=value pair
#              into a .env file. Cross-platform (Linux/macOS) compatible.
# Arguments:
#   - $1: Key name to update.
#   - $2: New value for the key.
#   - $3: Target .env file.
# ---------------------------------------------------------------------------
update_env_var() {
  local key="$1"
  local new_value="$2"
  local env_file="$3"

  if grep -q "^$key=" "$env_file"; then
    if [[ "$(uname)" == "Darwin" ]]; then
      sed -i '' -E "s|^($key=).*|\1$new_value|" "$env_file"
    else
      sed -i -E "s|^($key=).*|\1$new_value|" "$env_file"
    fi
  else
    echo "$key=$new_value" >> "$env_file"
  fi
}

# ---------------------------------------------------------------------------
# Function: validate_version
# Description: Validates that a given version string follows semantic
#              versioning format prefixed with 'v' (e.g., v1.2.3).
# Arguments:
#   - $1: Version string to validate.
# Exits:
#   - If invalid, prints an error and exits with status 1.
# ---------------------------------------------------------------------------
validate_version() {
    if ! [[ "$1" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "Error: Version must follow the semantic versioning format 'vX.Y.Z' (e.g., v1.0.0)."
        echo "Please refer to the documentation: https://semver.org/"
        exit 1
    fi
}

# -----------------------------------------------------------------------------
# Function: connect_containers_to_network
# Description: Connects multiple Docker containers to a single Docker network.
# Arguments:
#   - $1: Single network name
#   - $2: Comma-separated list of container names or IDs
# -----------------------------------------------------------------------------
connect_containers_to_network() {
  local network="$1"
  local containers_csv="$2"

  if [ -z "$network" ]; then
    echo "No network provided, skipping."
    return 0
  fi

  if [ -z "$containers_csv" ]; then
    echo "No containers provided, skipping."
    return 0
  fi

  # Split containers CSV into an array
  IFS=',' read -ra containers <<< "$containers_csv"

  # Loop over each container and connect to the network
  for container in "${containers[@]}"; do
    # Skip if container is empty or unset
    if [ -z "$container" ]; then
      continue
    fi

    # Try to connect but don't exit on failure
    if docker network connect "$network" "$container"; then
      echo "Connecting docker container '$container' to network '$network'..."
      continue
    fi
  done
}

# -----------------------------------------------------------------------------
# Function: disconnect_containers_from_network
# Description: Disconnects multiple Docker containers from a single Docker network.
# Arguments:
#   - $1: Single network name
#   - $2: Comma-separated list of container names or IDs
# -----------------------------------------------------------------------------
disconnect_containers_from_network() {
  local network="$1"
  local containers_csv="$2"

  if [ -z "$network" ]; then
    echo "No network provided, skipping."
    return 0
  fi

  if [ -z "$containers_csv" ]; then
    echo "No containers provided, skipping."
    return 0
  fi

  # Split containers CSV into an array
  IFS=',' read -ra containers <<< "$containers_csv"

  # Loop over each container and disconnect from the network
  for container in "${containers[@]}"; do
    # Skip if container is empty or unset
    if [ -z "$container" ]; then
      continue
    fi

    # Try to disconnect but don't exit on failure
    if docker network disconnect "$network" "$container"; then
      echo "Disconnected docker container '$container' from network '$network'."
    else
      echo "Warning: Failed to disconnect '$container' from network '$network'."
    fi
  done
}

# -----------------------------------------------------------------------------
# Function: backup_config
# Description: Creates a timestamped backup of the config.yaml file, appending
#              the version with dots replaced by dashes.
# Arguments:
#   - $1: Version string (e.g., "1.2.3")
# -----------------------------------------------------------------------------
backup_config() {
    local version="$1"
    into="storage/backup/config/config-$(date +%Y%m%d%H%M%S)-${version//./}.yaml.bk"
    cp config.yaml "$into"
    echo "Backup config.yaml file into $into"
}