import type { Entity } from '#/interfaces/api';
import type { BillingProvider } from '#/interfaces/billing';

export type BillingEvent = Entity & {
  provider: BillingProvider;
  providerEventId: string;
  organizationId?: string;
  eventType: string;
  payloadHash: string;
  processedAt?: string;
  failureReason?: string;
};

export type Invoice = Entity & {
  organizationId: string;
  subscriptionId: string;
  amountMinor: number;
  currency: string;
  status: string;
  invoiceUrl?: string;
  pdfUrl?: string;
  paidAt?: string;
};

export type Receipt = Entity & {
  organizationId: string;
  invoiceId: string;
  receiptNumber: string;
  amountMinor: number;
  currency: string;
  issuedAt: string;
};

export type BillingCredential = Entity & {
  organizationId: string;
  provider: BillingProvider;
  kind: string;
  rotatedAt: string;
  expiresAt?: string;
};
