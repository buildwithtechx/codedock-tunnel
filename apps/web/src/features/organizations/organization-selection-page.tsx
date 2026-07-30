import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Link, useNavigate } from '@tanstack/react-router';
import { useState } from 'react';
import { AuthPageShell } from '#/features/auth/components/auth-page-shell';
import {
  createOrganization,
  getOrganizations,
} from '#/features/auth/services/auth-service';

function rememberOrganization(slug: string) {
  window.localStorage.setItem('codedock_tunnel_last_organization', slug);
}

export function OrganizationSelectionPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [name, setName] = useState('');
  const [slug, setSlug] = useState('');
  const {
    data: organizations = [],
    isError,
    isLoading,
  } = useQuery({
    queryKey: ['organizations'],
    queryFn: getOrganizations,
  });
  const createMutation = useMutation({
    mutationFn: () => createOrganization(name.trim(), slug.trim()),
    onSuccess: async (organization) => {
      rememberOrganization(organization.slug);
      await queryClient.invalidateQueries({ queryKey: ['organizations'] });
      await navigate({
        to: '/$orgSlug',
        params: { orgSlug: organization.slug },
      });
    },
  });

  return (
    <AuthPageShell
      title="Choose a workspace"
      description="Select the organization you want to open."
      footer={null}
    >
      {isLoading && (
        <p className="text-sm text-white/45">Loading workspaces…</p>
      )}
      {isError && (
        <p className="text-sm leading-6 text-rose-200">
          We could not load your workspaces. Please refresh and try again.
        </p>
      )}
      {!isLoading && organizations.length === 0 && (
        <form
          className="space-y-3 text-left"
          onSubmit={(event) => {
            event.preventDefault();
            if (name.trim() && slug.trim()) createMutation.mutate();
          }}
        >
          <p className="text-sm leading-6 text-white/55">
            Create your first workspace to start opening tunnels.
          </p>
          <input
            aria-label="Workspace name"
            value={name}
            onChange={(event) => setName(event.target.value)}
            placeholder="Workspace name"
            maxLength={120}
            className="h-12 w-full rounded-xl border border-white/10 bg-[#111] px-4 text-sm text-white outline-none placeholder:text-white/30 focus:border-indigo-300/60 focus:ring-2 focus:ring-indigo-300/20"
          />
          <input
            aria-label="Workspace slug"
            value={slug}
            onChange={(event) => setSlug(event.target.value.toLowerCase())}
            placeholder="workspace-slug"
            maxLength={63}
            className="h-12 w-full rounded-xl border border-white/10 bg-[#111] px-4 text-sm text-white outline-none placeholder:text-white/30 focus:border-indigo-300/60 focus:ring-2 focus:ring-indigo-300/20"
          />
          {createMutation.isError && (
            <p className="text-sm text-rose-200">
              We could not create that workspace. Try another name or slug.
            </p>
          )}
          <button
            type="submit"
            disabled={createMutation.isPending || !name.trim() || !slug.trim()}
            className="h-12 w-full rounded-xl bg-indigo-300 px-4 text-sm font-semibold text-black transition-colors hover:bg-indigo-200 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-300 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {createMutation.isPending
              ? 'Creating workspace…'
              : 'Create workspace'}
          </button>
        </form>
      )}
      <div className="space-y-3 text-left">
        {organizations.map((organization) => (
          <Link
            key={organization.id}
            to="/$orgSlug"
            params={{ orgSlug: organization.slug }}
            onClick={() => rememberOrganization(organization.slug)}
            className="block rounded-xl border border-white/10 bg-[#111] px-5 py-4 transition-colors hover:border-indigo-300/45 hover:bg-indigo-300/10 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-300"
          >
            <span className="block font-medium text-white">
              {organization.name}
            </span>
            <span className="mt-1 block text-sm text-white/40">
              {organization.slug}
            </span>
          </Link>
        ))}
      </div>
    </AuthPageShell>
  );
}
