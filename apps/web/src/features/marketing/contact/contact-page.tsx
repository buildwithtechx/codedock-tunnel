import { Mail, MessageSquare, ShieldCheck } from 'lucide-react';
import { useState } from 'react';
import { MarketingContainer } from '#/components/layout';

export function ContactPage() {
  const [sent, setSent] = useState(false);
  return (
    <section className="pb-16 pt-28 sm:pt-32">
      <MarketingContainer className="grid max-w-6xl gap-14 lg:grid-cols-[0.75fr_1.25fr]">
        <div>
          <p className="text-sm font-medium text-cyan-300">Talk to us</p>
          <h1 className="mt-4 text-4xl font-semibold tracking-tight sm:text-5xl">
            Questions or feedback?
          </h1>
          <p className="mt-5 leading-8 text-white/50">
            Tell us what you are building. We read every message and use your
            feedback to shape the product.
          </p>
          <div className="mt-10 space-y-5 text-sm text-white/60">
            <p className="flex items-center gap-3">
              <Mail className="size-4 text-cyan-300" />
              hello@codedock-tunnel.dev
            </p>
            <p className="flex items-center gap-3">
              <MessageSquare className="size-4 text-cyan-300" />
              Product and integration questions welcome
            </p>
            <p className="flex items-center gap-3">
              <ShieldCheck className="size-4 text-cyan-300" />
              Never send credentials or private tunnel data
            </p>
          </div>
        </div>
        <form
          onSubmit={(event) => {
            event.preventDefault();
            setSent(true);
          }}
          className="rounded-2xl border border-white/10 bg-white/[0.03] p-6 sm:p-8"
        >
          <div className="grid gap-5 sm:grid-cols-2">
            <label className="grid gap-2 text-sm text-white/65">
              Name
              <input
                required
                name="name"
                className="rounded-lg border border-white/10 bg-black/20 px-3 py-3 text-white outline-none ring-indigo-300 focus:ring-2"
              />
            </label>
            <label className="grid gap-2 text-sm text-white/65">
              Email
              <input
                required
                type="email"
                name="email"
                className="rounded-lg border border-white/10 bg-black/20 px-3 py-3 text-white outline-none ring-indigo-300 focus:ring-2"
              />
            </label>
          </div>
          <label className="mt-5 grid gap-2 text-sm text-white/65">
            Subject
            <input
              required
              name="subject"
              className="rounded-lg border border-white/10 bg-black/20 px-3 py-3 text-white outline-none ring-indigo-300 focus:ring-2"
            />
          </label>
          <label className="mt-5 grid gap-2 text-sm text-white/65">
            Message
            <textarea
              required
              name="message"
              rows={6}
              className="resize-y rounded-lg border border-white/10 bg-black/20 px-3 py-3 text-white outline-none ring-indigo-300 focus:ring-2"
            />
          </label>
          <button
            type="submit"
            className="mt-6 rounded-lg bg-white px-5 py-3 text-sm font-semibold text-[#080b14]"
          >
            {sent ? 'Message ready to send' : 'Send message'}
          </button>
          {sent && (
            <p className="mt-3 text-sm text-emerald-300">
              Thanks. Connect your mail endpoint to deliver this form.
            </p>
          )}
        </form>
      </MarketingContainer>
    </section>
  );
}
