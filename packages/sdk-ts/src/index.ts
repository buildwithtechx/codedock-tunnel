import {
  encodeMessage,
  type OpenTunnel,
  type ProtocolEnvelope,
} from '@codedock-tunnel/protocol-ts';

export type {
  OpenTunnel,
  ProtocolEnvelope,
} from '@codedock-tunnel/protocol-ts';

export type TunnelClientOptions = {
  apiUrl: string;
  apiKey?: string;
  fetch?: typeof globalThis.fetch;
};

export type CreateTunnelRequest = OpenTunnel;

export type Tunnel = {
  id: string;
  publicUrl: string;
  status: string;
};

export class TunnelAPIError extends Error {
  readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = 'TunnelAPIError';
    this.status = status;
  }
}

export class TunnelClient {
  private readonly apiUrl: string;
  private readonly apiKey?: string;
  private readonly request: typeof globalThis.fetch;

  constructor(options: TunnelClientOptions) {
    this.apiUrl = options.apiUrl.replace(/\/$/, '');
    this.apiKey = options.apiKey;
    this.request = options.fetch ?? globalThis.fetch;
  }

  async createTunnel(request: CreateTunnelRequest): Promise<Tunnel> {
    return this.call<Tunnel>('/v1/tunnels', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  async closeTunnel(tunnelId: string): Promise<void> {
    await this.call<void>(`/v1/tunnels/${encodeURIComponent(tunnelId)}`, {
      method: 'DELETE',
    });
  }

  createOpenMessage(
    request: CreateTunnelRequest,
  ): ProtocolEnvelope<OpenTunnel> {
    return {
      version: 1,
      type: 'open_tunnel',
      payload: request,
    };
  }

  encodeOpenMessage(request: CreateTunnelRequest): string {
    return encodeMessage(this.createOpenMessage(request));
  }

  private async call<T>(path: string, init: RequestInit): Promise<T> {
    const headers = new Headers(init.headers);
    headers.set('Accept', 'application/json');
    if (init.body) {
      headers.set('Content-Type', 'application/json');
    }
    if (this.apiKey) {
      headers.set('Authorization', `Bearer ${this.apiKey}`);
    }

    const response = await this.request(`${this.apiUrl}${path}`, {
      ...init,
      headers,
    });
    if (!response.ok) {
      throw new TunnelAPIError(response.status, await response.text());
    }
    if (response.status === 204) {
      return undefined as T;
    }
    return (await response.json()) as T;
  }
}
