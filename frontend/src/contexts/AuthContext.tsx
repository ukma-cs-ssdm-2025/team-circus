import React, { useReducer, useEffect, type ReactNode } from 'react';
import { authService } from '../services/auth';
import { STORAGE_KEYS } from '../constants';
import { AuthContext } from './AuthContextBase';
import type { AuthContextType, AuthState, AuthUser, LoginRequest, RegisterRequest } from '../types/auth';

const loadStoredUser = (): AuthUser | null => {
  if (typeof window === 'undefined') {
    return null;
  }

  try {
    const raw = localStorage.getItem(STORAGE_KEYS.USER);
    if (!raw) {
      return null;
    }
    return JSON.parse(raw) as AuthUser;
  } catch (error) {
    console.warn('Failed to parse stored user', error);
    return null;
  }
};

const persistUser = (user: AuthUser | null): void => {
  if (typeof window === 'undefined') {
    return;
  }

  try {
    if (user) {
      localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(user));
    } else {
      localStorage.removeItem(STORAGE_KEYS.USER);
    }
  } catch (error) {
    console.warn('Failed to persist user', error);
  }
};

// Initial state
const storedUser = loadStoredUser();

const initialState: AuthState = {
  user: storedUser,
  isAuthenticated: Boolean(storedUser),
  isLoading: true,
};

// Action types
type AuthAction =
  | { type: 'AUTH_START' }
  | { type: 'AUTH_SUCCESS'; payload: AuthUser }
  | { type: 'AUTH_FAILURE' }
  | { type: 'LOGOUT' }
  | { type: 'SET_LOADING'; payload: boolean };

// Reducer
const authReducer = (state: AuthState, action: AuthAction): AuthState => {
  switch (action.type) {
    case 'AUTH_START':
      return {
        ...state,
        isLoading: true,
      };
    case 'AUTH_SUCCESS':
      return {
        ...state,
        user: action.payload,
        isAuthenticated: true,
        isLoading: false,
      };
    case 'AUTH_FAILURE':
      return {
        ...state,
        user: null,
        isAuthenticated: false,
        isLoading: false,
      };
    case 'LOGOUT':
      return {
        ...state,
        user: null,
        isAuthenticated: false,
        isLoading: false,
      };
    case 'SET_LOADING':
      return {
        ...state,
        isLoading: action.payload,
      };
    default:
      return state;
  }
};

// Provider component
interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [state, dispatch] = useReducer(authReducer, initialState);

  // Check authentication status on mount
  useEffect(() => {
    const checkAuthStatus = async () => {
      try {
        const stored = loadStoredUser();

        if (stored) {
          dispatch({ type: 'AUTH_SUCCESS', payload: stored });
          return;
        }

        dispatch({ type: 'AUTH_START' });
        const refreshSuccess = await authService.refreshToken();

        if (refreshSuccess) {
          const user: AuthUser = {
            uuid: 'placeholder-uuid',
            login: 'user',
            email: 'user@example.com',
            createdAt: new Date().toISOString(),
          };

          persistUser(user);
          dispatch({ type: 'AUTH_SUCCESS', payload: user });
        } else {
          persistUser(null);
          dispatch({ type: 'AUTH_FAILURE' });
        }
      } catch (error) {
        console.error('Auth check failed:', error);
        persistUser(null);
        dispatch({ type: 'AUTH_FAILURE' });
      }
    };

    checkAuthStatus();
  }, []);

  const login = async (credentials: LoginRequest): Promise<void> => {
    try {
      dispatch({ type: 'AUTH_START' });
      await authService.login(credentials);

      const stored = loadStoredUser();
      const user: AuthUser =
        stored?.login === credentials.login
          ? stored
          : {
              uuid: 'placeholder-uuid',
              login: credentials.login,
              email: 'user@example.com',
              createdAt: new Date().toISOString(),
            };

      persistUser(user);
      dispatch({ type: 'AUTH_SUCCESS', payload: user });
    } catch (error) {
      console.error('Login failed:', error);
      persistUser(null);
      dispatch({ type: 'AUTH_FAILURE' });
      throw error;
    }
  };

  const register = async (userData: RegisterRequest): Promise<void> => {
    try {
      dispatch({ type: 'AUTH_START' });
      const user = await authService.register(userData);
      await authService.login({ login: userData.login, password: userData.password });
      persistUser(user);
      dispatch({ type: 'AUTH_SUCCESS', payload: user });
    } catch (error) {
      console.error('Registration failed:', error);
      persistUser(null);
      dispatch({ type: 'AUTH_FAILURE' });
      throw error;
    }
  };

  const logout = async (): Promise<void> => {
    try {
      await authService.logout();
      persistUser(null);
      dispatch({ type: 'LOGOUT' });
    } catch (error) {
      console.error('Logout failed:', error);
      // Still dispatch logout even if the request fails
      persistUser(null);
      dispatch({ type: 'LOGOUT' });
    }
  };

  // Exposed refreshToken: will attempt refresh and update auth state accordingly
  const refreshToken = async (): Promise<boolean> => {
    try {
      const success = await authService.refreshToken();
      if (success) {
        const stored = loadStoredUser();
        const existingUser = state.user ?? stored;
        const user: AuthUser =
          existingUser ?? {
            uuid: 'placeholder-uuid',
            login: 'user',
            email: 'user@example.com',
            createdAt: new Date().toISOString(),
          };

        persistUser(user);
        dispatch({ type: 'AUTH_SUCCESS', payload: user });
        return true;
      } else {
        persistUser(null);
        dispatch({ type: 'AUTH_FAILURE' });
        return false;
      }
    } catch (error) {
      console.error('Token refresh failed:', error);
      persistUser(null);
      dispatch({ type: 'AUTH_FAILURE' });
      return false;
    }
  };

  const value: AuthContextType = {
    user: state.user,
    isAuthenticated: state.isAuthenticated,
    isLoading: state.isLoading,
    login,
    register,
    logout,
    refreshToken,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
