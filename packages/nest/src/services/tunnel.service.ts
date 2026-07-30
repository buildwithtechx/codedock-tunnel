import { type OpenTunnelAck, RelayConnection } from '@codedock/sdk';
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
  private startPromise?: Promise<OpenTunnelAck>;
  private generation = 0;

  constructor(
    @Inject(CODEDOCK_TUNNEL_OPTIONS)
    private readonly options: NestTunnelOptions,
  ) {
    this.connection = new RelayConnection(options);
    this.connection.on('tunnel_opened', (tunnel) => {
      this.tunnel = tunnel;
    });
    this.connection.on('disconnected', () => {
      this.tunnel = undefined;
    });
  }

  async onModuleInit(): Promise<void> {
    if (this.options.autoStart !== false) {
      await this.start();
    }
  }

  async onModuleDestroy(): Promise<void> {
    try {
      await this.stop('Nest module destroyed');
    } catch {
      return;
    }
  }

  async start(): Promise<OpenTunnelAck> {
    if (this.tunnel) return this.tunnel;
    if (this.startPromise) return this.startPromise;
    const startGeneration = ++this.generation;
    this.startPromise = this.connection
      .openTunnel({
        local_port: this.options.localPort,
        protocol: 'http',
        subdomain: this.options.subdomain,
        password: this.options.password,
      })
      .then((tunnel) => {
        if (startGeneration === this.generation) this.tunnel = tunnel;
        return tunnel;
      })
      .finally(() => {
        this.startPromise = undefined;
      });
    return this.startPromise;
  }

  async stop(reason?: string): Promise<void> {
    this.generation += 1;
    const tunnel = this.tunnel;
    let closeError: unknown;
    try {
      if (tunnel) await this.connection.closeTunnel(tunnel.tunnel_id, reason);
    } catch (error) {
      closeError = error;
    } finally {
      this.tunnel = undefined;
      this.connection.close();
    }
    if (closeError) throw closeError;
  }

  status(): OpenTunnelAck | undefined {
    return this.tunnel;
  }
}
