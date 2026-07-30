import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '#/components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/usage')({
  component: UsagePage,
});

function UsagePage() {
  return (
    <PagePlaceholder
      title="Usage"
      description="Track bandwidth, requests, connections, and retention."
    />
  );
}
