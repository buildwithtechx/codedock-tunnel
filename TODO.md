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
- [ ] Publish one reusable `sdk-ts` package for Node.js, TypeScript, browser, and framework integrations.
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

- [ ] Create `apps/web/` for the standalone React dashboard.
- [ ] Create `apps/desktop/` for the Tauri desktop application and `src-tauri/` Rust shell.
- [x] Create `cmd/server/` for the control-plane API server.
- [x] Create `cmd/tunnel/` for the public tunnel relay and data plane.
- [ ] Create `cmd/cli/` for the standalone CLI.
- [ ] Create `cmd/cron/` for cleanup, expiry, analytics rollups, certificate jobs, and scheduled maintenance.
- [ ] Create `cmd/check/` for custom-domain verification and edge readiness checks.
- [x] Organize domain models and DTOs under `internal/models/`.
- [x] Organize PostgreSQL persistence and repository interfaces under `internal/repositories/`.
- [x] Organize business logic and provider-neutral use cases under `internal/services/`.
- [x] Organize Fiber HTTP controllers under `internal/handlers/`.
- [x] Organize Fiber server setup, routes, middleware, and readiness wiring under `internal/http/`.
- [x] Organize relay, HTTP, TCP, UDP, and WebSocket engine primitives under `internal/engine/`.
- [x] Organize authentication and authorization helpers under `internal/auth/`.
- [x] Organize environment and runtime configuration under `internal/config/`.
- [x] Create the PostgreSQL infrastructure package under `internal/infra/postgres/`.
- [x] Organize Redis, billing providers, telemetry, storage, and distributed locks under `internal/infra/`.
- [x] Organize retryable cleanup, usage, billing, and reconciliation jobs under `internal/workers/`.
- [x] Keep repositories free of business policy and keep services independent from Fiber handlers.
- [x] Keep infrastructure adapters behind narrow interfaces consumed by services and workers.
- [x] Create `pkg/client/` for reusable Go client functionality.
- [x] Create `pkg/protocol/` for public Go protocol types and helpers.
- [x] Create `pkg/version/` for shared build version metadata.
- [x] Create `protocol/schema/` for language-neutral protocol definitions.
- [x] Generate protocol bindings under `protocol/generated/go/` and `protocol/generated/typescript/`.
- [x] Create `packages/protocol-ts/` for generated and hand-written TypeScript protocol types.
- [ ] Create `packages/sdk-ts/` as the framework-neutral Node.js and browser client.
- [ ] Create `packages/react/` as a thin React hooks and provider layer over `sdk-ts`.
- [ ] Create `packages/vite-plugin/` for local development tunnel integration with Vite.
- [ ] Create `packages/next/` for Next.js server, route, and development integration.
- [ ] Create `packages/nest/` for NestJS modules, providers, and tunnel lifecycle integration.
- [ ] Create `packages/express/` for Express middleware and tunnel lifecycle integration.
- [ ] Create `packages/tauri-bridge/` for desktop commands and Go CLI tunnel-client communication.
- [ ] Keep framework integrations as thin adapters over the same protocol and SDK lifecycle.
- [ ] Use peer dependencies for optional frameworks so the base SDK stays lightweight.
- [ ] Add package export maps, type declarations, and tree-shakable entrypoints.
- [ ] Keep framework adapters independently versioned while preserving compatibility with the core SDK.
- [ ] Keep `integrations/codedock/` as an optional external adapter.
- [ ] Keep database migrations in `migrations/` and integration tests in `tests/`.
- [x] Add independent Dockerfiles under `docker/` for API, tunnel, cron, and check.
- [x] Distribute the CLI as platform binaries instead of deploying a CLI service.
- [ ] Keep the root workspace configuration independent from the Codedock repository.
- [ ] Ensure the tunnel core has no imports of Codedock models, routes, authentication, or database packages.

## Product foundation

- [ ] Define the initial product boundary between `codedock-tunnel` and Codedock.
- [ ] Choose the first supported tunnel types: HTTP, HTTPS, and raw TCP.
- [ ] Define tunnel lifecycle states: creating, connecting, active, disconnected, expired, and revoked.
- [ ] Define the public URL and subdomain allocation strategy.
- [ ] Document the self-hosted and managed deployment models.
- [ ] Choose the project license and contribution policy.

## Go technology baseline

