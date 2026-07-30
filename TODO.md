# Codedock Tunnel TODO

The Go backend, CLI, protocol package, and current TypeScript framework adapters are implemented. Remaining work is dashboard functionality, desktop integration, optional integrations, and release operations.

## Standalone product principles

- [ ] Provide Codedock integration as an optional adapter, plugin, or external API client.
- [ ] Document the hosted standalone product and optional Codedock integration.

## Product surfaces

- [ ] Build a standalone tunnel dashboard with its own authentication and organization model.
- [ ] Keep the Tauri desktop application usable with standalone tunnel servers.

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
