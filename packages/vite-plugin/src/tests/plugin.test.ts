import { describe, expect, it } from 'vitest';
import { codedockTunnel } from '../services/plugin';

describe('codedockTunnel', () => {
  it('is enabled only for Vite development by default', () => {
    const plugin = codedockTunnel({
      relayUrl: 'wss://relay.test',
      agentToken: 'token',
    });
    expect(plugin.name).toBe('codedock-tunnel');
    expect(plugin.apply).toBe('serve');
  });

  it('does not configure a server when disabled', () => {
    const plugin = codedockTunnel({
      relayUrl: 'wss://relay.test',
      agentToken: 'token',
      enabled: false,
    });
    const configureServer = plugin.configureServer;
    if (typeof configureServer !== 'function') {
      throw new Error('configureServer hook is required');
    }
    expect(configureServer.call({} as never, {} as never)).toBeUndefined();
  });
});
