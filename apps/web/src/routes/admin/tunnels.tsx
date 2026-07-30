import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../components/layout/page-placeholder';

export const Route = createFileRoute('/admin/tunnels')({
  component: AdminTunnelsPage,
});

function AdminTunnelsPage() {
  return (
    <PagePlaceholder
      title="All tunnels"
      description="Monitor tunnel lifecycle and relay health across the platform."
    />
  );
}
