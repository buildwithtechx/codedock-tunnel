export function requireOrganization(context: unknown) {
  if (!context) throw new Error('organization route context is required');
  return context;
}

export function requirePlatformAdmin(context: unknown) {
  if (!context) throw new Error('platform admin route context is required');
  return context;
}
