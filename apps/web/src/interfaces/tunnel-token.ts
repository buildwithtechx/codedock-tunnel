import type { Entity } from '#/interfaces/api';

export type TunnelToken = Entity & {
  tunnelId: string;
  name: string;
  prefix: string;
  scopes: string;
  expiresAt?: string;
  lastUsedAt?: string;
  revokedAt?: string;
};
