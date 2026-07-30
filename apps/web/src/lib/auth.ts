import type { AuthSession } from '../interfaces/auth';
import { apiRequest } from './api-client';

export function getSession() {
  return apiRequest<AuthSession>('/api/v1/auth/session');
}

export function logout() {
  return apiRequest<void>('/api/v1/auth/logout', { method: 'POST' });
}
