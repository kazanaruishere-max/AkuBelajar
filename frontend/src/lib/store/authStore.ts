/**
 * AkuBelajar — Auth Store (Zustand)
 *
 * Manages authentication state:
 * - Access & refresh tokens (in-memory, NOT localStorage)
 * - Current user data
 * - Login/logout actions
 */

import { create } from 'zustand';

export type UserRole = 'super_admin' | 'teacher' | 'class_leader' | 'student';

export interface User {
  id: string;
  email: string;
  role: UserRole;
  school_id: string;
  is_first_login: boolean;
  profile?: {
    nisn?: string;
    nip?: string;
    phone_wa?: string;
    photo_url?: string;
    birth_date?: string;
    parent_name?: string;
  };
  last_login_at?: string;
}

interface AuthState {
  // State
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // Actions
  setUser: (user: User) => void;
  setTokens: (accessToken: string, refreshToken: string) => void;
  login: (user: User, accessToken: string, refreshToken: string) => void;
  logout: () => void;
  setLoading: (loading: boolean) => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  // Initial state
  user: null,
  accessToken: null,
  refreshToken: null,
  isAuthenticated: false,
  isLoading: true,

  // Actions
  setUser: (user) => set({ user }),

  setTokens: (accessToken, refreshToken) =>
    set({ accessToken, refreshToken }),

  login: (user, accessToken, refreshToken) =>
    set({
      user,
      accessToken,
      refreshToken,
      isAuthenticated: true,
      isLoading: false,
    }),

  logout: () =>
    set({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,
    }),

  setLoading: (isLoading) => set({ isLoading }),
}));
