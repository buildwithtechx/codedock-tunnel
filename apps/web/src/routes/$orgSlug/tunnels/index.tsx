import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../../components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/tunnels/')({
  component: TunnelsPage,
});

function TunnelsPage() {
  return (
    <PagePlaceholder
      title="Tunnels"
      description="Create, inspect, and manage organization tunnels."
    />
  );
}
