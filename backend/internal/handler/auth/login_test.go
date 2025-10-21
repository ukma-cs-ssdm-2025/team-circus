package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth/requests"
)

// MockUserRepository is a mock implementation of the user repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestNewLogInHandler(t *testing.T) {
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

	t.Run("SuccessfulLogin", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		handler := auth.NewLogInHandler(mockRepo)

		expectedUser := &domain.User{
			UUID:      uuid.New(),
			Login:     "testuser",
			Email:     "test@example.com",
			Password:  "testpassword123", // In real app, this would be hashed
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockRepo.On("GetByLogin", mock.Anything, "testuser").Return(expectedUser, nil)

		requestBody := requests.LogInRequest{
			Login:    "testuser",
			Password: "testpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		// Check that cookies are set
		cookies := w.Result().Cookies()
		assert.Len(t, cookies, 2)

		accessTokenCookie := findCookie(cookies, "accessToken")
		refreshTokenCookie := findCookie(cookies, "refreshToken")
		assert.NotNil(t, accessTokenCookie)
		assert.NotNil(t, refreshTokenCookie)

		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		handler := auth.NewLogInHandler(mockRepo)

		invalidJSON := `{"login": "testuser", "password": "testpassword123"` // Missing closing brace

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(invalidJSON))
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

		mockRepo.AssertNotCalled(t, "GetByLogin")
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		handler := auth.NewLogInHandler(mockRepo)

		// Return nil user (not found) without error
		mockRepo.On("GetByLogin", mock.Anything, "nonexistent").Return(nil, nil)

		requestBody := requests.LogInRequest{
			Login:    "nonexistent",
			Password: "testpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
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
		assert.Equal(t, "invalid credentials", response["error"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("WrongPassword", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		handler := auth.NewLogInHandler(mockRepo)

		expectedUser := &domain.User{
			UUID:      uuid.New(),
			Login:     "testuser",
			Email:     "test@example.com",
			Password:  "correctpassword",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockRepo.On("GetByLogin", mock.Anything, "testuser").Return(expectedUser, nil)

		requestBody := requests.LogInRequest{
			Login:    "testuser",
			Password: "wrongpassword",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
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
		assert.Equal(t, "invalid credentials", response["error"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		handler := auth.NewLogInHandler(mockRepo)

		mockRepo.On("GetByLogin", mock.Anything, "testuser").Return(nil, domain.ErrInternal)

		requestBody := requests.LogInRequest{
			Login:    "testuser",
			Password: "testpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
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
		assert.Equal(t, "failed to log in", response["error"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		handler := auth.NewLogInHandler(mockRepo)

		mockRepo.On("GetByLogin", mock.Anything, "testuser").Return(nil, errors.New("database connection failed"))

		requestBody := requests.LogInRequest{
			Login:    "testuser",
			Password: "testpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
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
		assert.Equal(t, "failed to log in", response["error"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("MissingSecretToken", func(t *testing.T) {
		// Arrange
		os.Unsetenv("SECRET_TOKEN")
		mockRepo := new(MockUserRepository)
		handler := auth.NewLogInHandler(mockRepo)

		expectedUser := &domain.User{
			UUID:      uuid.New(),
			Login:     "testuser",
			Email:     "test@example.com",
			Password:  "testpassword123",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockRepo.On("GetByLogin", mock.Anything, "testuser").Return(expectedUser, nil)

		requestBody := requests.LogInRequest{
			Login:    "testuser",
			Password: "testpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
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

		mockRepo.AssertExpectations(t)
	})
}

// Helper function to find a cookie by name
func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
