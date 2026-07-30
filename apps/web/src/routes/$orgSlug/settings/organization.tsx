import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../../components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/settings/organization')({
  component: OrganizationSettingsPage,
});

function OrganizationSettingsPage() {
  return (
    <PagePlaceholder
      title="Organization settings"
      description="Update organization identity, access, and lifecycle settings."
    />
  );
}
