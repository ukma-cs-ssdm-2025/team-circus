package group_test

import (
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
	"go.uber.org/zap"
)

type mockGetGroupService struct {
	mock.Mock
}

func (m *mockGetGroupService) GetByUUIDForUser(ctx context.Context, groupUUID, userUUID uuid.UUID) (*domain.Group, error) {
	args := m.Called(ctx, groupUUID, userUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Group), args.Error(1) //nolint:errcheck
}

type mockGetAllGroupsService struct {
	mock.Mock
}

func (m *mockGetAllGroupsService) GetAllForUser(ctx context.Context, userUUID uuid.UUID) ([]*domain.Group, error) {
	args := m.Called(ctx, userUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Group), args.Error(1) //nolint:errcheck
}

func TestNewGetGroupHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockGetGroupService, gin.HandlerFunc) {
		mockService := &mockGetGroupService{}
		handler := group.NewGetGroupHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulGetGroup", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		groupUUID := uuid.New()
		userUUID := uuid.New()
		expectedGroup := &domain.Group{
			UUID:      groupUUID,
			Name:      "Test Group",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockService.On("GetByUUIDForUser", mock.Anything, groupUUID, userUUID).Return(expectedGroup, nil)

		req := httptest.NewRequest("GET", "/groups/"+groupUUID.String(), nil)
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
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Group", response["name"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()

		req := httptest.NewRequest("GET", "/groups/invalid-uuid", nil)
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
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid uuid format", response["error"])

		mockService.AssertNotCalled(t, "GetByUUIDForUser")
	})

	t.Run("MissingUserContext", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		groupUUID := uuid.New()

		req := httptest.NewRequest("GET", "/groups/"+groupUUID.String(), nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user context missing", response["error"])

		mockService.AssertNotCalled(t, "GetByUUIDForUser")
	})

	t.Run("InvalidUserContext", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		groupUUID := uuid.New()

		req := httptest.NewRequest("GET", "/groups/"+groupUUID.String(), nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
		c.Set("user_uid", "not-a-uuid") // Invalid type

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid user context", response["error"])

		mockService.AssertNotCalled(t, "GetByUUIDForUser")
	})

	t.Run("GroupNotFound", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		groupUUID := uuid.New()
		userUUID := uuid.New()
		mockService.On("GetByUUIDForUser", mock.Anything, groupUUID, userUUID).Return(nil, domain.ErrGroupNotFound)

		req := httptest.NewRequest("GET", "/groups/"+groupUUID.String(), nil)
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
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "group not found", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("Forbidden", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		groupUUID := uuid.New()
		userUUID := uuid.New()
		mockService.On("GetByUUIDForUser", mock.Anything, groupUUID, userUUID).Return(nil, domain.ErrForbidden)

		req := httptest.NewRequest("GET", "/groups/"+groupUUID.String(), nil)
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
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "access forbidden", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		groupUUID := uuid.New()
		userUUID := uuid.New()
		mockService.On("GetByUUIDForUser", mock.Anything, groupUUID, userUUID).Return(nil, domain.ErrInternal)

		req := httptest.NewRequest("GET", "/groups/"+groupUUID.String(), nil)
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
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to get group", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		groupUUID := uuid.New()
		userUUID := uuid.New()
		mockService.On("GetByUUIDForUser", mock.Anything, groupUUID, userUUID).Return(nil, errors.New("database connection failed"))

		req := httptest.NewRequest("GET", "/groups/"+groupUUID.String(), nil)
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
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to get group", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestNewGetAllGroupsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockGetAllGroupsService, gin.HandlerFunc) {
		mockService := &mockGetAllGroupsService{}
		handler := group.NewGetAllGroupsHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulGetAllGroups", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()

		group1 := &domain.Group{
			UUID:      uuid.New(),
			Name:      "Group 1",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		group2 := &domain.Group{
			UUID:      uuid.New(),
			Name:      "Group 2",
			CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		}

		expectedGroups := []*domain.Group{group1, group2}
		mockService.On("GetAllForUser", mock.Anything, userUUID).Return(expectedGroups, nil)

		req := httptest.NewRequest("GET", "/groups", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "groups")

		groups := response["groups"].([]interface{}) //nolint:errcheck
		assert.Len(t, groups, 2)
		assert.Equal(t, "Group 1", groups[0].(map[string]interface{})["name"]) //nolint:errcheck
		assert.Equal(t, "Group 2", groups[1].(map[string]interface{})["name"]) //nolint:errcheck

		mockService.AssertExpectations(t)
	})

	t.Run("EmptyGroupsList", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		expectedGroups := []*domain.Group{}
		mockService.On("GetAllForUser", mock.Anything, userUUID).Return(expectedGroups, nil)

		req := httptest.NewRequest("GET", "/groups", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "groups")

		groups := response["groups"].([]interface{}) //nolint:errcheck
		assert.Len(t, groups, 0)

		mockService.AssertExpectations(t)
	})

	t.Run("MissingUserContext", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		req := httptest.NewRequest("GET", "/groups", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user context missing", response["error"])

		mockService.AssertNotCalled(t, "GetAllForUser")
	})

	t.Run("InvalidUserContext", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		req := httptest.NewRequest("GET", "/groups", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", "not-a-uuid") // Invalid type

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid user context", response["error"])

		mockService.AssertNotCalled(t, "GetAllForUser")
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		mockService.On("GetAllForUser", mock.Anything, userUUID).Return(nil, domain.ErrInternal)

		req := httptest.NewRequest("GET", "/groups", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to get groups", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		mockService.On("GetAllForUser", mock.Anything, userUUID).Return(nil, errors.New("database connection failed"))

		req := httptest.NewRequest("GET", "/groups", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to get groups", response["error"])

		mockService.AssertExpectations(t)
	})
}

