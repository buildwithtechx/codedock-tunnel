import { create } from 'zustand';
import type { AuthUser } from '#/interfaces/auth';

export type AuthState = {
  user: AuthUser | null;
  isAuthenticated: boolean;
  setUser: (user: AuthUser | null) => void;
  clear: () => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  setUser: (user) => set({ user, isAuthenticated: user !== null }),
  clear: () => set({ user: null, isAuthenticated: false }),
}));
