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

Provision a certificate containing `*.tunnel.codedock-tunnel.dev` and `tunnel.codedock-tunnel.dev`. Mount the certificate and key into the relay and set:

```text
CODEDOCK_REQUIRE_TLS=true
CODEDOCK_TLS_CERT_FILE=/run/secrets/tunnel/fullchain.pem
CODEDOCK_TLS_KEY_FILE=/run/secrets/tunnel/privkey.pem
CODEDOCK_TUNNEL_DOMAIN=tunnel.codedock-tunnel.dev
```

The relay reloads the certificate on the next TLS handshake after either mounted file changes. A deployment can renew the certificate atomically by writing new files beside the existing files and replacing them with a rename.

Wildcard certificates require DNS-01 validation. The current `autocert` helper is intended for individual host certificates and should not be used as the wildcard issuance mechanism.

## Custom domains

Custom domains are verified independently from the relay wildcard. Users may manage DNS themselves, or the API may use a configured provider adapter to create verification TXT records. Cloudflare is the first optional adapter; it is enabled only when `CODEDOCK_DNS_PROVIDER=cloudflare` and its credentials are present.

## Local development

Use the local example values:

```text
CODEDOCK_REQUIRE_TLS=false
CODEDOCK_TUNNEL_REQUIRE_TLS=false
CODEDOCK_TUNNEL_DOMAIN=tunnel.localhost
```
