import type { AuthSession } from '#/interfaces/auth';
import { apiRequest } from '#/lib/api-client';

export function getSession(): Promise<AuthSession> {
  return apiRequest<AuthSession>('/api/v1/auth/session');
}
