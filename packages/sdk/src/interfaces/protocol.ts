export type MessageType =
  | 'auth'
  | 'auth_response'
  | 'version_negotiate'
  | 'version_negotiate_ack'
  | 'flow_control'
  | 'open_tunnel'
  | 'open_tunnel_ack'
  | 'close_tunnel'
  | 'data'
  | 'heartbeat'
  | 'error'
  | 'http_request'
  | 'http_response'
  | 'tcp_data'
  | 'tcp_close'
  | 'udp_data'
  | 'udp_response';

export type ProtocolEnvelope<TPayload = unknown> = {
  version: number;
  type: MessageType;
  request_id?: string;
  payload?: TPayload;
};

export type AuthResponse = {
  authenticated: boolean;
  agent_id?: string;
  organization_id?: string;
  granted_capabilities?: string[];
  error?: string;
};

export type VersionNegotiateAck = {
  negotiated_version: number;
  supported_versions: number[];
  server_name?: string;
  server_version?: string;
};

export type OpenTunnel = {
  token?: string;
  tunnel_id?: string;
  local_port: number;
  subdomain?: string;
  protocol: 'http' | 'https' | 'tcp' | 'udp' | (string & {});
  custom_domain?: string;
  password?: string;
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
