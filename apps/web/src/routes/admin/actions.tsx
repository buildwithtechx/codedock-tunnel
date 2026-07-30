import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../components/layout/page-placeholder';

export const Route = createFileRoute('/admin/actions')({
  component: AdminActionsPage,
});

function AdminActionsPage() {
  return (
    <PagePlaceholder
      title="Admin actions"
      description="Run controlled platform operations with an audit trail."
    />
  );
}
