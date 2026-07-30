import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '#/components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/domains')({
  component: DomainsPage,
});

function DomainsPage() {
  return (
    <PagePlaceholder
      title="Domains"
      description="Configure and verify custom tunnel domains."
    />
  );
}
