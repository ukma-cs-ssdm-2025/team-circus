package group_test

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
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group/requests"
	"go.uber.org/zap"
)

type mockCreateGroupService struct {
	mock.Mock
}

func (m *mockCreateGroupService) Create(ctx context.Context, name string) (*domain.Group, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Group), args.Error(1) //nolint:errcheck
}

func TestNewCreateGroupHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockCreateGroupService, gin.HandlerFunc) {
		mockService := &mockCreateGroupService{}
		handler := group.NewCreateGroupHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulCreate", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		expectedGroup := &domain.Group{
			UUID:      uuid.New(),
			Name:      "Test Group",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockService.On("Create", mock.Anything, "Test Group").Return(expectedGroup, nil)

		requestBody := requests.CreateGroupRequest{
			Name: "Test Group",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/groups", bytes.NewBuffer(jsonBody))
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
		assert.Equal(t, "Test Group", response["name"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		invalidJSON := `{"name": "Test Group"` // Missing closing brace

		req := httptest.NewRequest("POST", "/groups", bytes.NewBufferString(invalidJSON))
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

		mockService.AssertNotCalled(t, "Create")
	})

	t.Run("ValidationFailed_EmptyName", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		requestBody := requests.CreateGroupRequest{
			Name: "",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/groups", bytes.NewBuffer(jsonBody))
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

		mockService.AssertNotCalled(t, "Create")
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		mockService.On("Create", mock.Anything, "Test Group").Return(nil, domain.ErrInternal)

		requestBody := requests.CreateGroupRequest{
			Name: "Test Group",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/groups", bytes.NewBuffer(jsonBody))
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
		assert.Equal(t, "failed to create group", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		mockService.On("Create", mock.Anything, "Test Group").Return(nil, errors.New("database connection failed"))

		requestBody := requests.CreateGroupRequest{
			Name: "Test Group",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/groups", bytes.NewBuffer(jsonBody))
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
		assert.Equal(t, "failed to create group", response["error"])

		mockService.AssertExpectations(t)
	})
}

