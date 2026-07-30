export type Subscription = {
  id: string;
  organizationId: string;
  planId: string;
  provider: 'polar' | 'paystack';
  status: string;
  currentPeriodEnd?: string;
};
