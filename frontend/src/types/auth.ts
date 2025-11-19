// Authentication types
export interface AuthUser {
	uuid: string;
	login: string;
	email: string;
	createdAt: string;
}

export interface LoginRequest {
	login: string;
	password: string;
}

export interface RegisterRequest {
	login: string;
	email: string;
	password: string;
}

export interface AuthContextType {
	user: AuthUser | null;
	isAuthenticated: boolean;
	isLoading: boolean;
	login: (credentials: LoginRequest) => Promise<void>;
	register: (userData: RegisterRequest) => Promise<void>;
	logout: () => Promise<void>;
	refreshToken: () => Promise<boolean>;
}

export interface AuthState {
	user: AuthUser | null;
	isAuthenticated: boolean;
	isLoading: boolean;
}
