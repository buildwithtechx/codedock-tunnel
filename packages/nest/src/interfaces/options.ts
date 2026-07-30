import type { RelayConnectionOptions } from '@codedock/sdk';

export const CODEDOCK_TUNNEL_OPTIONS = Symbol('CODEDOCK_TUNNEL_OPTIONS');

export type NestTunnelOptions = RelayConnectionOptions & {
  localPort: number;
  subdomain?: string;
  password?: string;
  autoStart?: boolean;
};
