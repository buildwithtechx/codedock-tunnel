import { create } from 'zustand';
import type { Tunnel } from '#/interfaces/tunnel';

export type TunnelState = {
  tunnels: Tunnel[];
  selectedTunnel: Tunnel | null;
  setTunnels: (tunnels: Tunnel[]) => void;
  setSelectedTunnel: (tunnel: Tunnel | null) => void;
  upsertTunnel: (tunnel: Tunnel) => void;
  removeTunnel: (tunnelId: string) => void;
  clear: () => void;
};

export const useTunnelStore = create<TunnelState>((set) => ({
  tunnels: [],
  selectedTunnel: null,
  setTunnels: (tunnels) => set({ tunnels }),
  setSelectedTunnel: (selectedTunnel) => set({ selectedTunnel }),
  upsertTunnel: (tunnel) =>
    set((state) => ({
      tunnels: [
        ...state.tunnels.filter((item) => item.id !== tunnel.id),
        tunnel,
      ],
      selectedTunnel:
        state.selectedTunnel?.id === tunnel.id ? tunnel : state.selectedTunnel,
    })),
  removeTunnel: (tunnelId) =>
    set((state) => ({
      tunnels: state.tunnels.filter((tunnel) => tunnel.id !== tunnelId),
      selectedTunnel:
        state.selectedTunnel?.id === tunnelId ? null : state.selectedTunnel,
    })),
  clear: () => set({ tunnels: [], selectedTunnel: null }),
}));
