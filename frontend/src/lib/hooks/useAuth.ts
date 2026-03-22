'use client';

import { useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore, type User } from '@/lib/store/authStore';
import { api } from '@/lib/api/client';

interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

/**
 * useAuth — React hook for auth operations.
 * Wraps Zustand store + API calls for login, logout, and refresh.
 */
export function useAuth() {
  const router = useRouter();
  const { user, accessToken, login: storeLogin, logout: clearAuth } = useAuthStore();

  const login = useCallback(async (email: string, password: string, rememberMe = false) => {
    const res = await api.post<LoginResponse>('/auth/login', {
      email,
      password,
      remember_me: rememberMe,
    });
    const { user, access_token, refresh_token } = res.data;
    storeLogin(user, access_token, refresh_token);
    return res.data;
  }, [storeLogin]);

  const logout = useCallback(async () => {
    try {
      const rt = useAuthStore.getState().refreshToken;
      if (rt) {
        await api.post('/auth/logout', { refresh_token: rt });
      }
    } catch {
      // Ignore logout API errors — just clear local state
    }
    clearAuth();
    router.push('/login');
  }, [clearAuth, router]);

  const refreshToken = useCallback(async () => {
    const currentRefresh = useAuthStore.getState().refreshToken;
    if (!currentRefresh) throw new Error('No refresh token');

    const res = await api.post<LoginResponse>('/auth/refresh', {
      refresh_token: currentRefresh,
    });
    const { user, access_token, refresh_token } = res.data;
    storeLogin(user, access_token, refresh_token);
    return res.data;
  }, [storeLogin]);

  return {
    user,
    isAuthenticated: !!accessToken,
    login,
    logout,
    refreshToken,
  };
}
