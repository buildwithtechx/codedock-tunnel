import { ArrowRight, Database, Gamepad2, Globe2 } from 'lucide-react';
import { MarketingContainer } from '#/components/layout';

const protocols = [
  {
    title: 'HTTP / HTTPS',
    text: 'Secure URLs for local web servers, APIs, webhooks, and OAuth callbacks.',
    command: 'codedock-tunnel http 3000',
    icon: Globe2,
    cardClass: 'hover:border-blue-300/30',
    iconClass: 'border-blue-300/20 bg-blue-300/10 text-blue-300',
  },
  {
    title: 'TCP tunnels',
    text: 'Expose databases, SSH, RDP, and other TCP services securely.',
    command: 'codedock-tunnel tcp 5432',
    icon: Database,
    cardClass: 'hover:border-violet-300/30',
    iconClass: 'border-violet-300/20 bg-violet-300/10 text-violet-300',
  },
  {
    title: 'UDP tunnels',
    text: 'Carry real-time services, game servers, and other UDP traffic.',
    command: 'codedock-tunnel udp 25565',
    icon: Gamepad2,
    cardClass: 'hover:border-amber-300/30',
    iconClass: 'border-amber-300/20 bg-amber-300/10 text-amber-300',
  },
];

export function ProtocolsSection() {
  return (
    <section className="py-16 sm:py-20">
      <MarketingContainer>
        <div className="mx-auto max-w-3xl text-center">
          <h2 className="text-4xl font-semibold tracking-[-0.05em] sm:text-5xl">
            Every protocol. One workflow.
          </h2>
          <p className="mt-6 text-lg leading-8 text-white/50">
            Codedock Tunnel is not limited to browser traffic. Use one CLI and
            one account for HTTP, HTTPS, TCP, and UDP services.
          </p>
        </div>
        <div className="mt-12 grid gap-5 md:grid-cols-3">
          {protocols.map(
            ({ title, text, command, icon: Icon, cardClass, iconClass }) => (
              <article
                key={title}
                className={`group relative overflow-hidden rounded-3xl border border-white/10 bg-white/[0.025] p-8 transition-transform hover:-translate-y-1 ${cardClass}`}
              >
                <div
                  className={`mb-7 flex size-12 items-center justify-center rounded-2xl border ${iconClass}`}
                >
                  <Icon className="size-6" />
                </div>
                <h3 className="text-xl font-semibold">{title}</h3>
                <p className="mt-3 min-h-16 text-sm leading-6 text-white/40">
                  {text}
                </p>
                <div className="mt-8 flex items-center justify-between rounded-xl border border-white/10 bg-black/30 p-3 font-mono text-xs text-white/55">
                  <span>{command}</span>
                  <ArrowRight
                    className={`size-4 opacity-0 transition-all group-hover:translate-x-1 group-hover:opacity-100 ${iconClass.split(' ').at(-1)}`}
                  />
                </div>
              </article>
            ),
          )}
        </div>
      </MarketingContainer>
    </section>
  );
}
