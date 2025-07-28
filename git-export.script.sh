#!/bin/bash
# Robust Git export to JSON (handles binary, diff, metadata safely)
# Usage:
#   bash export-commits-json.sh [--since=YYYY-MM-DD] [--until=YYYY-MM-DD] [--limit=N]
# Requires: jq

SINCE=""
UNTIL=""
LIMIT=""

for arg in "$@"; do
  case $arg in
    --since=*) SINCE="--since=${arg#*=}";;
    --until=*) UNTIL="--until=${arg#*=}";;
    --limit=*) LIMIT="-n ${arg#*=}";;
  esac
done

DELIM="@@@"

git log $SINCE $UNTIL $LIMIT --no-color --date=iso \
  --pretty=format:"%H$DELIM%an$DELIM%ad$DELIM%s" --patch |
awk -v RS="" -v FS="$DELIM" '
function clean(str) {
  gsub(/\x1B\[[0-9;]*[A-Za-z]/, "", str)  # remove ANSI
  gsub(/\\/,"\\\\",str)                   # escape backslashes
  gsub(/"/,"\\\"",str)                    # escape quotes
  return str
}
{
  # Split entire record into lines
  n = split($0, lines, "\n")

  # Initialize vars
  hash=""; author=""; date=""; message=""
  delete fileArr; fileCount=0
  delete diffArr; diffCount=0
  inBinary=0

  # Find the metadata line explicitly (line with 4 fields)
  for (i=1; i<=n; i++) {
    numFields = split(lines[i], meta, FS)
    if (numFields == 4 && hash == "") {
      hash = meta[1]
      author = meta[2]
      date = meta[3]
      message = meta[4]
      continue
    }

    # Process diff and files
    line = lines[i]
    if (index(line, "diff --git") == 1) {
      fileLine=line
      sub(/^.* b\//, "", fileLine)
      fileCount++
      fileArr[fileCount] = clean(fileLine)
    }
    if (index(line, "Binary files") == 1) {
      inBinary=1
    }
    if ((index(line, "@@") == 1 || index(line, "+") == 1 || index(line, "-") == 1) && !inBinary) {
      diffCount++
      diffArr[diffCount] = clean(line)
    }
  }

  # Build arrays safely
  fileJson=""
  for (j=1; j<=fileCount; j++) {
    fileJson = fileJson (fileJson ? "," : "") "\"" fileArr[j] "\""
  }

  diffJson=""
  for (k=1; k<=diffCount; k++) {
    diffJson = diffJson (diffJson ? "," : "") "\"" diffArr[k] "\""
  }

  # Output a JSON object per line
  printf "{"
  printf "\"commit\":\"%s\",", clean(hash)
  printf "\"author\":\"%s\",", clean(author)
  printf "\"date\":\"%s\",", clean(date)
  printf "\"message\":\"%s\",", clean(message)
  printf "\"files\":[%s],", fileJson
  printf "\"binary\":%s,", (inBinary ? "true" : "false")
  printf "\"diff\":[%s]", diffJson
  printf "}\n"
}
' | tr -d '\000-\037' | jq -s '.'
