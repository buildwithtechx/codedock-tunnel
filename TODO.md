# Codedock Tunnel TODO

## Standalone product principles

- [ ] Make `codedock-tunnel` usable without Codedock, a Codedock account, or a Codedock server.
- [ ] Treat Codedock as one optional integration among many, not as the tunnel system's control plane.
- [ ] Support an independent hosted service with its own accounts, organizations, billing, API, dashboard, and CLI.
- [ ] Support a fully self-hosted deployment with a user-owned relay, database, domain, and TLS certificates.
- [ ] Keep tunnel identity, credentials, sessions, routing, quotas, analytics, and audit history owned by `codedock-tunnel`.
- [ ] Define an integration-neutral public API for dashboards, CLIs, SDKs, and external platforms.
- [ ] Ensure the core tunnel server never imports or requires Codedock packages, models, routes, or authentication tokens.
- [ ] Provide Codedock integration as an optional adapter, plugin, or external API client.
- [ ] Allow users to run the tunnel CLI directly against any compatible tunnel server URL.
- [ ] Document the difference between the standalone product, self-hosted mode, and optional Codedock integration.

## Product surfaces

- [ ] Build a standalone tunnel dashboard with its own authentication and organization model.
- [ ] Build a standalone CLI that can create and manage tunnels without Codedock.
- [ ] Publish a reusable client SDK for Node.js, TypeScript, and Go integrations.
- [ ] Keep the Tauri desktop application usable with standalone tunnel servers.
- [ ] Support configuration of a custom server URL for local, self-hosted, and enterprise deployments.
- [ ] Define stable public API and protocol versioning independent of any Codedock release.

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
│   ├── tunnel-server/
│   ├── tunnel-agent/
│   ├── codedock-tunnel/
│   ├── tunnel-cron/
│   └── tunnel-internal-check/
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── routes/
│   │   └── server.go
│   ├── auth/
│   ├── control/
│   ├── engine/
│   │   ├── http_proxy/
│   │   ├── tcp_proxy/
│   │   ├── udp_proxy/
│   │   └── websocket_proxy/
│   ├── limits/
│   ├── relay/
│   ├── routing/
│   ├── sessions/
│   ├── storage/
│   ├── telemetry/
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
│   ├── sdk-ts/
│   ├── cli-ui/
│   └── codedock-adapter/
├── integrations/
│   └── codedock/
├── migrations/
├── tests/
│   ├── integration/
│   ├── protocol/
│   ├── security/
│   └── e2e/
├── deploy/
│   ├── docker/
│   ├── systemd/
│   └── kubernetes/
├── docs/
├── scripts/
├── go.mod
├── package.json
├── tsconfig.base.json
├── biome.json
├── Makefile
├── Dockerfile
├── README.md
└── TODO.md
```

- [ ] Create `apps/dashboard/` for the standalone React dashboard.
- [ ] Create `apps/desktop/` for the Tauri desktop application and `src-tauri/` Rust shell.
- [ ] Create `cmd/tunnel-server/` for the public relay and control-plane server.
- [ ] Create `cmd/tunnel-agent/` for the outbound local-network agent.
- [ ] Create `cmd/codedock-tunnel/` for the standalone CLI.
- [ ] Create `cmd/tunnel-cron/` for cleanup, expiry, analytics rollups, certificate jobs, and scheduled maintenance.
- [ ] Create `cmd/tunnel-internal-check/` for custom-domain verification and edge readiness checks.
- [ ] Organize server code under `internal/api/`, `internal/auth/`, `internal/control/`, `internal/engine/`, `internal/relay/`, `internal/routing/`, `internal/sessions/`, `internal/storage/`, and `internal/telemetry/`.
- [ ] Create `pkg/client/` for reusable Go client functionality.
- [ ] Create `pkg/protocol/` for public Go protocol types and helpers.
- [ ] Create `protocol/schema/` for language-neutral protocol definitions.
- [ ] Generate protocol bindings under `protocol/generated/go/` and `protocol/generated/typescript/`.
- [ ] Create `packages/sdk-ts/` for reusable TypeScript and Node.js clients.
- [ ] Keep `integrations/codedock/` as an optional external adapter.
- [ ] Keep database migrations in `migrations/` and integration tests in `tests/`.
- [ ] Add `deploy/` for Docker, systemd, and self-hosted deployment configurations.
- [ ] Keep the root workspace configuration independent from the Codedock repository.
- [ ] Ensure the tunnel core has no imports of Codedock models, routes, authentication, or database packages.

## Product foundation

- [ ] Define the initial product boundary between `codedock-tunnel` and Codedock.
- [ ] Choose the first supported tunnel types: HTTP, HTTPS, and raw TCP.
- [ ] Define tunnel lifecycle states: creating, connecting, active, disconnected, expired, and revoked.
- [ ] Define the public URL and subdomain allocation strategy.
- [ ] Document the self-hosted and managed deployment models.
- [ ] Choose the project license and contribution policy.

## Protocol

- [ ] Define a versioned control protocol for server, agent, CLI, dashboard, and desktop clients.
- [ ] Define authentication, authorization, and capability messages.
- [ ] Define tunnel open, accept, close, error, heartbeat, and reconnect messages.
- [ ] Define protocol version negotiation and backward compatibility rules.
- [ ] Define maximum frame size, idle timeout, connection timeout, and backpressure behavior.
- [ ] Add protocol conformance fixtures shared by Go and TypeScript clients.

## Go tunnel server

- [ ] Create the `tunnel-server` binary.
- [ ] Implement agent registration and authenticated session management.
- [ ] Implement persistent connections, heartbeats, reconnects, and graceful shutdown.
- [ ] Implement multiplexed stream handling.
- [ ] Implement HTTP and TCP forwarding.
- [ ] Implement tunnel routing by tunnel ID and hostname.
- [ ] Implement connection and stream cleanup after disconnects.
- [ ] Add bounded memory buffers and backpressure.
- [ ] Add configurable connection, bandwidth, and tunnel limits.
- [ ] Add structured logs, metrics, and health endpoints.

## Background jobs and verification

- [ ] Define jobs that can run inside the server for small self-hosted installations.
- [ ] Define jobs that can run in `cmd/tunnel-cron` for independent horizontal scaling.
- [ ] Expire inactive tunnels, sessions, credentials, reservations, and public hostnames.
- [ ] Roll up traffic and connection analytics without storing request secrets or bodies by default.
- [ ] Verify custom-domain ownership using DNS and HTTP challenges.
- [ ] Verify certificate readiness and edge routing before activating custom domains.
- [ ] Prevent internal-check jobs from being used to probe arbitrary private networks.
- [ ] Make all jobs idempotent and safe to retry.

## Go tunnel agent

- [ ] Create the `tunnel-agent` binary.
- [ ] Implement local target validation and loopback-only defaults.
- [ ] Implement outbound TLS connection to the tunnel server.
- [ ] Implement reconnect with bounded exponential backoff and jitter.
- [ ] Implement graceful shutdown and connection draining.
- [ ] Implement local HTTP, HTTPS, and TCP forwarding.
- [ ] Prevent access to link-local, metadata, and unintended private network targets by default.
- [ ] Add optional explicit target allowlists.
- [ ] Add service installation instructions for Linux, macOS, and Windows.

## CLI

- [ ] Create the `codedock-tunnel` CLI.
- [ ] Implement `login` and credential storage.
- [ ] Implement `create`, `list`, `inspect`, `start`, `stop`, and `revoke` commands.
- [ ] Implement `codedock-tunnel http 3000` for local development.
- [ ] Implement `codedock-tunnel tcp 5432` for explicitly authorized TCP access.
- [ ] Add human-readable and JSON output modes.
- [ ] Add shell completion for supported shells.
- [ ] Never place long-lived secrets in command-line arguments or public URLs.

## Security

- [ ] Use TLS for all control and data-plane connections.
- [ ] Use short-lived, scoped tunnel credentials.
- [ ] Support immediate credential and tunnel revocation.
- [ ] Separate user identity, agent identity, and tunnel capability tokens.
- [ ] Bind every stream to an authorized tunnel and target.
- [ ] Validate hostnames and reject ambiguous or malformed routing input.
- [ ] Add origin and host validation for HTTP tunnels.
- [ ] Add rate limiting for tunnel creation, authentication, and connection attempts.
- [ ] Add abuse controls for open proxies, port scanning, and excessive bandwidth.
- [ ] Redact credentials, tokens, cookies, and request bodies from logs.
- [ ] Threat-model the relay server, agent, public edge, and Codedock integration.
- [ ] Add security tests for token replay, cross-tunnel access, stale sessions, and SSRF.

## Public edge and routing

- [ ] Define the TLS certificate strategy for wildcard tunnel domains.
- [ ] Implement hostname-to-tunnel routing.
- [ ] Support custom domains only after ownership verification.
- [ ] Define HTTP header forwarding and stripping rules.
- [ ] Preserve client IP information safely without trusting user-supplied headers.
- [ ] Add connection draining during tunnel shutdown and certificate rotation.
- [ ] Document reverse-proxy and firewall requirements.

## Dashboard

- [ ] Create the tunnel dashboard workspace.
- [ ] Implement tunnel creation and target configuration.
- [ ] Display connection status, public URL, agent version, and last heartbeat.
- [ ] Add start, stop, rotate credentials, and revoke actions.
- [ ] Display active connections and bandwidth usage.
- [ ] Add tunnel expiration and access policy controls.
- [ ] Add audit history for tunnel and credential actions.
- [ ] Add accessible loading, empty, error, and disconnected states.

## Tauri desktop app

- [ ] Create the Tauri 2 application shell.
- [ ] Define how the desktop app starts and supervises the Go agent.
- [ ] Keep tunnel data-plane logic in Go rather than duplicating it in Rust or TypeScript.
- [ ] Store credentials using the native operating-system secret store.
- [ ] Add tray controls for tunnel status and quick start/stop.
- [ ] Add native notifications for disconnects, expiry, and authentication failures.
- [ ] Package the agent and desktop application for macOS, Windows, and Linux.

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

## Storage and operations

- [ ] Choose persistent storage for tunnel metadata and credentials.
- [ ] Store secret material encrypted at rest.
- [ ] Define cleanup for expired tunnels, sessions, and unused subdomains.
- [ ] Add database migrations and migration integrity checks.
- [ ] Add backups and restore procedures for control-plane metadata.
- [ ] Add Docker and systemd deployment examples.
- [ ] Add horizontal scaling design for multiple relay servers.
- [ ] Define agent-to-relay affinity and session handoff behavior.
- [ ] Add Prometheus-compatible metrics and operational dashboards.
- [ ] Add resource budgets to preserve the lightweight runtime target.

## Testing and release

- [ ] Add unit tests for protocol encoding, authentication, routing, and limits.
- [ ] Add integration tests using real local HTTP and TCP targets.
- [ ] Add reconnect, timeout, backpressure, and graceful-shutdown tests.
- [ ] Add cross-platform agent tests.
- [ ] Add fuzz tests for protocol frames and hostname routing.
- [ ] Add dependency and container vulnerability scanning.
- [ ] Add reproducible release builds and signed artifacts.
- [ ] Add compatibility tests before changing the protocol version.
- [ ] Document local development, staging, production, and incident-response workflows.
