import type { Entity } from '#/interfaces/api';

export type DomainStatus =
  | 'pending'
  | 'verified'
  | 'active'
  | 'failed'
  | 'revoked';

export type Domain = Entity & {
  organizationId: string;
  tunnelId?: string;
  hostname: string;
  status: DomainStatus;
  verificationMethod: string;
  verifiedAt?: string;
  certificateStatus: string;
  certificateExpiresAt?: string;
};
