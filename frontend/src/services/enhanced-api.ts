import { getApiUrl } from '../config/env';
import type { ApiResponse } from '../types';
import { authService } from './auth';

// Enhanced API Client with automatic token refresh
class EnhancedApiClient {
  constructor() {
    // Enhanced API client with token refresh capabilities
  }

  private async requestWithTokenRefresh<T>(
    endpoint: string,
    options: RequestInit = {},
    retryCount = 0
  ): Promise<ApiResponse<T>> {
    const url = getApiUrl(endpoint);
    
    const defaultHeaders = {
      'Content-Type': 'application/json',
    };

    // Get access token from cookies (set by backend)
    const accessToken = this.getCookieValue('Authorization');

    const config: RequestInit = {
      ...options,
      headers: {
        ...defaultHeaders,
        ...(accessToken && { 'Authorization': `Bearer ${accessToken}` }),
        ...options.headers,
      },
    };

    try {
      const response = await fetch(url, config);
      
      // If unauthorized and we have a refresh token, try to refresh
      if (response.status === 401 && authService.getRefreshToken() && retryCount === 0) {
        try {
          await authService.refreshAccessToken();
          // Retry the request with new token
          return this.requestWithTokenRefresh(endpoint, options, retryCount + 1);
        } catch (refreshError) {
          // Refresh failed, redirect to login or handle as needed
          console.error('Token refresh failed:', refreshError);
          // You might want to redirect to login page here
          throw new Error('Authentication failed');
        }
      }
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  private getCookieValue(name: string): string | null {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) {
      return parts.pop()?.split(';').shift() || null;
    }
    return null;
  }

  // GET request
  async get<T>(endpoint: string, options?: RequestInit): Promise<ApiResponse<T>> {
    return this.requestWithTokenRefresh<T>(endpoint, {
      method: 'GET',
      ...options,
    });
  }

  // POST request
  async post<T>(
    endpoint: string, 
    data?: any, 
    options?: RequestInit
  ): Promise<ApiResponse<T>> {
    return this.requestWithTokenRefresh<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
      ...options,
    });
  }

  // PUT request
  async put<T>(
    endpoint: string, 
    data?: any, 
    options?: RequestInit
  ): Promise<ApiResponse<T>> {
    return this.requestWithTokenRefresh<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
      ...options,
    });
  }

  // DELETE request
  async delete<T>(endpoint: string, options?: RequestInit): Promise<ApiResponse<T>> {
    return this.requestWithTokenRefresh<T>(endpoint, {
      method: 'DELETE',
      ...options,
    });
  }
}

// Export singleton instance
export const enhancedApiClient = new EnhancedApiClient();
