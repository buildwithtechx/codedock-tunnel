# Codedock Tunnel TODO

The Go backend, CLI, protocol package, and current TypeScript framework adapters are implemented. Remaining work is dashboard functionality, desktop integration, optional integrations, and release operations.

## Standalone product principles

- [ ] Provide Codedock integration as an optional adapter, plugin, or external API client.
- [ ] Document the hosted standalone product and optional Codedock integration.

## Product surfaces

- [ ] Build a standalone tunnel dashboard with its own authentication and organization model.
- [ ] Publish one reusable `sdk-ts` package for Node.js, TypeScript, browser, and framework integrations.
- [ ] Keep the Tauri desktop application usable with standalone tunnel servers.

## Repository structure

```text
codedock-tunnel/
├── apps/
│   ├── web/
│   │   ├── src/
│   │   │   ├── routes/
│   │   │   ├── components/
│   │   │   ├── features/
│   │   │   ├── hooks/
│   │   │   └── lib/
│   │   └── package.json
│   └── desktop/
│       ├── src/
│       ├── src-tauri/
│       └── package.json
├── cmd/
│   ├── server/
│   ├── tunnel/
│   ├── cli/
│   ├── cron/
│   └── check/
├── internal/
│   ├── models/
│   ├── repositories/
│   ├── services/
│   ├── handlers/
│   ├── http/
│   ├── engine/
│   ├── auth/
│   ├── config/
│   ├── infra/
│   │   ├── postgres/
│   │   ├── redis/
│   │   ├── billing/
│   │   ├── telemetry/
│   │   ├── storage/
│   │   └── locks/
│   └── workers/
├── pkg/
│   ├── client/
│   ├── protocol/
│   └── version/
├── protocol/
│   ├── schema/
│   └── generated/
│       ├── go/
│       └── typescript/
├── packages/
│   ├── protocol-ts/
│   │   ├── src/
│   │   └── package.json
│   ├── sdk-ts/
│   │   ├── src/
│   │   └── package.json
│   ├── react/
│   │   ├── src/
│   │   └── package.json
│   ├── vite-plugin/
│   │   ├── src/
│   │   └── package.json
│   ├── next/
│   │   ├── src/
│   │   └── package.json
│   ├── nest/
│   │   ├── src/
│   │   └── package.json
│   ├── express/
│   │   ├── src/
│   │   └── package.json
│   └── tauri-bridge/
│       ├── src/
│       └── package.json
├── integrations/
│   └── codedock/
├── migrations/
├── tests/
│   ├── integration/
│   ├── protocol/
│   ├── security/
│   └── e2e/
├── docker/
│   ├── Dockerfile.api
│   ├── Dockerfile.tunnel
│   ├── Dockerfile.cron
│   └── Dockerfile.check
├── docs/
├── scripts/
├── go.mod
├── package.json
├── tsconfig.base.json
├── biome.json
├── Makefile
├── README.md
└── TODO.md
```

- [x] Create `apps/web/` for the standalone React dashboard.
- [x] Create `apps/desktop/` for the Tauri desktop application and `src-tauri/` Rust shell.
- [x] Create `packages/protocol-ts/` as the TypeScript protocol contract and codec package.
- [x] Create `packages/sdk-ts/` as the framework-neutral Node.js and browser client.
- [x] Create `packages/react/` as a thin React hooks and provider layer over `sdk-ts`.
- [x] Create `packages/vite-plugin/` for local development tunnel integration with Vite.
- [x] Create `packages/next/` for Next.js server, route, and development integration.
- [x] Create `packages/nest/` for NestJS modules, providers, and tunnel lifecycle integration.
- [x] Create `packages/express/` for Express middleware and tunnel lifecycle integration.
- [ ] Create `packages/tauri-bridge/` for desktop commands and Go CLI tunnel-client communication.
- [x] Keep framework integrations as thin adapters over the same protocol and SDK lifecycle.
- [x] Use peer dependencies for optional frameworks so the base SDK stays lightweight.
- [x] Add package export maps, type declarations, and tree-shakable entrypoints.
- [x] Keep framework adapters independently versioned while preserving compatibility with the core SDK.
- [ ] Keep `integrations/codedock/` as an optional external adapter.

## Dashboard

- [ ] Create the tunnel dashboard workspace.
- [ ] Implement tunnel creation and target configuration.
- [ ] Display connection status, public URL, client version, and last heartbeat.
- [ ] Add start, stop, rotate credentials, and revoke actions.
- [ ] Display active connections and bandwidth usage.
- [ ] Add tunnel expiration and access policy controls.
- [ ] Add audit history for tunnel and credential actions.
- [ ] Add accessible loading, empty, error, and disconnected states.

## Tauri desktop app

- [ ] Create the Tauri 2 application shell.
- [ ] Define how the desktop app starts and supervises the Go CLI tunnel client.
- [ ] Keep tunnel data-plane logic in Go rather than duplicating it in Rust or TypeScript.
- [ ] Store credentials using the native operating-system secret store.
- [ ] Add tray controls for tunnel status and quick start/stop.
- [ ] Add native notifications for disconnects, expiry, and authentication failures.
- [ ] Package the CLI and desktop application for macOS, Windows, and Linux.

## Codedock integration

- [ ] Define an optional Codedock adapter that consumes the standalone tunnel API.
- [ ] Allow Codedock to create scoped tunnel credentials through the public integration API.
- [ ] Add optional tunnel metadata links to Codedock projects and services.
- [ ] Add tunnel status synchronization through polling or signed webhooks.
- [ ] Add optional tunnel creation after local-directory deployment.
- [ ] Add optional MCP tools for listing, creating, inspecting, and revoking tunnels.
- [ ] Keep Codedock usable when the tunnel service is unavailable.
- [ ] Ensure uninstalling or disabling Codedock integration never deletes standalone tunnel data.
- [ ] Add integration tests proving a standalone tunnel server works without Codedock.
