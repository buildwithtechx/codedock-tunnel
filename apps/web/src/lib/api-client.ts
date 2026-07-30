import { env } from '../env';
import type { ApiError } from '../interfaces/api';

const apiBaseURL = env.VITE_CODEDOCK_API_BASE_URL ?? 'http://localhost:8080';

export async function apiRequest<T>(
  path: string,
  init?: RequestInit,
): Promise<T> {
  const response = await fetch(`${apiBaseURL}${path}`, {
    ...init,
    credentials: 'include',
    headers: { Accept: 'application/json', ...init?.headers },
  });
  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as ApiError | null;
    throw new Error(
      body?.error ?? `API request failed with status ${response.status}`,
    );
  }
  if (response.status === 204) return undefined as T;
  return response.json() as Promise<T>;
}
