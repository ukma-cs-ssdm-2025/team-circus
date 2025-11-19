// API endpoints
export const API_ENDPOINTS = {
	BASE_URL: import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1",
	AUTH: {
		LOGIN: "/auth/login",
		REGISTER: "/signup",
		REFRESH: "/auth/refresh",
		LOGOUT: "/auth/logout",
	},
	USERS: {
		BASE: "/users",
		PROFILE: "/users/profile",
	},
	DOCUMENTS: {
		BASE: "/documents",
		SEARCH: "/documents/search",
		DETAIL: (uuid: string) => `/documents/${uuid}`,
	},
	GROUPS: {
		BASE: "/groups",
		DETAIL: (uuid: string) => `/groups/${uuid}`,
		MEMBERS: (groupUUID: string) => `/groups/${groupUUID}/members`,
		MEMBER: (groupUUID: string, userUUID: string) =>
			`/groups/${groupUUID}/members/${userUUID}`,
	},
} as const;

// Routes
export const ROUTES = {
	HOME: "/",
	LOGIN: "/login",
	REGISTER: "/register",
	PROFILE: "/profile",
	SETTINGS: "/settings",
	DOCUMENTS: "/documents",
	GROUPS: "/groups",
	GROUP_DETAILS: "/groups/:uuid",
	NOT_FOUND: "/404",
} as const;

// Local storage keys
export const STORAGE_KEYS = {
	TOKEN: "mcd_token",
	USER: "mcd_user",
	THEME: "mcd_theme",
	LANGUAGE: "mcd_language",
} as const;

// Theme
export const THEME = {
	LIGHT: "light",
	DARK: "dark",
} as const;

// User roles
export const USER_ROLES = {
	ADMIN: "admin",
	USER: "user",
} as const;

// Document permissions
export const DOCUMENT_PERMISSIONS = {
	READ: "read",
	WRITE: "write",
	DELETE: "delete",
	SHARE: "share",
} as const;

// Group member roles
export const MEMBER_ROLES = ["author", "editor", "viewer"] as const;

// Pagination
export const PAGINATION = {
	DEFAULT_PAGE_SIZE: 10,
	MAX_PAGE_SIZE: 100,
} as const;
