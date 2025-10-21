package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth/requests"
)

func TestNewRefreshTokenHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set up environment variable for testing
	originalSecret := os.Getenv("SECRET_TOKEN")
	defer func() {
		if originalSecret == "" {
			os.Unsetenv("SECRET_TOKEN")
		} else {
			os.Setenv("SECRET_TOKEN", originalSecret)
		}
	}()
	os.Setenv("SECRET_TOKEN", "test-secret-key")

	t.Run("SuccessfulRefresh", func(t *testing.T) {
		// Arrange
		handler := auth.NewRefreshTokenHandler(nil) // No repo needed for refresh

		// Create a valid refresh token
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   "123e4567-e89b-12d3-a456-426614174000",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		})

		refreshTokenString, err := refreshToken.SignedString([]byte("test-secret-key"))
		assert.NoError(t, err)

		requestBody := requests.RefreshTokenRequest{
			RefreshToken: refreshTokenString,
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "access_token")
		assert.Contains(t, response, "refresh_token")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Arrange
		handler := auth.NewRefreshTokenHandler(nil)

		invalidJSON := `{"refresh_token": "some-token"` // Missing closing brace

		req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid request format", response["error"])
	})

	t.Run("InvalidToken", func(t *testing.T) {
		// Arrange
		handler := auth.NewRefreshTokenHandler(nil)

		requestBody := requests.RefreshTokenRequest{
			RefreshToken: "invalid-token",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid access token", response["error"])
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		// Arrange
		// Make sure SECRET_TOKEN is set for this test
		os.Setenv("SECRET_TOKEN", "test-secret-key")
		handler := auth.NewRefreshTokenHandler(nil)

		// Create an expired refresh token
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   "123e4567-e89b-12d3-a456-426614174000",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
		})

		refreshTokenString, err := refreshToken.SignedString([]byte("test-secret-key"))
		assert.NoError(t, err)

		requestBody := requests.RefreshTokenRequest{
			RefreshToken: refreshTokenString,
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// The JWT library might reject expired tokens during parsing, so it returns "invalid access token"
		assert.Equal(t, "invalid access token", response["error"])
	})

	t.Run("InvalidUUIDInToken", func(t *testing.T) {
		// Arrange
		handler := auth.NewRefreshTokenHandler(nil)

		// Create a token with invalid UUID
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   "invalid-uuid",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		})

		refreshTokenString, err := refreshToken.SignedString([]byte("test-secret-key"))
		assert.NoError(t, err)

		requestBody := requests.RefreshTokenRequest{
			RefreshToken: refreshTokenString,
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid uuid format in access token", response["error"])
	})

	t.Run("MissingSecretToken", func(t *testing.T) {
		// Arrange
		os.Unsetenv("SECRET_TOKEN")
		handler := auth.NewRefreshTokenHandler(nil)

		requestBody := requests.RefreshTokenRequest{
			RefreshToken: "some-token",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "server misconfiguration", response["error"])
	})

	t.Run("WrongSigningMethod", func(t *testing.T) {
		// Arrange
		// Make sure SECRET_TOKEN is set for this test
		os.Setenv("SECRET_TOKEN", "test-secret-key")
		handler := auth.NewRefreshTokenHandler(nil)

		// Create a token with wrong signing method (RSA instead of HMAC)
		// This will fail to sign with HMAC key, but let's test the parsing
		refreshTokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjNlNDU2Ny1lODliLTEyZDMtYTQ1Ni00MjY2MTQxNzQwMDAiLCJleHAiOjE3MDAwMDAwMDB9.invalid"

		requestBody := requests.RefreshTokenRequest{
			RefreshToken: refreshTokenString,
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// The handler should detect the wrong signing method and return "invalid access token"
		assert.Equal(t, "invalid access token", response["error"])
	})
}
