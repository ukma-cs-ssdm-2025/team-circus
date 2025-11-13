import type { ApiError } from '../types';
import { HttpError } from '../services/httpError';

const defaultStatus = 0;

export const normalizeApiError = (
  error: unknown,
  fallbackMessage = 'Unknown error',
): ApiError => {
  if (error instanceof HttpError) {
    return {
      message: error.message || fallbackMessage,
      status: error.status,
      code: error.code,
      details: error.details,
    };
  }

  if (error instanceof Error) {
    return {
      message: error.message || fallbackMessage,
      status: defaultStatus,
    };
  }

  return {
    message: fallbackMessage,
    status: defaultStatus,
  };
};
