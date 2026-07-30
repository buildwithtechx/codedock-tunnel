import type { RelayConnectionOptions } from '@codedock-tunnel/sdk-ts';

export const CODEDOCK_TUNNEL_OPTIONS = Symbol('CODEDOCK_TUNNEL_OPTIONS');

export type NestTunnelOptions = RelayConnectionOptions & {
  localPort: number;
  subdomain?: string;
  autoStart?: boolean;
};
