# Codedock Tunnel

Codedock Tunnel is an independent tunneling platform for exposing local and private services through secure public endpoints.

It is designed to work without Codedock. Codedock is an optional integration that can create and manage tunnels through the public Codedock Tunnel API.

## Product surfaces

- Go tunnel server and public relay
- Go tunnel agent for machines behind NAT or firewalls
- Standalone `codedock-tunnel` CLI
- React dashboard for accounts, organizations, tunnels, access policies, and analytics
- Tauri desktop application
- Go and TypeScript client SDKs
- Self-hosted and managed deployment modes

## Repository layout

- `cmd/server` runs the relay and control-plane server.
- `cmd/agent` connects outward and forwards local traffic.
- `cmd/cli` provides the standalone CLI.
- `internal` contains private server, relay, routing, storage, and authentication code.
- `protocol` contains language-neutral protocol schemas and generated bindings.
- `apps/web` contains the standalone web interface.
- `apps/desktop` contains the Tauri desktop shell.
- `integrations/codedock` contains the optional Codedock adapter.

## Development

Requirements:

- Go 1.25 or newer
- Node.js 22 or newer
- npm 10 or newer

```sh
npm install
npm run dev
```

## Common commands

```sh
npm run build
npm run fmt
npm run typecheck
npm run test
make docker-build
```

## Deployment roles

The Go commands are independently deployable and do not all need to run on the same machine.

| Command | Deployment | Lifecycle |
| --- | --- | --- |
| `cmd/server` | Public VPS or container | Long-running relay and API |
| `cmd/agent` | User’s private machine or container | Long-running outbound connection |
| `cmd/cron` | Scheduled container, systemd timer, or Kubernetes CronJob | Short-lived maintenance jobs |
| `cmd/check` | Scheduled job or health-check container | Short-lived domain and edge checks |
| `cmd/cli` | User workstation binary or package manager | Invoked on demand |

Build the server image by default, or select another command with `TARGET`:

```sh
docker build --build-arg TARGET=server -t codedock-tunnel-server .
docker build --build-arg TARGET=agent -t codedock-tunnel-agent .
docker build --build-arg TARGET=cron -t codedock-tunnel-cron .
docker build --build-arg TARGET=check -t codedock-tunnel-check .
```

The CLI is normally distributed as a platform binary rather than deployed as a service.

## Standalone usage

The CLI must support a configurable tunnel server URL so users can connect to the hosted service, a private installation, or a local development server.

```sh
codedock-tunnel login --server https://tunnel.example.com
codedock-tunnel http 3000
```

## Design boundary

The tunnel core owns tunnel identity, credentials, sessions, routing, quotas, analytics, and audit history. It must not import Codedock models, routes, authentication, or database packages.

Codedock integrations communicate with the public tunnel API and may be disabled without affecting standalone tunnel operation.
