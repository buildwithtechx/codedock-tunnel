#!/usr/bin/env bash

set -euo pipefail

version="${VERSION:-dev}"
output_dir="${OUTPUT_DIR:-dist}"

rm -rf "$output_dir"
mkdir -p "$output_dir"

for target in server tunnel cron check cli; do
  GOOS="${GOOS:-linux}" \
  GOARCH="${GOARCH:-amd64}" \
  go build \
    -ldflags="-s -w -X main.version=${version}" \
    -o "$output_dir/$target" \
    "./cmd/$target"
done
