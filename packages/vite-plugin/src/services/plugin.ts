import { RelayConnection } from '@codedock-tunnel/sdk-ts';
import type { Plugin, ViteDevServer } from 'vite';
import type { CodedockTunnelPluginOptions } from '../interfaces/options';

export function codedockTunnel(options: CodedockTunnelPluginOptions): Plugin {
  let connection: RelayConnection | undefined;

  return {
    name: 'codedock-tunnel',
    apply: 'serve',
    configureServer(server: ViteDevServer) {
      if (options.enabled === false) {
        return;
      }
      return startTunnel(server, options, (nextConnection) => {
        connection = nextConnection;
      });
    },
    closeBundle() {
      connection?.close();
      connection = undefined;
    },
  };
}

async function startTunnel(
  server: ViteDevServer,
  options: CodedockTunnelPluginOptions,
  assign: (connection: RelayConnection) => void,
): Promise<void> {
  const connection = new RelayConnection(options);
  assign(connection);
  const localPort = options.localPort ?? server.config.server.port ?? 5173;
  const tunnel = await connection.openTunnel({
    local_port: localPort,
    protocol: 'http',
    subdomain: options.subdomain,
    password: options.password,
  });
  server.config.logger.info(`Codedock Tunnel: ${tunnel.public_url}`);
  server.httpServer?.once('close', () => connection.close());
}
