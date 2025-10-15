import { apiClient } from './api';
import { API_ENDPOINTS } from '../constants';

export interface LoginRequest {
  login: string;
  password: string;
}

export interface TokenResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

class AuthService {
  private refreshToken: string | null = null;

  constructor() {
    // Load refresh token from localStorage on initialization
    this.refreshToken = localStorage.getItem('refresh_token');
  }

  /**
   * Login user and store tokens
   */
  async login(credentials: LoginRequest): Promise<TokenResponse> {
    try {
      const response = await apiClient.post<TokenResponse>(
        API_ENDPOINTS.AUTH.LOGIN,
        credentials
      );

      if (response.access_token && response.refresh_token) {
        // Store refresh token in localStorage
        localStorage.setItem('refresh_token', response.refresh_token);
        this.refreshToken = response.refresh_token;
      }

      return response;
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    }
  }

  /**
   * Refresh access token using refresh token
   */
  async refreshAccessToken(): Promise<TokenResponse> {
    if (!this.refreshToken) {
      throw new Error('No refresh token available');
    }

    try {
      const response = await apiClient.post<TokenResponse>(
        API_ENDPOINTS.AUTH.REFRESH,
        { refresh_token: this.refreshToken } as RefreshTokenRequest
      );

      if (response.access_token && response.refresh_token) {
        // Update stored refresh token
        localStorage.setItem('refresh_token', response.refresh_token);
        this.refreshToken = response.refresh_token;
      }

      return response;
    } catch (error) {
      console.error('Token refresh failed:', error);
      // Clear invalid refresh token
      this.clearTokens();
      throw error;
    }
  }

  /**
   * Logout user and clear tokens
   */
  async logout(): Promise<void> {
    try {
      await apiClient.post(API_ENDPOINTS.AUTH.LOGOUT);
    } catch (error) {
      console.error('Logout request failed:', error);
    } finally {
      this.clearTokens();
    }
  }

  /**
   * Clear all stored tokens
   */
  private clearTokens(): void {
    localStorage.removeItem('refresh_token');
    this.refreshToken = null;
  }

  /**
   * Get current refresh token
   */
  getRefreshToken(): string | null {
    return this.refreshToken;
  }

  /**
   * Check if user has a valid refresh token
   */
  isLoggedIn(): boolean {
    return this.refreshToken !== null;
  }
}

// Export singleton instance
export const authService = new AuthService();
