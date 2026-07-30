import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '#/components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/agents')({
  component: AgentsPage,
});

function AgentsPage() {
  return (
    <PagePlaceholder
      title="Agents"
      description="Manage connected tunnel agents and their credentials."
    />
  );
}
