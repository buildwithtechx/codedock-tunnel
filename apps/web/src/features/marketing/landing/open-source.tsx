import { Heart, LockKeyhole, Users, Workflow } from 'lucide-react';
import { MarketingContainer } from '#/components/layout';

const values = [
  {
    title: 'Auditable security',
    text: 'Inspect the protocol, clients, and relay behavior in the open.',
    icon: LockKeyhole,
    iconClass: 'bg-blue-300/10 text-blue-300',
  },
  {
    title: 'Community driven',
    text: 'Build the developer surfaces you need and contribute upstream.',
    icon: Users,
    iconClass: 'bg-violet-300/10 text-violet-300',
  },
  {
    title: 'One protocol',
    text: 'Keep clients and integrations aligned around a public contract.',
    icon: Workflow,
    iconClass: 'bg-emerald-300/10 text-emerald-300',
  },
  {
    title: 'Built to last',
    text: 'Use a product foundation designed for real development workflows.',
    icon: Heart,
    iconClass: 'bg-rose-300/10 text-rose-300',
  },
];

export function OpenSourceSection() {
  return (
    <section className="py-24 sm:py-32">
      <MarketingContainer className="grid items-center gap-14 lg:grid-cols-2">
        <div>
          <h2 className="text-4xl font-semibold tracking-[-0.05em] sm:text-6xl">
            Open source
            <br />
            by design
          </h2>
          <p className="mt-6 max-w-xl text-lg leading-8 text-white/50">
            Codedock Tunnel is built in the open so developers can understand
            the connection path, audit the clients, and help shape what comes
            next.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <a
              href="https://github.com/codedock"
              target="_blank"
              rel="noreferrer"
              className="rounded-full bg-white px-5 py-3 text-sm font-semibold text-[#080b14]"
            >
              View on GitHub
            </a>
            <a
              href="mailto:hello@codedock-tunnel.dev"
              className="rounded-full border border-white/10 px-5 py-3 text-sm font-medium text-white/70 hover:bg-white/5"
            >
              Join the conversation
            </a>
          </div>
        </div>
        <div className="rounded-3xl border border-white/10 bg-[#0c1019] p-8">
          <div className="grid gap-8 sm:grid-cols-2">
            {values.map(({ title, text, icon: Icon, iconClass }) => (
              <div key={title}>
                <span
                  className={`mb-4 flex size-10 items-center justify-center rounded-xl ${iconClass}`}
                >
                  <Icon className="size-5" />
                </span>
                <h3 className="font-semibold">{title}</h3>
                <p className="mt-2 text-sm leading-6 text-white/40">{text}</p>
              </div>
            ))}
          </div>
        </div>
      </MarketingContainer>
    </section>
  );
}
