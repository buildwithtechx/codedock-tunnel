import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../../components/layout/page-placeholder';

export const Route = createFileRoute('/admin/users/')({
  component: AdminUsersPage,
});

function AdminUsersPage() {
  return (
    <PagePlaceholder
      title="Users"
      description="Manage platform users and account status."
    />
  );
}
