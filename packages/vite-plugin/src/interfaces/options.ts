import type { RelayConnectionOptions } from '@codedock-tunnel/sdk-ts';

export type CodedockTunnelPluginOptions = RelayConnectionOptions & {
  enabled?: boolean;
  localPort?: number;
  subdomain?: string;
};
