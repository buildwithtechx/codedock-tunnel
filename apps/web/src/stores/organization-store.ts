import { create } from 'zustand';
import type { Organization } from '../interfaces/organization';

type OrganizationStore = {
  organization?: Organization;
  setOrganization: (organization?: Organization) => void;
};

export const useOrganizationStore = create<OrganizationStore>((set) => ({
  setOrganization: (organization) => set({ organization }),
}));
