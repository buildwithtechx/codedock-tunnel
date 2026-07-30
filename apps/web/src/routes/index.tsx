import { createFileRoute, Link } from '@tanstack/react-router';

export const Route = createFileRoute('/')({ component: Home });

function Home() {
  return (
    <main className="min-h-screen bg-fd-background text-fd-foreground">
      <section className="mx-auto flex min-h-screen max-w-6xl flex-col justify-center px-6 py-20 lg:px-8">
        <div className="max-w-3xl">
          <div className="mb-8 flex items-center gap-3 text-sm font-medium text-fd-muted-foreground">
            <img src="/favicon.svg" alt="" className="size-9 rounded-xl" />
            <span>Codedock Tunnel</span>
          </div>
          <h1 className="max-w-3xl text-5xl font-semibold tracking-tight sm:text-7xl">
            Secure public access for local services.
          </h1>
          <p className="mt-6 max-w-2xl text-lg leading-8 text-fd-muted-foreground sm:text-xl">
            Open reliable HTTP, TCP, and HTTPS tunnels from your workstation or
            CI pipeline with one protocol, one CLI, and a developer-first SDK.
          </p>
          <div className="mt-10 flex flex-wrap gap-3">
            <Link
              to="/docs/installation"
              className="rounded-lg bg-fd-primary px-5 py-3 font-medium text-fd-primary-foreground transition-opacity hover:opacity-90"
            >
              Get started
            </Link>
            <Link
              to="/docs"
              className="rounded-lg border border-fd-border px-5 py-3 font-medium transition-colors hover:bg-fd-accent"
            >
              Read the docs
            </Link>
          </div>
        </div>
        <div className="mt-20 grid gap-4 border-t border-fd-border pt-8 sm:grid-cols-3">
          <div>
            <p className="font-medium">One client experience</p>
            <p className="mt-2 text-sm text-fd-muted-foreground">
              CLI, desktop app, and framework adapters share the same lifecycle.
            </p>
          </div>
          <div>
            <p className="font-medium">Built for real workflows</p>
            <p className="mt-2 text-sm text-fd-muted-foreground">
              Preview apps, receive webhooks, test OAuth, and run CI tunnels.
            </p>
          </div>
          <div>
            <p className="font-medium">Observable by default</p>
            <p className="mt-2 text-sm text-fd-muted-foreground">
              Track requests, bandwidth, connections, and tunnel health.
            </p>
          </div>
        </div>
      </section>
    </main>
  );
}
