import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/api-keys')({
  component: ApiKeysPage,
});

function ApiKeysPage() {
  return (
    <PagePlaceholder
      title="API keys"
      description="Create, scope, rotate, and revoke organization API keys."
    />
  );
}
