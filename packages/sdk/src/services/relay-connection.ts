import {
  type AuthResponse,
  type CloseTunnel,
  decodeMessage,
  type ErrorMessage,
  type HTTPRequest,
  type HTTPResponse,
  type MessageType,
  maxSupportedVersion,
  minSupportedVersion,
  type OpenTunnel,
  type OpenTunnelAck,
  type ProtocolEnvelope,
  type VersionNegotiateAck,
} from '@codedock/protocol-ts';
import type { WebSocketEvent, WebSocketLike } from '../interfaces/relay';
import { TunnelProtocolError, TunnelSDKError } from '../utils/errors';
import { RelayConnectionBase } from './relay-base';

export class RelayConnection extends RelayConnectionBase {
  private readonly managedTunnels = new Map<
    string,
    Omit<OpenTunnel, 'token'>
  >();

  async connect(): Promise<void> {
    if (this.authenticated) {
      return;
    }
    if (!this.connectPromise) {
      this.closedByUser = false;
      this.connectPromise = this.connectInternal().finally(() => {
        this.connectPromise = undefined;
      });
    }
    return this.connectPromise;
  }

  async openTunnel(request: Omit<OpenTunnel, 'token'>): Promise<OpenTunnelAck> {
    await this.connect();
    const response = await this.request<OpenTunnelAck>(
      'open_tunnel',
      { ...request, token: this.options.agentToken },
      'open_tunnel_ack',
    );
    this.managedTunnels.set(response.tunnel_id, request);
    this.emit('tunnel_opened', response);
    return response;
  }

  async closeTunnel(tunnelId: string, reason?: string): Promise<void> {
    this.managedTunnels.delete(tunnelId);
    await this.connect();
    const request = { tunnel_id: tunnelId, reason } satisfies CloseTunnel;
    this.send('close_tunnel', request);
    this.emit('tunnel_closed', request);
  }

  close(): void {
    this.closedByUser = true;
    this.authenticated = false;
    this.stopHeartbeat();
    this.clearReconnectTimer();
    this.rejectPending(new TunnelSDKError('relay connection closed'));
    this.socket?.close(1000, 'client closed');
    this.socket = undefined;
    this.managedTunnels.clear();
  }

  private async connectInternal(): Promise<void> {
    const factory = this.options.webSocket ?? defaultWebSocketFactory;
    const socket = factory(this.options.relayUrl);
    this.socket = socket;
    const opened = new Promise<void>((resolve, reject) => {
      this.openResolve = resolve;
      this.openReject = reject;
    });
    socket.addEventListener('open', this.handleOpen);
    socket.addEventListener('message', this.handleMessage);
    socket.addEventListener('close', this.handleClose);
    socket.addEventListener('error', this.handleError);
    try {
      await opened;
      const version = await this.request<VersionNegotiateAck>(
        'version_negotiate',
        {
          min_version: minSupportedVersion,
          max_version: maxSupportedVersion,
          client_name: this.options.clientName ?? 'codedock-sdk',
          client_version: this.options.clientVersion ?? '0.1.0',
        },
        'version_negotiate_ack',
      );
      const auth = await this.request<AuthResponse>(
        'auth',
        {
          token: this.options.agentToken,
          agent_id: this.options.agentId,
          requested_capabilities: ['http', 'https', 'tcp', 'udp'],
        },
        'auth_response',
      );
      if (!auth.authenticated) {
        throw new TunnelProtocolError(
          auth.error ?? 'relay authentication rejected',
        );
      }
      this.authenticated = true;
      this.reconnectAttempts = 0;
      this.startHeartbeat();
      this.emit('connected', version);
      this.emit('authenticated', auth);
      await this.reopenManagedTunnels();
    } catch (error) {
      this.cleanupSocket(socket);
      throw error;
    }
  }

  private request<TPayload>(
    type: MessageType,
    payload: unknown,
    expected: MessageType,
  ): Promise<TPayload> {
    const requestId = `${Date.now().toString(36)}-${(++this.requestCounter).toString(36)}`;
    return new Promise<TPayload>((resolve, reject) => {
      const timer = setTimeout(() => {
        this.pending.delete(requestId);
        reject(new TunnelProtocolError(`timed out waiting for ${expected}`));
      }, this.options.heartbeatIntervalMs);
      this.pending.set(requestId, {
        expected,
        resolve: (message) => resolve(message.payload as TPayload),
        reject,
        timer,
      });
      try {
        this.send(type, payload, requestId);
      } catch (error) {
        clearTimeout(timer);
        this.pending.delete(requestId);
        reject(
          error instanceof Error ? error : new TunnelSDKError(String(error)),
        );
      }
    });
  }

  private handleOpen = (): void => {
    this.openResolve?.();
    this.openResolve = undefined;
    this.openReject = undefined;
  };

  private handleMessage = (event: WebSocketEvent): void => {
    if (typeof event.data !== 'string') {
      this.fail(new TunnelProtocolError('relay returned a non-text frame'));
      return;
    }
    try {
      const message = decodeMessage(event.data);
      this.emit('message', message);
      if (message.request_id) {
        this.resolvePending(message);
      }
      if (message.type === 'http_request') {
        void this.forwardHTTPRequest(
          message.payload as HTTPRequest,
          message.request_id,
        );
      }
      if (message.type === 'error') {
        const error = message.payload as ErrorMessage;
        this.fail(new TunnelProtocolError(error.message, error.code));
      }
    } catch (error) {
      this.fail(
        error instanceof Error ? error : new TunnelProtocolError(String(error)),
      );
    }
  };