- [ ] Use `github.com/gofiber/fiber/v2` for the HTTP API and server middleware.
- [ ] Use `github.com/gofiber/contrib/websocket` for CLI and dashboard WebSocket connections.
- [x] Install `gorm.io/gorm` with `gorm.io/driver/postgres` for PostgreSQL persistence.
- [x] Define the initial PostgreSQL control-plane models and GORM migration registry.
- [ ] Use PostgreSQL as the primary control-plane database for hosted and self-hosted installations.
- [ ] Use PostgreSQL transactions, constraints, indexes, and explicit versioned migrations for durable state.
- [x] Use `github.com/redis/go-redis/v9` for ephemeral sessions, heartbeats, rate limits, presence, and relay coordination.
- [ ] Keep tunnel payloads out of Redis and PostgreSQL unless explicitly required for analytics or audit purposes.
- [x] Use the standard library `log/slog` with structured telemetry reporting.
- [ ] Use `github.com/go-playground/validator/v10` or equivalent request validation at API boundaries.
- [ ] Use `github.com/google/uuid` or a documented ID strategy consistently across database and protocol entities.
- [ ] Use `github.com/caarlos0/env/v11` for tagged environment configuration.
- [ ] Define dependency versions in `go.mod` and review them through automated vulnerability scanning.
- [ ] Document why each dependency is required and avoid duplicate HTTP, ORM, crypto, or WebSocket stacks.

## Identity, organizations, and access

- [x] Create standalone user, session, organization, organization-member, and role models.
- [ ] Implement Google and GitHub OAuth sign-in with revocable sessions.
- [x] Add CLI login sessions with one-time device codes and explicit expiration.
- [ ] Add organization membership roles and resource-level authorization.
- [ ] Keep Google and GitHub OAuth providers independent from Codedock.
- [x] Store session state in PostgreSQL with revocation support.
- [ ] Add account deletion, organization ownership transfer, and member removal flows.
- [ ] Audit authentication, membership, credential, billing, and tunnel-management actions.

## Plans, subscriptions, and billing

- [x] Create plan and subscription persistence models with provider-neutral billing fields.
- [x] Create plan definitions with tunnel, domain, member, bandwidth, retention, and connection limits.
- [ ] Store subscription status, provider, provider customer ID, provider subscription ID, product ID, billing interval, period end, cancellation state, trial state, and timestamps.
- [ ] Create a provider-neutral billing interface for checkout, portal access, cancellation, resumption, and subscription lookup.
- [x] Add a Polar billing adapter for international cards and subscription management.
- [x] Add a Paystack billing adapter for local currency and recurring payments.
- [ ] Keep billing-provider credentials and authorization codes encrypted at rest.
- [ ] Verify every provider webhook signature before processing events.
- [ ] Make webhook processing idempotent using provider event IDs and database constraints.
- [ ] Handle subscription created, active, updated, canceled, paused, past-due, expired, and revoked events.
- [ ] Implement plan entitlement checks in the API and relay before allocating tunnels or connections.
- [ ] Prevent client-provided plan or limit values from overriding server-side entitlements.
- [ ] Add grace periods, downgrade behavior, failed-payment handling, and free-plan fallback rules.
- [ ] Add billing history, invoices, receipts, and customer self-service portal links.
- [ ] Keep billing optional for self-hosted deployments while preserving local plan configuration.
- [ ] Use Zepto Mail as the only transactional email provider.

## Usage, analytics, and operations data

- [ ] Track tunnel count, active connections, bandwidth, request count, error rate, and retention usage per organization.
- [ ] Keep ephemeral presence, heartbeats, rate limits, and active tunnel indexes in Redis.
- [ ] Use Redis TTLs and periodic reconciliation to remove stale tunnel and organization presence entries.
- [x] Define durable control-plane and billing models for PostgreSQL through GORM.
- [ ] Store durable control-plane records and billing state in PostgreSQL through GORM.
- [ ] Support optional TimescaleDB or PostgreSQL time-series tables for high-volume usage analytics.
- [ ] Keep raw payloads and secrets out of analytics storage by default.
- [ ] Add retention policies based on the organization plan.
- [ ] Add usage snapshots and aggregation jobs in `cmd/cron`.
- [ ] Add dashboards for current usage, historical usage, plan limits, and billing status.

## Billing and operations jobs

- [x] Add reusable stale-session, API-key, and device-login cleanup jobs.
- [x] Add reusable usage aggregation and billing reconciliation worker contracts.
- [ ] Run stale-session cleanup, Redis index reconciliation, usage aggregation, and subscription maintenance in `cmd/cron`.
- [ ] Run recurring billing checks and provider reconciliation as retryable idempotent jobs.
- [x] Use Redis distributed locks so only one cron worker processes a billing or reconciliation batch.
- [ ] Add exponential backoff, dead-letter handling, and operator-visible job status.
- [ ] Run subscription maintenance immediately on startup only when protected by idempotency and locking.
- [ ] Add alerts for failed webhooks, repeated payment failures, stale presence growth, and quota inconsistencies.

