import { getApiUrl } from "../config/env";
import { authService } from "./auth";
import type { ApiResponse } from "../types";

// Enhanced API Client with automatic token refresh
class ApiClient {
	private isRefreshing = false;
	private refreshPromise: Promise<boolean> | null = null;

	constructor() {
		// API_BASE_URL is used in getApiUrl function
	}

	private async requestWithAuth<T>(
		endpoint: string,
		options: RequestInit = {},
	): Promise<ApiResponse<T>> {
		const url = getApiUrl(endpoint);

		const defaultHeaders = {
			"Content-Type": "application/json",
		};

		const config: RequestInit = {
			...options,
			credentials: "include", // Important for cookies
			headers: {
				...defaultHeaders,
				...options.headers,
			},
		};

		try {
			const response = await fetch(url, config);

			// If token is expired, try to refresh it
			if (response.status === 401) {
				const refreshed = await this.handleTokenRefresh();
				if (refreshed) {
					// Retry the original request with refreshed token
					return this.requestWithAuth<T>(endpoint, options);
				} else {
					// Refresh failed, user needs to login again
					throw new Error("Authentication expired. Please login again.");
				}
			}

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({}));
				throw new Error(
					errorData.error || `HTTP error! status: ${response.status}`,
				);
			}

			let parsed: unknown;

			try {
				parsed = await response.json();
			} catch {
				parsed = null;
			}

			const isObject = (value: unknown): value is Record<string, unknown> => {
				return typeof value === "object" && value !== null;
			};

			const parsedObject = isObject(parsed) ? parsed : null;
			const hasData = parsedObject !== null && "data" in parsedObject;
			const data = hasData ? (parsedObject.data as T) : (parsed as T);
			const success =
				parsedObject !== null && "success" in parsedObject
					? Boolean(parsedObject.success)
					: response.ok;
			const message =
				parsedObject !== null && typeof parsedObject.message === "string"
					? parsedObject.message
					: undefined;

			return {
				data,
				success,
				message,
			} satisfies ApiResponse<T>;
		} catch (error) {
			console.error("API request failed:", error);
			throw error;
		}
	}

	private async handleTokenRefresh(): Promise<boolean> {
		// If already refreshing, wait for the existing refresh to complete
		if (this.isRefreshing && this.refreshPromise) {
			return this.refreshPromise;
		}

		// Start a new refresh process
		this.isRefreshing = true;
		this.refreshPromise = authService.refreshToken();

		try {
			const result = await this.refreshPromise;
			return result;
		} finally {
			this.isRefreshing = false;
			this.refreshPromise = null;
		}
	}

	// GET request
	async get<T>(
		endpoint: string,
		options?: RequestInit,
	): Promise<ApiResponse<T>> {
		return this.requestWithAuth<T>(endpoint, {
			method: "GET",
			...options,
		});
	}

	// POST request
	async post<T, P = unknown>(
		endpoint: string,
		data?: P,
		options?: RequestInit,
	): Promise<ApiResponse<T>> {
		return this.requestWithAuth<T>(endpoint, {
			method: "POST",
			body: data !== undefined ? JSON.stringify(data) : undefined,
			...options,
		});
	}

	// PUT request
	async put<T, P = unknown>(
		endpoint: string,
		data?: P,
		options?: RequestInit,
	): Promise<ApiResponse<T>> {
		return this.requestWithAuth<T>(endpoint, {
			method: "PUT",
			body: data !== undefined ? JSON.stringify(data) : undefined,
			...options,
		});
	}

	// DELETE request
	async delete<T>(
		endpoint: string,
		options?: RequestInit,
	): Promise<ApiResponse<T>> {
		return this.requestWithAuth<T>(endpoint, {
			method: "DELETE",
			...options,
		});
	}
}

// Export singleton instance
export const apiClient = new ApiClient();

// Export the class for testing
export { ApiClient };
