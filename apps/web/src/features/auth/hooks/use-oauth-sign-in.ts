import { useState } from 'react';
import { startOAuthSignIn } from '#/features/auth/services/auth-service';
import type { OAuthProvider } from '#/interfaces/auth';

export function useOAuthSignIn(returnTo = '/') {
  const [provider, setProvider] = useState<OAuthProvider | null>(null);

  return {
    provider,
    signIn(providerName: OAuthProvider) {
      setProvider(providerName);
      startOAuthSignIn(providerName, returnTo);
    },
  };
}
