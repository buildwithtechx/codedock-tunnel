import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '#/components/layout/page-placeholder';

export const Route = createFileRoute('/login')({ component: LoginPage });

function LoginPage() {
  return (
    <PagePlaceholder
      title="Sign in"
      description="Continue with Google or GitHub to manage your tunnels."
    />
  );
}
