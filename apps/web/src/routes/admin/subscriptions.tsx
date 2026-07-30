import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../components/layout/page-placeholder';

export const Route = createFileRoute('/admin/subscriptions')({
  component: AdminSubscriptionsPage,
});

function AdminSubscriptionsPage() {
  return (
    <PagePlaceholder
      title="Subscriptions"
      description="Review plans, subscriptions, and billing state."
    />
  );
}
