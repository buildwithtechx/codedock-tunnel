import { Link } from '@tanstack/react-router';
import { ArrowRight, Check, Copy } from 'lucide-react';
import { AnimatePresence, motion } from 'motion/react';
import { useEffect, useState } from 'react';
import { MarketingContainer } from '#/components/layout';

const logs = [
  ['GET', '/api/health', '200'],
  ['POST', '/webhooks/stripe', '201'],
  ['GET', '/oauth/callback', '302'],
  ['GET', '/api/projects', '200'],
];

export function Hero() {
  const [copied, setCopied] = useState(false);
  const [hovered, setHovered] = useState(false);
  const [visibleLogs, setVisibleLogs] = useState(
    logs.slice(0, 3).map((log, id) => ({ log, id })),
  );

  useEffect(() => {
    if (!hovered) return;
    let index = 3;
    const timer = setInterval(() => {
      setVisibleLogs((current) => [
        ...current.slice(1),
        { log: logs[index++ % logs.length], id: index },
      ]);
    }, 800);
    return () => clearInterval(timer);
  }, [hovered]);

  async function copyCommand() {
    await navigator.clipboard.writeText(
      'curl -fsSL https://codedock-tunnel.dev/install.sh | bash',
    );
    setCopied(true);
    setTimeout(() => setCopied(false), 1800);
  }

  return (
    <section className="relative min-h-[880px] overflow-hidden border-b border-white/10 pt-28 sm:pt-32">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_22%,rgba(79,70,229,0.16),transparent_28%),radial-gradient(circle_at_12%_42%,rgba(8,145,178,0.12),transparent_24%)]" />
      <div className="absolute inset-x-[-20%] top-[31%] h-[390px] -rotate-[16deg] bg-[linear-gradient(105deg,transparent_5%,rgba(129,140,248,0.18)_32%,rgba(255,255,255,0.08)_53%,transparent_80%)] blur-[1px]" />
      <div className="absolute inset-0 opacity-40 [background-image:radial-gradient(circle_at_20%_40%,white_0_1px,transparent_1px),radial-gradient(circle_at_70%_25%,white_0_1px,transparent_1px),radial-gradient(circle_at_42%_68%,white_0_1px,transparent_1px)] [background-size:190px_170px,240px_220px,320px_260px]" />
      <MarketingContainer className="relative flex flex-col items-center">
        <div className="mt-10 inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/[0.04] px-3 py-1.5 text-xs text-white/55">
          <span className="size-1.5 rounded-full bg-cyan-300" />
          Open-source tunnel infrastructure for developers
        </div>
        <h1 className="mt-8 max-w-5xl text-center text-5xl font-semibold leading-[0.98] tracking-[-0.06em] sm:text-7xl lg:text-8xl">
          <motion.button
            type="button"
            className="relative inline-block cursor-default appearance-none border-0 bg-transparent p-0 text-inherit"
            onMouseEnter={() => setHovered(true)}
            onMouseLeave={() => setHovered(false)}
          >
            <motion.span
              animate={{ rotate: hovered ? -5 : 0, y: hovered ? -4 : 0 }}
              transition={{ type: 'spring', stiffness: 300, damping: 20 }}
              className="relative z-10 inline-block rounded-2xl border border-indigo-300/35 bg-indigo-300/15 px-4 py-1"
            >
              Expose
            </motion.span>
            <span className="pointer-events-none absolute inset-0 rounded-2xl border border-cyan-300/30 bg-[#090d17] px-4 py-1 font-mono text-xs leading-[5rem] sm:text-sm">
              {' '}
              <AnimatePresence mode="popLayout">
                {visibleLogs.map(({ log, id }) => (
                  <motion.span
                    key={id}
                    initial={{ opacity: 0, x: -10 }}
                    animate={{ opacity: hovered ? 1 : 0, x: 0 }}
                    exit={{ opacity: 0 }}
                    className="mr-4 inline-flex gap-2"
                  >
                    <span className="text-cyan-300">{log[0]}</span>
                    <span className="text-white/50">{log[1]}</span>
                    <span className="text-emerald-300">{log[2]}</span>
                  </motion.span>
                ))}
              </AnimatePresence>
            </span>
          </motion.button>{' '}
          your local service
          <br className="hidden sm:block" /> to the internet
        </h1>
        <p className="mt-8 max-w-2xl text-center text-lg leading-8 text-white/55 sm:text-xl">
          Codedock Tunnel gives your local apps a secure, observable public
          endpoint for previews, webhooks, OAuth callbacks, and CI workflows.
        </p>
        <div className="mt-9 flex w-full flex-col items-center justify-center gap-3 sm:flex-row">
          <Link
            to="/signup"
            className="group inline-flex w-full items-center justify-center gap-2 rounded-full bg-white px-7 py-4 text-base font-semibold text-[#080b14] transition-transform hover:-translate-y-0.5 sm:w-auto"
          >
            Get started free{' '}
            <ArrowRight className="size-5 transition-transform group-hover:translate-x-1" />
          </Link>
          <button
            type="button"
            onClick={copyCommand}
            className="group inline-flex w-full items-center justify-center gap-3 rounded-full border border-white/10 bg-white/[0.04] px-6 py-4 font-mono text-xs text-white/60 transition-colors hover:border-indigo-300/40 hover:bg-white/[0.08] sm:w-auto"
          >
            <span className="text-indigo-300">$</span> curl .../install.sh{' '}
            {copied ? (
              <Check className="size-4 text-emerald-300" />
            ) : (
              <Copy className="size-4 opacity-0 transition-opacity group-hover:opacity-100" />
            )}
          </button>
        </div>
        <TerminalWindow />
      </MarketingContainer>
    </section>
  );
}

function TerminalWindow() {
  return (
    <div className="mt-16 w-full max-w-4xl overflow-hidden rounded-2xl border border-white/10 bg-[#0b0d12] text-left font-mono text-sm shadow-2xl shadow-indigo-950/40">
      <div className="flex items-center gap-2 border-b border-white/10 bg-white/[0.03] px-5 py-3">
        <span className="size-3 rounded-full bg-red-400/80" />
        <span className="size-3 rounded-full bg-amber-300/80" />
        <span className="size-3 rounded-full bg-emerald-400/80" />
        <span className="ml-auto mr-auto text-xs text-white/30">
          user@codedock-cli
        </span>
      </div>
      <div className="grid gap-2 p-6 text-xs leading-6 sm:text-sm">
        <p>
          <span className="text-emerald-300">➜</span>{' '}
          <span className="text-cyan-300">~</span> codedock-tunnel 3000
        </p>
        <p className="text-cyan-300">Connecting to Codedock Tunnel...</p>
        <p className="text-emerald-300">Linked to local port 3000</p>
        <p className="text-fuchsia-300">
          Tunnel ready: https://quiet-moon.codedock-tunnel.dev
        </p>
        <p className="text-amber-300">
          Keep this process running to keep the tunnel active.
        </p>
        <div className="mt-3 grid gap-1 border-t border-white/10 pt-4 text-white/45">
          {[
            'GET     /api/health       200   12ms',
            'POST    /webhooks        201   45ms',
            'GET     /oauth/callback  302    8ms',
            'GET     /favicon.ico     200    2ms',
          ].map((line) => (
            <p key={line}>{line}</p>
          ))}
        </div>
      </div>
    </div>
  );
}
