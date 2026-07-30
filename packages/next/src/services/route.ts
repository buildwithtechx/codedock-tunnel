import type { NextTunnel } from '../interfaces/options';

export function createTunnelRoute(tunnel: NextTunnel) {
  return {
    GET: () => Response.json(tunnel.state()),
    DELETE: async () => {
      await tunnel.stop('route requested shutdown');
      return new Response(null, { status: 204 });
    },
  };
}
