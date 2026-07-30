import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '#/components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/billing')({
  component: BillingPage,
});

function BillingPage() {
  return (
    <PagePlaceholder
      title="Billing"
      description="Manage the organization plan, payment, and subscription state."
    />
  );
}
