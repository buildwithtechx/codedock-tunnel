import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '#/components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/settings/')({
  component: SettingsPage,
});

function SettingsPage() {
  return (
    <PagePlaceholder
      title="Settings"
      description="Manage organization and personal configuration."
    />
  );
}
