import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../components/layout/page-placeholder';

export const Route = createFileRoute('/pricing')({ component: PricingPage });

function PricingPage() {
  return (
    <PagePlaceholder
      title="Pricing"
      description="Choose the plan that fits your tunnel usage."
    />
  );
}
