import { createFileRoute } from '@tanstack/react-router';
import { PagePlaceholder } from '../components/layout/page-placeholder';

export const Route = createFileRoute('/signup')({ component: SignupPage });

function SignupPage() {
  return (
    <PagePlaceholder
      title="Create your account"
      description="Create a Codedock Tunnel account with your preferred OAuth provider."
    />
  );
}
