import { Link } from '@tanstack/react-router';
import { Check, Sparkles } from 'lucide-react';
import { useState } from 'react';
import { MarketingContainer } from '#/components/layout';

const plans = [
  {
    name: 'Free',
    price: 0,
    description: 'For trying the workflow',
    features: [
      '2 active tunnels',
      'Shared public hostnames',
      '3-day request retention',
      'Community support',
    ],
  },
  {
    name: 'Ray',
    price: 12,
    description: 'For independent builders',
    features: [
      '10 active tunnels',
      'Reserved subdomains',
      '14-day request retention',
      'Password protection',
    ],
  },
  {
    name: 'Beam',
    price: 29,
    description: 'For teams shipping previews',
    featured: true,
    features: [
      '50 active tunnels',
      'Custom domains',
      '30-day request retention',
      'Team workspaces',
    ],
  },
  {
    name: 'Pulse',
    price: 79,
    description: 'For high-volume workflows',
    features: [
      'Unlimited tunnel projects',
      'Priority relay capacity',
      '90-day request retention',
      'Advanced usage controls',
    ],
  },
];

export function PricingPage() {
  const [yearly, setYearly] = useState(false);
  return (
    <section className="border-b border-white/10 pt-36 pb-24 sm:pt-44">
      <MarketingContainer>
        <div className="mx-auto max-w-2xl text-center">
          <p className="text-sm font-medium text-cyan-300">
            Pricing that scales with your tunnels
          </p>
          <h1 className="mt-4 text-4xl font-semibold tracking-tight sm:text-6xl">
            Keep shipping. Pay for the traffic you need.
          </h1>
          <p className="mt-5 text-lg leading-8 text-white/50">
            Start with the essentials and add capacity when your previews,
            webhooks, and teams grow.
          </p>
          <div className="mt-8 inline-flex rounded-xl border border-white/10 bg-white/[0.04] p-1">
            <button
              type="button"
              onClick={() => setYearly(false)}
              className={`rounded-lg px-4 py-2 text-sm ${!yearly ? 'bg-white text-[#080b14]' : 'text-white/55'}`}
            >
              Monthly
            </button>
            <button
              type="button"
              onClick={() => setYearly(true)}
              className={`rounded-lg px-4 py-2 text-sm ${yearly ? 'bg-white text-[#080b14]' : 'text-white/55'}`}
            >
              Yearly <span className="ml-1 text-cyan-500">save 2 months</span>
            </button>
          </div>
        </div>
        <div className="mt-16 grid gap-4 lg:grid-cols-4">
          {plans.map((plan) => (
            <article
              key={plan.name}
              className={`relative flex flex-col rounded-2xl border p-6 ${plan.featured ? 'border-indigo-300/50 bg-indigo-300/10 shadow-xl shadow-indigo-950/30' : 'border-white/10 bg-white/[0.03]'}`}
            >
              {plan.featured && (
                <div className="absolute -top-3 left-6 inline-flex items-center gap-1 rounded-full bg-indigo-300 px-3 py-1 text-xs font-semibold text-[#080b14]">
                  <Sparkles className="size-3" /> Recommended
                </div>
              )}
              <h2 className="text-xl font-semibold">{plan.name}</h2>
              <p className="mt-2 min-h-12 text-sm leading-5 text-white/45">
                {plan.description}
              </p>
              <div className="mt-6 flex items-baseline gap-1">
                <span className="text-4xl font-semibold">
                  {plan.price === 0
                    ? 'Free'
                    : `$${yearly ? Math.round(plan.price * 10) : plan.price}`}
                </span>
                {plan.price > 0 && (
                  <span className="text-sm text-white/40">
                    /{yearly ? 'year' : 'month'}
                  </span>
                )}
              </div>
              <Link
                to="/signup"
                className={`mt-6 rounded-lg px-4 py-3 text-center text-sm font-semibold ${plan.featured ? 'bg-white text-[#080b14]' : 'border border-white/15 text-white hover:bg-white/5'}`}
              >
                {plan.price === 0 ? 'Start free' : 'Choose plan'}
              </Link>
              <ul className="mt-7 space-y-3 border-t border-white/10 pt-6">
                {plan.features.map((feature) => (
                  <li
                    key={feature}
                    className="flex gap-2 text-sm text-white/60"
                  >
                    <Check className="mt-0.5 size-4 shrink-0 text-cyan-300" />
                    {feature}
                  </li>
                ))}
              </ul>
            </article>
          ))}
        </div>
      </MarketingContainer>
    </section>
  );
}
