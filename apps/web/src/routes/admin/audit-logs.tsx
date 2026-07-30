import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '#/components/layout/page-placeholder';

export const Route = createFileRoute('/admin/audit-logs')({
  component: AdminAuditLogsPage,
});

function AdminAuditLogsPage() {
  return (
    <PagePlaceholder
      title="Audit logs"
      description="Review sensitive platform and account actions."
    />
  );
}
