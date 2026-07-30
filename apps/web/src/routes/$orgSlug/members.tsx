import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '#/components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/members')({
  component: MembersPage,
});

function MembersPage() {
  return (
    <PagePlaceholder
      title="Members"
      description="Invite teammates and manage organization roles."
    />
  );
}
