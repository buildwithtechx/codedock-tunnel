# syntax=docker/dockerfile:1

# Build one Go command at a time so each runtime can be deployed independently.
FROM golang:1.25-alpine AS go-build

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGET=server
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION}" -o /out/codedock-tunnel ./cmd/${TARGET}

# Keep the runtime image small and non-root.
FROM alpine:3.22

RUN addgroup -S codedock && adduser -S -G codedock codedock

WORKDIR /app

COPY --from=go-build /out/codedock-tunnel /app/codedock-tunnel

RUN mkdir -p /app/data && chown -R codedock:codedock /app

USER codedock

STOPSIGNAL SIGTERM
ENTRYPOINT ["/app/codedock-tunnel"]
