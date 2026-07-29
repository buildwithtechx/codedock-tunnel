FROM node:22-alpine AS dashboard-build

WORKDIR /workspace
COPY package.json package-lock.json* ./
COPY apps/dashboard/package.json apps/dashboard/package.json
COPY packages packages
RUN npm ci
COPY apps/dashboard apps/dashboard
COPY tsconfig.base.json biome.json ./
RUN npm run build:dashboard

FROM golang:1.25-alpine AS go-build

WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags='-s -w' -o /out/tunnel-server ./cmd/tunnel-server

FROM alpine:3.22

RUN addgroup -S codedock && adduser -S -G codedock codedock
WORKDIR /app
COPY --from=go-build /out/tunnel-server /app/tunnel-server
COPY --from=dashboard-build /workspace/apps/dashboard/dist /app/dashboard
RUN mkdir -p /app/data && chown -R codedock:codedock /app

USER codedock
ENV CODEDOCK_TUNNEL_DASHBOARD_DIR=/app/dashboard
EXPOSE 8080
ENTRYPOINT ["/app/tunnel-server"]
