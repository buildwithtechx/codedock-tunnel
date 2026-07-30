import { createFileRoute, notFound } from '@tanstack/react-router';
import { DocsLayout } from 'fumadocs-ui/layouts/docs';
import {
  DocsBody,
  DocsDescription,
  DocsPage,
  DocsTitle,
} from 'fumadocs-ui/page';
import { useMDXComponents } from '../../components/mdx';
import { source } from '../../lib/source';

export const Route = createFileRoute('/docs/$')({
  loader: async ({ params }) => {
    const slugs = params._splat ? params._splat.split('/') : [];
    const page = source.getPage(slugs);
    if (!page) throw notFound();

    return {
      title: page.data.title,
      description: page.data.description,
      slugs,
    };
  },
  notFoundComponent: () => (
    <div className="flex min-h-[50vh] flex-col items-center justify-center p-8 text-center">
      <h1 className="text-2xl font-bold">404 - Document Not Found</h1>
      <p className="mt-2 text-muted-foreground">
        The requested documentation page could not be found.
      </p>
    </div>
  ),
  component: Page,
});

function Page() {
  const { slugs } = Route.useLoaderData();
  const mdxComponents = useMDXComponents();
  const page = source.getPage(slugs);
  if (!page) return null;

  const MDX = page.data.body;

  return (
    <DocsLayout tree={source.pageTree} nav={{ title: 'Codedock Tunnel' }}>
      <DocsPage toc={page.data.toc}>
        <DocsTitle>{page.data.title}</DocsTitle>
        <DocsDescription>{page.data.description}</DocsDescription>
        <DocsBody>
          <MDX components={mdxComponents} />
        </DocsBody>
      </DocsPage>
    </DocsLayout>
  );
}
