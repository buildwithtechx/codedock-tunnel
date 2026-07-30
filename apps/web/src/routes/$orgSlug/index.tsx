import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/')({
  component: OrganizationOverviewPage,
});

function OrganizationOverviewPage() {
  return (
    <PagePlaceholder
      title="Organization overview"
      description="See tunnel health, recent activity, and usage for this organization."
    />
  );
}
