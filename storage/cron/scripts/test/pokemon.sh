#!/bin/bash

# Source the function
source /app/scripts/curl_debug.sh

in="$1"

case "$in" in
  test_a)
    curl_json_debug \
      --url "https://pokeapi.co/api/v2/pokemon?limit=3&offset=5" \
      --method GET
    ;;
  test_b)
    curl_json_debug \
      --url "https://pokeapi.co/api/v2/pokemon/pikachu" \
      --method GET
    ;;
  *)
    exit 1
    ;;
esac