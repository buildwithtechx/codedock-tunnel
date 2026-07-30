import { createFileRoute, Outlet, redirect } from '@tanstack/react-router';
import { requirePlatformAdmin } from '#/lib/route-guards';
import { getSession } from '#/lib/session';

export const Route = createFileRoute('/admin')({
  beforeLoad: async ({ context }) => {
    try {
      const session = await context.queryClient.fetchQuery({
        queryKey: ['auth', 'session'],
        queryFn: getSession,
        retry: false,
      });
      requirePlatformAdmin(session);
    } catch (error) {
      if (
        error instanceof Error &&
        error.message === 'platform admin access is required'
      ) {
        throw redirect({ to: '/' });
      }
      throw redirect({ to: '/login' });
    }
  },
  component: AdminLayout,
});

function AdminLayout() {
  return <Outlet />;
}
