#!/bin/bash

URL="http://localhost:8080/openapi.json"
DIR="docs"
OUTPUT="./$DIR/openapi.json"

mkdir -p "$DIR"

curl -s -o "$OUTPUT" "$URL"

if [ $? -eq 0 ]; then
  echo "Downloaded successfully: $OUTPUT"
else
  echo "Failed to download from $URL" >&2
  exit 1
fi