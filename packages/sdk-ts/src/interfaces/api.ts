export type TunnelAPIClientOptions = {
  apiUrl: string;
  apiKey?: string;
  fetch?: typeof globalThis.fetch;
  apiPrefix?: string;
};

export type Tunnel = {
  id: string;
  publicUrl?: string;
  public_url?: string;
  status: string;
  [key: string]: unknown;
};
