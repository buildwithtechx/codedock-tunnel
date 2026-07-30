import type { DynamicModule } from '@nestjs/common';
import { Module } from '@nestjs/common';
import {
  CODEDOCK_TUNNEL_OPTIONS,
  type NestTunnelOptions,
} from '../interfaces/options';
import { CodedockTunnelService } from './tunnel.service';

@Module({})
export class CodedockTunnelModule {
  static forRoot(options: NestTunnelOptions): DynamicModule {
    return {
      module: CodedockTunnelModule,
      providers: [
        { provide: CODEDOCK_TUNNEL_OPTIONS, useValue: options },
        CodedockTunnelService,
      ],
      exports: [CodedockTunnelService],
    };
  }
}
