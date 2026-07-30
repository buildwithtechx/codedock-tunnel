import type { OpenTunnel, OpenTunnelAck } from '@codedock-tunnel/protocol-ts';
import type {
  RelayConnection,
  RelayConnectionOptions,
} from '@codedock-tunnel/sdk-ts';

export type TunnelProviderProps = {
  options: RelayConnectionOptions;
  children: React.ReactNode;
};

export type TunnelStatus =
  | 'idle'
  | 'connecting'
  | 'connected'
  | 'error'
  | 'closed';

export type TunnelContextValue = {
  connection: RelayConnection;
  status: TunnelStatus;
  tunnel?: OpenTunnelAck;
  error?: Error;
  connect: () => Promise<void>;
  disconnect: () => void;
  openTunnel: (request: Omit<OpenTunnel, 'token'>) => Promise<OpenTunnelAck>;
  closeTunnel: (tunnelId?: string, reason?: string) => Promise<void>;
};
