import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../components/layout/page-placeholder';

export const Route = createFileRoute('/admin/usage')({
  component: AdminUsagePage,
});

function AdminUsagePage() {
  return (
    <PagePlaceholder
      title="Usage"
      description="Review platform-wide traffic, requests, errors, and retention."
    />
  );
}
