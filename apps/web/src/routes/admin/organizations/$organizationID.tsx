import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '#/components/layout/page-placeholder';

export const Route = createFileRoute('/admin/organizations/$organizationID')({
  component: AdminOrganizationPage,
});

function AdminOrganizationPage() {
  return (
    <PagePlaceholder
      title="Organization details"
      description="Inspect members, tunnels, billing, and organization activity."
    />
  );
}
