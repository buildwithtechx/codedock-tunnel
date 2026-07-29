export const protocolVersion = 1 as const;

export const messageTypes = [
  'open_tunnel',
  'open_tunnel_ack',
  'close_tunnel',
  'data',
  'heartbeat',
  'error',
] as const;

export type MessageType = (typeof messageTypes)[number];

export type ProtocolEnvelope<TPayload = unknown> = {
  version: typeof protocolVersion;
  type: MessageType;
  request_id?: string;
  payload?: TPayload;
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

export const encodeMessage = <TPayload>(
  message: ProtocolEnvelope<TPayload>,
): string => JSON.stringify(message);

export const decodeMessage = <TPayload>(
  value: string,
): ProtocolEnvelope<TPayload> => {
  const message = JSON.parse(value) as ProtocolEnvelope<TPayload>;
  if (message.version !== protocolVersion || !messageTypes.includes(message.type)) {
    throw new Error('unsupported tunnel protocol message');
  }
  return message;
};
