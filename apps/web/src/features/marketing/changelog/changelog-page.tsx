import { ArrowUpRight, CalendarDays } from 'lucide-react';
import { MarketingContainer } from '#/components/layout';

const entries = [
  {
    date: 'July 2026',
    tag: 'Platform',
    title: 'A shared protocol for every client',
    text: 'The CLI, desktop app, SDK, and framework integrations now align on one versioned tunnel lifecycle.',
  },
  {
    date: 'June 2026',
    tag: 'Developer experience',
    title: 'Framework adapters are ready',
    text: 'Vite, Next.js, NestJS, Express, and React integrations make local tunnel lifecycle part of your application workflow.',
  },
  {
    date: 'May 2026',
    tag: 'Observability',
    title: 'See the request path',
    text: 'Tunnel usage now includes request counts, bandwidth, response outcomes, and connection health.',
  },
];

export function ChangelogPage() {
  return (
    <section className="pb-16 pt-28 sm:pt-32">
      <MarketingContainer className="max-w-5xl">
        <div className="max-w-2xl">
          <p className="text-sm font-medium text-cyan-300">What’s new</p>
          <h1 className="mt-4 text-4xl font-semibold tracking-tight sm:text-5xl">
            Small releases. Better tunnels.
          </h1>
          <p className="mt-5 text-lg leading-8 text-white/50">
            Follow the work behind Codedock Tunnel as the protocol, clients, and
            developer surfaces evolve.
          </p>
        </div>
        <div className="mt-12 divide-y divide-white/10 border-y border-white/10">
          {entries.map((entry) => (
            <article
              key={entry.title}
              className="grid gap-5 py-10 md:grid-cols-[150px_1fr_120px]"
            >
              <div className="flex items-start gap-2 text-sm text-white/40">
                <CalendarDays className="size-4" />
                {entry.date}
              </div>
              <div>
                <span className="rounded-full border border-indigo-300/25 bg-indigo-300/10 px-2.5 py-1 text-xs text-indigo-200">
                  {entry.tag}
                </span>
                <h2 className="mt-5 text-2xl font-semibold tracking-tight">
                  {entry.title}
                </h2>
                <p className="mt-3 max-w-2xl leading-7 text-white/50">
                  {entry.text}
                </p>
              </div>
              <a
                href="https://github.com/codedock"
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center gap-1 text-sm text-white/45 hover:text-white"
              >
                View source <ArrowUpRight className="size-4" />
              </a>
            </article>
          ))}
        </div>
      </MarketingContainer>
    </section>
  );
}
