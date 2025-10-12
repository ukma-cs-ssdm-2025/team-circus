# Environment Variables Setup

## Створення .env файлу

Створіть файл `.env` в кореневій папці frontend проекту:

```bash
# API Configuration
VITE_API_BASE_URL=http://localhost:8080/api

# Environment
VITE_APP_ENV=development

# App Configuration
VITE_APP_NAME=MCD
VITE_APP_VERSION=1.0.0

# Feature flags
VITE_ENABLE_DEV_TOOLS=true
VITE_ENABLE_ANALYTICS=false
```

## Важливі правила

### 1. Префікс VITE_
Всі змінні оточення для фронтенду **ОБОВ'ЯЗКОВО** повинні починатися з `VITE_`.

### 2. Типи змінних
- **String**: `VITE_API_BASE_URL=http://localhost:8080/api`
- **Boolean**: `VITE_ENABLE_DEV_TOOLS=true`
- **Number**: `VITE_APP_VERSION=1.0.0`

### 3. Доступ в коді
```typescript
import { ENV } from './config/env';

// Використання
const apiUrl = ENV.API_BASE_URL;
const isDev = ENV.APP_ENV === 'development';
```

## Різні середовища

### Development (.env.development)
```bash
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_ENV=development
VITE_ENABLE_DEV_TOOLS=true
```

### Production (.env.production)
```bash
VITE_API_BASE_URL=https://api.mcd.com/api
VITE_APP_ENV=production
VITE_ENABLE_DEV_TOOLS=false
```

### Staging (.env.staging)
```bash
VITE_API_BASE_URL=https://staging-api.mcd.com/api
VITE_APP_ENV=staging
VITE_ENABLE_DEV_TOOLS=true
```

## Використання в компонентах

### 1. Прямий доступ до змінних
```typescript
import { ENV } from '../config/env';

const MyComponent = () => {
  return <div>API URL: {ENV.API_BASE_URL}</div>;
};
```

### 2. Використання API клієнта
```typescript
import { apiClient } from '../services/api';

const MyComponent = () => {
  const fetchData = async () => {
    try {
      const response = await apiClient.get('/users');
      console.log(response.data);
    } catch (error) {
      console.error('API Error:', error);
    }
  };
};
```

### 3. Використання хуків
```typescript
import { useApi } from '../hooks/useApi';

const MyComponent = () => {
  const { data, loading, error } = useApi('/users');
  
  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;
  
  return <div>{JSON.stringify(data)}</div>;
};
```

## Безпека

### ✅ Безпечно (префікс VITE_)
- `VITE_API_BASE_URL`
- `VITE_APP_NAME`
- `VITE_ENABLE_FEATURES`

### ❌ НЕ безпечно (буде доступно в браузері)
- `API_SECRET_KEY`
- `DATABASE_PASSWORD`
- `JWT_SECRET`

## Перевірка змінних

Додайте компонент `ApiStatus` для перевірки підключення:

```typescript
import { ApiStatus } from '../components/common';

const Header = () => {
  return (
    <div>
      <h1>MCD</h1>
      <ApiStatus />
    </div>
  );
};
```
