// Note: We use direct fetch calls instead of apiClient to avoid circular dependencies
import { API_ENDPOINTS } from '../constants';
import type { LoginRequest, RegisterRequest, AuthUser } from '../types/auth';

class AuthService {
  private async requestWithCredentials<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const response = await fetch(`${API_ENDPOINTS.BASE_URL}${endpoint}`, {
      ...options,
      credentials: 'include', // Important for cookies
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  async login(credentials: LoginRequest): Promise<void> {
    await this.requestWithCredentials(API_ENDPOINTS.AUTH.LOGIN, {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  }

  async register(userData: RegisterRequest): Promise<AuthUser> {
    const response = await this.requestWithCredentials<{ user: AuthUser }>(
      API_ENDPOINTS.AUTH.REGISTER,
      {
        method: 'POST',
        body: JSON.stringify(userData),
      }
    );
    return response.user;
  }

  async refreshToken(): Promise<boolean> {
    try {
      await this.requestWithCredentials('/auth/refresh', {
        method: 'POST',
      });
      return true;
    } catch (error) {
      console.error('Token refresh failed:', error);
      return false;
    }
  }

  async validateToken(): Promise<boolean> {
    try {
      await this.requestWithCredentials('/validate', {
        method: 'GET',
      });
      return true;
    } catch (error) {
      console.error('Token validation failed:', error);
      return false;
    }
  }

  async logout(): Promise<void> {
    // Clear cookies by setting them to expire in the past
    document.cookie = 'accessToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    document.cookie = 'refreshToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
  }
}

export const authService = new AuthService();
