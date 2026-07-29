FROM node:22-alpine AS web-build

WORKDIR /workspace
COPY package.json package-lock.json* ./
COPY apps/web/package.json apps/web/package.json
COPY packages packages
RUN npm ci
COPY apps/web apps/web
COPY tsconfig.base.json biome.json ./
RUN npm run build:web

FROM golang:1.25-alpine AS go-build

WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags='-s -w' -o /out/tunnel-server ./cmd/server

FROM alpine:3.22

RUN addgroup -S codedock && adduser -S -G codedock codedock
WORKDIR /app
COPY --from=go-build /out/tunnel-server /app/tunnel-server
COPY --from=web-build /workspace/apps/web/dist /app/web
RUN mkdir -p /app/data && chown -R codedock:codedock /app

USER codedock
ENV CODEDOCK_TUNNEL_WEB_DIR=/app/web
EXPOSE 8080
ENTRYPOINT ["/app/tunnel-server"]
