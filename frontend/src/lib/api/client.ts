/**
 * AkuBelajar — Type-safe API Client
 *
 * Handles:
 * - Base URL configuration
 * - Auto-attach Paseto token from authStore
 * - 401 → auto refresh token
 * - Standard error format parsing
 */

import { useAuthStore } from '@/lib/store/authStore';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

interface ApiError {
  code: string;
  message: string;
  details?: unknown;
}

interface ApiResponse<T> {
  data: T;
  meta?: {
    page: number;
    per_page: number;
    total: number;
    total_page: number;
  };
}

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    const token = useAuthStore.getState().accessToken;
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    return headers;
  }

  private async handleResponse<T>(response: Response): Promise<ApiResponse<T>> {
    if (response.status === 204) {
      return { data: {} as T };
    }

    const body = await response.json();

    if (!response.ok) {
      // Auto refresh on 401
      if (response.status === 401) {
        const refreshed = await this.tryRefresh();
        if (refreshed) {
          // Retry the original request would need to be handled by the caller
          throw new ApiRequestError('TOKEN_REFRESHED', 'Token refreshed, please retry', 401);
        }
        // Refresh failed — logout
        useAuthStore.getState().logout();
        throw new ApiRequestError('AUTH_EXPIRED', 'Sesi habis, silakan login kembali', 401);
      }

      const error = body.error as ApiError;
      throw new ApiRequestError(
        error?.code || 'UNKNOWN',
        error?.message || 'Terjadi kesalahan',
        response.status,
        error?.details,
      );
    }

    return body as ApiResponse<T>;
  }

  private async tryRefresh(): Promise<boolean> {
    const { refreshToken, setTokens } = useAuthStore.getState();
    if (!refreshToken) return false;

    try {
      const res = await fetch(`${this.baseURL}/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (!res.ok) return false;

      const body = await res.json();
      setTokens(body.data.access_token, body.data.refresh_token);
      return true;
    } catch {
      return false;
    }
  }

  async get<T>(path: string, params?: Record<string, string>): Promise<ApiResponse<T>> {
    const url = new URL(`${this.baseURL}${path}`);
    if (params) {
      Object.entries(params).forEach(([key, value]) => url.searchParams.set(key, value));
    }

    const response = await fetch(url.toString(), {
      method: 'GET',
      headers: this.getHeaders(),
      credentials: 'include',
    });

    return this.handleResponse<T>(response);
  }

  async post<T>(path: string, body?: unknown): Promise<ApiResponse<T>> {
    const response = await fetch(`${this.baseURL}${path}`, {
      method: 'POST',
      headers: this.getHeaders(),
      credentials: 'include',
      body: body ? JSON.stringify(body) : undefined,
    });

    return this.handleResponse<T>(response);
  }

  async put<T>(path: string, body?: unknown): Promise<ApiResponse<T>> {
    const response = await fetch(`${this.baseURL}${path}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      credentials: 'include',
      body: body ? JSON.stringify(body) : undefined,
    });

    return this.handleResponse<T>(response);
  }

  async delete<T>(path: string): Promise<ApiResponse<T>> {
    const response = await fetch(`${this.baseURL}${path}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
      credentials: 'include',
    });

    return this.handleResponse<T>(response);
  }

  async upload<T>(path: string, formData: FormData): Promise<ApiResponse<T>> {
    const headers: HeadersInit = {};
    const token = useAuthStore.getState().accessToken;
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    // Don't set Content-Type — browser will set multipart/form-data with boundary

    const response = await fetch(`${this.baseURL}${path}`, {
      method: 'POST',
      headers,
      credentials: 'include',
      body: formData,
    });

    return this.handleResponse<T>(response);
  }
}

export class ApiRequestError extends Error {
  code: string;
  statusCode: number;
  details?: unknown;

  constructor(code: string, message: string, statusCode: number, details?: unknown) {
    super(message);
    this.code = code;
    this.statusCode = statusCode;
    this.details = details;
    this.name = 'ApiRequestError';
  }
}

export const api = new ApiClient(API_BASE_URL);
