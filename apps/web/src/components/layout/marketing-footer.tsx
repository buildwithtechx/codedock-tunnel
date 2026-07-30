import { Link } from '@tanstack/react-router';
import { GitBranch, Mail } from 'lucide-react';
import { MarketingContainer } from './marketing-container';

export function MarketingFooter() {
  return (
    <footer className="border-t border-white/10 bg-[#070910] py-14 text-white">
      <MarketingContainer>
        <div className="grid gap-10 md:grid-cols-6">
          <div className="md:col-span-2">
            <Link to="/" className="flex items-center gap-3">
              <img src="/favicon.svg" alt="" className="size-9 rounded-xl" />
              <span className="font-semibold">Codedock Tunnel</span>
            </Link>
            <p className="mt-4 max-w-xs text-sm leading-6 text-white/45">
              Secure public access for local services, previews, webhooks, and
              private networks.
            </p>
            <p className="mt-7 text-xs text-white/35">
              © {new Date().getFullYear()} Codedock Tunnel. All rights reserved.
            </p>
          </div>
          <FooterGroup title="Product">
            <FooterLink to="/pricing">Pricing</FooterLink>
            <FooterLink to="/changelog">Changelog</FooterLink>
            <FooterLink to="/contact">Contact</FooterLink>
          </FooterGroup>
          <FooterGroup title="Developers">
            <Link
              to="/docs/$"
              params={{ _splat: '' }}
              className="text-sm text-white/45 transition-colors hover:text-white"
            >
              Documentation
            </Link>
            <FooterLink to="/plugins">Plugins</FooterLink>
            <Link
              to="/docs/$"
              params={{ _splat: 'cli' }}
              className="text-sm text-white/45 transition-colors hover:text-white"
            >
              CLI reference
            </Link>
          </FooterGroup>
          <FooterGroup title="Integrations">
            <FooterLink to="/plugins">React</FooterLink>
            <FooterLink to="/plugins">Vite</FooterLink>
            <FooterLink to="/plugins">Next.js</FooterLink>
          </FooterGroup>
          <FooterGroup title="Legal">
            <FooterLink to="/terms">Terms</FooterLink>
            <FooterLink to="/privacy">Privacy</FooterLink>
            <a
              href="mailto:hello@codedock-tunnel.dev"
              className="text-sm text-white/45 transition-colors hover:text-white"
            >
              Email us
            </a>
          </FooterGroup>
        </div>
        <div className="mt-12 flex items-center gap-4 border-t border-white/10 pt-6 text-white/45">
          <a
            href="https://github.com/codedock"
            target="_blank"
            rel="noreferrer"
            aria-label="Codedock on GitHub"
            className="transition-colors hover:text-white"
          >
            <GitBranch className="size-4" />
          </a>
          <a
            href="mailto:hello@codedock-tunnel.dev"
            aria-label="Email Codedock Tunnel"
            className="transition-colors hover:text-white"
          >
            <Mail className="size-4" />
          </a>
        </div>
      </MarketingContainer>
    </footer>
  );
}

function FooterGroup({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <div>
      <h2 className="mb-4 text-sm font-semibold text-white">{title}</h2>
      <div className="flex flex-col items-start gap-3">{children}</div>
    </div>
  );
}

function FooterLink({ children, ...props }: React.ComponentProps<typeof Link>) {
  return (
    <Link
      {...props}
      className="text-sm text-white/45 transition-colors hover:text-white"
    >
      {children}
    </Link>
  );
}
