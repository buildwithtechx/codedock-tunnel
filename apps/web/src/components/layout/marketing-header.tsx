import { Link } from '@tanstack/react-router';
import { ChevronDown, ExternalLink, LifeBuoy, Menu, X } from 'lucide-react';
import { useEffect, useState } from 'react';
import {
  SiExpress,
  SiGithub,
  SiNestjs,
  SiNextdotjs,
  SiVite,
} from 'react-icons/si';
import { MarketingContainer } from './marketing-container';

const mobileLinks = [
  { label: 'Pricing', to: '/pricing' as const },
  { label: 'Changelog', to: '/changelog' as const },
  { label: 'Plugins', to: '/plugins' as const },
  { label: 'Contact', to: '/contact' as const },
];

const pluginLinks = [
  { label: 'Vite', id: 'vite', icon: SiVite },
  { label: 'Next.js', id: 'next', icon: SiNextdotjs },
  { label: 'Express', id: 'express', icon: SiExpress },
  { label: 'NestJS', id: 'nest', icon: SiNestjs },
];

export function MarketingHeader() {
  const [open, setOpen] = useState(false);
  const [docsOpen, setDocsOpen] = useState(false);
  const [helpOpen, setHelpOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const updateScrolled = () => setScrolled(window.scrollY > 50);
    updateScrolled();
    window.addEventListener('scroll', updateScrolled, { passive: true });
    return () => window.removeEventListener('scroll', updateScrolled);
  }, []);

  return (
    <header
      className={`fixed inset-x-0 top-0 z-50 border-b transition-[background-color,border-color,box-shadow,backdrop-filter] duration-300 ${
        scrolled
          ? 'border-white/10 bg-black shadow-[0_10px_35px_rgba(0,0,0,0.32)] backdrop-blur-xl'
          : 'border-transparent bg-transparent'
      }`}
    >
      <MarketingContainer className="relative flex h-18 items-center justify-between">
        <Link
          to="/"
          className="flex items-center gap-3"
          onClick={() => setOpen(false)}
        >
          <img src="/favicon.svg" alt="" className="size-10 rounded-xl" />
          <span className="hidden font-semibold tracking-tight text-white sm:inline">
            Codedock Tunnel
          </span>
        </Link>
        <nav
          className="absolute left-1/2 hidden -translate-x-1/2 items-center gap-8 md:flex"
          aria-label="Main navigation"
        >
          <DropdownButton
            label="Docs"
            open={docsOpen}
            onOpen={() => {
              setDocsOpen(true);
              setHelpOpen(false);
            }}
            onClose={() => setDocsOpen(false)}
          >
            <div className="grid grid-cols-[150px_1fr] gap-6">
              <div className="flex flex-col gap-1 border-r border-white/10 pr-5">
                <Link
                  to="/docs/$"
                  params={{ _splat: '' }}
                  className="rounded-lg px-3 py-2 text-sm text-white/70 transition-colors hover:bg-indigo-300/10 hover:text-indigo-300"
                >
                  Getting started
                </Link>
                <Link
                  to="/docs/$"
                  params={{ _splat: 'cli' }}
                  className="rounded-lg px-3 py-2 text-sm text-white/70 transition-colors hover:bg-indigo-300/10 hover:text-indigo-300"
                >
                  CLI reference
                </Link>
                <Link
                  to="/plugins"
                  className="rounded-lg px-3 py-2 text-sm text-white/70 transition-colors hover:bg-indigo-300/10 hover:text-indigo-300"
                >
                  All integrations
                </Link>
              </div>
              <div className="grid grid-cols-2 gap-3">
                {pluginLinks.map(({ label, id, icon: Icon }) => (
                  <Link
                    key={label}
                    to="/plugins/$pluginId"
                    params={{ pluginId: id }}
                    className="flex size-16 items-center justify-center rounded-xl bg-white/[0.05] text-2xl text-white/60 transition-colors hover:bg-indigo-300/10 hover:text-indigo-300"
                    title={label}
                  >
                    <Icon />
                  </Link>
                ))}
              </div>
            </div>
          </DropdownButton>
          <Link
            to="/pricing"
            className="text-sm text-white/60 transition-colors hover:text-indigo-300"
          >
            Pricing
          </Link>
          <Link
            to="/changelog"
            className="text-sm text-white/60 transition-colors hover:text-indigo-300"
          >
            Changelog
          </Link>
          <DropdownButton
            label="Help"
            open={helpOpen}
            onOpen={() => {
              setHelpOpen(true);
              setDocsOpen(false);
            }}
            onClose={() => setHelpOpen(false)}
          >
            <div className="grid gap-1">
              <Link
                to="/contact"
                className="group flex items-center gap-3 rounded-lg px-3 py-3 text-sm text-white/70 transition-colors hover:bg-indigo-300/10 hover:text-indigo-300"
              >
                <LifeBuoy className="size-4 transition-colors group-hover:text-indigo-300" />
                <span>
                  <strong className="block font-medium text-white transition-colors group-hover:text-indigo-300">
                    Contact us
                  </strong>
                  <small className="text-white/40">
                    Questions and feedback
                  </small>
                </span>
              </Link>
              <a
                href="https://github.com/codedock"
                target="_blank"
                rel="noreferrer"
                className="group flex items-center gap-3 rounded-lg px-3 py-3 text-sm text-white/70 transition-colors hover:bg-indigo-300/10 hover:text-indigo-300"
              >
                <SiGithub className="size-4" />
                <span>
                  <strong className="block font-medium text-white transition-colors group-hover:text-indigo-300">
                    GitHub
                  </strong>
                  <small className="text-white/40">Read the source</small>
                </span>
                <ExternalLink className="ml-auto size-3 text-white/30" />
              </a>
            </div>
          </DropdownButton>
        </nav>
        <div className="hidden items-center gap-3 md:flex">
          <a
            href="https://github.com/codedock"
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/[0.04] px-3 py-2 text-xs text-white/60 transition-colors hover:border-indigo-300/40 hover:text-indigo-300"
          >
            <SiGithub className="size-4" /> Star on GitHub
          </a>
          <Link
            to="/login"
            className="rounded-full border border-white/10 px-4 py-2.5 text-sm font-medium text-white/70 transition-colors hover:border-indigo-300/40 hover:bg-indigo-300/10 hover:text-indigo-300"
          >
            Sign in
          </Link>
          <Link
            to="/signup"
            className="rounded-full bg-white px-5 py-2.5 text-sm font-semibold text-[#080b14] shadow-[0_0_24px_rgba(255,255,255,0.12)] transition-transform hover:-translate-y-0.5"
          >
            Start building
          </Link>
        </div>
        <button
          type="button"
          aria-label={open ? 'Close navigation' : 'Open navigation'}
          aria-expanded={open}
          onClick={() => setOpen((value) => !value)}
          className="rounded-lg border border-white/15 p-2 text-white md:hidden"
        >
          {open ? <X className="size-5" /> : <Menu className="size-5" />}
        </button>
      </MarketingContainer>
      {open && (
        <div className="border-t border-indigo-200/10 bg-black px-6 py-5 md:hidden">
          <nav
            className="mx-auto flex max-w-7xl flex-col gap-1"
            aria-label="Mobile navigation"
          >
            <Link
              to="/docs/$"
              params={{ _splat: '' }}
              onClick={() => setOpen(false)}
              className="rounded-lg px-3 py-3 text-sm text-white/70 transition-colors hover:bg-indigo-300/10 hover:text-indigo-300"
            >
              Documentation
            </Link>
            {mobileLinks.map((link) => (
              <Link
                key={link.label}
                to={link.to}
                onClick={() => setOpen(false)}
                className="rounded-lg px-3 py-3 text-sm text-white/70 transition-colors hover:bg-indigo-300/10 hover:text-indigo-300"
              >
                {link.label}
              </Link>
            ))}
            <div className="mt-3 grid grid-cols-2 gap-2 border-t border-white/10 pt-4">
              <Link
                to="/login"
                onClick={() => setOpen(false)}
                className="rounded-full border border-white/15 px-4 py-3 text-center text-sm text-white"
              >
                Sign in
              </Link>
              <Link
                to="/signup"
                onClick={() => setOpen(false)}
                className="rounded-full bg-white px-4 py-3 text-center text-sm font-semibold text-[#080b14]"
              >
                Start building
              </Link>
            </div>
          </nav>
        </div>
      )}
    </header>
  );
}

function DropdownButton({
  label,
  open,
  onOpen,
  onClose,
  children,
}: {
  label: string;
  open: boolean;
  onOpen: () => void;
  onClose: () => void;
  children: React.ReactNode;
}) {
  return (
    <nav
      className="relative"
      onMouseEnter={onOpen}
      onMouseLeave={onClose}
      onBlur={(event) => {
        if (!event.currentTarget.contains(event.relatedTarget)) onClose();
      }}
    >
      <button
        type="button"
        onClick={onOpen}
        onFocus={onOpen}
        className="inline-flex items-center gap-1 text-sm text-white/60 transition-colors hover:text-indigo-300"
      >
        {label}
        <ChevronDown
          className={`size-3.5 transition-transform ${open ? 'rotate-180' : ''}`}
        />
      </button>
      {open && (
        <div className="absolute left-1/2 top-full w-[420px] -translate-x-1/2 pt-5">
          <div className="rounded-2xl border border-white/10 bg-[#0a0a0a] p-5 shadow-2xl shadow-black/60">
            {children}
          </div>
        </div>
      )}
    </nav>
  );
}
