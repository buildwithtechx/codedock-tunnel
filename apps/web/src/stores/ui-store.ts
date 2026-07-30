import { create } from 'zustand';

export type UIState = {
  sidebarOpen: boolean;
  commandPaletteOpen: boolean;
  setSidebarOpen: (sidebarOpen: boolean) => void;
  setCommandPaletteOpen: (commandPaletteOpen: boolean) => void;
};

export const useUIStore = create<UIState>((set) => ({
  sidebarOpen: true,
  commandPaletteOpen: false,
  setSidebarOpen: (sidebarOpen) => set({ sidebarOpen }),
  setCommandPaletteOpen: (commandPaletteOpen) => set({ commandPaletteOpen }),
}));
