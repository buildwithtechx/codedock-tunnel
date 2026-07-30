import { createFileRoute } from '@tanstack/react-router';
import { ChangelogPage } from '#/features/marketing';
import { createSeo } from '#/lib/seo';

export const Route = createFileRoute('/_marketing/changelog')({
  head: () =>
    createSeo({
      title: 'Changelog — Codedock Tunnel',
      description:
        'Follow new releases and improvements across Codedock Tunnel.',
      path: '/changelog',
    }),
  component: ChangelogPage,
});
