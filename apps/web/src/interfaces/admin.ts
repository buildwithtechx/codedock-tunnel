import type { Entity } from '#/interfaces/api';

export type PlatformAdminRole = 'owner' | 'admin';

export type PlatformAdmin = Entity & {
  userId: string;
  name: string;
  role: PlatformAdminRole;
  active: boolean;
};

export type AdminUsage = {
  bandwidthBytes: number;
  requestCount: number;
  errorCount: number;
};
