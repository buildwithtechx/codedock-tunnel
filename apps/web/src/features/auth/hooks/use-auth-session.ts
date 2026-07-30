import { useQuery } from '@tanstack/react-query';
import { useEffect } from 'react';
import { getAuthSession } from '#/features/auth/services/auth-service';
import { ApiError } from '#/lib/api-client';
import { useAuthStore } from '#/stores/auth-store';

export function useAuthSession() {
  const setUser = useAuthStore((state) => state.setUser);
  const clear = useAuthStore((state) => state.clear);
  const query = useQuery({
    queryKey: ['auth', 'session'],
    queryFn: getAuthSession,
    enabled: typeof window !== 'undefined',
    retry: false,
    refetchOnMount: 'always',
  });

  useEffect(() => {
    if (query.data) setUser(query.data.user);
    else if (query.isError) clear();
  }, [clear, query.data, query.isError, setUser]);

  return {
    ...query,
    isAuthenticated: query.data !== undefined,
    isSessionUnavailable:
      query.isError &&
      (!(query.error instanceof ApiError) || query.error.status !== 401),
    user: query.data?.user ?? null,
  };
}
