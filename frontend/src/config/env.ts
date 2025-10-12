// Environment configuration
export const ENV = {
  // API Configuration
  API_BASE_URL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  
  // Environment
  NODE_ENV: import.meta.env.MODE,
  APP_ENV: import.meta.env.VITE_APP_ENV || 'development',
  
  // App Configuration
  APP_NAME: import.meta.env.VITE_APP_NAME || 'MCD',
  APP_VERSION: import.meta.env.VITE_APP_VERSION || '1.0.0',
  
  // Feature flags
  ENABLE_DEV_TOOLS: import.meta.env.VITE_ENABLE_DEV_TOOLS === 'true',
  ENABLE_ANALYTICS: import.meta.env.VITE_ENABLE_ANALYTICS === 'true',
} as const;

// Type for environment variables
export type EnvConfig = typeof ENV;

// Helper function to get API URL
export const getApiUrl = (endpoint: string = '') => {
  const baseUrl = ENV.API_BASE_URL.endsWith('/') 
    ? ENV.API_BASE_URL.slice(0, -1) 
    : ENV.API_BASE_URL;
  
  const cleanEndpoint = endpoint.startsWith('/') 
    ? endpoint 
    : `/${endpoint}`;
    
  return `${baseUrl}${cleanEndpoint}`;
};

// Helper function to check if we're in development
export const isDevelopment = () => ENV.NODE_ENV === 'development';

// Helper function to check if we're in production
export const isProduction = () => ENV.NODE_ENV === 'production';
