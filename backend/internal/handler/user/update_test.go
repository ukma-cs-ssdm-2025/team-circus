package user_test

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
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user/requests"
	"go.uber.org/zap"
)

type mockUpdateUserService struct {
	mock.Mock
}

func (m *mockUpdateUserService) Update(ctx context.Context, uuid uuid.UUID, login string, email string, password string) (*domain.User, error) {
	args := m.Called(ctx, uuid, login, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1) //nolint:errcheck
}

func TestNewUpdateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockUpdateUserService, gin.HandlerFunc) {
		mockService := &mockUpdateUserService{}
		handler := user.NewUpdateUserHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulUpdate", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		expectedUser := &domain.User{
			UUID:      userUUID,
			Login:     "updateduser",
			Email:     "updated@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockService.On("Update", mock.Anything, userUUID, "updateduser", "updated@example.com", "newpassword123").Return(expectedUser, nil)

		requestBody := requests.UpdateUserRequest{
			Login:    "updateduser",
			Email:    "updated@example.com",
			Password: "newpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/users/"+userUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: userUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "updateduser", response["login"])
		assert.Equal(t, "updated@example.com", response["email"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		requestBody := requests.UpdateUserRequest{
			Login:    "updateduser",
			Email:    "updated@example.com",
			Password: "newpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/users/invalid-uuid", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: "invalid-uuid"}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid uuid format", response["error"])

		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		invalidJSON := `{"login": "updateduser", "email": "updated@example.com"` // Missing closing brace

		req := httptest.NewRequest("PUT", "/users/"+userUUID.String(), bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: userUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid request format", response["error"])

		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("ValidationFailed_EmptyLogin", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		requestBody := requests.UpdateUserRequest{
			Login:    "",
			Email:    "updated@example.com",
			Password: "newpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/users/"+userUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: userUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "validation failed", response["error"])

		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("ValidationFailed_EmptyEmail", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		requestBody := requests.UpdateUserRequest{
			Login:    "updateduser",
			Email:    "",
			Password: "newpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/users/"+userUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: userUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "validation failed", response["error"])

		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("ValidationFailed_EmptyPassword", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		requestBody := requests.UpdateUserRequest{
			Login:    "updateduser",
			Email:    "updated@example.com",
			Password: "",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/users/"+userUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: userUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "validation failed", response["error"])

		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		mockService.On("Update", mock.Anything, userUUID, "updateduser", "updated@example.com", "newpassword123").Return(nil, domain.ErrUserNotFound)

		requestBody := requests.UpdateUserRequest{
			Login:    "updateduser",
			Email:    "updated@example.com",
			Password: "newpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/users/"+userUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: userUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user not found", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		mockService.On("Update", mock.Anything, userUUID, "updateduser", "updated@example.com", "newpassword123").Return(nil, domain.ErrInternal)

		requestBody := requests.UpdateUserRequest{
			Login:    "updateduser",
			Email:    "updated@example.com",
			Password: "newpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/users/"+userUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: userUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to update user", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		mockService.On("Update", mock.Anything, userUUID, "updateduser", "updated@example.com", "newpassword123").Return(nil, errors.New("database connection failed"))

		requestBody := requests.UpdateUserRequest{
			Login:    "updateduser",
			Email:    "updated@example.com",
			Password: "newpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/users/"+userUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: userUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to update user", response["error"])

		mockService.AssertExpectations(t)
	})
}

