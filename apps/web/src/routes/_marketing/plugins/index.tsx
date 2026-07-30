import { createFileRoute } from '@tanstack/react-router';
import { PluginsPage } from '#/features/marketing';
import { createSeo } from '#/lib/seo';

export const Route = createFileRoute('/_marketing/plugins/')({
  head: () =>
    createSeo({
      title: 'Plugins and SDKs — Codedock Tunnel',
      description:
        'Connect Codedock Tunnel to React, Vite, Next.js, NestJS, and Express.',
      path: '/plugins',
    }),
  component: PluginsPage,
});
