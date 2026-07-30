import type { AuthSession, OAuthProvider } from '#/interfaces/auth';
import { apiClient, getApiBaseURL } from '#/lib/api-client';

export function getAuthSession() {
  return apiClient.get<AuthSession>('/api/v1/auth/session');
}

export function startOAuthSignIn(provider: OAuthProvider, returnTo = '/') {
  const url = new URL(`/api/v1/auth/oauth/${provider}`, `${getApiBaseURL()}/`);
  url.searchParams.set('return_to', returnTo);
  window.location.assign(url);
}
