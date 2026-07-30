import { useEffect, useState } from 'react';

export function useAuthNotice() {
  const [notice, setNotice] = useState<string | null>(null);

  useEffect(() => {
    const error = new URLSearchParams(window.location.search).get('error');
    if (error === 'oauth_failed') {
      setNotice('We could not complete that sign-in. Please try again.');
    }
    if (error === 'oauth_start_failed') {
      setNotice(
        'That sign-in provider is unavailable right now. Please try again later.',
      );
    }
  }, []);

  return notice;
}
