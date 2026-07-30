export type AuthUser = {
  id: string;
  email: string;
  name: string;
  status: 'active' | 'disabled';
};

export type AuthSession = {
  user: AuthUser;
  expiresAt: string;
};
