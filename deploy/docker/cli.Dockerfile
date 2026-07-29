# syntax=docker/dockerfile:1

# Build the CLI for CI and automation environments.
FROM golang:1.25-alpine AS build

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION}" -o /out/codedock-tunnel ./cmd/cli

FROM alpine:3.22

RUN addgroup -S codedock && adduser -S -G codedock codedock

COPY --from=build /out/codedock-tunnel /usr/local/bin/codedock-tunnel

USER codedock
ENTRYPOINT ["/usr/local/bin/codedock-tunnel"]
