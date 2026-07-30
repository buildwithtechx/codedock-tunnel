import type { RelayConnectionOptions } from '@codedock/sdk';

export type CodedockTunnelPluginOptions = RelayConnectionOptions & {
  enabled?: boolean;
  localPort?: number;
  subdomain?: string;
  password?: string;
};
