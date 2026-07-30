export type Tunnel = {
  id: string;
  name: string;
  protocol: 'http' | 'https' | 'tcp' | 'udp';
  status: string;
  publicHostname: string;
  targetHost: string;
  targetPort: number;
};
