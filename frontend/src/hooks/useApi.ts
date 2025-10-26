import { useState, useEffect, useCallback } from 'react';
import { apiClient } from '../services/apiClient';
import type { ApiResponse, ApiError } from '../types';

// Hook state interface
interface UseApiState<T> {
  data: T | null;
  loading: boolean;
  error: ApiError | null;
}

// Hook return type
interface UseApiReturn<T> extends UseApiState<T> {
  refetch: () => Promise<void>;
  mutate: (data: T) => void;
  reset: () => void;
}

// Generic API hook
export function useApi<T>(
  endpoint: string,
  options?: {
    immediate?: boolean;
    onSuccess?: (data: T) => void;
    onError?: (error: ApiError) => void;
  }
): UseApiReturn<T> {
  const [state, setState] = useState<UseApiState<T>>({
    data: null,
    loading: false,
    error: null,
  });

  const { immediate = true, onSuccess, onError } = options || {};

  const fetchData = useCallback(async () => {
    setState(prev => ({ ...prev, loading: true, error: null }));
    
    try {
      const response = await apiClient.get<T>(endpoint);
      setState({
        data: response.data,
        loading: false,
        error: null,
      });
      onSuccess?.(response.data);
    } catch (error) {
      const apiError: ApiError = {
        message: error instanceof Error ? error.message : 'Unknown error',
        status: 500,
      };
      setState({
        data: null,
        loading: false,
        error: apiError,
      });
      onError?.(apiError);
    }
  }, [endpoint, onSuccess, onError]);

  const mutate = useCallback((data: T) => {
    setState(prev => ({ ...prev, data }));
  }, []);

  const reset = useCallback(() => {
    setState({
      data: null,
      loading: false,
      error: null,
    });
  }, []);

  useEffect(() => {
    if (immediate) {
      fetchData();
    }
  }, [fetchData, immediate]);

  return {
    ...state,
    refetch: fetchData,
    mutate,
    reset,
  };
}

// Hook for mutations (POST, PUT, DELETE)
export function useMutation<T, P = unknown>(
  endpoint: string,
  method: 'POST' | 'PUT' | 'DELETE' = 'POST'
) {
  const [state, setState] = useState<UseApiState<T>>({
    data: null,
    loading: false,
    error: null,
  });

  const mutate = useCallback(async (payload?: P) => {
    setState(prev => ({ ...prev, loading: true, error: null }));
    
    try {
      let response: ApiResponse<T>;
      
      switch (method) {
        case 'POST':
          response = await apiClient.post<T>(endpoint, payload);
          break;
        case 'PUT':
          response = await apiClient.put<T>(endpoint, payload);
          break;
        case 'DELETE':
          response = await apiClient.delete<T>(endpoint);
          break;
        default:
          throw new Error(`Unsupported method: ${method}`);
      }

      setState({
        data: response.data,
        loading: false,
        error: null,
      });
      
      return response.data;
    } catch (error) {
      const apiError: ApiError = {
        message: error instanceof Error ? error.message : 'Unknown error',
        status: 500,
      };
      setState({
        data: null,
        loading: false,
        error: apiError,
      });
      throw apiError;
    }
  }, [endpoint, method]);

  const reset = useCallback(() => {
    setState({
      data: null,
      loading: false,
      error: null,
    });
  }, []);

  return {
    ...state,
    mutate,
    reset,
  };
}
