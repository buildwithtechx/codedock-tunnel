import { describe, expect, it } from 'vitest';
import { CodedockTunnelModule } from '../services/tunnel.module';
import { CodedockTunnelService } from '../services/tunnel.service';

describe('CodedockTunnelModule', () => {
  it('registers its options and lifecycle service', () => {
    const dynamicModule = CodedockTunnelModule.forRoot({
      relayUrl: 'wss://relay.test',
      agentToken: 'token',
      localPort: 3000,
    });
    expect(dynamicModule.module).toBe(CodedockTunnelModule);
    expect(dynamicModule.exports).toContain(CodedockTunnelService);
    expect(dynamicModule.providers).toHaveLength(2);
  });
});
