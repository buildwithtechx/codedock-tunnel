import { type OpenTunnelAck, RelayConnection } from '@codedock-tunnel/sdk-ts';
import type { OnModuleDestroy, OnModuleInit } from '@nestjs/common';
import { Inject, Injectable } from '@nestjs/common';
import {
  CODEDOCK_TUNNEL_OPTIONS,
  type NestTunnelOptions,
} from '../interfaces/options';

@Injectable()
export class CodedockTunnelService implements OnModuleInit, OnModuleDestroy {
  private readonly connection: RelayConnection;
  private tunnel?: OpenTunnelAck;

  constructor(
    @Inject(CODEDOCK_TUNNEL_OPTIONS)
    private readonly options: NestTunnelOptions,
  ) {
    this.connection = new RelayConnection(options);
  }

  async onModuleInit(): Promise<void> {
    if (this.options.autoStart !== false) {
      await this.start();
    }
  }

  async onModuleDestroy(): Promise<void> {
    await this.stop('Nest module destroyed');
  }

  async start(): Promise<OpenTunnelAck> {
    this.tunnel = await this.connection.openTunnel({
      local_port: this.options.localPort,
      protocol: 'http',
      subdomain: this.options.subdomain,
      password: this.options.password,
    });
    return this.tunnel;
  }

  async stop(reason?: string): Promise<void> {
    if (this.tunnel) {
      await this.connection.closeTunnel(this.tunnel.tunnel_id, reason);
      this.tunnel = undefined;
    }
    this.connection.close();
  }

  status(): OpenTunnelAck | undefined {
    return this.tunnel;
  }
}
