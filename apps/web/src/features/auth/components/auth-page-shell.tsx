import { Link } from '@tanstack/react-router';
import type { ReactNode } from 'react';

type AuthPageShellProps = {
  title: string;
  description: string;
  footer: ReactNode;
  children: ReactNode;
};

export function AuthPageShell({
  title,
  description,
  footer,
  children,
}: AuthPageShellProps) {
  return (
    <main className="relative flex min-h-screen items-center justify-center overflow-hidden bg-[#080b14] px-6 py-12 text-white">
      <div className="absolute inset-x-0 top-0 h-80 bg-[radial-gradient(ellipse_at_top,rgba(99,102,241,0.18),transparent_70%)]" />
      <Link
        to="/"
        className="absolute left-6 top-6 z-10 flex items-center gap-3 sm:left-10 sm:top-9"
      >
        <img src="/favicon.svg" alt="" className="size-9 rounded-xl" />
        <span className="font-semibold tracking-tight">Codedock Tunnel</span>
      </Link>
      <section className="relative z-10 w-full max-w-md rounded-3xl border border-white/10 bg-black/35 p-7 shadow-2xl shadow-black/40 backdrop-blur sm:p-9">
        <p className="text-sm font-medium text-indigo-300">Codedock Tunnel</p>
        <h1 className="mt-4 text-3xl font-semibold tracking-[-0.04em]">
          {title}
        </h1>
        <p className="mt-3 leading-7 text-white/50">{description}</p>
        <div className="mt-8">{children}</div>
        <div className="mt-7 border-t border-white/10 pt-6 text-center text-sm text-white/45">
          {footer}
        </div>
      </section>
    </main>
  );
}