## Protocol

- [x] Define a versioned control protocol for server, CLI, dashboard, and desktop clients.
- [x] Define tunnel open, close, error, heartbeat, and data messages.
- [ ] Define authentication, authorization, and capability messages.
- [ ] Define protocol version negotiation and backward compatibility rules.
- [ ] Define maximum frame size, idle timeout, connection timeout, and backpressure behavior.
- [ ] Add protocol conformance fixtures shared by Go and TypeScript clients.

## Go tunnel relay

- [ ] Create the `codedock-tunnel-server` binary.
- [ ] Implement CLI tunnel registration and authenticated session management.
- [ ] Implement persistent connections, heartbeats, reconnects, and graceful shutdown.
- [x] Add a concurrency-safe session registry with takeover protection.
- [x] Add a request broker with request IDs, timeouts, cancellation, and response routing.
- [x] Add reusable HTTP, TCP, and UDP protocol payloads.
- [x] Implement an HTTP edge adapter with bounded bodies and hop-by-hop header filtering.
- [x] Implement hostname-to-tunnel session resolution with base-domain validation.
- [x] Add bounded stream copying with context cancellation.
- [x] Add O(1)-style port allocation and release for TCP/UDP listeners.
- [x] Add bandwidth accounting and limit enforcement primitives.
- [x] Add TCP connection and UDP packet state registries with cleanup operations.
- [ ] Implement multiplexed stream handling.
- [ ] Implement HTTP and TCP forwarding.
- [ ] Implement tunnel routing by tunnel ID and hostname.
- [ ] Implement connection and stream cleanup after disconnects.
- [ ] Add bounded memory buffers and backpressure.
- [ ] Add configurable connection, bandwidth, and tunnel limits.
- [ ] Add structured logs, metrics, and health endpoints.

## Background jobs and verification

- [ ] Define jobs that can run inside the server for small self-hosted installations.
- [ ] Define jobs that can run in `cmd/cron` for independent horizontal scaling.
- [ ] Expire inactive tunnels, sessions, credentials, reservations, and public hostnames.
- [ ] Roll up traffic and connection analytics without storing request secrets or bodies by default.
- [ ] Verify custom-domain ownership using DNS and HTTP challenges.
- [ ] Verify certificate readiness and edge routing before activating custom domains.
- [ ] Prevent verification jobs from being used to probe arbitrary private networks.
- [ ] Make all jobs idempotent and safe to retry.

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
- [ ] Separate user identity, client identity, and tunnel capability tokens.
- [ ] Bind every stream to an authorized tunnel and target.
- [ ] Validate hostnames and reject ambiguous or malformed routing input.
- [ ] Add origin and host validation for HTTP tunnels.
- [ ] Add rate limiting for tunnel creation, authentication, and connection attempts.
- [ ] Add abuse controls for open proxies, port scanning, and excessive bandwidth.
- [ ] Redact credentials, tokens, cookies, and request bodies from logs.
- [ ] Threat-model the relay server, CLI client, public edge, and Codedock integration.
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

## Storage and operations

- [x] Add a safe atomic local object storage adapter for tunnel artifacts and exports.
- [ ] Store secret material encrypted at rest.
- [ ] Define cleanup for expired tunnels, sessions, and unused subdomains.
- [ ] Add database migrations and migration integrity checks.
- [ ] Add backups and restore procedures for control-plane metadata.
- [ ] Add Docker deployment examples.
- [ ] Add horizontal scaling design for multiple relay servers.
- [ ] Define CLI-to-relay affinity and session handoff behavior.
- [ ] Add Prometheus-compatible metrics and operational dashboards.
- [ ] Add resource budgets to preserve the lightweight runtime target.

## Testing and release

- [ ] Add unit tests for protocol encoding, authentication, routing, and limits.
- [ ] Add integration tests using real local HTTP and TCP targets.
- [ ] Add reconnect, timeout, backpressure, and graceful-shutdown tests.
- [ ] Add cross-platform CLI tunnel tests.
- [ ] Add fuzz tests for protocol frames and hostname routing.
- [ ] Add dependency and container vulnerability scanning.
- [ ] Add reproducible release builds and signed artifacts.
- [ ] Add compatibility tests before changing the protocol version.
- [ ] Document local development, staging, production, and incident-response workflows.
