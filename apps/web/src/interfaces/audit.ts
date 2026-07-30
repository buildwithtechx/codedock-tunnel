import type { Entity } from '#/interfaces/api';

export type AuditEvent = Entity & {
  organizationId?: string;
  userId?: string;
  action: string;
  resourceType: string;
  resourceId?: string;
  ipAddress?: string;
  userAgent?: string;
  metadata: string;
  occurredAt: string;
};
