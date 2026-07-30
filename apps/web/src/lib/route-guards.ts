type OrganizationRouteContext = {
  organization?: { id?: string };
  isMember?: boolean;
};

type PlatformAdminRouteContext = {
  isPlatformAdmin?: boolean;
};

export function requireOrganization(context: OrganizationRouteContext) {
  if (!context.organization?.id || context.isMember !== true) {
    throw new Error('organization membership is required');
  }
  return context;
}

export function requirePlatformAdmin(context: PlatformAdminRouteContext) {
  if (context.isPlatformAdmin !== true) {
    throw new Error('platform admin access is required');
  }
  return context;
}
