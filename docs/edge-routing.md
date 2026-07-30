# Tunnel edge routing

The tunnel relay owns public routing. DNS automation is optional and is not part of the tunnel protocol.

## Hosted relay

Create these records once for the hosted installation:

```text
tunnel.codedock-tunnel.dev       A/AAAA <relay-address>
*.tunnel.codedock-tunnel.dev     A/AAAA <relay-address>
```

The relay then assigns URLs such as:

```text
https://abc123.tunnel.codedock-tunnel.dev
```

No DNS record is created for an individual tunnel.

When the relay is directly exposed on the public network, provision a certificate containing `*.tunnel.codedock-tunnel.dev` and `tunnel.codedock-tunnel.dev`. Mount the certificate and key into the relay and set:

```text
CODEDOCK_REQUIRE_TLS=true
CODEDOCK_TLS_CERT_FILE=/run/secrets/tunnel/fullchain.pem
CODEDOCK_TLS_KEY_FILE=/run/secrets/tunnel/privkey.pem
CODEDOCK_TUNNEL_DOMAIN=tunnel.codedock-tunnel.dev
```

The relay reloads the certificate on the next TLS handshake after either mounted file changes. A deployment can renew the certificate atomically by writing new files beside the existing files and replacing them with a rename.

Wildcard certificates require DNS-01 validation. The current `autocert` helper is intended for individual host certificates and should not be used as the wildcard issuance mechanism.

## Codedock PaaS or another reverse proxy

When Codedock PaaS or another reverse proxy owns the public domain, TLS belongs to that proxy. Configure the relay as an internal HTTP service and leave these values disabled:

```text
CODEDOCK_REQUIRE_TLS=false
CODEDOCK_TLS_CERT_FILE=
CODEDOCK_TLS_KEY_FILE=
```

Configure the proxy with the relay domain and forward WebSocket upgrades. The external client still uses `wss://`; TLS is terminated at the proxy and the internal hop uses HTTP/WebSocket. The relay container must not publish its internal port directly to the internet.

HTTP tunnels work through the normal HTTPS reverse-proxy route. Raw TCP and UDP tunnels require explicit public TCP/UDP port mappings or a proxy that supports TCP/UDP entrypoints; an ordinary HTTPS domain route cannot carry raw TCP or UDP traffic.

## Custom domains

Custom domains are verified independently from the relay wildcard. Users manage their DNS records themselves and provide either the DNS TXT challenge or HTTP challenge required by the dashboard.

Provider-specific DNS automation is intentionally outside the core service. A future user-authorized integration may create records through a customer-owned DNS token, but the hosted service does not store or use a global DNS provider credential.

## Local development

Use the local example values:

```text
CODEDOCK_REQUIRE_TLS=false
CODEDOCK_TUNNEL_REQUIRE_TLS=false
CODEDOCK_TUNNEL_DOMAIN=tunnel.localhost
```
