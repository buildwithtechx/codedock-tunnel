import { useQuery } from '@tanstack/react-query';
import { getSession } from '#/lib/auth';

export function useSession() {
  return useQuery({
    queryKey: ['auth', 'session'],
    queryFn: getSession,
    retry: false,
  });
}
