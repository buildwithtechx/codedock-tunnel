import type { RelayConnectionOptions } from '@codedock/sdk-ts';

export type CodedockTunnelPluginOptions = RelayConnectionOptions & {
  enabled?: boolean;
  localPort?: number;
  subdomain?: string;
  password?: string;
};
