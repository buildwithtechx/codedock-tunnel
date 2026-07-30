import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '#/components/layout/page-placeholder';

export const Route = createFileRoute('/admin/organizations/')({
  component: AdminOrganizationsPage,
});

function AdminOrganizationsPage() {
  return (
    <PagePlaceholder
      title="Organizations"
      description="Review organizations and their platform usage."
    />
  );
}
