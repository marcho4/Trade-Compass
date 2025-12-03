const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '/api';

let isRefreshing = false;

let failedQueue: Array<{
  resolve: (value: unknown) => void;
  reject: (reason?: unknown) => void;
}> = [];

const processQueue = (error: Error | null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(undefined);
    }
  });
  failedQueue = [];
};


async function refreshTokens(): Promise<boolean> {
  try {
    const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (response.ok) {
      return true;
    }
    return false;
  } catch {
    return false;
  }
}

async function fetchWithAuth(
  url: string,
  options: RequestInit = {}
): Promise<Response> {
  const defaultOptions: RequestInit = {
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  };

  const mergedOptions = { ...defaultOptions, ...options };
  
  let response = await fetch(`${API_BASE_URL}${url}`, mergedOptions);

  if (response.status === 401 && !url.includes('/auth/')) {
    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        failedQueue.push({ resolve, reject });
      }).then(() => {
        return fetch(`${API_BASE_URL}${url}`, mergedOptions);
      });
    }

    isRefreshing = true;

    const refreshSuccess = await refreshTokens();

    if (refreshSuccess) {
      isRefreshing = false;
      processQueue(null);
      response = await fetch(`${API_BASE_URL}${url}`, mergedOptions);
    } else {
      isRefreshing = false;
      processQueue(new Error('Refresh token failed'));
      if (typeof window !== 'undefined') {
        window.location.href = '/auth';
      }
    }
  }

  return response;
}

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
  /**
   * Получить URL для OAuth авторизации через Яндекс
   * Бэкенд генерирует state и сохраняет в cookie
   */
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

  /**
   * Авторизация по email/password
   */
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

  /**
   * Регистрация нового пользователя
   */
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

  /**
   * Получить данные текущего пользователя
   */
  async getCurrentUser(): Promise<UserResponse> {
    const response = await fetchWithAuth('/auth/me', {
      method: 'GET',
    });
    
    if (!response.ok) {
      throw new Error('Not authenticated');
    }
    
    return response.json();
  },

  /**
   * Выход из системы
   */
  async logout(): Promise<void> {
    await fetchWithAuth('/auth/logout', {
      method: 'POST',
    });
  },

  /**
   * Обновить токены
   */
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

export const api = {
  async get<T>(url: string): Promise<T> {
    const response = await fetchWithAuth(url);
    if (!response.ok) {
      throw new Error(`GET ${url} failed`);
    }
    return response.json();
  },

  async post<T>(url: string, data?: unknown): Promise<T> {
    const response = await fetchWithAuth(url, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
    if (!response.ok) {
      throw new Error(`POST ${url} failed`);
    }
    return response.json();
  },

  async put<T>(url: string, data: unknown): Promise<T> {
    const response = await fetchWithAuth(url, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      throw new Error(`PUT ${url} failed`);
    }
    return response.json();
  },

  async delete<T>(url: string): Promise<T> {
    const response = await fetchWithAuth(url, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error(`DELETE ${url} failed`);
    }
    return response.json();
  },
};

export default api;

