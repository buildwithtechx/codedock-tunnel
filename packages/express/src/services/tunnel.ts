import { RelayConnection } from '@codedock-tunnel/sdk-ts';
import type {
  ExpressTunnel,
  ExpressTunnelOptions,
  ExpressTunnelState,
} from '../interfaces/options';

export function createExpressTunnel(
  options: ExpressTunnelOptions,
): ExpressTunnel {
  const connection = new RelayConnection(options);
  let current: ExpressTunnelState = { status: 'idle' };

  return {
    start: async () => {
      if (options.autoStart === false) {
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
