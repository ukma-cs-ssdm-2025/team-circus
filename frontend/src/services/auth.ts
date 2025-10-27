// Note: We use direct fetch calls instead of apiClient to avoid circular dependencies
import { getApiUrl } from '../config/env';
import { HttpError, isRecord } from './httpError';
import { API_ENDPOINTS } from '../constants';
import type { LoginRequest, RegisterRequest, AuthUser } from '../types/auth';

class AuthService {
  private async requestWithCredentials<T>(
    endpoint: string,
    options: RequestInit = {},
  ): Promise<T | undefined> {
    try {
      const { headers: customHeaders, ...restOptions } = options;
      const hasJsonBody =
        restOptions.body !== undefined && restOptions.body !== null;

      const headers: HeadersInit = {
        Accept: 'application/json',
        ...(hasJsonBody ? { 'Content-Type': 'application/json' } : {}),
        ...(customHeaders ?? {}),
      };

      const response = await fetch(getApiUrl(endpoint), {
        ...restOptions,
        credentials: 'include', // Important for cookies
        headers,
      });

      if (!response.ok) {
        const { message, details, code } =
          await this.parseErrorResponse(response);
        throw new HttpError(message, response.status, details, code);
      }

      if (response.status === 204) {
        return undefined;
      }

      const contentType = response.headers.get('content-type');
      if (contentType?.includes('application/json')) {
        const data = (await response.json()) as T;
        return data;
      }

      return undefined;
    } catch (error) {
      if (error instanceof HttpError) {
        throw error;
      }
      if (error instanceof Error) {
        throw new HttpError(error.message, 0);
      }
      throw new HttpError('Unknown error', 0);
    }
  }

  async login(credentials: LoginRequest): Promise<void> {
    await this.requestWithCredentials(API_ENDPOINTS.AUTH.LOGIN, {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  }

  async register(userData: RegisterRequest): Promise<AuthUser> {
    const response = await this.requestWithCredentials<{
      uuid: string;
      login: string;
      email: string;
      created_at: string;
    }>(API_ENDPOINTS.AUTH.REGISTER, {
      method: 'POST',
      body: JSON.stringify(userData),
    });

    if (!response) {
      throw new HttpError('Empty response from registration endpoint', 0);
    }

    return {
      uuid: response.uuid,
      login: response.login,
      email: response.email,
      createdAt: response.created_at,
    };
  }

  async refreshToken(): Promise<boolean> {
    try {
      await this.requestWithCredentials(API_ENDPOINTS.AUTH.REFRESH, {
        method: 'POST',
      });
      return true;
    } catch (error) {
      console.error('Token refresh failed:', error);
      return false;
    }
  }

  async logout(): Promise<void> {
    await this.requestWithCredentials(API_ENDPOINTS.AUTH.LOGOUT, {
      method: 'POST',
    });
  }

  private async parseErrorResponse(
    response: Response,
  ): Promise<{ message: string; details?: unknown; code?: string }> {
    const fallbackMessage = `Request failed with status ${response.status}`;
    const contentType = response.headers.get('content-type') ?? '';

    if (contentType.includes('application/json')) {
      try {
        const json = await response.json();
        if (isRecord(json)) {
          const message =
            typeof json.error === 'string'
              ? json.error
              : typeof json.message === 'string'
                ? json.message
                : fallbackMessage;
          const code = typeof json.code === 'string' ? json.code : undefined;
          return { message, details: json, code };
        }
        return { message: fallbackMessage, details: json };
      } catch {
        return { message: fallbackMessage };
      }
    }

    try {
      const text = await response.text();
      const message = text.trim() || response.statusText || fallbackMessage;
      return { message, details: text };
    } catch {
      return { message: fallbackMessage };
    }
  }
}

export const authService = new AuthService();
