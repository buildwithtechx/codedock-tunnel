import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../../components/layout/page-placeholder';

export const Route = createFileRoute('/admin/users/$userId')({
  component: AdminUserPage,
});

function AdminUserPage() {
  return (
    <PagePlaceholder
      title="User details"
      description="Review identity, sessions, organizations, and activity."
    />
  );
}
