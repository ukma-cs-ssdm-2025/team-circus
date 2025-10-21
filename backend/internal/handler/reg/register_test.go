
package reg_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/reg"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/reg/requests"
)

// MockRegService is a mock implementation of the registration service
type MockRegService struct {
	mock.Mock
}

func (m *MockRegService) Register(ctx context.Context, login string, email string, password string) (*domain.User, error) {
	args := m.Called(ctx, login, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestNewRegHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("SuccessfulRegistration", func(t *testing.T) {
		// Arrange
		mockService := new(MockRegService)
		handler := reg.NewRegHandler(mockService)

		expectedUser := &domain.User{
			UUID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Login:     "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockService.On("Register", mock.Anything, "testuser", "test@example.com", "testpassword123").
			Return(expectedUser, nil)

		requestBody := requests.RegRequest{
			Login:    "testuser",
			Email:    "test@example.com",
			Password: "testpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", response["login"])
		assert.Equal(t, "test@example.com", response["email"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Arrange
		mockService := new(MockRegService)
		handler := reg.NewRegHandler(mockService)

		invalidJSON := `{"login": "testuser", "email": "test@example.com", "password": "testpassword123"` // Missing closing brace

		req := httptest.NewRequest("POST", "/signup", bytes.NewBufferString(invalidJSON))
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

		mockService.AssertNotCalled(t, "Register")
	})

	t.Run("ValidationFailure", func(t *testing.T) {
		// Arrange
		mockService := new(MockRegService)
		handler := reg.NewRegHandler(mockService)

		requestBody := requests.RegRequest{
			Login:    "", // Empty login should fail validation
			Email:    "test@example.com",
			Password: "testpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "validation failed", response["error"])
		assert.Contains(t, response, "details")

		mockService.AssertNotCalled(t, "Register")
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService := new(MockRegService)
		handler := reg.NewRegHandler(mockService)

		mockService.On("Register", mock.Anything, "testuser", "test@example.com", "testpassword123").
			Return(nil, domain.ErrInternal)

		requestBody := requests.RegRequest{
			Login:    "testuser",
			Email:    "test@example.com",
			Password: "testpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
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
		assert.Equal(t, "failed to register", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService := new(MockRegService)
		handler := reg.NewRegHandler(mockService)

		mockService.On("Register", mock.Anything, "testuser", "test@example.com", "testpassword123").
			Return(nil, errors.New("database connection failed"))

		requestBody := requests.RegRequest{
			Login:    "testuser",
			Email:    "test@example.com",
			Password: "testpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
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
		assert.Equal(t, "failed to register", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("MissingContentType", func(t *testing.T) {
		// Arrange
		mockService := new(MockRegService)
		handler := reg.NewRegHandler(mockService)

		// Send invalid JSON without Content-Type
		req := httptest.NewRequest("POST", "/signup", bytes.NewBufferString("invalid json"))
		// Don't set Content-Type header
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

		mockService.AssertNotCalled(t, "Register")
	})
}
