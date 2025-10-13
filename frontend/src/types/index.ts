// User types
export interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
  role: 'admin' | 'user';
  createdAt: Date;
  updatedAt: Date;
}

// Document types
export interface Document {
  id: string;
  title: string;
  content: string;
  authorId: string;
  groupId?: string;
  isPublic: boolean;
  createdAt: Date;
  updatedAt: Date;
}

// Group types
export interface Group {
  id: string;
  name: string;
  description?: string;
  ownerId: string;
  memberIds: string[];
  createdAt: Date;
  updatedAt: Date;
}

// Navigation types
export interface NavItem {
  label: string;
  path: string;
  icon?: string;
  children?: NavItem[];
}

// API Response types
export interface ApiResponse<T = unknown> {
  data: T;
  message?: string;
  success: boolean;
}

export interface ApiError {
  message: string;
  status: number;
  code?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// Component props types
export interface BaseComponentProps {
  className?: string;
  children?: React.ReactNode;
}

// Form types
export interface FormField {
  name: string;
  label: string;
  type: 'text' | 'email' | 'password' | 'textarea' | 'select';
  required?: boolean;
  placeholder?: string;
  options?: { value: string; label: string }[];
}

// Theme types
export type Theme = 'light' | 'dark';

export interface ThemeContextType {
  theme: Theme;
  toggleTheme: () => void;
}
