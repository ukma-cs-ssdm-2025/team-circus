import React, { createContext, useContext, useReducer, useEffect, type ReactNode } from 'react';
import { authService } from '../services/auth';
import type { AuthContextType, AuthState, AuthUser, LoginRequest, RegisterRequest } from '../types/auth';

// Initial state
const initialState: AuthState = {
  user: null,
  isAuthenticated: false,
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

// Create context
const AuthContext = createContext<AuthContextType | undefined>(undefined);

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
        dispatch({ type: 'AUTH_START' });

        // No validateToken available — try to refresh tokens instead.
        // If refresh succeeds, assume user is authenticated.
        const refreshSuccess = await authService.refreshToken();

        if (refreshSuccess) {
          // Replace this placeholder with a real user fetch when you have one (e.g. authService.getProfile()).
          const user: AuthUser = {
            uuid: 'placeholder-uuid',
            login: 'user',
            email: 'user@example.com',
            createdAt: new Date().toISOString(),
          };
          dispatch({ type: 'AUTH_SUCCESS', payload: user });
        } else {
          dispatch({ type: 'AUTH_FAILURE' });
        }
      } catch (error) {
        console.error('Auth check failed:', error);
        dispatch({ type: 'AUTH_FAILURE' });
      }
    };

    checkAuthStatus();
  }, []);

  const login = async (credentials: LoginRequest): Promise<void> => {
    try {
      dispatch({ type: 'AUTH_START' });
      await authService.login(credentials);

      // After successful login, set user (placeholder until you fetch real profile).
      const user: AuthUser = {
        uuid: 'placeholder-uuid',
        login: credentials.login,
        email: 'user@example.com',
        createdAt: new Date().toISOString(),
      };

      dispatch({ type: 'AUTH_SUCCESS', payload: user });
    } catch (error) {
      console.error('Login failed:', error);
      dispatch({ type: 'AUTH_FAILURE' });
      throw error;
    }
  };

  const register = async (userData: RegisterRequest): Promise<void> => {
    try {
      dispatch({ type: 'AUTH_START' });
      const user = await authService.register(userData);
      dispatch({ type: 'AUTH_SUCCESS', payload: user });
    } catch (error) {
      console.error('Registration failed:', error);
      dispatch({ type: 'AUTH_FAILURE' });
      throw error;
    }
  };

  const logout = async (): Promise<void> => {
    try {
      await authService.logout();
      dispatch({ type: 'LOGOUT' });
    } catch (error) {
      console.error('Logout failed:', error);
      // Still dispatch logout even if the request fails
      dispatch({ type: 'LOGOUT' });
    }
  };

  // Exposed refreshToken: will attempt refresh and update auth state accordingly
  const refreshToken = async (): Promise<boolean> => {
    try {
      const success = await authService.refreshToken();
      if (success) {
        const user: AuthUser = {
          uuid: 'placeholder-uuid',
          login: state.user?.login ?? 'user',
          email: state.user?.email ?? 'user@example.com',
          createdAt: state.user?.createdAt ?? new Date().toISOString(),
        };
        dispatch({ type: 'AUTH_SUCCESS', payload: user });
        return true;
      } else {
        dispatch({ type: 'AUTH_FAILURE' });
        return false;
      }
    } catch (error) {
      console.error('Token refresh failed:', error);
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

// Custom hook to use auth context
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
