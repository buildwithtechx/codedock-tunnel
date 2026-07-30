import type { Entity } from '#/interfaces/api';

export type AgentStatus = 'pending' | 'online' | 'offline' | 'revoked';

export type Agent = Entity & {
  organizationId: string;
  name: string;
  status: AgentStatus;
  version?: string;
  hostname?: string;
  platform?: string;
  lastSeenAt?: string;
  connectedAt?: string;
  revokedAt?: string;
  metadata: string;
};
