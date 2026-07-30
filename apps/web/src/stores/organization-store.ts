import { create } from 'zustand';
import type { Organization } from '#/interfaces/organization';

export type OrganizationState = {
  organizations: Organization[];
  organization: Organization | null;
  setOrganizations: (organizations: Organization[]) => void;
  setOrganization: (organization: Organization | null) => void;
  clear: () => void;
};

export const useOrganizationStore = create<OrganizationState>((set) => ({
  organizations: [],
  organization: null,
  setOrganizations: (organizations) => set({ organizations }),
  setOrganization: (organization) => set({ organization }),
  clear: () => set({ organizations: [], organization: null }),
}));
