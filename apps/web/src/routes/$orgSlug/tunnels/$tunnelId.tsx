import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../../components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/tunnels/$tunnelId')({
  component: TunnelDetailsPage,
});

function TunnelDetailsPage() {
  return (
    <PagePlaceholder
      title="Tunnel details"
      description="Inspect connection state, endpoint configuration, and request activity."
    />
  );
}
