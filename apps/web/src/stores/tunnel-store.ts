import { create } from 'zustand';
import type { Tunnel } from '../interfaces/tunnel';

type TunnelStore = {
  tunnels: Tunnel[];
  setTunnels: (tunnels: Tunnel[]) => void;
};

export const useTunnelStore = create<TunnelStore>((set) => ({
  tunnels: [],
  setTunnels: (tunnels) => set({ tunnels }),
}));
