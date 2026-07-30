import { useOrganizationStore } from '#/stores/organization-store';

export function useOrganization() {
  return useOrganizationStore((state) => state.organization);
}
