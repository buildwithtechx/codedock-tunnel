#!/usr/bin/env bash

set -euo pipefail

tag="${TAG:-dev}"

docker build -f docker/Dockerfile.api -t "codedock-api:$tag" .
docker build -f docker/Dockerfile.tunnel -t "codedock-tunnel-server:$tag" .
docker build -f docker/Dockerfile.cron -t "codedock-tunnel-cron:$tag" .
docker build -f docker/Dockerfile.check -t "codedock-tunnel-check:$tag" .
