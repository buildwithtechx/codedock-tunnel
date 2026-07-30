import { Check, Globe2 } from 'lucide-react';
import { MarketingContainer } from '#/components/layout';

export function CustomDomainSection() {
  return (
    <section className="py-16 sm:py-20">
      <MarketingContainer className="grid items-center gap-14 lg:grid-cols-2">
        <div>
          <h2 className="text-4xl font-semibold tracking-[-0.05em] sm:text-5xl">
            Memorable endpoints, every time.
          </h2>
          <p className="mt-6 max-w-xl text-lg leading-8 text-white/45">
            Use a memorable endpoint for demos, webhooks, and shared previews.
            Point your domain at the relay and keep routing in one place.
          </p>
          <div className="mt-8 space-y-4 text-sm text-white/60">
            {[
              'Reserved hostnames for repeatable previews',
              'Simple DNS target and certificate readiness',
              'Move traffic without changing your application',
            ].map((item) => (
              <div key={item} className="flex items-center gap-3">
                <span className="flex size-6 items-center justify-center rounded-full border border-indigo-300/30 bg-indigo-300/10">
                  <Check className="size-3 text-indigo-300" />
                </span>
                {item}
              </div>
            ))}
          </div>
        </div>
        <div className="relative">
          <div className="absolute inset-0 rounded-full bg-indigo-400/10 blur-3xl" />
          <div className="relative overflow-hidden rounded-3xl border border-white/10 bg-[#101421] p-5 shadow-2xl shadow-indigo-950/30">
            <div className="flex items-center gap-3 border-b border-white/10 pb-5">
              <span className="flex size-10 items-center justify-center rounded-xl bg-indigo-300/10">
                <Globe2 className="size-5 text-indigo-300" />
              </span>
              <div>
                <p className="font-medium">api.yourcompany.com</p>
                <span className="rounded-full bg-emerald-300/10 px-2 py-0.5 text-xs text-emerald-300">
                  Ready
                </span>
              </div>
            </div>
            <div className="mt-5 overflow-hidden rounded-2xl border border-white/10 bg-black/25 font-mono text-xs">
              {[
                ['Type', 'CNAME'],
                ['Name', 'api'],
                ['Value', 'edge.codedock-tunnel.dev'],
                ['TTL', 'Auto'],
              ].map(([key, value]) => (
                <div
                  key={key}
                  className="flex justify-between border-b border-white/5 px-4 py-4 last:border-0"
                >
                  <span className="text-white/35">{key}</span>
                  <span className="text-white/70">{value}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </MarketingContainer>
    </section>
  );
}
