// Note: We use direct fetch calls instead of apiClient to avoid circular dependencies
import { API_ENDPOINTS } from "../constants";
import type { AuthUser, LoginRequest, RegisterRequest } from "../types/auth";

class AuthService {
	private async requestWithCredentials<T>(
		endpoint: string,
		options: RequestInit = {},
	): Promise<T> {
		const response = await fetch(`${API_ENDPOINTS.BASE_URL}${endpoint}`, {
			...options,
			credentials: "include", // Important for cookies
			headers: {
				"Content-Type": "application/json",
				...options.headers,
			},
		});

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({}));
			throw new Error(
				errorData.error || `HTTP error! status: ${response.status}`,
			);
		}

		if (response.status === 204) {
			return undefined as T;
		}

		const contentType = response.headers.get("content-type");
		if (contentType?.includes("application/json")) {
			return response.json() as Promise<T>;
		}

		return undefined as T;
	}

	async login(credentials: LoginRequest): Promise<void> {
		await this.requestWithCredentials(API_ENDPOINTS.AUTH.LOGIN, {
			method: "POST",
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
			method: "POST",
			body: JSON.stringify(userData),
		});

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
				method: "POST",
			});
			return true;
		} catch (error) {
			console.error("Token refresh failed:", error);
			return false;
		}
	}

	async logout(): Promise<void> {
		await this.requestWithCredentials(API_ENDPOINTS.AUTH.LOGOUT, {
			method: "POST",
		});
	}
}

export const authService = new AuthService();
