import type { Entity } from '#/interfaces/api';

export type OAuthIdentity = Entity & {
  userId: string;
  provider: string;
  subject: string;
  email?: string;
};

export type DeviceLogin = Entity & {
  userId?: string;
  status: 'pending' | 'completed' | 'expired';
  expiresAt: string;
  completedAt?: string;
  ipAddress?: string;
};

export type APIKey = Entity & {
  userId: string;
  organizationId?: string;
  name: string;
  prefix: string;
  scopes: string;
  expiresAt?: string;
  lastUsedAt?: string;
  revokedAt?: string;
};
