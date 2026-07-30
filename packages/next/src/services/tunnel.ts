import { RelayConnection } from '@codedock/sdk';
import type {
  NextTunnel,
  NextTunnelOptions,
  NextTunnelState,
} from '../interfaces/options';

export function createNextTunnel(options: NextTunnelOptions): NextTunnel {
  const connection = new RelayConnection(options);
  let current: NextTunnelState = { status: 'idle' };

  return {
    start: async () => {
      if (options.enabled === false) {
        return current;
      }
      current = { status: 'connecting' };
      try {
        const tunnel = await connection.openTunnel({
          local_port: options.localPort,
          protocol: 'http',
          subdomain: options.subdomain,
          password: options.password,
        });
        current = {
          status: 'active',
          tunnelId: tunnel.tunnel_id,
          publicUrl: tunnel.public_url,
          publicPort: tunnel.public_port,
        };
      } catch (value) {
        current = { status: 'error', error: normalizeError(value) };
      }
      return current;
    },
    stop: async (reason) => {
      if (current.tunnelId) {
        await connection.closeTunnel(current.tunnelId, reason);
      }
      connection.close();
      current = { status: 'closed' };
    },
    state: () => current,
  };
}

function normalizeError(value: unknown): Error {
  return value instanceof Error ? value : new Error(String(value));
}
