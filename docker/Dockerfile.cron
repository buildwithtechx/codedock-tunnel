# syntax=docker/dockerfile:1

# Build scheduled maintenance jobs.
FROM golang:1.25-alpine AS build

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION}" -o /out/cron ./cmd/cron

# Run this image from a scheduler, not as a public service.
FROM alpine:3.22

RUN addgroup -S codedock && adduser -S -G codedock codedock

WORKDIR /app
COPY --from=build /out/cron /app/cron

USER codedock

ENTRYPOINT ["/app/cron"]
