import { create } from 'zustand';
import type { AuthUser } from '../interfaces/auth';

type AuthStore = {
  user?: AuthUser;
  setUser: (user?: AuthUser) => void;
};

export const useAuthStore = create<AuthStore>((set) => ({
  setUser: (user) => set({ user }),
}));
