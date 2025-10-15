import { getApiUrl } from '../config/env';
import type { ApiResponse } from '../types';

// API Client class
class ApiClient {
  constructor() {
    // API_BASE_URL is used in getApiUrl function
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = getApiUrl(endpoint);
    
    const defaultHeaders = {
      'Content-Type': 'application/json',
    };

    const config: RequestInit = {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
    };

    try {
      const response = await fetch(url, config);

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      let parsed: unknown;

      try {
        parsed = await response.json();
      } catch {
        parsed = null;
      }

      const isObject = (value: unknown): value is Record<string, unknown> => {
        return typeof value === 'object' && value !== null;
      };

      const parsedObject = isObject(parsed) ? parsed : null;
      const hasData = parsedObject !== null && 'data' in parsedObject;
      const data = hasData ? (parsedObject.data as T) : (parsed as T);
      const success = parsedObject !== null && 'success' in parsedObject
        ? Boolean(parsedObject.success)
        : response.ok;
      const message = parsedObject !== null && typeof parsedObject.message === 'string'
        ? parsedObject.message
        : undefined;

      return {
        data,
        success,
        message,
      } satisfies ApiResponse<T>;
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // GET request
  async get<T>(endpoint: string, options?: RequestInit): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'GET',
      ...options,
    });
  }

  // POST request
  async post<T, P = unknown>(
    endpoint: string,
    data?: P,
    options?: RequestInit
  ): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data !== undefined ? JSON.stringify(data) : undefined,
      ...options,
    });
  }

  // PUT request
  async put<T, P = unknown>(
    endpoint: string,
    data?: P,
    options?: RequestInit
  ): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data !== undefined ? JSON.stringify(data) : undefined,
      ...options,
    });
  }

  // DELETE request
  async delete<T>(endpoint: string, options?: RequestInit): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'DELETE',
      ...options,
    });
  }
}

// Export singleton instance
export const apiClient = new ApiClient();

// Export the class for testing
export { ApiClient };
