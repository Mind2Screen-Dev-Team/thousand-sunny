#!/bin/bash

# -----------------------------------------------------------------------------
# Function: curl_json_debug
# Description: Executes a curl HTTP request and outputs a detailed JSON log
#              including URL, method, query params, request/response headers,
#              request body, response body, status code, and latency.
#
# Arguments:
#   - --url <URL>:         Full request URL (with optional query params)
#   - --method <METHOD>:   HTTP method (GET, POST, PUT, DELETE, etc.)
#   - --data <BODY>:       Raw request body string (e.g., JSON)
#   - --header <HEADER>:   Request headers (repeatable, e.g., --header "Authorization: Bearer token")
# -----------------------------------------------------------------------------
# curl_debug.sh
curl_json_debug() {
  local url=""
  local method="GET"
  local body_data=""
  local headers=()

  # Parse arguments
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --url)
        url="$2"
        shift 2
        ;;
      --method)
        method="$2"
        shift 2
        ;;
      --data)
        body_data="$2"
        shift 2
        ;;
      --header)
        headers+=("$2")
        shift 2
        ;;
      *)
        echo "Unknown argument: $1"
        return 1
        ;;
    esac
  done

  # Unique temp file names
  local timestamp rand
  timestamp=$(($(date +%s%N)/1000000))
  rand=$RANDOM

  local res_body="/tmp/res_body_${timestamp}_${rand}.json"
  local res_headers="/tmp/res_headers_${timestamp}_${rand}.txt"
  local req_headers="/tmp/req_headers_${timestamp}_${rand}.txt"
  local timing="/tmp/timing_${timestamp}_${rand}.txt"
  local request_body="/tmp/request_body_${timestamp}_${rand}.json"

  # Prepare request headers file
  : > "$req_headers"
  for h in "${headers[@]}"; do
    echo "$h" >> "$req_headers"
  done

  # Save request body to temp file (if any)
  echo "$body_data" > "$request_body"

  # Extract query params (basic)
  local query_params
  if [[ "$url" == *"?"* ]]; then
    query_params=$(echo "$url" | awk -F '?' '{print $2}' | tr '&' '\n' | jq -Rn '
      [inputs | select(length > 0) | split("=") | {(.[0]): .[1]}] | add')
  else
    query_params="null"
  fi

  # Build curl args
  local curl_args=(-s -X "$method" -w "%{time_total}\n%{http_code}" -D "$res_headers" -o "$res_body")
  for h in "${headers[@]}"; do
    curl_args+=(-H "$h")
  done
  if [[ -n "$body_data" ]]; then
    curl_args+=(-d "$body_data")
  fi
  curl_args+=("$url")

  # Run curl
  curl "${curl_args[@]}" > "$timing"

  # Parse curl results
  local time_total http_code latency_ms
  time_total=$(head -n 1 "$timing")
  http_code=$(tail -n 1 "$timing")
  latency_ms=$(awk "BEGIN {printf \"%.0f\", $time_total * 1000}")

  # Format response headers (remove status line)
  local res_headers_clean formatted_res_headers formatted_req_headers
  res_headers_clean=$(sed '1d' "$res_headers" | grep ':' || echo '{}')
  formatted_res_headers=$(echo "$res_headers_clean" | jq -Rn '
    [inputs | select(length > 0) | split(": ") | {(.[0]): .[1]}] | add')
  formatted_req_headers=$(jq -Rn --rawfile headers "$req_headers" '
    [$headers | split("\n")[] | select(length > 0) | split(": ") | {(.[0]): .[1]}] | add')

  # Compose final JSON output
  jq -c -n \
    --arg url "$url" \
    --arg method "$method" \
    --argjson query_params "$query_params" \
    --argjson request_headers "$formatted_req_headers" \
    --slurpfile request_body "$request_body" \
    --arg status "$http_code" \
    --argjson response_headers "$formatted_res_headers" \
    --slurpfile response_body "$res_body" \
    --argjson latency_ms "$latency_ms" \
    '{
      url: $url,
      method: $method,
      query_params: $query_params,
      request_headers: $request_headers,
      request_body: ($request_body[0] // null),
      status: ($status | tonumber),
      response_headers: $response_headers,
      latency_ms: $latency_ms,
      response_body: ($response_body[0] // null)
    }'

  # Cleanup temp files
  rm -f "$res_body" "$res_headers" "$req_headers" "$timing" "$request_body"
}