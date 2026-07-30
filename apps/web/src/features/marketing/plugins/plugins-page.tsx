import { Link } from '@tanstack/react-router';
import { ArrowRight, Blocks, Box, Braces, Code2, Layers3 } from 'lucide-react';
import { MarketingContainer } from '#/components/layout';

const plugins = [
  {
    name: 'React',
    packageName: '@codedock/react',
    description: 'Hooks and providers for tunnel state in React applications.',
    icon: Braces,
  },
  {
    name: 'Vite',
    packageName: '@codedock/vite-plugin',
    description: 'Open a local development tunnel as Vite starts.',
    icon: Layers3,
  },
  {
    name: 'Next.js',
    packageName: '@codedock/next',
    description: 'Server and development lifecycle integration for Next.js.',
    icon: Code2,
  },
  {
    name: 'NestJS',
    packageName: '@codedock/nest',
    description: 'Inject tunnel management into NestJS modules and services.',
    icon: Box,
  },
  {
    name: 'Express',
    packageName: '@codedock/express',
    description: 'Add tunnel lifecycle and status to an Express server.',
    icon: Blocks,
  },
];

export function PluginsPage() {
  return (
    <section className="pt-36 pb-24 sm:pt-44">
      <MarketingContainer>
        <div className="max-w-2xl">
          <p className="text-sm font-medium text-cyan-300">Plugins and SDKs</p>
          <h1 className="mt-4 text-4xl font-semibold tracking-tight sm:text-6xl">
            Bring a public endpoint into your stack.
          </h1>
          <p className="mt-5 text-lg leading-8 text-white/50">
            Every adapter is a thin layer over the same framework-neutral SDK.
            Choose the integration that matches your application and keep the
            protocol out of your way.
          </p>
        </div>
        <div className="mt-14 grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {plugins.map(({ name, packageName, description, icon: Icon }) => (
            <article
              key={name}
              className="group rounded-2xl border border-white/10 bg-white/[0.03] p-6 transition-colors hover:border-indigo-300/40 hover:bg-indigo-300/[0.06]"
            >
              <div className="flex items-center justify-between">
                <Icon className="size-6 text-indigo-300" />
                <span className="font-mono text-xs text-white/35">
                  {packageName}
                </span>
              </div>
              <h2 className="mt-10 text-xl font-semibold">{name}</h2>
              <p className="mt-3 min-h-12 text-sm leading-6 text-white/50">
                {description}
              </p>
              <Link
                to="/docs/$"
                params={{
                  _splat:
                    name === 'Next.js'
                      ? 'next'
                      : name === 'NestJS'
                        ? 'nest'
                        : name.toLowerCase(),
                }}
                className="mt-7 inline-flex items-center gap-2 text-sm font-medium text-white/70 group-hover:text-white"
              >
                Read integration guide <ArrowRight className="size-4" />
              </Link>
            </article>
          ))}
        </div>
        <div className="mt-16 rounded-2xl border border-white/10 bg-[#0d1220] p-7 sm:p-10">
          <p className="font-mono text-sm text-cyan-300">
            one protocol, every surface
          </p>
          <h2 className="mt-4 text-2xl font-semibold">
            Install the adapter you need
          </h2>
          <p className="mt-3 max-w-2xl leading-7 text-white/50">
            Framework packages bring the core SDK as a dependency. You do not
            need to install the framework adapter and the SDK separately.
          </p>
        </div>
      </MarketingContainer>
    </section>
  );
}
