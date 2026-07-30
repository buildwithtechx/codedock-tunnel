import { Link } from '@tanstack/react-router';
import { Activity, Code2, Terminal } from 'lucide-react';
import { useEffect, useState } from 'react';
import { MarketingContainer } from '#/components/layout';

const traffic = [
  ['200 OK', 'GET', '/api/users', '12ms'],
  ['201 Created', 'POST', '/api/webhooks', '45ms'],
  ['401 Unauth', 'GET', '/admin', '8ms'],
  ['200 OK', 'GET', '/api/settings', '24ms'],
  ['500 Error', 'POST', '/api/checkout', '120ms'],
];

export function DeveloperExperience() {
  const [offset, setOffset] = useState(0);
  useEffect(() => {
    const timer = setInterval(() => setOffset((value) => value + 1), 1500);
    return () => clearInterval(timer);
  }, []);
  const liveTraffic = [...traffic, ...traffic].slice(
    offset % traffic.length,
    (offset % traffic.length) + 5,
  );
  return (
    <section className="bg-black py-20 sm:py-24">
      <MarketingContainer>
        <h2 className="mx-auto max-w-3xl text-center text-4xl font-bold tracking-[-0.05em] sm:text-6xl">
          Your workflow,
          <br />
          already connected
        </h2>
        <div className="mt-16 grid gap-8 lg:grid-cols-2">
          <div className="grid gap-8">
            <article className="group flex flex-col rounded-3xl border border-white/5 bg-white/[0.02] p-8 transition-colors hover:border-white/10">
              <div className="flex items-center gap-4">
                <span className="flex size-10 items-center justify-center rounded-full bg-indigo-300/10 transition-colors group-hover:bg-indigo-300/20">
                  <Terminal className="size-5 text-indigo-300" />
                </span>
                <h3 className="text-xl font-semibold">Online in one command</h3>
              </div>
              <p className="mt-6 text-white/45">
                One command and your local service has a public endpoint. No
                reverse proxy configuration required.
              </p>
              <div className="mt-8 rounded-2xl border border-white/10 bg-black/30 p-4 font-mono text-sm text-white/65">
                <span className="text-indigo-300">$</span> codedock-tunnel 3000
              </div>
            </article>
            <article className="group rounded-3xl border border-white/5 bg-white/[0.02] p-8 transition-colors hover:border-white/10">
              <div className="flex items-center gap-4">
                <span className="flex size-10 items-center justify-center rounded-full bg-indigo-300/10 transition-colors group-hover:bg-indigo-300/20">
                  <Code2 className="size-5 text-indigo-300" />
                </span>
                <h3 className="text-xl font-semibold">
                  Prefer your framework?
                </h3>
              </div>
              <p className="mt-6 text-white/45">
                Use a thin adapter in Vite, Next.js, NestJS, Express, or React
                while the same SDK handles the tunnel lifecycle.
              </p>
              <div className="mt-8 grid grid-cols-2 gap-2 font-mono text-xs sm:grid-cols-4">
                {[
                  '@codedock/react',
                  '@codedock/vite-plugin',
                  '@codedock/next',
                  '@codedock/express',
                ].map((name) => (
                  <Link
                    key={name}
                    to="/plugins"
                    className="rounded-lg border border-white/10 bg-black/25 p-3 text-white/55 transition-colors hover:border-cyan-300/35 hover:text-white"
                  >
                    {name}
                  </Link>
                ))}
              </div>
            </article>
          </div>
          <article className="group relative flex h-full flex-col overflow-hidden rounded-3xl border border-white/5 bg-white/[0.02] p-8 transition-colors hover:border-white/10">
            <div className="absolute -right-24 -top-24 size-64 rounded-full bg-indigo-400/10 blur-3xl" />
            <div className="relative">
              <div className="flex items-center gap-4">
                <span className="flex size-10 items-center justify-center rounded-full bg-indigo-300/10 transition-colors group-hover:bg-indigo-300/20">
                  <Activity className="size-5 text-indigo-300" />
                </span>
                <h3 className="text-xl font-semibold">Instant observability</h3>
              </div>
              <p className="mt-6 text-white/45">
                See live traffic as soon as the tunnel comes online, with
                status, path, duration, and response outcomes.
              </p>
            </div>
            <div className="relative mt-auto space-y-3 pt-10 font-mono text-xs">
              {liveTraffic.map(([status, method, path, time]) => (
                <div
                  key={`${status}-${path}-${offset}`}
                  className="grid grid-cols-[90px_48px_1fr_42px] gap-2 rounded-lg border border-white/5 bg-black/30 px-3 py-3 text-white/50"
                >
                  <span
                    className={
                      status.startsWith('5')
                        ? 'text-red-300'
                        : status.startsWith('4')
                          ? 'text-white/35'
                          : 'text-indigo-300'
                    }
                  >
                    {status}
                  </span>
                  <span>{method}</span>
                  <span>{path}</span>
                  <span className="text-right text-white/25">{time}</span>
                </div>
              ))}
            </div>
          </article>
        </div>
      </MarketingContainer>
    </section>
  );
}
