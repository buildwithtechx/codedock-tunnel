import { Link } from '@tanstack/react-router';
import { ArrowRight, Check, Copy, Terminal } from 'lucide-react';
import { useState } from 'react';
import { MarketingContainer } from '#/components/layout';
import type { PluginDefinition } from './plugin-data';

export function PluginDetailPage({ plugin }: { plugin: PluginDefinition }) {
  const [copied, setCopied] = useState(false);
  async function copyInstall() {
    await navigator.clipboard.writeText(plugin.install);
    setCopied(true);
    setTimeout(() => setCopied(false), 1800);
  }
  return (
    <>
      <section className="relative overflow-hidden pt-36 pb-24 sm:pt-44">
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_10%,rgba(99,102,241,0.18),transparent_32%)]" />
        <MarketingContainer className="relative text-center">
          <div
            className={`mx-auto flex size-20 items-center justify-center rounded-3xl border border-white/10 bg-white/[0.04] text-4xl ${plugin.colorClass}`}
          >
            <plugin.icon />
          </div>
          <p className="mt-8 text-sm font-medium text-cyan-300">
            {plugin.eyebrow}
          </p>
          <h1 className="mx-auto mt-4 max-w-4xl text-4xl font-semibold tracking-[-0.06em] sm:text-7xl">
            {plugin.headline}
          </h1>
          <p className="mx-auto mt-6 max-w-2xl text-lg leading-8 text-white/50">
            {plugin.description}
          </p>
          <div className="mt-9 flex flex-col items-center justify-center gap-3 sm:flex-row">
            <Link
              to="/signup"
              className="inline-flex items-center gap-2 rounded-full bg-white px-6 py-3.5 text-sm font-semibold text-[#080b14]"
            >
              Start a tunnel <ArrowRight className="size-4" />
            </Link>
            <Link
              to="/docs/$"
              params={{ _splat: plugin.docsSlug }}
              className="rounded-full border border-white/15 px-6 py-3.5 text-sm text-white/70 hover:bg-white/5"
            >
              Read the guide
            </Link>
          </div>
          <button
            type="button"
            onClick={copyInstall}
            className="group mt-7 inline-flex items-center gap-3 rounded-full border border-white/10 bg-white/[0.04] px-5 py-3 font-mono text-xs text-white/60 hover:border-indigo-300/35"
          >
            <span className="text-indigo-300">$</span>
            {plugin.install}
            {copied ? (
              <Check className="size-4 text-emerald-300" />
            ) : (
              <Copy className="size-4 opacity-0 group-hover:opacity-100" />
            )}
          </button>
        </MarketingContainer>
      </section>
      <section className="py-24">
        <MarketingContainer className="grid items-start gap-14 lg:grid-cols-[0.8fr_1.2fr]">
          <div>
            <h2 className="text-3xl font-semibold tracking-tight sm:text-5xl">
              A focused integration, not another runtime.
            </h2>
            <p className="mt-5 leading-7 text-white/45">
              The adapter owns the framework-specific details. Codedock SDK owns
              the protocol, session, reconnection, and tunnel lifecycle
              underneath.
            </p>
            <ul className="mt-8 space-y-3">
              {plugin.features.map((feature) => (
                <li
                  key={feature}
                  className="flex items-center gap-3 text-sm text-white/65"
                >
                  <Check className="size-4 text-cyan-300" />
                  {feature}
                </li>
              ))}
            </ul>
          </div>
          <CodePanel plugin={plugin} />
        </MarketingContainer>
      </section>
      <section className="bg-white/[0.02] py-24">
        <MarketingContainer>
          <div className="mx-auto max-w-2xl text-center">
            <p className="text-sm font-medium text-cyan-300">Fits your stack</p>
            <h2 className="mt-4 text-3xl font-semibold tracking-tight sm:text-5xl">
              Keep the tools you already trust.
            </h2>
          </div>
          <div className="mx-auto mt-12 grid max-w-4xl grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6">
            {plugin.technologies.map(({ label, icon: Icon }) => (
              <div
                key={label}
                className="flex flex-col items-center gap-3 rounded-2xl border border-white/10 bg-[#0d1220] px-4 py-5 text-center"
              >
                <Icon className="size-7 text-white/75" />
                <span className="text-xs text-white/50">{label}</span>
              </div>
            ))}
          </div>
        </MarketingContainer>
      </section>
      <section className="py-24">
        <MarketingContainer>
          <div className="mx-auto max-w-2xl text-center">
            <p className="text-sm font-medium text-indigo-300">
              Good fits for {plugin.name}
            </p>
            <h2 className="mt-4 text-3xl font-semibold tracking-tight sm:text-5xl">
              Build the workflow, not the tunnel plumbing.
            </h2>
          </div>
          <div className="mx-auto mt-12 grid max-w-5xl gap-4 md:grid-cols-3">
            {plugin.useCases.map((useCase) => (
              <div
                key={useCase}
                className="rounded-2xl border border-white/10 bg-white/[0.03] p-6"
              >
                <Terminal className="size-5 text-cyan-300" />
                <h3 className="mt-6 font-semibold">{useCase}</h3>
                <p className="mt-2 text-sm leading-6 text-white/40">
                  Keep your app local while the rest of the team, provider, or
                  browser reaches it through a controlled endpoint.
                </p>
              </div>
            ))}
          </div>
          <div className="mt-14 text-center">
            <Link
              to="/plugins"
              className="inline-flex items-center gap-2 text-sm text-white/60 hover:text-white"
            >
              View every integration <ArrowRight className="size-4" />
            </Link>
          </div>
        </MarketingContainer>
      </section>
    </>
  );
}

function CodePanel({ plugin }: { plugin: PluginDefinition }) {
  return (
    <div className="overflow-hidden rounded-2xl border border-white/10 bg-[#0b0d12] shadow-2xl">
      <div className="flex items-center gap-2 border-b border-white/10 bg-white/[0.03] px-4 py-3">
        <span className="size-2.5 rounded-full bg-red-400/80" />
        <span className="size-2.5 rounded-full bg-amber-300/80" />
        <span className="size-2.5 rounded-full bg-emerald-400/80" />
        <span className="ml-3 font-mono text-xs text-white/35">
          {plugin.fileName}
        </span>
      </div>
      <pre className="overflow-x-auto p-6 font-mono text-xs leading-7 text-white/70 sm:text-sm">
        <code>{plugin.code}</code>
      </pre>
    </div>
  );
}
