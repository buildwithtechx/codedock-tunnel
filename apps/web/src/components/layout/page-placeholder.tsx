type PagePlaceholderProps = {
  title: string;
  description?: string;
};

export function PagePlaceholder({ title, description }: PagePlaceholderProps) {
  return (
    <main className="mx-auto flex min-h-[60vh] max-w-6xl flex-col justify-center px-6 py-16">
      <h1 className="text-3xl font-semibold tracking-tight">{title}</h1>
      {description && (
        <p className="mt-3 max-w-2xl text-muted-foreground">{description}</p>
      )}
    </main>
  );
}
