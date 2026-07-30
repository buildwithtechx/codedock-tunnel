import { Cable } from 'lucide-react';
import {
  SiAstro,
  SiDrizzle,
  SiExpress,
  SiGraphql,
  SiMongodb,
  SiNestjs,
  SiNextdotjs,
  SiNodedotjs,
  SiPostgresql,
  SiPrisma,
  SiReact,
  SiRedis,
  SiSocketdotio,
  SiSolid,
  SiStripe,
  SiSvelte,
  SiSwagger,
  SiTailwindcss,
  SiTrpc,
  SiTypescript,
  SiVercel,
  SiVite,
  SiVuedotjs,
} from 'react-icons/si';

export type PluginId = 'sdk' | 'react' | 'vite' | 'next' | 'nest' | 'express';

export type PluginDefinition = {
  id: PluginId;
  name: string;
  packageName: string;
  eyebrow: string;
  headline: string;
  description: string;
  docsSlug: string;
  install: string;
  fileName: string;
  code: string;
  colorClass: string;
  icon: React.ComponentType<{ className?: string }>;
  features: string[];
  useCases: string[];
  technologies: {
    label: string;
    icon: React.ComponentType<{ className?: string }>;
  }[];
};

const typescript = { label: 'TypeScript', icon: SiTypescript };

