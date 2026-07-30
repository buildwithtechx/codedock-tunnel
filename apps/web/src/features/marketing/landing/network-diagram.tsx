import { Globe2, Laptop, Server, ShieldCheck, Sparkles } from 'lucide-react';
import { motion } from 'motion/react';
import { MarketingContainer } from '#/components/layout';

const nodes = [
  {
    label: 'Local service',
    detail: 'localhost:3000',
    icon: Laptop,
    color: 'text-amber-300',
  },
  {
    label: 'Encrypted session',
    detail: 'persistent connection',
    icon: ShieldCheck,
    color: 'text-emerald-300',
  },
  {
    label: 'Codedock relay',
    detail: 'routes your tunnel',
    icon: Server,
    color: 'text-violet-300',
  },
  {
    label: 'Public endpoint',
    detail: '{subdomain}.codedock-tunnel.dev',
    icon: Globe2,
    color: 'text-cyan-300',
  },
];

export function NetworkDiagram() {
  return (
    <section className="border-b border-white/10 py-28 sm:py-36">
      <MarketingContainer>
        <div className="text-center">
          <div className="relative inline-block">
            <h2 className="relative text-6xl font-semibold italic tracking-[-0.07em] text-transparent bg-gradient-to-b from-white via-white/90 to-white/40 bg-clip-text sm:text-8xl">
              It just works
            </h2>
            <Sparkles className="absolute -right-10 -top-8 size-9 text-amber-200" />
          </div>
        </div>
        <div className="relative mt-24">
          <div className="absolute left-0 right-0 top-8 hidden h-px bg-gradient-to-r from-transparent via-white/25 to-transparent md:block" />
          <motion.div
            className="absolute left-0 top-8 hidden size-1 -translate-y-1/2 rounded-full bg-white shadow-[0_0_15px_5px_rgba(255,255,255,0.7)] md:block"
            animate={{ left: ['0%', '100%'] }}
            transition={{ duration: 2.5, repeat: Infinity, ease: 'linear' }}
          />
          <div className="relative grid gap-10 md:grid-cols-4 md:gap-6">
            {nodes.map(({ label, detail, icon: Icon, color }) => (
              <div
                key={label}
                className="flex flex-col items-center text-center"
              >
                <div
                  className={`relative z-10 flex size-16 items-center justify-center rounded-2xl border border-white/10 bg-[#0c0f18] ${color}`}
                >
                  <Icon className="size-7" />
                </div>
                <h3 className="mt-6 font-semibold">{label}</h3>
                <p className="mt-2 text-sm text-white/40">{detail}</p>
              </div>
            ))}
          </div>
        </div>
      </MarketingContainer>
    </section>
  );
}