  private handleClose = (event: WebSocketEvent): void => {
    this.authenticated = false;
    this.stopHeartbeat();
    const error = new TunnelSDKError('relay connection closed');
    this.openReject?.(error);
    this.openReject = undefined;
    this.openResolve = undefined;
    this.rejectPending(new TunnelSDKError('relay connection closed'));
    const socket = this.socket;
    if (socket) {
      socket.removeEventListener('open', this.handleOpen);
      socket.removeEventListener('message', this.handleMessage);
      socket.removeEventListener('close', this.handleClose);
      socket.removeEventListener('error', this.handleError);
      this.socket = undefined;
    }
    this.emit('disconnected', event as CloseEvent);
    if (!this.closedByUser && this.options.reconnect) {
      this.scheduleReconnect();
    }
  };

  private handleError = (event: WebSocketEvent): void => {
    const error = new TunnelSDKError(event.message || 'relay WebSocket error');
    this.openReject?.(error);
    this.fail(error);
  };

  private cleanupSocket(socket: WebSocketLike): void {
    socket.removeEventListener('open', this.handleOpen);
    socket.removeEventListener('message', this.handleMessage);
    socket.removeEventListener('close', this.handleClose);
    socket.removeEventListener('error', this.handleError);
    if (this.socket === socket) {
      this.socket = undefined;
    }
    if (socket.readyState === 1 || socket.readyState === 0) {
      socket.close(1000, 'connection setup failed');
    }
  }

  private async reopenManagedTunnels(): Promise<void> {
    for (const [tunnelId, request] of this.managedTunnels) {
      try {
        const reopened = await this.request<OpenTunnelAck>(
          'open_tunnel',
          { ...request, tunnel_id: tunnelId, token: this.options.agentToken },
          'open_tunnel_ack',
        );
        this.emit('tunnel_opened', reopened);
      } catch (error) {
        this.managedTunnels.delete(tunnelId);
        this.emit(
          'error',
          error instanceof Error ? error : new TunnelSDKError(String(error)),
        );
      }
    }
  }

  private async forwardHTTPRequest(
    request: HTTPRequest,
    requestId?: string,
  ): Promise<void> {
    if (!this.options.localPort || !requestId || typeof fetch !== 'function') {
      return;
    }
    try {
      const headers = new Headers();
      for (const [name, values] of Object.entries(request.headers ?? {})) {
        headers.set(name, values.join(', '));
      }
      const response = await fetch(
        `http://127.0.0.1:${this.options.localPort}${request.path}`,
        {
          method: request.method,
          headers,
          body: request.body ? decodeBase64(request.body) : undefined,
        },
      );
      const responseHeaders: Record<string, string[]> = {};
      response.headers.forEach((value, name) => {
        responseHeaders[name] = [value];
      });
      this.send(
        'http_response',
        {
          status_code: response.status,
          headers: responseHeaders,
          body: encodeBase64(new Uint8Array(await response.arrayBuffer())),
        } satisfies HTTPResponse,
        requestId,
      );
    } catch (error) {
      this.send(
        'http_response',
        {
          status_code: 502,
          headers: {},
          error: error instanceof Error ? error.message : String(error),
        } satisfies HTTPResponse,
        requestId,
      );
    }
  }

  private resolvePending(message: ProtocolEnvelope): void {
    const requestId = message.request_id;
    if (!requestId) {
      return;
    }
    const pending = this.pending.get(requestId);
    if (!pending) {
      return;
    }
    this.pending.delete(requestId);
    clearTimeout(pending.timer);
    if (message.type !== pending.expected) {
      pending.reject(
        new TunnelProtocolError(
          `expected ${pending.expected}, received ${message.type}`,
        ),
      );
      return;
    }
    pending.resolve(message);
  }

  private startHeartbeat(): void {
    this.stopHeartbeat();
    this.heartbeatTimer = setInterval(() => {
      if (this.authenticated) {
        this.send('heartbeat', { timestamp: Math.floor(Date.now() / 1000) });
      }
    }, this.options.heartbeatIntervalMs);
  }

  private stopHeartbeat(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = undefined;
    }
  }

  private scheduleReconnect(): void {
    if (
      this.reconnectTimer ||
      this.reconnectAttempts >= this.options.maxReconnectAttempts
    ) {
      return;
    }
    const delay = this.options.reconnectDelayMs * 2 ** this.reconnectAttempts;
    this.reconnectAttempts += 1;
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = undefined;
      void this.connect().catch((error: unknown) =>
        this.fail(
          error instanceof Error ? error : new TunnelSDKError(String(error)),
        ),
      );
    }, delay);
  }

  private clearReconnectTimer(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = undefined;
    }
  }

  private rejectPending(error: Error): void {
    for (const [requestId, pending] of this.pending) {
      clearTimeout(pending.timer);
      pending.reject(error);
      this.pending.delete(requestId);
    }
  }

  private fail(error: Error): void {
    this.emit('error', error);
  }
}

function defaultWebSocketFactory(url: string): WebSocketLike {
  if (typeof WebSocket === 'undefined') {
    throw new TunnelSDKError(
      'WebSocket is not available; provide a webSocket factory',
    );
  }
  return new WebSocket(url);
}

function decodeBase64(value: string): Uint8Array {
  const binary = atob(value);
  return Uint8Array.from(binary, (character) => character.charCodeAt(0));
}

function encodeBase64(value: Uint8Array): string {
  let binary = '';
  for (const byte of value) {
    binary += String.fromCharCode(byte);
  }
  return btoa(binary);
}
