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

type mockUpdateGroupService struct {
	mock.Mock
}

func (m *mockUpdateGroupService) Update(ctx context.Context, userUUID, uuid uuid.UUID, name string) (*domain.Group, error) {
	args := m.Called(ctx, userUUID, uuid, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Group), args.Error(1) //nolint:errcheck
}

func TestNewUpdateGroupHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockUpdateGroupService, gin.HandlerFunc) {
		mockService := &mockUpdateGroupService{}
		handler := group.NewUpdateGroupHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulUpdate", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		groupUUID := uuid.New()
		expectedGroup := &domain.Group{
			UUID:      groupUUID,
			Name:      "Updated Group",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockService.On("Update", mock.Anything, userUUID, groupUUID, "Updated Group").Return(expectedGroup, nil)

		requestBody := requests.UpdateGroupRequest{
			Name: "Updated Group",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/groups/"+groupUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Group", response["name"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		requestBody := requests.UpdateGroupRequest{
			Name: "Updated Group",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/groups/invalid-uuid", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: "invalid-uuid"}}
		c.Set("user_uid", userUUID)

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
		groupUUID := uuid.New()
		invalidJSON := `{"name": "Updated Group"` // Missing closing brace

		req := httptest.NewRequest("PUT", "/groups/"+groupUUID.String(), bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
		c.Set("user_uid", userUUID)

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

	t.Run("MissingUserContext", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		groupUUID := uuid.New()
		requestBody := requests.UpdateGroupRequest{
			Name: "Updated Group",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/groups/"+groupUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("ValidationFailed_EmptyName", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		groupUUID := uuid.New()
		requestBody := requests.UpdateGroupRequest{
			Name: "",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/groups/"+groupUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
		c.Set("user_uid", userUUID)

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

	t.Run("GroupNotFound", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		groupUUID := uuid.New()
		mockService.On("Update", mock.Anything, userUUID, groupUUID, "Updated Group").Return(nil, domain.ErrGroupNotFound)

		requestBody := requests.UpdateGroupRequest{
			Name: "Updated Group",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/groups/"+groupUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "group not found", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		groupUUID := uuid.New()
		mockService.On("Update", mock.Anything, userUUID, groupUUID, "Updated Group").Return(nil, domain.ErrInternal)

		requestBody := requests.UpdateGroupRequest{
			Name: "Updated Group",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/groups/"+groupUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to update group", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		groupUUID := uuid.New()
		mockService.On("Update", mock.Anything, userUUID, groupUUID, "Updated Group").Return(nil, errors.New("database connection failed"))

		requestBody := requests.UpdateGroupRequest{
			Name: "Updated Group",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/groups/"+groupUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to update group", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("Forbidden", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		groupUUID := uuid.New()
		mockService.On("Update", mock.Anything, userUUID, groupUUID, "Updated Group").Return(nil, domain.ErrForbidden)

		requestBody := requests.UpdateGroupRequest{
			Name: "Updated Group",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/groups/"+groupUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusForbidden, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "access forbidden", response["error"])

		mockService.AssertExpectations(t)
	})
}
