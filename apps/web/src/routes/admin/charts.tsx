import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../components/layout/page-placeholder';

export const Route = createFileRoute('/admin/charts')({
  component: AdminChartsPage,
});

function AdminChartsPage() {
  return (
    <PagePlaceholder
      title="Charts"
      description="Explore operational and business trends."
    />
  );
}
