#!/bin/bash

# Source the function
source /app/scripts/curl_debug.sh

in="$1"

case "$in" in
  fetch_pikachu)
    curl_json_debug \
      --url "https://pokeapi.co/api/v2/pokemon/pikachu" \
      --method GET \
      --header "Accept: application/json" \
      --header "Content-Type: application/json" \
      >> /var/log/pokemon_pikachu_at_$(date '+%Y%m%d').log 2>&1
    ;;
  *)
    exit 0
    ;;
esac