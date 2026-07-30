import { create } from 'zustand';

type UIStore = {
  sidebarOpen: boolean;
  setSidebarOpen: (sidebarOpen: boolean) => void;
};

export const useUIStore = create<UIStore>((set) => ({
  sidebarOpen: true,
  setSidebarOpen: (sidebarOpen) => set({ sidebarOpen }),
}));
