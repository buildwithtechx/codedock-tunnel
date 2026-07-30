# Codedock Tunnel TODO

The Go backend, CLI, protocol package, and current TypeScript framework adapters are implemented. Remaining work is dashboard functionality, desktop integration, optional integrations, and release operations.

## Standalone product principles

- [ ] Provide Codedock integration as an optional adapter, plugin, or external API client.
- [ ] Document the hosted standalone product and optional Codedock integration.

## Product surfaces

- [ ] Build a standalone tunnel dashboard with its own authentication and organization model.
- [ ] Keep the Tauri desktop application usable with standalone tunnel servers.

## Route structure

Dashboard routes should remain browser-only and organization-scoped:

```text
/
├── pricing
├── docs/*
├── login
├── signup
├── cli/login
├── admin/
│   ├── index
│   ├── users
│   ├── users/$userId
│   ├── organizations
│   ├── organizations/$organizationID
│   ├── tunnels
│   ├── subscriptions
│   ├── usage
│   ├── charts
│   ├── audit-logs
│   └── actions
└── $orgSlug/
    ├── index
    ├── tunnels/
    │   ├── index
    │   └── $tunnelId
    ├── agents
    ├── domains
    ├── requests
    ├── usage
    ├── billing
    ├── members
    ├── api-keys
    └── settings/
        ├── index
        ├── profile
        └── organization
```

The dashboard application should keep page routing separate from reusable
application code:

```text
apps/web/src/
├── routes/
│   ├── __root.tsx
│   ├── index.tsx
│   ├── login.tsx
│   ├── signup.tsx
│   └── $orgSlug/
├── components/
│   ├── ui/
│   ├── layout/
│   ├── navigation/
│   ├── feedback/
│   └── data-display/
├── features/
│   ├── auth/
│   ├── admin/
│   │   ├── users/
│   │   ├── organizations/
│   │   ├── tunnels/
│   │   ├── subscriptions/
│   │   ├── usage/
│   │   ├── charts/
│   │   └── actions/
│   ├── organizations/
│   ├── tunnels/
│   ├── agents/
│   ├── domains/
│   ├── usage/
│   ├── billing/
│   ├── audit-logs/
│   └── api-keys/
├── interfaces/
│   ├── api.ts
│   ├── auth.ts
│   ├── organization.ts
│   ├── tunnel.ts
│   └── billing.ts
├── hooks/
├── stores/
│   ├── auth-store.ts
│   ├── organization-store.ts
│   ├── tunnel-store.ts
│   └── ui-store.ts
├── lib/
│   ├── api-client.ts
│   ├── auth.ts
│   ├── query-client.ts
│   ├── route-guards.ts
│   └── utils.ts
├── integrations/
│   └── telemetry/
└── env.ts
```

- [ ] Keep route files focused on loaders, route metadata, and page composition.
- [ ] Protect `/admin/*` with a separate platform-admin authorization guard.
- [ ] Keep admin features and components isolated from organization-member features.
- [ ] Keep domain behavior inside `features/`, not inside route files.
- [ ] Keep shared visual primitives inside `components/ui/`.
- [ ] Keep API contracts and frontend DTOs inside `interfaces/`.
- [ ] Keep cross-feature clients, guards, query setup, and utilities inside `lib/`.
- [ ] Keep global client state in focused Zustand stores under `stores/`.
- [ ] Keep feature-specific hooks beside their feature unless shared by multiple domains.

The control-plane API should remain versioned under `/api/v1`:

```text
/api/v1/
├── auth/
│   ├── oauth/:provider
│   ├── oauth/:provider/callback
│   ├── device/start
│   ├── device/poll
│   ├── device/complete
│   ├── session
│   └── logout
├── account
├── organizations
│   ├── :organizationID
│   ├── :organizationID/members
│   ├── :organizationID/members/:memberID
│   └── :organizationID/transfer
├── organizations/:organizationID/
│   ├── tunnels
│   ├── agents
│   ├── domains
│   ├── usage/events
│   ├── usage/snapshot
│   ├── usage/requests
│   ├── billing/{status,checkout,portal,cancel,resume,invoices}
│   ├── api-keys
│   ├── audit-logs
│   └── settings
├── tunnels/:tunnelID/
│   ├── status
│   ├── revoke
│   ├── connections
│   └── requests
├── domains/:domainID/verify
├── agents/:agentID/
│   ├── heartbeat
│   └── revoke
├── admin/
│   ├── users
│   ├── users/:userID
│   ├── organizations
│   ├── organizations/:organizationID
│   ├── tunnels
│   ├── subscriptions
│   ├── usage
│   ├── charts
│   └── actions
└── webhooks/billing/:provider
```

- [x] Keep liveness, readiness, and metrics outside the versioned API at `/healthz`, `/readyz`, and `/metrics`.
- [x] Keep relay WebSocket transport separate at `wss://tunnel.codedock-tunnel.dev/v1/connect`.
- [ ] Add organization detail, member listing/removal, profile, API-key, audit-log, request-log, and invoice routes.
- [ ] Add platform-admin authorization and admin users, organizations, tunnels, subscriptions, usage, charts, and action routes.
- [ ] Move agent heartbeats away from browser-session authentication to agent or relay authentication.
- [ ] Keep internal health, usage ingestion, relay handoff, and agent authentication routes private.
- [ ] Route wildcard public tunnel traffic through `*.tunnel.codedock-tunnel.dev`, not through control-plane handlers.

- [ ] Keep `integrations/codedock/` as an optional external adapter.

## Future SDKs

- [ ] Create an official Laravel/PHP SDK over the public HTTP and relay APIs.
- [ ] Create an official Rust SDK for native services and desktop tooling.
- [ ] Create an official Go SDK for Go services and tunnel-aware workers.
- [ ] Create an official Angular adapter over the shared client contract.
- [ ] Keep every SDK aligned with the versioned protocol and authentication model.
- [ ] Publish language-specific SDKs only after compatibility, security, and conformance tests pass.

```text
packages/
├── php/
│   ├── src/
│   │   ├── Client/
│   │   ├── Contracts/
│   │   ├── Exceptions/
│   │   ├── Resources/
│   │   └── Laravel/
│   │       ├── Console/
│   │       ├── Facades/
│   │       ├── Http/
│   │       ├── Services/
│   │       └── CodedockServiceProvider.php
│   ├── config/
│   ├── tests/
│   ├── composer.json
│   └── README.md
├── rust/
│   ├── src/
│   │   ├── client.rs
│   │   ├── error.rs
│   │   ├── models.rs
│   │   └── protocol.rs
│   ├── tests/
│   ├── Cargo.toml
│   └── README.md
├── go/
│   ├── client/
│   ├── protocol/
│   ├── internal/
│   ├── tests/
│   ├── go.mod
│   └── README.md
└── angular/
    ├── src/
    │   ├── guards/
    │   ├── interceptors/
    │   ├── models/
    │   ├── services/
    │   └── tokens/
    ├── tests/
    ├── package.json
    └── README.md
```

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