export const pluginDefinitions: PluginDefinition[] = [
  {
    id: 'sdk',
    name: 'TypeScript SDK',
    packageName: '@codedock/sdk',
    eyebrow: 'Framework-neutral foundation',
    headline: 'Put tunnel lifecycle inside your application.',
    description:
      'A small browser and Node.js client for applications that need direct control over authentication, tunnel creation, status, and shutdown.',
    docsSlug: 'sdk',
    install: 'npm install @codedock/sdk',
    fileName: 'server.ts',
    code: "import { CodedockClient } from '@codedock/sdk';\n\nconst client = new CodedockClient({\n  apiKey: process.env.CODEDOCK_API_KEY,\n});\n\nawait client.openTunnel({ protocol: 'http', localPort: 3000 });",
    colorClass: 'text-cyan-300',
    icon: Cable,
    features: [
      'Browser and Node.js support',
      'Typed protocol lifecycle',
      'Reconnect-aware state',
      'Fetch-based transport',
    ],
    useCases: [
      'Preview environments',
      'Custom developer tooling',
      'CI pipeline jobs',
    ],
    technologies: [typescript, { label: 'Node.js', icon: SiNodedotjs }],
  },
  {
    id: 'react',
    name: 'React',
    packageName: '@codedock/react',
    eyebrow: 'React integration',
    headline: 'Tunnel state that feels native to React.',
    description:
      'Provider and hooks for showing connection state, public URLs, and tunnel controls directly in your React application.',
    docsSlug: 'react',
    install: 'npm install @codedock/react',
    fileName: 'app.tsx',
    code: "import { CodedockTunnelProvider, useTunnel } from '@codedock/react';\n\nfunction PreviewStatus() {\n  const { tunnel, status } = useTunnel();\n  return <span>{tunnel?.publicUrl ?? status}</span>;\n}",
    colorClass: 'text-sky-300',
    icon: SiReact,
    features: [
      'Provider and hooks',
      'Reactive tunnel status',
      'Typed lifecycle actions',
      'Works with React 18+',
    ],
    useCases: [
      'Preview dashboards',
      'Internal developer portals',
      'Live connection controls',
    ],
    technologies: [
      { label: 'React', icon: SiReact },
      typescript,
      { label: 'Next.js', icon: SiNextdotjs },
    ],
  },
  {
    id: 'vite',
    name: 'Vite',
    packageName: '@codedock/vite-plugin',
    eyebrow: 'Vite integration',
    headline: 'Turn every Vite preview into a shareable URL.',
    description:
      'The development server integration opens a tunnel when Vite is ready and keeps the local target aligned with the running server.',
    docsSlug: 'vite',
    install: 'npm install -D @codedock/vite-plugin',
    fileName: 'vite.config.ts',
    code: "import { defineConfig } from 'vite';\nimport react from '@vitejs/plugin-react';\nimport codedock from '@codedock/vite-plugin';\n\nexport default defineConfig({\n  plugins: [react(), codedock()],\n});",
    colorClass: 'text-indigo-300',
    icon: SiVite,
    features: [
      'Starts with the dev server',
      'Dynamic port awareness',
      'HMR-friendly lifecycle',
      'React, Vue, Svelte, and Solid',
    ],
    useCases: ['Design reviews', 'Webhook callbacks', 'Remote QA sessions'],
    technologies: [
      { label: 'React', icon: SiReact },
      { label: 'Vue', icon: SiVuedotjs },
      { label: 'Svelte', icon: SiSvelte },
      { label: 'Solid', icon: SiSolid },
      { label: 'Astro', icon: SiAstro },
    ],
  },
  {
    id: 'next',
    name: 'Next.js',
    packageName: '@codedock/next',
    eyebrow: 'Next.js integration',
    headline: 'Give your Next.js app an edge it can reach.',
    description:
      'A lifecycle wrapper for Next.js development and server workflows, with the same tunnel controls as the CLI and SDK.',
    docsSlug: 'next',
    install: 'npm install @codedock/next',
    fileName: 'next.config.ts',
    code: "import withCodedock from '@codedock/next';\n\nexport default withCodedock({\n  reactStrictMode: true,\n});",
    colorClass: 'text-white',
    icon: SiNextdotjs,
    features: [
      'App and Pages Router support',
      'Development lifecycle hooks',
      'Server-friendly configuration',
      'Typed configuration',
    ],
    useCases: ['Preview deployments', 'OAuth callback testing', 'Client demos'],
    technologies: [
      { label: 'Vercel', icon: SiVercel },
      { label: 'Prisma', icon: SiPrisma },
      { label: 'Tailwind', icon: SiTailwindcss },
      typescript,
      { label: 'tRPC', icon: SiTrpc },
      { label: 'Drizzle', icon: SiDrizzle },
    ],
  },
  {
    id: 'nest',
    name: 'NestJS',
    packageName: '@codedock/nest',
    eyebrow: 'NestJS integration',
    headline: 'Expose your Nest application without extra plumbing.',
    description:
      'A Nest module and service that make tunnel lifecycle part of your application bootstrap and shutdown flow.',
    docsSlug: 'nest',
    install: 'npm install @codedock/nest',
    fileName: 'app.module.ts',
    code: "import { CodedockModule } from '@codedock/nest';\n\n@Module({\n  imports: [CodedockModule.forRoot({ localPort: 3000 })],\n})\nexport class AppModule {}",
    colorClass: 'text-rose-300',
    icon: SiNestjs,
    features: [
      'Module-based setup',
      'Lifecycle-aware service',
      'Automatic shutdown cleanup',
      'TypeScript-first API',
    ],
    useCases: ['Webhook development', 'Team API previews', 'Staging callbacks'],
    technologies: [
      { label: 'PostgreSQL', icon: SiPostgresql },
      { label: 'MongoDB', icon: SiMongodb },
      { label: 'Redis', icon: SiRedis },
      { label: 'GraphQL', icon: SiGraphql },
      { label: 'Swagger', icon: SiSwagger },
    ],
  },
  {
    id: 'express',
    name: 'Express',
    packageName: '@codedock/express',
    eyebrow: 'Express integration',
    headline: 'Make an Express server reachable in one middleware.',
    description:
      'A lightweight lifecycle wrapper for Express servers, designed for APIs, webhooks, and services that already own their HTTP process.',
    docsSlug: 'express',
    install: 'npm install @codedock/express',
    fileName: 'server.ts',
    code: "import express from 'express';\nimport { codedockTunnel } from '@codedock/express';\n\nconst app = express();\ncodedockTunnel(app, { localPort: 3000 });\napp.listen(3000);",
    colorClass: 'text-amber-300',
    icon: SiExpress,
    features: [
      'Minimal middleware setup',
      'Idempotent lifecycle',
      'Status endpoint helpers',
      'Express 4 and 5 support',
    ],
    useCases: [
      'Webhook inspection',
      'Partner API demos',
      'Local service sharing',
    ],
    technologies: [
      { label: 'MongoDB', icon: SiMongodb },
      { label: 'PostgreSQL', icon: SiPostgresql },
      { label: 'Redis', icon: SiRedis },
      { label: 'Socket.IO', icon: SiSocketdotio },
      { label: 'Stripe', icon: SiStripe },
    ],
  },
];

export function getPluginDefinition(pluginId: string) {
  return pluginDefinitions.find((plugin) => plugin.id === pluginId);
}
