# Refresh Token Implementation

This document describes the refresh token implementation that has been added to the application.

## Overview

The refresh token system provides secure authentication without storing tokens in the database. It uses JWT tokens with different expiration times for access and refresh tokens.

## Backend Implementation

### Token Types

- **Access Token**: Short-lived (15 minutes) for API requests
- **Refresh Token**: Long-lived (7 days) for obtaining new access tokens

### Endpoints

- `POST /api/v1/auth/login` - Login and get both tokens
- `POST /api/v1/auth/refresh` - Refresh access token using refresh token
- `POST /api/v1/auth/logout` - Logout and clear tokens

### Token Structure

Both tokens are JWTs with the following claims:
- `sub`: User UUID
- `exp`: Expiration timestamp
- `iat`: Issued at timestamp
- `iss`: Token type ("access" or "refresh")

## Frontend Implementation

### AuthService

The `AuthService` class provides methods for:
- `login(credentials)` - Login and store refresh token
- `refreshAccessToken()` - Refresh access token
- `logout()` - Logout and clear tokens
- `isLoggedIn()` - Check if user has valid refresh token

### Enhanced API Client

The `EnhancedApiClient` automatically:
- Includes access token in requests (from cookies)
- Handles 401 responses by refreshing tokens
- Retries failed requests with new tokens

## Usage Examples

### Basic Login
```typescript
import { authService } from './services/auth';

const response = await authService.login({
  login: 'username',
  password: 'password'
});

console.log('Access token:', response.access_token);
console.log('Refresh token:', response.refresh_token);
```

### Making Authenticated Requests
```typescript
import { enhancedApiClient } from './services/enhanced-api';

// This will automatically handle token refresh if needed
const data = await enhancedApiClient.get('/api/v1/users/profile');
```

### Manual Token Refresh
```typescript
import { authService } from './services/auth';

try {
  const newTokens = await authService.refreshAccessToken();
  console.log('New access token:', newTokens.access_token);
} catch (error) {
  console.error('Refresh failed:', error);
  // Redirect to login page
}
```

## Security Features

1. **HttpOnly Cookies**: Access tokens are stored in httpOnly cookies for security
2. **Short-lived Access Tokens**: 15-minute expiration reduces exposure window
3. **Long-lived Refresh Tokens**: 7-day expiration balances security and user experience
4. **Token Type Validation**: Refresh endpoint only accepts refresh tokens
5. **Automatic Cleanup**: Invalid tokens are cleared automatically

## Configuration

Make sure to set the `SECRET_TOKEN` environment variable on the backend for JWT signing.

## Migration Notes

The existing login endpoint now returns both tokens in the response body, while still setting the access token as an httpOnly cookie. This maintains backward compatibility while enabling refresh token functionality.
