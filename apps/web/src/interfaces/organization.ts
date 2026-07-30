import type { Entity } from '#/interfaces/api';

export type Organization = Entity & {
  name: string;
  slug: string;
  ownerId: string;
  settings: string;
};

export type MemberRole = 'owner' | 'admin' | 'member' | 'viewer';

export type OrganizationMember = Entity & {
  organizationId: string;
  userId: string;
  role: MemberRole;
};

export type OrganizationInvitation = Entity & {
  organizationId: string;
  inviterId: string;
  email: string;
  role: Exclude<MemberRole, 'owner'>;
  expiresAt: string;
  acceptedAt?: string;
};
