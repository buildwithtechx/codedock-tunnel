#!/usr/bin/env bash

set -euo pipefail

tag="${TAG:-dev}"

docker build -f deploy/docker/api.Dockerfile -t "codedock-api:$tag" .
docker build -f deploy/docker/tunnel.Dockerfile -t "codedock-tunnel-server:$tag" .
docker build -f deploy/docker/cron.Dockerfile -t "codedock-tunnel-cron:$tag" .
docker build -f deploy/docker/check.Dockerfile -t "codedock-tunnel-check:$tag" .
docker build -f deploy/docker/cli.Dockerfile -t "codedock-tunnel-cli:$tag" .
