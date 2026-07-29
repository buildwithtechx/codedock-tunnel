export const protocolVersion = 1 as const;
export const minSupportedVersion = 1 as const;
export const maxSupportedVersion = 1 as const;

export const defaultMaxFrameSize = 16 * 1024 * 1024;
export const absoluteMaxFrameSize = 32 * 1024 * 1024;
export const defaultIdleTimeoutSeconds = 60;
export const defaultConnectionTimeoutSeconds = 30;

export const messageTypes = [
  'auth',
  'auth_response',
  'version_negotiate',
  'version_negotiate_ack',
  'flow_control',
  'open_tunnel',
  'open_tunnel_ack',
  'close_tunnel',
  'data',
  'heartbeat',
  'error',
  'http_request',
  'http_response',
  'tcp_data',
  'tcp_close',
  'udp_data',
  'udp_response',
] as const;

export type MessageType = (typeof messageTypes)[number];

export type ProtocolEnvelope<TPayload = unknown> = {
  version: number;
  type: MessageType;
  request_id?: string;
  payload?: TPayload;
};

export type AuthRequest = {
  token: string;
  agent_id?: string;
  requested_capabilities?: string[];
};

export type AuthResponse = {
  authenticated: boolean;
  agent_id?: string;
  organization_id?: string;
  granted_capabilities?: string[];
  error?: string;
};

export type VersionNegotiate = {
  min_version: number;
  max_version: number;
  client_name?: string;
  client_version?: string;
};

export type VersionNegotiateAck = {
  negotiated_version: number;
  supported_versions: number[];
  server_name?: string;
  server_version?: string;
};

export type FlowControl = {
  stream_id: string;
  action: 'pause' | 'resume';
  window_size?: number;
};

export type OpenTunnel = {
  token: string;
  local_port: number;
  subdomain?: string;
  protocol: string;
  custom_domain?: string;
};

export type OpenTunnelAck = {
  tunnel_id: string;
  public_url: string;
  public_port?: number;
};

export type CloseTunnel = {
  tunnel_id: string;
  reason?: string;
};

export type Data = {
  tunnel_id: string;
  stream_id: string;
  data: string;
};

export type Heartbeat = {
  timestamp: number;
};

export type ErrorMessage = {
  code: string;
  message: string;
};

export type HTTPRequest = {
  method: string;
  path: string;
  headers: Record<string, string[]>;
  body?: string;
};

export type HTTPResponse = {
  status_code: number;
  headers: Record<string, string[]>;
  body?: string;
  error?: string;
};

export type TCPData = {
  connection_id: string;
  data: string;
};

export type TCPClose = {
  connection_id: string;
  reason?: string;
};

export type UDPData = {
  packet_id: string;
  source_address: string;
  source_port: number;
  data: string;
};

export type UDPResponse = {
  packet_id: string;
  target_address: string;
  target_port: number;
  data: string;
};

export const encodeMessage = <TPayload>(
  message: ProtocolEnvelope<TPayload>,
): string => {
  const result = JSON.stringify(message);
  if (result.length > absoluteMaxFrameSize) {
    throw new Error(
      `frame size ${result.length} exceeds maximum allowed frame size ${absoluteMaxFrameSize}`,
    );
  }
  return result;
};

export const decodeMessage = <TPayload>(
  value: string,
): ProtocolEnvelope<TPayload> => {
  if (value.length > absoluteMaxFrameSize) {
    throw new Error(
      `frame size ${value.length} exceeds maximum allowed frame size ${absoluteMaxFrameSize}`,
    );
  }
  const message = JSON.parse(value) as ProtocolEnvelope<TPayload>;
  if (
    message.version < minSupportedVersion ||
    message.version > maxSupportedVersion ||
    !messageTypes.includes(message.type)
  ) {
    throw new Error('unsupported tunnel protocol message');
  }
  return message;
};

export const negotiateVersion = (
  req: VersionNegotiate,
): VersionNegotiateAck => {
  if (
    req.max_version < minSupportedVersion ||
    req.min_version > maxSupportedVersion
  ) {
    throw new Error(
      `incompatible protocol version request: ${req.min_version}-${req.max_version}`,
    );
  }
  const negotiated = Math.min(req.max_version, maxSupportedVersion);
  return {
    negotiated_version: negotiated,
    supported_versions: [protocolVersion],
    server_name: 'codedock-tunnel',
    server_version: '0.1.0',
  };
};
