import { fetchWithAuth, API_BASE_URL } from './http-client';

export interface AuthResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

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
    const response = await fetchWithAuth('/auth/yandex/url', {
      method: 'GET',
    });
    if (!response.ok) {
      throw new Error('Failed to get Yandex OAuth URL');
    }
    const data = await response.json();
    return data.url;
  },

  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response = await fetchWithAuth('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Login failed');
    }

    return response.json();
  },

  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await fetchWithAuth('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Registration failed');
    }

    return response.json();
  },

  async getCurrentUser(): Promise<UserResponse> {
    const response = await fetchWithAuth('/auth/me', {
      method: 'GET',
    });

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

  async refresh(): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Refresh failed');
    }

    return response.json();
  },
};
