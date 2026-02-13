import { fetchWithAuth } from './http-client';

export interface UserResponse {
  id: number;
  name: string;
  status: string;
  email?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export const authApi = {
  async getYandexAuthUrl(): Promise<string> {
    const response = await fetchWithAuth('/auth/yandex/url');

    if (!response.ok) {
      throw new Error('Failed to get Yandex OAuth URL');
    }

    const data = await response.json();
    return data.url;
  },

  async login(credentials: LoginRequest): Promise<void> {
    const response = await fetchWithAuth('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Login failed');
    }
  },

  async register(data: RegisterRequest): Promise<void> {
    const response = await fetchWithAuth('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Registration failed');
    }
  },

  async getCurrentUser(): Promise<UserResponse> {
    const response = await fetchWithAuth('/auth/me');

    if (!response.ok) {
      throw new Error('Not authenticated');
    }

    return response.json();
  },

  async logout(): Promise<void> {
    await fetchWithAuth('/auth/logout', {
      method: 'POST',
    });
  },
};
