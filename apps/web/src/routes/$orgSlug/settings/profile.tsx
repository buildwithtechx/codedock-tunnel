import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../../../components/layout/page-placeholder';

export const Route = createFileRoute('/$orgSlug/settings/profile')({
  component: ProfileSettingsPage,
});

function ProfileSettingsPage() {
  return (
    <PagePlaceholder
      title="Profile settings"
      description="Update your profile and connected identity details."
    />
  );
}
