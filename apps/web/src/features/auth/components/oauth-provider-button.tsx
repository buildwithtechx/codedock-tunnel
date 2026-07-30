import { SiGithub, SiGoogle } from 'react-icons/si';
import type { OAuthProvider } from '#/interfaces/auth';

type OAuthProviderButtonProps = {
  provider: OAuthProvider;
  label: string;
  loading: boolean;
  onClick: (provider: OAuthProvider) => void;
};

export function OAuthProviderButton({
  provider,
  label,
  loading,
  onClick,
}: OAuthProviderButtonProps) {
  const Icon = provider === 'google' ? SiGoogle : SiGithub;
  const providerName = provider === 'google' ? 'Google' : 'GitHub';

  return (
    <button
      type="button"
      disabled={loading}
      onClick={() => onClick(provider)}
      className="flex h-12 w-full items-center justify-center gap-3 rounded-xl border border-white/10 bg-white/[0.04] text-sm font-medium text-white transition-colors hover:border-indigo-300/45 hover:bg-indigo-300/10 disabled:opacity-60"
    >
      <Icon className="size-4" />
      {loading ? `Connecting to ${providerName}…` : `${label} ${providerName}`}
    </button>
  );
}
