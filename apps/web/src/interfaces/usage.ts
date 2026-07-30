import type { Entity } from '#/interfaces/api';

export type UsageSnapshot = Entity & {
  organizationId: string;
  periodStart: string;
  periodEnd: string;
  tunnelCount: number;
  activeConnections: number;
  bandwidthBytes: number;
  requestCount: number;
  errorCount: number;
};

export type UsageEvent = Entity & {
  organizationId: string;
  tunnelId?: string;
  eventType: string;
  bytes: number;
  connections: number;
  method?: string;
  path?: string;
  statusCode?: number;
  durationMillis?: number;
  responseBytes?: number;
  clientIp?: string;
  occurredAt: string;
};
