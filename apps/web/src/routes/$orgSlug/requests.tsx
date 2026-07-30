import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/requests')({
  component: RequestsPage,
});

function RequestsPage() {
  return (
    <PagePlaceholder
      title="Requests"
      description="Inspect request activity flowing through organization tunnels."
    />
  );
}
