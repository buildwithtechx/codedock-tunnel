# syntax=docker/dockerfile:1

# Build the control-plane API.
FROM golang:1.25-alpine AS build

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION}" -o /out/api ./cmd/server

# Run as a non-root user.
FROM alpine:3.22

RUN addgroup -S codedock && adduser -S -G codedock codedock

WORKDIR /app
COPY --from=build /out/api /app/api

USER codedock

EXPOSE 8080
STOPSIGNAL SIGTERM
ENTRYPOINT ["/app/api"]
